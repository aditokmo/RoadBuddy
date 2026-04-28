package postgres

import (
	"backend/internal/domain/auth"
	"backend/internal/domain/user"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

func (ar *AuthRepository) Create(ctx context.Context, u *user.User) error {
	query := `
        INSERT INTO users (
            id, first_name, last_name, email, password_hash, 
            role, is_email_verified, is_phone_verified, is_id_verified, 
            is_disabled, created_at, updated_at, version
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err := ar.db.Exec(ctx, query,
		u.ID,
		u.FirstName,
		u.LastName,
		u.Email,
		u.HashedPassword,
		u.Role,
		u.IsEmailVerified,
		u.IsPhoneVerified,
		u.IsIDVerified,
		u.IsDisabled,
		u.CreatedAt,
		u.UpdatedAt,
		u.Version,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return auth.ErrEmailTaken
		}
		return fmt.Errorf("Failed to create user: %w", err)
	}

	return nil
}

func (ar *AuthRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
        SELECT 
            id, first_name, last_name, email, password_hash, 
            role, is_email_verified, is_disabled, rating_average 
        FROM users 
        WHERE email = $1`

	u := &user.User{}

	err := ar.db.QueryRow(ctx, query, email).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.HashedPassword,
		&u.Role,
		&u.IsEmailVerified,
		&u.IsDisabled,
		&u.RatingAverage,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return u, nil
}

func (ar *AuthRepository) CreateSession(ctx context.Context, s *auth.Session) error {
	query := `
	INSERT INTO sessions (id, user_id, refresh_token_hash, user_agent, ip_address, expires_at, created_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7) 
	ON CONFLICT (user_id, user_agent) 
	DO UPDATE SET
	refresh_token_hash = EXCLUDED.refresh_token_hash,
	ip_address = EXCLUDED.ip_address,
	expires_at = EXCLUDED.expires_at,
	created_at = EXCLUDED.created_at
	`

	_, err := ar.db.Exec(ctx, query,
		s.ID,
		s.UserID,
		s.RefreshTokenHash,
		s.UserAgent,
		s.IPAddress,
		s.ExpiresAt,
		s.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("Failed to create session: %w", err)
	}

	return nil
}

func (ar *AuthRepository) DeleteSession(ctx context.Context, refreshTokenHash string) error {
	query := `DELETE FROM sessions WHERE refresh_token_hash = $1`

	result, err := ar.db.Exec(ctx, query, refreshTokenHash)
	if err != nil {
		return fmt.Errorf("Failed to delete session: %w", err)
	}

	if result.RowsAffected() == 0 {
		return auth.ErrSessionNotFound
	}

	return nil
}

func (ar *AuthRepository) DeleteAllUserSessions(ctx context.Context, userID string) error {
	query := `DELETE FROM sessions WHERE user_id = $1`

	_, err := ar.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("Failed to delete user sessions: %w", err)
	}

	return nil
}

func (ar *AuthRepository) GetSessionByRefreshToken(ctx context.Context, refreshTokenHash string) (*auth.Session, error) {
	query := `SELECT id, user_id, refresh_token_hash, user_agent, ip_address, expires_at, created_at FROM sessions WHERE refresh_token_hash = $1`

	s := &auth.Session{}

	err := ar.db.QueryRow(ctx, query, refreshTokenHash).Scan(
		&s.ID,
		&s.UserID,
		&s.RefreshTokenHash,
		&s.UserAgent,
		&s.IPAddress,
		&s.ExpiresAt,
		&s.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, auth.ErrSessionNotFound
		}
		return nil, err
	}

	return s, nil
}
