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

// entityImplementation represents a schemaless entity
type entityImplementation struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

var _ EntityInterface = (*entityImplementation)(nil)

// == CONSTRUCTORS ==========================================================

// NewEntity creates a new entity with default values
func NewEntity() EntityInterface {
	o := &entityImplementation{}
	o.SetEntityType("")
	o.SetEntityHandle("")
	o.SetID(GenerateShortID())
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	return o
}

// NewEntityFromExistingData creates a new entity from existing data
func NewEntityFromExistingData(data map[string]string) EntityInterface {
	o := &entityImplementation{}
	o.Hydrate(data)
	return o
}

// == SETTERS AND GETTERS =====================================================

func (o *entityImplementation) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *entityImplementation) SetID(id string) EntityInterface {
	o.Set(COLUMN_ID, id)
	return o
}

func (o *entityImplementation) EntityType() string {
	return o.Get(COLUMN_ENTITY_TYPE)
}

func (o *entityImplementation) SetEntityType(entityType string) EntityInterface {
	o.Set(COLUMN_ENTITY_TYPE, entityType)
	return o
}

func (o *entityImplementation) EntityHandle() string {
	return o.Get(COLUMN_ENTITY_HANDLE)
}

func (o *entityImplementation) SetEntityHandle(handle string) EntityInterface {
	o.Set(COLUMN_ENTITY_HANDLE, handle)
	return o
}

func (o *entityImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *entityImplementation) SetCreatedAt(createdAt string) EntityInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *entityImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

func (o *entityImplementation) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *entityImplementation) SetUpdatedAt(updatedAt string) EntityInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *entityImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt())
}

// == DYNAMIC ATTRIBUTES ======================================================

// GetAttribute retrieves an attribute by key
func (o *entityImplementation) GetAttribute(key string) string {
	return o.Get(key)
}

// SetAttribute sets an attribute value
func (o *entityImplementation) SetAttribute(key string, value string) EntityInterface {
	o.Set(key, value)
	return o
}

// GetAllAttributes returns all dynamic attributes (excludes system columns)
func (o *entityImplementation) GetAllAttributes() map[string]string {
	systemColumns := map[string]bool{
		COLUMN_ID:            true,
		COLUMN_ENTITY_TYPE:   true,
		COLUMN_ENTITY_HANDLE: true,
		COLUMN_CREATED_AT:    true,
		COLUMN_UPDATED_AT:    true,
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

// attributeImplementation represents a single attribute of an entity
type attributeImplementation struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

var _ AttributeInterface = (*attributeImplementation)(nil)

// == CONSTRUCTORS ==========================================================

// NewAttribute creates a new attribute with default values
func NewAttribute() AttributeInterface {
	o := &attributeImplementation{}
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
	o := &attributeImplementation{}
	o.Hydrate(data)
	return o
}

// == SETTERS AND GETTERS =====================================================

func (o *attributeImplementation) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *attributeImplementation) SetID(id string) AttributeInterface {
	o.Set(COLUMN_ID, id)
	return o
}

func (o *attributeImplementation) EntityID() string {
	return o.Get(COLUMN_ENTITY_ID)
}

func (o *attributeImplementation) SetEntityID(entityID string) AttributeInterface {
	o.Set(COLUMN_ENTITY_ID, entityID)
	return o
}

func (o *attributeImplementation) AttributeKey() string {
	return o.Get(COLUMN_ATTRIBUTE_KEY)
}

func (o *attributeImplementation) SetAttributeKey(key string) AttributeInterface {
	o.Set(COLUMN_ATTRIBUTE_KEY, key)
	return o
}

func (o *attributeImplementation) AttributeValue() string {
	return o.Get(COLUMN_ATTRIBUTE_VALUE)
}

func (o *attributeImplementation) SetAttributeValue(value string) AttributeInterface {
	o.Set(COLUMN_ATTRIBUTE_VALUE, value)
	return o
}

func (o *attributeImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *attributeImplementation) SetCreatedAt(createdAt string) AttributeInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *attributeImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

func (o *attributeImplementation) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *attributeImplementation) SetUpdatedAt(updatedAt string) AttributeInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *attributeImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt())
}

// == TYPE CONVERSIONS ======================================================

