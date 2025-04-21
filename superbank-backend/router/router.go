package router

import (
	"bank-backend/module/auth"
	"bank-backend/module/customer"
	"bank-backend/module/dashboard"
	"bank-backend/module/deposit"
	"bank-backend/module/pocket"
	"bank-backend/module/user"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	route := app.Group("/api/v1")

	auth.SetupAuthRoutes(route)
	user.SetupUserRoutes(route)
	customer.SetupCustomerRoutes(route)
	deposit.SetupDepositRoutes(route)
	pocket.SetupPocketRoutes(route)
	dashboard.SetupDashboardRoutes(route)
}
