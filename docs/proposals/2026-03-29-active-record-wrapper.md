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
	Save() error
	Trash() (bool, error)
	Delete() (bool, error)
	Err() error
	GetEntity() entitystore.EntityInterface
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
	New(entityType string) ActiveEntityInterface
	FindByID(entityID string) (ActiveEntityInterface, error)
	List(query entitystore.EntityQueryInterface) ([]ActiveEntityInterface, error)
	Count(query entitystore.EntityQueryInterface) (int64, error)
	Trash(entityID string) (bool, error)
	Delete(entityID string) (bool, error)
	WrapEntity(entity entitystore.EntityInterface) (ActiveEntityInterface, error)
	GetStore() entitystore.StoreInterface
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

func (s *activeStoreImplementation) New(entityType string) ActiveEntityInterface {
	entity := entitystore.NewEntity() // Assuming NewEntity exists
	entity.SetType(entityType)
	wrapped, _ := newActiveEntity(s.ctx, s.store, entity) // entity is never nil
	return wrapped
}

func (s *activeStoreImplementation) FindByID(entityID string) (ActiveEntityInterface, error) {
	ent, err := s.store.EntityFindByID(s.ctx, entityID)
	if err != nil {
		return nil, err
	}
	return newActiveEntity(s.ctx, s.store, ent)
}

func (s *activeStoreImplementation) List(query entitystore.EntityQueryInterface) ([]ActiveEntityInterface, error) {
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

func (s *activeStoreImplementation) Count(query entitystore.EntityQueryInterface) (int64, error) {
	return s.store.EntityCount(s.ctx, query)
}

func (s *activeStoreImplementation) Trash(entityID string) (bool, error) {
	return s.store.EntityTrash(s.ctx, entityID)
}

func (s *activeStoreImplementation) Delete(entityID string) (bool, error) {
	return s.store.EntityDelete(s.ctx, entityID)
}

func (s *activeStoreImplementation) WrapEntity(entity entitystore.EntityInterface) (ActiveEntityInterface, error) {
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
func (e *activeEntityImplementation) Save() error {
	if e.err != nil {
		return e.err
	}
	if e.persisted {
		return e.store.EntityUpdate(e.ctx, e.entity)
	}
	if err := e.store.EntityCreate(e.ctx, e.entity); err != nil {
		return err
	}
	e.persisted = true
	for _, op := range e.pendingOps {
		if err := op(e.ctx, e.store, e.entity.ID()); err != nil {
			e.err = err
			return err
		}
	}
	e.pendingOps = nil
	return nil
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

## Benefits

1. **Zero Core Pollution:** The core `EntityInterface` and `StoreInterface` remain 100% clean, decoupled from one another, and easy to mock.
2. **Opt-In Overhead:** The memory overhead of storing the interface pointer bindings only applies to developers who explicitly choose to use the `activestore` package.
3. **Safe Entity Lifecycles:** A developer can never call `.SetString()` on an unbound entity because an `ActiveEntityInterface` cannot be safely instantiated without explicitly passing the store parameter.
4. **Best of Both Worlds:** Developers who prefer strict Data Mapper patterns use the core package. Developers prototyping or writing UI scripts can use the `activestore` package to save typing and boilerplate.

## Drawbacks & Considerations

- **Staged Writes for New Entities:** Attributes on an unpersisted entity cannot be written to the store until the entity row exists. `New()` entities must stage attribute writes in memory and flush them inside `Save()` after `EntityCreate` succeeds. Entities obtained via `FindByID`/`WrapEntity`/`List` can write immediately.
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

err := products.New("product").
	SetString("name", "Laptop").
	SetFloat("price", 1299.99).
	SetInt("stock", 50).
	Save()

existing, _ := products.FindByID(entityID)
name, _, _ := existing.GetString("name")
```

## Migration

There are no backward compatibility issues or breaking changes. This is an entirely new, opt-in wrapper package built securely on top of the existing, stable public API.
