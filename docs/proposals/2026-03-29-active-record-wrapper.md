# Proposal: Active Record Entity Wrapper Package

**Date:** 2026-03-29  
**Status:** PROPOSED  
**Author:** AI Assistant  

## Summary

This proposal introduces an optional `activestore` sub-package that provides an object-oriented wrapper around standard `entitystore` entities. This wrapper approach satisfies the desire for an intuitive, fluent API (e.g., `entity.SetString().Save()`) without tightly coupling the core `EntityInterface` to the `StoreInterface` or introducing memory overhead for users who prefer the pure Data Mapper approach.

## Problem Statement

A previously declined proposal suggested adding `.Save()`, `.SetString()`, and other store-bound methods directly into the core `EntityInterface` and `entityImplementation`. This was rejected because:
1. It violated the separation of concerns (mixing pure data structs with I/O).
2. It introduced a "detached entity" problem for newly created entities.
3. It created memory overhead for all entities.
4. It made core interface serialization, caching, and mocking more difficult.

However, developers still find the object-oriented or "Active Record" style API highly desirable for rapid development, UI scripting, reducing boilerplate, and simpler code generation.

## Proposed Solution

Create a new package alongside the `entitystore` module: `package activestore`. 

This package defines an `ActiveEntityInterface` interface that wraps both an `EntityInterface` and a `StoreInterface`:

```go
package activestore

import (
	"context"
	"errors"
	"github.com/yourorg/entitystore" // adjust path
)

// ActiveEntityInterface defines the rich, store-aware API for entities.
// Mutating methods return the ActiveEntityInterface for chaining; errors are
// accumulated internally and surfaced by Save() or Err().
type ActiveEntityInterface interface {
	SetString(key, value string) ActiveEntityInterface
	SetInt(key string, value int64) ActiveEntityInterface
	SetFloat(key string, value float64) ActiveEntityInterface
	GetString(key string) (string, bool, error)
	GetAttributes() ([]entitystore.AttributeInterface, error)
	Prefetch() error
	Save() (ActiveEntityInterface, error)
	Trash() (bool, error)
	Delete() (bool, error)
	Err() error
	GetEntity() entitystore.EntityInterface

	// Relationships (requires RelationshipsEnabled)
	RelateTo(relatedEntityID, relationshipType string) error
	RelateToOrdered(relatedEntityID, relationshipType string, sequence int) error
	Unrelate(relatedEntityID, relationshipType string) error
	Related(relationshipType string) ([]ActiveEntityInterface, error)

	// Taxonomies (requires TaxonomiesEnabled)
	AssignTerm(taxonomyID, termID string) error
	RemoveTerm(taxonomyID, termID string) error
	Terms(taxonomyID string) ([]ActiveTaxonomyTermInterface, error)
}

// activeEntityImplementation is the concrete wrapper.
type activeEntityImplementation struct {
	ctx        context.Context
	entity     entitystore.EntityInterface
	store      entitystore.StoreInterface
	persisted  bool        // false until Save() creates the row
	pendingOps []pendingOp // attribute writes staged for unpersisted entities
	err        error       // first accumulated error, deferred to Save()/Err()
}

// newActiveEntity binds an entity to a store context. Entities are only
// created through ActiveStoreInterface, so a store is always present.
func newActiveEntity(ctx context.Context, store entitystore.StoreInterface, entity entitystore.EntityInterface) (ActiveEntityInterface, error) {
	if entity == nil {
		return nil, errors.New("entity cannot be nil")
	}
	return &activeEntityImplementation{
		ctx:    ctx,
		entity: entity,
		store:  store,
	}, nil
}
```

### Store-Level Wrapper

For discoverability and a single entry point, the package can also wrap the store itself. `ActiveStoreInterface` exposes store operations that return `ActiveEntityInterface` values directly:

