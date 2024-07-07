package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/cyberguru1/stage-two/handlers"
)

func SetupOrgRoutes(group fiber.Router, handlers *handlers.Handlers) {
	useRoute := group.Group("/organisaitons")
	useRoute.Get("/", handlers.GetOrgs)
	useRoute.Get("/:orgId", handlers.GetOrg)
	useRoute.Post("/", handlers.OrgRegister)
	useRoute.Post("/:orgId/users", handlers.UserOrgRegister)


}
