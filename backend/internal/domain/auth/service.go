package auth

import (
	"backend/internal/domain/user"
	"backend/pkg/crypto"
	"context"
	"errors"
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
	emailPort      EmailPort
}

func NewService(authRepo Repository, userRepo user.Repository, tokens TokenProvider, passwordHasher PasswordHasher, tokenHasher TokenHasher, emailPort EmailPort) *Service {
	return &Service{
		authRepo:       authRepo,
		userRepo:       userRepo,
		tokens:         tokens,
		passwordHasher: passwordHasher,
		tokenHasher:    tokenHasher,
		emailPort:      emailPort,
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

	// Create User
	if err := s.authRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	// Email Verification
	verificationToken, err := crypto.GenerateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("Generate Secure Verification Token: %w", err)
	}

	tokenExpiry := time.Now().UTC().Add(24 * time.Hour)
	err = s.authRepo.SaveVerificationToken(ctx, newUser.ID, verificationToken, tokenExpiry)
	if err != nil {
		return nil, fmt.Errorf("Saving Verification Token: %w", err)
	}

	err = s.emailPort.SendEmailVerification(ctx, newUser.Email, verificationToken)
	if err != nil {
		fmt.Printf("Failed to send verification email, please try again: %v\n", err)
	}

	// Generate Access and Refresh tokens
	tokens, err := s.tokens.GenerateTokens(newUser)
	if err != nil {
		return nil, fmt.Errorf("Generating tokens: %w", err)
	}

	return &Token{
		Access:             tokens.AccessToken,
		Refresh:            tokens.RefreshToken,
		AccessTokenExpiry:  tokens.AccessTokenExpiry,
		RefreshTokenExpiry: tokens.RefreshTokenExpiry,
	}, nil
}

func (s *Service) VerifyEmail(ctx context.Context, token string) error {
	if token == "" {
		return ErrInvalidVerificationToken
	}

	userID, err := s.authRepo.GetUserIDByVerificationToken(ctx, token)
	if err != nil {
		return ErrInvalidVerificationToken
	}

	currentUser, err := s.userRepo.GetById(ctx, userID)
	if err != nil {
		return fmt.Errorf("Fetching user for email verification: %w", err)
	}

	if currentUser.IsEmailVerified {
		_ = s.authRepo.DeleteVerificationToken(ctx, token)
		return nil
	}

	err = s.userRepo.UpdateEmailVerificationStatus(ctx, currentUser.ID, true)
	if err != nil {
		return fmt.Errorf("Updating user verification status: %w", err)
	}

	_ = s.authRepo.DeleteVerificationToken(ctx, token)

	return nil
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

	if !existingUser.IsEmailVerified {
		return nil, ErrEmailNotVerified
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
		Access:             tokens.AccessToken,
		Refresh:            tokens.RefreshToken,
		AccessTokenExpiry:  tokens.AccessTokenExpiry,
		RefreshTokenExpiry: tokens.RefreshTokenExpiry,
	}, nil
}

func (s *Service) RefreshAccessToken(ctx context.Context, rawRefreshToken string) (*Token, error) {
	if rawRefreshToken == "" {
		return nil, ErrInvalidRefreshToken
	}

	hashedToken := s.tokenHasher.HashToken(rawRefreshToken)

	session, err := s.authRepo.GetSessionByRefreshToken(ctx, hashedToken)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, fmt.Errorf("Fetching session: %w", err)
	}

	if session.IsExpired() {
		_ = s.authRepo.DeleteSession(ctx, hashedToken)
		return nil, ErrExpiredSession
	}

	currentUser, err := s.userRepo.GetById(ctx, session.UserID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			_ = s.authRepo.DeleteSession(ctx, hashedToken)
			return nil, ErrInvalidRefreshToken
		}
		return nil, fmt.Errorf("Fetching user for session: %w", err)
	}

	if currentUser.IsDisabled {
		err := s.authRepo.DeleteAllUserSessions(ctx, currentUser.ID)
		if err != nil {
			return nil, fmt.Errorf("Deleting user sessions: %w", err)
		}
		return nil, ErrAccountDisabled
	}

	tokens, err := s.tokens.GenerateTokens(&currentUser)
	if err != nil {
		return nil, fmt.Errorf("Generating tokens: %w", err)
	}

	newTokenHash := s.tokenHasher.HashToken(tokens.RefreshToken)
	newSession := &Session{
		ID:               uuid.NewString(),
		UserID:           currentUser.ID,
		RefreshTokenHash: newTokenHash,
		UserAgent:        session.UserAgent,
		IPAddress:        session.IPAddress,
		ExpiresAt:        tokens.RefreshTokenExpiry,
		CreatedAt:        time.Now().UTC(),
	}

	if err := s.authRepo.CreateSession(ctx, newSession); err != nil {
		return nil, fmt.Errorf("Creating new session: %w", err)
	}

	if err := s.authRepo.DeleteSession(ctx, hashedToken); err != nil {
		_ = s.authRepo.DeleteSession(ctx, newTokenHash)
		return nil, fmt.Errorf("Deleting old session: %w", err)
	}

	return &Token{
		Access:             tokens.AccessToken,
		Refresh:            tokens.RefreshToken,
		AccessTokenExpiry:  tokens.AccessTokenExpiry,
		RefreshTokenExpiry: tokens.RefreshTokenExpiry,
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

func (s *Service) ValidateToken(token string) (*JWTPayload, error) {
	return s.tokens.ValidateAccessToken(token)
}
