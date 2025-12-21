package queries

import "time"

type CategoryResponse struct {
	ID        string
	Name      string
	Type      int
	TypeName  string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy string
	UpdatedBy string
}