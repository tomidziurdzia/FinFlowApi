package domain

type CategoryRepository interface {
	Create(category *Category) error
	GetByID(id string, userID string) (*Category, error)
	List(userID string) ([]*Category, error)
	Update(category *Category) error
	Delete(id string, userID string) error
}