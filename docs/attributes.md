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

### Direct on Entity (Recommended)

```go
entity := store.EntityCreateWithType("person")
entity.SetString("name", "John Doe")
entity.SetInt("age", 30)
entity.SetFloat("salary", 75000.50)
entity.SetInterface("tags", []string{"developer", "golang"})
```

### Using Store Methods

```go
// Create attribute object
attr := store.NewAttribute()
attr.SetEntityID(entity.ID())
attr.SetKey("email")
attr.SetString("john@example.com")
store.AttributeCreate(ctx, attr)
```

### Shortcut Methods

```go
// Set directly via store
store.AttributeSetString(ctx, entity.ID(), "status", "active")
store.AttributeSetInt(ctx, entity.ID(), "login_count", 5)
store.AttributeSetFloat(ctx, entity.ID(), "rating", 4.5)
```

## Retrieving Attributes

### Via Entity (Recommended)

```go
entity, _ := store.EntityFindByID(ctx, "86ccrtsgx")

// Get with default values
name := entity.GetString("name", "Unknown")
age, _ := entity.GetInt("age", 0)
salary, _ := entity.GetFloat("salary", 0.0)
tags := entity.GetInterface("tags", []string{}).([]string)
```

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
    email := attr.GetString()
}
```

## Updating Attributes

Setting an attribute with the same key updates the existing value:

```go
entity.SetString("status", "active")   // Creates or updates
entity.SetString("status", "inactive") // Updates existing
```

## Deleting Attributes

### Hard Delete

```go
deleted, err := store.AttributeDelete(ctx, "86ccrtsgx")
```

### Soft Delete (Trash)

```go
trashed, err := store.AttributeTrash(ctx, "86ccrtsgx")
```

### Delete by Entity ID

Delete all attributes for an entity:

```go
deletedCount, err := store.AttributeDeleteByEntityID(ctx, entity.ID())
```

## Attribute Methods

### Getters

| Method | Description |
|--------|-------------|
| `ID() string` | Returns attribute ID |
| `EntityID() string` | Returns parent entity ID |
| `Key() string` | Returns attribute key |
| `GetString() string` | Returns value as string |
| `GetInt() (int64, error)` | Returns value as int |
| `GetFloat() (float64, error)` | Returns value as float |
| `GetInterface() interface{}` | Returns JSON deserialized value |
| `CreatedAtCarbon() *carbon.Carbon` | Get creation timestamp |
| `UpdatedAtCarbon() *carbon.Carbon` | Get update timestamp |

### Setters (Fluent Interface)

| Method | Description |
|--------|-------------|
| `SetEntityID(id string) AttributeInterface` | Set parent entity ID |
| `SetKey(key string) AttributeInterface` | Set attribute key |
| `SetString(value string) bool` | Set string value |
| `SetInt(value int64) bool` | Set int value |
| `SetFloat(value float64) bool` | Set float value |
| `SetInterface(value interface{}) bool` | Set interface{} value (JSON) |

## Store Methods

| Method | Description |
|--------|-------------|
| `AttributeCount(ctx, opts) (int64, error)` | Count attributes |
| `AttributeCreate(ctx, attr) error` | Create attribute |
| `AttributeDelete(ctx, id string) (bool, error)` | Hard delete |
| `AttributeDeleteByEntityID(ctx, entityID string) (int64, error)` | Delete by entity |
| `AttributeFind(ctx, entityID, key string) (AttributeInterface, error)` | Find attribute |
| `AttributeGetFloat(ctx, entityID, key string) (float64, bool, error)` | Get float value |
| `AttributeGetInt(ctx, entityID, key string) (int64, bool, error)` | Get int value |
| `AttributeGetString(ctx, entityID, key string) (string, bool, error)` | Get string value |
| `AttributeList(ctx, opts) ([]AttributeInterface, error)` | List attributes |
| `AttributeRestore(ctx, id string) (bool, error)` | Restore from trash |
| `AttributeSetFloat(ctx, entityID, key string, value float64) error` | Upsert float |
| `AttributeSetInt(ctx, entityID, key string, value int64) error` | Upsert int |
| `AttributeSetInterface(ctx, entityID, key string, value interface{}) error` | Upsert interface{} |
| `AttributeSetString(ctx, entityID, key, value string) error` | Upsert string |
| `AttributeTrash(ctx, id string) (bool, error)` | Soft delete |
| `AttributeTrashList(ctx, opts) ([]AttributeTrashInterface, error)` | List trashed |
| `AttributeUpdate(ctx, attr) error` | Update attribute |

## Query Options

```go
type AttributeQueryOptions struct {
    EntityID    string
    EntityIDs   []string
    Key         string
    Keys        []string
    Search      string
    OrderBy     string
    SortOrder   string
    Limit       int
    Offset      int
}
```

## Trash Attribute Methods

Trashed attributes have the same getters plus:

- `DeletedAtCarbon() *carbon.Carbon` - When attribute was deleted
- `DeletedBy() string` - Who deleted the attribute

## Best Practices

1. **Use entity methods** for simple attribute operations
2. **Use store methods** when working with attributes independently
3. **Use SetInterface for complex types** - arrays, maps, structs
4. **Handle errors** when parsing int/float values
5. **Use consistent key naming** across your application
