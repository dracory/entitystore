# API Reference

Complete reference for Entity Store interfaces and methods.

## StoreInterface

Main entry point for all operations.

### Entity Operations

```go
// Create
EntityCreate(ctx context.Context, entity EntityInterface) error
EntityCreateWithType(ctx context.Context, entityType string) (EntityInterface, error)
EntityCreateWithTypeAndAttributes(ctx context.Context, entityType string, attributes map[string]string) (EntityInterface, error)

// Read
EntityFindByID(ctx context.Context, id string) (EntityInterface, error)
EntityFindByHandle(ctx context.Context, entityType string, entityHandle string) (EntityInterface, error)
EntityFindByAttribute(ctx context.Context, entityType string, attributeKey string, attributeValue string) (EntityInterface, error)
EntityList(ctx context.Context, options EntityQueryOptions) ([]EntityInterface, error)
EntityListByAttribute(ctx context.Context, entityType string, attributeKey string, attributeValue string) ([]EntityInterface, error)
EntityCount(ctx context.Context, options EntityQueryOptions) (int64, error)
EntityAttributeList(ctx context.Context, entityID string) ([]AttributeInterface, error)

// Update
EntityUpdate(ctx context.Context, entity EntityInterface) error

// Delete
EntityDelete(ctx context.Context, id string) (bool, error)
EntityTrash(ctx context.Context, id string) (bool, error)
EntityRestore(ctx context.Context, id string) (bool, error)
EntityTrashList(ctx context.Context, options EntityQueryOptions) ([]EntityTrashInterface, error)
```

### Attribute Operations

```go
// CRUD
AttributeCreate(ctx context.Context, attr AttributeInterface) error
AttributeFind(ctx context.Context, entityID string, attributeKey string) (AttributeInterface, error)
AttributeFindByHandle(ctx context.Context, entityType string, entityHandle string, attributeKey string) (AttributeInterface, error)
AttributeList(ctx context.Context, options AttributeQueryOptions) ([]AttributeInterface, error)
AttributeCount(ctx context.Context, options AttributeQueryOptions) (int64, error)
AttributeUpdate(ctx context.Context, attr AttributeInterface) error
AttributeDelete(ctx context.Context, id string) (bool, error)
AttributeTrash(ctx context.Context, id string) (bool, error)
AttributeRestore(ctx context.Context, id string) (bool, error)
AttributeTrashList(ctx context.Context, options AttributeQueryOptions) ([]AttributeTrashInterface, error)

// Shortcuts
AttributeCreateWithKeyAndValue(ctx context.Context, entityID string, attributeKey string, attributeValue string) (AttributeInterface, error)
AttributeSetString(ctx context.Context, entityID string, attributeKey string, attributeValue string) error
AttributeSetInt(ctx context.Context, entityID string, attributeKey string, attributeValue int64) error
AttributeSetFloat(ctx context.Context, entityID string, attributeKey string, attributeValue float64) error
AttributeSetInterface(ctx context.Context, entityID string, attributeKey string, attributeValue interface{}) error
AttributesSet(ctx context.Context, entityID string, attributes map[string]string) error
```

### Relationship Operations (Requires `RelationshipsEnabled: true`)

```go
// CRUD
RelationshipCreate(ctx context.Context, relationship RelationshipInterface) error
RelationshipCreateByOptions(ctx context.Context, options RelationshipOptions) (RelationshipInterface, error)
RelationshipFind(ctx context.Context, relationshipID string) (RelationshipInterface, error)
RelationshipFindByEntities(ctx context.Context, entityID string, relatedEntityID string, relationshipType string) (RelationshipInterface, error)
RelationshipList(ctx context.Context, options RelationshipQueryOptions) ([]RelationshipInterface, error)
RelationshipCount(ctx context.Context, options RelationshipQueryOptions) (int64, error)
RelationshipDelete(ctx context.Context, relationshipID string) (bool, error)
RelationshipDeleteAll(ctx context.Context, entityID string) error

// Trash
RelationshipTrash(ctx context.Context, relationshipID string, deletedBy string) (bool, error)
RelationshipRestore(ctx context.Context, relationshipID string) (bool, error)
RelationshipTrashList(ctx context.Context, options RelationshipQueryOptions) ([]RelationshipTrashInterface, error)

// Queries
RelationshipListRelated(ctx context.Context, relatedEntityID string, relationshipType string) ([]RelationshipInterface, error)
```

