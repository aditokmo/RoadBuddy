package auth

import (
	"backend/internal/domain/user"
	"context"
)

type Repository interface {
	Create(ctx context.Context, user *user.User) error
	GetByEmail(ctx context.Context, email string) (*user.User, error)
}

type TokenProvider interface {
	GenerateTokens(user *user.User) (*TokenPair, error)
	ValidateAccessToken(token string) (*Claims, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashed string, plain string) bool
}
