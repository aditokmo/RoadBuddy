package http

import (
	"backend/internal/domain/auth"
	"backend/internal/domain/user"
	"time"
)

type userResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	AvatarURL   string    `json:"avatar_url"`
	Role        user.Role `json:"role"`
	IsVerified  bool      `json:"is_verified"`
	CreatedAt   string    `json:"created_at"`
}

func toUserResponse(u user.User) userResponse {
	return userResponse{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		Role:       u.Role,
		IsVerified: u.IsVerified,
		CreatedAt:  u.CreatedAt.Format(time.RFC3339),
	}
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func toTokenResponse(token *auth.Token) tokenResponse {
	return tokenResponse{
		AccessToken:  token.Access,
		RefreshToken: token.Refresh,
	}
}
