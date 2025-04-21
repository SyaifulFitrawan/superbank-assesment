package middleware

import (
	"bank-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func Authorize(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ResponseError(401, "Unauthorize"))
	}

	return c.Next()
}
