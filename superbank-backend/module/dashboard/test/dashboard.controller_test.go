package dashboard

import (
	"bank-backend/module/dashboard"
	"context"
	"encoding/json" // Tambahkan import errors untuk error kustom
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDashboardService struct {
	mock.Mock
}

func (m *mockDashboardService) GetDashboard(ctx context.Context) (dashboard.DashboardResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).(dashboard.DashboardResponse), args.Error(1)
}

func setupApp(ctrl *dashboard.DashboardController) *fiber.App {
	app := fiber.New()
	app.Get("/dashboard", ctrl.GetDashboard)
	return app
}

func TestGetDashboard(t *testing.T) {
	t.Run("should return 200 and dashboard data on success", func(t *testing.T) {
		mockSvc := new(mockDashboardService)
		ctrl := dashboard.NewDashboardController(mockSvc)
		app := setupApp(ctrl)

		expectedResponse := dashboard.DashboardResponse{
			Total: dashboard.DashboardTotalCounts{
				TotalCustomers: 100,
				TotalDeposits:  500,
				TotalPockets:   300,
			},
			Type: []dashboard.AccountType{
				{AccountType: "Silver", Count: 200},
				{AccountType: "Gold", Count: 150},
				{AccountType: "Platinum", Count: 150},
			},
			Deposit: []dashboard.CustomerDepositOrPocketGroup{
				{RangeLabel: "0-1", Count: 50},
				{RangeLabel: "2-3", Count: 60},
				{RangeLabel: "4-5", Count: 70},
				{RangeLabel: "6+", Count: 80},
			},
			Pocket: []dashboard.CustomerDepositOrPocketGroup{
				{RangeLabel: "0-1", Count: 50},
				{RangeLabel: "2-3", Count: 60},
				{RangeLabel: "4-5", Count: 70},
				{RangeLabel: "6+", Count: 80},
			},
		}

		mockSvc.On("GetDashboard", mock.Anything).Return(expectedResponse, nil)

		req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&body)

		assert.Equal(t, true, body["meta"].(map[string]interface{})["success"])
		assert.Equal(t, "Get dashboard successfully", body["meta"].(map[string]interface{})["message"])
	})

	t.Run("should return 500 when service returns error", func(t *testing.T) {
		mockSvc := new(mockDashboardService)
		ctrl := dashboard.NewDashboardController(mockSvc)
		app := setupApp(ctrl)

		mockSvc.On("GetDashboard", mock.Anything).Return(dashboard.DashboardResponse{}, errors.New("failed")).Once()

		req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var body map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		meta, ok := body["meta"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, false, meta["success"])
		assert.Equal(t, "Get dashboard failed", meta["message"])
		mockSvc.AssertExpectations(t)
	})
}
