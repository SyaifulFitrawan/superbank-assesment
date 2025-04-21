package pocket

import (
	"bank-backend/model"
	"bank-backend/module/pocket"
	"context"
	"errors"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPocketRepo struct {
	mock.Mock
}

func (m *mockPocketRepo) Create(ctx context.Context, pocket *model.Pocket) error {
	args := m.Called(ctx, pocket)
	return args.Error(0)
}

func (m *mockPocketRepo) Detail(ctx context.Context, id string) (*model.Pocket, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Pocket), args.Error(1)
}

func (m *mockPocketRepo) Update(ctx context.Context, id string, pocket *model.Pocket) error {
	args := m.Called(ctx, id, pocket)
	return args.Error(0)
}

func (m *mockPocketRepo) Deactivated(ctx context.Context, id string, input any) error {
	args := m.Called(ctx, id, input)
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

func TestCreatePocket(t *testing.T) {
	pocketRepo := new(mockPocketRepo)
	customerRepo := new(mockCustomerRepo)

	service := pocket.NewPocketService(pocketRepo, customerRepo)

	input := pocket.CreatePocketRequest{
		CustomerID:   uuid.NewV4().String(),
		Name:         "Holiday Fund",
		TargetAmount: func() *float64 { f := 100000.0; return &f }(),
		TargetDate:   func() *string { s := "2025-12-31"; return &s }(),
	}

	t.Run("should create pocket successfully", func(t *testing.T) {
		pocketRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Pocket")).Return(nil).Once()

		result, err := service.Create(context.TODO(), input)

		assert.NoError(t, err)
		assert.Equal(t, input.Name, result.Name)
		assert.Equal(t, *input.TargetAmount, *result.TargetAmount)

		pocketRepo.AssertExpectations(t)
	})

	t.Run("should fail with invalid target date format", func(t *testing.T) {
		inputInvalidDate := input
		invalidDate := "2025-31-12"
		inputInvalidDate.TargetDate = &invalidDate

		result, err := service.Create(context.TODO(), inputInvalidDate)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "invalid target date format", err.Error())
	})

	t.Run("should fail when repository create returns error", func(t *testing.T) {
		pocketRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Pocket")).Return(errors.New("create failed")).Once()

		result, err := service.Create(context.TODO(), input)

		assert.Nil(t, result)
		assert.EqualError(t, err, "create failed")

		pocketRepo.AssertExpectations(t)
	})
}

