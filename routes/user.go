package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/cyberguru1/stage-two/config"
	"github.com/cyberguru1/stage-two/handlers"
	"github.com/cyberguru1/stage-two/middleware"
)

func SetupUserRoutes(group fiber.Router, handlers *handlers.Handlers) {
	conf := config.New()
	useRoute := group.Group("/users")
	// useRoute.Use(middleware.IsAuthorize(conf)) // group way
	useRoute.Get("/:id", middleware.IsAuthorize(conf), handlers.GetUser) // specific
}
