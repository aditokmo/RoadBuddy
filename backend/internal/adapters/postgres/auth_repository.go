package postgres

import (
	"backend/internal/domain/auth"
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

func (ar *AuthRepository) Create(ctx context.Context, user *auth.User) error {
	query := "INSERT INTO users (id, name, email, password_hash, role, is_verified, is_disabled, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := ar.db.Exec(ctx, query, user.ID, user.Name, user.Email, user.HashedPassword, user.Role, user.IsVerified, user.IsDisabled, user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return auth.ErrEmailTaken
		}
		return fmt.Errorf("Failed to create user: %w", err)
	}

	return nil
}

func (ar *AuthRepository) GetByEmail(ctx context.Context, email string) (*auth.User, error) {
	query := "SELECT id, name, email, password_hash, role, is_verified, is_disabled, created_at FROM users WHERE email = $1"

	user := &auth.User{}

	err := ar.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.HashedPassword,
		&user.Role,
		&user.IsVerified,
		&user.IsDisabled,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, auth.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}
