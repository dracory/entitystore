package entitystore

import (
	"github.com/dracory/dataobject"
	"github.com/dromara/carbon/v2"
)

// == TYPE ===================================================================

// attributeTrashImplementation represents a trashed attribute backed by a map[string]string
type attributeTrashImplementation struct {
	dataobject.DataObject
}

// == INTERFACE COMPLIANCE ===================================================

var _ AttributeTrashInterface = (*attributeTrashImplementation)(nil)

// == CONSTRUCTORS ===========================================================

// NewAttributeTrash creates a new attribute trash record with default values
func NewAttributeTrash() AttributeTrashInterface {
	o := &attributeTrashImplementation{}
	o.SetID(GenerateShortID())
	o.SetEntityID("")
	o.SetAttributeKey("")
	o.SetAttributeValue("")
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetDeletedBy("")
	return o
}

// NewAttributeTrashFromExistingData creates an attribute trash record from a raw data map
func NewAttributeTrashFromExistingData(data map[string]string) AttributeTrashInterface {
	o := &attributeTrashImplementation{}
	o.Hydrate(data)
	return o
}

// == GETTERS & SETTERS ======================================================

func (o *attributeTrashImplementation) EntityID() string {
	return o.Get(COLUMN_ENTITY_ID)
}

func (o *attributeTrashImplementation) SetEntityID(entityID string) AttributeTrashInterface {
	o.Set(COLUMN_ENTITY_ID, entityID)
	return o
}

func (o *attributeTrashImplementation) AttributeKey() string {
	return o.Get(COLUMN_ATTRIBUTE_KEY)
}

func (o *attributeTrashImplementation) SetAttributeKey(key string) AttributeTrashInterface {
	o.Set(COLUMN_ATTRIBUTE_KEY, key)
	return o
}

func (o *attributeTrashImplementation) AttributeValue() string {
	return o.Get(COLUMN_ATTRIBUTE_VALUE)
}

func (o *attributeTrashImplementation) SetAttributeValue(value string) AttributeTrashInterface {
	o.Set(COLUMN_ATTRIBUTE_VALUE, value)
	return o
}

func (o *attributeTrashImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *attributeTrashImplementation) SetCreatedAt(createdAt string) AttributeTrashInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *attributeTrashImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt(), carbon.UTC)
}

func (o *attributeTrashImplementation) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *attributeTrashImplementation) SetUpdatedAt(updatedAt string) AttributeTrashInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *attributeTrashImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt(), carbon.UTC)
}

func (o *attributeTrashImplementation) DeletedAt() string {
	return o.Get(COLUMN_DELETED_AT)
}

func (o *attributeTrashImplementation) SetDeletedAt(deletedAt string) AttributeTrashInterface {
	o.Set(COLUMN_DELETED_AT, deletedAt)
	return o
}

func (o *attributeTrashImplementation) DeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.DeletedAt(), carbon.UTC)
}

func (o *attributeTrashImplementation) DeletedBy() string {
	return o.Get(COLUMN_DELETED_BY)
}

func (o *attributeTrashImplementation) SetDeletedBy(deletedBy string) AttributeTrashInterface {
	o.Set(COLUMN_DELETED_BY, deletedBy)
	return o
}
