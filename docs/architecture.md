# Architecture

Understanding Entity Store's design patterns and database structure.

## Overview

Entity Store implements the **EAV (Entity-Attribute-Value)** pattern, a flexible approach to storing schemaless data in relational databases.

## Design Patterns

### EAV Pattern

Traditional relational databases require predefined schemas. EAV decouples entity structure from storage:

```
Traditional Table:
| id | name  | email           | age |
|----|-------|-----------------|-----|
| 1  | John  | john@email.com  | 30  |

EAV Tables:
entities:           attributes:
| id | type   |     | entity_id | key   | value           |
|----|--------|     |-----------|-------|-----------------|
| 1  | person |     | 1         | name  | John            |
                    | 1         | email | john@email.com  |
                    | 1         | age   | 30              |
```

**Benefits:**
- No schema migrations for new attributes
- Dynamic attribute addition/removal
- Consistent storage format
- Full SQL query capability

### DataObject Pattern

All domain objects embed `dataobject.DataObject`:

```go
type entityImplementation struct {
    dataobject.DataObject  // Embedded map[string]string storage
}
```

**Features:**
- Map-based internal storage
- Fluent setter interface (chainable)
- Automatic dirty tracking
- Hydration from database rows

### Repository Pattern

Store methods act as repositories:

```go
type StoreInterface interface {
    // Entity repository
    EntityCreate(ctx context.Context, entity EntityInterface) error
    EntityFindByID(ctx context.Context, id string) (EntityInterface, error)
    EntityList(ctx context.Context, options EntityQueryOptions) ([]EntityInterface, error)
    
    // Attribute repository  
    AttributeSetString(ctx context.Context, entityID, key, value string) error
    AttributeFind(ctx context.Context, entityID, key string) (AttributeInterface, error)
    
    // ...
}
```

## Database Schema

### Core Tables

<img src="images/entitystore-database-schema.svg" width="800" />

#### entities

Primary entity storage.

| Column | Type | Description |
|--------|------|-------------|
| id | string(9) | Primary key (short ID) |
| entity_type | string(40) | Entity classification |
| entity_handle | string(60) | Unique handle/slug |
| created_at | datetime | Creation timestamp |
| updated_at | datetime | Last update timestamp |

#### attributes

Typed key-value storage.

| Column | Type | Description |
|--------|------|-------------|
| id | string(9) | Primary key |
| entity_id | string(9) | FK to entities |
| attribute_key | string(255) | Attribute name |
| attribute_value | text | Value (string/JSON) |
| created_at | datetime | Creation timestamp |
| updated_at | datetime | Last update timestamp |

### Trash Tables

Soft delete support via parallel `_trash` tables:

- `entities_trash` - Soft-deleted entities
- `attributes_trash` - Soft-deleted attributes
- `relationships_trash` - Soft-deleted relationships
- `taxonomies_trash` - Soft-deleted taxonomies
- `taxonomy_terms_trash` - Soft-deleted terms

Additional columns:
- `deleted_at` - When entity was deleted
- `deleted_by` - Who deleted the entity

### Optional Tables

#### relationships

Links entities together.

| Column | Type | Description |
|--------|------|-------------|
| id | string(9) | Primary key |
| entity_id | string(9) | Source entity FK |
| related_entity_id | string(9) | Target entity FK |
| relationship_type | string(50) | belongs_to/has_many/many_to_many |
| parent_id | string(9) | Hierarchical parent |
| sequence | int | Display order |
| metadata | text | JSON metadata |
| created_at | datetime | Creation timestamp |

#### taxonomies, taxonomy_terms, entity_taxonomies

See [Taxonomies](taxonomies.md) documentation.

## Type System

### Attribute Types

Values stored as strings with type conversion:

| Type | Storage | Conversion |
|------|---------|------------|
| String | Plain text | Direct |
| Integer | String | `strconv.ParseInt` |
| Float | String | `strconv.ParseFloat` |
| Interface | JSON | `json.Marshal/Unmarshal` |

### ID Generation

Short IDs (9 characters by default) using `GenerateShortID()`:

- URL-safe characters
- Collision-resistant
- Human-readable
- Database index friendly

## Query Builder

Entity Store uses `goqu` for SQL generation:

```go
func (st *storeImplementation) EntityQuery(options EntityQueryOptions) *goqu.SelectDataset {
    q := goqu.Dialect(st.dbDriverName).From(st.entityTableName)
    
    if options.EntityType != "" {
        q = q.Where(goqu.C(COLUMN_ENTITY_TYPE).Eq(options.EntityType))
    }
    
    if options.SortOrder == "asc" {
        q = q.Order(goqu.I(options.SortBy).Asc())
    }
    
    return q.Offset(uint(options.Offset)).Limit(uint(options.Limit))
}
```

