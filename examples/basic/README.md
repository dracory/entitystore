# Basic EntityStore Example

This example demonstrates the core functionality of EntityStore for creating, reading, updating, and deleting entities with schemaless attributes.

## What This Example Shows

### 1. Store Initialization
- Creating an EntityStore with SQLite (in-memory database)
- Enabling automatic migration (`AutomigrateEnabled`)

### 2. Entity Management
- Creating entities with types (e.g., "person", "product")
- Setting attributes on entities (strings, integers)
- Finding entities by ID
- Listing all entities with filtering options
- Counting entities (total and by type)

### 3. Attribute Operations
- Storing typed attributes (string, int, float)
- Retrieving attributes via the store
- Converting attribute values to specific types
- Listing all attributes for an entity

### 4. Soft Deletes
- Moving entities to trash (soft delete)
- Entities are not permanently deleted but moved to `entities_trash` table

## Running the Example

```bash
go run examples/basic/main.go
```

## Running Tests

```bash
go test ./examples/basic/... -v
```

## Key Concepts

**EAV Pattern**: Entity-Attribute-Value pattern stores data in three tables:
- `entities` - Stores entity IDs and types
- `attributes` - Stores key-value pairs linked to entities
- `entities_trash` - Stores soft-deleted entities

**In-Memory Attributes**: The `GetAttribute()` and `SetAttribute()` methods on entities only work with in-memory data. To persist attributes, use `store.AttributeSetString()` and `store.AttributeFind()`.

## Code Highlights

```go
// Create entity with attributes
person, _ := store.EntityCreateWithTypeAndAttributes(ctx, "person", map[string]string{
    "name": "John Doe",
    "age":  "30",
})

// Find and retrieve attributes
found, _ := store.EntityFindByID(ctx, person.ID())
attr, _ := store.AttributeFind(ctx, found.ID(), "name")
fmt.Println(attr.AttributeValue()) // "John Doe"

// Soft delete
store.EntityTrash(ctx, person.ID())
```
