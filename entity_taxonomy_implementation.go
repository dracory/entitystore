package entitystore

import (
	"github.com/dracory/dataobject"
	"github.com/dromara/carbon/v2"
)

// == TYPE ===================================================================

// entityTaxonomyImplementation represents an entity-taxonomy assignment backed by a map[string]string
type entityTaxonomyImplementation struct {
	dataobject.DataObject
}

// == INTERFACE COMPLIANCE ===================================================

var _ EntityTaxonomyInterface = (*entityTaxonomyImplementation)(nil)

// == CONSTRUCTORS ===========================================================

// NewEntityTaxonomy creates a new entity-taxonomy assignment with default values
// Generates a new short ID and sets timestamps to the current time
func NewEntityTaxonomy() EntityTaxonomyInterface {
	o := &entityTaxonomyImplementation{}
	o.SetID(GenerateShortID())
	o.SetEntityID("")
	o.SetTaxonomyID("")
	o.SetTermID("")
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	return o
}

// NewEntityTaxonomyFromExistingData creates an entity-taxonomy assignment from a raw data map (e.g., from DB rows)
// Used internally when hydrating entity-taxonomy assignments from database results
func NewEntityTaxonomyFromExistingData(data map[string]string) EntityTaxonomyInterface {
	o := &entityTaxonomyImplementation{}
	o.Hydrate(data)
	return o
}

// == GETTERS & SETTERS ======================================================

func (o *entityTaxonomyImplementation) GetID() string {
	return o.ID()
}

func (o *entityTaxonomyImplementation) GetEntityID() string {
	return o.Get(COLUMN_ENTITY_ID)
}

func (o *entityTaxonomyImplementation) SetEntityID(entityID string) EntityTaxonomyInterface {
	o.Set(COLUMN_ENTITY_ID, entityID)
	return o
}

func (o *entityTaxonomyImplementation) GetTaxonomyID() string {
	return o.Get(COLUMN_TAXONOMY_ID)
}

func (o *entityTaxonomyImplementation) SetTaxonomyID(taxonomyID string) EntityTaxonomyInterface {
	o.Set(COLUMN_TAXONOMY_ID, taxonomyID)
	return o
}

func (o *entityTaxonomyImplementation) GetTermID() string {
	return o.Get(COLUMN_TERM_ID)
}

func (o *entityTaxonomyImplementation) SetTermID(termID string) EntityTaxonomyInterface {
	o.Set(COLUMN_TERM_ID, termID)
	return o
}

func (o *entityTaxonomyImplementation) GetCreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *entityTaxonomyImplementation) SetCreatedAt(createdAt string) EntityTaxonomyInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *entityTaxonomyImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetCreatedAt(), carbon.UTC)
}
