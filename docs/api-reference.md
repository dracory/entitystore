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
EntityList(ctx context.Context, query EntityQueryInterface) ([]EntityInterface, error)
EntityListByAttribute(ctx context.Context, entityType string, attributeKey string, attributeValue string) ([]EntityInterface, error)
EntityCount(ctx context.Context, query EntityQueryInterface) (int64, error)
EntityAttributeList(ctx context.Context, entityID string) ([]AttributeInterface, error)

// Update
EntityUpdate(ctx context.Context, entity EntityInterface) error

// Delete
EntityDelete(ctx context.Context, id string) (bool, error)
EntityTrash(ctx context.Context, id string) (bool, error)
EntityRestore(ctx context.Context, id string) error
```

### Attribute Operations

```go
// CRUD
AttributeCreate(ctx context.Context, attr AttributeInterface) error
AttributeFind(ctx context.Context, entityID string, attributeKey string) (AttributeInterface, error)
AttributeFindByHandle(ctx context.Context, entityType string, entityHandle string, attributeKey string) (AttributeInterface, error)
AttributeList(ctx context.Context, query AttributeQueryInterface) ([]AttributeInterface, error)
AttributeCount(ctx context.Context, query AttributeQueryInterface) (int64, error)
AttributeGroupBy(ctx context.Context, query AttributeQueryInterface) (map[string]any, error)
AttributeUpdate(ctx context.Context, attr AttributeInterface) error
AttributeDelete(ctx context.Context, id string) error
AttributeTrash(ctx context.Context, id string, deletedBy string) error
AttributeRestore(ctx context.Context, id string) error
AttributesDeleteByEntityID(ctx context.Context, entityID string) error

// Shortcuts
AttributeCreateWithKeyAndValue(ctx context.Context, entityID string, attributeKey string, attributeValue string) (AttributeInterface, error)
AttributeSetString(ctx context.Context, entityID string, attributeKey string, attributeValue string) error
AttributeSetInt(ctx context.Context, entityID string, attributeKey string, attributeValue int64) error
AttributeSetFloat(ctx context.Context, entityID string, attributeKey string, attributeValue float64) error
AttributeSetInterface(ctx context.Context, entityID string, attributeKey string, attributeValue interface{}) error
AttributesSet(ctx context.Context, entityID string, attributes map[string]string) error

// Getters
AttributeGetString(ctx context.Context, entityID string, attributeKey string) (value string, exists bool, err error)
AttributeGetInt(ctx context.Context, entityID string, attributeKey string) (value int64, exists bool, err error)
AttributeGetFloat(ctx context.Context, entityID string, attributeKey string) (value float64, exists bool, err error)
```

### Relationship Operations (Requires `RelationshipsEnabled: true`)

```go
// CRUD
RelationshipCreate(ctx context.Context, relationship RelationshipInterface) error
RelationshipCreateByOptions(ctx context.Context, options RelationshipOptions) (RelationshipInterface, error)
RelationshipFind(ctx context.Context, relationshipID string) (RelationshipInterface, error)
RelationshipFindByEntities(ctx context.Context, entityID string, relatedEntityID string, relationshipType string) (RelationshipInterface, error)
RelationshipList(ctx context.Context, query RelationshipQueryInterface) ([]RelationshipInterface, error)
RelationshipCount(ctx context.Context, query RelationshipQueryInterface) (int64, error)
RelationshipDelete(ctx context.Context, relationshipID string) (bool, error)
RelationshipDeleteAll(ctx context.Context, entityID string) error

// Trash
RelationshipTrash(ctx context.Context, relationshipID string, deletedBy string) (bool, error)
RelationshipRestore(ctx context.Context, relationshipID string) (bool, error)
RelationshipTrashList(ctx context.Context, query RelationshipQueryInterface) ([]RelationshipTrashInterface, error)

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
TaxonomyList(ctx context.Context, query TaxonomyQueryInterface) ([]TaxonomyInterface, error)
TaxonomyCount(ctx context.Context, query TaxonomyQueryInterface) (int64, error)
TaxonomyUpdate(ctx context.Context, taxonomy TaxonomyInterface) error
TaxonomyDelete(ctx context.Context, taxonomyID string) (bool, error)
TaxonomyTrash(ctx context.Context, taxonomyID string, deletedBy string) (bool, error)
TaxonomyRestore(ctx context.Context, taxonomyID string) (bool, error)
TaxonomyTrashList(ctx context.Context, query TaxonomyQueryInterface) ([]TaxonomyTrashInterface, error)

