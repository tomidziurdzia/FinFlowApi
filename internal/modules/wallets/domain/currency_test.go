package domain

import (
	"testing"
)

func TestCurrency_String(t *testing.T) {
	tests := []struct {
		name     string
		currency Currency
		expected string
	}{
		// Fiat currencies
		{"USD", CurrencyUSD, "USD"},
		{"EUR", CurrencyEUR, "EUR"},
		{"ARS", CurrencyARS, "ARS"},
		{"ARG", CurrencyARG, "ARG"},
		{"BRL", CurrencyBRL, "BRL"},
		{"MXN", CurrencyMXN, "MXN"},
		{"GBP", CurrencyGBP, "GBP"},
		{"JPY", CurrencyJPY, "JPY"},
		{"CAD", CurrencyCAD, "CAD"},
		{"AUD", CurrencyAUD, "AUD"},
		{"CHF", CurrencyCHF, "CHF"},
		{"DKK", CurrencyDKK, "DKK"},
		// Cryptocurrencies
		{"BTC", CurrencyBTC, "BTC"},
		{"ETH", CurrencyETH, "ETH"},
		{"USDT", CurrencyUSDT, "USDT"},
		{"USDC", CurrencyUSDC, "USDC"},
		{"NEXO", CurrencyNEXO, "NEXO"},
		{"BNB", CurrencyBNB, "BNB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.currency.String()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestIsValidCurrency(t *testing.T) {
	tests := []struct {
		name     string
		currency string
		expected bool
	}{
		// Fiat currencies
		{"USD", "USD", true},
		{"EUR", "EUR", true},
		{"ARS", "ARS", true},
		{"ARG", "ARG", true},
		{"BRL", "BRL", true},
		{"MXN", "MXN", true},
		{"GBP", "GBP", true},
		{"JPY", "JPY", true},
		{"CAD", "CAD", true},
		{"AUD", "AUD", true},
		{"CHF", "CHF", true},
		{"DKK", "DKK", true},
		// Cryptocurrencies
		{"BTC", "BTC", true},
		{"ETH", "ETH", true},
		{"USDT", "USDT", true},
		{"USDC", "USDC", true},
		{"NEXO", "NEXO", true},
		{"BNB", "BNB", true},
		// Invalid
		{"Invalid", "INVALID", false},
		{"Empty", "", false},
		{"Lowercase", "usd", false},
		{"Lowercase crypto", "btc", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidCurrency(tt.currency)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}