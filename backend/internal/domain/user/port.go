package user

import (
	"context"
)

type Service interface {
	GetUsers(ctx context.Context) ([]User, error)
	GetUserById(ctx context.Context, id string) (User, error)
}

type Repository interface {
	GetAll(ctx context.Context) ([]User, error)
	GetById(ctx context.Context, id string) (User, error)
}