```go
// ActiveStoreInterface is a store-aware facade that returns ActiveEntityInterface results.
type ActiveStoreInterface interface {
	GetStore() entitystore.StoreInterface

	EntityCreate(entityType string) ActiveEntityInterface
	EntityFindByID(entityID string) (ActiveEntityInterface, error)
	EntityList(query entitystore.EntityQueryInterface) ([]ActiveEntityInterface, error)
	EntityCount(query entitystore.EntityQueryInterface) (int64, error)
	EntityTrash(entityID string) (bool, error)
	EntityDelete(entityID string) (bool, error)
	EntityWrap(entity entitystore.EntityInterface) (ActiveEntityInterface, error)

	// Relationships (requires RelationshipsEnabled)
	RelationshipCreate(options entitystore.RelationshipOptions) (ActiveRelationshipInterface, error)
	RelationshipFindByID(relationshipID string) (ActiveRelationshipInterface, error)
	RelationshipList(query entitystore.RelationshipQueryInterface) ([]ActiveRelationshipInterface, error)
	RelationshipCount(query entitystore.RelationshipQueryInterface) (int64, error)
	RelationshipTrash(relationshipID, deletedBy string) (bool, error)
	RelationshipRestore(relationshipID string) (bool, error)
	RelationshipDelete(relationshipID string) (bool, error)
	RelationshipDeleteAll(entityID string) error

	// Taxonomies (requires TaxonomiesEnabled)
	TaxonomyCreate(options entitystore.TaxonomyOptions) (ActiveTaxonomyInterface, error)
	TaxonomyFindByID(taxonomyID string) (ActiveTaxonomyInterface, error)
	TaxonomyFindBySlug(slug string) (ActiveTaxonomyInterface, error)
	TaxonomyList(query entitystore.TaxonomyQueryInterface) ([]ActiveTaxonomyInterface, error)
	TaxonomyCount(query entitystore.TaxonomyQueryInterface) (int64, error)
	TaxonomyUpdate(taxonomy entitystore.TaxonomyInterface) error
	TaxonomyTrash(taxonomyID, deletedBy string) (bool, error)
	TaxonomyRestore(taxonomyID string) (bool, error)
	TaxonomyDelete(taxonomyID string) (bool, error)

	// Taxonomy terms (requires TaxonomiesEnabled)
	TermCreate(options entitystore.TaxonomyTermOptions) (ActiveTaxonomyTermInterface, error)
	TermFindByID(termID string) (ActiveTaxonomyTermInterface, error)
	TermFindBySlug(taxonomyID, slug string) (ActiveTaxonomyTermInterface, error)
	TermList(query entitystore.TaxonomyTermQueryInterface) ([]ActiveTaxonomyTermInterface, error)
	TermCount(query entitystore.TaxonomyTermQueryInterface) (int64, error)
	TermUpdate(term entitystore.TaxonomyTermInterface) error
	TermTrash(termID, deletedBy string) (bool, error)
	TermRestore(termID string) (bool, error)
	TermDelete(termID string) (bool, error)
}

// activeStoreImplementation is the concrete store wrapper.
type activeStoreImplementation struct {
	ctx   context.Context
	store entitystore.StoreInterface
}

// New converts a standard store into an ActiveStoreInterface — the
// single entry point for the package.
func New(ctx context.Context, store entitystore.StoreInterface) (ActiveStoreInterface, error) {
	if store == nil {
		return nil, errors.New("store cannot be nil")
	}
	return &activeStoreImplementation{ctx: ctx, store: store}, nil
}

func (s *activeStoreImplementation) EntityCreate(entityType string) ActiveEntityInterface {
	entity := entitystore.NewEntity() // Assuming NewEntity exists
	entity.SetType(entityType)
	wrapped, _ := newActiveEntity(s.ctx, s.store, entity) // entity is never nil
	return wrapped
}

func (s *activeStoreImplementation) EntityFindByID(entityID string) (ActiveEntityInterface, error) {
	ent, err := s.store.EntityFindByID(s.ctx, entityID)
	if err != nil {
		return nil, err
	}
	return newActiveEntity(s.ctx, s.store, ent)
}

func (s *activeStoreImplementation) EntityList(query entitystore.EntityQueryInterface) ([]ActiveEntityInterface, error) {
	entities, err := s.store.EntityList(s.ctx, query)
	if err != nil {
		return nil, err
	}
	active := make([]ActiveEntityInterface, 0, len(entities))
	for _, ent := range entities {
		wrapped, err := newActiveEntity(s.ctx, s.store, ent)
		if err != nil {
			return nil, err
		}
		active = append(active, wrapped)
	}
	return active, nil
}

func (s *activeStoreImplementation) EntityCount(query entitystore.EntityQueryInterface) (int64, error) {
	return s.store.EntityCount(s.ctx, query)
}

func (s *activeStoreImplementation) EntityTrash(entityID string) (bool, error) {
	return s.store.EntityTrash(s.ctx, entityID)
}

func (s *activeStoreImplementation) EntityDelete(entityID string) (bool, error) {
	return s.store.EntityDelete(s.ctx, entityID)
}

func (s *activeStoreImplementation) EntityWrap(entity entitystore.EntityInterface) (ActiveEntityInterface, error) {
	return newActiveEntity(s.ctx, s.store, entity)
}

func (s *activeStoreImplementation) GetStore() entitystore.StoreInterface {
	return s.store
}
```

