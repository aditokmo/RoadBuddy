package http

import (
	"backend/internal/domain/auth"
	"log/slog"
)

type Handlers struct {
	Auth *AuthHandler
}

func NewHandlers(authService auth.Services, logger *slog.Logger) *Handlers {
	return &Handlers{
		Auth: NewAuthHandler(authService, logger),
	}
}
