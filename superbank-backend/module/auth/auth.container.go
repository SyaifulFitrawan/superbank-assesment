package auth

import (
	"bank-backend/config"
	"bank-backend/module/user"
	"bank-backend/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthContainer struct {
	Service    AuthService
	Controller *AuthController
}

func NewAuthContainer() *AuthContainer {
	utils := utils.NewJWTGenerator()

	repo := user.NewUserRepository(config.DB)
	service := NewAuthService(repo, utils)
	controller := NewAuthController(service)

	return &AuthContainer{
		Service:    service,
		Controller: controller,
	}
}

func SetupAuthRoutes(app fiber.Router) {
	container := NewAuthContainer()
	app.Post("/login", container.Controller.LoginHandler)
	app.Post("/register", container.Controller.LoginHandler)
}
