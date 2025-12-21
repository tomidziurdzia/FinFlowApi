package postgres

import (
	"context"
	"errors"
	"fmt"

	"fin-flow-api/internal/modules/categories/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(category *domain.Category) error {
	query := `
		INSERT INTO categories (id, user_id, name, type, created_at, modified_at, created_by, modified_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.pool.Exec(
		context.Background(),
		query,
		category.ID,
		category.UserID,
		category.Name,
		category.Type.Value(),
		category.CreatedAt,
		category.ModifiedAt,
		category.CreatedBy,
		category.ModifiedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(id string, userID string) (*domain.Category, error) {
	query := `
		SELECT id, user_id, name, type, created_at, modified_at, created_by, modified_by
		FROM categories
		WHERE id = $1 AND user_id = $2
	`

	var category domain.Category
	var typeValue int
	err := r.pool.QueryRow(context.Background(), query, id, userID).Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&typeValue,
		&category.CreatedAt,
		&category.ModifiedAt,
		&category.CreatedBy,
		&category.ModifiedBy,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	category.Type = domain.CategoryType(typeValue)

	return &category, nil
}

func (r *Repository) List(userID string) ([]*domain.Category, error) {
	query := `
		SELECT id, user_id, name, type, created_at, modified_at, created_by, modified_by
		FROM categories
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var category domain.Category
		var typeValue int
		err := rows.Scan(
			&category.ID,
			&category.UserID,
			&category.Name,
			&typeValue,
			&category.CreatedAt,
			&category.ModifiedAt,
			&category.CreatedBy,
			&category.ModifiedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		category.Type = domain.CategoryType(typeValue)
		categories = append(categories, &category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate categories: %w", err)
	}

	return categories, nil
}

func (r *Repository) Update(category *domain.Category) error {
	query := `
		UPDATE categories
		SET name = $2, type = $3, modified_at = $4, modified_by = $5
		WHERE id = $1 AND user_id = $6
	`

	result, err := r.pool.Exec(
		context.Background(),
		query,
		category.ID,
		category.Name,
		category.Type.Value(),
		category.ModifiedAt,
		category.ModifiedBy,
		category.UserID,
	)

	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

func (r *Repository) Delete(id string, userID string) error {
	query := `DELETE FROM categories WHERE id = $1 AND user_id = $2`

	result, err := r.pool.Exec(context.Background(), query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}