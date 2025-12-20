package domain

type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByAuthID(authID string) (*User, error)
	Update(user *User) error
	Delete(id string) error
	List() ([]*User, error)
}