### Taxonomy Operations (Requires `TaxonomiesEnabled: true`)

```go
// Taxonomy CRUD
TaxonomyCreate(ctx context.Context, taxonomy TaxonomyInterface) error
TaxonomyCreateByOptions(ctx context.Context, options TaxonomyOptions) (TaxonomyInterface, error)
TaxonomyFind(ctx context.Context, taxonomyID string) (TaxonomyInterface, error)
TaxonomyFindBySlug(ctx context.Context, slug string) (TaxonomyInterface, error)
TaxonomyList(ctx context.Context, options TaxonomyQueryOptions) ([]TaxonomyInterface, error)
TaxonomyCount(ctx context.Context, options TaxonomyQueryOptions) (int64, error)
TaxonomyUpdate(ctx context.Context, taxonomy TaxonomyInterface) error
TaxonomyDelete(ctx context.Context, taxonomyID string) (bool, error)
TaxonomyTrash(ctx context.Context, taxonomyID string, deletedBy string) (bool, error)
TaxonomyRestore(ctx context.Context, taxonomyID string) (bool, error)
TaxonomyTrashList(ctx context.Context, options TaxonomyQueryOptions) ([]TaxonomyTrashInterface, error)

// TaxonomyTerm CRUD
TaxonomyTermCreate(ctx context.Context, term TaxonomyTermInterface) error
TaxonomyTermCreateByOptions(ctx context.Context, options TaxonomyTermOptions) (TaxonomyTermInterface, error)
TaxonomyTermFind(ctx context.Context, termID string) (TaxonomyTermInterface, error)
TaxonomyTermFindBySlug(ctx context.Context, taxonomyID string, slug string) (TaxonomyTermInterface, error)
TaxonomyTermList(ctx context.Context, options TaxonomyTermQueryOptions) ([]TaxonomyTermInterface, error)
TaxonomyTermCount(ctx context.Context, options TaxonomyTermQueryOptions) (int64, error)
TaxonomyTermUpdate(ctx context.Context, term TaxonomyTermInterface) error
TaxonomyTermDelete(ctx context.Context, termID string) (bool, error)
TaxonomyTermTrash(ctx context.Context, termID string, deletedBy string) (bool, error)
TaxonomyTermRestore(ctx context.Context, termID string) (bool, error)
TaxonomyTermTrashList(ctx context.Context, options TaxonomyTermQueryOptions) ([]TaxonomyTermTrashInterface, error)

// Entity Assignment
EntityTaxonomyAssign(ctx context.Context, entityID string, taxonomyID string, termID string) error
EntityTaxonomyRemove(ctx context.Context, entityID string, taxonomyID string, termID string) error
EntityTaxonomyList(ctx context.Context, options EntityTaxonomyQueryOptions) ([]EntityTaxonomyInterface, error)
EntityTaxonomyCount(ctx context.Context, options EntityTaxonomyQueryOptions) (int64, error)
```

### Utility Methods

```go
// Table names
GetAttributeTableName() string
GetAttributeTrashTableName() string
GetDB() *sql.DB
GetEntityTableName() string
GetEntityTrashTableName() string
GetRelationshipTableName() string
GetRelationshipTrashTableName() string
GetTaxonomyTableName() string
GetTaxonomyTrashTableName() string
GetTaxonomyTermTableName() string
GetTaxonomyTermTrashTableName() string
GetEntityTaxonomyTableName() string

// Migration
AutoMigrate(ctx context.Context) error
```

## EntityInterface

Primary domain object interface.

