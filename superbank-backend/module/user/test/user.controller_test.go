package user

import (
	"bank-backend/model"
	"bank-backend/module/user"
	"bank-backend/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

var HashPasswordFunc = bcrypt.GenerateFromPassword

func setupTestApp(t *testing.T) (*fiber.App, *MockUserService) {
	utils.InitValidator()
	app := fiber.New()
	mockService := new(MockUserService)
	controller := user.NewUserController(mockService)

	app.Post("/create", controller.CreateUserHandler)
	app.Get("/list", controller.ListUserHandler)
	app.Get("/detail/:id", controller.DetailUserHandler)
	app.Delete("/delete/:id", controller.DeleteUserHandler)

	return app, mockService
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(ctx context.Context, input *model.User) (*model.User, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) List(ctx context.Context, page, limit int, search string) ([]model.User, *utils.Paginator, error) {
	args := m.Called(ctx, page, limit, search)
	return args.Get(0).([]model.User), args.Get(1).(*utils.Paginator), args.Error(2)
}

func (m *MockUserService) Detail(ctx context.Context, id string) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateUserHandler(t *testing.T) {
	app, service := setupTestApp(t)

	t.Run("should return 200 and user data on success", func(t *testing.T) {
		reqBody := user.UserRequest{
			Email:    "test@example.com",
			Username: "test",
		}
		userModel := &model.User{
			ID:       uuid.NewV4(),
			Email:    "test@example.com",
			Username: "test",
		}
		service.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(userModel, nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should return 400 for invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader([]byte("invalid-json")))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 500 when service returns error", func(t *testing.T) {
		reqBody := user.UserRequest{
			Email:    "test@example.com",
			Username: "test",
		}
		service.ExpectedCalls = nil
		service.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil, errors.New("something went wrong")).Once()

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestListUserHandler(t *testing.T) {
	app, service := setupTestApp(t)

	t.Run("should return 200 and list of users", func(t *testing.T) {
		users := []model.User{
			{ID: uuid.NewV4(), Email: "test1@example.com", Username: "test1"},
			{ID: uuid.NewV4(), Email: "test2@example.com", Username: "test2"},
		}
		service.On("List", mock.Anything, 1, 10, "").Return(users, &utils.Paginator{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/list", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should return 200 with empty user list", func(t *testing.T) {
		service.ExpectedCalls = nil
		service.On("List", mock.Anything, 1, 10, "").Return([]model.User{}, &utils.Paginator{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/list?page=0&limit=0", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should handle service error gracefully", func(t *testing.T) {
		service.ExpectedCalls = nil
		service.On("List", mock.Anything, 1, 10, "").Return([]model.User{}, &utils.Paginator{}, errors.New("db error"))

		req := httptest.NewRequest(http.MethodGet, "/list", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode) // controller does not return error status
	})
}

func TestDetailUserHandler(t *testing.T) {
	app, service := setupTestApp(t)

	t.Run("should return 200 and user detail on success", func(t *testing.T) {
		userModel := &model.User{
			ID:       uuid.NewV4(),
			Email:    "test@example.com",
			Username: "test",
		}
		service.On("Detail", mock.Anything, mock.AnythingOfType("string")).Return(userModel, nil)

		req := httptest.NewRequest(http.MethodGet, "/detail/"+userModel.ID.String(), nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should return 404 when user not found", func(t *testing.T) {
		service.ExpectedCalls = nil
		service.On("Detail", mock.Anything, mock.AnythingOfType("string")).Return((*model.User)(nil), errors.New("not found")).Once()

		req := httptest.NewRequest(http.MethodGet, "/detail/some-non-existent-id", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestDeleteUserHandler(t *testing.T) {
	app, service := setupTestApp(t)

	t.Run("should return 200 on successful delete", func(t *testing.T) {
		userID := uuid.NewV4().String()
		service.On("Delete", mock.Anything, userID).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/delete/"+userID, nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should return 500 when delete fails", func(t *testing.T) {
		userID := uuid.NewV4().String()
		service.On("Delete", mock.Anything, userID).Return(errors.New("delete error"))

		req := httptest.NewRequest(http.MethodDelete, "/delete/"+userID, nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