### Delegation & Sugar Methods

The `activeEntityImplementation` implements all the "sugar" methods that delegate their actual storage operations to the `store` while operating on the underlying `entity` ID:

```go
// SetString stages a string attribute. If the entity is already persisted
// the write goes to the store immediately; otherwise it is staged until Save().
// Any error is accumulated and returned by Save()/Err().
func (e *activeEntityImplementation) SetString(key, value string) ActiveEntityInterface {
	if e.err != nil {
		return e
	}
	e.err = e.store.AttributeSetString(e.ctx, e.entity.ID(), key, value)
	return e
}

// GetString gets a string attribute from the entity via the store.
func (e *activeEntityImplementation) GetString(key string) (string, bool, error) {
	return e.store.AttributeGetString(e.ctx, e.entity.ID(), key)
}

// Err returns the first error accumulated by the fluent setters, if any.
func (e *activeEntityImplementation) Err() error {
	return e.err
}

// Save persists the underlying entity to the store, then flushes any
// attribute writes that were staged while the entity was unpersisted.
// Returns the entity itself so it can be the last call in a fluent chain.
func (e *activeEntityImplementation) Save() (ActiveEntityInterface, error) {
	if e.err != nil {
		return e, e.err
	}
	if e.persisted {
		return e, e.store.EntityUpdate(e.ctx, e.entity)
	}
	if err := e.store.EntityCreate(e.ctx, e.entity); err != nil {
		return e, err
	}
	e.persisted = true
	for _, op := range e.pendingOps {
		if err := op(e.ctx, e.store, e.entity.ID()); err != nil {
			e.err = err
			return e, err
		}
	}
	e.pendingOps = nil
	return e, nil
}

// Trash moves the entity to the trash table.
func (e *activeEntityImplementation) Trash() (bool, error) {
	return e.store.EntityTrash(e.ctx, e.entity.ID())
}

// Delete permanently removes the entity.
func (e *activeEntityImplementation) Delete() (bool, error) {
	return e.store.EntityDelete(e.ctx, e.entity.ID())
}

// GetEntity returns the underlying pure data entity if the developer
// needs to pass it to a core method that expects EntityInterface.
func (e *activeEntityImplementation) GetEntity() entitystore.EntityInterface {
	return e.entity
}
```

### Store Helper Methods (Optional)

The only package-level constructor is `New(ctx, store)` returning `ActiveStoreInterface` — everything else hangs off it. Entity-level helpers:

```go
// GetAttributes returns all attributes wrapped or directly from store.
func (e *activeEntityImplementation) GetAttributes() ([]entitystore.AttributeInterface, error) {
	return e.store.EntityAttributeList(e.ctx, e.entity.ID())
}
```

### Relationships Sugar

The core API forces a three-step dance: build `RelationshipOptions` with two raw IDs, later list `RelationshipInterface` rows, then fetch each related entity by ID. Entity-centric methods collapse this:

```go
// On ActiveEntityInterface — the entity's own ID is implicit:
RelateTo(relatedEntityID, relationshipType string) error
RelateToOrdered(relatedEntityID, relationshipType string, sequence int) error
Unrelate(relatedEntityID, relationshipType string) error
Related(relationshipType string) ([]ActiveEntityInterface, error)
```

`Related(type)` is the main win — one call performs `RelationshipList` + `EntityFindByID` per row + wraps each result into `ActiveEntityInterface`.

Relationships themselves also get a wrapper so the developer never falls back to raw objects mid-flow:

