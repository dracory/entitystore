package entitystore

import "context"

// EntityAttributeList list all attributes of an entity
func (st *storeImplementation) EntityAttributeList(ctx context.Context, entityID string) (attributes []Attribute, err error) {
	return st.AttributeList(ctx, AttributeQueryOptions{
		EntityID: entityID,
	})
}
