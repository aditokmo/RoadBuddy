package user

import (
	"context"
)

type Service interface {
	GetUsers(ctx context.Context) ([]User, error)
}

type Repository interface {
	GetAll(ctx context.Context) ([]User, error)
}
