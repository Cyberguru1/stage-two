package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/cyberguru1/stage-two/handlers"
)

func SetupApiV1(app *fiber.App, handlers *handlers.Handlers) {
	apiRoute := app.Group("/api")
	authRoute := app.Group("/auth")
	SetupUserRoutes(apiRoute, handlers)
	SetupAuthRoutes(authRoute, handlers)
	SetupOrgRoutes(apiRoute, handlers)

}
