# Proposal: Entity Object-Like Methods

**Date:** 2026-03-29  
**Status:** DECLINED  
**Author:** AI Assistant  

## Decision

This proposal has been **DECLINED** because the current architecture (Data Mapper / Repository pattern) is cleaner, more flexible, and memory-friendly. 

Specifically:
1. **Separation of Concerns:** The current design keeps entities as lightweight, pure data objects while delegating all database I/O and persistence logic to the Store. This avoids mixing I/O into domain models.
2. **Predictability:** Masking network or database calls behind simple getters and setters (e.g., `entity.SetString()`) hides the cost of these operations, increasing the risk of N+1 query problems and unexpected side effects.
3. **Memory Overhead & Detached Entities:** Tying entities to a `StoreInterface` adds unnecessary memory overhead per instantiated object and introduces complex lifecycle states (e.g., entities retrieved from the database being "bound" versus manually created entities being "unbound" and throwing runtime errors).
4. **Serialization Flexibility:** Pure data structs are much easier to serialize, cache, mock in unit tests, and pass across application layers without breaking or carrying internal database connections.

The current architecture provides a more explicit, transparent, and scalable foundation as the project grows.

## Summary

This proposal introduces object-like methods on entities (e.g., `entity.SetFloat()`, `entity.GetString()`) by adding a store reference to the entity implementation. This allows for a more intuitive, object-oriented API while maintaining backward compatibility.

## Current State

Currently, entities are pure data objects without any store reference. All attribute operations require going through the store:

```go
// Current pattern - requires store reference and context every time
person, _ := store.EntityCreateWithType(ctx, "person")
_ = store.AttributeSetString(ctx, person.ID(), "name", "John")
_ = store.AttributeSetInt(ctx, person.ID(), "age", 30)
_ = store.AttributeSetFloat(ctx, person.ID(), "price", 99.99)

// Reading attributes
name, exists, _ := store.AttributeGetString(ctx, person.ID(), "name")
age, exists, _ := store.AttributeGetInt(ctx, person.ID(), "age")
```

## Proposed Solution

### API Design

The proposed API allows direct attribute operations on entity instances:

```go
// Proposed pattern - more intuitive, object-oriented
person, _ := store.EntityCreateWithType(ctx, "person")
_ = person.SetString(ctx, "name", "John")
_ = person.SetInt(ctx, "age", 30)
_ = person.SetFloat(ctx, "price", 99.99)

// Reading attributes
name, exists, _ := person.GetString(ctx, "name")
age, exists, _ := person.GetInt(ctx, "age")
```

### Interface Changes

Add the following methods to `EntityInterface`:

```go
type EntityInterface interface {
    // ... existing methods ...
    
    // Attribute setters
    SetString(ctx context.Context, key string, value string) error
    SetInt(ctx context.Context, key string, value int64) error
    SetFloat(ctx context.Context, key string, value float64) error
    
    // Attribute getters
    GetString(ctx context.Context, key string) (value string, exists bool, err error)
    GetInt(ctx context.Context, key string) (value int64, exists bool, err error)
    GetFloat(ctx context.Context, key string) (value float64, exists bool, err error)
    
    // Bulk operations
    SetAttributes(ctx context.Context, attributes map[string]string) error
    GetAttributes(ctx context.Context) ([]AttributeInterface, error)
    
    // Entity lifecycle (delegated to store)
    Save(ctx context.Context) error          // Upsert entity
    Delete(ctx context.Context) (bool, error) // Permanent delete
    Trash(ctx context.Context) (bool, error)  // Soft delete
    
    // Internal method (not exported, used by store)
    setStore(store StoreInterface)
}
```

### Implementation Details

#### 1. Modify `entityImplementation` to Store Store Reference

```go
// entityImplementation represents a schemaless entity backed by a map[string]string
type entityImplementation struct {
    dataobject.DataObject
    st StoreInterface  // Store reference for object-like methods
}
```

#### 2. Add Internal Store Setter Method

```go
// setStore sets the store reference (internal use only)
func (o *entityImplementation) setStore(store StoreInterface) {
    o.st = store
}
```

#### 3. Implement Object-Like Methods

