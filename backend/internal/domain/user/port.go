package user

import (
	"context"
)

type Repository interface {
	GetAll(ctx context.Context) ([]User, error)
	GetById(ctx context.Context, id string) (User, error)
	UpdateEmailVerificationStatus(ctx context.Context, userID string, isVerified bool) error
}