```go
// ActiveRelationshipInterface wraps a RelationshipInterface with
// navigation and lifecycle helpers.
type ActiveRelationshipInterface interface {
	GetRelationship() entitystore.RelationshipInterface
	GetEntity() (ActiveEntityInterface, error)        // source entity, hydrated
	GetRelatedEntity() (ActiveEntityInterface, error) // target entity, hydrated
	Trash(deletedBy string) (bool, error)
	Restore() (bool, error)
	Delete() (bool, error)
}

// On ActiveStoreInterface — the full CRUD set mirrors core names:
RelationshipCreate(options entitystore.RelationshipOptions) (ActiveRelationshipInterface, error)
RelationshipFindByID(relationshipID string) (ActiveRelationshipInterface, error)
RelationshipList(query entitystore.RelationshipQueryInterface) ([]ActiveRelationshipInterface, error)
RelationshipCount(query entitystore.RelationshipQueryInterface) (int64, error)
RelationshipTrash(relationshipID, deletedBy string) (bool, error)
RelationshipRestore(relationshipID string) (bool, error)
RelationshipDelete(relationshipID string) (bool, error)
RelationshipDeleteAll(entityID string) error
```

`RelationshipCreate` takes `RelationshipOptions` and delegates to `RelationshipCreateByOptions` — callers should not have to build a `RelationshipInterface` first.

```go
post.RelateTo(authorID, "written_by")
post.RelateToOrdered(tagID, "has_tag", 3)

authors, _ := post.Related("written_by")   // hydrated, wrapped entities
tags, _ := post.Related("has_tag")
post.Unrelate(tagID, "has_tag")
```

### Taxonomy Sugar

Taxonomies juggle three objects (taxonomy, term, assignment) and three IDs. Assignment/removal becomes entity-centric; vocabulary creation stays options-struct based:

Taxonomies and terms get their own wrappers, so navigation (terms of a taxonomy, entities under a term, parent/children) becomes method calls instead of query construction:

```go
// ActiveTaxonomyInterface wraps a TaxonomyInterface.
type ActiveTaxonomyInterface interface {
	GetTaxonomy() entitystore.TaxonomyInterface
	Terms() ([]ActiveTaxonomyTermInterface, error)  // terms in this taxonomy
	Entities() ([]ActiveEntityInterface, error)     // assigned entities
	Trash(deletedBy string) (bool, error)
	Restore() (bool, error)
	Delete() (bool, error)
}

// ActiveTaxonomyTermInterface wraps a TaxonomyTermInterface.
type ActiveTaxonomyTermInterface interface {
	GetTerm() entitystore.TaxonomyTermInterface
	GetTaxonomy() (ActiveTaxonomyInterface, error)
	Parent() (ActiveTaxonomyTermInterface, error)   // hierarchical nav
	Children() ([]ActiveTaxonomyTermInterface, error)
	Entities() ([]ActiveEntityInterface, error)     // entities tagged with this term
	Trash(deletedBy string) (bool, error)
	Restore() (bool, error)
	Delete() (bool, error)
}

// On ActiveStoreInterface — vocabulary management, full CRUD mirrors
// core names. Create/Find return the active wrappers:
TaxonomyCreate(options entitystore.TaxonomyOptions) (ActiveTaxonomyInterface, error)
TaxonomyFindByID(taxonomyID string) (ActiveTaxonomyInterface, error)
TaxonomyFindBySlug(slug string) (ActiveTaxonomyInterface, error)
TaxonomyList(query entitystore.TaxonomyQueryInterface) ([]ActiveTaxonomyInterface, error)
TaxonomyCount(query entitystore.TaxonomyQueryInterface) (int64, error)
TaxonomyUpdate(taxonomy entitystore.TaxonomyInterface) error
TaxonomyTrash(taxonomyID, deletedBy string) (bool, error)
TaxonomyRestore(taxonomyID string) (bool, error)
TaxonomyDelete(taxonomyID string) (bool, error)

TermCreate(options entitystore.TaxonomyTermOptions) (ActiveTaxonomyTermInterface, error)
TermFindByID(termID string) (ActiveTaxonomyTermInterface, error)
TermFindBySlug(taxonomyID, slug string) (ActiveTaxonomyTermInterface, error)
TermList(query entitystore.TaxonomyTermQueryInterface) ([]ActiveTaxonomyTermInterface, error)
TermCount(query entitystore.TaxonomyTermQueryInterface) (int64, error)
TermUpdate(term entitystore.TaxonomyTermInterface) error
TermTrash(termID, deletedBy string) (bool, error)
TermRestore(termID string) (bool, error)
TermDelete(termID string) (bool, error)

// On ActiveEntityInterface — the entity's ID is implicit:
AssignTerm(taxonomyID, termID string) error
RemoveTerm(taxonomyID, termID string) error
Terms(taxonomyID string) ([]ActiveTaxonomyTermInterface, error)
```

