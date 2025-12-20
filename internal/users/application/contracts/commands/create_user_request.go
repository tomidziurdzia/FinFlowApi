package commands

// CreateUserRequest represents the contract for creating a user
type CreateUserRequest struct {
	AuthID    string
	FirstName string
	LastName  string
	Email     string
	Password  string
}

