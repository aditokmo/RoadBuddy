package http

import (
	"backend/internal/adapters/render"
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
	var req auth.RegisterInput
	if err := render.Decode(w, r, &req); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			render.Error(w, http.StatusRequestEntityTooLarge,
				"request body must not exceed 1MB",
				"payload_too_large",
			)
			return
		}
		render.Error(w, http.StatusBadRequest,
			"request body contains invalid JSON",
			"invalid_json",
		)
		return
	}

	if err := req.ValidateRegister(); err != nil {
		render.Error(w, http.StatusBadRequest, err.Error(), "validation_failed")
		return
	}

	token, err := h.service.Register(r.Context(), req)
	if err != nil {
		h.logger.Error("Registration failed", "email", req.Email, "error", err)
		h.handleAuthError(w, err)
		return
	}

	SetAuthCookies(w, token.Access, token.Refresh, token.AccessTokenExpiry, token.RefreshTokenExpiry)

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginInput
	if err := render.Decode(w, r, &req); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			render.Error(w, http.StatusRequestEntityTooLarge, "Request body must not exceed 1MB", "payload_too_large")
			return
		}
		render.Error(w, http.StatusBadRequest, "Request body contains invalid JSON", "invalid_json")
		return
	}

	if err := req.ValidateLogin(); err != nil {
		render.Error(w, http.StatusBadRequest, err.Error(), "validation_failed")
		return
	}

	headers := auth.LoginHeaders{
		UserAgent: r.Header.Get("User-Agent"),
		IPAddress: r.RemoteAddr,
	}

	token, err := h.service.Login(r.Context(), req, headers)
	if err != nil {
		h.logger.Error("Login failed", "email", req.Email, "error", err)
		h.handleAuthError(w, err)
		return
	}

	SetAuthCookies(w, token.Access, token.Refresh, token.AccessTokenExpiry, token.RefreshTokenExpiry)

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		render.Error(w, http.StatusUnauthorized, "Missing refresh_token cookie", "missing_refresh_token")
		return
	}

	token, err := h.service.RefreshAccessToken(r.Context(), refreshToken.Value)
	if err != nil {
		h.logger.Error("Token refresh failed", "error", err)
		h.handleAuthError(w, err)
		return
	}

	SetAuthCookies(w, token.Access, token.Refresh, token.AccessTokenExpiry, token.RefreshTokenExpiry)

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		render.Error(w, http.StatusBadRequest, "Verification token is required", "missing_verification_token")
		return
	}

	ctx := r.Context()
	err := h.service.VerifyEmail(ctx, token)
	if err != nil {
		if err == auth.ErrInvalidVerificationToken {
			render.Error(w, http.StatusBadRequest, err.Error(), "invalid_verification_token")
			return
		}

		render.Error(w, http.StatusInternalServerError, "Internal server error", "invalid_verification_token_server_error")
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"message": "Email successfully verified! You can now log in."})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		render.Error(w, http.StatusUnauthorized, "Missing refresh_token", "missing_refresh_token")
		return
	}

	if err := h.service.Logout(r.Context(), refreshToken.Value); err != nil {
		h.handleAuthError(w, err)
		return
	}

	ClearAuthCookies(w)

	w.WriteHeader(http.StatusNoContent)
}
