package entitystore

import (
	"context"
	"database/sql"

	"github.com/dracory/dataobject"
	"github.com/dromara/carbon/v2"
)

// == ENTITY INTERFACE =======================================================

// EntityInterface defines the contract for schemaless entities
// Entities are the core data objects in the store, identified by a unique ID
// and associated with a type and optional handle.
type EntityInterface interface {
	dataobject.DataObjectInterface

	// GetID returns the unique identifier for this entity
	GetID() string
	// GetType returns the entity type (e.g., "user", "product")
	GetType() string
	// GetHandle returns the optional human-readable handle/slug
	GetHandle() string
	// GetCreatedAt returns the creation timestamp in UTC
	GetCreatedAt() string
	// GetCreatedAtCarbon returns the creation timestamp as a Carbon instance
	GetCreatedAtCarbon() *carbon.Carbon
	// GetUpdatedAt returns the last update timestamp in UTC
	GetUpdatedAt() string
	// GetUpdatedAtCarbon returns the last update timestamp as a Carbon instance
	GetUpdatedAtCarbon() *carbon.Carbon

	// SetType sets the entity type and returns the entity for method chaining
	SetType(entityType string) EntityInterface
	// SetHandle sets the entity handle and returns the entity for method chaining
	SetHandle(handle string) EntityInterface
	// SetCreatedAt sets the creation timestamp and returns the entity for method chaining
	SetCreatedAt(createdAt string) EntityInterface
	// SetUpdatedAt sets the update timestamp and returns the entity for method chaining
	SetUpdatedAt(updatedAt string) EntityInterface

	// GetTempKey retrieves a temporary in-memory attribute by key
	// These attributes are not persisted to the database
	GetTempKey(key string) string
	// SetTempKey sets a temporary in-memory attribute value
	// These attributes are not persisted to the database
	SetTempKey(key string, value string) EntityInterface
	// GetTempKeys returns all temporary in-memory attributes (excludes system columns)
	GetTempKeys() map[string]string
}

// == ATTRIBUTE INTERFACE ====================================================

// AttributeInterface defines the contract for persisted entity key-value attributes
// Attributes store typed values associated with entities in a separate table,
// allowing schemaless data storage while maintaining queryability.
type AttributeInterface interface {
	dataobject.DataObjectInterface

	// GetID returns the unique identifier for this attribute record
	GetID() string
	// GetEntityID returns the ID of the entity this attribute belongs to
	GetEntityID() string
	// GetKey returns the attribute name/key
	GetKey() string
	// GetValue returns the attribute value as a string
	GetValue() string
	// GetCreatedAt returns the creation timestamp in UTC
	GetCreatedAt() string
	// GetCreatedAtCarbon returns the creation timestamp as a Carbon instance
	GetCreatedAtCarbon() *carbon.Carbon
	// GetUpdatedAt returns the last update timestamp in UTC
	GetUpdatedAt() string
	// GetUpdatedAtCarbon returns the last update timestamp as a Carbon instance
	GetUpdatedAtCarbon() *carbon.Carbon

	// SetEntityID sets the associated entity ID and returns the attribute for chaining
	SetEntityID(entityID string) AttributeInterface
	// SetKey sets the attribute key/name and returns the attribute for chaining
	SetKey(key string) AttributeInterface
	// SetValue sets the attribute value as a string and returns the attribute for chaining
	SetValue(value string) AttributeInterface
	// SetCreatedAt sets the creation timestamp and returns the attribute for chaining
	SetCreatedAt(createdAt string) AttributeInterface
	// SetUpdatedAt sets the update timestamp and returns the attribute for chaining
	SetUpdatedAt(updatedAt string) AttributeInterface

	// GetInt parses and returns the value as an int64
	GetInt() (int64, error)
	// GetFloat parses and returns the value as a float64
	GetFloat() (float64, error)
	// SetInt stores an int64 value, converting it to a string
	SetInt(value int64) AttributeInterface
	// SetFloat stores a float64 value, converting it to a string
	SetFloat(value float64) AttributeInterface
}

// == RELATIONSHIP INTERFACE ==================================================

