package auth

import (
	"backend/internal/domain/user"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repository Repository
	tokens     TokenProvider
	hasher     PasswordHasher
}

func NewService(repository Repository, tokens TokenProvider, hasher PasswordHasher) *Service {
	return &Service{
		repository: repository,
		tokens:     tokens,
		hasher:     hasher,
	}
}

func (s *Service) Register(ctx context.Context, u UserCredentials) (*Token, error) {
	if err := u.ValidateRegister(); err != nil {
		return nil, err
	}

	hashedPassword, err := s.hasher.Hash(u.Password)
	if err != nil {
		return nil, fmt.Errorf("Hashing password: %w", err)
	}

	newUser := &user.User{
		ID:             uuid.NewString(),
		Name:           u.Name,
		Email:          u.Email,
		HashedPassword: hashedPassword,
		Role:           user.RolePassenger,
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

func (s *Service) Login(ctx context.Context, u UserCredentials) (*Token, error) {
	existingUser, err := s.repository.GetByEmail(ctx, u.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !s.hasher.Compare(existingUser.HashedPassword, u.Password) {
		return nil, ErrInvalidCredentials
	}

	if existingUser.IsDisabled {
		return nil, ErrAccountDisabled
	}

	tokens, err := s.tokens.GenerateTokens(existingUser)
	if err != nil {
		return nil, fmt.Errorf("Generating tokens: %w", err)
	}

	return &Token{
		Access:  tokens.AccessToken,
		Refresh: tokens.RefreshToken,
	}, nil
}