// TaxonomyTerm CRUD
TaxonomyTermCreate(ctx context.Context, term TaxonomyTermInterface) error
TaxonomyTermCreateByOptions(ctx context.Context, options TaxonomyTermOptions) (TaxonomyTermInterface, error)
TaxonomyTermFind(ctx context.Context, termID string) (TaxonomyTermInterface, error)
TaxonomyTermFindBySlug(ctx context.Context, taxonomyID string, slug string) (TaxonomyTermInterface, error)
TaxonomyTermList(ctx context.Context, query TaxonomyTermQueryInterface) ([]TaxonomyTermInterface, error)
TaxonomyTermCount(ctx context.Context, query TaxonomyTermQueryInterface) (int64, error)
TaxonomyTermUpdate(ctx context.Context, term TaxonomyTermInterface) error
TaxonomyTermDelete(ctx context.Context, termID string) (bool, error)
TaxonomyTermTrash(ctx context.Context, termID string, deletedBy string) (bool, error)
TaxonomyTermRestore(ctx context.Context, termID string) (bool, error)
TaxonomyTermTrashList(ctx context.Context, query TaxonomyTermQueryInterface) ([]TaxonomyTermTrashInterface, error)

// Entity Assignment
EntityTaxonomyAssign(ctx context.Context, entityID string, taxonomyID string, termID string) error
EntityTaxonomyRemove(ctx context.Context, entityID string, taxonomyID string, termID string) error
EntityTaxonomyList(ctx context.Context, query EntityTaxonomyQueryInterface) ([]EntityTaxonomyInterface, error)
EntityTaxonomyCount(ctx context.Context, query EntityTaxonomyQueryInterface) (int64, error)
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
    GetID() string
    GetType() string
    GetHandle() string
    GetCreatedAt() string
    GetCreatedAtCarbon() *carbon.Carbon
    GetUpdatedAt() string
    GetUpdatedAtCarbon() *carbon.Carbon
    
    // Core setters (fluent)
    SetType(entityType string) EntityInterface
    SetHandle(handle string) EntityInterface
    SetCreatedAt(createdAt string) EntityInterface
    SetUpdatedAt(updatedAt string) EntityInterface
    
    // Temporary attributes (in-memory only, not persisted)
    GetTempKey(key string) string
    SetTempKey(key string, value string) EntityInterface
    GetTempKeys() map[string]string
}
```

### Embedded DataObject Methods

`EntityInterface` embeds `dataobject.DataObjectInterface`, which provides
in-memory key-value storage (`Get`, `Set`, `ID`, etc.). Note that entity
attributes persisted in the attributes table are accessed through the store
(`AttributeGetString`, `AttributeSetString`, ...) — in-memory values are not
persisted automatically.

## AttributeInterface

Key-value attribute interface.

```go
type AttributeInterface interface {
    dataobject.DataObjectInterface
    
    // Core getters
    GetID() string
    GetEntityID() string
    GetKey() string
    GetValue() string
    GetCreatedAt() string
    GetCreatedAtCarbon() *carbon.Carbon
    GetUpdatedAt() string
    GetUpdatedAtCarbon() *carbon.Carbon
    
    // Core setters (fluent)
    SetEntityID(entityID string) AttributeInterface
    SetKey(key string) AttributeInterface
    SetValue(value string) AttributeInterface
    SetCreatedAt(createdAt string) AttributeInterface
    SetUpdatedAt(updatedAt string) AttributeInterface
    
    // Type conversion
    GetInt() (int64, error)
    GetFloat() (float64, error)
    SetInt(value int64) AttributeInterface
    SetFloat(value float64) AttributeInterface
}
```

## RelationshipInterface

Entity relationship interface.

```go
type RelationshipInterface interface {
    dataobject.DataObjectInterface
    
    // Core getters
    GetID() string
    GetEntityID() string
    GetRelatedEntityID() string
    GetRelationshipType() string
    GetParentID() string
    GetSequence() int
    GetMetadata() string
    GetCreatedAt() string
    GetCreatedAtCarbon() *carbon.Carbon
    
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
    GetID() string
    GetName() string
    GetSlug() string
    GetDescription() string
    GetParentID() string
    GetEntityTypes() []string
    GetCreatedAt() string
    GetCreatedAtCarbon() *carbon.Carbon
    GetUpdatedAt() string
    GetUpdatedAtCarbon() *carbon.Carbon
    
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
    GetID() string
    GetTaxonomyID() string
    GetName() string
    GetSlug() string
    GetParentID() string
    GetSortOrder() int
    GetCreatedAt() string
    GetCreatedAtCarbon() *carbon.Carbon
    GetUpdatedAt() string
    GetUpdatedAtCarbon() *carbon.Carbon
    
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
    GetID() string
    GetEntityID() string
    GetTaxonomyID() string
    GetTermID() string
    GetCreatedAt() string
    GetCreatedAtCarbon() *carbon.Carbon
    
    // Core setters (fluent)
    SetEntityID(entityID string) EntityTaxonomyInterface
    SetTaxonomyID(taxonomyID string) EntityTaxonomyInterface
    SetTermID(termID string) EntityTaxonomyInterface
    SetCreatedAt(createdAt string) EntityTaxonomyInterface
}
```

## Fluent Queries

List, count, and trash-list methods accept fluent query interfaces built by
constructor functions. Queries are validated by the store before execution;
passing `nil` or an invalid query returns an error.

```go
entities, err := store.EntityList(ctx, entitystore.EntityQuery().
    WithEntityType("person").
    WithSortBy("created_at").
    WithSortOrder("desc").
    WithLimit(10))
