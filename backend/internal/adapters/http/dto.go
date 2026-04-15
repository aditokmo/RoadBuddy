package http

import "backend/internal/domain/auth"

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
