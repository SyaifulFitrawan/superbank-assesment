package auth

import (
	"bank-backend/module/auth"
	"bytes"
	"context"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Login(ctx context.Context, email, password string) (*auth.LoginResponse, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.LoginResponse), args.Error(1)
}

func setupApp(ctrl *auth.AuthController) *fiber.App {
	app := fiber.New()
	app.Post("/login", ctrl.LoginHandler)
	return app
}

func TestLoginHandler_Success(t *testing.T) {
	mockSvc := new(mockAuthService)
	ctrl := auth.NewAuthController(mockSvc)
	app := setupApp(ctrl)

	reqBody := `{"email": "test@example.com", "password": "secret"}`
	mockSvc.On("Login", mock.Anything, "test@example.com", "secret").
		Return(&auth.LoginResponse{Token: "mock-token"}, nil)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestLoginHandler_InvalidBody(t *testing.T) {
	mockSvc := new(mockAuthService)
	ctrl := auth.NewAuthController(mockSvc)
	app := setupApp(ctrl)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("invalid-json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestLoginHandler_LoginFailed(t *testing.T) {
	mockSvc := new(mockAuthService)
	ctrl := auth.NewAuthController(mockSvc)
	app := setupApp(ctrl)

	reqBody := `{"email": "wrong@example.com", "password": "wrongpass"}`
	mockSvc.On("Login", mock.Anything, "wrong@example.com", "wrongpass").
		Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}
