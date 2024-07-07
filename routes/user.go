package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/cyberguru1/stage-two/handlers"
)

func SetupUserRoutes(group fiber.Router, handlers *handlers.Handlers) {
	useRoute := group.Group("/users")
	useRoute.Get("/:id", handlers.GetUser)
}
