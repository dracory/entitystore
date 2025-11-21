package entitystore

import (
	"context"
	"errors"
)

// EntityFindByHandle finds an entity by handle
func (st *storeImplementation) EntityFindByHandle(ctx context.Context, entityType string, entityHandle string) (*Entity, error) {
	if entityType == "" {
		return nil, errors.New("entity type cannot be empty")
	}

	if entityHandle == "" {
		return nil, errors.New("entity handle cannot be empty")
	}

	list, err := st.EntityList(ctx, EntityQueryOptions{
		EntityType:   entityType,
		EntityHandle: entityHandle,
		Limit:        1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return &list[0], nil
	}

	return nil, nil
}