// RelationshipInterface defines the contract for entity relationships
// Relationships link entities together with types like "belongs_to", "has_many", "many_to_many".
type RelationshipInterface interface {
	dataobject.DataObjectInterface

	// GetID returns the unique identifier for this relationship
	GetID() string
	// GetEntityID returns the ID of the source entity
	GetEntityID() string
	// GetRelatedEntityID returns the ID of the target/related entity
	GetRelatedEntityID() string
	// GetRelationshipType returns the relationship type (e.g., "belongs_to", "has_many", "many_to_many")
	GetRelationshipType() string
	// GetParentID returns the parent relationship ID for hierarchical relationships
	GetParentID() string
	// GetSequence returns the sort order for this relationship
	GetSequence() int
	// GetMetadata returns JSON metadata associated with this relationship
	GetMetadata() string
	// GetCreatedAt returns the creation timestamp in UTC
	GetCreatedAt() string
	// GetCreatedAtCarbon returns the creation timestamp as a Carbon instance
	GetCreatedAtCarbon() *carbon.Carbon

	// SetEntityID sets the source entity ID and returns the relationship for chaining
	SetEntityID(entityID string) RelationshipInterface
	// SetRelatedEntityID sets the target entity ID and returns the relationship for chaining
	SetRelatedEntityID(relatedID string) RelationshipInterface
	// SetRelationshipType sets the relationship type and returns the relationship for chaining
	SetRelationshipType(relType string) RelationshipInterface
	// SetParentID sets the parent relationship ID and returns the relationship for chaining
	SetParentID(parentID string) RelationshipInterface
	// SetSequence sets the sort order and returns the relationship for chaining
	SetSequence(sequence int) RelationshipInterface
	// SetMetadata sets the JSON metadata and returns the relationship for chaining
	SetMetadata(metadata string) RelationshipInterface
	// SetCreatedAt sets the creation timestamp and returns the relationship for chaining
	SetCreatedAt(createdAt string) RelationshipInterface
}

// == TRASH INTERFACES =======================================================

// EntityTrashInterface defines the contract for trashed (soft-deleted) entities
// Trashed entities are moved to a separate table and can be restored.
type EntityTrashInterface interface {
	dataobject.DataObjectInterface

	// GetID returns the unique identifier for this trashed entity
	GetID() string
	// GetType returns the entity type
	GetType() string
	// GetHandle returns the entity handle
	GetHandle() string
	// GetCreatedAt returns the original creation timestamp in UTC
	GetCreatedAt() string
	// GetCreatedAtCarbon returns the original creation timestamp as a Carbon instance
	GetCreatedAtCarbon() *carbon.Carbon
	// GetUpdatedAt returns the last update timestamp in UTC
	GetUpdatedAt() string
	// GetUpdatedAtCarbon returns the last update timestamp as a Carbon instance
	GetUpdatedAtCarbon() *carbon.Carbon
	// GetDeletedAt returns the deletion timestamp in UTC
	GetDeletedAt() string
	// GetDeletedAtCarbon returns the deletion timestamp as a Carbon instance
	GetDeletedAtCarbon() *carbon.Carbon
	// GetDeletedBy returns the identifier of who/what deleted this entity
	GetDeletedBy() string

	// SetType sets the entity type and returns the trashed entity for chaining
	SetType(entityType string) EntityTrashInterface
	// SetHandle sets the entity handle and returns the trashed entity for chaining
	SetHandle(handle string) EntityTrashInterface
	// SetCreatedAt sets the original creation timestamp and returns the trashed entity for chaining
	SetCreatedAt(createdAt string) EntityTrashInterface
	// SetUpdatedAt sets the update timestamp and returns the trashed entity for chaining
	SetUpdatedAt(updatedAt string) EntityTrashInterface
	// SetDeletedAt sets the deletion timestamp and returns the trashed entity for chaining
	SetDeletedAt(deletedAt string) EntityTrashInterface
	// SetDeletedBy sets the deleter identifier and returns the trashed entity for chaining
	SetDeletedBy(deletedBy string) EntityTrashInterface
}

