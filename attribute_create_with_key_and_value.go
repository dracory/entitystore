package entitystore

import (
	"context"
)

// AttributeCreateWithKeyAndValue creates a new attribute with the given key and value.
// The ID and timestamps are auto-assigned.
func (st *storeImplementation) AttributeCreateWithKeyAndValue(ctx context.Context, entityID string, attributeKey string, attributeValue string) (AttributeInterface, error) {
	attr := NewAttribute()
	attr.SetEntityID(entityID)
	attr.SetAttributeKey(attributeKey)
	attr.SetAttributeValue(attributeValue)

	if err := st.AttributeCreate(ctx, attr); err != nil {
		return nil, err
	}

	return attr, nil
}
