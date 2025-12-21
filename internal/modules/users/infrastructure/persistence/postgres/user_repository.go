package postgres

import (
	"context"
	"errors"
	"fmt"

	"fin-flow-api/internal/modules/users/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (id, auth_id, first_name, last_name, email, password, created_at, modified_at, created_by, modified_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.pool.Exec(
		context.Background(),
		query,
		user.ID,
		user.AuthID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.ModifiedAt,
		user.CreatedBy,
		user.ModifiedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(id string) (*domain.User, error) {
	query := `
		SELECT id, auth_id, first_name, last_name, email, password, created_at, modified_at, created_by, modified_by
		FROM users
		WHERE id = $1
	`

	var user domain.User
	var authID *string
	err := r.pool.QueryRow(context.Background(), query, id).Scan(
		&user.ID,
		&authID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.ModifiedAt,
		&user.CreatedBy,
		&user.ModifiedBy,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if authID != nil {
		user.AuthID = *authID
	}

	return &user, nil
}

func (r *Repository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, auth_id, first_name, last_name, email, password, created_at, modified_at, created_by, modified_by
		FROM users
		WHERE email = $1
	`

	var user domain.User
	var authID *string
	err := r.pool.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&authID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.ModifiedAt,
		&user.CreatedBy,
		&user.ModifiedBy,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if authID != nil {
		user.AuthID = *authID
	}

	return &user, nil
}

func (r *Repository) GetByAuthID(authID string) (*domain.User, error) {
	query := `
		SELECT id, auth_id, first_name, last_name, email, password, created_at, modified_at, created_by, modified_by
		FROM users
		WHERE auth_id = $1
	`

	var user domain.User
	var authIDValue *string
	err := r.pool.QueryRow(context.Background(), query, authID).Scan(
		&user.ID,
		&authIDValue,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.ModifiedAt,
		&user.CreatedBy,
		&user.ModifiedBy,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if authIDValue != nil {
		user.AuthID = *authIDValue
	}

	return &user, nil
}

func (r *Repository) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET first_name = $2, last_name = $3, email = $4, modified_at = $5, modified_by = $6
		WHERE id = $1
	`

	result, err := r.pool.Exec(
		context.Background(),
		query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.ModifiedAt,
		user.ModifiedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *Repository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.pool.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *Repository) List() ([]*domain.User, error) {
	query := `
		SELECT id, auth_id, first_name, last_name, email, password, created_at, modified_at, created_by, modified_by
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		var authID *string
		err := rows.Scan(
			&user.ID,
			&authID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.ModifiedAt,
			&user.CreatedBy,
			&user.ModifiedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		if authID != nil {
			user.AuthID = *authID
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate users: %w", err)
	}

	return users, nil
}