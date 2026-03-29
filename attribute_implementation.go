package entitystore

import (
	"github.com/dracory/dataobject"
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
)

// == TYPE ===================================================================

// attributeImplementation represents a single persisted key-value attribute of an entity
type attributeImplementation struct {
	dataobject.DataObject
}

// == INTERFACE COMPLIANCE ===================================================

var _ AttributeInterface = (*attributeImplementation)(nil)

// == CONSTRUCTORS ===========================================================

// NewAttribute creates a new attribute with default values
func NewAttribute() AttributeInterface {
	o := &attributeImplementation{}
	o.SetID(GenerateShortID())
	o.SetEntityID("")
	o.SetKey("")
	o.SetValue("")
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	return o
}

// NewAttributeFromExistingData creates an attribute from a raw data map (e.g. from DB rows)
func NewAttributeFromExistingData(data map[string]string) AttributeInterface {
	o := &attributeImplementation{}
	o.Hydrate(data)
	return o
}

// == GETTERS & SETTERS ======================================================

func (o *attributeImplementation) GetEntityID() string {
	return o.Get(COLUMN_ENTITY_ID)
}

func (o *attributeImplementation) SetEntityID(entityID string) AttributeInterface {
	o.Set(COLUMN_ENTITY_ID, entityID)
	return o
}

func (o *attributeImplementation) GetKey() string {
	return o.Get(COLUMN_ATTRIBUTE_KEY)
}

func (o *attributeImplementation) SetKey(key string) AttributeInterface {
	o.Set(COLUMN_ATTRIBUTE_KEY, key)
	return o
}

func (o *attributeImplementation) GetValue() string {
	return o.Get(COLUMN_ATTRIBUTE_VALUE)
}

func (o *attributeImplementation) SetValue(value string) AttributeInterface {
	o.Set(COLUMN_ATTRIBUTE_VALUE, value)
	return o
}

func (o *attributeImplementation) GetCreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *attributeImplementation) SetCreatedAt(createdAt string) AttributeInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *attributeImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetCreatedAt(), carbon.UTC)
}

func (o *attributeImplementation) GetUpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *attributeImplementation) SetUpdatedAt(updatedAt string) AttributeInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *attributeImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetUpdatedAt(), carbon.UTC)
}

// == TYPE CONVERSIONS =======================================================

// GetInt returns the attribute value parsed as int64
func (o *attributeImplementation) GetInt() (int64, error) {
	return cast.ToInt64E(o.GetValue())
}

// GetFloat returns the attribute value parsed as float64
func (o *attributeImplementation) GetFloat() (float64, error) {
	return cast.ToFloat64E(o.GetValue())
}

// SetInt sets the attribute value from an int64
func (o *attributeImplementation) SetInt(value int64) AttributeInterface {
	o.SetValue(cast.ToString(value))
	return o
}

// SetFloat sets the attribute value from a float64
func (o *attributeImplementation) SetFloat(value float64) AttributeInterface {
	o.SetValue(cast.ToString(value))
	return o
}
