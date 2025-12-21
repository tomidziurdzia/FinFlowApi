package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"fin-flow-api/internal/modules/wallets/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(wallet *domain.Wallet) error {
	query := `
		INSERT INTO wallets (id, user_id, name, type, balance, currency, created_at, modified_at, created_by, modified_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.pool.Exec(
		context.Background(),
		query,
		wallet.ID,
		wallet.UserID,
		wallet.Name,
		wallet.Type.Value(),
		wallet.Balance,
		wallet.Currency.String(),
		wallet.CreatedAt,
		wallet.ModifiedAt,
		wallet.CreatedBy,
		wallet.ModifiedBy,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pgErr.ConstraintName, "name") {
					return fmt.Errorf("wallet name already exists")
				}
				return fmt.Errorf("duplicate entry")
			case "23503": // foreign_key_violation
				return fmt.Errorf("invalid user reference")
			}
		}
		return fmt.Errorf("failed to create wallet: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(id string, userID string) (*domain.Wallet, error) {
	checkQuery := `SELECT user_id FROM wallets WHERE id = $1`
	var walletUserID string
	err := r.pool.QueryRow(context.Background(), checkQuery, id).Scan(&walletUserID)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("wallet not found")
		}
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}
	
	if walletUserID != userID {
		return nil, fmt.Errorf("unauthorized access to wallet")
	}

	query := `
		SELECT id, user_id, name, type, balance, currency, created_at, modified_at, created_by, modified_by
		FROM wallets
		WHERE id = $1 AND user_id = $2
	`

	var wallet domain.Wallet
	var typeValue int
	var currencyStr string
	err = r.pool.QueryRow(context.Background(), query, id, userID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Name,
		&typeValue,
		&wallet.Balance,
		&currencyStr,
		&wallet.CreatedAt,
		&wallet.ModifiedAt,
		&wallet.CreatedBy,
		&wallet.ModifiedBy,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("wallet not found")
		}
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	wallet.Type = domain.WalletType(typeValue)
	wallet.Currency = domain.Currency(currencyStr)

	return &wallet, nil
}

func (r *Repository) List(userID string) ([]*domain.Wallet, error) {
	query := `
		SELECT id, user_id, name, type, balance, currency, created_at, modified_at, created_by, modified_by
		FROM wallets
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list wallets: %w", err)
	}
	defer rows.Close()

	var wallets []*domain.Wallet
	for rows.Next() {
		var wallet domain.Wallet
		var typeValue int
		var currencyStr string
		err := rows.Scan(
			&wallet.ID,
			&wallet.UserID,
			&wallet.Name,
			&typeValue,
			&wallet.Balance,
			&currencyStr,
			&wallet.CreatedAt,
			&wallet.ModifiedAt,
			&wallet.CreatedBy,
			&wallet.ModifiedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan wallet: %w", err)
		}
		wallet.Type = domain.WalletType(typeValue)
		wallet.Currency = domain.Currency(currencyStr)
		wallets = append(wallets, &wallet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate wallets: %w", err)
	}

	return wallets, nil
}

func (r *Repository) Update(wallet *domain.Wallet) error {
	checkQuery := `SELECT user_id FROM wallets WHERE id = $1`
	var walletUserID string
	err := r.pool.QueryRow(context.Background(), checkQuery, wallet.ID).Scan(&walletUserID)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("wallet not found")
		}
		return fmt.Errorf("failed to update wallet: %w", err)
	}
	
	if walletUserID != wallet.UserID {
		return fmt.Errorf("unauthorized access to wallet")
	}

	query := `
		UPDATE wallets
		SET name = $2, type = $3, balance = $4, currency = $5, modified_at = $6, modified_by = $7
		WHERE id = $1 AND user_id = $8
	`

	result, err := r.pool.Exec(
		context.Background(),
		query,
		wallet.ID,
		wallet.Name,
		wallet.Type.Value(),
		wallet.Balance,
		wallet.Currency.String(),
		wallet.ModifiedAt,
		wallet.ModifiedBy,
		wallet.UserID,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pgErr.ConstraintName, "name") {
					return fmt.Errorf("wallet name already exists")
				}
				return fmt.Errorf("duplicate entry")
			}
		}
		return fmt.Errorf("failed to update wallet: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("wallet not found")
	}

	return nil
}

func (r *Repository) Delete(id string, userID string) error {
	checkQuery := `SELECT user_id FROM wallets WHERE id = $1`
	var walletUserID string
	err := r.pool.QueryRow(context.Background(), checkQuery, id).Scan(&walletUserID)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("wallet not found")
		}
		return fmt.Errorf("failed to delete wallet: %w", err)
	}
	
	if walletUserID != userID {
		return fmt.Errorf("unauthorized access to wallet")
	}

	query := `DELETE FROM wallets WHERE id = $1 AND user_id = $2`

	result, err := r.pool.Exec(context.Background(), query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete wallet: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("wallet not found")
	}

	return nil
}