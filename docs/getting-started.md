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

// Create a person entity
person := store.EntityCreateWithType("person")
person.SetString("name", "John Doe")
person.SetInt("age", 30)

// Persist to database
err := store.EntityCreate(ctx, person)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Created person with ID:", person.ID())
```

### With Multiple Attributes

```go
// Create with attributes map
attrs := map[string]string{
    "name":  "Jane Doe",
    "email": "jane@example.com",
    "role":  "admin",
}

user := store.EntityCreateWithTypeAndAttributes("user", attrs)
store.EntityCreate(ctx, user)
```

## Working with Attributes

Attributes store typed data for entities:

### Setting Attributes

```go
entity := store.EntityCreateWithType("product")

// String
entity.SetString("name", "Laptop")
entity.SetString("sku", "LAP-001")

// Numbers
entity.SetInt("stock", 50)
entity.SetFloat("price", 999.99)

// Complex data (JSON-serialized)
tags := []string{"electronics", "computers"}
entity.SetInterface("tags", tags)

specs := map[string]string{
    "cpu": "Intel i7",
    "ram": "16GB",
}
entity.SetInterface("specs", specs)

store.EntityCreate(ctx, entity)
```

### Getting Attributes

```go
// Retrieve entity
product, _ := store.EntityFindByID(ctx, "abc123xyz")

// Get with defaults
name := product.GetString("name", "Unknown")
stock, _ := product.GetInt("stock", 0)
price, _ := product.GetFloat("price", 0.0)

// Get complex data
tags := product.GetInterface("tags", []string{}).([]string)
specs := product.GetInterface("specs", map[string]string{}).(map[string]string)
```

## Querying Entities

### Find by ID

```go
entity, err := store.EntityFindByID(ctx, "abc123xyz")
if err != nil {
    log.Fatal(err)
}
if entity != nil {
    fmt.Println("Found:", entity.GetString("name", ""))
}
```

### List by Type

```go
people, err := store.EntityList(ctx, entitystore.EntityQueryOptions{
    EntityType: "person",
    Limit:      10,
    Offset:     0,
    SortBy:     "created_at",
    SortOrder:  "desc",
})
if err != nil {
    log.Fatal(err)
}

for _, person := range people {
    fmt.Println(person.GetString("name", ""))
}
```

### Search

```go
results, err := store.EntityList(ctx, entitystore.EntityQueryOptions{
    EntityType: "person",
    Search:     "john", // Searches across attributes
    Limit:      20,
})
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
count, err := store.EntityCount(ctx, entitystore.EntityQueryOptions{
    EntityType: "product",
})
fmt.Printf("Total products: %d\n", count)
```

## Updating Entities

```go
// Retrieve
entity, _ := store.EntityFindByID(ctx, "abc123xyz")

// Modify
entity.SetString("status", "active")
entity.SetInt("login_count", 5)

// Persist changes
err := store.EntityUpdate(ctx, entity)
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

// List trashed entities
trashed, _ := store.EntityTrashList(ctx, entitystore.EntityQueryOptions{
    Limit: 10,
})

// Restore
restored, _ := store.EntityRestore(ctx, "abc123xyz")
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
    person := store.EntityCreateWithType("person")
    person.SetString("name", "Alice Smith")
    person.SetInt("age", 28)
    store.EntityCreate(ctx, person)
    
    // Read
    found, _ := store.EntityFindByID(ctx, person.ID())
    fmt.Printf("Found: %s (age %d)\n", 
        found.GetString("name", ""),
        mustInt(found.GetInt("age", 0)))
    
    // Update
    found.SetString("name", "Alice Johnson")
    store.EntityUpdate(ctx, found)
    
    // List
    people, _ := store.EntityList(ctx, entitystore.EntityQueryOptions{
        EntityType: "person",
    })
    fmt.Printf("Total people: %d\n", len(people))
    
    // Soft delete
    store.EntityTrash(ctx, person.ID())
    fmt.Println("Person moved to trash")
}

func mustInt(i int64, _ error) int64 { return i }
```

## Next Steps

- [Entities](entities.md) - Complete entity documentation
- [Attributes](attributes.md) - Working with typed attributes
- [Relationships](entity-relationships.md) - Linking entities
- [Taxonomies](taxonomies.md) - Categorizing entities
- [Architecture](architecture.md) - Understanding the design
