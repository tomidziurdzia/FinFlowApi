package domain

import "errors"

type Currency string

const (
	// Fiat currencies
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyARS Currency = "ARS"
	CurrencyARG Currency = "ARG"
	CurrencyBRL Currency = "BRL"
	CurrencyMXN Currency = "MXN"
	CurrencyGBP Currency = "GBP"
	CurrencyJPY Currency = "JPY"
	CurrencyCAD Currency = "CAD"
	CurrencyAUD Currency = "AUD"
	CurrencyCHF Currency = "CHF"
	CurrencyDKK Currency = "DKK"
	
	// Cryptocurrencies
	CurrencyBTC  Currency = "BTC"
	CurrencyETH  Currency = "ETH"
	CurrencyUSDT Currency = "USDT"
	CurrencyUSDC Currency = "USDC"
	CurrencyNEXO Currency = "NEXO"
	CurrencyBNB  Currency = "BNB"
)

var ErrInvalidCurrency = errors.New("invalid currency")

func (c Currency) String() string {
	return string(c)
}

func IsValidCurrency(currency string) bool {
	validCurrencies := map[string]bool{
		// Fiat currencies
		string(CurrencyUSD): true,
		string(CurrencyEUR): true,
		string(CurrencyARS): true,
		string(CurrencyARG): true,
		string(CurrencyBRL): true,
		string(CurrencyMXN): true,
		string(CurrencyGBP): true,
		string(CurrencyJPY): true,
		string(CurrencyCAD): true,
		string(CurrencyAUD): true,
		string(CurrencyCHF): true,
		string(CurrencyDKK): true,
		// Cryptocurrencies
		string(CurrencyBTC):  true,
		string(CurrencyETH):  true,
		string(CurrencyUSDT): true,
		string(CurrencyUSDC): true,
		string(CurrencyNEXO): true,
		string(CurrencyBNB):  true,
	}
	return validCurrencies[currency]
}

func GetAllCurrencies() map[string][]string {
	return map[string][]string{
		"fiat": {
			string(CurrencyUSD),
			string(CurrencyEUR),
			string(CurrencyARS),
			string(CurrencyARG),
			string(CurrencyBRL),
			string(CurrencyMXN),
			string(CurrencyGBP),
			string(CurrencyJPY),
			string(CurrencyCAD),
			string(CurrencyAUD),
			string(CurrencyCHF),
			string(CurrencyDKK),
		},
		"crypto": {
			string(CurrencyBTC),
			string(CurrencyETH),
			string(CurrencyUSDT),
			string(CurrencyUSDC),
			string(CurrencyNEXO),
			string(CurrencyBNB),
		},
	}
}