// AttributeTrashInterface defines the contract for trashed (soft-deleted) attributes
// Trashed attributes are moved to a separate table and can be restored.
type AttributeTrashInterface interface {
	dataobject.DataObjectInterface

	// GetID returns the unique identifier for this trashed attribute
	GetID() string
	// GetEntityID returns the ID of the entity this attribute belonged to
	GetEntityID() string
	// GetKey returns the attribute key/name
	GetKey() string
	// GetValue returns the attribute value
	GetValue() string
	// GetCreatedAt returns the original creation timestamp in UTC
	GetCreatedAt() string
	// GetCreatedAtCarbon returns the original creation timestamp as a Carbon instance
	GetCreatedAtCarbon() *carbon.Carbon
	// GetUpdatedAt returns the last update timestamp in UTC
	GetUpdatedAt() string
	// GetUpdatedAtCarbon returns the last update timestamp as a Carbon instance
	GetUpdatedAtCarbon() *carbon.Carbon
	// GetDeletedAt returns the deletion timestamp in UTC
	GetDeletedAt() string
	// GetDeletedAtCarbon returns the deletion timestamp as a Carbon instance
	GetDeletedAtCarbon() *carbon.Carbon
	// GetDeletedBy returns the identifier of who/what deleted this attribute
	GetDeletedBy() string

	// SetEntityID sets the associated entity ID and returns the trashed attribute for chaining
	SetEntityID(entityID string) AttributeTrashInterface
	// SetKey sets the attribute key and returns the trashed attribute for chaining
	SetKey(key string) AttributeTrashInterface
	// SetValue sets the attribute value and returns the trashed attribute for chaining
	SetValue(value string) AttributeTrashInterface
	// SetCreatedAt sets the original creation timestamp and returns the trashed attribute for chaining
	SetCreatedAt(createdAt string) AttributeTrashInterface
	// SetUpdatedAt sets the update timestamp and returns the trashed attribute for chaining
	SetUpdatedAt(updatedAt string) AttributeTrashInterface
	// SetDeletedAt sets the deletion timestamp and returns the trashed attribute for chaining
	SetDeletedAt(deletedAt string) AttributeTrashInterface
	// SetDeletedBy sets the deleter identifier and returns the trashed attribute for chaining
	SetDeletedBy(deletedBy string) AttributeTrashInterface
}

// RelationshipTrashInterface defines the contract for trashed (soft-deleted) relationships
// Trashed relationships are moved to a separate table and can be restored.
type RelationshipTrashInterface interface {
	dataobject.DataObjectInterface

	// GetID returns the unique identifier for this trashed relationship
	GetID() string
	// GetEntityID returns the ID of the source entity
	GetEntityID() string
	// GetRelatedEntityID returns the ID of the target/related entity
	GetRelatedEntityID() string
	// GetRelationshipType returns the relationship type
	GetRelationshipType() string
	// GetParentID returns the parent relationship ID
	GetParentID() string
	// GetSequence returns the sort order
	GetSequence() int
	// GetMetadata returns the JSON metadata
	GetMetadata() string
	// GetCreatedAt returns the original creation timestamp in UTC
	GetCreatedAt() string
	// GetCreatedAtCarbon returns the original creation timestamp as a Carbon instance
	GetCreatedAtCarbon() *carbon.Carbon
	// GetDeletedAt returns the deletion timestamp in UTC
	GetDeletedAt() string
	// GetDeletedAtCarbon returns the deletion timestamp as a Carbon instance
	GetDeletedAtCarbon() *carbon.Carbon
	// GetDeletedBy returns the identifier of who/what deleted this relationship
	GetDeletedBy() string

	// SetEntityID sets the source entity ID and returns the trashed relationship for chaining
	SetEntityID(entityID string) RelationshipTrashInterface
	// SetRelatedEntityID sets the target entity ID and returns the trashed relationship for chaining
	SetRelatedEntityID(relatedID string) RelationshipTrashInterface
	// SetRelationshipType sets the relationship type and returns the trashed relationship for chaining
	SetRelationshipType(relType string) RelationshipTrashInterface
	// SetParentID sets the parent relationship ID and returns the trashed relationship for chaining
	SetParentID(parentID string) RelationshipTrashInterface
	// SetSequence sets the sort order and returns the trashed relationship for chaining
	SetSequence(sequence int) RelationshipTrashInterface
	// SetMetadata sets the JSON metadata and returns the trashed relationship for chaining
	SetMetadata(metadata string) RelationshipTrashInterface
	// SetCreatedAt sets the original creation timestamp and returns the trashed relationship for chaining
	SetCreatedAt(createdAt string) RelationshipTrashInterface
	// SetDeletedAt sets the deletion timestamp and returns the trashed relationship for chaining
	SetDeletedAt(deletedAt string) RelationshipTrashInterface
	// SetDeletedBy sets the deleter identifier and returns the trashed relationship for chaining
	SetDeletedBy(deletedBy string) RelationshipTrashInterface
}

// == TAXONOMY INTERFACES =====================================================

