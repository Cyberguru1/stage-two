package handlers

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gofiber/fiber/v2"

	"github.com/cyberguru1/stage-two/ent/user"
	"github.com/cyberguru1/stage-two/middleware"
	"github.com/cyberguru1/stage-two/utils"
)

type FieldErr struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (r registerReq) validate() ([]byte, error) {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.Firstname, validation.Required, validation.Length(3, 30)),
		validation.Field(&r.Lastname, validation.Required, validation.Length(3, 30)),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required, validation.Length(8, 30)),
		validation.Field(&r.Phone, validation.Length(0, 30)),
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

func (h *Handlers) UserRegister(ctx *fiber.Ctx) error {
	var request registerReq

	if err := ctx.BodyParser(&request); err != nil {
		if err = ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":     "Bad request",
			"message":    "Registration unsuccessful",
			"statusCode": 400,
		}); err != nil {
			ctx.Status(http.StatusInternalServerError)
		}
	}

	// validate users input and return json in expected format
	if errRes, err := request.validate(); errRes != nil || err != nil {
		ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": json.RawMessage(errRes),
		})
		return nil
	}

	if exist, _ := h.Client.User.Query().Where(user.Email(request.Email)).Only(ctx.Context()); exist != nil {
		ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":     "Bad request",
			"message":    "Registration unsuccessful",
			"statusCode": 422,
		})
	}

	hashpassword, err := utils.HashPassword(request.Password)

	if err != nil {
		utils.Errorf("Failed hash user password: ", err)
		return nil
	}

	if _, err := h.Client.User.Create().
		SetEmail(request.Email).
		SetFirstName(request.Firstname).
		SetLastName(request.Lastname).
		SetPassword(hashpassword).
		SetPhone(request.Phone).
		Save(ctx.Context()); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":     "Bad request",
			"message":    "Registration unsuccessful",
			"statusCode": 400,
		})
		utils.Errorf("Fail to create user: ", err)
		return nil
	}

	// Create a default organisation for the user
	organisation, err := h.Client.Organisation.
		Create().
		SetName(request.Firstname + "'s Organisation").
		Save(ctx.Context())
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":     "Bad request",
			"message":    "Registration unsuccessful",
			"statusCode": 400,
		})
		utils.Errorf("Fail to create user: ", err)
		return nil
	}

	// Associate the user with the organisation
	if _, err = h.Client.User.
		Update().
		AddOrganisations(organisation).
		Save(ctx.Context()); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":     "Bad request",
			"message":    "Registration unsuccessful",
			"statusCode": 400,
		})
		utils.Errorf("Fail to create user: ", err)
		return nil
	}

	u, err := h.Client.User.Query().Where(user.Email(request.Email)).Only(ctx.Context())

	if err != nil {
		_ = ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "Bad request",
			"message":    "Registration unsuccessful",
			"statusCode": 400,
		})

		return nil
	}

	token, err := middleware.ClaimToken(u.Userid)

	if err != nil {
		utils.Errorf("Token generation error: ", err)
		return nil
	}

	user := map[string]interface{}{
		"userId":    u.Userid,
		"firstName": u.FirstName,
		"lastName":  u.LastName,
		"email":     u.Email,
		"phone":     u.Phone,
	}

	data := map[string]interface{}{
		"accessToken": token,
		"user":        user,
	}

	_ = ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Registration successful",
		"data":    data,
	})

	return nil
}

func (h *Handlers) UserLogin(ctx *fiber.Ctx) error {
	var request loginReq

	err := ctx.BodyParser(&request)

	if err != nil {
		err = ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"status":     "Bad request",
			"message":    "Authentication failed",
			"statusCode": 401,
		})

		if err != nil {
			ctx.Status(http.StatusInternalServerError)
		}
		return nil
	}

	u, err := h.Client.User.Query().Where(user.Email(request.Email)).Only(ctx.Context())

	if err != nil {
		_ = ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":     "Bad request",
			"message":    "Authentication failed",
			"statusCode": 401,
		})

		return nil
	}

	if err = utils.ComparePassword(request.Password, u.Password); err != nil {
		_ = ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":     "Bad request",
			"message":    "Authentication failed",
			"statusCode": 401,
		})
	}

	token, err := middleware.ClaimToken(u.Userid)

	if err != nil {
		utils.Errorf("Token generation error: ", err)
		return nil
	}

	user := map[string]interface{}{
		"userId":    u.Userid,
		"firstName": u.FirstName,
		"lastName":  u.LastName,
		"email":     u.Email,
		"phone":     u.Phone,
	}

	data := map[string]interface{}{
		"accessToken": token,
		"user":        user,
	}

	_ = ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Login successful",
		"data":    data,
	})

	return nil
}
