package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/cyberguru1/stage-two/ent/user"
	"github.com/cyberguru1/stage-two/middleware"
)

func (h *Handlers) GetUser(ctx *fiber.Ctx) error {

	passedId := ctx.Params("id")

	userId, err := middleware.GetUserIdFromContext(ctx)

	uid, _ := uuid.Parse(userId)
	pid, _ := uuid.Parse(passedId)

	u, err := h.Client.User.Query().Where(user.Userid(uid)).Only(ctx.Context())

	if err != nil || passedId != userId {

		orgs, err := u.QueryOrganisations().
			All(ctx.Context())

		isMember := false
		for _, org := range orgs {
			exists, _ := org.QueryUsers().Where(user.Userid(pid)).Only(ctx.Context())
			if exists != nil {
				isMember = true
				break
			}
		}
		if err != nil || isMember == false {
			_ = ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":     "Bad Request",
				"message":    "Client error",
				"statusCode": 401,
			})

			return nil
		}
		u, _ = h.Client.User.Query().Where(user.Userid(pid)).Only(ctx.Context())
	}

	// u, err := h.Client.User.Query().Where(user.Userid(uid)).Only(ctx.Context())

	if err != nil {
		_ = ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "Bad Request",
			"message":    "Client error",
			"statusCode": 404,
		})

		return nil
	}

	user := map[string]interface{}{
		"userId":    u.Userid,
		"firstName": u.FirstName,
		"lastName":  u.LastName,
		"email":     u.Email,
		"phone":     u.Phone,
	}

	_ = ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"data":    user,
	})

	return nil
}
