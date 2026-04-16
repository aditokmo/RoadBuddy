package http

import (
	"backend/internal/domain/user"
	"log/slog"
	"net/http"
)

type UserHandler struct {
	service user.Service
	logger  *slog.Logger
}

func NewUserHandler(service user.Service, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

func (h UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers(r.Context())
	if err != nil {
		h.handleUserError(w, err)
		return
	}

	response := make([]userResponse, 0, len(users))
	for _, u := range users {
		response = append(response, toUserResponse(&u))
	}

	writeJSON(w, http.StatusOK, response)
}
