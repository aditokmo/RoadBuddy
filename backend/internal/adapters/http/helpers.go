package http

import (
	"backend/internal/adapters/render"
	"backend/internal/domain/auth"
	"backend/internal/domain/user"
	"errors"
	"log/slog"
	"net/http"
)

func SetAuthCookies(w http.ResponseWriter, access, refresh string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    access,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearAuthCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *AuthHandler) handleAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, auth.ErrInvalidCredentials):
		render.Error(w, http.StatusUnauthorized, "Invalid email or password", "invalid_credentials")

	case errors.Is(err, auth.ErrEmailTaken):
		render.Error(w, http.StatusConflict, "Email address is already used", "email_taken")

	case errors.Is(err, auth.ErrAccountDisabled):
		render.Error(w, http.StatusForbidden, "This account has been disabled", "account_disabled")

	case errors.Is(err, auth.ErrWeakPassword):
		render.Error(w, http.StatusBadRequest,
			"password must be at least 8 characters, include an uppercase letter, a number, and a special character",
			"weak_password")

	case errors.Is(err, auth.ErrInvalidEmail):
		render.Error(w, http.StatusBadRequest, "Email address is not valid", "invalid_email")

	case errors.Is(err, auth.ErrInvalidToken):
		render.Error(w, http.StatusUnauthorized, "Token is not valid", "invalid_token")

	case errors.Is(err, auth.ErrExpiredToken):
		render.Error(w, http.StatusUnauthorized, "Token has expired", "expired_token")

	case errors.Is(err, user.ErrUserNotFound):
		render.Error(w, http.StatusNotFound, "User not found", "user_not_found")

	default:
		h.logger.Error("Unhandled auth error", "error", err)
		render.Error(w, http.StatusInternalServerError, "An unexpected error occurred", "internal_error")
	}
}

func (h *UserHandler) handleUserError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, user.ErrUserNotFound):
		render.Error(w, http.StatusNotFound, "user not found", "user_not_found")

	default:
		h.logger.Error("unhandled user error", slog.String("error", err.Error()))
		render.Error(w, http.StatusInternalServerError, "an unexpected error occurred", "internal_error")
	}
}
