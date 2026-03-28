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
func NewTaxonomyTrash() TaxonomyTrashInterface {
	o := &taxonomyTrashImplementation{}
	o.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetDeletedBy("")
	return o
}

// NewTaxonomyTrashFromExistingData creates a taxonomy trash record from a raw data map
func NewTaxonomyTrashFromExistingData(data map[string]string) TaxonomyTrashInterface {
	o := &taxonomyTrashImplementation{}
	o.Hydrate(data)
	return o
}

// == GETTERS & SETTERS ======================================================

func (o *taxonomyTrashImplementation) Name() string {
	return o.Get(COLUMN_NAME)
}

func (o *taxonomyTrashImplementation) SetName(name string) TaxonomyTrashInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

func (o *taxonomyTrashImplementation) Slug() string {
	return o.Get(COLUMN_SLUG)
}

func (o *taxonomyTrashImplementation) SetSlug(slug string) TaxonomyTrashInterface {
	o.Set(COLUMN_SLUG, slug)
	return o
}

func (o *taxonomyTrashImplementation) Description() string {
	return o.Get(COLUMN_DESCRIPTION)
}

func (o *taxonomyTrashImplementation) SetDescription(desc string) TaxonomyTrashInterface {
	o.Set(COLUMN_DESCRIPTION, desc)
	return o
}

func (o *taxonomyTrashImplementation) ParentID() string {
	return o.Get(COLUMN_PARENT_ID)
}

func (o *taxonomyTrashImplementation) SetParentID(parentID string) TaxonomyTrashInterface {
	o.Set(COLUMN_PARENT_ID, parentID)
	return o
}

func (o *taxonomyTrashImplementation) EntityTypes() []string {
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

func (o *taxonomyTrashImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *taxonomyTrashImplementation) SetCreatedAt(createdAt string) TaxonomyTrashInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *taxonomyTrashImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt(), carbon.UTC)
}

func (o *taxonomyTrashImplementation) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *taxonomyTrashImplementation) SetUpdatedAt(updatedAt string) TaxonomyTrashInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *taxonomyTrashImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt(), carbon.UTC)
}

func (o *taxonomyTrashImplementation) DeletedAt() string {
	return o.Get(COLUMN_DELETED_AT)
}

func (o *taxonomyTrashImplementation) SetDeletedAt(deletedAt string) TaxonomyTrashInterface {
	o.Set(COLUMN_DELETED_AT, deletedAt)
	return o
}

func (o *taxonomyTrashImplementation) DeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.DeletedAt(), carbon.UTC)
}

func (o *taxonomyTrashImplementation) DeletedBy() string {
	return o.Get(COLUMN_DELETED_BY)
}

func (o *taxonomyTrashImplementation) SetDeletedBy(deletedBy string) TaxonomyTrashInterface {
	o.Set(COLUMN_DELETED_BY, deletedBy)
	return o
}
