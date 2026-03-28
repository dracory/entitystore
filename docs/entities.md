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
person := store.EntityCreateWithType("person")
person.SetString("name", "John Doe")
person.SetInt("age", 30)
```

### Create with Type and Attributes

```go
attributes := map[string]interface{}{
    "name": "Jane Doe",
    "age":  25,
}
entity := store.EntityCreateWithTypeAndAttributes("person", attributes)
```

### Create from Entity Object

```go
entity := store.NewEntity()
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
name := entity.GetString("name", "Unknown")
```

### Find by Attribute

```go
// Find first entity matching attribute
entity, err := store.EntityFindByAttribute(ctx, "person", "email", "john@example.com")
```

### List Entities

```go
// List all entities of type "person"
entities, err := store.EntityList(ctx, entitystore.EntityQueryOptions{
    EntityType: "person",
    Limit:      10,
    Offset:     0,
})

// Search entities
results, err := store.EntityList(ctx, entitystore.EntityQueryOptions{
    EntityType: "person",
    Search:     "john",
    OrderBy:    "created_at",
    SortOrder:  "desc",
})
```

## Entity Methods

### Getters

| Method | Description |
|--------|-------------|
| `ID() string` | Returns entity unique ID |
| `Type() string` | Returns entity type |
| `Handle() string` | Returns entity handle (slug) |
| `GetString(key, defaultValue string) string` | Get attribute as string |
| `GetInt(key string, defaultValue int64) (int64, error)` | Get attribute as int |
| `GetFloat(key string, defaultValue float64) (float64, error)` | Get attribute as float |
| `GetInterface(key string, defaultValue interface{}) interface{}` | Get attribute as interface{} |
| `GetAttribute(key string) AttributeInterface` | Get attribute object |
| `CreatedAtCarbon() *carbon.Carbon` | Get creation timestamp |
| `UpdatedAtCarbon() *carbon.Carbon` | Get update timestamp |

### Setters (Fluent Interface)

| Method | Description |
|--------|-------------|
| `SetType(t string) EntityInterface` | Set entity type |
| `SetHandle(h string) EntityInterface` | Set entity handle |
| `SetString(key, value string) bool` | Set string attribute |
| `SetInt(key string, value int64) bool` | Set int attribute |
| `SetFloat(key string, value float64) bool` | Set float attribute |
| `SetInterface(key string, value interface{}) bool` | Set interface{} attribute |

## Updating Entities

```go
entity, _ := store.EntityFindByID(ctx, "86ccrtsgx")
entity.SetString("name", "Updated Name")
store.EntityUpdate(ctx, entity)
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
restored, err := store.EntityRestore(ctx, "86ccrtsgx")
```

### List Trashed Entities

```go
trashed, err := store.EntityTrashList(ctx, entitystore.EntityQueryOptions{
    Limit: 10,
})
```

## Store Methods

| Method | Description |
|--------|-------------|
| `EntityCount(ctx, opts) (int64, error)` | Count entities matching options |
| `EntityCreate(ctx, entity) error` | Create from entity object |
| `EntityCreateWithType(type string) EntityInterface` | Create with type only |
| `EntityCreateWithTypeAndAttributes(type string, attrs map) EntityInterface` | Create with attributes |
| `EntityDelete(ctx, id string) (bool, error)` | Hard delete entity |
| `EntityFindByID(ctx, id string) (EntityInterface, error)` | Find by ID |
| `EntityFindByAttribute(ctx, type, key, value string) (EntityInterface, error)` | Find by attribute |
| `EntityList(ctx, opts) ([]EntityInterface, error)` | List entities |
| `EntityListByAttribute(ctx, type, key, value string) ([]EntityInterface, error)` | List by attribute |
| `EntityRestore(ctx, id string) (bool, error)` | Restore from trash |
| `EntityTrash(ctx, id string) (bool, error)` | Soft delete entity |
| `EntityTrashList(ctx, opts) ([]EntityTrashInterface, error)` | List trashed entities |
| `EntityUpdate(ctx, entity) error` | Update entity |

## Query Options

```go
type EntityQueryOptions struct {
    EntityType string
    IDs        []string
    Handle     string
    Search     string
    OrderBy    string
    SortOrder  string
    Limit      int
    Offset     int
}
```

## Trash Entity Methods

Trashed entities have the same getters as regular entities plus:

- `DeletedAtCarbon() *carbon.Carbon` - When entity was deleted
- `DeletedBy() string` - Who deleted the entity
