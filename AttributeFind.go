package entitystore

import (
	"context"
	"errors"
)

// AttributeFind finds an entity by ID
func (st *storeImplementation) AttributeFind(ctx context.Context, entityID string, attributeKey string) (*Attribute, error) {
	if entityID == "" {
		return nil, errors.New("entity id cannot be empty")
	}

	if attributeKey == "" {
		return nil, errors.New("attribute key cannot be empty")
	}

	list, err := st.AttributeList(ctx, AttributeQueryOptions{
		EntityID:     entityID,
		AttributeKey: attributeKey,
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
