package auth

import (
	"bank-backend/model"
	"bank-backend/module/auth"
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

type mockJWTGenerator struct {
	mock.Mock
}

func (m *mockJWTGenerator) Generate(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func hash(pw string) string {
	b, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b)
}

func TestLogin(t *testing.T) {
	mockRepo := new(mockUserRepo)
	mockJWT := new(mockJWTGenerator)

	password := "password123"
	hashedPassword := hash(password)
	userID := uuid.NewV4()

	mockUser := &model.User{
		ID:       userID,
		Email:    "test@example.com",
		Password: hashedPassword,
	}

	mockRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(mockUser, nil)
	mockJWT.On("Generate", userID.String()).Return("mocked-token", nil)

	svc := auth.NewAuthService(mockRepo, mockJWT)

	res, err := svc.Login(context.Background(), "test@example.com", password)

	assert.NoError(t, err)
	assert.Equal(t, "mocked-token", res.Token)
	mockRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
}

func TestLogin_InvalidEmail(t *testing.T) {
	mockRepo := new(mockUserRepo)
	mockJWT := new(mockJWTGenerator)

	mockRepo.On("FindByEmail", mock.Anything, "notfound@example.com").Return(nil, errors.New("not found"))

	svc := auth.NewAuthService(mockRepo, mockJWT)

	res, err := svc.Login(context.Background(), "notfound@example.com", "secret")

	assert.Nil(t, res)
	assert.EqualError(t, err, "invalid credentials")
	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	mockRepo := new(mockUserRepo)
	mockJWT := new(mockJWTGenerator)

	hashed := hash("correct-password")

	mockUser := &model.User{
		ID:       uuid.NewV4(),
		Email:    "test@example.com",
		Password: hashed,
	}

	mockRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(mockUser, nil)

	svc := auth.NewAuthService(mockRepo, mockJWT)

	res, err := svc.Login(context.Background(), "test@example.com", "wrong-password")

	assert.Nil(t, res)
	assert.EqualError(t, err, "invalid credentials")
	mockRepo.AssertExpectations(t)
}

func TestLogin_GenerateTokenError(t *testing.T) {
	mockRepo := new(mockUserRepo)
	mockJWT := new(mockJWTGenerator)

	password := "secret"
	hashedPassword := hash(password)
	userID := uuid.NewV4()

	mockUser := &model.User{
		ID:       userID,
		Email:    "test@example.com",
		Password: hashedPassword,
	}

	mockRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(mockUser, nil)
	mockJWT.On("Generate", userID.String()).Return("", errors.New("jwt error"))

	svc := auth.NewAuthService(mockRepo, mockJWT)

	res, err := svc.Login(context.Background(), "test@example.com", "secret")

	assert.Nil(t, res)
	assert.EqualError(t, err, "jwt error")
	mockRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
}
