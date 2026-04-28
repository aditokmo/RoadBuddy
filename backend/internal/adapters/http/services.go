package http

import (
	"backend/internal/domain/auth"
	"backend/internal/domain/user"
	"log/slog"
)

type Services struct {
	Auth *auth.Service
	User *user.Service
}

func NewServices(authService *auth.Service, UserService *user.Service, logger *slog.Logger) *Services {
	return &Services{
		Auth: authService,
		User: UserService,
	}
}
