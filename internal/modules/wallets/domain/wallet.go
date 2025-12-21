package domain

import (
	"fin-flow-api/internal/shared/domain"
)

type Wallet struct {
	domain.Entity

	ID       string
	UserID   string
	Name     string
	Type     WalletType
	Balance  float64
	Currency Currency
}

func NewWallet(id, userID, name string, walletType WalletType, balance float64, currency Currency, createdBy string) *Wallet {
	return &Wallet{
		Entity:   domain.NewEntity(id, createdBy),
		ID:       id,
		UserID:   userID,
		Name:     name,
		Type:     walletType,
		Balance:  balance,
		Currency: currency,
	}
}