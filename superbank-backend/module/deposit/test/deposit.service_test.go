package deposit

import (
	"bank-backend/model"
	"bank-backend/module/deposit"
	"context"
	"errors"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDepositRepo struct {
	mock.Mock
}

func (m *mockDepositRepo) Create(ctx context.Context, d *model.Deposit) error {
	args := m.Called(ctx, d)
	return args.Error(0)
}

func (m *mockDepositRepo) FindMatureUnwithdraw(ctx context.Context) ([]model.Deposit, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Deposit), args.Error(1)
}

func (m *mockDepositRepo) Update(ctx context.Context, id string, d *model.Deposit) error {
	args := m.Called(ctx, id, d)
	return args.Error(0)
}

type mockCustomerRepo struct {
	mock.Mock
}

func (m *mockCustomerRepo) Create(ctx context.Context, customer *model.Customer) error {
	return nil
}

func (m *mockCustomerRepo) List(ctx context.Context, limit, offset int, search string) ([]model.Customer, int64, error) {
	args := m.Called(ctx, limit, offset, search)
	return args.Get(0).([]model.Customer), args.Get(1).(int64), args.Error(2)
}

func (m *mockCustomerRepo) Detail(ctx context.Context, id string) (*model.Customer, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Customer), args.Error(1)
}

func (m *mockCustomerRepo) Update(ctx context.Context, id string, c *model.Customer) error {
	args := m.Called(ctx, id, c)
	return args.Error(0)
}

func (m *mockCustomerRepo) AddBalance(ctx context.Context, id string, amount float64) error {
	args := m.Called(ctx, id, amount)
	return args.Error(0)
}

