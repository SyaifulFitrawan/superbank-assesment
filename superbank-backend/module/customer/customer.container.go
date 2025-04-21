package customer

import (
	"bank-backend/config"
	"bank-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

type CustomerContainer struct {
	Repository CustomerRepository
	Service    CustomerService
	Controller *CustomerController
}

func NewCustomerContainer() *CustomerContainer {
	repo := NewCustomerRepository(config.DB)
	service := NewCustomerService(repo)
	controller := NewCustomerController(service)

	return &CustomerContainer{
		Repository: repo,
		Service:    service,
		Controller: controller,
	}
}

func SetupCustomerRoutes(app fiber.Router) {
	container := NewCustomerContainer()
	controller := container.Controller

	customer := app.Group("/customer", middleware.Authorize)

	customer.Post("/create", controller.CreateCustomerHandler)
	customer.Get("/list", controller.ListCustomerHandler)
	customer.Get("/detail/:id", controller.DetailCustomerHandler)
	customer.Put("/update/:id", controller.UpdateCustomerHandler)
	customer.Delete("/delete/:id", controller.DeleteCustomerHandler)
}
