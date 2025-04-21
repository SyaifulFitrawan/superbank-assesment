package dashboard

import (
	"bank-backend/utils"
	"context"

	"github.com/gofiber/fiber/v2"
)

type DashboardController struct {
	service DashboardService
}

func NewDashboardController(service DashboardService) *DashboardController {
	return &DashboardController{service: service}
}

func (ctrl *DashboardController) GetDashboard(ctx *fiber.Ctx) error {
	logger := utils.NewLogger()

	dashboard, err := ctrl.service.GetDashboard(context.Background())
	if err != nil {
		logger.Error("Get dashboard failed: %v", err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			utils.ResponseError(fiber.StatusInternalServerError, "Get dashboard failed"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.ResponseSuccess(fiber.StatusOK, "Get dashboard successfully", &dashboard, nil))
}