```go
// SetFloat sets an attribute with float64 value
func (o *entityImplementation) SetFloat(ctx context.Context, key string, value float64) error {
    if o.st == nil {
        return errors.New("entity not bound to store")
    }
    return o.st.AttributeSetFloat(ctx, o.ID(), key, value)
}

// SetInt sets an attribute with int64 value
func (o *entityImplementation) SetInt(ctx context.Context, key string, value int64) error {
    if o.st == nil {
        return errors.New("entity not bound to store")
    }
    return o.st.AttributeSetInt(ctx, o.ID(), key, value)
}

// SetString sets an attribute with string value
func (o *entityImplementation) SetString(ctx context.Context, key string, value string) error {
    if o.st == nil {
        return errors.New("entity not bound to store")
    }
    return o.st.AttributeSetString(ctx, o.ID(), key, value)
}

// GetFloat retrieves a float64 attribute value
func (o *entityImplementation) GetFloat(ctx context.Context, key string) (float64, bool, error) {
    if o.st == nil {
        return 0, false, errors.New("entity not bound to store")
    }
    return o.st.AttributeGetFloat(ctx, o.ID(), key)
}

// GetInt retrieves an int64 attribute value
func (o *entityImplementation) GetInt(ctx context.Context, key string) (int64, bool, error) {
    if o.st == nil {
        return 0, false, errors.New("entity not bound to store")
    }
    return o.st.AttributeGetInt(ctx, o.ID(), key)
}

// GetString retrieves a string attribute value
func (o *entityImplementation) GetString(ctx context.Context, key string) (string, bool, error) {
    if o.st == nil {
        return "", false, errors.New("entity not bound to store")
    }
    return o.st.AttributeGetString(ctx, o.ID(), key)
}

// SetAttributes sets multiple attributes at once
func (o *entityImplementation) SetAttributes(ctx context.Context, attributes map[string]string) error {
    if o.st == nil {
        return errors.New("entity not bound to store")
    }
    return o.st.AttributesSet(ctx, o.ID(), attributes)
}

// GetAttributes retrieves all attributes for this entity
func (o *entityImplementation) GetAttributes(ctx context.Context) ([]AttributeInterface, error) {
    if o.st == nil {
        return nil, errors.New("entity not bound to store")
    }
    return o.st.EntityAttributeList(ctx, o.ID())
}

// Save persists the entity to the store (create or update)
func (o *entityImplementation) Save(ctx context.Context) error {
    if o.st == nil {
        return errors.New("entity not bound to store")
    }
    // Check if entity exists
    existing, err := o.st.EntityFindByID(ctx, o.ID())
    if err != nil {
        return err
    }
    if existing == nil {
        return o.st.EntityCreate(ctx, o)
    }
    return o.st.EntityUpdate(ctx, o)
}

// Delete permanently removes the entity
func (o *entityImplementation) Delete(ctx context.Context) (bool, error) {
    if o.st == nil {
        return false, errors.New("entity not bound to store")
    }
    return o.st.EntityDelete(ctx, o.ID())
}

// Trash soft-deletes the entity
func (o *entityImplementation) Trash(ctx context.Context) (bool, error) {
    if o.st == nil {
        return false, errors.New("entity not bound to store")
    }
    return o.st.EntityTrash(ctx, o.ID())
}
```

#### 4. Modify Store Methods to Inject Store Reference

Update all store methods that return entities to inject the store reference:

```go
// In store_entities.go - EntityFindByID
func (st *storeImplementation) EntityFindByID(ctx context.Context, entityID string) (EntityInterface, error) {
    // ... existing code ...
    if len(list) > 0 {
        entity := list[0]
        // Inject store reference for object-like methods
        if impl, ok := entity.(*entityImplementation); ok {
            impl.setStore(st)
        }
        return entity, nil
    }
    return nil, nil
}

// In store_entities.go - EntityList
func (st *storeImplementation) EntityList(ctx context.Context, options EntityQueryOptions) ([]EntityInterface, error) {
    // ... existing query code ...
    var list []EntityInterface
    for _, m := range entityMaps {
        entity := NewEntityFromExistingData(m)
        // Inject store reference for object-like methods
        if impl, ok := entity.(*entityImplementation); ok {
            impl.setStore(st)
        }
        list = append(list, entity)
    }
    return list, nil
}

// Similar updates for:
// - EntityCreateWithType
// - EntityCreateWithTypeAndAttributes
// - EntityFindByHandle
// - EntityFindByAttribute
// - EntityListByAttribute
```

#### 5. New Entity Constructor with Store Binding

Add a helper for creating entities directly bound to a store:

```go
// NewEntityWithStore creates a new entity bound to a specific store
func NewEntityWithStore(store StoreInterface) EntityInterface {
    o := &entityImplementation{}
    o.SetType("")
    o.SetHandle("")
    o.SetID(GenerateShortID())
    o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
    o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
    o.setStore(store)
    return o
}
```

