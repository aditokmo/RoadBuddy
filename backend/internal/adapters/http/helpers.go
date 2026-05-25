package http

import (
	"backend/internal/adapters/render"
	"backend/internal/domain/auth"
	"backend/internal/domain/user"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

func SetAuthCookies(w http.ResponseWriter, access, refresh string, accessExpiry, refreshExpiry time.Time) {
	setCookie(w, "access_token", access, accessExpiry, true)
	setCookie(w, "refresh_token", refresh, refreshExpiry, true)
}

func setCookie(w http.ResponseWriter, name, value string, expiry time.Time, httpOnly bool) {
	maxAge := int(time.Until(expiry).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: httpOnly,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
		Expires:  expiry,
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
		Expires:  time.Unix(0, 0),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
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

	case errors.Is(err, auth.ErrInvalidRefreshToken):
		render.Error(w, http.StatusUnauthorized, "Refresh token is not valid", "invalid_refresh_token")

	case errors.Is(err, auth.ErrExpiredSession):
		render.Error(w, http.StatusUnauthorized, "Refresh token has expired", "refresh_token_expired")

	case errors.Is(err, auth.ErrInvalidSession):
		render.Error(w, http.StatusUnauthorized, "Refresh token is not valid", "invalid_refresh_token")

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
