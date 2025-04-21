package pocket

import (
	"bank-backend/database"
	"bank-backend/model"
	"bank-backend/utils"
	"context"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PocketController struct {
	service PocketService
	db      *gorm.DB
}

func NewPocketController(service PocketService, db *gorm.DB) *PocketController {
	return &PocketController{service: service, db: db}
}

func (ctrl *PocketController) CreatePocketHandler(ctx *fiber.Ctx) error {
	var req CreatePocketRequest
	logger := utils.NewLogger()

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseError(400, "Invalid request body"))
	}

	message := utils.ValidateStruct(req)
	if len(message) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseError(400, message))
	}

	result, err := ctrl.service.Create(context.Background(), req)
	if err != nil {
		logger.Error("Create Pocket Failed %v", err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError(500, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.ResponseSuccess(200, "Create Pocket Successful", result, nil))
}

func (ctrl *PocketController) TopUpPocketHandler(ctx *fiber.Ctx) error {
	var req TopUpOrWithdrawPocketRequest
	id := ctx.Params("id")
	logger := utils.NewLogger()

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseError(400, "Invalid request body"))
	}

	message := utils.ValidateStruct(req)
	if len(message) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseError(400, message))
	}

	err := database.WithTransaction(ctrl.db, func(tx *gorm.DB) error {
		ctxWithTx := database.NewContext(context.Background(), tx)

		err := ctrl.service.TopUp(ctxWithTx, id, req)
		if err != nil {
			logger.Error("Topup pocket failed %v", err.Error())
			return err
		}

		return nil
	})

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError(500, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(
		utils.ResponseSuccess[*model.Pocket](fiber.StatusOK, "Topup pocket successfully", nil, nil),
	)
}

func (ctrl *PocketController) WithDrawnHandler(ctx *fiber.Ctx) error {
	var req TopUpOrWithdrawPocketRequest
	id := ctx.Params("id")
	logger := utils.NewLogger()

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseError(400, "Invalid request body"))
	}

	message := utils.ValidateStruct(req)
	if len(message) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseError(400, message))
	}

	err := database.WithTransaction(ctrl.db, func(tx *gorm.DB) error {
		ctxWithTx := database.NewContext(context.Background(), tx)

		err := ctrl.service.Withdrawn(ctxWithTx, id, req)
		if err != nil {
			logger.Error("Withdrawn pocket failed %v", err.Error())
			return err
		}

		return nil
	})

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError(500, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(
		utils.ResponseSuccess[*model.Pocket](fiber.StatusOK, "Withdrawn pocket successfully", nil, nil),
	)
}

func (ctrl *PocketController) DeactivatedHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	logger := utils.NewLogger()

	err := database.WithTransaction(ctrl.db, func(tx *gorm.DB) error {
		ctxWithTx := database.NewContext(context.Background(), tx)

		err := ctrl.service.Deactivated(ctxWithTx, id)
		if err != nil {
			logger.Error("Deactivated pocket failed %v", err.Error())
			return err
		}

		return nil
	})

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError(500, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(
		utils.ResponseSuccess[*model.Pocket](fiber.StatusOK, "Deactivated pocket successfully", nil, nil),
	)
}
