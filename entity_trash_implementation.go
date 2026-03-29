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
// Generates a new short ID and sets timestamps to the current time
func NewEntityTrash() EntityTrashInterface {
	o := &entityTrashImplementation{}
	o.SetType("")
	o.SetHandle("")
	o.SetID(GenerateShortID())
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetDeletedBy("")
	return o
}

// NewEntityTrashFromExistingData creates an entity trash record from a raw data map (e.g., from DB rows)
// Used internally when hydrating trashed entities from database results
func NewEntityTrashFromExistingData(data map[string]string) EntityTrashInterface {
	o := &entityTrashImplementation{}
	o.Hydrate(data)
	return o
}

// == GETTERS & SETTERS ======================================================

func (o *entityTrashImplementation) GetID() string {
	return o.ID()
}

func (o *entityTrashImplementation) GetType() string {
	return o.Get(COLUMN_ENTITY_TYPE)
}

func (o *entityTrashImplementation) SetType(entityType string) EntityTrashInterface {
	o.Set(COLUMN_ENTITY_TYPE, entityType)
	return o
}

func (o *entityTrashImplementation) GetHandle() string {
	return o.Get(COLUMN_ENTITY_HANDLE)
}

func (o *entityTrashImplementation) SetHandle(handle string) EntityTrashInterface {
	o.Set(COLUMN_ENTITY_HANDLE, handle)
	return o
}

func (o *entityTrashImplementation) GetCreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *entityTrashImplementation) SetCreatedAt(createdAt string) EntityTrashInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *entityTrashImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetCreatedAt(), carbon.UTC)
}

func (o *entityTrashImplementation) GetUpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *entityTrashImplementation) SetUpdatedAt(updatedAt string) EntityTrashInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *entityTrashImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetUpdatedAt(), carbon.UTC)
}

func (o *entityTrashImplementation) GetDeletedAt() string {
	return o.Get(COLUMN_DELETED_AT)
}

func (o *entityTrashImplementation) SetDeletedAt(deletedAt string) EntityTrashInterface {
	o.Set(COLUMN_DELETED_AT, deletedAt)
	return o
}

func (o *entityTrashImplementation) GetDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetDeletedAt(), carbon.UTC)
}

func (o *entityTrashImplementation) GetDeletedBy() string {
	return o.Get(COLUMN_DELETED_BY)
}

func (o *entityTrashImplementation) SetDeletedBy(deletedBy string) EntityTrashInterface {
	o.Set(COLUMN_DELETED_BY, deletedBy)
	return o
}