```go
type EntityInterface interface {
    dataobject.DataObjectInterface
    
    // Core getters
    EntityType() string
    EntityHandle() string
    CreatedAt() string
    CreatedAtCarbon() *carbon.Carbon
    UpdatedAt() string
    UpdatedAtCarbon() *carbon.Carbon
    
    // Core setters (fluent)
    SetEntityType(entityType string) EntityInterface
    SetEntityHandle(handle string) EntityInterface
    SetCreatedAt(createdAt string) EntityInterface
    SetUpdatedAt(updatedAt string) EntityInterface
    
    // Dynamic attributes (in-memory only)
    GetAttribute(key string) string
    SetAttribute(key string, value string) EntityInterface
    GetAllAttributes() map[string]string
}
```

### Convenience Methods (via type assertion)

```go
// String getter/setter
GetString(key string, defaultValue string) string
SetString(key string, value string) bool

// Int getter/setter
GetInt(key string, defaultValue int64) (int64, error)
SetInt(key string, value int64) bool

// Float getter/setter
GetFloat(key string, defaultValue float64) (float64, error)
SetFloat(key string, value float64) bool

// Interface getter/setter (JSON)
GetInterface(key string, defaultValue interface{}) interface{}
SetInterface(key string, value interface{}) bool
```

## AttributeInterface

Key-value attribute interface.

```go
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
    
    // Core setters (fluent)
    SetEntityID(entityID string) AttributeInterface
    SetAttributeKey(key string) AttributeInterface
    SetAttributeValue(value string) AttributeInterface
    SetCreatedAt(createdAt string) AttributeInterface
    SetUpdatedAt(updatedAt string) AttributeInterface
    
    // Type conversion
    GetInt() (int64, error)
    GetFloat() (float64, error)
    GetInterface() interface{}
    SetInt(value int64) AttributeInterface
    SetFloat(value float64) AttributeInterface
    SetInterface(value interface{}) AttributeInterface
}
```

## RelationshipInterface

Entity relationship interface.

```go
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
    
    // Core setters (fluent)
    SetEntityID(entityID string) RelationshipInterface
    SetRelatedEntityID(relatedID string) RelationshipInterface
    SetRelationshipType(relType string) RelationshipInterface
    SetParentID(parentID string) RelationshipInterface
    SetSequence(sequence int) RelationshipInterface
    SetMetadata(metadata string) RelationshipInterface
    SetCreatedAt(createdAt string) RelationshipInterface
}
```

## TaxonomyInterface

Taxonomy classification system interface.

```go
type TaxonomyInterface interface {
    dataobject.DataObjectInterface
    
    // Core getters
    Name() string
    Slug() string
    Description() string
    ParentID() string
    EntityTypes() []string
    CreatedAt() string
    CreatedAtCarbon() *carbon.Carbon
    UpdatedAt() string
    UpdatedAtCarbon() *carbon.Carbon
    
    // Core setters (fluent)
    SetName(name string) TaxonomyInterface
    SetSlug(slug string) TaxonomyInterface
    SetDescription(desc string) TaxonomyInterface
    SetParentID(parentID string) TaxonomyInterface
    SetEntityTypes(types []string) TaxonomyInterface
    SetCreatedAt(createdAt string) TaxonomyInterface
    SetUpdatedAt(updatedAt string) TaxonomyInterface
}
```

## TaxonomyTermInterface

Taxonomy term (category) interface.

```go
type TaxonomyTermInterface interface {
    dataobject.DataObjectInterface
    
    // Core getters
    TaxonomyID() string
    Name() string
    Slug() string
    ParentID() string
    SortOrder() int
    CreatedAt() string
    CreatedAtCarbon() *carbon.Carbon
    UpdatedAt() string
    UpdatedAtCarbon() *carbon.Carbon
    
    // Core setters (fluent)
    SetTaxonomyID(taxonomyID string) TaxonomyTermInterface
    SetName(name string) TaxonomyTermInterface
    SetSlug(slug string) TaxonomyTermInterface
    SetParentID(parentID string) TaxonomyTermInterface
    SetSortOrder(order int) TaxonomyTermInterface
    SetCreatedAt(createdAt string) TaxonomyTermInterface
    SetUpdatedAt(updatedAt string) TaxonomyTermInterface
}
```

