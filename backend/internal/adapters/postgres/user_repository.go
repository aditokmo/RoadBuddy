package postgres

import (
	"backend/internal/domain/user"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (ur UserRepository) GetAll(ctx context.Context) ([]user.User, error) {
	query := "SELECT id, name, email, role, is_verified, is_disabled, created_at FROM users ORDER BY created_at DESC"

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
			&u.Name,
			&u.Email,
			&u.Role,
			&u.IsVerified,
			&u.IsDisabled,
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
