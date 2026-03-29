# Proposal: Active Record Entity Wrapper Package

**Date:** 2026-03-29  
**Status:** PROPOSED  
**Author:** AI Assistant  

## Summary

This proposal introduces an optional `activerecord` (or `rich`) sub-package that provides an object-oriented wrapper around standard `entitystore` entities. This wrapper approach satisfies the desire for an intuitive, fluent API (e.g., `entity.SetString().Save()`) without tightly coupling the core `EntityInterface` to the `StoreInterface` or introducing memory overhead for users who prefer the pure Data Mapper approach.

## Problem Statement

A previously declined proposal suggested adding `.Save()`, `.SetString()`, and other store-bound methods directly into the core `EntityInterface` and `entityImplementation`. This was rejected because:
1. It violated the separation of concerns (mixing pure data structs with I/O).
2. It introduced a "detached entity" problem for newly created entities.
3. It created memory overhead for all entities.
4. It made core interface serialization, caching, and mocking more difficult.

However, developers still find the object-oriented or "Active Record" style API highly desirable for rapid development, UI scripting, reducing boilerplate, and simpler code generation.

## Proposed Solution

Create a new package alongside the `entitystore` module, for example, `package activerecord`. 

This package defines an `ActiveEntity` interface that wraps both an `EntityInterface` and a `StoreInterface`:

```go
package activerecord

import (
	"context"
	"github.com/yourorg/entitystore" // adjust path
)

// ActiveEntity defines the rich, store-aware API for entities.
type ActiveEntity interface {
	SetString(ctx context.Context, key, value string) error
	GetString(ctx context.Context, key string) (string, bool, error)
	GetAttributes(ctx context.Context) ([]entitystore.AttributeInterface, error)
	Save(ctx context.Context) error
	Target() entitystore.EntityInterface
}

// activeEntityImpl is the concrete wrapper.
type activeEntityImpl struct {
	entity entitystore.EntityInterface
	store  entitystore.StoreInterface
}

// Wrap converts a standard entity into an ActiveEntity.
func Wrap(store entitystore.StoreInterface, entity entitystore.EntityInterface) ActiveEntity {
	if entity == nil {
		return nil
	}
	return &activeEntityImpl{
		entity: entity,
		store:  store,
	}
}

// New creates a fresh ActiveEntity bound to the store.
func New(store entitystore.StoreInterface) ActiveEntity {
	return &activeEntityImpl{
		entity: entitystore.NewEntity(), // Assuming NewEntity exists
		store:  store,
	}
}
```

### Delegation & Sugar Methods

The `activeEntityImpl` implements all the "sugar" methods that delegate their actual storage operations to the `store` while operating on the underlying `entity` ID:

```go
// SetString sets a string attribute on the entity.
func (e *activeEntityImpl) SetString(ctx context.Context, key, value string) error {
	return e.store.AttributeSetString(ctx, e.entity.ID(), key, value)
}

// GetString gets a string attribute from the entity via the store.
func (e *activeEntityImpl) GetString(ctx context.Context, key string) (string, bool, error) {
	return e.store.AttributeGetString(ctx, e.entity.ID(), key)
}

// Save persists the underlying entity to the store.
func (e *activeEntityImpl) Save(ctx context.Context) error {
	// Simple example logic:
	existing, err := e.store.EntityFindByID(ctx, e.entity.ID())
	if err != nil {
		return err
	}
	if existing == nil {
		return e.store.EntityCreate(ctx, e.entity)
	}
	return e.store.EntityUpdate(ctx, e.entity)
}

// Target returns the underlying pure data entity if the developer
// needs to pass it to a core method that expects EntityInterface.
func (e *activeEntityImpl) Target() entitystore.EntityInterface {
	return e.entity
}
```

### Store Helper Methods (Optional)

The wrapper package can also provide helpers to fetch entities directly as `ActiveEntity`:

```go
// FindByID fetches an entity and wraps it.
func FindByID(ctx context.Context, store entitystore.StoreInterface, id string) (ActiveEntity, error) {
	ent, err := store.EntityFindByID(ctx, id)
	if err != nil || ent == nil {
		return nil, err
	}
	return Wrap(store, ent), nil
}

// GetAttributes returns all attributes wrapped or directly from store.
func (e *activeEntityImpl) GetAttributes(ctx context.Context) ([]entitystore.AttributeInterface, error) {
	return e.store.EntityAttributeList(ctx, e.entity.ID())
}
```

## Benefits

1. **Zero Core Pollution:** The core `EntityInterface` and `StoreInterface` remain 100% clean, decoupled from one another, and easy to mock.
2. **Opt-In Overhead:** The memory overhead of storing the interface pointer bindings only applies to developers who explicitly choose to use the `activerecord` package.
3. **Safe Entity Lifecycles:** A developer can never call `.SetString()` on an unbound entity because an `ActiveEntity` cannot be safely instantiated without explicitly passing the store parameter.
4. **Best of Both Worlds:** Developers who prefer strict Data Mapper patterns use the core package. Developers prototyping or writing UI scripts can use the `activerecord` package to save typing and boilerplate.

## Drawbacks & Considerations

- **Package Discoverability:** Developers need to be aware that the `activerecord` package exists; otherwise, they might complain the core API is too verbose.
- **Double API Surface:** The maintainers of this library will effectively support two public APIs. The `ActiveEntity` interface will need roughly the same amount of methods as the core `StoreInterface` simply to delegate them. However, since it consists of simple passthrough methods, testing and maintenance should be trivial.

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
product := activerecord.New(store)
product.Target().SetType("product")
_ = product.Save(ctx)
_ = product.SetString(ctx, "name", "Laptop")
_ = product.SetInt(ctx, "stock", 50)
```

## Migration

There are no backward compatibility issues or breaking changes. This is an entirely new, opt-in wrapper package built securely on top of the existing, stable public API.
