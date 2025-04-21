package user

import (
	"bank-backend/model"
	"bank-backend/module/user"
	"context"
	"errors"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) Create(ctx context.Context, input *model.User) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *mockUserRepo) List(ctx context.Context, limit, offset int, search string) ([]model.User, int64, error) {
	args := m.Called(ctx, limit, offset, search)
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

func (m *mockUserRepo) Detail(ctx context.Context, id string) (*model.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *mockUserRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*model.User), args.Error(1)
}

func TestCreate(t *testing.T) {
	t.Run("should return success", func(t *testing.T) {
		mockRepo := new(mockUserRepo)
		svc := user.NewUserService(mockRepo)

		user := &model.User{
			ID:       uuid.NewV4(),
			Email:    "test@example.com",
			Password: "secret123",
		}

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

		res, err := svc.Create(context.Background(), user)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEqual(t, "secret123", res.Password)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return missing fields", func(t *testing.T) {
		svc := user.NewUserService(nil)

		res, err := svc.Create(context.Background(), &model.User{
			Email:    "",
			Password: "",
		})

		assert.Nil(t, res)
		assert.EqualError(t, err, "email and password are required")
	})

	t.Run("should return error when hashing fails", func(t *testing.T) {
		originalFunc := bcrypt.GenerateFromPassword
		defer func() { user.HashPasswordFunc = originalFunc }()

		user.HashPasswordFunc = func(password []byte, cost int) ([]byte, error) {
			return nil, errors.New("hashing failed")
		}

		svc := user.NewUserService(&mockUserRepo{})
		input := &model.User{
			Email:    "fail@hash.com",
			Password: "12345678",
		}

		res, err := svc.Create(context.Background(), input)

		assert.Nil(t, res)
		assert.EqualError(t, err, "hashing failed")
	})

	t.Run("should return repo error", func(t *testing.T) {
		mockRepo := new(mockUserRepo)
		svc := user.NewUserService(mockRepo)

		user := &model.User{
			ID:       uuid.NewV4(),
			Email:    "fail@example.com",
			Password: "secret123",
		}

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

		res, err := svc.Create(context.Background(), user)

		assert.Nil(t, res)
		assert.EqualError(t, err, "db error")
		mockRepo.AssertExpectations(t)
	})
}

func TestList(t *testing.T) {
	t.Run("should return users with pagination successfully", func(t *testing.T) {
		mockRepo := new(mockUserRepo)
		svc := user.NewUserService(mockRepo)

		page := 2
		limit := 2
		search := "john"
		offset := (page - 1) * limit

		mockUsers := []model.User{
			{ID: uuid.NewV4(), Email: "john3@example.com"},
			{ID: uuid.NewV4(), Email: "john4@example.com"},
		}
		totalCount := int64(5)

		mockRepo.On("List", mock.Anything, limit, offset, search).Return(mockUsers, totalCount, nil)

		users, paginator, err := svc.List(context.Background(), page, limit, search)

		assert.NoError(t, err)
		assert.Equal(t, mockUsers, users)
		assert.NotNil(t, paginator)
		assert.Equal(t, int(totalCount), paginator.ItemCount)
		assert.Equal(t, 3, paginator.PageCount)
		assert.Equal(t, page, paginator.Page)
		assert.True(t, paginator.HasPrevPage)
		assert.True(t, paginator.HasNextPage)
		assert.NotNil(t, paginator.PrevPage)
		assert.Equal(t, 1, *paginator.PrevPage)
		assert.NotNil(t, paginator.NextPage)
		assert.Equal(t, 3, *paginator.NextPage)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repo fails", func(t *testing.T) {
		mockRepo := new(mockUserRepo)
		svc := user.NewUserService(mockRepo)

		page := 1
		limit := 2
		search := "error"
		offset := (page - 1) * limit

		mockRepo.On("List", mock.Anything, limit, offset, search).Return([]model.User{}, int64(0), errors.New("list failed"))

		users, paginator, err := svc.List(context.Background(), page, limit, search)

		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Nil(t, paginator)
		assert.EqualError(t, err, "list failed")

		mockRepo.AssertExpectations(t)
	})
}

func TestDetail(t *testing.T) {
	t.Run("should return user details successfully", func(t *testing.T) {
		mockRepo := new(mockUserRepo)
		svc := user.NewUserService(mockRepo)

		userID := uuid.NewV4()
		expectedUser := &model.User{
			ID:    userID,
			Email: "test@example.com",
		}

		mockRepo.On("Detail", mock.Anything, userID.String()).Return(expectedUser, nil)

		res, err := svc.Detail(context.Background(), userID.String())

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, expectedUser.ID, res.ID)
		assert.Equal(t, expectedUser.Email, res.Email)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		mockRepo := new(mockUserRepo)
		svc := user.NewUserService(mockRepo)

		userID := uuid.NewV4().String()

		mockRepo.On("Detail", mock.Anything, userID).Return((*model.User)(nil), errors.New("user not found"))

		res, err := svc.Detail(context.Background(), userID)

		assert.Nil(t, res)
		assert.EqualError(t, err, "user not found")

		mockRepo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	t.Run("should delete user successfully", func(t *testing.T) {
		mockRepo := new(mockUserRepo)
		svc := user.NewUserService(mockRepo)

		userID := uuid.NewV4().String()

		mockRepo.On("Delete", mock.Anything, userID).Return(nil)

		err := svc.Delete(context.Background(), userID)

		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when deletion fails", func(t *testing.T) {
		mockRepo := new(mockUserRepo)
		svc := user.NewUserService(mockRepo)

		userID := uuid.NewV4().String()

		mockRepo.On("Delete", mock.Anything, userID).Return(errors.New("delete failed"))

		err := svc.Delete(context.Background(), userID)

		assert.EqualError(t, err, "delete failed")

		mockRepo.AssertExpectations(t)
	})
}
