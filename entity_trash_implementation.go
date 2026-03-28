package entitystore

import (
	"github.com/dracory/dataobject"
	"github.com/dromara/carbon/v2"
)

// == TYPE ===================================================================

// entityTrashImplementation represents a trashed entity backed by a map[string]string
type entityTrashImplementation struct {
	dataobject.DataObject
}

// == INTERFACE COMPLIANCE ===================================================

var _ EntityTrashInterface = (*entityTrashImplementation)(nil)

// == CONSTRUCTORS ===========================================================

// NewEntityTrash creates a new entity trash record with default values
func NewEntityTrash() EntityTrashInterface {
	o := &entityTrashImplementation{}
	o.SetEntityType("")
	o.SetEntityHandle("")
	o.SetID(GenerateShortID())
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetDeletedBy("")
	return o
}

// NewEntityTrashFromExistingData creates an entity trash record from a raw data map
func NewEntityTrashFromExistingData(data map[string]string) EntityTrashInterface {
	o := &entityTrashImplementation{}
	o.Hydrate(data)
	return o
}

// == GETTERS & SETTERS ======================================================

func (o *entityTrashImplementation) EntityType() string {
	return o.Get(COLUMN_ENTITY_TYPE)
}

func (o *entityTrashImplementation) SetEntityType(entityType string) EntityTrashInterface {
	o.Set(COLUMN_ENTITY_TYPE, entityType)
	return o
}

func (o *entityTrashImplementation) EntityHandle() string {
	return o.Get(COLUMN_ENTITY_HANDLE)
}

func (o *entityTrashImplementation) SetEntityHandle(handle string) EntityTrashInterface {
	o.Set(COLUMN_ENTITY_HANDLE, handle)
	return o
}

func (o *entityTrashImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *entityTrashImplementation) SetCreatedAt(createdAt string) EntityTrashInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *entityTrashImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt(), carbon.UTC)
}

func (o *entityTrashImplementation) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *entityTrashImplementation) SetUpdatedAt(updatedAt string) EntityTrashInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *entityTrashImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt(), carbon.UTC)
}

func (o *entityTrashImplementation) DeletedAt() string {
	return o.Get(COLUMN_DELETED_AT)
}

func (o *entityTrashImplementation) SetDeletedAt(deletedAt string) EntityTrashInterface {
	o.Set(COLUMN_DELETED_AT, deletedAt)
	return o
}

func (o *entityTrashImplementation) DeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.DeletedAt(), carbon.UTC)
}

func (o *entityTrashImplementation) DeletedBy() string {
	return o.Get(COLUMN_DELETED_BY)
}

func (o *entityTrashImplementation) SetDeletedBy(deletedBy string) EntityTrashInterface {
	o.Set(COLUMN_DELETED_BY, deletedBy)
	return o
}
