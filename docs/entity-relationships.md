# Entity Relationships

Entity relationships allow you to link entities together using different relationship types (belongs_to, has_many, many_to_many). This feature is optional and disabled by default.

## Database Schema

<img src="images/entity-relationships-schema.svg" />

## Setup

```go
store, err := entitystore.NewStore(entitystore.NewStoreOptions{
    DB:                      db,
    EntityTableName:         "entities_entity",
    AttributeTableName:      "entities_attribute",
    RelationshipsEnabled:    true,                         // Enable relationship support
    RelationshipTableName:     "entities_relationships",     // Optional: custom table name
    AutomigrateEnabled:      true,
})
```

## Relationship Types

| Type | Description |
|------|-------------|
| `RELATIONSHIP_TYPE_BELONGS_TO` | Entity belongs to one parent |
| `RELATIONSHIP_TYPE_HAS_MANY` | Entity has many children |
| `RELATIONSHIP_TYPE_MANY_MANY` | Entities linked bidirectionally |

## Basic Usage

### Create a Relationship

```go
// Create entities
author, _ := store.EntityCreateWithType(ctx, "author")
store.AttributeSetString(ctx, author.ID(), "name", "John Doe")

book, _ := store.EntityCreateWithType(ctx, "book")
store.AttributeSetString(ctx, book.ID(), "title", "Go Programming")

// Link book to author
rel, _ := store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
    EntityID:         book.ID(),
    RelatedEntityID:  author.ID(),
    RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
})
```

### Query Relationships

```go
// Find all books by author
relationships, _ := store.RelationshipListRelated(ctx, author.ID(), entitystore.RELATIONSHIP_TYPE_BELONGS_TO)
for _, rel := range relationships {
    book, _ := store.EntityFindByID(ctx, rel.GetEntityID())
    title, _, _ := store.AttributeGetString(ctx, book.ID(), "title")
    fmt.Println(title)
}

// Or with a fluent query (validated by the store)
relationships, _ = store.RelationshipList(ctx, entitystore.RelationshipQuery().
    WithRelatedEntityID(author.ID()).
    WithRelationshipType(entitystore.RELATIONSHIP_TYPE_BELONGS_TO))
```

## Hierarchical Relationships

Use `parent_id` and `sequence` for tree structures:

```go
// Create nested categories
electronics, _ := store.EntityCreateWithType(ctx, "category")
store.AttributeSetString(ctx, electronics.ID(), "name", "Electronics")

phones, _ := store.EntityCreateWithType(ctx, "category")
store.AttributeSetString(ctx, phones.ID(), "name", "Phones")

// Create relationship with parent_id
store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
    EntityID:         phones.ID(),
    RelatedEntityID:  electronics.ID(),
    RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
    ParentID:         relElectronics.ID(),  // Child of electronics
    Sequence:         1,                     // Order within parent
})
```

## Store Methods

### CRUD Operations

- `RelationshipCreate(ctx, relationship RelationshipInterface) error`
- `RelationshipCreateByOptions(ctx, options RelationshipOptions) (RelationshipInterface, error)`
- `RelationshipFind(ctx, relationshipID string) (RelationshipInterface, error)`
- `RelationshipFindByEntities(ctx, entityID, relatedEntityID, relationshipType string) (RelationshipInterface, error)`
- `RelationshipList(ctx, query RelationshipQueryInterface) ([]RelationshipInterface, error)`
- `RelationshipListRelated(ctx, relatedEntityID string, relationshipType string) ([]RelationshipInterface, error)`
- `RelationshipCount(ctx, query RelationshipQueryInterface) (int64, error)`
- `RelationshipDelete(ctx, relationshipID string) (bool, error)`
- `RelationshipDeleteAll(ctx, entityID string) error`

### Trash Operations

- `RelationshipTrash(ctx, relationshipID string, deletedBy string) (bool, error)`
- `RelationshipRestore(ctx, relationshipID string) (bool, error)`
- `RelationshipTrashList(ctx, query RelationshipQueryInterface) ([]RelationshipTrashInterface, error)`

## Relationship Object Methods

- `GetEntityID() string` / `SetEntityID(id string) RelationshipInterface`
- `GetRelatedEntityID() string` / `SetRelatedEntityID(id string) RelationshipInterface`
- `GetRelationshipType() string` / `SetRelationshipType(t string) RelationshipInterface`
- `GetParentID() string` / `SetParentID(id string) RelationshipInterface`
- `GetSequence() int` / `SetSequence(n int) RelationshipInterface`
- `GetMetadata() string` / `SetMetadata(json string) RelationshipInterface`

`RelationshipList`, `RelationshipCount`, and `RelationshipTrashList` accept a
`RelationshipQueryInterface` built via `RelationshipQuery()`, with filters for
ID(s), entity/related-entity IDs, type, parent ID, pagination, sorting, and
inclusive created/updated UTC time bounds.
