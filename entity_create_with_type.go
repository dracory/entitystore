package entitystore

import (
	"context"
	"time"

	"github.com/dracory/uid"
)

// EntityCreateWithType quick shortcut method
// to create an entity by providing only the type
// NB. The ID will be auto-assigned
func (st *storeImplementation) EntityCreateWithType(ctx context.Context, entityType string) (*Entity, error) {
	entity := st.NewEntity(NewEntityOptions{
		ID:        uid.HumanUid(),
		Type:      entityType,
		Handle:    "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	err := st.EntityCreate(ctx, &entity)

	if err != nil {
		return &entity, err
	}

	return &entity, nil
}