## EntityTaxonomyInterface

Entity-taxonomy assignment interface.

```go
type EntityTaxonomyInterface interface {
    dataobject.DataObjectInterface
    
    // Core getters
    EntityID() string
    TaxonomyID() string
    TermID() string
    CreatedAt() string
    CreatedAtCarbon() *carbon.Carbon
    
    // Core setters (fluent)
    SetEntityID(entityID string) EntityTaxonomyInterface
    SetTaxonomyID(taxonomyID string) EntityTaxonomyInterface
    SetTermID(termID string) EntityTaxonomyInterface
    SetCreatedAt(createdAt string) EntityTaxonomyInterface
}
```

## Query Options

### EntityQueryOptions

```go
type EntityQueryOptions struct {
    ID           string
    IDs          []string
    EntityType   string
    EntityHandle string
    Limit        uint64
    Offset       uint64
    Search       string
    SortBy       string
    SortOrder    string // asc / desc
    CountOnly    bool
}
```

### AttributeQueryOptions

```go
type AttributeQueryOptions struct {
    ID        string
    EntityID  string
    EntityIDs []string
    Key       string
    Keys      []string
    Limit     uint64
    Offset    uint64
    OrderBy   string
    SortOrder string
    CountOnly bool
}
```

### RelationshipQueryOptions

```go
type RelationshipQueryOptions struct {
    ID               string
    IDs              []string
    EntityID         string
    RelatedEntityID  string
    RelationshipType string
    ParentID         string
    Limit            uint64
    Offset           uint64
    OrderBy          string
    SortOrder        string
    CountOnly        bool
}
```

type RelationshipOptions struct {
    EntityID         string
    RelatedEntityID  string
    RelationshipType string
    ParentID         string
    Sequence         int
    Metadata         string
}
```

### TaxonomyQueryOptions

```go
type TaxonomyQueryOptions struct {
    ID         string
    Slug       string
    EntityType string
    Limit      uint64
    Offset     uint64
    OrderBy    string
    SortOrder  string
    CountOnly  bool
}
```

type TaxonomyOptions struct {
    Name        string
    Slug        string
    Description string
    ParentID    string
    EntityTypes []string
}
```

### TaxonomyTermQueryOptions

```go
type TaxonomyTermQueryOptions struct {
    ID         string
    TaxonomyID string
    ParentID   string
    Slug       string
    Limit      uint64
    Offset     uint64
    OrderBy    string
    SortOrder  string
    CountOnly  bool
}
```

