package utils

import (
	"bank-backend/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	SecretKey string
	ExpiredIn time.Duration
	RefreshIn time.Duration
}

type JWTGenerator interface {
	Generate(userID string) (string, error)
}

func LoadJWTConfig() *JWTConfig {
	return &JWTConfig{
		SecretKey: config.Env("JWT_SECRET", "supersecret"),
		ExpiredIn: time.Minute * 5,
		RefreshIn: time.Hour * 24 * 7,
	}
}

type defaultJWTGenerator struct{}

func (g *defaultJWTGenerator) Generate(userID string) (string, error) {
	jwtConfig := LoadJWTConfig()

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(jwtConfig.ExpiredIn).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtConfig.SecretKey))
}

func NewJWTGenerator() JWTGenerator {
	return &defaultJWTGenerator{}
}
