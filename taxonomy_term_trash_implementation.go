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
// Sets default deleted_at timestamp and empty deleted_by
func NewTaxonomyTermTrash() TaxonomyTermTrashInterface {
	o := &taxonomyTermTrashImplementation{}
	o.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetDeletedBy("")
	return o
}

// NewTaxonomyTermTrashFromExistingData creates a taxonomy term trash record from a raw data map (e.g., from DB rows)
// Used internally when hydrating trashed taxonomy terms from database results
func NewTaxonomyTermTrashFromExistingData(data map[string]string) TaxonomyTermTrashInterface {
	o := &taxonomyTermTrashImplementation{}
	o.Hydrate(data)
	return o
}

// == GETTERS & SETTERS ======================================================

func (o *taxonomyTermTrashImplementation) GetID() string {
	return o.ID()
}

func (o *taxonomyTermTrashImplementation) GetTaxonomyID() string {
	return o.Get(COLUMN_TAXONOMY_ID)
}

func (o *taxonomyTermTrashImplementation) SetTaxonomyID(taxonomyID string) TaxonomyTermTrashInterface {
	o.Set(COLUMN_TAXONOMY_ID, taxonomyID)
	return o
}

func (o *taxonomyTermTrashImplementation) GetName() string {
	return o.Get(COLUMN_NAME)
}

func (o *taxonomyTermTrashImplementation) SetName(name string) TaxonomyTermTrashInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

func (o *taxonomyTermTrashImplementation) GetSlug() string {
	return o.Get(COLUMN_SLUG)
}

func (o *taxonomyTermTrashImplementation) SetSlug(slug string) TaxonomyTermTrashInterface {
	o.Set(COLUMN_SLUG, slug)
	return o
}

func (o *taxonomyTermTrashImplementation) GetParentID() string {
	return o.Get(COLUMN_PARENT_ID)
}

func (o *taxonomyTermTrashImplementation) SetParentID(parentID string) TaxonomyTermTrashInterface {
	o.Set(COLUMN_PARENT_ID, parentID)
	return o
}

func (o *taxonomyTermTrashImplementation) GetSortOrder() int {
	val, _ := strconv.Atoi(o.Get(COLUMN_SORT_ORDER))
	return val
}

func (o *taxonomyTermTrashImplementation) SetSortOrder(order int) TaxonomyTermTrashInterface {
	o.Set(COLUMN_SORT_ORDER, strconv.Itoa(order))
	return o
}

func (o *taxonomyTermTrashImplementation) GetCreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *taxonomyTermTrashImplementation) SetCreatedAt(createdAt string) TaxonomyTermTrashInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *taxonomyTermTrashImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetCreatedAt(), carbon.UTC)
}

func (o *taxonomyTermTrashImplementation) GetUpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *taxonomyTermTrashImplementation) SetUpdatedAt(updatedAt string) TaxonomyTermTrashInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *taxonomyTermTrashImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetUpdatedAt(), carbon.UTC)
}

func (o *taxonomyTermTrashImplementation) GetDeletedAt() string {
	return o.Get(COLUMN_DELETED_AT)
}

func (o *taxonomyTermTrashImplementation) SetDeletedAt(deletedAt string) TaxonomyTermTrashInterface {
	o.Set(COLUMN_DELETED_AT, deletedAt)
	return o
}

func (o *taxonomyTermTrashImplementation) GetDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetDeletedAt(), carbon.UTC)
}

func (o *taxonomyTermTrashImplementation) GetDeletedBy() string {
	return o.Get(COLUMN_DELETED_BY)
}

func (o *taxonomyTermTrashImplementation) SetDeletedBy(deletedBy string) TaxonomyTermTrashInterface {
	o.Set(COLUMN_DELETED_BY, deletedBy)
	return o
}
