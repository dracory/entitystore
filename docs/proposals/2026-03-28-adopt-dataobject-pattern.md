# Adopt dataobject Pattern for Entitystore Proposal

**Date:** 2026-03-28
**Status:** Draft
**Author:** AI Assistant
**Repository:** github.com/dracory/entitystore

---

## 1. Executive Summary

**Problem:** `entitystore` uses a different entity implementation pattern than `cmsstore`:

| Aspect | entitystore (current) | cmsstore |
|--------|----------------------|----------|
| **Base struct** | Manual struct with private fields | `dataobject.DataObject` |
| **Data storage** | Explicit struct fields | `map[string]string` |
| **Getters/Setters** | Direct field access | `o.Get()` / `o.Set()` |
| **Constructor** | Minimal defaults | Rich defaults with status, timestamps |
| **Hydration** | `NewEntityFromMap()` | `o.Hydrate(data)` |

**Solution:** Migrate `entitystore` to use `dataobject.DataObject` pattern matching `cmsstore`.

**Impact:**
- Consistency between `entitystore` and `cmsstore`
- Simpler entity implementation (no explicit struct fields)
- Easier to add new attributes without code changes
- Built-in serialization/deserialization
- **Breaking change requiring migration**

---

## 2. Current State (As-Is)

### 2.1 Entitystore Entity (Current)

**File:** `entity.go`

```go
package entitystore

import "time"

// Entity this is the type for an Entity
type Entity struct {
	id           string
	entityType   string
	entityHandle string
	createdAt    time.Time
	updatedAt    time.Time
	st           *storeImplementation
}

// Explicit getter
func (e *Entity) ID() string {
	return e.id
}

// Fluent setter
func (e *Entity) SetID(id string) *Entity {
	e.id = id
	return e
}

// ... more explicit getters/setters for each field
```

### 2.2 Attribute (Current)

**File:** `attribute.go`

```go
package entitystore

import "time"

// Attribute type
type Attribute struct {
	id             string
	entityID       string
	attributeKey   string
	attributeValue string
	createdAt      time.Time
	updatedAt      time.Time
	st             *storeImplementation
}

// ... explicit getters/setters for each field
```

### 2.3 Problems with Current Pattern

1. **Verbose** - Every field needs explicit getter/setter
2. **Rigid** - Adding a field requires code changes
3. **Inconsistent** - Different from `cmsstore` pattern
4. **Boilerplate** - Lots of repetitive code

---

## 3. Proposed Solution (To-Be)

### 3.1 New Entity Pattern

**File:** `entity.go` (revised)

```go
package entitystore

import (
	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

// == TYPE ===================================================================

// entity represents a schemaless entity
type entity struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

var _ EntityInterface = (*entity)(nil)

// == CONSTRUCTORS ==========================================================

// NewEntity creates a new entity with default values
func NewEntity() EntityInterface {
	o := &entity{}
	o.SetEntityType("")
	o.SetEntityHandle("")
	o.SetID(GenerateShortID())
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetSoftDeletedAt(sb.MAX_DATETIME)
	return o
}

// NewEntityFromExistingData creates a new entity from existing data
func NewEntityFromExistingData(data map[string]string) EntityInterface {
	o := &entity{}
	o.Hydrate(data)
	return o
}

// == METHODS ===============================================================

func (o *entity) IsActive() bool {
	return o.Status() == ENTITY_STATUS_ACTIVE
}

func (o *entity) IsInactive() bool {
	return o.Status() == ENTITY_STATUS_INACTIVE
}

func (o *entity) IsSoftDeleted() bool {
	return o.SoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// == SETTERS AND GETTERS =====================================================

func (o *entity) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *entity) SetID(id string) EntityInterface {
	o.Set(COLUMN_ID, id)
	return o
}

func (o *entity) EntityType() string {
	return o.Get(COLUMN_ENTITY_TYPE)
}

func (o *entity) SetEntityType(entityType string) EntityInterface {
	o.Set(COLUMN_ENTITY_TYPE, entityType)
	return o
}

func (o *entity) EntityHandle() string {
	return o.Get(COLUMN_ENTITY_HANDLE)
}

func (o *entity) SetEntityHandle(handle string) EntityInterface {
	o.Set(COLUMN_ENTITY_HANDLE, handle)
	return o
}

func (o *entity) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *entity) SetStatus(status string) EntityInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *entity) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *entity) SetCreatedAt(createdAt string) EntityInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *entity) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

func (o *entity) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *entity) SetUpdatedAt(updatedAt string) EntityInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *entity) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt())
}

func (o *entity) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

func (o *entity) SetSoftDeletedAt(softDeletedAt string) EntityInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

func (o *entity) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt())
}

// == DYNAMIC ATTRIBUTES ======================================================

// GetAttribute retrieves an attribute by key
func (o *entity) GetAttribute(key string) string {
	return o.Get(key)
}

// SetAttribute sets an attribute value
func (o *entity) SetAttribute(key string, value string) EntityInterface {
	o.Set(key, value)
	return o
}

// GetAllAttributes returns all dynamic attributes (excludes system columns)
func (o *entity) GetAllAttributes() map[string]string {
	systemColumns := map[string]bool{
		COLUMN_ID:              true,
		COLUMN_ENTITY_TYPE:     true,
		COLUMN_ENTITY_HANDLE:   true,
		COLUMN_STATUS:          true,
		COLUMN_CREATED_AT:      true,
		COLUMN_UPDATED_AT:      true,
		COLUMN_SOFT_DELETED_AT: true,
	}
	
	attrs := make(map[string]string)
	for k, v := range o.Data() {
		if !systemColumns[k] {
			attrs[k] = v
		}
	}
	return attrs
}
```

