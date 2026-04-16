package http

import (
	"backend/internal/domain/auth"
	"backend/internal/domain/user"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func decodeJson(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()
	return decode.Decode(dst)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		slog.Error("Failed to marshal JSON response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonData)
}

func writeError(w http.ResponseWriter, status int, message, code string) {
	writeJSON(w, status, errorResponse{Error: message, Code: code})
}

func (h *AuthHandler) handleAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, auth.ErrInvalidCredentials):
		writeError(w, http.StatusUnauthorized, "Invalid email or password", "invalid_credentials")

	case errors.Is(err, auth.ErrEmailTaken):
		writeError(w, http.StatusConflict, "Email address is already used", "email_taken")

	case errors.Is(err, auth.ErrAccountDisabled):
		writeError(w, http.StatusForbidden, "This account has been disabled", "account_disabled")

	case errors.Is(err, auth.ErrWeakPassword):
		writeError(w, http.StatusBadRequest,
			"password must be at least 8 characters, include an uppercase letter, a number, and a special character",
			"weak_password")

	case errors.Is(err, auth.ErrInvalidEmail):
		writeError(w, http.StatusBadRequest, "Email address is not valid", "invalid_email")

	case errors.Is(err, auth.ErrInvalidToken):
		writeError(w, http.StatusUnauthorized, "Token is not valid", "invalid_token")

	case errors.Is(err, auth.ErrExpiredToken):
		writeError(w, http.StatusUnauthorized, "Token has expired", "expired_token")

	case errors.Is(err, user.ErrUserNotFound):
		writeError(w, http.StatusNotFound, "User not found", "user_not_found")

	default:
		h.logger.Error("Unhandled auth error", "error", err)
		writeError(w, http.StatusInternalServerError, "An unexpected error occurred", "internal_error")
	}
}

func (h *UserHandler) handleUserError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, user.ErrUserNotFound):
		writeError(w, http.StatusNotFound, "user not found", "user_not_found")

	default:
		h.logger.Error("unhandled user error", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "an unexpected error occurred", "internal_error")
	}
}
