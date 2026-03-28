package entitystore

import "context"

// EntityAttributeList lists all attributes of an entity
func (st *storeImplementation) EntityAttributeList(ctx context.Context, entityID string) ([]AttributeInterface, error) {
	return st.AttributeList(ctx, AttributeQueryOptions{
		EntityID: entityID,
	})
}
