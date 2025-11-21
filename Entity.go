package entitystore

import (
	"context"
	"log"
	"time"
)

// Entity this is the type for an Entity
type Entity struct {
	id           string
	entityType   string
	entityHandle string
	createdAt    time.Time
	updatedAt    time.Time
	st           *storeImplementation
}

func (e *Entity) ToMap() map[string]any {
	entry := map[string]any{}
	entry[COLUMN_ID] = e.ID()
	entry[COLUMN_ENTITY_TYPE] = e.Type()
	entry[COLUMN_ENTITY_HANDLE] = e.Handle()
	entry[COLUMN_CREATED_AT] = e.CreatedAt()
	entry[COLUMN_UPDATED_AT] = e.UpdatedAt()
	return entry
}

func (e *Entity) ID() string {
	return e.id
}

func (e *Entity) Type() string {
	return e.entityType
}

func (e *Entity) Handle() string {
	return e.entityHandle
}

func (e *Entity) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Entity) UpdatedAt() time.Time {
	return e.updatedAt
}

func (e *Entity) SetID(id string) *Entity {
	e.id = id
	return e
}

func (e *Entity) SetType(entityType string) *Entity {
	e.entityType = entityType
	return e
}

func (e *Entity) SetHandle(handle string) *Entity {
	e.entityHandle = handle
	return e
}

func (e *Entity) SetCreatedAt(createdAt time.Time) *Entity {
	e.createdAt = createdAt
	return e
}

func (e *Entity) SetUpdatedAt(updatedAt time.Time) *Entity {
	e.updatedAt = updatedAt
	return e
}

// GetInt the value of the attribute as string or the default value if it does not exist
func (e *Entity) GetInt(attributeKey string, defaultValue int64) (int64, error) {
	attr, err := e.GetAttribute(attributeKey)

	if err != nil {
		return defaultValue, err
	}

	if attr == nil {
		return defaultValue, nil
	}

	return attr.GetInt()
}

// GetAttribute return specified attribute
func (e *Entity) GetAttribute(attributeKey string) (*Attribute, error) {
	return e.st.AttributeFind(context.Background(), e.ID(), attributeKey)
}

// GetAttributes all the attributes of the entity
func (e *Entity) GetAttributes() ([]Attribute, error) {
	return e.st.EntityAttributeList(context.Background(), e.ID())
}

// GetFloat the value of the attribute as float or the default value if it does not exist
func (e *Entity) GetFloat(attributeKey string, defaultValue float64) (float64, error) {
	attr, err := e.GetAttribute(attributeKey)

	if err != nil {
		if e.st.GetDebug() {
			log.Println(err)
		}
		return defaultValue, err
	}

	if attr == nil {
		return defaultValue, nil
	}

	return attr.GetFloat()
}

// GetString the value of the attribute as string or the default value if it does not exist
func (e *Entity) GetString(attributeKey string, defaultValue string) (string, error) {
	attr, err := e.GetAttribute(attributeKey)

	if err != nil {
		if e.st.GetDebug() {
			log.Println(err)
		}
		return defaultValue, err
	}

	if attr == nil {
		return defaultValue, nil
	}

	return attr.GetString(), nil
}

// SetAll upserts the attributes
func (e *Entity) SetAll(attributes map[string]string) error {
	return e.st.AttributesSet(context.Background(), e.ID(), attributes)
}

// SetFloat sets an attribute with float value
func (e *Entity) SetFloat(attributeKey string, attributeValue float64) error {
	return e.st.AttributeSetFloat(context.Background(), e.ID(), attributeKey, attributeValue)
}

// SetInt sets an attribute with int value
func (e *Entity) SetInt(attributeKey string, attributeValue int64) error {
	return e.st.AttributeSetInt(context.Background(), e.ID(), attributeKey, attributeValue)
}

// SetString sets an attribute with string value
func (e *Entity) SetString(attributeKey string, attributeValue string) error {
	return e.st.AttributeSetString(context.Background(), e.ID(), attributeKey, attributeValue)
}
