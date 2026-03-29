package entitystore

import (
	"encoding/json"

	"github.com/dracory/dataobject"
	"github.com/dromara/carbon/v2"
)

// == TYPE ===================================================================

// taxonomyTrashImplementation represents a trashed taxonomy backed by a map[string]string
type taxonomyTrashImplementation struct {
	dataobject.DataObject
}

// == INTERFACE COMPLIANCE ===================================================

var _ TaxonomyTrashInterface = (*taxonomyTrashImplementation)(nil)

// == CONSTRUCTORS ===========================================================

// NewTaxonomyTrash creates a new taxonomy trash record with default values
// Sets default deleted_at timestamp and empty deleted_by
func NewTaxonomyTrash() TaxonomyTrashInterface {
	o := &taxonomyTrashImplementation{}
	o.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetDeletedBy("")
	return o
}

// NewTaxonomyTrashFromExistingData creates a taxonomy trash record from a raw data map (e.g., from DB rows)
// Used internally when hydrating trashed taxonomies from database results
func NewTaxonomyTrashFromExistingData(data map[string]string) TaxonomyTrashInterface {
	o := &taxonomyTrashImplementation{}
	o.Hydrate(data)
	return o
}

// == GETTERS & SETTERS ======================================================

func (o *taxonomyTrashImplementation) GetID() string {
	return o.ID()
}

func (o *taxonomyTrashImplementation) GetName() string {
	return o.Get(COLUMN_NAME)
}

func (o *taxonomyTrashImplementation) SetName(name string) TaxonomyTrashInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

func (o *taxonomyTrashImplementation) GetSlug() string {
	return o.Get(COLUMN_SLUG)
}

func (o *taxonomyTrashImplementation) SetSlug(slug string) TaxonomyTrashInterface {
	o.Set(COLUMN_SLUG, slug)
	return o
}

func (o *taxonomyTrashImplementation) GetDescription() string {
	return o.Get(COLUMN_DESCRIPTION)
}

func (o *taxonomyTrashImplementation) SetDescription(desc string) TaxonomyTrashInterface {
	o.Set(COLUMN_DESCRIPTION, desc)
	return o
}

func (o *taxonomyTrashImplementation) GetParentID() string {
	return o.Get(COLUMN_PARENT_ID)
}

func (o *taxonomyTrashImplementation) SetParentID(parentID string) TaxonomyTrashInterface {
	o.Set(COLUMN_PARENT_ID, parentID)
	return o
}

func (o *taxonomyTrashImplementation) GetEntityTypes() []string {
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

func (o *taxonomyTrashImplementation) SetEntityTypes(types []string) TaxonomyTrashInterface {
	if types == nil {
		types = []string{}
	}
	data, _ := json.Marshal(types)
	o.Set(COLUMN_ENTITY_TYPES, string(data))
	return o
}

func (o *taxonomyTrashImplementation) GetCreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *taxonomyTrashImplementation) SetCreatedAt(createdAt string) TaxonomyTrashInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *taxonomyTrashImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetCreatedAt(), carbon.UTC)
}

func (o *taxonomyTrashImplementation) GetUpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *taxonomyTrashImplementation) SetUpdatedAt(updatedAt string) TaxonomyTrashInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *taxonomyTrashImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetUpdatedAt(), carbon.UTC)
}

func (o *taxonomyTrashImplementation) GetDeletedAt() string {
	return o.Get(COLUMN_DELETED_AT)
}

func (o *taxonomyTrashImplementation) SetDeletedAt(deletedAt string) TaxonomyTrashInterface {
	o.Set(COLUMN_DELETED_AT, deletedAt)
	return o
}

func (o *taxonomyTrashImplementation) GetDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetDeletedAt(), carbon.UTC)
}

func (o *taxonomyTrashImplementation) GetDeletedBy() string {
	return o.Get(COLUMN_DELETED_BY)
}

func (o *taxonomyTrashImplementation) SetDeletedBy(deletedBy string) TaxonomyTrashInterface {
	o.Set(COLUMN_DELETED_BY, deletedBy)
	return o
}