### 3.2 New Attribute Pattern

**File:** `attribute.go` (revised)

```go
package entitystore

import (
	"github.com/dracory/dataobject"
	"github.com/dromara/carbon/v2"
)

// == TYPE ===================================================================

// attribute represents a single attribute of an entity
type attribute struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

var _ AttributeInterface = (*attribute)(nil)

// == CONSTRUCTORS ==========================================================

// NewAttribute creates a new attribute with default values
func NewAttribute() AttributeInterface {
	o := &attribute{}
	o.SetEntityID("")
	o.SetAttributeKey("")
	o.SetAttributeValue("")
	o.SetID(GenerateShortID())
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	return o
}

// NewAttributeFromExistingData creates a new attribute from existing data
func NewAttributeFromExistingData(data map[string]string) AttributeInterface {
	o := &attribute{}
	o.Hydrate(data)
	return o
}

// == SETTERS AND GETTERS =====================================================

func (o *attribute) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *attribute) SetID(id string) AttributeInterface {
	o.Set(COLUMN_ID, id)
	return o
}

func (o *attribute) EntityID() string {
	return o.Get(COLUMN_ENTITY_ID)
}

func (o *attribute) SetEntityID(entityID string) AttributeInterface {
	o.Set(COLUMN_ENTITY_ID, entityID)
	return o
}

func (o *attribute) AttributeKey() string {
	return o.Get(COLUMN_ATTRIBUTE_KEY)
}

func (o *attribute) SetAttributeKey(key string) AttributeInterface {
	o.Set(COLUMN_ATTRIBUTE_KEY, key)
	return o
}

func (o *attribute) AttributeValue() string {
	return o.Get(COLUMN_ATTRIBUTE_VALUE)
}

func (o *attribute) SetAttributeValue(value string) AttributeInterface {
	o.Set(COLUMN_ATTRIBUTE_VALUE, value)
	return o
}

func (o *attribute) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *attribute) SetCreatedAt(createdAt string) AttributeInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *attribute) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

func (o *attribute) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *attribute) SetUpdatedAt(updatedAt string) AttributeInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *attribute) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt())
}

// == TYPE CONVERSIONS ======================================================

// GetInt returns the value as int64
func (o *attribute) GetInt() (int64, error) {
	return strconv.ParseInt(o.AttributeValue(), 10, 64)
}

// GetFloat returns the value as float64
func (o *attribute) GetFloat() (float64, error) {
	return strconv.ParseFloat(o.AttributeValue(), 64)
}

// SetInt sets an int64 value
func (o *attribute) SetInt(value int64) AttributeInterface {
	o.SetAttributeValue(strconv.FormatInt(value, 10))
	return o
}

// SetFloat sets a float64 value
func (o *attribute) SetFloat(value float64) AttributeInterface {
	o.SetAttributeValue(strconv.FormatFloat(value, 'f', 30, 64))
	return o
}
```

### 3.3 New Interfaces

**File:** `interfaces.go` (append)

