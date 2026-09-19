# ActiveStore Example

This example demonstrates the `activestore` package — an Active Record style wrapper around EntityStore that provides a fluent, object-oriented API.

## What This Example Shows

### 1. Store Initialization
- Creating the core `entitystore` store with SQLite (in-memory)
- Wrapping it once with `activestore.New(ctx, store)` — the single entry point

### 2. Fluent Entity Creation
- `active.EntityCreate("product")` creates an unpersisted entity
- `SetString`/`SetInt`/`SetFloat` chain and stage attribute writes
- `Save()` creates the row and flushes staged attributes
- Errors accumulate and surface via `Save()` or `Err()`

### 3. Working with Persisted Entities
- `EntityFindByID` returns an `ActiveEntityInterface` that writes immediately
- `GetString`/`GetAttributes` read attributes via the store
- `GetEntity()` unwraps the underlying `EntityInterface` when needed

### 4. Queries and Lifecycle
- `EntityList(query)`/`EntityCount(query)` accept the same `EntityQueryInterface` as the core store
- `Trash()`/`Delete()` work on the entity itself — no ID juggling

### 5. Relationships (requires `RelationshipsEnabled`)
- `RelateTo(id, type)`/`Unrelate(id, type)` create and remove links
- `Related(type)` returns hydrated `ActiveEntityInterface` values — no manual `RelationshipList` + `EntityFindByID` loop

### 6. Taxonomies (requires `TaxonomiesEnabled`)
- `TaxonomyCreate`/`TermCreate` return active wrappers with navigation methods
- `AssignTerm`/`RemoveTerm`/`Terms` work on the entity — no ID juggling
- `term.Entities()` lists everything tagged with that term

## Running the Example

```bash
go run examples/activestore/main.go
```

## Running Tests

```bash
go test ./examples/activestore/... -v
```

## Code Highlights

```go
active, _ := activestore.New(ctx, store)

// Fluent create — attributes staged until Save()
err := active.EntityCreate("product").
    SetString("name", "Laptop").
    SetFloat("price", 1299.99).
    SetInt("stock", 50).
    Save()

// Fetch and update — writes apply immediately
product, _ := active.EntityFindByID(id)
product.SetInt("stock", 45)
product.Trash()

// Relationships — hydrated in one call
_ = post.RelateTo(author.GetEntity().ID(), "written_by")
authors, _ := post.Related("written_by")

// Taxonomies — navigate instead of querying
cat, _ := active.TaxonomyCreate(entitystore.TaxonomyOptions{Slug: "categories"})
term, _ := active.TermCreate(entitystore.TaxonomyTermOptions{TaxonomyID: cat.GetTaxonomy().GetID(), Slug: "go"})
_ = post.AssignTerm(cat.GetTaxonomy().GetID(), term.GetTerm().GetID())
posts, _ := term.Entities()
```

## Comparison with Core API

```go
// Core (Data Mapper)
product := entitystore.NewEntity()
product.SetType("product")
_ = store.EntityCreate(ctx, product)
_ = store.AttributeSetString(ctx, product.ID(), "name", "Laptop")

// ActiveStore (Active Record)
_ = active.EntityCreate("product").SetString("name", "Laptop").Save()
```
