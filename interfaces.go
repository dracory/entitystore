package entitystore

import (
	"context"
	"database/sql"

	"github.com/dracory/dataobject"
	"github.com/dromara/carbon/v2"
)

// == ENTITY INTERFACE =======================================================

// EntityInterface defines the contract for schemaless entities
type EntityInterface interface {
	dataobject.DataObjectInterface

	// Core getters
	GetType() string
	GetHandle() string
	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon
	GetUpdatedAt() string
	GetUpdatedAtCarbon() *carbon.Carbon

	// Core setters (fluent) — ID() / SetID() come from DataObjectInterface
	SetType(entityType string) EntityInterface
	SetHandle(handle string) EntityInterface
	SetCreatedAt(createdAt string) EntityInterface
	SetUpdatedAt(updatedAt string) EntityInterface

	// Dynamic / extra attributes (in-memory only, not persisted)
	GetTemp(key string) string
	SetTemp(key string, value string) EntityInterface
	GetAllTemp() map[string]string
}

// == ATTRIBUTE INTERFACE ====================================================

// AttributeInterface defines the contract for persisted entity key-value attributes
type AttributeInterface interface {
	dataobject.DataObjectInterface

	// Core getters
	GetEntityID() string
	GetKey() string
	GetValue() string
	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon
	GetUpdatedAt() string
	GetUpdatedAtCarbon() *carbon.Carbon

	// Core setters (fluent) — ID() / SetID() come from DataObjectInterface
	SetEntityID(entityID string) AttributeInterface
	SetKey(key string) AttributeInterface
	SetValue(value string) AttributeInterface
	SetCreatedAt(createdAt string) AttributeInterface
	SetUpdatedAt(updatedAt string) AttributeInterface

	// Type-conversion helpers
	GetInt() (int64, error)
	GetFloat() (float64, error)
	SetInt(value int64) AttributeInterface
	SetFloat(value float64) AttributeInterface
}

// == RELATIONSHIP INTERFACE ==================================================

// RelationshipInterface defines the contract for entity relationships
type RelationshipInterface interface {
	dataobject.DataObjectInterface

	// Core getters
	GetEntityID() string
	GetRelatedEntityID() string
	GetRelationshipType() string
	GetParentID() string
	GetSequence() int
	GetMetadata() string
	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon

	// Core setters (fluent) — ID() / SetID() come from DataObjectInterface
	SetEntityID(entityID string) RelationshipInterface
	SetRelatedEntityID(relatedID string) RelationshipInterface
	SetRelationshipType(relType string) RelationshipInterface
	SetParentID(parentID string) RelationshipInterface
	SetSequence(sequence int) RelationshipInterface
	SetMetadata(metadata string) RelationshipInterface
	SetCreatedAt(createdAt string) RelationshipInterface
}

// == TRASH INTERFACES =======================================================

// EntityTrashInterface defines the contract for trashed entities
type EntityTrashInterface interface {
	dataobject.DataObjectInterface

	// Core getters
	GetType() string
	GetHandle() string
	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon
	GetUpdatedAt() string
	GetUpdatedAtCarbon() *carbon.Carbon
	GetDeletedAt() string
	GetDeletedAtCarbon() *carbon.Carbon
	GetDeletedBy() string

	// Core setters (fluent) — ID() / SetID() come from DataObjectInterface
	SetType(entityType string) EntityTrashInterface
	SetHandle(handle string) EntityTrashInterface
	SetCreatedAt(createdAt string) EntityTrashInterface
	SetUpdatedAt(updatedAt string) EntityTrashInterface
	SetDeletedAt(deletedAt string) EntityTrashInterface
	SetDeletedBy(deletedBy string) EntityTrashInterface
}

// AttributeTrashInterface defines the contract for trashed attributes
type AttributeTrashInterface interface {
	dataobject.DataObjectInterface

	// Core getters
	GetEntityID() string
	GetKey() string
	GetValue() string
	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon
	GetUpdatedAt() string
	GetUpdatedAtCarbon() *carbon.Carbon
	GetDeletedAt() string
	GetDeletedAtCarbon() *carbon.Carbon
	GetDeletedBy() string

	// Core setters (fluent) — ID() / SetID() come from DataObjectInterface
	SetEntityID(entityID string) AttributeTrashInterface
	SetKey(key string) AttributeTrashInterface
	SetValue(value string) AttributeTrashInterface
	SetCreatedAt(createdAt string) AttributeTrashInterface
	SetUpdatedAt(updatedAt string) AttributeTrashInterface
	SetDeletedAt(deletedAt string) AttributeTrashInterface
	SetDeletedBy(deletedBy string) AttributeTrashInterface
}

