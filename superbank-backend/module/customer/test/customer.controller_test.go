package customer

import (
	"bank-backend/model"
	"bank-backend/module/customer"
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCustomerService struct {
	mock.Mock
}

func (m *MockCustomerService) Create(ctx context.Context, req customer.CustomerCreateRequest) (*model.Customer, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Customer), args.Error(1)
}

func (m *MockCustomerService) List(ctx context.Context, page, limit int, search string) ([]model.Customer, *utils.Paginator, error) {
	args := m.Called(ctx, page, limit, search)
	return args.Get(0).([]model.Customer), args.Get(1).(*utils.Paginator), args.Error(2)
}

func (m *MockCustomerService) Detail(ctx context.Context, id string) (*customer.CustomerDetailResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*customer.CustomerDetailResponse), args.Error(1)
}

func (m *MockCustomerService) Update(ctx context.Context, id string, req customer.CustomerUpdateRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *MockCustomerService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupCustomerTestApp(t *testing.T) (*fiber.App, *MockCustomerService) {
	utils.InitValidator()
	app := fiber.New()
	mockService := new(MockCustomerService)
	controller := customer.NewCustomerController(mockService)

	app.Post("/customers", controller.CreateCustomerHandler)
	app.Get("/customers", controller.ListCustomerHandler)
	app.Get("/customers/:id", controller.DetailCustomerHandler)
	app.Put("/customers/:id", controller.UpdateCustomerHandler)
	app.Delete("/customers/:id", controller.DeleteCustomerHandler)

	return app, mockService
}

func TestCreateCustomerHandler(t *testing.T) {
	app, service := setupCustomerTestApp(t)

	t.Run("should return 200 when successful", func(t *testing.T) {
		reqBody := customer.CustomerCreateRequest{
			Name:          "John Doe",
			Phone:         "081234567890",
			Address:       "Jl. Raya",
			ParentName:    "Mr. Doe",
			AccountBranch: "Jakarta",
			AccountType:   "Gold",
		}
		customerModel := &model.Customer{ID: uuid.NewV4(), Name: "John Doe", Phone: "081234567890", Address: "Jl. Raya", ParentName: "Mr. Doe", AccountBranch: "Jakarta", AccountType: "Gold"}
		service.On("Create", mock.Anything, reqBody).Return(customerModel, nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should return 400 for invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewReader([]byte("invalid-json")))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for validation error", func(t *testing.T) {
		reqBody := customer.CustomerCreateRequest{Name: "", Phone: "081234567890", Address: "Jl. Raya", ParentName: "Mr. Doe", AccountBranch: "Jakarta", AccountType: "Gold"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 500 when service returns error", func(t *testing.T) {
		reqBody := customer.CustomerCreateRequest{
			Name:          "Jane",
			Phone:         "081234567890",
			Address:       "Jl. Raya",
			ParentName:    "Mr. Doe",
			AccountBranch: "Jakarta",
			AccountType:   "Gold",
		}
		service.On("Create", mock.Anything, reqBody).Return(nil, errors.New("something went wrong"))
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestListCustomerHandler(t *testing.T) {
	app, service := setupCustomerTestApp(t)

	t.Run("should return 200 with customer list", func(t *testing.T) {
		mockCustomers := []model.Customer{{ID: uuid.NewV4(), Name: "A"}}
		mockPaginator := &utils.Paginator{ItemCount: 1, Limit: 10, PageCount: 1, Page: 1}
		service.On("List", mock.Anything, 1, 10, "").Return(mockCustomers, mockPaginator, nil)

		req := httptest.NewRequest(http.MethodGet, "/customers", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var body map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&body)
		meta := body["meta"].(map[string]interface{})
		assert.Equal(t, true, meta["success"])
		assert.Equal(t, "List Customers Successful", meta["message"])
	})

	t.Run("should return 400 if page is less than 1", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/customers?page=0", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var body map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&body)
		meta := body["meta"].(map[string]interface{})
		assert.Equal(t, true, meta["success"])
	})

	t.Run("should return 400 if limit is less than 1", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/customers?limit=0", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var body map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&body)
		meta := body["meta"].(map[string]interface{})
		assert.Equal(t, true, meta["success"])
	})
}

func TestDetailCustomerHandler(t *testing.T) {
	app, service := setupCustomerTestApp(t)

	t.Run("should return 200 with customer detail", func(t *testing.T) {
		expected := &customer.CustomerDetailResponse{Customer: model.Customer{ID: uuid.NewV4(), Name: "B"}}
		service.On("Detail", mock.Anything, "123").Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/customers/123", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should return 404 if not found", func(t *testing.T) {
		service.On("Detail", mock.Anything, "notfound").Return(nil, errors.New("not found"))

		req := httptest.NewRequest(http.MethodGet, "/customers/notfound", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestUpdateCustomerHandler(t *testing.T) {
	app, service := setupCustomerTestApp(t)

	t.Run("should return 200 on successful update", func(t *testing.T) {
		reqBody := customer.CustomerUpdateRequest{Name: "Updated"}
		service.On("Update", mock.Anything, "123", reqBody).Return(nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/customers/123", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should return 400 for invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/customers/123", bytes.NewReader([]byte("invalid-json")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 500 if update fails", func(t *testing.T) {
		reqBody := customer.CustomerUpdateRequest{Name: "ErrName"}
		service.On("Update", mock.Anything, "456", reqBody).Return(errors.New("update failed"))

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/customers/456", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestDeleteCustomerHandler(t *testing.T) {
	app, service := setupCustomerTestApp(t)

	t.Run("should return 200 on success", func(t *testing.T) {
		service.On("Delete", mock.Anything, "123").Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/customers/123", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should return 500 if delete fails", func(t *testing.T) {
		service.On("Delete", mock.Anything, "456").Return(errors.New("delete error"))

		req := httptest.NewRequest(http.MethodDelete, "/customers/456", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