```go
// EntityInterface defines the contract for entities
type EntityInterface interface {
	dataobject.DataObjectInterface
	
	// Getters
	ID() string
	EntityType() string
	EntityHandle() string
	Status() string
	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon
	SoftDeletedAt() string
	SoftDeletedAtCarbon() *carbon.Carbon
	
	// Setters
	SetID(id string) EntityInterface
	SetEntityType(entityType string) EntityInterface
	SetEntityHandle(handle string) EntityInterface
	SetStatus(status string) EntityInterface
	SetCreatedAt(createdAt string) EntityInterface
	SetUpdatedAt(updatedAt string) EntityInterface
	SetSoftDeletedAt(softDeletedAt string) EntityInterface
	
	// Dynamic attributes
	GetAttribute(key string) string
	SetAttribute(key string, value string) EntityInterface
	GetAllAttributes() map[string]string
	
	// Status checks
	IsActive() bool
	IsInactive() bool
	IsSoftDeleted() bool
}

// AttributeInterface defines the contract for attributes
type AttributeInterface interface {
	dataobject.DataObjectInterface
	
	// Getters
	ID() string
	EntityID() string
	AttributeKey() string
	AttributeValue() string
	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon
	
	// Setters
	SetID(id string) AttributeInterface
	SetEntityID(entityID string) AttributeInterface
	SetAttributeKey(key string) AttributeInterface
	SetAttributeValue(value string) AttributeInterface
	SetCreatedAt(createdAt string) AttributeInterface
	SetUpdatedAt(updatedAt string) AttributeInterface
	
	// Type conversions
	GetInt() (int64, error)
	GetFloat() (float64, error)
	SetInt(value int64) AttributeInterface
	SetFloat(value float64) AttributeInterface
}
```

### 3.4 Comparison: Before vs After

| Aspect | Before (Current) | After (dataobject) |
|--------|-----------------|-------------------|
| **Entity struct** | `type Entity struct { id string; entityType string; ... }` | `type entity struct { dataobject.DataObject }` |
| **Get ID** | `return e.id` | `return o.Get(COLUMN_ID)` |
| **Set ID** | `e.id = id; return e` | `o.Set(COLUMN_ID, id); return o` |
| **Add field** | Add to struct + getter + setter | Just add getter + setter |
| **Hydrate** | Manual mapping in `NewEntityFromMap` | `o.Hydrate(data)` |
| **ToMap** | Manual `map[string]any` construction | `o.Data()` |
| **Consistency** | Different from cmsstore | Matches cmsstore exactly |

---

## 4. Files to Modify

| File | Changes | Lines (Est) |
|------|---------|-------------|
| `entity.go` | Rewrite using dataobject | 120 |
| `attribute.go` | Rewrite using dataobject | 100 |
| `new_entity.go` | Remove (merged into entity.go) | -20 |
| `new_attribute.go` | Remove (merged into attribute.go) | -20 |
| `entity_attribute_list.go` | Update to use interface | 20 |
| `interfaces.go` | Add EntityInterface, AttributeInterface | 60 |
| `consts.go` | Add status constants | 10 |
| `go.mod` | Add dataobject dependency | 1 |
| `id_helpers.go` | Add GenerateShortID() | 50 |
| `store_implementation.go` | Update to return interfaces | 30 |
| `entity_create.go` | Update for new pattern | 20 |
| `attribute_create.go` | Update for new pattern | 20 |
| **Total** | | **~391** |

---

## 5. Dependencies to Add

**File:** `go.mod`

```
require github.com/dracory/dataobject v0.x.x
```

This is the same dependency used by `cmsstore`.

---

## 6. Implementation Plan

### Phase 1: Setup (1 day)

1. Add `dataobject` dependency to `go.mod`
2. Create `id_helpers.go` with `GenerateShortID()`
3. Update `consts.go` with status constants
4. Define `EntityInterface` and `AttributeInterface`

### Phase 2: Rewrite Entity (2 days)

1. Rewrite `entity.go` using `dataobject.DataObject`
2. Add all getters/setters with `o.Get()` / `o.Set()` pattern
3. Implement `IsActive()`, `IsInactive()`, `IsSoftDeleted()`
4. Add dynamic attribute methods
5. Remove `new_entity.go` (merge into `entity.go`)
6. Write/update tests

### Phase 3: Rewrite Attribute (2 days)

