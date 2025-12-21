package http

import "time"

type CategoryResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      int       `json:"type"`
	TypeName  string    `json:"type_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
}