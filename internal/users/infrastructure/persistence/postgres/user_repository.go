package postgres

import (
	"context"
	"errors"
	"fmt"

	"fin-flow-api/internal/users/domain"

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
		INSERT INTO users (id, first_name, last_name, email, password, created_at, modified_at, created_by, modified_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.pool.Exec(
		context.Background(),
		query,
		user.ID,
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
		return fmt.Errorf("error al crear usuario: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(id string) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, created_at, modified_at, created_by, modified_by
		FROM users
		WHERE id = $1
	`

	var user domain.User
	err := r.pool.QueryRow(context.Background(), query, id).Scan(
		&user.ID,
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
			return nil, fmt.Errorf("usuario no encontrado")
		}
		return nil, fmt.Errorf("error al obtener usuario: %w", err)
	}

	return &user, nil
}

func (r *Repository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, created_at, modified_at, created_by, modified_by
		FROM users
		WHERE email = $1
	`

	var user domain.User
	err := r.pool.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
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
			return nil, fmt.Errorf("usuario no encontrado")
		}
		return nil, fmt.Errorf("error al obtener usuario: %w", err)
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
		return fmt.Errorf("error al actualizar usuario: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("usuario no encontrado")
	}

	return nil
}

func (r *Repository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.pool.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar usuario: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("usuario no encontrado")
	}

	return nil
}

func (r *Repository) List() ([]*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, created_at, modified_at, created_by, modified_by
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("error al listar usuarios: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
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
			return nil, fmt.Errorf("error al escanear usuario: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar usuarios: %w", err)
	}

	return users, nil
}

