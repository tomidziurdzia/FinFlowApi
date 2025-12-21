package http

type CurrenciesResponse struct {
	Fiat   []string `json:"fiat"`
	Crypto []string `json:"crypto"`
}