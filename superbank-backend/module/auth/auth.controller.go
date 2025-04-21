package auth

import (
	"bank-backend/utils"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	service AuthService
}

func NewAuthController(service AuthService) *AuthController {
	return &AuthController{service: service}
}

func (ctrl *AuthController) LoginHandler(ctx *fiber.Ctx) error {
	var req LoginRequest
	logger := utils.NewLogger()

	if err := ctx.BodyParser(&req); err != nil {
		logger.Error("Invalid request: %v", err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseError(400, "Invalid request payload"))
	}

	token, err := ctrl.service.Login(context.Background(), req.Email, req.Password)
	if err != nil {
		logger.Error("Login failed: %v", err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.ResponseError(401, "Invalid email or password"))
	}

	token.Token = fmt.Sprintf("Bearer %s", token.Token)

	return ctx.Status(fiber.StatusOK).JSON(utils.ResponseSuccess(200, "Login Successful", token, nil))
}