func (m *mockCustomerRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestDepositService_Create(t *testing.T) {
	ctx := context.Background()
	mockDeposit := new(mockDepositRepo)
	mockCustomer := new(mockCustomerRepo)
	service := deposit.NewDepositService(mockDeposit, mockCustomer)

	t.Run("should create deposit successfully", func(t *testing.T) {
		mockDeposit.Mock = mock.Mock{}
		mockCustomer.Mock = mock.Mock{}

		customerID := uuid.NewV4().String()
		startDate := "2025-04-01"

		mockCustomer.On("Detail", ctx, customerID).Return(&model.Customer{
			ID:      uuid.FromStringOrNil(customerID),
			Name:    "Alice",
			Balance: 2000,
		}, nil)

		mockCustomer.On("Update", ctx, customerID, mock.Anything).Return(nil)
		mockDeposit.On("Create", ctx, mock.Anything).Return(nil)

		input := deposit.CreateDepositRequest{
			CustomerID:   customerID,
			Amount:       1000,
			InterestRate: 0.05,
			TermMonths:   6,
			StartDate:    startDate,
			Note:         "Initial deposit",
		}

		result, err := service.Create(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, input.Amount, result.Amount)
		assert.False(t, result.IsWithdrawn)

		mockCustomer.AssertExpectations(t)
		mockDeposit.AssertExpectations(t)
	})

	t.Run("should fail when customer not found", func(t *testing.T) {
		mockDeposit.Mock = mock.Mock{}
		mockCustomer.Mock = mock.Mock{}

		mockCustomer.On("Detail", mock.Anything, "invalid-id").Return(&model.Customer{}, errors.New("not found"))

		input := deposit.CreateDepositRequest{
			CustomerID: "invalid-id",
			Amount:     1000,
			StartDate:  "2025-01-01",
		}

		result, err := service.Create(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockCustomer.AssertExpectations(t)
	})

	t.Run("should fail when balance is not enough", func(t *testing.T) {
		mockDeposit.Mock = mock.Mock{}
		mockCustomer.Mock = mock.Mock{}

		mockCustomer.On("Detail", ctx, "low-balance").Return(&model.Customer{
			ID:      uuid.FromStringOrNil("low-balance"),
			Name:    "Bob",
			Balance: 100,
		}, nil)

		input := deposit.CreateDepositRequest{
			CustomerID: "low-balance",
			Amount:     1000,
			StartDate:  "2025-01-01",
		}

		result, err := service.Create(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not have enough balance")
		mockCustomer.AssertExpectations(t)
	})

	t.Run("should fail with invalid start date", func(t *testing.T) {
		mockDeposit.Mock = mock.Mock{}
		mockCustomer.Mock = mock.Mock{}

		mockCustomer.On("Detail", ctx, "bad-date").Return(&model.Customer{
			ID:      uuid.FromStringOrNil("bad-date"),
			Name:    "Charlie",
			Balance: 3000,
		}, nil)

		input := deposit.CreateDepositRequest{
			CustomerID: "bad-date",
			Amount:     1000,
			StartDate:  "wrong-format",
		}

		result, err := service.Create(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid start date")
		mockCustomer.AssertExpectations(t)
	})

	t.Run("should return error if customer update fails", func(t *testing.T) {
		mockDeposit := new(mockDepositRepo)
		mockCustomer := new(mockCustomerRepo)
		service := deposit.NewDepositService(mockDeposit, mockCustomer)

		customerID := uuid.NewV4().String()
		startDate := "2025-04-01"

		mockCustomer.On("Detail", ctx, customerID).Return(&model.Customer{
			ID:      uuid.FromStringOrNil(customerID),
			Name:    "Alice",
			Balance: 2000,
		}, nil)

		mockCustomer.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("update failed"))

		input := deposit.CreateDepositRequest{
			CustomerID:   customerID,
			Amount:       1000,
			InterestRate: 0.05,
			TermMonths:   6,
			StartDate:    startDate,
			Note:         "Initial deposit",
		}

		result, err := service.Create(ctx, input)

		assert.Nil(t, result)
		assert.EqualError(t, err, "update failed")

		mockCustomer.AssertExpectations(t)
		mockDeposit.AssertExpectations(t)
	})

	t.Run("should return error if deposit creation fails", func(t *testing.T) {
		mockDeposit := new(mockDepositRepo)
		mockCustomer := new(mockCustomerRepo)
		service := deposit.NewDepositService(mockDeposit, mockCustomer)

		customerID := uuid.NewV4().String()
		startDate := "2025-04-01"

		mockCustomer.On("Detail", ctx, customerID).Return(&model.Customer{
			ID:      uuid.FromStringOrNil(customerID),
			Name:    "Alice",
			Balance: 2000,
		}, nil)

		mockCustomer.On("Update", ctx, customerID, mock.Anything).Return(nil)
		mockDeposit.On("Create", ctx, mock.Anything).Return(errors.New("create failed"))

		input := deposit.CreateDepositRequest{
			CustomerID:   customerID,
			Amount:       1000,
			InterestRate: 0.05,
			TermMonths:   6,
			StartDate:    startDate,
			Note:         "Initial deposit",
		}

		result, err := service.Create(ctx, input)

		assert.Nil(t, result)
		assert.EqualError(t, err, "create failed")

		mockCustomer.AssertExpectations(t)
		mockDeposit.AssertExpectations(t)
	})
}

func TestDepositService_ProcessMatureDeposits(t *testing.T) {
	ctx := context.Background()

	t.Run("should process mature deposits successfully", func(t *testing.T) {
		mockDeposit := new(mockDepositRepo)
		mockCustomer := new(mockCustomerRepo)
		service := deposit.NewDepositService(mockDeposit, mockCustomer)

		customerID := uuid.NewV4()
		depositID := uuid.NewV4()

		deposits := []model.Deposit{
			{
				ID:           depositID,
				CustomerID:   customerID,
				Amount:       1000,
				InterestRate: 0.1,
				TermMonths:   12,
				Customer:     model.Customer{ID: customerID, Name: "John"},
			},
		}

		mockDeposit.On("FindMatureUnwithdraw", ctx).Return(deposits, nil)
		mockCustomer.On("AddBalance", ctx, customerID.String(), 1100.0).Return(nil)
		mockDeposit.On("Update", ctx, depositID.String(), mock.Anything).Return(nil)

		err := service.ProcessMatureDeposits(ctx)
		assert.NoError(t, err)

		mockDeposit.AssertExpectations(t)
		mockCustomer.AssertExpectations(t)
	})

	t.Run("should return nil if FindMatureUnwithdraw returns error", func(t *testing.T) {
		mockDeposit := new(mockDepositRepo)
		mockCustomer := new(mockCustomerRepo)
		service := deposit.NewDepositService(mockDeposit, mockCustomer)

		mockDeposit.On("FindMatureUnwithdraw", mock.Anything).Return([]model.Deposit{}, errors.New("not found"))

		err := service.ProcessMatureDeposits(ctx)
		assert.NoError(t, err)

		mockDeposit.AssertExpectations(t)
		mockCustomer.AssertExpectations(t)
	})

	t.Run("should skip deposit if AddBalance fails", func(t *testing.T) {
		mockDeposit := new(mockDepositRepo)
		mockCustomer := new(mockCustomerRepo)
		service := deposit.NewDepositService(mockDeposit, mockCustomer)

		customerID := uuid.NewV4()
		depositID := uuid.NewV4()

		deposits := []model.Deposit{
			{
				ID:           depositID,
				CustomerID:   customerID,
				Amount:       1000,
				InterestRate: 0.1,
				TermMonths:   12,
				Customer:     model.Customer{ID: customerID, Name: "John"},
			},
		}

		mockDeposit.On("FindMatureUnwithdraw", ctx).Return(deposits, nil)
		mockCustomer.On("AddBalance", ctx, customerID.String(), 1100.0).Return(errors.New("failed to add"))

		err := service.ProcessMatureDeposits(ctx)
		assert.NoError(t, err)

		mockDeposit.AssertExpectations(t)
		mockCustomer.AssertExpectations(t)
	})

	t.Run("should continue processing even if Update fails", func(t *testing.T) {
		mockDeposit := new(mockDepositRepo)
		mockCustomer := new(mockCustomerRepo)
		service := deposit.NewDepositService(mockDeposit, mockCustomer)

		customerID := uuid.NewV4()
		depositID := uuid.NewV4()

		deposits := []model.Deposit{
			{
				ID:           depositID,
				CustomerID:   customerID,
				Amount:       1000,
				InterestRate: 0.1,
				TermMonths:   12,
				Customer:     model.Customer{ID: customerID, Name: "John"},
			},
		}

		mockDeposit.On("FindMatureUnwithdraw", ctx).Return(deposits, nil)
		mockCustomer.On("AddBalance", ctx, customerID.String(), 1100.0).Return(nil)
		mockDeposit.On("Update", ctx, depositID.String(), mock.Anything).Return(errors.New("update error"))

		err := service.ProcessMatureDeposits(ctx)
		assert.NoError(t, err)

		mockDeposit.AssertExpectations(t)
		mockCustomer.AssertExpectations(t)
	})
}
