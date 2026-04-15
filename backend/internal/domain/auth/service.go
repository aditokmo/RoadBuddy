package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	repository Repository
	tokens     TokenProvider
	hasher     PasswordHasher
}

func NewService(repository Repository, tokens TokenProvider, hasher PasswordHasher) Services {
	return &AuthService{
		repository: repository,
		tokens:     tokens,
		hasher:     hasher,
	}
}

func (s *AuthService) Register(ctx context.Context, user UserCredentials) (*Token, error) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	hashedPassword, err := s.hasher.Hash(user.Password)
	if err != nil {
		return nil, fmt.Errorf("Hashing password: %w", err)
	}

	newUser := &User{
		ID:             uuid.NewString(),
		Name:           user.Name,
		Email:          user.Email,
		HashedPassword: hashedPassword,
		Role:           RolePassenger,
		IsVerified:     false,
		IsDisabled:     false,
		CreatedAt:      time.Now().UTC(),
	}

	if err := s.repository.Create(ctx, newUser); err != nil {
		return nil, err
	}

	tokens, err := s.tokens.GenerateTokens(newUser)
	if err != nil {
		return nil, fmt.Errorf("Generating tokens: %w", err)
	}

	return &Token{
		Access:  tokens.AccessToken,
		Refresh: tokens.RefreshToken,
	}, nil
}
