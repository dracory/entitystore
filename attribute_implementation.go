package entitystore

import (
	"strconv"

	"github.com/dracory/dataobject"
	"github.com/dromara/carbon/v2"
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
	o.SetAttributeKey("")
	o.SetAttributeValue("")
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

func (o *attributeImplementation) EntityID() string {
	return o.Get(COLUMN_ENTITY_ID)
}

func (o *attributeImplementation) SetEntityID(entityID string) AttributeInterface {
	o.Set(COLUMN_ENTITY_ID, entityID)
	return o
}

func (o *attributeImplementation) AttributeKey() string {
	return o.Get(COLUMN_ATTRIBUTE_KEY)
}

func (o *attributeImplementation) SetAttributeKey(key string) AttributeInterface {
	o.Set(COLUMN_ATTRIBUTE_KEY, key)
	return o
}

func (o *attributeImplementation) AttributeValue() string {
	return o.Get(COLUMN_ATTRIBUTE_VALUE)
}

func (o *attributeImplementation) SetAttributeValue(value string) AttributeInterface {
	o.Set(COLUMN_ATTRIBUTE_VALUE, value)
	return o
}

func (o *attributeImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *attributeImplementation) SetCreatedAt(createdAt string) AttributeInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *attributeImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt(), carbon.UTC)
}

func (o *attributeImplementation) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *attributeImplementation) SetUpdatedAt(updatedAt string) AttributeInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *attributeImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt(), carbon.UTC)
}

// == TYPE CONVERSIONS =======================================================

// GetInt returns the attribute value parsed as int64
func (o *attributeImplementation) GetInt() (int64, error) {
	return strconv.ParseInt(o.AttributeValue(), 10, 64)
}

// GetFloat returns the attribute value parsed as float64
func (o *attributeImplementation) GetFloat() (float64, error) {
	return strconv.ParseFloat(o.AttributeValue(), 64)
}

// SetInt sets the attribute value from an int64
func (o *attributeImplementation) SetInt(value int64) AttributeInterface {
	o.SetAttributeValue(strconv.FormatInt(value, 10))
	return o
}

// SetFloat sets the attribute value from a float64
func (o *attributeImplementation) SetFloat(value float64) AttributeInterface {
	o.SetAttributeValue(strconv.FormatFloat(value, 'f', 30, 64))
	return o
}