// TaxonomyInterface defines the contract for taxonomies (classification systems)
// Taxonomies are used to categorize entities (e.g., "Categories", "Tags", "Topics")
// with support for hierarchical terms and entity type restrictions.
type TaxonomyInterface interface {
	dataobject.DataObjectInterface

	// GetID returns the unique identifier for this taxonomy
	GetID() string
	// GetName returns the display name of the taxonomy
	GetName() string
	// GetSlug returns the URL-friendly identifier for the taxonomy
	GetSlug() string
	// GetDescription returns the optional description of the taxonomy
	GetDescription() string
	// GetParentID returns the parent taxonomy ID for hierarchical taxonomies
	GetParentID() string
	// GetEntityTypes returns the list of entity types that can use this taxonomy
	GetEntityTypes() []string
	// GetCreatedAt returns the creation timestamp in UTC
	GetCreatedAt() string
	// GetCreatedAtCarbon returns the creation timestamp as a Carbon instance
	GetCreatedAtCarbon() *carbon.Carbon
	// GetUpdatedAt returns the last update timestamp in UTC
	GetUpdatedAt() string
	// GetUpdatedAtCarbon returns the last update timestamp as a Carbon instance
	GetUpdatedAtCarbon() *carbon.Carbon

	// SetName sets the display name and returns the taxonomy for chaining
	SetName(name string) TaxonomyInterface
	// SetSlug sets the URL-friendly identifier and returns the taxonomy for chaining
	SetSlug(slug string) TaxonomyInterface
	// SetDescription sets the description and returns the taxonomy for chaining
	SetDescription(desc string) TaxonomyInterface
	// SetParentID sets the parent taxonomy ID and returns the taxonomy for chaining
	SetParentID(parentID string) TaxonomyInterface
	// SetEntityTypes sets the allowed entity types and returns the taxonomy for chaining
	SetEntityTypes(types []string) TaxonomyInterface
	// SetCreatedAt sets the creation timestamp and returns the taxonomy for chaining
	SetCreatedAt(createdAt string) TaxonomyInterface
	// SetUpdatedAt sets the update timestamp and returns the taxonomy for chaining
	SetUpdatedAt(updatedAt string) TaxonomyInterface
}

// TaxonomyTermInterface defines the contract for taxonomy terms
// Terms are the actual categories/tags within a taxonomy (e.g., "Electronics" in "Categories")
// with support for hierarchical parent-child relationships.
type TaxonomyTermInterface interface {
	dataobject.DataObjectInterface

	// GetID returns the unique identifier for this taxonomy term
	GetID() string
	// GetTaxonomyID returns the ID of the taxonomy this term belongs to
	GetTaxonomyID() string
	// GetName returns the display name of the term
	GetName() string
	// GetSlug returns the URL-friendly identifier for the term
	GetSlug() string
	// GetParentID returns the parent term ID for hierarchical terms
	GetParentID() string
	// GetSortOrder returns the sort order for this term
	GetSortOrder() int
	// GetCreatedAt returns the creation timestamp in UTC
	GetCreatedAt() string
	// GetCreatedAtCarbon returns the creation timestamp as a Carbon instance
	GetCreatedAtCarbon() *carbon.Carbon
	// GetUpdatedAt returns the last update timestamp in UTC
	GetUpdatedAt() string
	// GetUpdatedAtCarbon returns the last update timestamp as a Carbon instance
	GetUpdatedAtCarbon() *carbon.Carbon

	// SetTaxonomyID sets the parent taxonomy ID and returns the term for chaining
	SetTaxonomyID(taxonomyID string) TaxonomyTermInterface
	// SetName sets the display name and returns the term for chaining
	SetName(name string) TaxonomyTermInterface
	// SetSlug sets the URL-friendly identifier and returns the term for chaining
	SetSlug(slug string) TaxonomyTermInterface
	// SetParentID sets the parent term ID and returns the term for chaining
	SetParentID(parentID string) TaxonomyTermInterface
	// SetSortOrder sets the sort order and returns the term for chaining
	SetSortOrder(order int) TaxonomyTermInterface
	// SetCreatedAt sets the creation timestamp and returns the term for chaining
	SetCreatedAt(createdAt string) TaxonomyTermInterface
	// SetUpdatedAt sets the update timestamp and returns the term for chaining
	SetUpdatedAt(updatedAt string) TaxonomyTermInterface
}

