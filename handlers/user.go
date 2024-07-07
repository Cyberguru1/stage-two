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

	if err != nil || passedId != userId {
		_ = ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":     "Bad Request",
			"message":    "Client error",
			"statusCode": 400,
		})
		return nil
	}

	uid, _ := uuid.Parse(userId)

	u, err := h.Client.User.Query().Where(user.Userid(uid)).Only(ctx.Context())

	if err != nil {
		_ = ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "Bad Request",
			"message":    "Client error",
			"statusCode": 400,
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
