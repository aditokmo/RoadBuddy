package auth

import (
	"backend/internal/domain/user"
	"errors"
	"net/mail"
	"strings"
	"time"
	"unicode"
)

var (
	ErrInvalidCredentials  = errors.New("Invalid credentials")
	ErrSessionNotFound     = errors.New("Session not found")
	ErrInvalidToken        = errors.New("Invalid token")
	ErrInvalidEmail        = errors.New("Invalid email")
	ErrExpiredToken        = errors.New("Expired token")
	ErrExpiredSession      = errors.New("Session has expired")
	ErrInvalidSession      = errors.New("Invalid session")
	ErrEmailTaken          = errors.New("Email is taken")
	ErrAccountDisabled     = errors.New("Account has been disabled")
	ErrWeakPassword        = errors.New("Weak password")
	ErrInvalidRefreshToken = errors.New("Invalid refresh token")
)

type Session struct {
	ID               string
	UserID           string
	RefreshTokenHash string
	UserAgent        string
	IPAddress        string
	ExpiresAt        time.Time
	CreatedAt        time.Time
}

type RegisterInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginHeaders struct {
	UserAgent string
	IPAddress string
}

type Token struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

type TokenPair struct {
	AccessToken        string
	RefreshToken       string
	AccessTokenExpiry  time.Time
	RefreshTokenExpiry time.Time
}

type Claims struct {
	UserID string
	Email  string
	Role   user.Role
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (user *RegisterInput) ValidateRegister() error {
	if err := user.ValidateCommon(); err != nil {
		return err
	}

	if !isStrongPassword(user.Password) {
		return ErrWeakPassword
	}

	return nil
}

func (user *LoginInput) ValidateLogin() error {
	return user.ValidateCommon()
}

func (user *LoginInput) ValidateCommon() error {
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	return nil
}

func (user *RegisterInput) ValidateCommon() error {
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	if user.Email == "" {
		return errors.New("Email is required")
	}

	if _, err := mail.ParseAddress(user.Email); err != nil {
		return ErrInvalidEmail
	}

	if user.Password == "" {
		return errors.New("Password is required")
	}

	return nil
}

func isStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, r := range password {
		if unicode.IsUpper(r) {
			hasUpper = true
		}
		if unicode.IsDigit(r) {
			hasDigit = true
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			hasSpecial = true
		}
	}

	return hasUpper && hasDigit && hasSpecial
}
