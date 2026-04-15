package auth

import "context"

type Services interface {
	Register(ctx context.Context, user UserCredentials) (*Token, error)
}

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type TokenProvider interface {
	GenerateTokens(user *User) (*TokenPair, error)
	ValidateAccessToken(token string) (*Claims, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashed string, plain string) bool
}
