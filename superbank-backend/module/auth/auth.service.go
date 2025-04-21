package auth

import (
	"bank-backend/module/user"
	"bank-backend/utils"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type authServiceImpl struct {
	userRepo     user.UserRepository
	jwtGenerator utils.JWTGenerator
}

func NewAuthService(userRepo user.UserRepository, jwtGen utils.JWTGenerator) AuthService {
	return &authServiceImpl{
		userRepo:     userRepo,
		jwtGenerator: jwtGen,
	}
}

func (s *authServiceImpl) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.jwtGenerator.Generate(u.ID.String())
	if err != nil {
		return nil, err
	}

	user := UserResponse{
		ID:       u.ID.String(),
		Email:    u.Email,
		Username: u.Username,
	}

	return &LoginResponse{User: user, Token: token}, nil
}
