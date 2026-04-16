package user

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound = errors.New("User not found")
)

type User struct {
	ID             string
	Name           string
	Email          string
	HashedPassword string
	Role           Role
	IsVerified     bool
	IsDisabled     bool
	CreatedAt      time.Time
}

type Role string

const (
	RolePassenger Role = "passenger"
	RoleDriver    Role = "driver"
)
