package middleware

import (
	"backend/internal/adapters/render"
	"backend/internal/domain/auth"
	"context"
	"errors"
	"net/http"
)

type claimsKey struct{}

func JWT(service *auth.Service) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := extractAccessToken(r)
			if tokenString == "" {
				render.Error(w, http.StatusUnauthorized,
					"Access token is required", "missing_token")
				return
			}

			claims, err := service.ValidateToken(tokenString)
			if err != nil {
				switch {
				case errors.Is(err, auth.ErrExpiredToken):
					render.Error(w, http.StatusUnauthorized,
						"Access token has expired", "token_expired")
				default:
					render.Error(w, http.StatusUnauthorized,
						"Invalid access token", "token_invalid")
				}
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractAccessToken(r *http.Request) string {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}
