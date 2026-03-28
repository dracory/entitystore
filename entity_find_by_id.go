package entitystore

import (
	"context"
	"errors"
)

// EntityFindByID finds an entity by its ID
func (st *storeImplementation) EntityFindByID(ctx context.Context, entityID string) (EntityInterface, error) {
	if entityID == "" {
		return nil, errors.New("entity ID cannot be empty")
	}

	list, err := st.EntityList(ctx, EntityQueryOptions{
		ID:    entityID,
		Limit: 1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}
