package domain

import (
	"fin-flow-api/internal/shared/domain"
)

type Category struct {
	domain.Entity

	ID         string
	UserID     string
	Name       string
	Type       CategoryType
}

func NewCategory(id, userID, name string, categoryType CategoryType, createdBy string) *Category {
	return &Category{
		Entity: domain.NewEntity(id, createdBy),
		ID:     id,
		UserID: userID,
		Name:   name,
		Type:   categoryType,
	}
}