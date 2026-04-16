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
	ErrInvalidCredentials = errors.New("Invalid credentials")
	ErrInvalidToken       = errors.New("Invalid token")
	ErrInvalidEmail       = errors.New("Invalid email")
	ErrExpiredToken       = errors.New("Expired token")
	ErrEmailTaken         = errors.New("Email is taken")
	ErrAccountDisabled    = errors.New("Account has been disabled")
	ErrWeakPassword       = errors.New("Weak password")
)

type UserCredentials struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
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

func (user *UserCredentials) Validate() error {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	if strings.TrimSpace(user.Name) == "" {
		return errors.New("name is required")
	}

	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email is required")
	}

	if _, err := mail.ParseAddress(user.Email); err != nil {
		return ErrInvalidEmail
	}

	if !isStrongPassword(user.Password) {
		return ErrWeakPassword
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
