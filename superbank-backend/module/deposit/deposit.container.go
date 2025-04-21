package deposit

import (
	"bank-backend/config"
	"bank-backend/middleware"
	"bank-backend/module/customer"

	"github.com/gofiber/fiber/v2"
)

type DepositContainer struct {
	Repository DepositRepository
	Service    DepositService
	Controller *DepositController
}

func NewDepositContainer() *DepositContainer {
	depositRepo := NewDepositRepository(config.DB)
	customerRepo := customer.NewCustomerRepository(config.DB)

	repo := NewDepositRepository(config.DB)
	service := NewDepositService(depositRepo, customerRepo)
	controller := NewDepositController(service, config.DB)

	return &DepositContainer{
		Repository: repo,
		Service:    service,
		Controller: controller,
	}
}

func SetupDepositRoutes(app fiber.Router) {
	container := NewDepositContainer()
	controller := container.Controller

	deposit := app.Group("/deposit", middleware.Authorize)

	deposit.Post("/create", controller.CreateDepositHandler)
}
