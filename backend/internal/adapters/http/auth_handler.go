package http

import (
	"backend/internal/domain/auth"
	"errors"
	"log/slog"
	"net/http"
)

type AuthHandler struct {
	service *auth.Service
	logger  *slog.Logger
}

func NewAuthHandler(service *auth.Service, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		service: service,
		logger:  logger,
	}
}

func (h *AuthHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req auth.UserCredentials
	if err := decodeJson(w, r, &req); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			writeError(w, http.StatusRequestEntityTooLarge,
				"request body must not exceed 1MB",
				"payload_too_large",
			)
			return
		}
		writeError(w, http.StatusBadRequest,
			"request body contains invalid JSON",
			"invalid_json",
		)
		return
	}

	if err := req.ValidateRegister(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "validation_failed")
		return
	}

	token, err := h.service.Register(r.Context(), req)
	if err != nil {
		h.logger.Error("Registration failed", "email", req.Email, "error", err)
		h.handleAuthError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toTokenResponse(token))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.UserCredentials
	if err := decodeJson(w, r, &req); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			writeError(w, http.StatusRequestEntityTooLarge, "Request body must not exceed 1MB", "payload_too_large")
			return
		}
		writeError(w, http.StatusBadRequest, "Request body contains invalid JSON", "invalid_json")
		return
	}

	if err := req.ValidateLogin(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "validation_failed")
		return
	}

	token, err := h.service.Login(r.Context(), req)
	if err != nil {
		h.logger.Error("Login failed", "email", req.Email, "error", err)
		h.handleAuthError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toTokenResponse(token))
}
