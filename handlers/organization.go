package handlers

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/cyberguru1/stage-two/ent"
	"github.com/cyberguru1/stage-two/ent/organisation"
	"github.com/cyberguru1/stage-two/ent/user"
	"github.com/cyberguru1/stage-two/middleware"
	"github.com/cyberguru1/stage-two/utils"
)

func (r orgRegisterReq) validate() ([]byte, error) {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(3, 250)),
	)

	if err == nil {
		return nil, nil
	}

	var FieldErrs []FieldErr

	if ve, ok := err.(validation.Errors); ok {
		for field, err := range ve {
			FieldErrs = append(FieldErrs, FieldErr{
				Field:   field,
				Message: err.Error(),
			})
		}
	}

	// marshal the map to JSON
	erroJSON, JsonErr := json.Marshal(FieldErrs)
	if JsonErr != nil {
		return nil, JsonErr
	}

	return erroJSON, nil
}

func (r orgAddUserReq) validate() ([]byte, error) {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.UserId, validation.Required, validation.Length(3, 250)),
	)

	if err == nil {
		return nil, nil
	}

	var FieldErrs []FieldErr

	if ve, ok := err.(validation.Errors); ok {
		for field, err := range ve {
			FieldErrs = append(FieldErrs, FieldErr{
				Field:   field,
				Message: err.Error(),
			})
		}
	}

	// marshal the map to JSON
	erroJSON, JsonErr := json.Marshal(FieldErrs)
	if JsonErr != nil {
		return nil, JsonErr
	}

	return erroJSON, nil
}

func (h *Handlers) GetOrgs(ctx *fiber.Ctx) error {

	userId, err := middleware.GetUserIdFromContext(ctx)

	if err != nil || userId == "" {
		_ = ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":     "Bad Request",
			"message":    "Client error",
			"statusCode": 400,
		})
		return nil
	}

	uid, _ := uuid.Parse(userId)

	u, err := h.Client.User.Query().WithOrganisations().Where(user.Userid(uid)).Only(ctx.Context())

	if err != nil {
		_ = ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "Bad Request",
			"message":    "Client error",
			"statusCode": 404,
		})

		return nil
	}

	var orgs []orgReq

	for _, v := range u.Edges.Organisations {
		orgs = append(orgs, orgReq{
			Name:        v.Name,
			OrgId:       v.Orgid.String(),
			Description: v.Description,
		})
	}

	outOrgs, _ := json.Marshal(orgs)
	_ = ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"data":    map[string]interface{}{"organisations": json.RawMessage(outOrgs)},
	})

	return nil
}

func (h *Handlers) GetOrg(ctx *fiber.Ctx) error {

	userId, err := middleware.GetUserIdFromContext(ctx)

	uid, _ := uuid.Parse(userId)

	if err != nil || userId == "" {
		_ = ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":     "Bad Request",
			"message":    "Client error",
			"statusCode": 400,
		})
		return nil
	}

	orgId, _ := uuid.Parse(ctx.Params("orgId"))

	outOrgs, err := h.Client.User.Query().
		WithOrganisations(func(q *ent.OrganisationQuery) {
			q.Where(organisation.Orgid(orgId))
		}).Where(user.Userid(uid)).
		Only(ctx.Context())

	org := outOrgs.Edges.Organisations

	if err != nil || len(org) == 0 {
		_ = ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "Bad Request",
			"message":    "Client error",
			"statusCode": 404,
		})

		return nil
	}
	o := org[0]
	data := map[string]interface{}{
		"orgId":       o.Orgid,
		"name":        o.Name,
		"description": o.Description,
	}

	_ = ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"data":    data,
	})

	return nil
}

func (h *Handlers) OrgRegister(ctx *fiber.Ctx) error {

	userId, err := middleware.GetUserIdFromContext(ctx)

	if err != nil || userId == "" {
		_ = ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":     "Bad Request",
			"message":    "Client error",
			"statusCode": 400,
		})
		return nil
	}

	uid, _ := uuid.Parse(userId)

	userQuery, err := h.Client.User.Query().Where(user.Userid(uid)).Only(ctx.Context())

	if err != nil {
		_ = ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "Bad Request",
			"message":    "Client error",
			"statusCode": 404,
		})

		return nil
	}

	var request orgRegisterReq

	if err := ctx.BodyParser(&request); err != nil {
		if err = ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":     "Bad request",
			"message":    "Client error",
			"statusCode": 400,
		}); err != nil {
			ctx.Status(http.StatusInternalServerError)
		}
		return nil
	}

	// validate users input and return json in expected format
	if errRes, err := request.validate(); errRes != nil || err != nil {
		ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": json.RawMessage(errRes),
		})
		return nil
	}

	if exist, _ := h.Client.Organisation.
		Query().
		Where(organisation.Name(request.Name)).
		Only(ctx.Context()); exist != nil {
		ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":     "Bad request",
			"message":    "Registration unsuccessful",
			"statusCode": 422,
		})
	}

	org, err := h.Client.Organisation.
		Create().
		SetName(request.Name).
		SetDescription(request.Description).
		Save(ctx.Context())

	if err != nil {
		ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":     "Bad request",
			"message":    "Registration unsuccessful",
			"statusCode": 422,
		})
		return nil
	}

	if _, err = userQuery.
		Update().
		AddOrganisations(org).
		Save(ctx.Context()); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":     "Bad request",
			"message":    "Client error",
			"statusCode": 400,
		})
		utils.Errorf("Fail to associate user with org: ", err)
		return nil
	}

	data := map[string]interface{}{
		"name":        org.Name,
		"description": org.Description,
	}

	_ = ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Organisation created successfully",
		"data":    data,
	})

	return nil
}

func (h *Handlers) UserOrgRegister(ctx *fiber.Ctx) error {

	var request orgAddUserReq

	if err := ctx.BodyParser(&request); err != nil {
		if err = ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":     "Bad request",
			"message":    "Client error",
			"statusCode": 400,
		}); err != nil {
			ctx.Status(http.StatusInternalServerError)
		}
		return nil
	}

	// validate users input and return json in expected format
	if errRes, err := request.validate(); errRes != nil || err != nil {
		ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": json.RawMessage(errRes),
		})
		return nil
	}

	uid, _ := uuid.Parse(request.UserId)
	orgId, _ := uuid.Parse(ctx.Params("orgId"))

	orgQuery, err := h.Client.Organisation.Query().Where(organisation.Orgid(orgId)).Only(ctx.Context())

	if err != nil {
		_ = ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "Bad Request",
			"message":    "Organization not found",
			"statusCode": 404,
		})

		return nil
	}

	userQuery, err := h.Client.User.Query().Where(user.Userid(uid)).Only(ctx.Context())

	_, err = userQuery.
		Update().
		AddOrganisations(orgQuery).
		Save(ctx.Context())

	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":     "Bad request",
			"message":    "Client error",
			"statusCode": 400,
		})
		utils.Errorf("Fail to associate user with org: ", err)
		return nil
	}

	ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User added to organisation successfully",
	})

	return nil
}
