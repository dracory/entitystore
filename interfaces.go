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
	EntityType() string
	EntityHandle() string
	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon

	// Core setters (fluent) — ID() / SetID() come from DataObjectInterface
	SetEntityType(entityType string) EntityInterface
	SetEntityHandle(handle string) EntityInterface
	SetCreatedAt(createdAt string) EntityInterface
	SetUpdatedAt(updatedAt string) EntityInterface

	// Dynamic / extra attributes (in-memory)
	GetAttribute(key string) string
	SetAttribute(key string, value string) EntityInterface
	GetAllAttributes() map[string]string
}

// == ATTRIBUTE INTERFACE ====================================================

// AttributeInterface defines the contract for persisted entity key-value attributes
type AttributeInterface interface {
	dataobject.DataObjectInterface

	// Core getters
	EntityID() string
	AttributeKey() string
	AttributeValue() string
	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon

	// Core setters (fluent) — ID() / SetID() come from DataObjectInterface
	SetEntityID(entityID string) AttributeInterface
	SetAttributeKey(key string) AttributeInterface
	SetAttributeValue(value string) AttributeInterface
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
	EntityID() string
	RelatedEntityID() string
	RelationshipType() string
	ParentID() string
	Sequence() int
	Metadata() string
	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon

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
	EntityType() string
	EntityHandle() string
	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon
	DeletedAt() string
	DeletedAtCarbon() *carbon.Carbon
	DeletedBy() string

	// Core setters (fluent) — ID() / SetID() come from DataObjectInterface
	SetEntityType(entityType string) EntityTrashInterface
	SetEntityHandle(handle string) EntityTrashInterface
	SetCreatedAt(createdAt string) EntityTrashInterface
	SetUpdatedAt(updatedAt string) EntityTrashInterface
	SetDeletedAt(deletedAt string) EntityTrashInterface
	SetDeletedBy(deletedBy string) EntityTrashInterface
}

// AttributeTrashInterface defines the contract for trashed attributes
type AttributeTrashInterface interface {
	dataobject.DataObjectInterface

	// Core getters
	EntityID() string
	AttributeKey() string
	AttributeValue() string
	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon
	DeletedAt() string
	DeletedAtCarbon() *carbon.Carbon
	DeletedBy() string

	// Core setters (fluent) — ID() / SetID() come from DataObjectInterface
	SetEntityID(entityID string) AttributeTrashInterface
	SetAttributeKey(key string) AttributeTrashInterface
	SetAttributeValue(value string) AttributeTrashInterface
	SetCreatedAt(createdAt string) AttributeTrashInterface
	SetUpdatedAt(updatedAt string) AttributeTrashInterface
	SetDeletedAt(deletedAt string) AttributeTrashInterface
	SetDeletedBy(deletedBy string) AttributeTrashInterface
}

// RelationshipTrashInterface defines the contract for trashed relationships
type RelationshipTrashInterface interface {
	dataobject.DataObjectInterface

	// Core getters
	EntityID() string
	RelatedEntityID() string
	RelationshipType() string
	ParentID() string
	Sequence() int
	Metadata() string
	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	DeletedAt() string
	DeletedAtCarbon() *carbon.Carbon
	DeletedBy() string

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
}
