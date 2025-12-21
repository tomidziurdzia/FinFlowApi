package domain

import (
	"testing"
)

func TestNewWallet(t *testing.T) {
	id := "wallet-1"
	userID := "user-1"
	name := "Main Account"
	walletType := WalletTypeBank
	balance := 1000.50
	currency := CurrencyUSD
	createdBy := "system"

	wallet := NewWallet(id, userID, name, walletType, balance, currency, createdBy)

	if wallet == nil {
		t.Fatal("NewWallet returned nil")
	}

	if wallet.ID != id {
		t.Errorf("expected ID %s, got %s", id, wallet.ID)
	}

	if wallet.UserID != userID {
		t.Errorf("expected UserID %s, got %s", userID, wallet.UserID)
	}

	if wallet.Name != name {
		t.Errorf("expected Name %s, got %s", name, wallet.Name)
	}

	if wallet.Type != walletType {
		t.Errorf("expected Type %v, got %v", walletType, wallet.Type)
	}

	if wallet.Balance != balance {
		t.Errorf("expected Balance %f, got %f", balance, wallet.Balance)
	}

	if wallet.Currency != currency {
		t.Errorf("expected Currency %s, got %s", currency, wallet.Currency)
	}

	if wallet.CreatedBy != createdBy {
		t.Errorf("expected CreatedBy %s, got %s", createdBy, wallet.CreatedBy)
	}
}

func TestWalletType_String(t *testing.T) {
	tests := []struct {
		name     string
		walletType WalletType
		expected string
	}{
		{"Bank", WalletTypeBank, "Bank"},
		{"Cash", WalletTypeCash, "Cash"},
		{"CreditCard", WalletTypeCreditCard, "CreditCard"},
		{"DebitCard", WalletTypeDebitCard, "DebitCard"},
		{"Savings", WalletTypeSavings, "Savings"},
		{"Investment", WalletTypeInvestment, "Investment"},
		{"Other", WalletTypeOther, "Other"},
		{"Unknown", WalletType(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.walletType.String()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestWalletType_Value(t *testing.T) {
	tests := []struct {
		name     string
		walletType WalletType
		expected int
	}{
		{"Bank", WalletTypeBank, 0},
		{"Cash", WalletTypeCash, 1},
		{"CreditCard", WalletTypeCreditCard, 2},
		{"DebitCard", WalletTypeDebitCard, 3},
		{"Savings", WalletTypeSavings, 4},
		{"Investment", WalletTypeInvestment, 5},
		{"Other", WalletTypeOther, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.walletType.Value()
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestIsValidWalletType(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		expected bool
	}{
		{"Bank", 0, true},
		{"Cash", 1, true},
		{"CreditCard", 2, true},
		{"DebitCard", 3, true},
		{"Savings", 4, true},
		{"Investment", 5, true},
		{"Other", 6, true},
		{"Invalid negative", -1, false},
		{"Invalid high", 99, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidWalletType(tt.value)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}