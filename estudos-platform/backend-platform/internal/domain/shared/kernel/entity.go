package kernel

import (
	"time"

	"github.com/google/uuid"
)

// Entity é a interface base de toda entidade de domínio.
type Entity interface {
	ID() uuid.UUID
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

// BaseEntity encapsula identidade e timestamps comuns a todas as entidades.
type BaseEntity struct {
	id        uuid.UUID
	createdAt time.Time
	updatedAt time.Time
}

// NewBaseEntity gera ID e timestamps para uma entidade nova.
func NewBaseEntity() BaseEntity {
	now := time.Now().UTC()
	return BaseEntity{id: uuid.New(), createdAt: now, updatedAt: now}
}

// NewBaseEntityWithID reconstrói entidade a partir do banco (sem gerar ID novo).
func NewBaseEntityWithID(id uuid.UUID, createdAt, updatedAt time.Time) BaseEntity {
	return BaseEntity{id: id, createdAt: createdAt, updatedAt: updatedAt}
}

func (e BaseEntity) ID() uuid.UUID        { return e.id }
func (e BaseEntity) CreatedAt() time.Time { return e.createdAt }
func (e BaseEntity) UpdatedAt() time.Time { return e.updatedAt }

func (e *BaseEntity) Touch() {
	e.updatedAt = time.Now().UTC()
}