// RelationshipTrashInterface defines the contract for trashed relationships
type RelationshipTrashInterface interface {
	dataobject.DataObjectInterface

	// Core getters
	GetEntityID() string
	GetRelatedEntityID() string
	GetRelationshipType() string
	GetParentID() string
	GetSequence() int
	GetMetadata() string
	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon
	GetDeletedAt() string
	GetDeletedAtCarbon() *carbon.Carbon
	GetDeletedBy() string

	// Core setters (fluent) — ID() / SetID() come from DataObjectInterface
	SetEntityID(entityID string) RelationshipTrashInterface
	SetRelatedEntityID(relatedID string) RelationshipTrashInterface
	SetRelationshipType(relType string) RelationshipTrashInterface
	SetParentID(parentID string) RelationshipTrashInterface
	SetSequence(sequence int) RelationshipTrashInterface
	SetMetadata(metadata string) RelationshipTrashInterface
	SetCreatedAt(createdAt string) RelationshipTrashInterface
	SetDeletedAt(deletedAt string) RelationshipTrashInterface
	SetDeletedBy(deletedBy string) RelationshipTrashInterface
}

// == TAXONOMY INTERFACES =====================================================

// TaxonomyInterface defines the contract for taxonomies (classification systems)
type TaxonomyInterface interface {
	dataobject.DataObjectInterface

	// Core getters
	GetName() string
	GetSlug() string
	GetDescription() string
	GetParentID() string
	GetEntityTypes() []string
	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon
	GetUpdatedAt() string
	GetUpdatedAtCarbon() *carbon.Carbon

	// Core setters (fluent) — ID() / SetID() come from DataObjectInterface
	SetName(name string) TaxonomyInterface
	SetSlug(slug string) TaxonomyInterface
	SetDescription(desc string) TaxonomyInterface
	SetParentID(parentID string) TaxonomyInterface
	SetEntityTypes(types []string) TaxonomyInterface
	SetCreatedAt(createdAt string) TaxonomyInterface
	SetUpdatedAt(updatedAt string) TaxonomyInterface
}

// TaxonomyTermInterface defines the contract for taxonomy terms
type TaxonomyTermInterface interface {
	dataobject.DataObjectInterface

	// Core getters
	GetTaxonomyID() string
	GetName() string
	GetSlug() string
	GetParentID() string
	GetSortOrder() int
	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon
	GetUpdatedAt() string
	GetUpdatedAtCarbon() *carbon.Carbon

	// Core setters (fluent) — ID() / SetID() come from DataObjectInterface
	SetTaxonomyID(taxonomyID string) TaxonomyTermInterface
	SetName(name string) TaxonomyTermInterface
	SetSlug(slug string) TaxonomyTermInterface
	SetParentID(parentID string) TaxonomyTermInterface
	SetSortOrder(order int) TaxonomyTermInterface
	SetCreatedAt(createdAt string) TaxonomyTermInterface
	SetUpdatedAt(updatedAt string) TaxonomyTermInterface
}

// EntityTaxonomyInterface defines the contract for entity-taxonomy assignments
type EntityTaxonomyInterface interface {
	dataobject.DataObjectInterface

	// Core getters
	GetEntityID() string
	GetTaxonomyID() string
	GetTermID() string
	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon

	// Core setters (fluent) — ID() / SetID() come from DataObjectInterface
	SetEntityID(entityID string) EntityTaxonomyInterface
	SetTaxonomyID(taxonomyID string) EntityTaxonomyInterface
	SetTermID(termID string) EntityTaxonomyInterface
	SetCreatedAt(createdAt string) EntityTaxonomyInterface
}

// TaxonomyTrashInterface defines the contract for trashed taxonomies
type TaxonomyTrashInterface interface {
	dataobject.DataObjectInterface

	// Core getters
	GetName() string
	GetSlug() string
	GetDescription() string
	GetParentID() string
	GetEntityTypes() []string
	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon
	GetUpdatedAt() string
	GetUpdatedAtCarbon() *carbon.Carbon
	GetDeletedAt() string
	GetDeletedAtCarbon() *carbon.Carbon
	GetDeletedBy() string

	// Core setters (fluent) — ID() / SetID() come from DataObjectInterface
	SetName(name string) TaxonomyTrashInterface
	SetSlug(slug string) TaxonomyTrashInterface
	SetDescription(desc string) TaxonomyTrashInterface
	SetParentID(parentID string) TaxonomyTrashInterface
	SetEntityTypes(types []string) TaxonomyTrashInterface
	SetCreatedAt(createdAt string) TaxonomyTrashInterface
	SetUpdatedAt(updatedAt string) TaxonomyTrashInterface
	SetDeletedAt(deletedAt string) TaxonomyTrashInterface
	SetDeletedBy(deletedBy string) TaxonomyTrashInterface
}

