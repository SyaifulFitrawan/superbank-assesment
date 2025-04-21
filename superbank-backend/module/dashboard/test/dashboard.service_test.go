package dashboard

import (
	"bank-backend/module/dashboard"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDashboardRepo struct {
	mock.Mock
}

func (m *MockDashboardRepo) GetTotals(ctx context.Context) (dashboard.DashboardTotalCounts, error) {
	args := m.Called(ctx)
	return args.Get(0).(dashboard.DashboardTotalCounts), args.Error(1)
}

func (m *MockDashboardRepo) CountByAccountType(ctx context.Context) ([]dashboard.AccountType, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dashboard.AccountType), args.Error(1)
}

func (m *MockDashboardRepo) GetCustomerDepositGroups(ctx context.Context) ([]dashboard.CustomerDepositOrPocketGroup, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dashboard.CustomerDepositOrPocketGroup), args.Error(1)
}

func (m *MockDashboardRepo) GetCustomerPocketGroups(ctx context.Context) ([]dashboard.CustomerDepositOrPocketGroup, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dashboard.CustomerDepositOrPocketGroup), args.Error(1)
}

func TestDashboardService_GetDashboard(t *testing.T) {
	t.Run("should return dashboard data successfully", func(t *testing.T) {
		mockRepo := new(MockDashboardRepo)
		ctx := context.Background()

		mockRepo.On("GetTotals", ctx).Return(dashboard.DashboardTotalCounts{
			TotalCustomers: 100,
			TotalDeposits:  200,
			TotalPockets:   50,
		}, nil)

		mockRepo.On("CountByAccountType", ctx).Return([]dashboard.AccountType{
			{AccountType: "Regular", Count: 70},
			{AccountType: "Premium", Count: 30},
		}, nil)

		mockRepo.On("GetCustomerDepositGroups", ctx).Return([]dashboard.CustomerDepositOrPocketGroup{
			{RangeLabel: "0 - 1M", Count: 40},
		}, nil)

		mockRepo.On("GetCustomerPocketGroups", ctx).Return([]dashboard.CustomerDepositOrPocketGroup{
			{RangeLabel: "0 - 1M", Count: 25},
		}, nil)

		service := dashboard.NewDashboardService(mockRepo)
		result, err := service.GetDashboard(ctx)

		assert.NoError(t, err)
		assert.Equal(t, int64(100), result.Total.TotalCustomers)
		assert.Equal(t, "Regular", result.Type[0].AccountType)
		assert.Equal(t, int64(25), result.Pocket[0].Count)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when GetTotals fails", func(t *testing.T) {
		mockRepo := new(MockDashboardRepo)
		ctx := context.Background()

		mockRepo.On("GetTotals", ctx).Return(dashboard.DashboardTotalCounts{}, errors.New("failed to get totals"))

		service := dashboard.NewDashboardService(mockRepo)
		_, err := service.GetDashboard(ctx)

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to get totals")
	})

	t.Run("should return error when CountByAccountType fails", func(t *testing.T) {
		mockRepo := new(MockDashboardRepo)
		ctx := context.Background()

		mockRepo.On("GetTotals", ctx).Return(dashboard.DashboardTotalCounts{}, nil)
		mockRepo.On("CountByAccountType", ctx).Return([]dashboard.AccountType{}, assert.AnError)

		service := dashboard.NewDashboardService(mockRepo)
		result, err := service.GetDashboard(ctx)

		assert.Error(t, err)
		assert.Empty(t, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when GetCustomerDepositGroups fails", func(t *testing.T) {
		mockRepo := new(MockDashboardRepo)
		ctx := context.Background()

		mockRepo.On("GetTotals", ctx).Return(dashboard.DashboardTotalCounts{}, nil)
		mockRepo.On("CountByAccountType", ctx).Return([]dashboard.AccountType{}, nil)
		mockRepo.On("GetCustomerDepositGroups", ctx).Return([]dashboard.CustomerDepositOrPocketGroup{}, assert.AnError)

		service := dashboard.NewDashboardService(mockRepo)
		result, err := service.GetDashboard(ctx)

		assert.Error(t, err)
		assert.Empty(t, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when GetCustomerPocketGroups fails", func(t *testing.T) {
		mockRepo := new(MockDashboardRepo)
		ctx := context.Background()

		mockRepo.On("GetTotals", ctx).Return(dashboard.DashboardTotalCounts{}, nil)
		mockRepo.On("CountByAccountType", ctx).Return([]dashboard.AccountType{}, nil)
		mockRepo.On("GetCustomerDepositGroups", ctx).Return([]dashboard.CustomerDepositOrPocketGroup{}, nil)
		mockRepo.On("GetCustomerPocketGroups", ctx).Return([]dashboard.CustomerDepositOrPocketGroup{}, assert.AnError)

		service := dashboard.NewDashboardService(mockRepo)
		result, err := service.GetDashboard(ctx)

		assert.Error(t, err)
		assert.Empty(t, result)

		mockRepo.AssertExpectations(t)
	})
}
