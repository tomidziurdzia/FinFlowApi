package http

import "time"

type WalletResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      int       `json:"type"`
	TypeName  string    `json:"type_name"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
}