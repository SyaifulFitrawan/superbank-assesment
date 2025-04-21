package pocket

import (
	"bank-backend/model"
	"bank-backend/module/pocket"
	"bank-backend/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MockPocketService implements PocketService interface
type MockPocketService struct {
	mock.Mock
}

func (m *MockPocketService) Create(ctx context.Context, req pocket.CreatePocketRequest) (*model.Pocket, error) {
	args := m.Called(ctx, req)
	deposit, ok := args.Get(0).(*model.Pocket)
	if !ok && args.Get(0) != nil {
		panic("MockPocketService.Create: expected *model.Pocket return value")
	}
	return deposit, args.Error(1)
}

func (m *MockPocketService) TopUp(ctx context.Context, id string, req pocket.TopUpOrWithdrawPocketRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *MockPocketService) Withdrawn(ctx context.Context, id string, req pocket.TopUpOrWithdrawPocketRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *MockPocketService) Deactivated(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// SetupTestApp initializes a new fiber app and sets up the routes
func setupTestApp(t *testing.T) (*fiber.App, *gorm.DB, *MockPocketService) {
	t.Helper()

	app := fiber.New()
	mockService := new(MockPocketService)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	ctrl := pocket.NewPocketController(mockService, db)

	app.Post("/pocket", ctrl.CreatePocketHandler)
	app.Patch("/pocket/topup/:id", ctrl.TopUpPocketHandler)
	app.Patch("/pocket/withdraw/:id", ctrl.WithDrawnHandler)
	app.Patch("/pocket/deactivate/:id", ctrl.DeactivatedHandler)

	return app, db, mockService
}

func TestCreatePocketHandler(t *testing.T) {
	utils.InitValidator()
	app, _, mockService := setupTestApp(t)

	t.Run("should_create_pocket_successfully", func(t *testing.T) {
		utils.InitValidator()
		validUUID := uuid.NewV4().String()

		reqBody := pocket.CreatePocketRequest{
			CustomerID: validUUID,
			Name:       "My Pocket",
		}

		mockPocket := &model.Pocket{
			ID:         uuid.NewV4(),
			CustomerID: uuid.FromStringOrNil(validUUID),
			Name:       "My Pocket",
		}

		mockService.On("Create", mock.Anything, reqBody).Return(mockPocket, nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/pocket", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("should return error when invalid body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/pocket", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 when validation fails", func(t *testing.T) {
		utils.InitValidator()
		reqBody := pocket.CreatePocketRequest{
			CustomerID: "", // invalid UUID
			Name:       "", // empty name to fail validation
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/pocket", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 500 when service returns error", func(t *testing.T) {
		utils.InitValidator()
		validUUID := uuid.NewV4().String()

		reqBody := pocket.CreatePocketRequest{
			CustomerID: validUUID,
			Name:       "Pocket Error",
		}

		mockService.On("Create", mock.Anything, reqBody).Return(nil, errors.New("create error"))

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/pocket", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}

func TestTopUpPocketHandler(t *testing.T) {
	utils.InitValidator()
	app, _, mockService := setupTestApp(t)

	t.Run("should top up pocket successfully", func(t *testing.T) {
		reqBody := pocket.TopUpOrWithdrawPocketRequest{Amount: 5000}
		mockService.On("TopUp", mock.Anything, "pocket123", reqBody).Return(nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PATCH", "/pocket/topup/pocket123", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("should return error when top up failed", func(t *testing.T) {
		reqBody := pocket.TopUpOrWithdrawPocketRequest{Amount: 5000}
		mockService.On("TopUp", mock.Anything, "fail-pocket", reqBody).Return(errors.New("topup failed"))

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PATCH", "/pocket/topup/fail-pocket", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("should return 400 when top up body is invalid", func(t *testing.T) {
		req := httptest.NewRequest("PATCH", "/pocket/topup/any-id", bytes.NewReader([]byte("invalid-json")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 when top up body fails validation", func(t *testing.T) {
		reqBody := pocket.TopUpOrWithdrawPocketRequest{Amount: 0}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("PATCH", "/pocket/topup/any-id", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

func TestWithDrawnHandler(t *testing.T) {
	utils.InitValidator()
	app, _, mockService := setupTestApp(t)

	t.Run("should withdraw pocket successfully", func(t *testing.T) {
		reqBody := pocket.TopUpOrWithdrawPocketRequest{Amount: 3000}
		mockService.On("Withdrawn", mock.Anything, "pocket123", reqBody).Return(nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PATCH", "/pocket/withdraw/pocket123", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("should return error when withdraw failed", func(t *testing.T) {
		reqBody := pocket.TopUpOrWithdrawPocketRequest{Amount: 3000}
		mockService.On("Withdrawn", mock.Anything, "fail-pocket", reqBody).Return(errors.New("withdraw failed"))

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PATCH", "/pocket/withdraw/fail-pocket", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("should return 400 when withdraw body is invalid", func(t *testing.T) {
		req := httptest.NewRequest("PATCH", "/pocket/withdraw/any-id", bytes.NewReader([]byte("invalid-json")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 when withdraw body fails validation", func(t *testing.T) {
		reqBody := pocket.TopUpOrWithdrawPocketRequest{Amount: 0}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("PATCH", "/pocket/withdraw/any-id", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

func TestDeactivatedHandler(t *testing.T) {
	utils.InitValidator()
	app, _, mockService := setupTestApp(t)

	t.Run("should deactivate pocket successfully", func(t *testing.T) {
		mockService.On("Deactivated", mock.Anything, "pocket123").Return(nil)

		req := httptest.NewRequest("PATCH", "/pocket/deactivate/pocket123", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("should return error when deactivate failed", func(t *testing.T) {
		mockService.On("Deactivated", mock.Anything, "fail-pocket").Return(errors.New("deactivation failed"))

		req := httptest.NewRequest("PATCH", "/pocket/deactivate/fail-pocket", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}
