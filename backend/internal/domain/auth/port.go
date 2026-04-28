package auth

import (
	"backend/internal/domain/user"
	"context"
)

type Repository interface {
	Create(ctx context.Context, user *user.User) error
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	CreateSession(ctx context.Context, s *Session) error
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*Session, error)
	DeleteSession(ctx context.Context, refreshToken string) error
	DeleteAllUserSessions(ctx context.Context, userID string) error
}

type TokenProvider interface {
	GenerateTokens(user *user.User) (*TokenPair, error)
	ValidateAccessToken(token string) (*Claims, error)
}

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	Compare(hashed, plain string) bool
}

type TokenHasher interface {
	HashToken(token string) string
}
