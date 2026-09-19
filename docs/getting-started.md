# Getting Started

Complete guide to installing and using Entity Store in your Go application.

## Table of Contents

- [Installation](#installation)
- [Basic Setup](#basic-setup)
- [Creating Your First Entity](#creating-your-first-entity)
- [Working with Attributes](#working-with-attributes)
- [Querying Entities](#querying-entities)
- [Soft Deletes](#soft-deletes)
- [Next Steps](#next-steps)

## Installation

```bash
go get -u github.com/dracory/entitystore
```

Entity Store requires a SQL database. It works with any database supported by the `goqu` query builder:

- SQLite
- PostgreSQL
- MySQL
- SQL Server

## Basic Setup

### 1. Database Connection

```go
package main

import (
    "database/sql"
    "log"
    
    "github.com/dracory/entitystore"
    _ "modernc.org/sqlite" // Or your preferred driver
)

func main() {
    // Open database connection
    db, err := sql.Open("sqlite", "app.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Create store
    store, err := entitystore.NewStore(entitystore.NewStoreOptions{
        DB:                 db,
        EntityTableName:    "entities",
        AttributeTableName: "attributes",
        AutomigrateEnabled: true, // Auto-create tables
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Use store...
}
```

### 2. Required Options

| Option | Description | Required |
|--------|-------------|----------|
| `DB` | Database connection | Yes |
| `EntityTableName` | Table for entities | Yes |
| `AttributeTableName` | Table for attributes | Yes |
| `AutomigrateEnabled` | Auto-create tables | Recommended |

### 3. Optional Features

Enable additional features:

```go
store, err := entitystore.NewStore(entitystore.NewStoreOptions{
    DB:                   db,
    EntityTableName:      "entities",
    AttributeTableName:   "attributes",
    RelationshipsEnabled: true, // Enable relationships
    TaxonomiesEnabled:    true, // Enable taxonomies
    AutomigrateEnabled:   true,
})
```

## Creating Your First Entity

Entities are the core objects in Entity Store. Each entity has a type and a set of attributes.

### Simple Creation

```go
ctx := context.Background()

// Create and persist a person entity
person, err := store.EntityCreateWithType(ctx, "person")
if err != nil {
    log.Fatal(err)
}

// Set attributes
store.AttributeSetString(ctx, person.ID(), "name", "John Doe")
store.AttributeSetInt(ctx, person.ID(), "age", 30)

fmt.Println("Created person with ID:", person.ID())
```

### With Multiple Attributes

```go
// Create with attributes map
user, err := store.EntityCreateWithTypeAndAttributes(ctx, "user", map[string]string{
    "name":  "Jane Doe",
    "email": "jane@example.com",
    "role":  "admin",
})
```

## Working with Attributes

Attributes store typed data for entities:

### Setting Attributes

```go
entity, _ := store.EntityCreateWithType(ctx, "product")

// String
store.AttributeSetString(ctx, entity.ID(), "name", "Laptop")
store.AttributeSetString(ctx, entity.ID(), "sku", "LAP-001")

// Numbers
store.AttributeSetInt(ctx, entity.ID(), "stock", 50)
store.AttributeSetFloat(ctx, entity.ID(), "price", 999.99)

// Multiple attributes at once
store.AttributesSet(ctx, entity.ID(), map[string]string{
    "color": "silver",
    "brand": "ACME",
})
```

### Getting Attributes

```go
// Retrieve entity
product, _ := store.EntityFindByID(ctx, "abc123xyz")

// Get attributes (returns value, exists flag, and error)
name, exists, err := store.AttributeGetString(ctx, product.ID(), "name")
stock, _, _ := store.AttributeGetInt(ctx, product.ID(), "stock")
price, _, _ := store.AttributeGetFloat(ctx, product.ID(), "price")
```

## Querying Entities

### Find by ID

```go
entity, err := store.EntityFindByID(ctx, "abc123xyz")
if err != nil {
    log.Fatal(err)
}
if entity != nil {
    name, _, _ := store.AttributeGetString(ctx, entity.ID(), "name")
    fmt.Println("Found:", name)
}
```

### List by Type

```go
people, err := store.EntityList(ctx, entitystore.EntityQuery().
    WithEntityType("person").
    WithLimit(10).
    WithOffset(0).
    WithSortBy("created_at").
    WithSortOrder("desc"))
if err != nil {
    log.Fatal(err)
}

for _, person := range people {
    name, _, _ := store.AttributeGetString(ctx, person.ID(), "name")
    fmt.Println(name)
}
```

### Search

```go
results, err := store.EntityList(ctx, entitystore.EntityQuery().
    WithEntityType("person").
    WithSearch("john"). // Searches across attributes
    WithLimit(20))
```

### Find by Attribute

```go
// Find first entity with matching attribute
admin, err := store.EntityFindByAttribute(ctx, "user", "role", "admin")

// Find all entities with matching attribute
admins, err := store.EntityListByAttribute(ctx, "user", "role", "admin")
```

### Count

```go
count, err := store.EntityCount(ctx, entitystore.EntityQuery().
    WithEntityType("product"))
fmt.Printf("Total products: %d\n", count)
```

## Updating Entities

```go
// Retrieve
entity, _ := store.EntityFindByID(ctx, "abc123xyz")

// Modify core fields and persist
entity.SetHandle("new-handle")
err := store.EntityUpdate(ctx, entity)

// Modify attributes
store.AttributeSetString(ctx, entity.ID(), "status", "active")
store.AttributeSetInt(ctx, entity.ID(), "login_count", 5)
```

## Deleting Entities

### Hard Delete

Permanently removes the entity and all its attributes:

```go
deleted, err := store.EntityDelete(ctx, "abc123xyz")
if deleted {
    fmt.Println("Entity permanently deleted")
}
```

### Soft Delete (Trash)

Moves entity to trash bin for potential recovery:

```go
trashed, err := store.EntityTrash(ctx, "abc123xyz")
if trashed {
    fmt.Println("Entity moved to trash")
}

// Restore
err = store.EntityRestore(ctx, "abc123xyz")
```

## Complete Example

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    
    "github.com/dracory/entitystore"
    _ "modernc.org/sqlite"
)

func main() {
    db, _ := sql.Open("sqlite", "example.db")
    defer db.Close()
    
    store, _ := entitystore.NewStore(entitystore.NewStoreOptions{
        DB:                 db,
        EntityTableName:    "entities",
        AttributeTableName: "attributes",
        AutomigrateEnabled: true,
    })
    
    ctx := context.Background()
    
    // Create
    person, _ := store.EntityCreateWithType(ctx, "person")
    store.AttributeSetString(ctx, person.ID(), "name", "Alice Smith")
    store.AttributeSetInt(ctx, person.ID(), "age", 28)
    
    // Read
    found, _ := store.EntityFindByID(ctx, person.ID())
    name, _, _ := store.AttributeGetString(ctx, found.ID(), "name")
    age, _, _ := store.AttributeGetInt(ctx, found.ID(), "age")
    fmt.Printf("Found: %s (age %d)\n", name, age)
    
    // Update
    store.AttributeSetString(ctx, found.ID(), "name", "Alice Johnson")
    
    // List
    people, _ := store.EntityList(ctx, entitystore.EntityQuery().
        WithEntityType("person"))
    fmt.Printf("Total people: %d\n", len(people))
    
    // Soft delete
    store.EntityTrash(ctx, person.ID())
    fmt.Println("Person moved to trash")
}
```

## Next Steps

- [Entities](entities.md) - Complete entity documentation
- [Attributes](attributes.md) - Working with typed attributes
- [Relationships](entity-relationships.md) - Linking entities
- [Taxonomies](taxonomies.md) - Categorizing entities
- [Architecture](architecture.md) - Understanding the design
