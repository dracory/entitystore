package entitystore

import (
	"strconv"

	"github.com/dracory/dataobject"
	"github.com/dromara/carbon/v2"
)

// == TYPE ===================================================================

// taxonomyTermImplementation represents a taxonomy term backed by a map[string]string
type taxonomyTermImplementation struct {
	dataobject.DataObject
}

// == INTERFACE COMPLIANCE ===================================================

var _ TaxonomyTermInterface = (*taxonomyTermImplementation)(nil)

// == CONSTRUCTORS ===========================================================

// NewTaxonomyTerm creates a new taxonomy term with default values
func NewTaxonomyTerm() TaxonomyTermInterface {
	o := &taxonomyTermImplementation{}
	o.SetID(GenerateShortID())
	o.SetTaxonomyID("")
	o.SetName("")
	o.SetSlug("")
	o.SetParentID("")
	o.SetSortOrder(0)
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	return o
}

// NewTaxonomyTermFromExistingData creates a taxonomy term from a raw data map (e.g. from DB rows)
func NewTaxonomyTermFromExistingData(data map[string]string) TaxonomyTermInterface {
	o := &taxonomyTermImplementation{}
	o.Hydrate(data)
	return o
}

// == GETTERS & SETTERS ======================================================

func (o *taxonomyTermImplementation) TaxonomyID() string {
	return o.Get(COLUMN_TAXONOMY_ID)
}

func (o *taxonomyTermImplementation) SetTaxonomyID(taxonomyID string) TaxonomyTermInterface {
	o.Set(COLUMN_TAXONOMY_ID, taxonomyID)
	return o
}

func (o *taxonomyTermImplementation) Name() string {
	return o.Get(COLUMN_NAME)
}

func (o *taxonomyTermImplementation) SetName(name string) TaxonomyTermInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

func (o *taxonomyTermImplementation) Slug() string {
	return o.Get(COLUMN_SLUG)
}

func (o *taxonomyTermImplementation) SetSlug(slug string) TaxonomyTermInterface {
	o.Set(COLUMN_SLUG, slug)
	return o
}

func (o *taxonomyTermImplementation) ParentID() string {
	return o.Get(COLUMN_PARENT_ID)
}

func (o *taxonomyTermImplementation) SetParentID(parentID string) TaxonomyTermInterface {
	o.Set(COLUMN_PARENT_ID, parentID)
	return o
}

func (o *taxonomyTermImplementation) SortOrder() int {
	val, _ := strconv.Atoi(o.Get(COLUMN_SORT_ORDER))
	return val
}

func (o *taxonomyTermImplementation) SetSortOrder(order int) TaxonomyTermInterface {
	o.Set(COLUMN_SORT_ORDER, strconv.Itoa(order))
	return o
}

func (o *taxonomyTermImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *taxonomyTermImplementation) SetCreatedAt(createdAt string) TaxonomyTermInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *taxonomyTermImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt(), carbon.UTC)
}

func (o *taxonomyTermImplementation) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *taxonomyTermImplementation) SetUpdatedAt(updatedAt string) TaxonomyTermInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *taxonomyTermImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt(), carbon.UTC)
}
