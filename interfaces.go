package entitystore

import (
	"context"
	"database/sql"
)

type StoreInterface interface {
	AutoMigrate(ctx context.Context) error

	GetAttributeTableName() string
	GetAttributeTrashTableName() string
	GetDB() *sql.DB
	GetEntityTableName() string
	GetEntityTrashTableName() string

	// AttributeCount(entityID string) uint64
	AttributeCreate(ctx context.Context, attr *Attribute) error
	AttributeCreateWithKeyAndValue(ctx context.Context, entityID string, attributeKey string, attributeValue string) (*Attribute, error)
	AttributeFind(ctx context.Context, entityID string, attributeKey string) (*Attribute, error)
	AttributeFindByHandle(ctx context.Context, entityID string, attributeKey string, attributeValue string) (*Attribute, error)
	AttributeList(ctx context.Context, options AttributeQueryOptions) ([]Attribute, error)
	AttributesSet(ctx context.Context, entityID string, attributes map[string]string) error
	AttributeSetFloat(ctx context.Context, entityID string, attributeKey string, attributeValue float64) error
	AttributeSetInt(ctx context.Context, entityID string, attributeKey string, attributeValue int64) error
	AttributeSetString(ctx context.Context, entityID string, attributeKey string, attributeValue string) error
	// AttributeTrash(attr *Attribute) error

	EntityAttributeList(ctx context.Context, entityID string) ([]Attribute, error)
	EntityCount(ctx context.Context, options EntityQueryOptions) (int64, error)
	EntityCreate(ctx context.Context, entity *Entity) error
	EntityCreateWithType(ctx context.Context, entityType string) (*Entity, error)
	EntityCreateWithTypeAndAttributes(ctx context.Context, entityType string, attributes map[string]string) (*Entity, error)
	EntityDelete(ctx context.Context, entityID string) (bool, error)
	EntityFindByAttribute(ctx context.Context, entityType string, attributeKey string, attributeValue string) (*Entity, error)
	EntityFindByHandle(ctx context.Context, entityType string, entityHandle string) (*Entity, error)
	EntityFindByID(ctx context.Context, entityID string) (*Entity, error)
	EntityList(ctx context.Context, options EntityQueryOptions) ([]Entity, error)
	EntityListByAttribute(ctx context.Context, entityType string, attributeKey string, attributeValue string) ([]Entity, error)
	EntityTrash(ctx context.Context, entityID string) (bool, error)
	EntityUpdate(ctx context.Context, entity Entity) error

	NewAttribute(opts NewAttributeOptions) Attribute
	NewAttributeFromMap(entityMap map[string]string) Attribute

	NewEntity(opts NewEntityOptions) Entity
	NewEntityFromMap(entityMap map[string]string) Entity
}
