package customer

import (
	"bank-backend/model"
	"bank-backend/module/customer"
	"context"
	"errors"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCustomerRepo struct {
	mock.Mock
}

func (m *mockCustomerRepo) Create(ctx context.Context, input *model.Customer) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *mockCustomerRepo) List(ctx context.Context, limit, offset int, search string) ([]model.Customer, int64, error) {
	args := m.Called(ctx, limit, offset, search)
	return args.Get(0).([]model.Customer), args.Get(1).(int64), args.Error(2)
}

func (m *mockCustomerRepo) Detail(ctx context.Context, id string) (*model.Customer, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Customer), args.Error(1)
}

func (m *mockCustomerRepo) Update(ctx context.Context, id string, input *model.Customer) error {
	args := m.Called(ctx, id, input)
	return args.Error(0)
}

func (m *mockCustomerRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockCustomerRepo) AddBalance(ctx context.Context, customerID string, amount float64) error {
	args := m.Called(ctx, customerID, amount)
	return args.Error(0)
}

func TestCustomerService_Create(t *testing.T) {
	t.Run("should return success", func(t *testing.T) {
		mockRepo := new(mockCustomerRepo)
		svc := customer.NewCustomerService(mockRepo)

		input := customer.CustomerCreateRequest{
			Name:          "John",
			Phone:         "08123",
			Address:       "Jakarta",
			ParentName:    "Doe",
			AccountBranch: "001",
			AccountType:   "Savings",
		}

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

		res, err := svc.Create(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, input.Name, res.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return repo error", func(t *testing.T) {
		mockRepo := new(mockCustomerRepo)
		svc := customer.NewCustomerService(mockRepo)

		input := customer.CustomerCreateRequest{Name: "Error"}

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("create failed"))

		res, err := svc.Create(context.Background(), input)

		assert.Nil(t, res)
		assert.EqualError(t, err, "create failed")
		mockRepo.AssertExpectations(t)
	})
}

func TestCustomerService_List(t *testing.T) {
	t.Run("should return customers and paginator", func(t *testing.T) {
		mockRepo := new(mockCustomerRepo)
		svc := customer.NewCustomerService(mockRepo)

		page := 2
		limit := 2
		search := "john"
		offset := (page - 1) * limit

		mockCustomers := []model.Customer{{Name: "John Doe"}, {Name: "Jane Doe"}}
		total := int64(5)

		mockRepo.On("List", mock.Anything, limit, offset, search).Return(mockCustomers, total, nil)

		customers, paginator, err := svc.List(context.Background(), page, limit, search)

		assert.NoError(t, err)
		assert.Equal(t, mockCustomers, customers)
		assert.NotNil(t, paginator)
		assert.Equal(t, int(total), paginator.ItemCount)
		assert.Equal(t, 3, paginator.PageCount)
		assert.True(t, paginator.HasPrevPage)
		assert.True(t, paginator.HasNextPage)
		assert.Equal(t, 1, *paginator.PrevPage)
		assert.Equal(t, 3, *paginator.NextPage)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error on repo failure", func(t *testing.T) {
		mockRepo := new(mockCustomerRepo)
		svc := customer.NewCustomerService(mockRepo)

		mockRepo.On("List", mock.Anything, 10, 0, "").Return([]model.Customer{}, int64(0), errors.New("list error"))

		customers, paginator, err := svc.List(context.Background(), 1, 10, "")

		assert.Nil(t, customers)
		assert.Nil(t, paginator)
		assert.EqualError(t, err, "list error")

		mockRepo.AssertExpectations(t)
	})
}

func TestCustomerService_Detail(t *testing.T) {
	t.Run("should return customer detail", func(t *testing.T) {
		mockRepo := new(mockCustomerRepo)
		svc := customer.NewCustomerService(mockRepo)

		customer := &model.Customer{
			ID:      uuid.NewV4(),
			Name:    "John",
			Phone:   "08123",
			Address: "Jakarta",
			Deposits: []model.Deposit{
				{ID: uuid.NewV4(), Amount: 100},
			},
			Pockets: []model.Pocket{
				{ID: uuid.NewV4(), Balance: 50},
			},
		}

		mockRepo.On("Detail", mock.Anything, "123").Return(customer, nil)

		res, err := svc.Detail(context.Background(), "123")

		assert.NoError(t, err)
		assert.Equal(t, customer.Name, res.Customer.Name)
		assert.Equal(t, customer.Deposits, res.Deposits)
		assert.Equal(t, customer.Pockets, res.Pockets)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when not found", func(t *testing.T) {
		mockRepo := new(mockCustomerRepo)
		svc := customer.NewCustomerService(mockRepo)

		mockRepo.On("Detail", mock.Anything, "notfound").Return(&model.Customer{}, errors.New("not found"))

		res, err := svc.Detail(context.Background(), "notfound")

		assert.Nil(t, res)
		assert.EqualError(t, err, "not found")
		mockRepo.AssertExpectations(t)
	})
}

func TestCustomerService_Update(t *testing.T) {
	t.Run("should update successfully", func(t *testing.T) {
		mockRepo := new(mockCustomerRepo)
		svc := customer.NewCustomerService(mockRepo)

		input := customer.CustomerUpdateRequest{
			Name:          "Updated Name",
			Phone:         "09876",
			Address:       "Bandung",
			ParentName:    "New Parent",
			AccountBranch: "002",
			AccountType:   "Checking",
		}

		mockRepo.On("Update", mock.Anything, "123", mock.Anything).Return(nil)

		err := svc.Update(context.Background(), "123", input)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		mockRepo := new(mockCustomerRepo)
		svc := customer.NewCustomerService(mockRepo)

		mockRepo.On("Update", mock.Anything, "fail", mock.Anything).Return(errors.New("update failed"))

		err := svc.Update(context.Background(), "fail", customer.CustomerUpdateRequest{Name: "Fail"})

		assert.EqualError(t, err, "update failed")
		mockRepo.AssertExpectations(t)
	})
}

func TestCustomerService_Delete(t *testing.T) {
	t.Run("should delete successfully", func(t *testing.T) {
		mockRepo := new(mockCustomerRepo)
		svc := customer.NewCustomerService(mockRepo)

		mockRepo.On("Delete", mock.Anything, "123").Return(nil)

		err := svc.Delete(context.Background(), "123")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when delete fails", func(t *testing.T) {
		mockRepo := new(mockCustomerRepo)
		svc := customer.NewCustomerService(mockRepo)

		mockRepo.On("Delete", mock.Anything, "fail").Return(errors.New("delete failed"))

		err := svc.Delete(context.Background(), "fail")

		assert.EqualError(t, err, "delete failed")
		mockRepo.AssertExpectations(t)
	})
}
