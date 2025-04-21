package user

import (
	"bank-backend/config"
	"bank-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

type UserContainer struct {
	Repository UserRepository
	Service    UserService
	Controller *UserController
}

func NewUserContainer() *UserContainer {
	repo := NewUserRepository(config.DB)
	service := NewUserService(repo)
	controller := NewUserController(service)

	return &UserContainer{
		Repository: repo,
		Service:    service,
		Controller: controller,
	}
}

func SetupUserRoutes(app fiber.Router) {
	container := NewUserContainer()
	controller := container.Controller

	user := app.Group("user", middleware.Authorize)

	user.Post("/create", controller.CreateUserHandler)
	user.Get("/list", controller.ListUserHandler)
	user.Get("/detail/:id", controller.DetailUserHandler)
	user.Delete("/delete/:id", controller.DeleteUserHandler)
}
