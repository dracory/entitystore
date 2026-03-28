package entitystore

import (
	"strconv"

	"github.com/dracory/dataobject"
	"github.com/dromara/carbon/v2"
)

// == TYPE ===================================================================

// taxonomyTermTrashImplementation represents a trashed taxonomy term backed by a map[string]string
type taxonomyTermTrashImplementation struct {
	dataobject.DataObject
}

// == INTERFACE COMPLIANCE ===================================================

var _ TaxonomyTermTrashInterface = (*taxonomyTermTrashImplementation)(nil)

// == CONSTRUCTORS ===========================================================

// NewTaxonomyTermTrash creates a new taxonomy term trash record with default values
func NewTaxonomyTermTrash() TaxonomyTermTrashInterface {
	o := &taxonomyTermTrashImplementation{}
	o.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetDeletedBy("")
	return o
}

// NewTaxonomyTermTrashFromExistingData creates a taxonomy term trash record from a raw data map
func NewTaxonomyTermTrashFromExistingData(data map[string]string) TaxonomyTermTrashInterface {
	o := &taxonomyTermTrashImplementation{}
	o.Hydrate(data)
	return o
}

// == GETTERS & SETTERS ======================================================

func (o *taxonomyTermTrashImplementation) TaxonomyID() string {
	return o.Get(COLUMN_TAXONOMY_ID)
}

func (o *taxonomyTermTrashImplementation) SetTaxonomyID(taxonomyID string) TaxonomyTermTrashInterface {
	o.Set(COLUMN_TAXONOMY_ID, taxonomyID)
	return o
}

func (o *taxonomyTermTrashImplementation) Name() string {
	return o.Get(COLUMN_NAME)
}

func (o *taxonomyTermTrashImplementation) SetName(name string) TaxonomyTermTrashInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

func (o *taxonomyTermTrashImplementation) Slug() string {
	return o.Get(COLUMN_SLUG)
}

func (o *taxonomyTermTrashImplementation) SetSlug(slug string) TaxonomyTermTrashInterface {
	o.Set(COLUMN_SLUG, slug)
	return o
}

func (o *taxonomyTermTrashImplementation) ParentID() string {
	return o.Get(COLUMN_PARENT_ID)
}

func (o *taxonomyTermTrashImplementation) SetParentID(parentID string) TaxonomyTermTrashInterface {
	o.Set(COLUMN_PARENT_ID, parentID)
	return o
}

func (o *taxonomyTermTrashImplementation) SortOrder() int {
	val, _ := strconv.Atoi(o.Get(COLUMN_SORT_ORDER))
	return val
}

func (o *taxonomyTermTrashImplementation) SetSortOrder(order int) TaxonomyTermTrashInterface {
	o.Set(COLUMN_SORT_ORDER, strconv.Itoa(order))
	return o
}

func (o *taxonomyTermTrashImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *taxonomyTermTrashImplementation) SetCreatedAt(createdAt string) TaxonomyTermTrashInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *taxonomyTermTrashImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt(), carbon.UTC)
}

func (o *taxonomyTermTrashImplementation) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *taxonomyTermTrashImplementation) SetUpdatedAt(updatedAt string) TaxonomyTermTrashInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *taxonomyTermTrashImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt(), carbon.UTC)
}

func (o *taxonomyTermTrashImplementation) DeletedAt() string {
	return o.Get(COLUMN_DELETED_AT)
}

func (o *taxonomyTermTrashImplementation) SetDeletedAt(deletedAt string) TaxonomyTermTrashInterface {
	o.Set(COLUMN_DELETED_AT, deletedAt)
	return o
}

func (o *taxonomyTermTrashImplementation) DeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.DeletedAt(), carbon.UTC)
}

func (o *taxonomyTermTrashImplementation) DeletedBy() string {
	return o.Get(COLUMN_DELETED_BY)
}

func (o *taxonomyTermTrashImplementation) SetDeletedBy(deletedBy string) TaxonomyTermTrashInterface {
	o.Set(COLUMN_DELETED_BY, deletedBy)
	return o
}
