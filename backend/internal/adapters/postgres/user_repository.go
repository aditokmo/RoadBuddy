package postgres

import (
	"backend/internal/domain/user"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) GetAll(ctx context.Context) ([]user.User, error) {
	query := `
        SELECT 
            id, first_name, last_name, email, phone_number, profile_image_url, 
            rating_average, rating_count, role, is_email_verified, is_id_verified, 
            is_disabled, version, updated_at, created_at 
        FROM users 
        ORDER BY created_at DESC`

	rows, err := ur.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Querying users: %w", err)
	}

	defer rows.Close()

	var users []user.User

	for rows.Next() {
		var u user.User
		err := rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.PhoneNumber,
			&u.ProfileImageURL,
			&u.RatingAverage,
			&u.RatingCount,
			&u.Role,
			&u.IsEmailVerified,
			&u.IsIDVerified,
			&u.IsDisabled,
			&u.Version,
			&u.UpdatedAt,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("Scanning user row: %w", err)
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Processing rows: %w", err)
	}

	return users, nil
}

func (ur *UserRepository) GetById(ctx context.Context, id string) (user.User, error) {
	var u user.User
	query := `
        SELECT 
            id, first_name, last_name, email, phone_number, profile_image_url, 
            rating_average, rating_count, role, is_email_verified, is_id_verified, is_phone_verified, 
            is_disabled, version, updated_at, created_at 
        FROM users 
        WHERE id = $1`

	err := ur.db.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.PhoneNumber,
		&u.ProfileImageURL,
		&u.RatingAverage,
		&u.RatingCount,
		&u.Role,
		&u.IsEmailVerified,
		&u.IsPhoneVerified,
		&u.IsIDVerified,
		&u.IsDisabled,
		&u.Version,
		&u.UpdatedAt,
		&u.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return user.User{}, user.ErrUserNotFound
	}
	if err != nil {
		return user.User{}, err
	}

	return u, nil
}
