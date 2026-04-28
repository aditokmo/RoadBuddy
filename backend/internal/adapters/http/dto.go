package http

import (
	"backend/internal/domain/user"
	"time"
)

type userResponse struct {
	ID              string    `json:"id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Email           string    `json:"email"`
	PhoneNumber     string    `json:"phone_number"`
	ProfileImageURL string    `json:"profile_image_url"`
	RatingAverage   float64   `json:"rating_average"`
	RatingCount     int       `json:"rating_count"`
	Role            user.Role `json:"role"`
	IsEmailVerified bool      `json:"is_email_verified"`
	IsIDVerified    bool      `json:"is_id_verified"`
	IsPhoneVerified bool      `json:"is_phone_verified"`
	IsDisabled      bool      `json:"is_disabled"`
	Version         int       `json:"version"`
	UpdatedAt       string    `json:"updated_at"`
	CreatedAt       string    `json:"created_at"`
}

func toUserResponse(u user.User) userResponse {
	return userResponse{
		ID:              u.ID,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		Email:           u.Email,
		PhoneNumber:     u.PhoneNumber,
		ProfileImageURL: u.ProfileImageURL,
		RatingAverage:   u.RatingAverage,
		RatingCount:     u.RatingCount,
		Role:            u.Role,
		IsEmailVerified: u.IsEmailVerified,
		IsPhoneVerified: u.IsPhoneVerified,
		IsIDVerified:    u.IsIDVerified,
		IsDisabled:      u.IsDisabled,
		Version:         u.Version,
		UpdatedAt:       u.UpdatedAt.Format(time.RFC3339),
		CreatedAt:       u.CreatedAt.Format(time.RFC3339),
	}
}
