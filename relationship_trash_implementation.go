package entitystore

import (
	"strconv"

	"github.com/dracory/dataobject"
	"github.com/dromara/carbon/v2"
)

// == TYPE ===================================================================

// relationshipTrashImplementation represents a trashed relationship backed by a map[string]string
type relationshipTrashImplementation struct {
	dataobject.DataObject
}

// == INTERFACE COMPLIANCE ===================================================

var _ RelationshipTrashInterface = (*relationshipTrashImplementation)(nil)

// == CONSTRUCTORS ===========================================================

// NewRelationshipTrash creates a new relationship trash record with default values
func NewRelationshipTrash() RelationshipTrashInterface {
	o := &relationshipTrashImplementation{}
	o.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetDeletedBy("")
	return o
}

// NewRelationshipTrashFromExistingData creates a relationship trash record from a raw data map
func NewRelationshipTrashFromExistingData(data map[string]string) RelationshipTrashInterface {
	o := &relationshipTrashImplementation{}
	o.Hydrate(data)
	return o
}

// == GETTERS & SETTERS ======================================================

func (o *relationshipTrashImplementation) GetEntityID() string {
	return o.Get(COLUMN_ENTITY_ID)
}

func (o *relationshipTrashImplementation) SetEntityID(entityID string) RelationshipTrashInterface {
	o.Set(COLUMN_ENTITY_ID, entityID)
	return o
}

func (o *relationshipTrashImplementation) GetRelatedEntityID() string {
	return o.Get(COLUMN_RELATED_ENTITY_ID)
}

func (o *relationshipTrashImplementation) SetRelatedEntityID(relatedID string) RelationshipTrashInterface {
	o.Set(COLUMN_RELATED_ENTITY_ID, relatedID)
	return o
}

func (o *relationshipTrashImplementation) GetRelationshipType() string {
	return o.Get(COLUMN_RELATIONSHIP_TYPE)
}

func (o *relationshipTrashImplementation) SetRelationshipType(relType string) RelationshipTrashInterface {
	o.Set(COLUMN_RELATIONSHIP_TYPE, relType)
	return o
}

func (o *relationshipTrashImplementation) GetParentID() string {
	return o.Get(COLUMN_PARENT_ID)
}

func (o *relationshipTrashImplementation) SetParentID(parentID string) RelationshipTrashInterface {
	o.Set(COLUMN_PARENT_ID, parentID)
	return o
}

func (o *relationshipTrashImplementation) GetSequence() int {
	val, _ := strconv.Atoi(o.Get(COLUMN_SEQUENCE))
	return val
}

func (o *relationshipTrashImplementation) SetSequence(sequence int) RelationshipTrashInterface {
	o.Set(COLUMN_SEQUENCE, strconv.Itoa(sequence))
	return o
}

func (o *relationshipTrashImplementation) GetMetadata() string {
	return o.Get(COLUMN_METADATA)
}

func (o *relationshipTrashImplementation) SetMetadata(metadata string) RelationshipTrashInterface {
	o.Set(COLUMN_METADATA, metadata)
	return o
}

func (o *relationshipTrashImplementation) GetCreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *relationshipTrashImplementation) SetCreatedAt(createdAt string) RelationshipTrashInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *relationshipTrashImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetCreatedAt(), carbon.UTC)
}

func (o *relationshipTrashImplementation) GetDeletedAt() string {
	return o.Get(COLUMN_DELETED_AT)
}

func (o *relationshipTrashImplementation) SetDeletedAt(deletedAt string) RelationshipTrashInterface {
	o.Set(COLUMN_DELETED_AT, deletedAt)
	return o
}

func (o *relationshipTrashImplementation) GetDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetDeletedAt(), carbon.UTC)
}

func (o *relationshipTrashImplementation) GetDeletedBy() string {
	return o.Get(COLUMN_DELETED_BY)
}

func (o *relationshipTrashImplementation) SetDeletedBy(deletedBy string) RelationshipTrashInterface {
	o.Set(COLUMN_DELETED_BY, deletedBy)
	return o
}
