package dashboard

import (
	"bank-backend/config"
	"bank-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

type DashboardContainer struct {
	Service    DashboardService
	Controller *DashboardController
}

func NewDashboardContainer() *DashboardContainer {
	repo := NewDashboardRepository(config.DB)
	service := NewDashboardService(repo)
	controller := NewDashboardController(service)

	return &DashboardContainer{
		Service:    service,
		Controller: controller,
	}
}

func SetupDashboardRoutes(app fiber.Router) {
	container := NewDashboardContainer()
	app.Get("/dashboard", middleware.Authorize, container.Controller.GetDashboard)
}