func TestTopUpPocket(t *testing.T) {
	pocketRepo := new(mockPocketRepo)
	customerRepo := new(mockCustomerRepo)
	service := pocket.NewPocketService(pocketRepo, customerRepo)

	basePocket := func() *model.Pocket {
		return &model.Pocket{
			ID:         uuid.NewV4(),
			CustomerID: uuid.NewV4(),
			Name:       "Vacation",
			IsActive:   true,
			Balance:    1000,
		}
	}

	baseCustomer := func(customerID uuid.UUID) *model.Customer {
		return &model.Customer{
			ID:      customerID,
			Name:    "John",
			Balance: 2000,
		}
	}

	input := pocket.TopUpOrWithdrawPocketRequest{
		Amount: 1000,
	}

	t.Run("should top up pocket successfully", func(t *testing.T) {
		p := basePocket()
		c := baseCustomer(p.CustomerID)

		pocketRepo.On("Detail", mock.Anything, p.ID.String()).Return(p, nil)
		customerRepo.On("Detail", mock.Anything, p.CustomerID.String()).Return(c, nil)
		pocketRepo.On("Update", mock.Anything, p.ID.String(), mock.AnythingOfType("*model.Pocket")).Return(nil)
		customerRepo.On("Update", mock.Anything, c.ID.String(), mock.AnythingOfType("*model.Customer")).Return(nil)

		err := service.TopUp(context.TODO(), p.ID.String(), input)

		assert.NoError(t, err)
	})

	t.Run("should fail when pocket is inactive", func(t *testing.T) {
		p := basePocket()
		p.IsActive = false

		pocketRepo.On("Detail", mock.Anything, p.ID.String()).Return(p, nil)

		err := service.TopUp(context.TODO(), p.ID.String(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "inactive")
	})

	t.Run("should return error when pocket not found", func(t *testing.T) {
		p := basePocket()

		pocketRepo.On("Detail", mock.Anything, p.ID.String()).Return(&model.Pocket{}, errors.New("pocket not found"))

		err := service.TopUp(context.TODO(), p.ID.String(), input)

		assert.Error(t, err)
		assert.Equal(t, "pocket not found", err.Error())
	})

	t.Run("should return error when customer not found", func(t *testing.T) {
		p := basePocket()

		pocketRepo.On("Detail", mock.Anything, p.ID.String()).Return(p, nil)
		customerRepo.On("Detail", mock.Anything, p.CustomerID.String()).Return(&model.Customer{}, errors.New("customer not found"))

		err := service.TopUp(context.TODO(), p.ID.String(), input)

		assert.Error(t, err)
		assert.Equal(t, "customer not found", err.Error())
	})

	t.Run("should return error when customer has insufficient balance", func(t *testing.T) {
		p := basePocket()
		c := baseCustomer(p.CustomerID)
		c.Balance = 500

		pocketRepo.On("Detail", mock.Anything, p.ID.String()).Return(p, nil)
		customerRepo.On("Detail", mock.Anything, c.ID.String()).Return(c, nil)

		err := service.TopUp(context.TODO(), p.ID.String(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not have enough balance")
	})

	t.Run("should return error when updating pocket fails", func(t *testing.T) {
		p := basePocket()
		c := baseCustomer(p.CustomerID)

		pocketRepo.On("Detail", mock.Anything, p.ID.String()).Return(p, nil)
		customerRepo.On("Detail", mock.Anything, c.ID.String()).Return(c, nil)
		pocketRepo.On("Update", mock.Anything, p.ID.String(), mock.Anything).Return(errors.New("update failed"))

		err := service.TopUp(context.TODO(), p.ID.String(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update failed")
	})

	t.Run("should return error when updating customer fails", func(t *testing.T) {
		p := basePocket()
		c := baseCustomer(p.CustomerID)

		pocketRepo.On("Detail", mock.Anything, p.ID.String()).Return(p, nil)
		customerRepo.On("Detail", mock.Anything, c.ID.String()).Return(c, nil)
		pocketRepo.On("Update", mock.Anything, p.ID.String(), mock.Anything).Return(nil)
		customerRepo.On("Update", mock.Anything, c.ID.String(), mock.Anything).Return(errors.New("update customer failed"))

		err := service.TopUp(context.TODO(), p.ID.String(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update customer failed")
	})
}

func TestWithdrawPocket(t *testing.T) {
	pocketRepo := new(mockPocketRepo)
	customerRepo := new(mockCustomerRepo)
	service := pocket.NewPocketService(pocketRepo, customerRepo)

	basePocket := func() *model.Pocket {
		return &model.Pocket{
			ID:         uuid.NewV4(),
			CustomerID: uuid.NewV4(),
			Name:       "Emergency Fund",
			IsActive:   true,
			Balance:    1500,
		}
	}

	baseCustomer := func(customerID uuid.UUID) *model.Customer {
		return &model.Customer{
			ID:      customerID,
			Name:    "Alice",
			Balance: 500,
		}
	}

	input := pocket.TopUpOrWithdrawPocketRequest{
		Amount: 1000,
	}

	t.Run("should withdraw from pocket successfully", func(t *testing.T) {
		p := basePocket()
		c := baseCustomer(p.CustomerID)

		pocketRepo.On("Detail", mock.Anything, p.ID.String()).Return(p, nil)
		customerRepo.On("Detail", mock.Anything, p.CustomerID.String()).Return(c, nil)
		pocketRepo.On("Update", mock.Anything, p.ID.String(), mock.AnythingOfType("*model.Pocket")).Return(nil)
		customerRepo.On("Update", mock.Anything, c.ID.String(), mock.AnythingOfType("*model.Customer")).Return(nil)

		err := service.Withdrawn(context.TODO(), p.ID.String(), input)

		assert.NoError(t, err)
	})

	t.Run("should fail when pocket is inactive", func(t *testing.T) {
		p := basePocket()
		p.IsActive = false

		pocketRepo.On("Detail", mock.Anything, p.ID.String()).Return(p, nil)

		err := service.Withdrawn(context.TODO(), p.ID.String(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "inactive")
	})

	t.Run("should return error when pocket not found", func(t *testing.T) {
		p := basePocket()

		pocketRepo.On("Detail", mock.Anything, p.ID.String()).Return(&model.Pocket{}, errors.New("pocket not found"))

		err := service.Withdrawn(context.TODO(), p.ID.String(), input)

		assert.Error(t, err)
		assert.Equal(t, "pocket not found", err.Error())
	})

	t.Run("should return error when customer not found", func(t *testing.T) {
		p := basePocket()

		pocketRepo.On("Detail", mock.Anything, p.ID.String()).Return(p, nil)
		customerRepo.On("Detail", mock.Anything, p.CustomerID.String()).Return(&model.Customer{}, errors.New("customer not found"))

		err := service.Withdrawn(context.TODO(), p.ID.String(), input)

		assert.Error(t, err)
		assert.Equal(t, "customer not found", err.Error())
	})

	t.Run("should return error when pocket has insufficient balance", func(t *testing.T) {
		p := basePocket()
		p.Balance = 500
		c := baseCustomer(p.CustomerID)

		pocketRepo.On("Detail", mock.Anything, p.ID.String()).Return(p, nil)
		customerRepo.On("Detail", mock.Anything, c.ID.String()).Return(c, nil)

		err := service.Withdrawn(context.TODO(), p.ID.String(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not have enough balance")
	})

	t.Run("should return error when updating pocket fails", func(t *testing.T) {
		p := basePocket()
		c := baseCustomer(p.CustomerID)

		pocketRepo.On("Detail", mock.Anything, p.ID.String()).Return(p, nil)
		customerRepo.On("Detail", mock.Anything, c.ID.String()).Return(c, nil)
		pocketRepo.On("Update", mock.Anything, p.ID.String(), mock.Anything).Return(errors.New("update pocket failed"))

		err := service.Withdrawn(context.TODO(), p.ID.String(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update pocket failed")
	})

	t.Run("should return error when updating customer fails", func(t *testing.T) {
		p := basePocket()
		c := baseCustomer(p.CustomerID)

		pocketRepo.On("Detail", mock.Anything, p.ID.String()).Return(p, nil)
		customerRepo.On("Detail", mock.Anything, c.ID.String()).Return(c, nil)
		pocketRepo.On("Update", mock.Anything, p.ID.String(), mock.Anything).Return(nil)
		customerRepo.On("Update", mock.Anything, c.ID.String(), mock.Anything).Return(errors.New("update customer failed"))

		err := service.Withdrawn(context.TODO(), p.ID.String(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update customer failed")
	})
}

func TestDeactivatePocket(t *testing.T) {
	pocketRepo := new(mockPocketRepo)
	customerRepo := new(mockCustomerRepo)
	service := pocket.NewPocketService(pocketRepo, customerRepo)

	basePocket := func() *model.Pocket {
		return &model.Pocket{
			ID:         uuid.NewV4(),
			CustomerID: uuid.NewV4(),
			Name:       "Holiday Fund",
			IsActive:   true,
			Balance:    1000,
		}
	}

	baseCustomer := func(customerID uuid.UUID) *model.Customer {
		return &model.Customer{
			ID:      customerID,
			Name:    "Bob",
			Balance: 300,
		}
	}

	t.Run("should deactivate pocket with balance transfer", func(t *testing.T) {
		p := basePocket()
		c := baseCustomer(p.CustomerID)

		pocketRepo.On("Detail", mock.Anything, p.ID.String()).Return(p, nil)
		customerRepo.On("Detail", mock.Anything, c.ID.String()).Return(c, nil)
		customerRepo.On("Update", mock.Anything, c.ID.String(), mock.AnythingOfType("*model.Customer")).Return(nil)
		pocketRepo.On("Deactivated", mock.Anything, p.ID.String(), mock.AnythingOfType("map[string]interface {}")).Return(nil)

		err := service.Deactivated(context.TODO(), p.ID.String())
		assert.NoError(t, err)
	})

	t.Run("should deactivate pocket with zero balance", func(t *testing.T) {
		p := basePocket()
		p.Balance = 0

		pocketRepo.On("Detail", mock.Anything, p.ID.String()).Return(p, nil)
		pocketRepo.On("Deactivated", mock.Anything, p.ID.String(), mock.AnythingOfType("map[string]interface {}")).Return(nil)

		err := service.Deactivated(context.TODO(), p.ID.String())
		assert.NoError(t, err)
	})

	t.Run("should return error when pocket not found", func(t *testing.T) {
		p := basePocket()

		pocketRepo.On("Detail", mock.Anything, p.ID.String()).Return(&model.Pocket{}, errors.New("pocket not found"))

		err := service.Deactivated(context.TODO(), p.ID.String())
		assert.Error(t, err)
		assert.Equal(t, "pocket not found", err.Error())
	})

	t.Run("should return error when pocket already inactive", func(t *testing.T) {
		p := basePocket()
		p.IsActive = false

		pocketRepo.On("Detail", mock.Anything, p.ID.String()).Return(p, nil)

		err := service.Deactivated(context.TODO(), p.ID.String())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "has been inactive")
	})

	t.Run("should return error when customer not found", func(t *testing.T) {
		p := basePocket()

		pocketRepo.On("Detail", mock.Anything, p.ID.String()).Return(p, nil)
		customerRepo.On("Detail", mock.Anything, p.CustomerID.String()).Return(&model.Customer{}, errors.New("customer not found"))

		err := service.Deactivated(context.TODO(), p.ID.String())
		assert.Error(t, err)
		assert.Equal(t, "customer not found", err.Error())
	})

	t.Run("should return error when update customer fails", func(t *testing.T) {
		p := basePocket()
		c := baseCustomer(p.CustomerID)

		pocketRepo.On("Detail", mock.Anything, p.ID.String()).Return(p, nil)
		customerRepo.On("Detail", mock.Anything, c.ID.String()).Return(c, nil)
		customerRepo.On("Update", mock.Anything, c.ID.String(), mock.Anything).Return(errors.New("update customer failed"))

		err := service.Deactivated(context.TODO(), p.ID.String())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update customer failed")
	})

	t.Run("should return error when pocket deactivation fails", func(t *testing.T) {
		p := basePocket()
		c := baseCustomer(p.CustomerID)

		pocketRepo.On("Detail", mock.Anything, p.ID.String()).Return(p, nil)
		customerRepo.On("Detail", mock.Anything, c.ID.String()).Return(c, nil)
		customerRepo.On("Update", mock.Anything, c.ID.String(), mock.Anything).Return(nil)
		pocketRepo.On("Deactivated", mock.Anything, p.ID.String(), mock.Anything).Return(errors.New("deactivation failed"))

		err := service.Deactivated(context.TODO(), p.ID.String())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "deactivation failed")
	})
}