// EntityTaxonomyInterface defines the contract for entity-taxonomy assignments
// Links entities to taxonomy terms for categorization (e.g., Product X is in "Electronics" category).
type EntityTaxonomyInterface interface {
	dataobject.DataObjectInterface

	// GetID returns the unique identifier for this entity-taxonomy assignment
	GetID() string
	// GetEntityID returns the ID of the categorized entity
	GetEntityID() string
	// GetTaxonomyID returns the ID of the taxonomy
	GetTaxonomyID() string
	// GetTermID returns the ID of the taxonomy term assigned to the entity
	GetTermID() string
	// GetCreatedAt returns the creation timestamp in UTC
	GetCreatedAt() string
	// GetCreatedAtCarbon returns the creation timestamp as a Carbon instance
	GetCreatedAtCarbon() *carbon.Carbon

	// SetEntityID sets the entity ID and returns the assignment for chaining
	SetEntityID(entityID string) EntityTaxonomyInterface
	// SetTaxonomyID sets the taxonomy ID and returns the assignment for chaining
	SetTaxonomyID(taxonomyID string) EntityTaxonomyInterface
	// SetTermID sets the term ID and returns the assignment for chaining
	SetTermID(termID string) EntityTaxonomyInterface
	// SetCreatedAt sets the creation timestamp and returns the assignment for chaining
	SetCreatedAt(createdAt string) EntityTaxonomyInterface
}

// TaxonomyTrashInterface defines the contract for trashed (soft-deleted) taxonomies
// Trashed taxonomies are moved to a separate table and can be restored.
type TaxonomyTrashInterface interface {
	dataobject.DataObjectInterface

	// GetID returns the unique identifier for this trashed taxonomy
	GetID() string
	// GetName returns the display name
	GetName() string
	// GetSlug returns the URL-friendly identifier
	GetSlug() string
	// GetDescription returns the description
	GetDescription() string
	// GetParentID returns the parent taxonomy ID
	GetParentID() string
	// GetEntityTypes returns the list of allowed entity types
	GetEntityTypes() []string
	// GetCreatedAt returns the original creation timestamp in UTC
	GetCreatedAt() string
	// GetCreatedAtCarbon returns the original creation timestamp as a Carbon instance
	GetCreatedAtCarbon() *carbon.Carbon
	// GetUpdatedAt returns the last update timestamp in UTC
	GetUpdatedAt() string
	// GetUpdatedAtCarbon returns the last update timestamp as a Carbon instance
	GetUpdatedAtCarbon() *carbon.Carbon
	// GetDeletedAt returns the deletion timestamp in UTC
	GetDeletedAt() string
	// GetDeletedAtCarbon returns the deletion timestamp as a Carbon instance
	GetDeletedAtCarbon() *carbon.Carbon
	// GetDeletedBy returns the identifier of who/what deleted this taxonomy
	GetDeletedBy() string

	// SetName sets the display name and returns the trashed taxonomy for chaining
	SetName(name string) TaxonomyTrashInterface
	// SetSlug sets the URL-friendly identifier and returns the trashed taxonomy for chaining
	SetSlug(slug string) TaxonomyTrashInterface
	// SetDescription sets the description and returns the trashed taxonomy for chaining
	SetDescription(desc string) TaxonomyTrashInterface
	// SetParentID sets the parent taxonomy ID and returns the trashed taxonomy for chaining
	SetParentID(parentID string) TaxonomyTrashInterface
	// SetEntityTypes sets the allowed entity types and returns the trashed taxonomy for chaining
	SetEntityTypes(types []string) TaxonomyTrashInterface
	// SetCreatedAt sets the original creation timestamp and returns the trashed taxonomy for chaining
	SetCreatedAt(createdAt string) TaxonomyTrashInterface
	// SetUpdatedAt sets the update timestamp and returns the trashed taxonomy for chaining
	SetUpdatedAt(updatedAt string) TaxonomyTrashInterface
	// SetDeletedAt sets the deletion timestamp and returns the trashed taxonomy for chaining
	SetDeletedAt(deletedAt string) TaxonomyTrashInterface
	// SetDeletedBy sets the deleter identifier and returns the trashed taxonomy for chaining
	SetDeletedBy(deletedBy string) TaxonomyTrashInterface
}

