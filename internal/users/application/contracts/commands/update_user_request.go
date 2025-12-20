package commands

// UpdateUserRequest represents the contract for updating a user
type UpdateUserRequest struct {
	UserID    string
	FirstName string
	LastName  string
	Email     string
}