```

Every query interface exposes three method families plus validation:

- `With<Field>(value)` — sets a filter and returns the query for chaining
- `Get<Field>()` — returns the current value (empty if unset)
- `Has<Field>() bool` — reports whether the field was explicitly set
- `Validate() error` — checks the query; called by store methods automatically

Presence flags distinguish an unset field from an explicitly set invalid value:
an empty query is valid, but `WithID("")`, `WithIDs([]string{})`, an unknown
`SortOrder` (only `asc`/`desc` are allowed), or an empty time bound fails
validation.

### Time Range Filters

All query types support inclusive UTC timestamp bounds in
`"YYYY-MM-DD HH:MM:SS"` format:

- `WithCreatedAtGte(ts)` / `WithCreatedAtLte(ts)` — `created_at >= ts` / `created_at <= ts`
- `WithUpdatedAtGte(ts)` / `WithUpdatedAtLte(ts)` — `updated_at >= ts` / `updated_at <= ts`

Both bounds are inclusive and can be combined for a closed range.

### EntityQuery()

`WithID`, `WithIDs`, `WithEntityType`, `WithEntityHandle`, `WithSearch`,
`WithLimit`, `WithOffset`, `WithSortBy`, `WithSortOrder`, `WithCountOnly`,
`WithCreatedAtGte/Lte`, `WithUpdatedAtGte/Lte`, `WithPrefetchAttributes`

`WithSearch` matches against entity attributes via an escaped `LIKE` pattern.

`WithPrefetchAttributes(keys...)` eagerly loads attributes for all returned
entities in a single batch query. The values are stored as in-memory
attributes on each entity (readable via `GetTempKey(key)`), avoiding an N+1
query pattern. Pass specific keys to load only those attributes, or call it
with no arguments to load all attributes.

### AttributeQuery()

`WithID`, `WithIDs`, `WithEntityID`, `WithEntityType`, `WithEntityHandle`,
`WithAttributeKey`, `WithAttributeKeys`, `WithAttributeKeyLike`,
`WithAttributeKeyStartsWith`, `WithAttributeKeyEndsWith`,
`WithAttributeKeyContains`, `WithLimit`, `WithOffset`,
`WithSortBy`, `WithSortOrder`, `WithCountOnly`,
`WithCreatedAtGte/Lte`, `WithUpdatedAtGte/Lte`

`WithEntityType` and `WithEntityHandle` join the entities table; filters apply
to `List`, `Count`, and `TrashList` uniformly via shared helpers.

The key-pattern filters differ in wildcard handling:

- `WithAttributeKeyLike(pattern)` — raw SQL `LIKE` pattern; `%` and `_` act as
  wildcards
- `WithAttributeKeyStartsWith`/`EndsWith`/`Contains` — literal text; any `%`,
  `_`, or `\` in the value is escaped so it matches verbatim

All `LIKE` clauses use an explicit `ESCAPE '\'` for consistent behavior across
SQLite, MySQL, and PostgreSQL.

### RelationshipQuery()

`WithID`, `WithIDs`, `WithEntityID`, `WithEntityIDs`, `WithRelatedEntityID`,
`WithRelatedEntityIDs`, `WithRelationshipType`, `WithParentID`, `WithLimit`,
`WithOffset`, `WithSortBy`, `WithSortOrder`, `WithCountOnly`,
`WithCreatedAtGte/Lte`, `WithUpdatedAtGte/Lte`

### TaxonomyQuery()

`WithID`, `WithIDs`, `WithSlug`, `WithParentID`, `WithEntityType`,
`WithEntityTypes`, `WithLimit`, `WithOffset`, `WithSortBy`, `WithSortOrder`,
`WithCountOnly`, `WithCreatedAtGte/Lte`, `WithUpdatedAtGte/Lte`

### TaxonomyTermQuery()

`WithID`, `WithIDs`, `WithTaxonomyID`, `WithSlug`, `WithParentID`,
`WithLimit`, `WithOffset`, `WithSortBy`, `WithSortOrder`, `WithCountOnly`,
`WithCreatedAtGte/Lte`, `WithUpdatedAtGte/Lte`

### EntityTaxonomyQuery()

`WithID`, `WithEntityID`, `WithEntityIDs`, `WithTaxonomyID`, `WithTermID`,
`WithTermIDs`, `WithLimit`, `WithOffset`, `WithSortBy`, `WithSortOrder`,
`WithCountOnly`, `WithCreatedAtGte/Lte`, `WithUpdatedAtGte/Lte`

## Create Options

`RelationshipCreateByOptions`, `TaxonomyCreateByOptions`, and
`TaxonomyTermCreateByOptions` take plain options structs (distinct from the
fluent queries above):

```go
type RelationshipOptions struct {
    EntityID         string
    RelatedEntityID  string
    RelationshipType string
    ParentID         string
    Sequence         int
    Metadata         string
}

type TaxonomyOptions struct {
    Name        string
    Slug        string
    Description string
    ParentID    string
    EntityTypes []string
}

type TaxonomyTermOptions struct {
    TaxonomyID string
    Name       string
    Slug       string
    ParentID   string
    SortOrder  int
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
