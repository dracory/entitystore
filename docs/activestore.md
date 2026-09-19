# ActiveStore

`activestore` is an optional Active Record style wrapper around EntityStore. It binds a `StoreInterface` to rich, store-aware objects so attribute writes, persistence, relationships, and taxonomy operations can be performed with a fluent, object-oriented API.

Use it when you want less boilerplate (prototyping, UI scripts, codegen). Use the core package when you want strict Data Mapper separation.

## Entry Point

Everything hangs off a single constructor:

```go
store, _ := entitystore.NewStore(entitystore.NewStoreOptions{ /* ... */ })
active, _ := activestore.New(ctx, store) // ActiveStoreInterface
```

`ctx` is bound once at construction — no need to pass it to every call.

## Entities

```go
// Fluent create — setters chain, attributes are staged until Save()
product, err := active.EntityCreate("product").
    SetString("name", "Laptop").
    SetFloat("price", 1299.99).
    SetInt("stock", 50).
    Save()

// Fetch — returns a persisted ActiveEntityInterface; writes apply immediately
product, _ := active.EntityFindByID(id)
product.SetString("name", "Tablet")
product.Err()   // first accumulated error from setters, if any

// Queries mirror the core names
list, _ := active.EntityList(entitystore.EntityQuery().WithEntityType("product"))
count, _ := active.EntityCount(entitystore.EntityQuery().WithEntityType("product"))

// Lifecycle on the entity itself — no ID juggling
product.Trash()
product.Delete()

// Escape hatches back to the raw layer
product.GetEntity()  // entitystore.EntityInterface
active.GetStore()    // entitystore.StoreInterface
```

### Key semantics

- **Staged writes** — `EntityCreate` returns an unpersisted entity. Setters stage their writes in memory; `Save()` inserts the row then flushes them. Persisted entities (from `EntityFindByID`/`EntityList`/`EntityWrap`) write immediately.
- **Deferred errors** — fluent setters can't return `(entity, error)` and still chain, so the first error is accumulated and surfaced by `Save()` or `Err()`. Always check one of them.
- **`Save()` returns the entity** — `product, err := active.EntityCreate("x").SetString(...).Save()` works as a complete expression.
- **`Prefetch()`** — loads all attributes in one query into a per-entity cache; subsequent `GetString` reads from memory and setters keep the cache in sync. Without it, every `GetString` is a DB round-trip. Opt-in because the cache can go stale if another process writes directly.
- **Not found is `(nil, nil)`** — matching the core store convention, `Find*` methods (`EntityFindByID`, `RelationshipFindByID`, `TaxonomyFindByID/Slug`, `TermFindByID/Slug`) and navigation methods (`GetEntity`, `GetRelatedEntity`, `GetTaxonomy`, `Parent`) return `(nil, nil)` when the record doesn't exist — a missing record is not an error. Hydration methods (`Related`, `Terms`, `Entities`) silently skip dangling references whose target was deleted.

## Relationships

Requires `RelationshipsEnabled: true` in `NewStoreOptions`.

```go
post.RelateTo(author, "written_by")        // object variant
post.RelateToOrdered(tag, "has_tag", 3)
post.Unrelate(tag, "has_tag")

authors, _ := post.Related("written_by")   // hydrated []ActiveEntityInterface

// Raw-ID variants when you only have IDs
post.RelateToID(authorID, "written_by")
post.UnrelateID(authorID, "written_by")
```

Store-level CRUD mirrors the core names and returns `ActiveRelationshipInterface` (with `GetEntity()`/`GetRelatedEntity()` hydration):

```go
rel, _ := active.RelationshipCreate(entitystore.RelationshipOptions{...})
related, _ := rel.GetRelatedEntity()
list, _ := active.RelationshipList(entitystore.RelationshipQuery().WithEntityID(id))
_ = active.RelationshipCount / RelationshipTrash / RelationshipRestore /
    RelationshipDelete / RelationshipDeleteAll / RelationshipFindByID
```

## Taxonomies

Requires `TaxonomiesEnabled: true` in `NewStoreOptions`.

```go
cat, _ := active.TaxonomyCreate(entitystore.TaxonomyOptions{Slug: "categories"})
term, _ := active.TermCreate(entitystore.TaxonomyTermOptions{
    TaxonomyID: cat.GetTaxonomy().GetID(), Slug: "go",
})

post.AssignTerm(cat, term)            // object variants
terms, _ := post.Terms(cat)
post.RemoveTerm(cat, term)

post.AssignTermByID(catID, termID)    // raw-ID variants
post.RemoveTermByID(catID, termID)
post.TermsByID(catID)

// Navigation
terms, _ := cat.Terms()               // terms in a taxonomy
entities, _ := term.Entities()        // entities tagged with a term
parent, _ := term.Parent()
children, _ := term.Children()
tax, _ := term.GetTaxonomy()
```

Store-level CRUD mirrors the core names: `TaxonomyCreate/FindByID/FindBySlug/List/Count/Update/Trash/Restore/Delete` and `TermCreate/FindByID/FindBySlug/List/Count/Update/Trash/Restore/Delete` — returning `ActiveTaxonomyInterface` / `ActiveTaxonomyTermInterface` where applicable.

## Naming Convention

- **Store-level** methods mirror core `StoreInterface` names: `Entity*`, `Relationship*`, `Taxonomy*`, `Term*`
- **Entity-level** methods use domain sugar names: `RelateTo`, `Related`, `AssignTerm`, `Terms`
- **`*ID`-suffixed** variants accept raw IDs instead of active wrappers
- **`Get*`** methods are escape hatches returning the wrapped raw interface: `GetEntity`, `GetRelationship`, `GetTaxonomy`, `GetTerm`, `GetStore`

## Examples

- [`examples/activestore/`](../examples/activestore/) — entities, attributes, lifecycle
- [`examples/activestore-relationships/`](../examples/activestore-relationships/) — relationships
- [`examples/activestore-taxonomy/`](../examples/activestore-taxonomy/) — taxonomies

## Design Proposal

See [`docs/proposals/2026-03-29-active-record-wrapper.md`](proposals/2026-03-29-active-record-wrapper.md) for the full design rationale.
