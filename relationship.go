package entitystore

import (
	"strconv"

	"github.com/dracory/dataobject"
	"github.com/dromara/carbon/v2"
)

// == TYPE ===================================================================

// relationshipImplementation represents a relationship between entities backed by a map[string]string
type relationshipImplementation struct {
	dataobject.DataObject
}

// == INTERFACE COMPLIANCE ===================================================

var _ RelationshipInterface = (*relationshipImplementation)(nil)

// == CONSTRUCTORS ===========================================================

// NewRelationship creates a new relationship with default values
// Sets default sequence to 0
func NewRelationship() RelationshipInterface {
	o := &relationshipImplementation{}
	o.SetSequence(0) // Default sequence to 0
	return o
}

// NewRelationshipFromExistingData creates a relationship from a raw data map (e.g., from DB rows)
// Used internally when hydrating relationships from database results
func NewRelationshipFromExistingData(data map[string]string) RelationshipInterface {
	o := &relationshipImplementation{}
	o.Hydrate(data)
	return o
}

// == GETTERS & SETTERS ======================================================

func (o *relationshipImplementation) GetID() string {
	return o.ID()
}

func (o *relationshipImplementation) GetEntityID() string {
	return o.Get(COLUMN_ENTITY_ID)
}

func (o *relationshipImplementation) SetEntityID(entityID string) RelationshipInterface {
	o.Set(COLUMN_ENTITY_ID, entityID)
	return o
}

func (o *relationshipImplementation) GetRelatedEntityID() string {
	return o.Get(COLUMN_RELATED_ENTITY_ID)
}

func (o *relationshipImplementation) SetRelatedEntityID(relatedID string) RelationshipInterface {
	o.Set(COLUMN_RELATED_ENTITY_ID, relatedID)
	return o
}

func (o *relationshipImplementation) GetRelationshipType() string {
	return o.Get(COLUMN_RELATIONSHIP_TYPE)
}

func (o *relationshipImplementation) SetRelationshipType(relType string) RelationshipInterface {
	o.Set(COLUMN_RELATIONSHIP_TYPE, relType)
	return o
}

func (o *relationshipImplementation) GetParentID() string {
	return o.Get(COLUMN_PARENT_ID)
}

func (o *relationshipImplementation) SetParentID(parentID string) RelationshipInterface {
	o.Set(COLUMN_PARENT_ID, parentID)
	return o
}

func (o *relationshipImplementation) GetSequence() int {
	val, _ := strconv.Atoi(o.Get(COLUMN_SEQUENCE))
	return val
}

func (o *relationshipImplementation) SetSequence(sequence int) RelationshipInterface {
	o.Set(COLUMN_SEQUENCE, strconv.Itoa(sequence))
	return o
}

func (o *relationshipImplementation) GetMetadata() string {
	return o.Get(COLUMN_METADATA)
}

func (o *relationshipImplementation) SetMetadata(metadata string) RelationshipInterface {
	o.Set(COLUMN_METADATA, metadata)
	return o
}

func (o *relationshipImplementation) GetCreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *relationshipImplementation) SetCreatedAt(createdAt string) RelationshipInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *relationshipImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetCreatedAt(), carbon.UTC)
}
