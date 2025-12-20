package domain

import (
	"fin-flow-api/internal/shared/domain"
)

type User struct {
	domain.Entity
	
	ID        string
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func NewUser(id, firstName, lastName, email, password, createdBy string) *User {
	return &User{
		Entity:    domain.NewEntity(id, createdBy),
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	}
}