// TaxonomyTermTrashInterface defines the contract for trashed (soft-deleted) taxonomy terms
// Trashed taxonomy terms are moved to a separate table and can be restored.
type TaxonomyTermTrashInterface interface {
	dataobject.DataObjectInterface

	// GetID returns the unique identifier for this trashed taxonomy term
	GetID() string
	// GetTaxonomyID returns the ID of the taxonomy this term belonged to
	GetTaxonomyID() string
	// GetName returns the display name
	GetName() string
	// GetSlug returns the URL-friendly identifier
	GetSlug() string
	// GetParentID returns the parent term ID
	GetParentID() string
	// GetSortOrder returns the sort order
	GetSortOrder() int
	// GetCreatedAt returns the original creation timestamp in UTC
	GetCreatedAt() string
	// GetCreatedAtCarbon returns the original creation timestamp as a Carbon instance
	GetCreatedAtCarbon() *carbon.Carbon
	// GetUpdatedAt returns the last update timestamp in UTC
	GetUpdatedAt() string
	// GetUpdatedAtCarbon returns the last update timestamp as a Carbon instance
	GetUpdatedAtCarbon() *carbon.Carbon
	// GetDeletedAt returns the deletion timestamp in UTC
	GetDeletedAt() string
	// GetDeletedAtCarbon returns the deletion timestamp as a Carbon instance
	GetDeletedAtCarbon() *carbon.Carbon
	// GetDeletedBy returns the identifier of who/what deleted this term
	GetDeletedBy() string

	// SetTaxonomyID sets the taxonomy ID and returns the trashed term for chaining
	SetTaxonomyID(taxonomyID string) TaxonomyTermTrashInterface
	// SetName sets the display name and returns the trashed term for chaining
	SetName(name string) TaxonomyTermTrashInterface
	// SetSlug sets the URL-friendly identifier and returns the trashed term for chaining
	SetSlug(slug string) TaxonomyTermTrashInterface
	// SetParentID sets the parent term ID and returns the trashed term for chaining
	SetParentID(parentID string) TaxonomyTermTrashInterface
	// SetSortOrder sets the sort order and returns the trashed term for chaining
	SetSortOrder(order int) TaxonomyTermTrashInterface
	// SetCreatedAt sets the original creation timestamp and returns the trashed term for chaining
	SetCreatedAt(createdAt string) TaxonomyTermTrashInterface
	// SetUpdatedAt sets the update timestamp and returns the trashed term for chaining
	SetUpdatedAt(updatedAt string) TaxonomyTermTrashInterface
	// SetDeletedAt sets the deletion timestamp and returns the trashed term for chaining
	SetDeletedAt(deletedAt string) TaxonomyTermTrashInterface
	// SetDeletedBy sets the deleter identifier and returns the trashed term for chaining
	SetDeletedBy(deletedBy string) TaxonomyTermTrashInterface
}

// == STORE INTERFACE ========================================================

