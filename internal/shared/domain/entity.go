package domain

import "time"

type Entity struct {
	ID        string
	CreatedAt time.Time
	ModifiedAt time.Time
	CreatedBy string
	ModifiedBy string
}

func NewEntity(id, createdBy string) Entity {
	now := time.Now()
	return Entity{
		ID:        id,
		CreatedAt: now,
		ModifiedAt: now,
		CreatedBy: createdBy,
		ModifiedBy: createdBy,
	}
}

func (e *Entity) UpdateModified(modifiedBy string) {
	e.ModifiedAt = time.Now()
	e.ModifiedBy = modifiedBy
}