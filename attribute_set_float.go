package entitystore

import (
	"context"
	"strconv"
)

// AttributeSetFloat creates a new attribute or updates existing
func (st *storeImplementation) AttributeSetFloat(ctx context.Context, entityID string, attributeKey string, attributeValue float64) error {
	attributeValueAsString := strconv.FormatFloat(attributeValue, 'f', 30, 64)
	return st.AttributeSetString(ctx, entityID, attributeKey, attributeValueAsString)
}
