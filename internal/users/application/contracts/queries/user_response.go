package queries

import "time"

// UserResponse represents the contract for user data in responses
type UserResponse struct {
	ID        string
	AuthID    string
	FirstName string
	LastName  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy string
	UpdatedBy string
}

