package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/cyberguru1/stage-two/config"
	"github.com/cyberguru1/stage-two/handlers"
	"github.com/cyberguru1/stage-two/middleware"
)

func SetupOrgRoutes(group fiber.Router, handlers *handlers.Handlers) {
	conf := config.New()
	useRoute := group.Group("/organisations")
	useRoute.Get("/",middleware.IsAuthorize(conf), handlers.GetOrgs)
	useRoute.Get("/:orgId",middleware.IsAuthorize(conf), handlers.GetOrg)
	useRoute.Post("/", middleware.IsAuthorize(conf), handlers.OrgRegister)
	useRoute.Post("/:orgId/users", handlers.UserOrgRegister)


}