1. Rewrite `attribute.go` using `dataobject.DataObject`
2. Add all getters/setters with `o.Get()` / `o.Set()` pattern
3. Implement type conversion methods
4. Remove `new_attribute.go` (merge into `attribute.go`)
5. Write/update tests

### Phase 4: Update Store Methods (2 days)

1. Update `store_implementation.go` to return interfaces
2. Update `entity_create.go` for new pattern
3. Update `attribute_create.go` for new pattern
4. Update `entity_attribute_list.go` for new pattern
5. Run full test suite

### Phase 5: Documentation (1 day)

1. Update `README.md` with breaking changes
2. Update usage examples
3. Document migration path
4. Tag v2.0.0

**Total: ~8 days**

---

## 7. Migration Guide

### 7.1 Breaking Changes

| Aspect | v1.x | v2.0.0 |
|--------|------|--------|
| **Entity type** | `Entity` struct | `EntityInterface` |
| **Attribute type** | `Attribute` struct | `AttributeInterface` |
| **ID generation** | `uid.HumanUid()` | `GenerateShortID()` (9-char) |
| **Constructor** | `NewEntity(opts)` | `NewEntity()` |
| **Hydration** | `NewEntityFromMap(data)` | `NewEntityFromExistingData(data)` |
| **Data access** | `entity.id` | `entity.Get(COLUMN_ID)` |

### 7.2 Code Migration Examples

**Before:**
```go
// Create entity
entity := store.NewEntity(NewEntityOptions{
    Type: "product",
})
entity.SetString("name", "iPhone")

// Access field
id := entity.ID()

// Check status
if entity.IsActive() // Not available
```

**After:**
```go
// Create entity
entity := entitystore.NewEntity()
entity.SetEntityType("product")
entity.SetAttribute("name", "iPhone")

// Access field (same API)
id := entity.ID()

// Check status (new)
if entity.IsActive()
```

### 7.3 Database Migration

Same as short ID migration - requires export/import or staying on v1.x:

```sql
-- Option 1: New installations - use v2.0.0 directly
-- Option 2: Existing data - export/import with new IDs
-- Option 3: Stay on v1.x for existing projects
```

---

## 8. Benefits

### 8.1 Code Quality

- **Less boilerplate** - No explicit struct fields
- **Easier to extend** - Add fields without struct changes
- **Consistent** - Matches `cmsstore` exactly
- **DRY** - Reuse `dataobject` functionality

### 8.2 Maintainability

- **Single source** - Changes to dataobject benefit all packages
- **Testable** - Interface-based design
- **Documented** - Follows established cmsstore pattern

### 8.3 Developer Experience

```go
// Adding a new field is now trivial:

// 1. Add constant
const COLUMN_ENTITY_NAME = "entity_name"

// 2. Add getter/setter (no struct changes!)
func (o *entity) EntityName() string {
    return o.Get(COLUMN_ENTITY_NAME)
}

func (o *entity) SetEntityName(name string) EntityInterface {
    o.Set(COLUMN_ENTITY_NAME, name)
    return o
}

// Done! No struct modifications needed
```

---

## 9. Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Breaking API changes | High | Major version bump, migration guide |
| Performance (map vs struct) | Low | Negligible for typical use |
| Learning curve | Medium | Document with examples |
| Dependency addition | Low | dataobject is lightweight |
| Test failures | Medium | Comprehensive test updates |

---

## 10. Conclusion

### Recommendation

**Proceed with implementation** as part of v2.0.0.

This change:
1. **Aligns** entitystore with cmsstore architecture
2. **Simplifies** entity implementation significantly
3. **Enables** faster feature development
4. **Reduces** maintenance burden

### Combined v2.0.0 Changes

This proposal should be implemented alongside:
1. **Short IDs** (separate proposal) - 9-char IDs
2. **dataobject Pattern** (this proposal) - Map-based entities
3. **Relationships** (separate proposal) - Entity linking
4. **Taxonomy** (separate proposal) - Categorization

All four changes together create a cohesive v2.0.0 that modernizes entitystore to match cmsstore standards.

### Next Actions

1. **Review** this proposal alongside other v2.0.0 proposals
2. **Create** v2.0.0 feature branch
3. **Implement** in order: short IDs → dataobject → relationships → taxonomy
4. **Release** v2.0.0 with combined changelog

---

**End of Proposal**
