package user

import (
	"context"
)

type Repository interface {
	GetAll(ctx context.Context) ([]User, error)
	GetById(ctx context.Context, id string) (User, error)
}
