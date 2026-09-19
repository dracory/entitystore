# Entities

Entities are the primary objects in entitystore. Each entity represents a single object (e.g., a person, a book, an order) with a type and any number of attributes.

## Overview

- Entities have a **type** (e.g., "person", "book", "order")
- Entities have a unique **ID** (short ID, 9-15 characters)
- Entities can have **attributes** (key-value pairs)
- Entities support **soft delete** via trash bin

## Creating Entities

### Create with Type Only

```go
person, err := store.EntityCreateWithType(ctx, "person")
store.AttributeSetString(ctx, person.ID(), "name", "John Doe")
store.AttributeSetInt(ctx, person.ID(), "age", 30)
```

### Create with Type and Attributes

```go
entity, err := store.EntityCreateWithTypeAndAttributes(ctx, "person", map[string]string{
    "name": "Jane Doe",
    "age":  "25",
})
```

### Create from Entity Object

```go
entity := entitystore.NewEntity()
entity.SetType("person")
entity.SetHandle("john-doe")
store.EntityCreate(ctx, entity)
```

## Retrieving Entities

### Find by ID

```go
entity, err := store.EntityFindByID(ctx, "86ccrtsgx")
if err != nil {
    log.Fatal(err)
}
name, exists, err := store.AttributeGetString(ctx, entity.ID(), "name")
```

### Find by Attribute

```go
// Find first entity matching attribute
entity, err := store.EntityFindByAttribute(ctx, "person", "email", "john@example.com")
```

### List Entities

```go
// List all entities of type "person"
entities, err := store.EntityList(ctx, entitystore.EntityQuery().
    WithEntityType("person").
    WithLimit(10))

// Search entities (searches across attributes)
results, err := store.EntityList(ctx, entitystore.EntityQuery().
    WithEntityType("person").
    WithSearch("john").
    WithSortBy("created_at").
    WithSortOrder("desc"))

// Entities created within a UTC datetime range (inclusive bounds)
recent, err := store.EntityList(ctx, entitystore.EntityQuery().
    WithCreatedAtGte("2026-01-01 00:00:00").
    WithCreatedAtLte("2026-02-01 00:00:00"))
```

## Entity Methods

### Getters

| Method | Description |
|--------|-------------|
| `GetID() string` | Returns entity unique ID |
| `GetType() string` | Returns entity type |
| `GetHandle() string` | Returns entity handle (slug) |
| `Get(key string) string` | Get in-memory attribute (via embedded DataObject) |
| `GetTempKey(key string) string` | Get temporary in-memory attribute |
| `GetTempKeys() map[string]string` | Get all temporary attributes |
| `GetCreatedAtCarbon() *carbon.Carbon` | Get creation timestamp |
| `GetUpdatedAtCarbon() *carbon.Carbon` | Get update timestamp |

### Setters (Fluent Interface)

| Method | Description |
|--------|-------------|
| `SetType(t string) EntityInterface` | Set entity type |
| `SetHandle(h string) EntityInterface` | Set entity handle |
| `SetTempKey(key, value string) EntityInterface` | Set temporary in-memory attribute |

Persisted attributes live in the attributes table — use
`AttributeSetString`/`AttributeSetInt`/`AttributeSetFloat` and
`AttributeGetString`/`AttributeGetInt`/`AttributeGetFloat` to read and write
them by entity ID.

## Updating Entities

```go
entity, _ := store.EntityFindByID(ctx, "86ccrtsgx")
entity.SetHandle("new-handle")
store.EntityUpdate(ctx, entity)

// Update an attribute
store.AttributeSetString(ctx, entity.ID(), "name", "Updated Name")
```

## Deleting Entities

### Hard Delete

Permanently removes entity and all its attributes:

```go
deleted, err := store.EntityDelete(ctx, "86ccrtsgx")
```

### Soft Delete (Trash)

Moves entity to trash bin (requires trash to be enabled):

```go
trashed, err := store.EntityTrash(ctx, "86ccrtsgx")
```

### Restore from Trash

```go
err := store.EntityRestore(ctx, "86ccrtsgx")
```

## Store Methods

| Method | Description |
|--------|-------------|
| `EntityCount(ctx, query) (int64, error)` | Count entities matching query |
| `EntityCreate(ctx, entity) error` | Create from entity object |
| `EntityCreateWithType(ctx, type string) (EntityInterface, error)` | Create with type only |
| `EntityCreateWithTypeAndAttributes(ctx, type string, attrs map) (EntityInterface, error)` | Create with attributes |
| `EntityDelete(ctx, id string) (bool, error)` | Hard delete entity |
| `EntityFindByID(ctx, id string) (EntityInterface, error)` | Find by ID |
| `EntityFindByAttribute(ctx, type, key, value string) (EntityInterface, error)` | Find by attribute |
| `EntityList(ctx, query) ([]EntityInterface, error)` | List entities |
| `EntityListByAttribute(ctx, type, key, value string) ([]EntityInterface, error)` | List by attribute |
| `EntityRestore(ctx, id string) error` | Restore from trash |
| `EntityTrash(ctx, id string) (bool, error)` | Soft delete entity |
| `EntityUpdate(ctx, entity) error` | Update entity |

## Fluent Queries

`EntityList` and `EntityCount` accept an `EntityQueryInterface` built with the
fluent `EntityQuery()` constructor. `With*` setters chain, `Get*`/`Has*`
accessors expose the current state, and `Validate()` is run by the store
(nil or invalid queries are rejected).

```go
entities, err := store.EntityList(ctx, entitystore.EntityQuery().
    WithEntityType("person").
    WithIDs([]string{"id1", "id2"}).
    WithLimit(10).
    WithOffset(0).
    WithSortBy("created_at").
    WithSortOrder("desc")) // asc / desc
```

Available filters: `WithID`, `WithIDs`, `WithEntityType`, `WithEntityHandle`,
`WithSearch`, `WithLimit`, `WithOffset`, `WithSortBy`, `WithSortOrder`,
`WithCountOnly`, and inclusive UTC time bounds `WithCreatedAtGte/Lte` and
`WithUpdatedAtGte/Lte` in `"YYYY-MM-DD HH:MM:SS"` format.

## Trash Entity Methods

Trashed entities have the same getters as regular entities plus:

- `GetDeletedAtCarbon() *carbon.Carbon` - When entity was deleted
- `DeletedBy() string` - Who deleted the entity