type TaxonomyTermOptions struct {
    TaxonomyID string
    Name       string
    Slug       string
    ParentID   string
    SortOrder  int
}
```

### EntityTaxonomyQueryOptions

```go
type EntityTaxonomyQueryOptions struct {
    ID         string
    EntityID   string
    TaxonomyID string
    TermID     string
    Limit      uint64
    Offset     uint64
    OrderBy    string
    SortOrder  string
    CountOnly  bool
}
```

## Constants

### Relationship Types

```go
const (
    RELATIONSHIP_TYPE_BELONGS_TO = "belongs_to"   // Entity belongs to one parent
    RELATIONSHIP_TYPE_HAS_MANY   = "has_many"     // Entity has many children
    RELATIONSHIP_TYPE_MANY_MANY  = "many_to_many" // Bidirectional link
)
```

### Column Names

```go
const (
    COLUMN_ID                = "id"
    COLUMN_ENTITY_TYPE       = "entity_type"
    COLUMN_ENTITY_HANDLE     = "entity_handle"
    COLUMN_ENTITY_ID         = "entity_id"
    COLUMN_ATTRIBUTE_KEY     = "attribute_key"
    COLUMN_ATTRIBUTE_VALUE   = "attribute_value"
    COLUMN_CREATED_AT        = "created_at"
    COLUMN_UPDATED_AT        = "updated_at"
    COLUMN_DELETED_AT        = "deleted_at"
    COLUMN_DELETED_BY        = "deleted_by"
    COLUMN_RELATED_ENTITY_ID = "related_entity_id"
    COLUMN_RELATIONSHIP_TYPE = "relationship_type"
    COLUMN_PARENT_ID         = "parent_id"
    COLUMN_SEQUENCE          = "sequence"
    COLUMN_METADATA          = "metadata"
    COLUMN_NAME              = "name"
    COLUMN_SLUG              = "slug"
    COLUMN_DESCRIPTION       = "description"
    COLUMN_ENTITY_TYPES      = "entity_types"
    COLUMN_TAXONOMY_ID       = "taxonomy_id"
    COLUMN_TERM_ID           = "term_id"
    COLUMN_SORT_ORDER        = "sort_order"
)
```

### Default Table Names

```go
const (
    DEFAULT_RELATIONSHIP_TABLE_NAME       = "entities_relationships"
    DEFAULT_RELATIONSHIP_TRASH_TABLE_NAME = "entities_relationships_trash"
    DEFAULT_TAXONOMY_TABLE_NAME           = "entities_taxonomies"
    DEFAULT_TAXONOMY_TERM_TABLE_NAME    = "entities_taxonomy_terms"
    DEFAULT_ENTITY_TAXONOMY_TABLE_NAME    = "entities_entity_taxonomies"
    DEFAULT_TAXONOMY_TRASH_TABLE_NAME     = "entities_taxonomies_trash"
    DEFAULT_TAXONOMY_TERM_TRASH_TABLE_NAME = "entities_taxonomy_terms_trash"
)
```

## NewStoreOptions

```go
type NewStoreOptions struct {
    // Required
    DB                 *sql.DB
    EntityTableName    string
    AttributeTableName string
    
    // Optional - trash tables
    EntityTrashTableName    string // Defaults to EntityTableName + "_trash"
    AttributeTrashTableName string // Defaults to AttributeTableName + "_trash"
    
    // Optional - relationships
    RelationshipsEnabled       bool
    RelationshipTableName    string // Defaults to "entities_relationships"
    RelationshipTrashTableName string // Defaults to "entities_relationships_trash"
    
    // Optional - taxonomies
    TaxonomiesEnabled          bool
    TaxonomyTableName          string // Defaults to "entities_taxonomies"
    TaxonomyTrashTableName     string // Defaults to "entities_taxonomies_trash"
    TaxonomyTermTableName      string // Defaults to "entities_taxonomy_terms"
    TaxonomyTermTrashTableName string // Defaults to "entities_taxonomy_terms_trash"
    EntityTaxonomyTableName    string // Defaults to "entities_entity_taxonomies"
    
    // Database
    Database     sb.DatabaseInterface // Alternative to DB
    DbDriverName string              // Auto-detected if not set
    
    // Features
    AutomigrateEnabled bool // Auto-create tables
    DebugEnabled       bool // Log SQL queries
}
```

## Helper Functions

```go
// ID Generation
GenerateShortID() string

// Constructor Functions
NewEntity() EntityInterface
NewEntityFromExistingData(data map[string]string) EntityInterface
NewAttribute() AttributeInterface
NewAttributeFromExistingData(data map[string]string) AttributeInterface
NewRelationship() RelationshipInterface
NewRelationshipFromExistingData(data map[string]string) RelationshipInterface
NewTaxonomy() TaxonomyInterface
NewTaxonomyFromExistingData(data map[string]string) TaxonomyInterface
NewTaxonomyTerm() TaxonomyTermInterface
NewTaxonomyTermFromExistingData(data map[string]string) TaxonomyTermInterface
NewEntityTaxonomy() EntityTaxonomyInterface
NewEntityTaxonomyFromExistingData(data map[string]string) EntityTaxonomyInterface
```

## Error Handling

Most methods return errors for:
- Database connection failures
- Constraint violations
- Not found (nil results, not errors)
- Invalid parameters

Pattern:
```go
entity, err := store.EntityFindByID(ctx, id)
if err != nil {
    // Handle error (database error)
    return err
}
if entity == nil {
    // Handle not found
    return errors.New("entity not found")
}
```