// StoreInterface defines the contract for the entity store
// It provides CRUD operations for entities, attributes, relationships, and taxonomies.
// All methods accept a context.Context for cancellation and timeout control.
type StoreInterface interface {
	// AutoMigrate creates or updates database tables to match the current schema
	AutoMigrate(ctx context.Context) error

	// GetAttributeTableName returns the configured attributes table name
	GetAttributeTableName() string
	// GetAttributeTrashTableName returns the configured trashed attributes table name
	GetAttributeTrashTableName() string
	// GetDB returns the underlying *sql.DB connection
	GetDB() *sql.DB
	// GetEntityTableName returns the configured entities table name
	GetEntityTableName() string
	// GetEntityTrashTableName returns the configured trashed entities table name
	GetEntityTrashTableName() string
	// GetRelationshipTableName returns the configured relationships table name
	GetRelationshipTableName() string
	// GetRelationshipTrashTableName returns the configured trashed relationships table name
	GetRelationshipTrashTableName() string

	// AttributeCreate persists a new attribute record
	AttributeCreate(ctx context.Context, attr AttributeInterface) error
	// AttributeCreateWithKeyAndValue creates an attribute with the given key/value for an entity
	AttributeCreateWithKeyAndValue(ctx context.Context, entityID string, attributeKey string, attributeValue string) (AttributeInterface, error)
	// AttributeFind retrieves an attribute by entity ID and key
	AttributeFind(ctx context.Context, entityID string, attributeKey string) (AttributeInterface, error)
	// AttributeFindByHandle retrieves an attribute by entity type, handle, and attribute key
	AttributeFindByHandle(ctx context.Context, entityType string, entityHandle string, attributeKey string) (AttributeInterface, error)
	// AttributeList retrieves attributes matching the given query options
	AttributeList(ctx context.Context, options AttributeQueryOptions) ([]AttributeInterface, error)
	// AttributesSet creates or updates multiple attributes for an entity at once
	AttributesSet(ctx context.Context, entityID string, attributes map[string]string) error
	// AttributeSetFloat stores a float64 value as an attribute
	AttributeSetFloat(ctx context.Context, entityID string, attributeKey string, attributeValue float64) error
	// AttributeSetInt stores an int64 value as an attribute
	AttributeSetInt(ctx context.Context, entityID string, attributeKey string, attributeValue int64) error
	// AttributeSetString stores a string value as an attribute
	AttributeSetString(ctx context.Context, entityID string, attributeKey string, attributeValue string) error
	// AttributeUpdate updates an existing attribute record
	AttributeUpdate(ctx context.Context, attr AttributeInterface) error

	// EntityAttributeList retrieves all attributes for a given entity
	EntityAttributeList(ctx context.Context, entityID string) ([]AttributeInterface, error)
	// EntityCount counts entities matching the given query options
	EntityCount(ctx context.Context, options EntityQueryOptions) (int64, error)
	// EntityCreate persists a new entity record
	EntityCreate(ctx context.Context, entity EntityInterface) error
	// EntityCreateWithType creates a new entity with the given type
	EntityCreateWithType(ctx context.Context, entityType string) (EntityInterface, error)
	// EntityCreateWithTypeAndAttributes creates a new entity with the given type and attributes
	EntityCreateWithTypeAndAttributes(ctx context.Context, entityType string, attributes map[string]string) (EntityInterface, error)
	// EntityDelete permanently removes an entity by ID
	EntityDelete(ctx context.Context, entityID string) (bool, error)
	// EntityFindByAttribute finds an entity by type and attribute key/value
	EntityFindByAttribute(ctx context.Context, entityType string, attributeKey string, attributeValue string) (EntityInterface, error)
	// EntityFindByHandle finds an entity by its type and handle
	EntityFindByHandle(ctx context.Context, entityType string, entityHandle string) (EntityInterface, error)
	// EntityFindByID finds an entity by its unique ID
	EntityFindByID(ctx context.Context, entityID string) (EntityInterface, error)
	// EntityList retrieves entities matching the given query options
	EntityList(ctx context.Context, options EntityQueryOptions) ([]EntityInterface, error)
	// EntityListByAttribute finds all entities of a type with a specific attribute value
	EntityListByAttribute(ctx context.Context, entityType string, attributeKey string, attributeValue string) ([]EntityInterface, error)
	// EntityTrash soft-deletes an entity by moving it to the trash table
	EntityTrash(ctx context.Context, entityID string) (bool, error)
	// EntityUpdate updates an existing entity record
	EntityUpdate(ctx context.Context, entity EntityInterface) error

	// RelationshipCreate persists a new relationship record
	RelationshipCreate(ctx context.Context, relationship RelationshipInterface) error
	// RelationshipCreateByOptions creates a relationship using the provided options
	RelationshipCreateByOptions(ctx context.Context, options RelationshipOptions) (RelationshipInterface, error)
	// RelationshipCount counts relationships matching the given query options
	RelationshipCount(ctx context.Context, options RelationshipQueryOptions) (int64, error)
	// RelationshipDelete permanently removes a relationship by ID
	RelationshipDelete(ctx context.Context, relationshipID string) (bool, error)
	// RelationshipDeleteAll removes all relationships for a given entity
	RelationshipDeleteAll(ctx context.Context, entityID string) error
	// RelationshipFind retrieves a relationship by its ID
	RelationshipFind(ctx context.Context, relationshipID string) (RelationshipInterface, error)
	// RelationshipFindByEntities finds a relationship by source, target, and type
	RelationshipFindByEntities(ctx context.Context, entityID string, relatedEntityID string, relationshipType string) (RelationshipInterface, error)
	// RelationshipList retrieves relationships matching the given query options
	RelationshipList(ctx context.Context, options RelationshipQueryOptions) ([]RelationshipInterface, error)
	// RelationshipListRelated retrieves relationships where the given entity is the target
	RelationshipListRelated(ctx context.Context, relatedEntityID string, relationshipType string) ([]RelationshipInterface, error)
	// RelationshipRestore restores a trashed relationship
	RelationshipRestore(ctx context.Context, relationshipID string) (bool, error)
	// RelationshipTrash soft-deletes a relationship by moving it to the trash table
	RelationshipTrash(ctx context.Context, relationshipID string, deletedBy string) (bool, error)
	// RelationshipTrashList retrieves trashed relationships matching the query options
	RelationshipTrashList(ctx context.Context, options RelationshipQueryOptions) ([]RelationshipTrashInterface, error)

	// TaxonomyCreate persists a new taxonomy record
	TaxonomyCreate(ctx context.Context, taxonomy TaxonomyInterface) error
	// TaxonomyCreateByOptions creates a taxonomy using the provided options
	TaxonomyCreateByOptions(ctx context.Context, options TaxonomyOptions) (TaxonomyInterface, error)
	// TaxonomyCount counts taxonomies matching the given query options
	TaxonomyCount(ctx context.Context, options TaxonomyQueryOptions) (int64, error)
	// TaxonomyDelete permanently removes a taxonomy by ID
	TaxonomyDelete(ctx context.Context, taxonomyID string) (bool, error)
	// TaxonomyFind retrieves a taxonomy by its ID
	TaxonomyFind(ctx context.Context, taxonomyID string) (TaxonomyInterface, error)
	// TaxonomyFindBySlug finds a taxonomy by its slug
	TaxonomyFindBySlug(ctx context.Context, slug string) (TaxonomyInterface, error)
	// TaxonomyList retrieves taxonomies matching the given query options
	TaxonomyList(ctx context.Context, options TaxonomyQueryOptions) ([]TaxonomyInterface, error)
	// TaxonomyRestore restores a trashed taxonomy
	TaxonomyRestore(ctx context.Context, taxonomyID string) (bool, error)
	// TaxonomyTrash soft-deletes a taxonomy by moving it to the trash table
	TaxonomyTrash(ctx context.Context, taxonomyID string, deletedBy string) (bool, error)
	// TaxonomyTrashList retrieves trashed taxonomies matching the query options
	TaxonomyTrashList(ctx context.Context, options TaxonomyQueryOptions) ([]TaxonomyTrashInterface, error)
	// TaxonomyUpdate updates an existing taxonomy record
	TaxonomyUpdate(ctx context.Context, taxonomy TaxonomyInterface) error

	// TaxonomyTermCreate persists a new taxonomy term record
	TaxonomyTermCreate(ctx context.Context, term TaxonomyTermInterface) error
	// TaxonomyTermCreateByOptions creates a taxonomy term using the provided options
	TaxonomyTermCreateByOptions(ctx context.Context, options TaxonomyTermOptions) (TaxonomyTermInterface, error)
	// TaxonomyTermCount counts taxonomy terms matching the given query options
	TaxonomyTermCount(ctx context.Context, options TaxonomyTermQueryOptions) (int64, error)
	// TaxonomyTermDelete permanently removes a taxonomy term by ID
	TaxonomyTermDelete(ctx context.Context, termID string) (bool, error)
	// TaxonomyTermFind retrieves a taxonomy term by its ID
	TaxonomyTermFind(ctx context.Context, termID string) (TaxonomyTermInterface, error)
	// TaxonomyTermFindBySlug finds a taxonomy term by its taxonomy ID and slug
	TaxonomyTermFindBySlug(ctx context.Context, taxonomyID string, slug string) (TaxonomyTermInterface, error)
	// TaxonomyTermList retrieves taxonomy terms matching the given query options
	TaxonomyTermList(ctx context.Context, options TaxonomyTermQueryOptions) ([]TaxonomyTermInterface, error)
	// TaxonomyTermRestore restores a trashed taxonomy term
	TaxonomyTermRestore(ctx context.Context, termID string) (bool, error)
	// TaxonomyTermTrash soft-deletes a taxonomy term by moving it to the trash table
	TaxonomyTermTrash(ctx context.Context, termID string, deletedBy string) (bool, error)
	// TaxonomyTermTrashList retrieves trashed taxonomy terms matching the query options
	TaxonomyTermTrashList(ctx context.Context, options TaxonomyTermQueryOptions) ([]TaxonomyTermTrashInterface, error)
	// TaxonomyTermUpdate updates an existing taxonomy term record
	TaxonomyTermUpdate(ctx context.Context, term TaxonomyTermInterface) error

	// EntityTaxonomyAssign assigns an entity to a taxonomy term
	EntityTaxonomyAssign(ctx context.Context, entityID string, taxonomyID string, termID string) error
	// EntityTaxonomyCount counts entity-taxonomy assignments matching the given query options
	EntityTaxonomyCount(ctx context.Context, options EntityTaxonomyQueryOptions) (int64, error)
	// EntityTaxonomyList retrieves entity-taxonomy assignments matching the given query options
	EntityTaxonomyList(ctx context.Context, options EntityTaxonomyQueryOptions) ([]EntityTaxonomyInterface, error)
	// EntityTaxonomyRemove removes an entity from a taxonomy term
	EntityTaxonomyRemove(ctx context.Context, entityID string, taxonomyID string, termID string) error

	// GetTaxonomyTableName returns the configured taxonomies table name
	GetTaxonomyTableName() string
	// GetTaxonomyTrashTableName returns the configured trashed taxonomies table name
	GetTaxonomyTrashTableName() string
	// GetTaxonomyTermTableName returns the configured taxonomy terms table name
	GetTaxonomyTermTableName() string
	// GetTaxonomyTermTrashTableName returns the configured trashed taxonomy terms table name
	GetTaxonomyTermTrashTableName() string
	// GetEntityTaxonomyTableName returns the configured entity-taxonomy assignments table name
	GetEntityTaxonomyTableName() string
}
