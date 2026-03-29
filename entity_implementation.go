package entitystore

import (
	"github.com/dracory/dataobject"
	"github.com/dromara/carbon/v2"
)

// == TYPE ===================================================================

// entityImplementation represents a schemaless entity backed by a map[string]string
type entityImplementation struct {
	dataobject.DataObject
}

// == INTERFACE COMPLIANCE ===================================================

var _ EntityInterface = (*entityImplementation)(nil)

// == CONSTRUCTORS ===========================================================

// NewEntity creates a new entity with default values
func NewEntity() EntityInterface {
	o := &entityImplementation{}
	o.SetType("")
	o.SetHandle("")
	o.SetID(GenerateShortID())
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	return o
}

// NewEntityFromExistingData creates an entity from a raw data map (e.g. from DB rows)
func NewEntityFromExistingData(data map[string]string) EntityInterface {
	o := &entityImplementation{}
	o.Hydrate(data)
	return o
}

// == GETTERS & SETTERS ======================================================

func (o *entityImplementation) GetType() string {
	return o.Get(COLUMN_ENTITY_TYPE)
}

func (o *entityImplementation) SetType(entityType string) EntityInterface {
	o.Set(COLUMN_ENTITY_TYPE, entityType)
	return o
}

func (o *entityImplementation) GetHandle() string {
	return o.Get(COLUMN_ENTITY_HANDLE)
}

func (o *entityImplementation) SetHandle(handle string) EntityInterface {
	o.Set(COLUMN_ENTITY_HANDLE, handle)
	return o
}

func (o *entityImplementation) GetCreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *entityImplementation) SetCreatedAt(createdAt string) EntityInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *entityImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetCreatedAt(), carbon.UTC)
}

func (o *entityImplementation) GetUpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *entityImplementation) SetUpdatedAt(updatedAt string) EntityInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *entityImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetUpdatedAt(), carbon.UTC)
}

// == DYNAMIC ATTRIBUTES =====================================================

// GetTempKey retrieves an in-memory attribute by key
func (o *entityImplementation) GetTempKey(key string) string {
	return o.Get(key)
}

// SetTempKey sets an in-memory attribute value
func (o *entityImplementation) SetTempKey(key string, value string) EntityInterface {
	o.Set(key, value)
	return o
}

// GetTempKeys returns all dynamic attributes (excludes system columns)
func (o *entityImplementation) GetTempKeys() map[string]string {
	systemColumns := map[string]bool{
		COLUMN_ID:            true,
		COLUMN_ENTITY_TYPE:   true,
		COLUMN_ENTITY_HANDLE: true,
		COLUMN_CREATED_AT:    true,
		COLUMN_UPDATED_AT:    true,
	}

	attrs := make(map[string]string)
	for k, v := range o.Data() {
		if !systemColumns[k] {
			attrs[k] = v
		}
	}
	return attrs
}