## Transaction Support

Operations use context for cancellation and can participate in transactions:

```go
// Automatic transaction handling
store.EntityCreate(ctx, entity)

// Manual transaction
tx, _ := db.BeginTx(ctx, nil)
// ... use transaction context
```

## Interfaces

Core interfaces define contracts:

```go
// EntityInterface - Primary domain object
type EntityInterface interface {
    dataobject.DataObjectInterface
    EntityType() string
    SetEntityType(string) EntityInterface
    // ...
}

// StoreInterface - Repository operations
type StoreInterface interface {
    AutoMigrate(ctx context.Context) error
    EntityCreate(ctx context.Context, entity EntityInterface) error
    // ...
}
```

## Extension Points

### Custom Entity Types

Create domain-specific wrappers:

```go
type Product struct {
    entitystore.EntityInterface
}

func (p *Product) Name() string {
    return p.GetString("name", "")
}

func (p *Product) SetName(name string) {
    p.SetString("name", name)
}
```

### Query Extensions

Extend query options for custom filtering:

```go
type ProductQueryOptions struct {
    entitystore.EntityQueryOptions
    MinPrice float64
    MaxPrice float64
}
```

## Performance Considerations

### Indexing

Core indexes provided:
- `entities.id` - Primary key
- `entities.entity_type` - Type filtering
- `entities.entity_handle` - Handle lookups
- `attributes.entity_id` - Entity attribute queries
- `attributes.attribute_key` - Key-based queries

### Query Optimization

- Use `CountOnly: true` for existence checks
- Limit result sets with `Limit`
- Use `IDs` filter for batch operations
- Avoid full table scans with proper filtering

### Caching Strategy

Entity Store does not implement caching. Recommended approach:

```go
type CachedStore struct {
    store entitystore.StoreInterface
    cache cache.Interface
}

func (cs *CachedStore) EntityFindByID(ctx context.Context, id string) (entitystore.EntityInterface, error) {
    // Check cache first
    if cached := cs.cache.Get(id); cached != nil {
        return cached.(entitystore.EntityInterface), nil
    }
    
    // Load from store
    entity, err := cs.store.EntityFindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Populate cache
    cs.cache.Set(id, entity, ttl)
    return entity, nil
}
```

## Data Integrity

### Referential Integrity

Entity Store maintains integrity through:
- Foreign key constraints (database level)
- Soft delete cascade prevention
- Assignment validation

### Validation

Built-in validations:
- Required fields (ID, type)
- Slug uniqueness (taxonomies, terms)
- Entity existence before assignment

## Security Considerations

### SQL Injection

All queries use parameterized statements via `goqu`:

```go
// Safe - uses parameterization
q.Where(goqu.C("name").Eq(userInput))

// Unsafe - never do this
q.Where(goqu.L("name = '" + userInput + "'"))
```

### Access Control

Entity Store does not implement authorization. Recommended layers:

1. **Application-level** - Check permissions before store calls
2. **Service-level** - Wrap store with permission checks
3. **Database-level** - Row-level security (PostgreSQL)

## Testing Strategy

### Unit Testing

Mock the store interface:

```go
type mockStore struct {
    entities map[string]entitystore.EntityInterface
}

func (m *mockStore) EntityFindByID(ctx context.Context, id string) (entitystore.EntityInterface, error) {
    return m.entities[id], nil
}
```

### Integration Testing

Use test utilities:

```go
func TestEntityCreate(t *testing.T) {
    store := testutils.InitStore(t)
    ctx := context.Background()
    
    entity := store.EntityCreateWithType("test")
    err := store.EntityCreate(ctx, entity)
    
    assert.NoError(t, err)
    assert.NotEmpty(t, entity.ID())
}
```

## Migration Strategy

### Schema Evolution

Entity Store handles schema through `AutoMigrate()`:

1. **New tables** - Created automatically
2. **New columns** - Added on migration
3. **Existing data** - Preserved

### Version Compatibility

- Semantic versioning for API changes
- Deprecation warnings before breaking changes
- Migration guides for major versions

## Debugging

Enable debug mode for SQL logging:

```go
store, _ := entitystore.NewStore(entitystore.NewStoreOptions{
    // ...
    DebugEnabled: true, // Logs all SQL queries
})
```

## Related Documentation

- [Getting Started](getting-started.md) - Installation and first steps
- [Entities](entities.md) - Working with entities
- [Attributes](attributes.md) - Typed attribute storage
- [Relationships](entity-relationships.md) - Entity linking
- [Taxonomies](taxonomies.md) - Categorization system