## Files to Modify

1. **entity_implementation.go**
   - Add `st StoreInterface` field
   - Add `setStore()` internal method
   - Implement all object-like methods (SetFloat, SetInt, SetString, GetFloat, GetInt, GetString, etc.)

2. **interfaces.go**
   - Add new methods to `EntityInterface`
   - Add unexported `setStore(store StoreInterface)` method

3. **store_entities.go**
   - Update all methods that return entities to inject store reference
   - Update `NewEntityFromExistingData` calls

4. **entity_implementation_test.go**
   - Add tests for new object-like methods
   - Test store binding behavior
   - Test error cases (unbound entity)

## Backward Compatibility

This change is fully backward compatible:

- All existing store methods continue to work unchanged
- Entities returned from existing store methods gain new capabilities
- Entities created via `NewEntity()` will not have store binding (object methods return error)
- No changes to database schema

## Usage Examples

### Basic Attribute Operations

```go
ctx := context.Background()

// Create entity through store
product, _ := store.EntityCreateWithType(ctx, "product")

// Set attributes directly on entity
_ = product.SetString(ctx, "name", "Laptop")
_ = product.SetFloat(ctx, "price", 999.99)
_ = product.SetInt(ctx, "stock", 42)

// Get attributes directly from entity
name, exists, _ := product.GetString(ctx, "name")
price, exists, _ := product.GetFloat(ctx, "price")
stock, exists, _ := product.GetInt(ctx, "stock")
```

### Fluent Entity Creation

```go
// Create and configure entity in one chain
product := entitystore.NewEntityWithStore(store).
    SetType("product").
    SetHandle("laptop-pro-2024")

// Set attributes
_ = product.SetString(ctx, "name", "Laptop Pro")
_ = product.SetFloat(ctx, "price", 1499.99)

// Persist to database
_ = product.Save(ctx)
```

### Working with Retrieved Entities

```go
// Find existing entity - automatically bound to store
product, _ := store.EntityFindByID(ctx, "abc123xyz")

// Can use object methods immediately
_ = product.SetFloat(ctx, "sale_price", 1299.99)
```

### List Processing

```go
// All entities in list are store-bound
products, _ := store.EntityList(ctx, entitystore.EntityQueryOptions{
    EntityType: "product",
})

for _, product := range products {
    // Each product can use object methods
    price, _, _ := product.GetFloat(ctx, "price")
    _ = product.SetFloat(ctx, "price", price * 0.9) // 10% discount
}
```

## Benefits

1. **More Intuitive API**: `entity.SetString()` is more natural than `store.AttributeSetString(ctx, entity.ID(), key, value)`
2. **Reduced Boilerplate**: No need to pass `ctx`, `store`, and `entity.ID()` for every operation
3. **Encapsulation**: Entity carries its own context (store binding)
4. **Discoverability**: IDE autocomplete shows available entity operations
5. **Chain-Friendly**: Fluent interface enables method chaining
6. **Backward Compatible**: Existing code continues to work unchanged

## Considerations

### Context Handling

All object-like methods accept `context.Context` as the first parameter. This maintains:
- Proper request-scoped context handling
- Timeout and cancellation support
- Consistency with store methods

### Unbound Entities

Entities created via `NewEntity()` (without store) will have `st == nil`. Object methods return descriptive errors:

```go
entity := entitystore.NewEntity() // No store binding
err := entity.SetString(ctx, "name", "test")
// err: "entity not bound to store"
```

### Thread Safety

The store reference is set once during entity retrieval and never modified. The store implementation itself should handle concurrent access (which it already does).

## Migration Path

No migration required. This is a pure API addition. Existing code patterns continue to work.

Developers can gradually adopt the new pattern:

```go
// Old pattern (still works)
_ = store.AttributeSetString(ctx, person.ID(), "name", "John")

// New pattern (equivalent)
_ = person.SetString(ctx, "name", "John")
```

## Testing Strategy

1. **Unit tests** for each new entity method
2. **Integration tests** verifying store binding works correctly
3. **Error case tests** for unbound entities
4. **Backward compatibility tests** ensuring old patterns still work

## Conclusion

This proposal adds object-like methods to entities by storing a reference to the store in the entity implementation. This enables a more intuitive, object-oriented API while maintaining full backward compatibility and proper context handling.

The implementation is straightforward and follows the existing patterns in the codebase. All changes are additive and do not break existing functionality.
