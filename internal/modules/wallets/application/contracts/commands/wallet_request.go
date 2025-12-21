package commands

type WalletRequest struct {
	Name     string  `json:"name"`
	Type     int     `json:"type"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}