package entitystore

import (
	"context"
)

// EntityCreateWithType is a shortcut to create an entity by providing only the type.
// The ID and timestamps are auto-assigned.
func (st *storeImplementation) EntityCreateWithType(ctx context.Context, entityType string) (EntityInterface, error) {
	entity := NewEntity()
	entity.SetEntityType(entityType)

	if err := st.EntityCreate(ctx, entity); err != nil {
		return entity, err
	}

	return entity, nil
}
