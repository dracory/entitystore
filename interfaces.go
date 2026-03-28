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
	Status() string
	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon
	SoftDeletedAt() string
	SoftDeletedAtCarbon() *carbon.Carbon

	// Core setters (fluent) — ID() / SetID() come from DataObjectInterface
	SetEntityType(entityType string) EntityInterface
	SetEntityHandle(handle string) EntityInterface
	SetStatus(status string) EntityInterface
	SetCreatedAt(createdAt string) EntityInterface
	SetUpdatedAt(updatedAt string) EntityInterface
	SetSoftDeletedAt(softDeletedAt string) EntityInterface

	// Dynamic / extra attributes (in-memory)
	GetAttribute(key string) string
	SetAttribute(key string, value string) EntityInterface
	GetAllAttributes() map[string]string

	// Status helpers
	IsActive() bool
	IsInactive() bool
	IsSoftDeleted() bool
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

// == STORE INTERFACE ========================================================

type StoreInterface interface {
	AutoMigrate(ctx context.Context) error

	GetAttributeTableName() string
	GetAttributeTrashTableName() string
	GetDB() *sql.DB
	GetEntityTableName() string
	GetEntityTrashTableName() string

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
}