// TaxonomyTermTrashInterface defines the contract for trashed taxonomy terms
type TaxonomyTermTrashInterface interface {
	dataobject.DataObjectInterface

	// Core getters
	GetTaxonomyID() string
	GetName() string
	GetSlug() string
	GetParentID() string
	GetSortOrder() int
	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon
	GetUpdatedAt() string
	GetUpdatedAtCarbon() *carbon.Carbon
	GetDeletedAt() string
	GetDeletedAtCarbon() *carbon.Carbon
	GetDeletedBy() string

	// Core setters (fluent) — ID() / SetID() come from DataObjectInterface
	SetTaxonomyID(taxonomyID string) TaxonomyTermTrashInterface
	SetName(name string) TaxonomyTermTrashInterface
	SetSlug(slug string) TaxonomyTermTrashInterface
	SetParentID(parentID string) TaxonomyTermTrashInterface
	SetSortOrder(order int) TaxonomyTermTrashInterface
	SetCreatedAt(createdAt string) TaxonomyTermTrashInterface
	SetUpdatedAt(updatedAt string) TaxonomyTermTrashInterface
	SetDeletedAt(deletedAt string) TaxonomyTermTrashInterface
	SetDeletedBy(deletedBy string) TaxonomyTermTrashInterface
}

// == STORE INTERFACE ========================================================

