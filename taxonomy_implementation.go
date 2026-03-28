package entitystore

import (
	"encoding/json"

	"github.com/dracory/dataobject"
	"github.com/dromara/carbon/v2"
)

// == TYPE ===================================================================

// taxonomyImplementation represents a taxonomy backed by a map[string]string
type taxonomyImplementation struct {
	dataobject.DataObject
}

// == INTERFACE COMPLIANCE ===================================================

var _ TaxonomyInterface = (*taxonomyImplementation)(nil)

// == CONSTRUCTORS ===========================================================

// NewTaxonomy creates a new taxonomy with default values
func NewTaxonomy() TaxonomyInterface {
	o := &taxonomyImplementation{}
	o.SetID(GenerateShortID())
	o.SetName("")
	o.SetSlug("")
	o.SetDescription("")
	o.SetParentID("")
	o.SetEntityTypes([]string{})
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	return o
}

// NewTaxonomyFromExistingData creates a taxonomy from a raw data map (e.g. from DB rows)
func NewTaxonomyFromExistingData(data map[string]string) TaxonomyInterface {
	o := &taxonomyImplementation{}
	o.Hydrate(data)
	return o
}

// == GETTERS & SETTERS ======================================================

func (o *taxonomyImplementation) Name() string {
	return o.Get(COLUMN_NAME)
}

func (o *taxonomyImplementation) SetName(name string) TaxonomyInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

func (o *taxonomyImplementation) Slug() string {
	return o.Get(COLUMN_SLUG)
}

func (o *taxonomyImplementation) SetSlug(slug string) TaxonomyInterface {
	o.Set(COLUMN_SLUG, slug)
	return o
}

func (o *taxonomyImplementation) Description() string {
	return o.Get(COLUMN_DESCRIPTION)
}

func (o *taxonomyImplementation) SetDescription(desc string) TaxonomyInterface {
	o.Set(COLUMN_DESCRIPTION, desc)
	return o
}

func (o *taxonomyImplementation) ParentID() string {
	return o.Get(COLUMN_PARENT_ID)
}

func (o *taxonomyImplementation) SetParentID(parentID string) TaxonomyInterface {
	o.Set(COLUMN_PARENT_ID, parentID)
	return o
}

func (o *taxonomyImplementation) EntityTypes() []string {
	typesStr := o.Get(COLUMN_ENTITY_TYPES)
	if typesStr == "" {
		return []string{}
	}
	var types []string
	if err := json.Unmarshal([]byte(typesStr), &types); err != nil {
		return []string{}
	}
	return types
}

func (o *taxonomyImplementation) SetEntityTypes(types []string) TaxonomyInterface {
	if types == nil {
		types = []string{}
	}
	data, _ := json.Marshal(types)
	o.Set(COLUMN_ENTITY_TYPES, string(data))
	return o
}

func (o *taxonomyImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *taxonomyImplementation) SetCreatedAt(createdAt string) TaxonomyInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *taxonomyImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt(), carbon.UTC)
}

func (o *taxonomyImplementation) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *taxonomyImplementation) SetUpdatedAt(updatedAt string) TaxonomyInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *taxonomyImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt(), carbon.UTC)
}
