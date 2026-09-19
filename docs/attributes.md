# Attributes

Attributes store typed data for entities. Each attribute is a key-value pair linked to a specific entity.

## Overview

- Attributes belong to a specific entity
- Attribute keys are strings
- Values can be: string, int, float, or interface{} (serialized to JSON)
- Attributes support **soft delete** via trash bin

## Attribute Types

| Type | Go Type | Storage |
|------|---------|---------|
| String | `string` | Plain text |
| Integer | `int64` | String representation |
| Float | `float64` | String representation |
| Interface | `interface{}` | JSON serialized |

## Creating Attributes

### Using Store Setters (Recommended)

```go
entity, _ := store.EntityCreateWithType(ctx, "person")
store.AttributeSetString(ctx, entity.ID(), "name", "John Doe")
store.AttributeSetInt(ctx, entity.ID(), "age", 30)
store.AttributeSetFloat(ctx, entity.ID(), "salary", 75000.50)
```

`AttributesSet` sets several attributes in one call:

```go
err := store.AttributesSet(ctx, entity.ID(), map[string]string{
    "name": "John Doe",
    "age":  "30",
})
```

### Using Attribute Objects

```go
attr := entitystore.NewAttribute()
attr.SetEntityID(entity.ID()).SetKey("email").SetValue("john@example.com")
err := store.AttributeCreate(ctx, attr)
```

## Retrieving Attributes

### Using Store Getters

```go
// Get attribute values directly (returns value, exists flag, and error)
name, exists, err := store.AttributeGetString(ctx, entity.ID(), "name")
if err != nil {
    // Handle database error
}
if !exists {
    // Handle missing attribute
}

// Get typed values
age, exists, err := store.AttributeGetInt(ctx, entity.ID(), "age")
rating, exists, err := store.AttributeGetFloat(ctx, entity.ID(), "rating")
```

### Direct Attribute Access

```go
// Get attribute object
attr, _ := store.AttributeFind(ctx, entity.ID(), "email")
if attr != nil {
    email := attr.GetValue()
    asInt, err := attr.GetInt()
    asFloat, err := attr.GetFloat()
}
```

## Updating Attributes

Setting an attribute with the same key updates the existing value:

```go
store.AttributeSetString(ctx, entity.ID(), "status", "active")   // Creates or updates
store.AttributeSetString(ctx, entity.ID(), "status", "inactive") // Updates existing
```

## Deleting Attributes

### Hard Delete

```go
err := store.AttributeDelete(ctx, "86ccrtsgx")
```

### Soft Delete (Trash)

```go
err := store.AttributeTrash(ctx, "86ccrtsgx", "admin")
```

### Delete by Entity ID

Delete all attributes for an entity:

```go
err := store.AttributesDeleteByEntityID(ctx, entity.ID())
```

## Attribute Methods

### Getters

| Method | Description |
|--------|-------------|
| `GetID() string` | Returns attribute ID |
| `GetEntityID() string` | Returns parent entity ID |
| `GetKey() string` | Returns attribute key |
| `GetValue() string` | Returns value as string |
| `GetInt() (int64, error)` | Returns value as int |
| `GetFloat() (float64, error)` | Returns value as float |
| `GetCreatedAtCarbon() *carbon.Carbon` | Get creation timestamp |
| `GetUpdatedAtCarbon() *carbon.Carbon` | Get update timestamp |

### Setters (Fluent Interface)

| Method | Description |
|--------|-------------|
| `SetEntityID(id string) AttributeInterface` | Set parent entity ID |
| `SetKey(key string) AttributeInterface` | Set attribute key |
| `SetValue(value string) AttributeInterface` | Set string value |
| `SetInt(value int64) AttributeInterface` | Set int value |
| `SetFloat(value float64) AttributeInterface` | Set float value |

## Store Methods

| Method | Description |
|--------|-------------|
| `AttributeCount(ctx, query) (int64, error)` | Count attributes |
| `AttributeCreate(ctx, attr) error` | Create attribute |
| `AttributeDelete(ctx, id string) error` | Hard delete |
| `AttributesDeleteByEntityID(ctx, entityID string) error` | Delete by entity |
| `AttributeFind(ctx, entityID, key string) (AttributeInterface, error)` | Find attribute |
| `AttributeGetFloat(ctx, entityID, key string) (float64, bool, error)` | Get float value |
| `AttributeGetInt(ctx, entityID, key string) (int64, bool, error)` | Get int value |
| `AttributeGetString(ctx, entityID, key string) (string, bool, error)` | Get string value |
| `AttributeList(ctx, query) ([]AttributeInterface, error)` | List attributes |
| `AttributeRestore(ctx, id string) error` | Restore from trash |
| `AttributeSetFloat(ctx, entityID, key string, value float64) error` | Upsert float |
| `AttributeSetInt(ctx, entityID, key string, value int64) error` | Upsert int |
| `AttributeSetInterface(ctx, entityID, key string, value interface{}) error` | Upsert interface{} |
| `AttributeSetString(ctx, entityID, key, value string) error` | Upsert string |
| `AttributeTrash(ctx, id string, deletedBy string) error` | Soft delete |
| `AttributeUpdate(ctx, attr) error` | Update attribute |

## Fluent Queries

`AttributeList` and `AttributeCount` accept an `AttributeQueryInterface` built
with the fluent `AttributeQuery()` constructor. `With*` setters chain,
`Get*`/`Has*` accessors expose state, and the store calls `Validate()` before
executing (nil or invalid queries are rejected).

```go
attrs, err := store.AttributeList(ctx, entitystore.AttributeQuery().
    WithEntityID(entity.ID()).
    WithAttributeKeys([]string{"name", "age"}).
    WithSortBy("created_at").
    WithSortOrder("desc").
    WithLimit(10))
```

Available filters: `WithID`, `WithIDs`, `WithEntityID`, `WithEntityType`,
`WithEntityHandle`, `WithAttributeKey`, `WithAttributeKeys`, `WithLimit`,
`WithOffset`, `WithSortBy`, `WithSortOrder`, `WithCountOnly`, and inclusive
UTC time bounds `WithCreatedAtGte/Lte` and `WithUpdatedAtGte/Lte` in
`"YYYY-MM-DD HH:MM:SS"` format.

`WithEntityType`/`WithEntityHandle` join the entities table. The same filters
are applied uniformly to `AttributeList`, `AttributeCount`, and trash listing
via shared helpers, so count and list results never drift apart.

## Trash Attribute Methods

Trashed attributes have the same getters plus:

- `GetDeletedAtCarbon() *carbon.Carbon` - When attribute was deleted
- `GetDeletedBy() string` - Who deleted the attribute

## Best Practices

1. **Use store setters** (`AttributeSetString`/`SetInt`/`SetFloat`) for simple attribute operations
2. **Use attribute objects** when working with attributes independently
3. **Handle errors** when parsing int/float values
4. **Use consistent key naming** across your application
