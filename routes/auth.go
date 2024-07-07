package routes


import (
	"github.com/gofiber/fiber/v2"

	"github.com/cyberguru1/stage-two/handlers"
)

func SetupAuthRoutes(group fiber.Router, handlers *handlers.Handlers) {
	group.Post("/register", handlers.UserRegister)
	group.Post("/login", handlers.UserLogin)
}
