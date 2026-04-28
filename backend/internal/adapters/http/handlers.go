package http

import (
	"backend/internal/domain/auth"
	"backend/internal/domain/user"
	"log/slog"
)

type Handlers struct {
	Auth *AuthHandler
	User *UserHandler
}

func NewHandlers(authService *auth.Service, UserService *user.Service, logger *slog.Logger) *Handlers {
	return &Handlers{
		Auth: NewAuthHandler(authService, logger),
		User: NewUserHandler(UserService, logger),
	}
}
