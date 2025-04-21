package user

import (
	"bank-backend/model"
	"bank-backend/utils"
	"context"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	service UserService
}

func NewUserController(service UserService) *UserController {
	return &UserController{service: service}
}

func (ctrl *UserController) CreateUserHandler(ctx *fiber.Ctx) error {
	var req UserRequest
	logger := utils.NewLogger()

	if err := ctx.BodyParser(&req); err != nil {
		logger.Error("Invalid request payload: %v", err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseError(400, "Invalid request payload"))
	}

	user := &model.User{
		Email:    req.Email,
		Username: req.Username,
		Password: "password",
	}

	result, err := ctrl.service.Create(context.Background(), user)

	if err != nil {
		logger.Error("Create User Failed: %v", err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError(500, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.ResponseSuccess(200, "Create User Successfully", &result, nil))
}

func (ctrl *UserController) ListUserHandler(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}

	limit := ctx.QueryInt("limit", 10)
	if limit < 1 {
		limit = 10
	}

	search := ctx.Query("search", "")

	users, paginator, _ := ctrl.service.List(context.Background(), page, limit, search)

	return ctx.Status(fiber.StatusOK).JSON(
		utils.ResponseSuccess(fiber.StatusOK, "List Users Successful", &users, paginator),
	)
}

func (ctrl *UserController) DetailUserHandler(ctx *fiber.Ctx) error {
	logger := utils.NewLogger()
	id := ctx.Params("id")

	user, err := ctrl.service.Detail(context.Background(), id)
	if err != nil {
		logger.Error("Get User Detail Failed: %v", err.Error())
		return ctx.Status(fiber.StatusNotFound).JSON(
			utils.ResponseError(fiber.StatusNotFound, "User not found"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		utils.ResponseSuccess(fiber.StatusOK, "User detail fetched successfully", user, nil),
	)
}

func (ctrl *UserController) DeleteUserHandler(ctx *fiber.Ctx) error {
	logger := utils.NewLogger()
	id := ctx.Params("id")

	if err := ctrl.service.Delete(context.Background(), id); err != nil {
		logger.Error("Delete User Failed: %v", err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			utils.ResponseError(fiber.StatusInternalServerError, "Failed to delete customer"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		utils.ResponseSuccess[*model.User](fiber.StatusOK, "User deleted successfully", nil, nil),
	)
}