type StoreInterface interface {
	AutoMigrate(ctx context.Context) error

	GetAttributeTableName() string
	GetAttributeTrashTableName() string
	GetDB() *sql.DB
	GetEntityTableName() string
	GetEntityTrashTableName() string
	GetRelationshipTableName() string
	GetRelationshipTrashTableName() string

	// Attribute CRUD
	AttributeCreate(ctx context.Context, attr AttributeInterface) error
	AttributeCreateWithKeyAndValue(ctx context.Context, entityID string, attributeKey string, attributeValue string) (AttributeInterface, error)
	AttributeFind(ctx context.Context, entityID string, attributeKey string) (AttributeInterface, error)
	AttributeFindByHandle(ctx context.Context, entityType string, entityHandle string, attributeKey string) (AttributeInterface, error)
	AttributeList(ctx context.Context, options AttributeQueryOptions) ([]AttributeInterface, error)
	AttributesSet(ctx context.Context, entityID string, attributes map[string]string) error
	AttributeSetFloat(ctx context.Context, entityID string, attributeKey string, attributeValue float64) error
	AttributeSetInt(ctx context.Context, entityID string, attributeKey string, attributeValue int64) error
	AttributeSetString(ctx context.Context, entityID string, attributeKey string, attributeValue string) error
	AttributeUpdate(ctx context.Context, attr AttributeInterface) error

	// Entity CRUD + helpers
	EntityAttributeList(ctx context.Context, entityID string) ([]AttributeInterface, error)
	EntityCount(ctx context.Context, options EntityQueryOptions) (int64, error)
	EntityCreate(ctx context.Context, entity EntityInterface) error
	EntityCreateWithType(ctx context.Context, entityType string) (EntityInterface, error)
	EntityCreateWithTypeAndAttributes(ctx context.Context, entityType string, attributes map[string]string) (EntityInterface, error)
	EntityDelete(ctx context.Context, entityID string) (bool, error)
	EntityFindByAttribute(ctx context.Context, entityType string, attributeKey string, attributeValue string) (EntityInterface, error)
	EntityFindByHandle(ctx context.Context, entityType string, entityHandle string) (EntityInterface, error)
	EntityFindByID(ctx context.Context, entityID string) (EntityInterface, error)
	EntityList(ctx context.Context, options EntityQueryOptions) ([]EntityInterface, error)
	EntityListByAttribute(ctx context.Context, entityType string, attributeKey string, attributeValue string) ([]EntityInterface, error)
	EntityTrash(ctx context.Context, entityID string) (bool, error)
	EntityUpdate(ctx context.Context, entity EntityInterface) error

	// Relationship CRUD + helpers
	RelationshipCreate(ctx context.Context, relationship RelationshipInterface) error
	RelationshipCreateByOptions(ctx context.Context, options RelationshipOptions) (RelationshipInterface, error)
	RelationshipCount(ctx context.Context, options RelationshipQueryOptions) (int64, error)
	RelationshipDelete(ctx context.Context, relationshipID string) (bool, error)
	RelationshipDeleteAll(ctx context.Context, entityID string) error
	RelationshipFind(ctx context.Context, relationshipID string) (RelationshipInterface, error)
	RelationshipFindByEntities(ctx context.Context, entityID string, relatedEntityID string, relationshipType string) (RelationshipInterface, error)
	RelationshipList(ctx context.Context, options RelationshipQueryOptions) ([]RelationshipInterface, error)
	RelationshipListRelated(ctx context.Context, relatedEntityID string, relationshipType string) ([]RelationshipInterface, error)
	RelationshipRestore(ctx context.Context, relationshipID string) (bool, error)
	RelationshipTrash(ctx context.Context, relationshipID string, deletedBy string) (bool, error)
	RelationshipTrashList(ctx context.Context, options RelationshipQueryOptions) ([]RelationshipTrashInterface, error)

	// Taxonomy CRUD + helpers
	TaxonomyCreate(ctx context.Context, taxonomy TaxonomyInterface) error
	TaxonomyCreateByOptions(ctx context.Context, options TaxonomyOptions) (TaxonomyInterface, error)
	TaxonomyCount(ctx context.Context, options TaxonomyQueryOptions) (int64, error)
	TaxonomyDelete(ctx context.Context, taxonomyID string) (bool, error)
	TaxonomyFind(ctx context.Context, taxonomyID string) (TaxonomyInterface, error)
	TaxonomyFindBySlug(ctx context.Context, slug string) (TaxonomyInterface, error)
	TaxonomyList(ctx context.Context, options TaxonomyQueryOptions) ([]TaxonomyInterface, error)
	TaxonomyRestore(ctx context.Context, taxonomyID string) (bool, error)
	TaxonomyTrash(ctx context.Context, taxonomyID string, deletedBy string) (bool, error)
	TaxonomyTrashList(ctx context.Context, options TaxonomyQueryOptions) ([]TaxonomyTrashInterface, error)
	TaxonomyUpdate(ctx context.Context, taxonomy TaxonomyInterface) error

	// TaxonomyTerm CRUD + helpers
	TaxonomyTermCreate(ctx context.Context, term TaxonomyTermInterface) error
	TaxonomyTermCreateByOptions(ctx context.Context, options TaxonomyTermOptions) (TaxonomyTermInterface, error)
	TaxonomyTermCount(ctx context.Context, options TaxonomyTermQueryOptions) (int64, error)
	TaxonomyTermDelete(ctx context.Context, termID string) (bool, error)
	TaxonomyTermFind(ctx context.Context, termID string) (TaxonomyTermInterface, error)
	TaxonomyTermFindBySlug(ctx context.Context, taxonomyID string, slug string) (TaxonomyTermInterface, error)
	TaxonomyTermList(ctx context.Context, options TaxonomyTermQueryOptions) ([]TaxonomyTermInterface, error)
	TaxonomyTermRestore(ctx context.Context, termID string) (bool, error)
	TaxonomyTermTrash(ctx context.Context, termID string, deletedBy string) (bool, error)
	TaxonomyTermTrashList(ctx context.Context, options TaxonomyTermQueryOptions) ([]TaxonomyTermTrashInterface, error)
	TaxonomyTermUpdate(ctx context.Context, term TaxonomyTermInterface) error

	// EntityTaxonomy CRUD + helpers
	EntityTaxonomyAssign(ctx context.Context, entityID string, taxonomyID string, termID string) error
	EntityTaxonomyCount(ctx context.Context, options EntityTaxonomyQueryOptions) (int64, error)
	EntityTaxonomyList(ctx context.Context, options EntityTaxonomyQueryOptions) ([]EntityTaxonomyInterface, error)
	EntityTaxonomyRemove(ctx context.Context, entityID string, taxonomyID string, termID string) error

	// Getters for table names
	GetTaxonomyTableName() string
	GetTaxonomyTrashTableName() string
	GetTaxonomyTermTableName() string
	GetTaxonomyTermTrashTableName() string
	GetEntityTaxonomyTableName() string
}
