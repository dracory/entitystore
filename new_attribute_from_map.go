package entitystore

import "github.com/dromara/carbon/v2"

func (st *storeImplementation) NewAttributeFromMap(attributeMap map[string]string) Attribute {
	opts := NewAttributeOptions{}

	if id, exists := attributeMap[COLUMN_ID]; exists {
		opts.ID = id
	}
	if entityID, exists := attributeMap[COLUMN_ENTITY_ID]; exists {
		opts.EntityID = entityID
	}
	if attributeKey, exists := attributeMap[COLUMN_ATTRIBUTE_KEY]; exists {
		opts.AttributeKey = attributeKey
	}
	if attributeValue, exists := attributeMap[COLUMN_ATTRIBUTE_VALUE]; exists {
		opts.AttributeValue = attributeValue
	}
	if createdAt, exists := attributeMap[COLUMN_CREATED_AT]; exists {
		opts.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	}
	if updatedAt, exists := attributeMap[COLUMN_UPDATED_AT]; exists {
		opts.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).StdTime()
	}

	return st.NewAttribute(opts)
}
