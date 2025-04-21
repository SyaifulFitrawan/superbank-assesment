package deposit

import (
	"bank-backend/database"
	"bank-backend/model"
	"bank-backend/utils"
	"context"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type DepositController struct {
	service DepositService
	db      *gorm.DB
}

func NewDepositController(service DepositService, db *gorm.DB) *DepositController {
	return &DepositController{service: service, db: db}
}

func (ctrl *DepositController) CreateDepositHandler(ctx *fiber.Ctx) error {
	var req CreateDepositRequest
	logger := utils.NewLogger()

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseError(400, "Invalid request body"))
	}

	message := utils.ValidateStruct(req)
	if len(message) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseError(400, message))
	}

	var result *model.Deposit

	err := database.WithTransaction(ctrl.db, func(tx *gorm.DB) error {
		ctxWithTx := database.NewContext(context.Background(), tx)

		res, err := ctrl.service.Create(ctxWithTx, req)
		if err != nil {
			return err
		}

		result = res
		return nil
	})

	if err != nil {
		logger.Error("Create deposit failed %v", err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError(500, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.ResponseSuccess(200, "Create deposit successfully", result, nil))
}
