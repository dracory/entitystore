package entitystore

import "context"

// AttributeSetString creates a new entity
func (st *storeImplementation) AttributeSetString(ctx context.Context, entityID string, attributeKey string, attributeValue string) error {
	attr, err := st.AttributeFind(ctx, entityID, attributeKey)

	if err != nil {
		return err
	}

	if attr == nil {
		attr, err := st.AttributeCreateWithKeyAndValue(ctx, entityID, attributeKey, attributeValue)
		if err != nil {
			return err
		}
		if attr != nil {
			return nil
		}
		return err
	}

	attr.SetString(attributeValue)

	return st.AttributeUpdate(*attr)
}
