package auth

import (
	"context"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (*LoginResponse, error)
}
