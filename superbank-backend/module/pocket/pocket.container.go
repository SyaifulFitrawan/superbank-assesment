package pocket

import (
	"bank-backend/config"
	"bank-backend/middleware"
	"bank-backend/module/customer"

	"github.com/gofiber/fiber/v2"
)

type PocketContainer struct {
	Repository PocketRepository
	Service    PocketService
	Controller *PocketController
}

func NewPocketContainer() *PocketContainer {
	pocketRepo := NewPocketRepository(config.DB)
	customerRepo := customer.NewCustomerRepository(config.DB)

	repo := NewPocketRepository(config.DB)
	service := NewPocketService(pocketRepo, customerRepo)
	controller := NewPocketController(service, config.DB)

	return &PocketContainer{
		Repository: repo,
		Service:    service,
		Controller: controller,
	}
}

func SetupPocketRoutes(app fiber.Router) {
	container := NewPocketContainer()
	controller := container.Controller

	pocket := app.Group("/pocket", middleware.Authorize)

	pocket.Post("/create", controller.CreatePocketHandler)
	pocket.Put("/topup/:id", controller.TopUpPocketHandler)
	pocket.Put("/withdrawn/:id", controller.WithDrawnHandler)
	pocket.Put("/deactive/:id", controller.DeactivatedHandler)
}
