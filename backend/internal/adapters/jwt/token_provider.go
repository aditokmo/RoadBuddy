package jwt

import (
	"backend/internal/domain/auth"
	"backend/internal/domain/user"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenProvider struct {
	secret          []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewTokenProvider(secret string, accessTTL, refreshTTL time.Duration) (*TokenProvider, error) {
	if secret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}

	return &TokenProvider{
		secret:          []byte(secret),
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}, nil
}

func (p *TokenProvider) GenerateTokens(user *user.User) (*auth.TokenPair, error) {
	now := time.Now().UTC()
	accessExpiry := now.Add(p.accessTokenTTL)
	refreshExpiry := now.Add(p.refreshTokenTTL)

	accessToken, err := p.signToken(user, accessExpiry, "access")
	if err != nil {
		return nil, err
	}

	refreshToken, err := p.signToken(user, refreshExpiry, "refresh")
	if err != nil {
		return nil, err
	}

	return &auth.TokenPair{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		AccessTokenExpiry:  accessExpiry,
		RefreshTokenExpiry: refreshExpiry,
	}, nil
}

func (p *TokenProvider) ValidateAccessToken(token string) (*auth.Claims, error) {
	parsedToken, err := jwt.Parse(token, func(parsedToken *jwt.Token) (interface{}, error) {
		if _, ok := parsedToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrInvalidToken
		}
		return p.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, auth.ErrExpiredToken
		}
		return nil, auth.ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, auth.ErrInvalidToken
	}

	tokenType, _ := claims["type"].(string)
	if tokenType != "access" {
		return nil, auth.ErrInvalidToken
	}

	userID, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)
	roleValue, _ := claims["role"].(string)

	return &auth.Claims{
		UserID: userID,
		Email:  email,
		Role:   user.Role(roleValue),
	}, nil
}

func (p *TokenProvider) signToken(user *user.User, expiresAt time.Time, tokenType string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  string(user.Role),
		"type":  tokenType,
		"exp":   expiresAt.Unix(),
		"iat":   time.Now().UTC().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.secret)
}
