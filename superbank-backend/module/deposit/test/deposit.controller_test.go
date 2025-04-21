package deposit

import (
	"bank-backend/model"
	"bank-backend/module/deposit"
	"bank-backend/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestApp(t *testing.T) *fiber.App {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Deposit{})
	require.NoError(t, err)

	mockService := &MockDepositService{}

	controller := deposit.NewDepositController(mockService, db)

	app := fiber.New()
	app.Post("/deposit", controller.CreateDepositHandler)

	return app
}

type MockDepositService struct{}

func (m *MockDepositService) Create(ctx context.Context, req deposit.CreateDepositRequest) (*model.Deposit, error) {
	return &model.Deposit{
		ID:         uuid.NewV4(),
		CustomerID: uuid.FromStringOrNil(req.CustomerID),
		Amount:     req.Amount,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

func (m *MockDepositService) ProcessMatureDeposits(ctx context.Context) error {
	return nil
}

type MockDepositServiceWithError struct{}

func (m *MockDepositServiceWithError) Create(ctx context.Context, req deposit.CreateDepositRequest) (*model.Deposit, error) {
	return nil, fmt.Errorf("failed to create deposit")
}

func (m *MockDepositServiceWithError) ProcessMatureDeposits(ctx context.Context) error {
	return nil
}

func TestCreateDepositHandler(t *testing.T) {
	utils.InitValidator()

	t.Run("should return 200 and deposit data on success", func(t *testing.T) {
		app := setupTestApp(t)

		validUUID := uuid.NewV4()
		reqBody := fmt.Sprintf(`{
			"customer_id": "%s",
			"amount": 1000,
			"interest_rate": 5,
			"term_months": 12,
			"start_date": "2023-04-01",
			"note": "Test deposit"
		}`, validUUID.String())

		req := httptest.NewRequest(http.MethodPost, "/deposit", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		meta := body["meta"].(map[string]interface{})
		assert.Equal(t, true, meta["success"])
		assert.Equal(t, "Create deposit successfully", meta["message"])

		data := body["data"].(map[string]interface{})
		assert.NotEmpty(t, data["id"])
		assert.Equal(t, float64(1000), data["amount"])
	})

	t.Run("should return 400 when request body is invalid JSON", func(t *testing.T) {
		app := setupTestApp(t)

		req := httptest.NewRequest(http.MethodPost, "/deposit", bytes.NewBufferString(`invalid-json`))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var body map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&body)

		meta := body["meta"].(map[string]interface{})
		assert.Equal(t, false, meta["success"])
		assert.Equal(t, "Invalid request body", meta["message"])
	})

	t.Run("should return 400 when request fails validation", func(t *testing.T) {
		app := setupTestApp(t)

		reqBody := `{
			"customer_id": "",
			"amount": -1000,
			"interest_rate": 5,
			"term_months": 0,
			"start_date": "",
			"note": ""
		}`

		req := httptest.NewRequest(http.MethodPost, "/deposit", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var body map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&body)

		meta := body["meta"].(map[string]interface{})
		assert.Equal(t, false, meta["success"])
		assert.NotEmpty(t, meta["message"])
	})

	t.Run("should return 500 when service returns an error", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		require.NoError(t, err)
		_ = db.AutoMigrate(&model.Deposit{})

		controller := deposit.NewDepositController(&MockDepositServiceWithError{}, db)
		app := fiber.New()
		app.Post("/deposit", controller.CreateDepositHandler)

		validUUID := uuid.NewV4()
		reqBody := fmt.Sprintf(`{
			"customer_id": "%s",
			"amount": 1000,
			"interest_rate": 5,
			"term_months": 12,
			"start_date": "2023-04-01",
			"note": "Test deposit"
		}`, validUUID.String())

		req := httptest.NewRequest(http.MethodPost, "/deposit", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var body map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&body)

		meta := body["meta"].(map[string]interface{})
		assert.Equal(t, false, meta["success"])
		assert.Equal(t, "failed to create deposit", meta["message"])
	})
}
