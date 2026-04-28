package user

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound = errors.New("User not found")
)

type User struct {
	ID                       string
	FirstName                string
	LastName                 string
	Email                    string
	HashedPassword           string
	DateOfBirth              time.Time
	PhoneNumber              string
	ProfileImageURL          string
	RatingAverage            float64
	RatingCount              int
	Role                     Role
	IsEmailVerified          bool
	IsPhoneVerified          bool
	IsIDVerified             bool
	IsDisabled               bool
	EmailVerificationToken   string
	EmailVerificationExpiry  time.Time
	PasswordResetToken       string
	PasswordResetTokenExpiry time.Time
	LastSeenAt               time.Time
	Version                  int
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

type Role string

const (
	RolePassenger Role = "passenger"
	RoleDriver    Role = "driver"
)
