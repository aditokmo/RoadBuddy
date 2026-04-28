package auth

import (
	"backend/internal/domain/user"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	authRepo       Repository
	userRepo       user.Repository
	tokens         TokenProvider
	passwordHasher PasswordHasher
	tokenHasher    TokenHasher
}

func NewService(authRepo Repository, tokens TokenProvider, passwordHasher PasswordHasher, tokenHasher TokenHasher) *Service {
	return &Service{
		authRepo:       authRepo,
		tokens:         tokens,
		passwordHasher: passwordHasher,
		tokenHasher:    tokenHasher,
	}
}

func (s *Service) Register(ctx context.Context, u RegisterInput) (*Token, error) {
	if err := u.ValidateRegister(); err != nil {
		return nil, err
	}

	hashedPassword, err := s.passwordHasher.HashPassword(u.Password)
	if err != nil {
		return nil, fmt.Errorf("Hashing password: %w", err)
	}

	newUser := &user.User{
		ID:              uuid.NewString(),
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		Email:           u.Email,
		HashedPassword:  hashedPassword,
		Role:            user.RolePassenger,
		IsEmailVerified: false,
		IsPhoneVerified: false,
		IsIDVerified:    false,
		IsDisabled:      false,
		CreatedAt:       time.Now().UTC(),
	}

	if err := s.authRepo.Create(ctx, newUser); err != nil {
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

func (s *Service) Login(ctx context.Context, u LoginInput, headers LoginHeaders) (*Token, error) {
	existingUser, err := s.authRepo.GetByEmail(ctx, u.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !s.passwordHasher.Compare(existingUser.HashedPassword, u.Password) {
		return nil, ErrInvalidCredentials
	}

	if existingUser.IsDisabled {
		return nil, ErrAccountDisabled
	}

	tokens, err := s.tokens.GenerateTokens(existingUser)
	if err != nil {
		return nil, fmt.Errorf("Generating tokens: %w", err)
	}

	newSession := &Session{
		ID:               uuid.NewString(),
		UserID:           existingUser.ID,
		RefreshTokenHash: s.tokenHasher.HashToken(tokens.RefreshToken),
		UserAgent:        headers.UserAgent,
		IPAddress:        headers.IPAddress,
		ExpiresAt:        tokens.RefreshTokenExpiry,
		CreatedAt:        time.Now().UTC(),
	}

	err = s.authRepo.CreateSession(ctx, newSession)
	if err != nil {
		return nil, fmt.Errorf("Creating session: %w", err)
	}

	return &Token{
		Access:  tokens.AccessToken,
		Refresh: tokens.RefreshToken,
	}, nil
}

func (s *Service) RefreshAccessToken(ctx context.Context, rawRefreshToken string) (*Token, error) {
	hashedToken := s.tokenHasher.HashToken(rawRefreshToken)

	session, err := s.authRepo.GetSessionByRefreshToken(ctx, hashedToken)
	if err != nil {
		return nil, ErrInvalidSession
	}

	if session.IsExpired() {
		_ = s.authRepo.DeleteSession(ctx, hashedToken)
		return nil, ErrExpiredSession
	}

	user, err := s.userRepo.GetById(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("Fetching user for session: %w", err)
	}

	// If account is disabled delete all users sessions for all devices
	if user.IsDisabled {
		err := s.authRepo.DeleteAllUserSessions(ctx, user.ID)
		if err != nil {
			return nil, fmt.Errorf("Deleting user sessions: %w", err)
		}
		return nil, ErrAccountDisabled
	}

	tokens, err := s.tokens.GenerateTokens(&user)
	if err != nil {
		return nil, fmt.Errorf("Generating tokens: %w", err)
	}

	if err := s.authRepo.DeleteSession(ctx, hashedToken); err != nil {
		return nil, fmt.Errorf("Deleting old session: %w", err)
	}

	newTokenHash := s.tokenHasher.HashToken(tokens.RefreshToken)
	newSession := &Session{
		ID:               uuid.NewString(),
		UserID:           user.ID,
		RefreshTokenHash: newTokenHash,
		UserAgent:        session.UserAgent,
		IPAddress:        session.IPAddress,
		ExpiresAt:        tokens.RefreshTokenExpiry,
		CreatedAt:        time.Now().UTC(),
	}

	if err := s.authRepo.CreateSession(ctx, newSession); err != nil {
		return nil, fmt.Errorf("Creating new session: %w", err)
	}

	return &Token{
		Access:  tokens.AccessToken,
		Refresh: tokens.RefreshToken,
	}, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	hashedToken := s.tokenHasher.HashToken(refreshToken)

	if err := s.authRepo.DeleteSession(ctx, hashedToken); err != nil {
		if err == ErrSessionNotFound {
			return nil
		}

		return fmt.Errorf("Deleting refresh token: %w", err)
	}

	return nil
}

func (s *Service) ValidateToken(token string) (*Claims, error) {
	return s.tokens.ValidateAccessToken(token)
}