```go
cat, _ := products.TaxonomyCreate(entitystore.TaxonomyOptions{Name: "Categories", Slug: "categories"})
term, _ := products.TermCreate(entitystore.TaxonomyTermOptions{TaxonomyID: cat.GetTaxonomy().ID(), Name: "Laptops", Slug: "laptops"})

product.AssignTerm(cat.GetTaxonomy().ID(), term.GetTerm().ID())
terms, _ := product.Terms(cat.GetTaxonomy().ID())
laptops, _ := term.Entities()                    // entities tagged "laptops"
product.RemoveTerm(cat.GetTaxonomy().ID(), term.GetTerm().ID())
```

Both features are opt-in at the core store level (`RelationshipsEnabled`, `TaxonomiesEnabled`); the wrappers delegate, so a disabled feature surfaces the store's own error naturally.

## Benefits

1. **Zero Core Pollution:** The core `EntityInterface` and `StoreInterface` remain 100% clean, decoupled from one another, and easy to mock.
2. **Opt-In Overhead:** The memory overhead of storing the interface pointer bindings only applies to developers who explicitly choose to use the `activestore` package.
3. **Safe Entity Lifecycles:** A developer can never call `.SetString()` on an unbound entity because an `ActiveEntityInterface` cannot be safely instantiated without explicitly passing the store parameter.
4. **Best of Both Worlds:** Developers who prefer strict Data Mapper patterns use the core package. Developers prototyping or writing UI scripts can use the `activestore` package to save typing and boilerplate.

## Drawbacks & Considerations

- **Staged Writes for New Entities:** Attributes on an unpersisted entity cannot be written to the store until the entity row exists. `EntityCreate()` entities must stage attribute writes in memory and flush them inside `Save()` after the row insert succeeds. Entities obtained via `EntityFindByID`/`EntityWrap`/`EntityList` can write immediately.
- **Deferred Errors:** The fluent setters accumulate the first error instead of returning it (required for chaining). Developers must check `Save()` or `Err()` — a silent failure mode if they forget.
- **Prefetch Staleness:** `Prefetch()` caches all attributes in memory so `GetString` avoids per-call DB round-trips. The cache is updated by the wrapper's own setters but goes stale if another process writes to the store directly; it is opt-in per entity.
- **Package Discoverability:** Developers need to be aware that the `activestore` package exists; otherwise, they might complain the core API is too verbose.
- **Double API Surface:** The maintainers of this library will effectively support two public APIs. The `ActiveEntityInterface` interface will need roughly the same amount of methods as the core `StoreInterface` simply to delegate them. However, since it consists of simple passthrough methods, testing and maintenance should be trivial.

## Usage Comparison

```go
// ------------------------------------------------
// 1. Pure approach (Core Data Mapper)
// ------------------------------------------------
ctx := context.Background()
product := entitystore.NewEntity()
product.SetType("product")
_ = store.EntityCreate(ctx, product)
_ = store.AttributeSetString(ctx, product.ID(), "name", "Laptop")
_ = store.AttributeSetInt(ctx, product.ID(), "stock", 50)


// ------------------------------------------------
// 2. Active Record Wrapper approach (Rich API)
// ------------------------------------------------
ctx := context.Background()
products, _ := activestore.New(ctx, store)

product, err := products.EntityCreate("product").
	SetString("name", "Laptop").
	SetFloat("price", 1299.99).
	SetInt("stock", 50).
	Save()

existing, _ := products.EntityFindByID(entityID)
name, _, _ := existing.GetString("name")
```

## Migration

There are no backward compatibility issues or breaking changes. This is an entirely new, opt-in wrapper package built securely on top of the existing, stable public API.
