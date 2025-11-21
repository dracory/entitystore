package entitystore

import (
	"context"
	"time"

	"github.com/dracory/uid"
)

// AttributeCreateWithKeyAndValue shortcut to create a new attribute
// by providing only the key and value
// NN. The ID will be auto-assigned
func (st *storeImplementation) AttributeCreateWithKeyAndValue(ctx context.Context, entityID string, attributeKey string, attributeValue string) (*Attribute, error) {
	newAttribute := st.NewAttribute(NewAttributeOptions{
		ID:             uid.HumanUid(),
		EntityID:       entityID,
		AttributeKey:   attributeKey,
		AttributeValue: attributeValue,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

	err := st.AttributeCreate(ctx, &newAttribute)

	if err != nil {
		return nil, err
	}

	return &newAttribute, nil
}
