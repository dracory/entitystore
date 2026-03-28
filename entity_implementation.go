package entitystore

import (
	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
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
	o.SetEntityType("")
	o.SetEntityHandle("")
	o.SetStatus(ENTITY_STATUS_ACTIVE)
	o.SetID(GenerateShortID())
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetSoftDeletedAt(sb.MAX_DATETIME)
	return o
}

// NewEntityFromExistingData creates an entity from a raw data map (e.g. from DB rows)
func NewEntityFromExistingData(data map[string]string) EntityInterface {
	o := &entityImplementation{}
	o.Hydrate(data)
	return o
}

// == STATUS HELPERS =========================================================

func (o *entityImplementation) IsActive() bool {
	return o.Status() == ENTITY_STATUS_ACTIVE
}

func (o *entityImplementation) IsInactive() bool {
	return o.Status() == ENTITY_STATUS_INACTIVE
}

func (o *entityImplementation) IsSoftDeleted() bool {
	return o.SoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// == GETTERS & SETTERS ======================================================

func (o *entityImplementation) EntityType() string {
	return o.Get(COLUMN_ENTITY_TYPE)
}

func (o *entityImplementation) SetEntityType(entityType string) EntityInterface {
	o.Set(COLUMN_ENTITY_TYPE, entityType)
	return o
}

func (o *entityImplementation) EntityHandle() string {
	return o.Get(COLUMN_ENTITY_HANDLE)
}

func (o *entityImplementation) SetEntityHandle(handle string) EntityInterface {
	o.Set(COLUMN_ENTITY_HANDLE, handle)
	return o
}

func (o *entityImplementation) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *entityImplementation) SetStatus(status string) EntityInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *entityImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *entityImplementation) SetCreatedAt(createdAt string) EntityInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *entityImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt(), carbon.UTC)
}

func (o *entityImplementation) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *entityImplementation) SetUpdatedAt(updatedAt string) EntityInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *entityImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt(), carbon.UTC)
}

func (o *entityImplementation) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

func (o *entityImplementation) SetSoftDeletedAt(softDeletedAt string) EntityInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

func (o *entityImplementation) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt(), carbon.UTC)
}

// == DYNAMIC ATTRIBUTES =====================================================

// GetAttribute retrieves an in-memory attribute by key
func (o *entityImplementation) GetAttribute(key string) string {
	return o.Get(key)
}

// SetAttribute sets an in-memory attribute value
func (o *entityImplementation) SetAttribute(key string, value string) EntityInterface {
	o.Set(key, value)
	return o
}

// GetAllAttributes returns all dynamic attributes (excludes system columns)
func (o *entityImplementation) GetAllAttributes() map[string]string {
	systemColumns := map[string]bool{
		COLUMN_ID:              true,
		COLUMN_ENTITY_TYPE:     true,
		COLUMN_ENTITY_HANDLE:   true,
		COLUMN_STATUS:          true,
		COLUMN_CREATED_AT:      true,
		COLUMN_UPDATED_AT:      true,
		COLUMN_SOFT_DELETED_AT: true,
	}

	attrs := make(map[string]string)
	for k, v := range o.Data() {
		if !systemColumns[k] {
			attrs[k] = v
		}
	}
	return attrs
}
