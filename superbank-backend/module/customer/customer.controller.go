package customer

import (
	"bank-backend/model"
	"bank-backend/utils"
	"context"

	"github.com/gofiber/fiber/v2"
)

type CustomerController struct {
	service CustomerService
}

func NewCustomerController(service CustomerService) *CustomerController {
	return &CustomerController{service: service}
}

func (ctrl *CustomerController) CreateCustomerHandler(ctx *fiber.Ctx) error {
	var req CustomerCreateRequest
	logger := utils.NewLogger()

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseError(400, "Invalid request body"))
	}

	messages := utils.ValidateStruct(req)
	if len(messages) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseError(400, messages))
	}

	result, err := ctrl.service.Create(context.Background(), req)
	if err != nil {
		logger.Error("Create Customer Failed %v", err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError(500, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.ResponseSuccess(200, "Create Customer Successful", result, nil))
}

func (ctrl *CustomerController) ListCustomerHandler(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}

	limit := ctx.QueryInt("limit", 10)
	if limit < 1 {
		limit = 10
	}

	search := ctx.Query("search", "")

	customers, paginator, _ := ctrl.service.List(context.Background(), page, limit, search)

	return ctx.Status(fiber.StatusOK).JSON(
		utils.ResponseSuccess(fiber.StatusOK, "List Customers Successful", &customers, paginator),
	)
}

func (ctrl *CustomerController) DetailCustomerHandler(ctx *fiber.Ctx) error {
	logger := utils.NewLogger()
	id := ctx.Params("id")

	customer, err := ctrl.service.Detail(context.Background(), id)
	if err != nil {
		logger.Error("Get Customer Detail Failed: %v", err.Error())
		return ctx.Status(fiber.StatusNotFound).JSON(
			utils.ResponseError(fiber.StatusNotFound, "Customer not found"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		utils.ResponseSuccess(fiber.StatusOK, "Customer detail fetched successfully", customer, nil),
	)
}

func (ctrl *CustomerController) UpdateCustomerHandler(ctx *fiber.Ctx) error {
	logger := utils.NewLogger()
	id := ctx.Params("id")

	var req CustomerUpdateRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseError(400, "Invalid request body"))
	}

	if err := ctrl.service.Update(context.Background(), id, req); err != nil {
		logger.Error("Update Customer Failed: %v", err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			utils.ResponseError(fiber.StatusInternalServerError, "Failed to update customer"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		utils.ResponseSuccess[*model.Customer](fiber.StatusOK, "Customer updated successfully", nil, nil),
	)
}

func (ctrl *CustomerController) DeleteCustomerHandler(ctx *fiber.Ctx) error {
	logger := utils.NewLogger()
	id := ctx.Params("id")

	if err := ctrl.service.Delete(context.Background(), id); err != nil {
		logger.Error("Delete Customer Failed: %v", err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			utils.ResponseError(fiber.StatusInternalServerError, "Failed to delete customer"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		utils.ResponseSuccess[*model.Customer](fiber.StatusOK, "Customer deleted successfully", nil, nil),
	)
}