// GetInt returns the value as int64
func (o *attributeImplementation) GetInt() (int64, error) {
	return strconv.ParseInt(o.AttributeValue(), 10, 64)
}

// GetFloat returns the value as float64
func (o *attributeImplementation) GetFloat() (float64, error) {
	return strconv.ParseFloat(o.AttributeValue(), 64)
}

// SetInt sets an int64 value
func (o *attributeImplementation) SetInt(value int64) AttributeInterface {
	o.SetAttributeValue(strconv.FormatInt(value, 10))
	return o
}

// SetFloat sets a float64 value
func (o *attributeImplementation) SetFloat(value float64) AttributeInterface {
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
	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon
	
	// Setters
	SetID(id string) EntityInterface
	SetEntityType(entityType string) EntityInterface
	SetEntityHandle(handle string) EntityInterface
	SetCreatedAt(createdAt string) EntityInterface
	SetUpdatedAt(updatedAt string) EntityInterface
	
	// Dynamic attributes
	GetAttribute(key string) string
	SetAttribute(key string, value string) EntityInterface
	GetAllAttributes() map[string]string
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

Following cmsstore pattern exactly - **8 files per entity**:

### 4.1 Entity Files (8 files x 4 entities = 32 files)

| File | Purpose | Lines (Est) |
|------|---------|-------------|
| `entity_implementation.go` | Struct with dataobject, getters, setters | 130 |
| `entity_implementation_test.go` | Tests for entity implementation | 100 |
| `entity_query.go` | Query builder (EntityFindByID, EntityList, etc.) | 150 |
| `entity_query_interface.go` | Query interface definitions | 60 |
| `entity_query_test.go` | Query tests | 100 |
| `entity_table_create_sql.go` | SQL schema for entities table | 40 |
| `store_entities.go` | Store CRUD methods (EntityCreate, EntityUpdate, etc.) | 200 |
| `store_entities_test.go` | Store method tests | 150 |
| **Entity Subtotal** | | **930** |

### 4.2 Attribute Files (8 files)

| File | Purpose | Lines (Est) |
|------|---------|-------------|
| `attribute_implementation.go` | Struct with dataobject, getters, setters | 110 |
| `attribute_implementation_test.go` | Tests for attribute implementation | 80 |
| `attribute_query.go` | Query builder (AttributeFind, AttributeList, etc.) | 120 |
| `attribute_query_interface.go` | Query interface definitions | 50 |
| `attribute_query_test.go` | Query tests | 80 |
| `attribute_table_create_sql.go` | SQL schema for attributes table | 40 |
| `store_attributes.go` | Store CRUD methods (AttributeCreate, AttributeUpdate, etc.) | 150 |
| `store_attributes_test.go` | Store method tests | 120 |
| **Attribute Subtotal** | | **750** |

### 4.3 EntityTrash Files (8 files)

| File | Purpose | Lines (Est) |
|------|---------|-------------|
| `entity_trash_implementation.go` | Struct with dataobject, getters, setters | 130 |
| `entity_trash_implementation_test.go` | Tests for entity trash implementation | 80 |
| `entity_trash_query.go` | Query builder (EntityTrashFindByID, etc.) | 100 |
| `entity_trash_query_interface.go` | Query interface definitions | 40 |
| `entity_trash_query_test.go` | Query tests | 60 |
| `entity_trash_table_create_sql.go` | SQL schema for entities_trash table | 40 |
| `store_entities_trash.go` | Store CRUD methods (EntityTrash, EntityRestore, etc.) | 120 |
| `store_entities_trash_test.go` | Store method tests | 80 |
| **EntityTrash Subtotal** | | **650** |

### 4.4 AttributeTrash Files (8 files)

| File | Purpose | Lines (Est) |
|------|---------|-------------|
| `attribute_trash_implementation.go` | Struct with dataobject, getters, setters | 110 |
| `attribute_trash_implementation_test.go` | Tests for attribute trash implementation | 70 |
| `attribute_trash_query.go` | Query builder (AttributeTrashFindByID, etc.) | 90 |
| `attribute_trash_query_interface.go` | Query interface definitions | 35 |
| `attribute_trash_query_test.go` | Query tests | 50 |
| `attribute_trash_table_create_sql.go` | SQL schema for attributes_trash table | 40 |
| `store_attributes_trash.go` | Store CRUD methods (AttributeTrash, AttributeRestore, etc.) | 100 |
| `store_attributes_trash_test.go` | Store method tests | 70 |
| **AttributeTrash Subtotal** | | **565** |

### 4.5 Support Files (5 files)

| File | Purpose | Lines (Est) |
|------|---------|-------------|
| `interfaces.go` | EntityInterface, AttributeInterface, EntityTrashInterface, AttributeTrashInterface | 100 |
| `consts.go` | Column constants for all 4 entities | 30 |
| `id_helpers.go` | GenerateShortID(), NormalizeID(), IsShortID() | 60 |
| `go.mod` | Add dataobject dependency | 1 |
| `store_implementation.go` | Update AutoMigrate to call new SQL files | 30 |
| **Support Subtotal** | | **221** |

### 4.6 Total Effort

| Category | Files | Lines (Est) |
|----------|-------|-------------|
| **New/Modified** | 37 | 3,116 |
| **Removed** | 16 | -800 |
| **Net Total** | **37** | **~2,316** |

**Implementation: ~14-16 days** (4 entities x 8 files each + support).

### 4.7 Files to Remove (consolidated into new structure)

| File | Reason |
|------|--------|
| `entity.go` | Replaced by `entity_implementation.go` |
| `new_entity.go` | Merged into `entity_implementation.go` |
| `entity_create.go` | Moved to `store_entities.go` |
| `entity_list.go` | Moved to `entity_query.go` |
| `entity_find_by_id.go` | Moved to `entity_query.go` |
| `entity_update.go` | Moved to `store_entities.go` |
| `entity_delete.go` | Moved to `store_entities.go` |
| `entity_trash.go` | Replaced by `entity_trash_implementation.go` |
| `entity_trash_create.go` | Moved to `store_entities_trash.go` |
| `entity_trash_list.go` | Moved to `entity_trash_query.go` |
| `entity_trash_find.go` | Moved to `entity_trash_query.go` |
| `entity_trash_restore.go` | Moved to `store_entities_trash.go` |
| `attribute_trash.go` | Replaced by `attribute_trash_implementation.go` |
| `attribute_trash_create.go` | Moved to `store_attributes_trash.go` |
| `attribute_trash_list.go` | Moved to `attribute_trash_query.go` |
| `attribute_trash_find.go` | Moved to `attribute_trash_query.go` |
| `attribute_trash_restore.go` | Moved to `store_attributes_trash.go` |
| `attribute.go` | Replaced by `attribute_implementation.go` |
| `new_attribute.go` | Merged into `attribute_implementation.go` |
| `attribute_create.go` | Moved to `store_attributes.go` |
| `attribute_list.go` | Moved to `attribute_query.go` |
| `attribute_find.go` | Moved to `attribute_query.go` |
| `attribute_update.go` | Moved to `store_attributes.go` |
| `attribute_delete.go` | Moved to `store_attributes.go` |
| `attribute_trash.go` | Moved to `store_attributes.go` |

### 4.5 Total Effort

| Category | Files | Lines (Est) |
|----------|-------|-------------|
| **New/Modified** | 21 | 1,941 |
| **Removed** | 16 | -800 |
| **Net Total** | **21** | **~1,141** |

**Implementation: ~10-12 days** (was ~8 days with old estimate)

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
3. Update `consts.go` with column constants for all 4 entities
4. Define interfaces: `EntityInterface`, `AttributeInterface`, `EntityTrashInterface`, `AttributeTrashInterface`

### Phase 2: Rewrite Entity (2 days)

1. Rewrite `entity_implementation.go` using `dataobject.DataObject`
2. Add getters/setters with `o.Get()` / `o.Set()` pattern
3. Remove `new_entity.go` (merge into `entity_implementation.go`)
4. Write/update tests

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
| **EntityTrash type** | `EntityTrash` struct | `EntityTrashInterface` |
| **AttributeTrash type** | `AttributeTrash` struct | `AttributeTrashInterface` |

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

// Check status (trash tables handle deletion, not soft delete)
// Use EntityTrash() to move to trash, EntityRestore() to restore
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
