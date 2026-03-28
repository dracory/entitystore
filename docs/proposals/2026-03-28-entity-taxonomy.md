# Entity Taxonomy Support Proposal

**Date:** 2026-03-28
**Status:** ⏳ PENDING (waiting for dataobject pattern stabilization)
**Author:** AI Assistant
**Repository:** github.com/dracory/entitystore

---

## 1. Executive Summary

**Problem:** `entitystore` provides excellent EAV storage but lacks native support for categorizing/classifying entities (taxonomy).

**Current Workaround:** Store taxonomy IDs as string attributes, query separately.

```go
// Current workaround - manual attribute storage
entity.SetAttribute("category_id", "cat_123")
entity.SetAttribute("tag_ids", `["tag_1", "tag_2"]`)

// Manual filtering in application code
categoryID := entity.GetAttribute("category_id")
// Query all products, filter in code - inefficient
```

**Solution:** Add native taxonomy support to `entitystore`.

**Impact:**
- All projects using `entitystore` get categorization for free
- Enables taxonomy-based queries ("find all products in Electronics")
- Supports hierarchical categories, tags, labels, etc.

**Status:** ⏳ **PENDING** - Waiting for dataobject pattern implementation to stabilize before adding taxonomy

---

## 2. Current State (As-Is) ✅ IMPLEMENTED

### 2.1 Entitystore Architecture (dataobject Pattern)

```
entitystore/
├── entity_implementation.go          # Entity with dataobject.DataObject
├── entity_implementation_test.go     # Entity tests
├── entity_query.go                   # Entity query builder
├── entity_query_interface.go         # EntityQueryInterface
├── entity_table_create_sql.go        # Entities table SQL
├── attribute_implementation.go       # Attribute with dataobject.DataObject
├── attribute_implementation_test.go    # Attribute tests
├── attribute_query.go                # Attribute query builder
├── attribute_query_interface.go      # AttributeQueryInterface
├── attribute_table_create_sql.go     # Attributes table SQL
├── entity_trash_implementation.go  # EntityTrash with dataobject
├── entity_trash_implementation_test.go
├── entity_trash_query_interface.go
├── entity_trash_table_create_sql.go
├── attribute_trash_implementation.go # AttributeTrash with dataobject
├── attribute_trash_implementation_test.go
├── attribute_trash_query_interface.go
├── attribute_trash_table_create_sql.go
├── store_entities.go                 # Entity CRUD methods
├── store_entities_test.go            # Entity store tests
├── store_attributes.go               # Attribute CRUD methods
├── store_attributes_test.go            # Attribute store tests
├── store_entities_trash.go           # Entity trash/restore
├── store_attributes_trash.go         # Attribute trash/restore
├── interfaces.go                     # All entity interfaces
├── consts.go                         # Column constants
├── id_helpers.go                     # GenerateShortID()
└── new.go                            # Store initialization
```

### 2.2 Current Database Schema (Short IDs + dataobject)

```sql
-- Entities table (exists with dataobject pattern)
CREATE TABLE entities (
    id varchar(15) PRIMARY KEY,           -- Short ID (9-15 chars)
    entity_type varchar(40) NOT NULL,
    entity_handle varchar(60),
    created_at datetime NOT NULL,
    updated_at datetime NOT NULL
);

-- Attributes table (exists with dataobject pattern)
CREATE TABLE attributes (
    id varchar(15) PRIMARY KEY,           -- Short ID
    entity_id varchar(15) NOT NULL,       -- References entities.id
    attribute_key varchar(255) NOT NULL,
    attribute_value text,
    created_at datetime NOT NULL,
    updated_at datetime NOT NULL
);

-- Trash tables (separate, no soft delete)
CREATE TABLE entities_trash (
    id varchar(15) PRIMARY KEY,
    entity_type varchar(40) NOT NULL,
    entity_handle varchar(60),
    created_at datetime NOT NULL,
    updated_at datetime NOT NULL,
    deleted_at datetime NOT NULL,
    deleted_by varchar(15)
);

CREATE TABLE attributes_trash (
    id varchar(15) PRIMARY KEY,
    entity_id varchar(15) NOT NULL,
    attribute_key varchar(255) NOT NULL,
    attribute_value text,
    created_at datetime NOT NULL,
    updated_at datetime NOT NULL,
    deleted_at datetime NOT NULL,
    deleted_by varchar(15)
);
```

### 2.3 Current Limitations

**Problems with attribute-based taxonomy:**
- No referential integrity
- No hierarchical support (parent/child categories)
- Complex queries require parsing JSON attributes
- No taxonomy metadata (description, slug, etc.)
- No validation of taxonomy assignments

---

## 3. Proposed Solution (To-Be) ⏳ PENDING

**Note:** This proposal is pending stabilization of the dataobject pattern implementation. Once the core entitystore architecture (entities, attributes, trash tables with dataobject) is stable, taxonomy can be added following the same pattern.

### 3.1 New Types (dataobject Pattern)

Following the same pattern as entities and attributes:

**File:** `taxonomy_implementation.go`

```go
package entitystore

import "github.com/dracory/dataobject"

// Taxonomy is a categorization system (categories, tags, etc.)
type taxonomyImplementation struct {
	*dataobject.DataObject
}

// Column constants for taxonomies
const (
	COLUMN_NAME        = "name"
	COLUMN_SLUG        = "slug"
	COLUMN_DESCRIPTION = "description"
	COLUMN_PARENT_ID   = "parent_id"
	COLUMN_ENTITY_TYPES = "entity_types"
	COLUMN_TAXONOMY_ID = "taxonomy_id"
	COLUMN_SORT_ORDER  = "sort_order"
	COLUMN_TERM_ID     = "term_id"
)

// NewTaxonomy creates a new taxonomy instance
func NewTaxonomy() TaxonomyInterface {
	return &taxonomyImplementation{
		DataObject: dataobject.NewDataObject(),
	}
}

// ID returns the taxonomy ID
func (o *taxonomyImplementation) ID() string {
	return o.Get(COLUMN_ID)
}

// SetID sets the taxonomy ID (fluent)
func (o *taxonomyImplementation) SetID(id string) TaxonomyInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// Name returns the taxonomy name
func (o *taxonomyImplementation) Name() string {
	return o.Get(COLUMN_NAME)
}

// SetName sets the taxonomy name (fluent)
func (o *taxonomyImplementation) SetName(name string) TaxonomyInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

// Slug returns the taxonomy slug
func (o *taxonomyImplementation) Slug() string {
	return o.Get(COLUMN_SLUG)
}

// SetSlug sets the taxonomy slug (fluent)
func (o *taxonomyImplementation) SetSlug(slug string) TaxonomyInterface {
	o.Set(COLUMN_SLUG, slug)
	return o
}

// Description returns the taxonomy description
func (o *taxonomyImplementation) Description() string {
	return o.Get(COLUMN_DESCRIPTION)
}

// SetDescription sets the taxonomy description (fluent)
func (o *taxonomyImplementation) SetDescription(desc string) TaxonomyInterface {
	o.Set(COLUMN_DESCRIPTION, desc)
	return o
}

// ParentID returns the parent taxonomy ID
func (o *taxonomyImplementation) ParentID() string {
	return o.Get(COLUMN_PARENT_ID)
}

// SetParentID sets the parent taxonomy ID (fluent)
func (o *taxonomyImplementation) SetParentID(parentID string) TaxonomyInterface {
	o.Set(COLUMN_PARENT_ID, parentID)
	return o
}

// EntityTypes returns the entity types this taxonomy applies to
func (o *taxonomyImplementation) EntityTypes() []string {
	// Parse from stored JSON or comma-separated
	return o.GetAttributeStrings(COLUMN_ENTITY_TYPES)
}

// SetEntityTypes sets the entity types (fluent)
func (o *taxonomyImplementation) SetEntityTypes(types []string) TaxonomyInterface {
	o.SetAttributeStrings(COLUMN_ENTITY_TYPES, types)
	return o
}

// CreatedAt returns creation timestamp
func (o *taxonomyImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

// SetCreatedAt sets creation timestamp (fluent)
func (o *taxonomyImplementation) SetCreatedAt(createdAt string) TaxonomyInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

// UpdatedAt returns update timestamp
func (o *taxonomyImplementation) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

// SetUpdatedAt sets update timestamp (fluent)
func (o *taxonomyImplementation) SetUpdatedAt(updatedAt string) TaxonomyInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}
```

**File:** `taxonomy_term_implementation.go`

```go
package entitystore

import "github.com/dracory/dataobject"

// taxonomyTermImplementation implements TaxonomyTermInterface
type taxonomyTermImplementation struct {
	*dataobject.DataObject
}

// NewTaxonomyTerm creates a new taxonomy term instance
func NewTaxonomyTerm() TaxonomyTermInterface {
	return &taxonomyTermImplementation{
		DataObject: dataobject.NewDataObject(),
	}
}

// ID returns the term ID
func (o *taxonomyTermImplementation) ID() string {
	return o.Get(COLUMN_ID)
}

// SetID sets the term ID (fluent)
func (o *taxonomyTermImplementation) SetID(id string) TaxonomyTermInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// TaxonomyID returns the parent taxonomy ID
func (o *taxonomyTermImplementation) TaxonomyID() string {
	return o.Get(COLUMN_TAXONOMY_ID)
}

// SetTaxonomyID sets the parent taxonomy ID (fluent)
func (o *taxonomyTermImplementation) SetTaxonomyID(taxonomyID string) TaxonomyTermInterface {
	o.Set(COLUMN_TAXONOMY_ID, taxonomyID)
	return o
}

// Name returns the term name
func (o *taxonomyTermImplementation) Name() string {
	return o.Get(COLUMN_NAME)
}

// SetName sets the term name (fluent)
func (o *taxonomyTermImplementation) SetName(name string) TaxonomyTermInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

// Slug returns the term slug
func (o *taxonomyTermImplementation) Slug() string {
	return o.Get(COLUMN_SLUG)
}

// SetSlug sets the term slug (fluent)
func (o *taxonomyTermImplementation) SetSlug(slug string) TaxonomyTermInterface {
	o.Set(COLUMN_SLUG, slug)
	return o
}

// ParentID returns the parent term ID (for hierarchy)
func (o *taxonomyTermImplementation) ParentID() string {
	return o.Get(COLUMN_PARENT_ID)
}

// SetParentID sets the parent term ID (fluent)
func (o *taxonomyTermImplementation) SetParentID(parentID string) TaxonomyTermInterface {
	o.Set(COLUMN_PARENT_ID, parentID)
	return o
}

// SortOrder returns the sort order
func (o *taxonomyTermImplementation) SortOrder() int {
	return o.GetInt(COLUMN_SORT_ORDER)
}

// SetSortOrder sets the sort order (fluent)
func (o *taxonomyTermImplementation) SetSortOrder(order int) TaxonomyTermInterface {
	o.Set(COLUMN_SORT_ORDER, order)
	return o
}

// CreatedAt returns creation timestamp
func (o *taxonomyTermImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

// SetCreatedAt sets creation timestamp (fluent)
func (o *taxonomyTermImplementation) SetCreatedAt(createdAt string) TaxonomyTermInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

// UpdatedAt returns update timestamp
func (o *taxonomyTermImplementation) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

// SetUpdatedAt sets update timestamp (fluent)
func (o *taxonomyTermImplementation) SetUpdatedAt(updatedAt string) TaxonomyTermInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}
```

**File:** `entity_taxonomy_implementation.go`

```go
package entitystore

import "github.com/dracory/dataobject"

// entityTaxonomyImplementation implements EntityTaxonomyInterface
// This links entities to taxonomy terms (the assignment table)
type entityTaxonomyImplementation struct {
	*dataobject.DataObject
}

// NewEntityTaxonomy creates a new entity-taxonomy link instance
func NewEntityTaxonomy() EntityTaxonomyInterface {
	return &entityTaxonomyImplementation{
		DataObject: dataobject.NewDataObject(),
	}
}

// ID returns the link ID
func (o *entityTaxonomyImplementation) ID() string {
	return o.Get(COLUMN_ID)
}

// SetID sets the link ID (fluent)
func (o *entityTaxonomyImplementation) SetID(id string) EntityTaxonomyInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// EntityID returns the entity ID
func (o *entityTaxonomyImplementation) EntityID() string {
	return o.Get(COLUMN_ENTITY_ID)
}

// SetEntityID sets the entity ID (fluent)
func (o *entityTaxonomyImplementation) SetEntityID(entityID string) EntityTaxonomyInterface {
	o.Set(COLUMN_ENTITY_ID, entityID)
	return o
}

// TaxonomyID returns the taxonomy ID
func (o *entityTaxonomyImplementation) TaxonomyID() string {
	return o.Get(COLUMN_TAXONOMY_ID)
}

// SetTaxonomyID sets the taxonomy ID (fluent)
func (o *entityTaxonomyImplementation) SetTaxonomyID(taxonomyID string) EntityTaxonomyInterface {
	o.Set(COLUMN_TAXONOMY_ID, taxonomyID)
	return o
}

// TermID returns the term ID
func (o *entityTaxonomyImplementation) TermID() string {
	return o.Get(COLUMN_TERM_ID)
}

// SetTermID sets the term ID (fluent)
func (o *entityTaxonomyImplementation) SetTermID(termID string) EntityTaxonomyInterface {
	o.Set(COLUMN_TERM_ID, termID)
	return o
}

// CreatedAt returns creation timestamp
func (o *entityTaxonomyImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

// SetCreatedAt sets creation timestamp (fluent)
func (o *entityTaxonomyImplementation) SetCreatedAt(createdAt string) EntityTaxonomyInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}
```

**Options structs:**

**Store Options (enable/disable feature):**

```go
type NewStoreOptions struct {
    DB                  *sql.DB
    EntityTableName     string
    AttributeTableName  string
    
    // Feature flags - taxonomies are optional
    TaxonomiesEnabled bool   // Enable taxonomy support
    TaxonomyTableName string // Default: "entities_taxonomies"
    TaxonomyTermTableName string // Default: "entities_taxonomy_terms"
    EntityTaxonomyTableName string // Default: "entities_entity_taxonomies"
    
    // Trash tables for soft delete
    TaxonomyTrashTableName string // Default: "entities_taxonomies_trash"
    TaxonomyTermTrashTableName string // Default: "entities_taxonomy_terms_trash"
}
```

type TaxonomyOptions struct {
    Name        string
    Slug        string
    Description string
    ParentID    string
    EntityTypes []string
}

type TaxonomyTermOptions struct {
    TaxonomyID string
    Name       string
    Slug       string
    ParentID   string
    SortOrder  int
}
```

### 3.1b Interface Definitions

**File:** `interfaces.go` (append new interfaces)

```go
// TaxonomyInterface defines the taxonomy (classification system) type
type TaxonomyInterface interface {
	dataobject.DataObjectInterface

	// Core fields
	ID() string
	SetID(id string) TaxonomyInterface
	Name() string
	SetName(name string) TaxonomyInterface
	Slug() string
	SetSlug(slug string) TaxonomyInterface
	Description() string
	SetDescription(desc string) TaxonomyInterface

	// Hierarchy and scope
	ParentID() string
	SetParentID(parentID string) TaxonomyInterface
	EntityTypes() []string
	SetEntityTypes(types []string) TaxonomyInterface

	// Timestamps
	CreatedAt() string
	SetCreatedAt(createdAt string) TaxonomyInterface
	UpdatedAt() string
	SetUpdatedAt(updatedAt string) TaxonomyInterface
}

// TaxonomyTermInterface defines a term within a taxonomy
type TaxonomyTermInterface interface {
	dataobject.DataObjectInterface

	// Core fields
	ID() string
	SetID(id string) TaxonomyTermInterface
	TaxonomyID() string
	SetTaxonomyID(taxonomyID string) TaxonomyTermInterface
	Name() string
	SetName(name string) TaxonomyTermInterface
	Slug() string
	SetSlug(slug string) TaxonomyTermInterface

	// Hierarchy and ordering
	ParentID() string
	SetParentID(parentID string) TaxonomyTermInterface
	SortOrder() int
	SetSortOrder(order int) TaxonomyTermInterface

	// Timestamps
	CreatedAt() string
	SetCreatedAt(createdAt string) TaxonomyTermInterface
	UpdatedAt() string
	SetUpdatedAt(updatedAt string) TaxonomyTermInterface
}

// EntityTaxonomyInterface links entities to taxonomy terms
type EntityTaxonomyInterface interface {
	dataobject.DataObjectInterface

	// Core fields
	ID() string
	SetID(id string) EntityTaxonomyInterface
	EntityID() string
	SetEntityID(entityID string) EntityTaxonomyInterface
	TaxonomyID() string
	SetTaxonomyID(taxonomyID string) EntityTaxonomyInterface
	TermID() string
	SetTermID(termID string) EntityTaxonomyInterface

	// Timestamps
	CreatedAt() string
	SetCreatedAt(createdAt string) EntityTaxonomyInterface
}
```

### 3.2 Store Interface Extensions

**File:** `interfaces.go` (append to StoreInterface)

```go
type StoreInterface interface {
    // ... existing methods ...
    
    // ==========================================
    // Taxonomies
    // ==========================================
    
    // TaxonomyCreate creates a new taxonomy
    TaxonomyCreate(ctx context.Context, opts TaxonomyOptions) (*Taxonomy, error)
    
    // TaxonomyFindByID finds taxonomy by ID
    TaxonomyFindByID(ctx context.Context, id string) (*Taxonomy, error)
    
    // TaxonomyFindBySlug finds taxonomy by slug
    TaxonomyFindBySlug(ctx context.Context, slug string) (*Taxonomy, error)
    
    // TaxonomyList lists all taxonomies
    TaxonomyList(ctx context.Context) ([]Taxonomy, error)
    
    // TaxonomyUpdate updates a taxonomy
    TaxonomyUpdate(ctx context.Context, tax Taxonomy) error
    
    // TaxonomyDelete deletes a taxonomy and all its terms
    TaxonomyDelete(ctx context.Context, id string) error
    
    // TaxonomyTrash soft-deletes a taxonomy
    TaxonomyTrash(ctx context.Context, id string, deletedBy string) error
    
    // TaxonomyRestore restores a taxonomy from trash
    TaxonomyRestore(ctx context.Context, id string) error
    
    // TaxonomyListTrash lists deleted taxonomies
    TaxonomyListTrash(ctx context.Context) ([]TaxonomyTrash, error)
    
    // ==========================================
    // Taxonomy Terms
    // ==========================================
    
    // TaxonomyTermCreate creates a term within a taxonomy
    TaxonomyTermCreate(ctx context.Context, opts TaxonomyTermOptions) (*TaxonomyTerm, error)
    
    // TaxonomyTermFindByID finds term by ID
    TaxonomyTermFindByID(ctx context.Context, id string) (*TaxonomyTerm, error)
    
    // TaxonomyTermFindBySlug finds term by slug within a taxonomy
    TaxonomyTermFindBySlug(ctx context.Context, taxonomyID, slug string) (*TaxonomyTerm, error)
    
    // TaxonomyTermList lists all terms in a taxonomy
    TaxonomyTermList(ctx context.Context, taxonomyID string) ([]TaxonomyTerm, error)
    
    // TaxonomyTermListTree returns hierarchical tree of terms
    TaxonomyTermListTree(ctx context.Context, taxonomyID string, parentID string) ([]TaxonomyTerm, error)
    
    // TaxonomyTermUpdate updates a term
    TaxonomyTermUpdate(ctx context.Context, term TaxonomyTerm) error
    
    // TaxonomyTermDelete deletes a term
    TaxonomyTermDelete(ctx context.Context, id string) error
    
    // TaxonomyTermTrash soft-deletes a term
    TaxonomyTermTrash(ctx context.Context, id string, deletedBy string) error
    
    // TaxonomyTermRestore restores a term from trash
    TaxonomyTermRestore(ctx context.Context, id string) error
    
    // TaxonomyTermListTrash lists deleted terms
    TaxonomyTermListTrash(ctx context.Context, taxonomyID string) ([]TaxonomyTermTrash, error)
    
    // ==========================================
    // Entity-Taxonomy Associations
    // ==========================================
    
    // EntityTaxonomyAssign assigns an entity to a taxonomy term
    EntityTaxonomyAssign(ctx context.Context, entityID, taxonomyID, termID string) error
    
    // EntityTaxonomyRemove removes an entity from a taxonomy term
    EntityTaxonomyRemove(ctx context.Context, entityID, taxonomyID, termID string) error
    
    // EntityTaxonomyList lists all taxonomy terms assigned to an entity
    EntityTaxonomyList(ctx context.Context, entityID string) ([]EntityTaxonomy, error)
    
    // EntityTaxonomyListByTaxonomy lists entities by taxonomy term
    EntityTaxonomyListByTaxonomy(ctx context.Context, taxonomyID, termID string) ([]EntityTaxonomy, error)
}
```

### 3.2b Query Interface Definitions

**File:** `taxonomy_query_interface.go`

```go
// TaxonomyQueryInterface provides query capabilities for taxonomies
type TaxonomyQueryInterface interface {
	// Select enables method chaining
	Select() TaxonomyQueryInterface

	// Filters
	SetID(id string) TaxonomyQueryInterface
	SetSlug(slug string) TaxonomyQueryInterface
	SetParentID(parentID string) TaxonomyQueryInterface

	// Execution
	Count(ctx context.Context) (int64, error)
	List(ctx context.Context) ([]TaxonomyInterface, error)
	First(ctx context.Context) (TaxonomyInterface, error)
}
```

**File:** `taxonomy_term_query_interface.go`

```go
// TaxonomyTermQueryInterface provides query capabilities for taxonomy terms
type TaxonomyTermQueryInterface interface {
	// Select enables method chaining
	Select() TaxonomyTermQueryInterface

	// Filters
	SetID(id string) TaxonomyTermQueryInterface
	SetTaxonomyID(taxonomyID string) TaxonomyTermQueryInterface
	SetSlug(slug string) TaxonomyTermQueryInterface
	SetParentID(parentID string) TaxonomyTermQueryInterface

	// Execution
	Count(ctx context.Context) (int64, error)
	List(ctx context.Context) ([]TaxonomyTermInterface, error)
	ListTree(ctx context.Context, parentID string) ([]TaxonomyTermInterface, error)
	First(ctx context.Context) (TaxonomyTermInterface, error)
}
```

**File:** `entity_taxonomy_query_interface.go`

```go
// EntityTaxonomyQueryInterface provides query capabilities for entity-taxonomy links
type EntityTaxonomyQueryInterface interface {
	// Select enables method chaining
	Select() EntityTaxonomyQueryInterface

	// Filters
	SetID(id string) EntityTaxonomyQueryInterface
	SetEntityID(entityID string) EntityTaxonomyQueryInterface
	SetTaxonomyID(taxonomyID string) EntityTaxonomyQueryInterface
	SetTermID(termID string) EntityTaxonomyQueryInterface

	// Execution
	Count(ctx context.Context) (int64, error)
	List(ctx context.Context) ([]EntityTaxonomyInterface, error)
	First(ctx context.Context) (EntityTaxonomyInterface, error)
}
```

### 3.3 Database Schema Additions

Uses **9-char short IDs** (varchar(9)) for space efficiency, matching cmsstore pattern:

**File:** `store_implementation.go` (extend SqlCreateTable)

```go
func (st *storeImplementation) SqlCreateTable() ([]string, error) {
    // ... existing entity and attribute tables ...
    
    sqlArray := existingSQL
    
    // Only create taxonomy tables if enabled
    if st.taxonomiesEnabled {
        sqlTaxonomy := `
        CREATE TABLE IF NOT EXISTS ` + st.taxonomyTableName + ` (
            id varchar(9) NOT NULL PRIMARY KEY,
            name varchar(255) NOT NULL,
            slug varchar(255) NOT NULL,
            description text,
            parent_id varchar(9),
            entity_types text,
            created_at datetime NOT NULL,
            updated_at datetime NOT NULL,
            UNIQUE KEY unique_slug (slug)
        );`
        
        sqlTaxonomyTerm := `
        CREATE TABLE IF NOT EXISTS ` + st.taxonomyTermTableName + ` (
            id varchar(9) NOT NULL PRIMARY KEY,
            taxonomy_id varchar(9) NOT NULL,
            name varchar(255) NOT NULL,
            slug varchar(255) NOT NULL,
            parent_id varchar(9),
            sort_order int DEFAULT 0,
            created_at datetime NOT NULL,
            updated_at datetime NOT NULL,
            UNIQUE KEY unique_taxonomy_slug (taxonomy_id, slug),
            INDEX idx_taxonomy (taxonomy_id),
            INDEX idx_parent (parent_id)
        );`
        
        sqlEntityTaxonomy := `
        CREATE TABLE IF NOT EXISTS ` + st.entityTaxonomyTableName + ` (
            id varchar(9) NOT NULL PRIMARY KEY,
            entity_id varchar(9) NOT NULL,
            taxonomy_id varchar(9) NOT NULL,
            term_id varchar(9) NOT NULL,
            created_at datetime NOT NULL,
            UNIQUE KEY unique_entity_term (entity_id, taxonomy_id, term_id),
            INDEX idx_entity (entity_id),
            INDEX idx_taxonomy (taxonomy_id),
            INDEX idx_term (term_id)
        );`
        
        sqlArray = append(sqlArray, sqlTaxonomy, sqlTaxonomyTerm, sqlEntityTaxonomy)
        
        // Create trash tables for soft delete
        sqlTaxonomyTrash := `
        CREATE TABLE IF NOT EXISTS ` + st.taxonomyTrashTableName + ` (
            id varchar(9) NOT NULL PRIMARY KEY,
            name varchar(255) NOT NULL,
            slug varchar(255) NOT NULL,
            description text,
            parent_id varchar(9),
            entity_types text,
            created_at datetime NOT NULL,
            updated_at datetime NOT NULL,
            deleted_at datetime NOT NULL,
            deleted_by varchar(9)
        );`
        
        sqlTaxonomyTermTrash := `
        CREATE TABLE IF NOT EXISTS ` + st.taxonomyTermTrashTableName + ` (
            id varchar(9) NOT NULL PRIMARY KEY,
            taxonomy_id varchar(9) NOT NULL,
            name varchar(255) NOT NULL,
            slug varchar(255) NOT NULL,
            parent_id varchar(9),
            sort_order int DEFAULT 0,
            created_at datetime NOT NULL,
            updated_at datetime NOT NULL,
            deleted_at datetime NOT NULL,
            deleted_by varchar(9)
        );`
        
        sqlArray = append(sqlArray, sqlTaxonomyTrash, sqlTaxonomyTermTrash)
    }
    
    return sqlArray, nil
}
```

### 3.4 Implementation Files (dataobject Pattern)

Following the same 8-file pattern as entities and attributes for each taxonomy entity:

| File | Purpose | Lines (Est) |
|------|---------|-------------|
| `taxonomy_implementation.go` | dataobject-based struct with getters/setters | 130 |
| `taxonomy_implementation_test.go` | Tests for taxonomy implementation | 100 |
| `taxonomy_query.go` | Query builder for taxonomies | 120 |
| `taxonomy_query_interface.go` | Query interface definitions | 50 |
| `taxonomy_query_test.go` | Query tests | 80 |
| `taxonomy_table_create_sql.go` | SQL schema using sb builder | 40 |
| `store_taxonomies.go` | Store CRUD methods | 150 |
| `store_taxonomies_test.go` | Store method tests | 100 |
| **Taxonomy Subtotal** | | **~770** |

Plus similar 8-file sets for:
- TaxonomyTerm (8 files)
- EntityTaxonomy (8 files)

| `taxonomy_trash_implementation.go` | Trash type for soft delete | 80 |
| `taxonomy_trash_table_create_sql.go` | Trash table SQL schema | 40 |
| `store_taxonomies_trash.go` | Trash/restore store methods | 80 |
| `taxonomy_term_trash_implementation.go` | Term trash type | 80 |
| `taxonomy_term_trash_table_create_sql.go` | Term trash table SQL | 40 |
| `store_taxonomy_terms_trash.go` | Term trash/restore methods | 80 |
| **Trash Subtotal** | | **~360** |

**Prerequisite:** This should only be implemented AFTER the dataobject pattern for entities and attributes is stable.

---

## 4. Usage Examples

### 4.1 Enabling Taxonomies (Optional Feature)

```go
// Create store WITH taxonomy support
store, _ := entitystore.NewStore(entitystore.NewStoreOptions{
    DB:                      db,
    EntityTableName:         "entities",
    AttributeTableName:      "attributes",
    TaxonomiesEnabled:       true,                      // Enable taxonomies
    TaxonomyTableName:       "my_taxonomies",           // Optional: custom table names
    TaxonomyTermTableName:   "my_taxonomy_terms",
    EntityTaxonomyTableName: "my_entity_taxonomies",
})

// Create store WITHOUT taxonomies (default behavior)
storeMinimal, _ := entitystore.NewStore(entitystore.NewStoreOptions{
    DB:                  db,
    EntityTableName:     "entities",
    AttributeTableName:  "attributes",
    // TaxonomiesEnabled defaults to false
})
```

### 4.2 Basic Taxonomy (Categories)

```go
store, _ := entitystore.NewStore(entitystore.NewStoreOptions{
    DB: db,
    EntityTableName: "entities",
    AttributeTableName: "attributes",
})

// Create product categories taxonomy
categories, _ := store.TaxonomyCreate(ctx, entitystore.TaxonomyOptions{
    Name: "Product Categories",
    Slug: "product_categories",
    EntityTypes: []string{"product"},
})

// Create category terms
electronics, _ := store.TaxonomyTermCreate(ctx, entitystore.TaxonomyTermOptions{
    TaxonomyID: categories.ID(),
    Name: "Electronics",
    Slug: "electronics",
})

phones, _ := store.TaxonomyTermCreate(ctx, entitystore.TaxonomyTermOptions{
    TaxonomyID: categories.ID(),
    Name: "Phones",
    Slug: "phones",
    ParentID: electronics.ID(), // Hierarchical
})

// Assign product to category
entity := store.EntityCreateWithType(ctx, "product")
entity.SetAttribute("name", "iPhone 15")
store.EntityTaxonomyAssign(ctx, entity.ID(), categories.ID(), phones.ID())

// Find all products in Electronics category
assignments, _ := store.EntityTaxonomyListByTaxonomy(ctx, categories.ID(), electronics.ID())
for _, assignment := range assignments {
    entity, _ := store.EntityFindByID(ctx, assignment.EntityID())
    fmt.Println(entity.GetAttribute("name"))
}
```

### 4.2 Tags (Flat Taxonomy)

```go
// Create tags taxonomy (no hierarchy)
tags, _ := store.TaxonomyCreate(ctx, entitystore.TaxonomyOptions{
    Name: "Tags",
    Slug: "tags",
    EntityTypes: []string{"post", "product"},
})

// Create tag terms
golang, _ := store.TaxonomyTermCreate(ctx, entitystore.TaxonomyTermOptions{
    TaxonomyID: tags.ID(),
    Name: "Golang",
    Slug: "golang",
})

tutorial, _ := store.TaxonomyTermCreate(ctx, entitystore.TaxonomyTermOptions{
    TaxonomyID: tags.ID(),
    Name: "Tutorial",
    Slug: "tutorial",
})

// Assign multiple tags to post
post := store.EntityCreateWithType(ctx, "post")
store.EntityTaxonomyAssign(ctx, post.ID(), tags.ID(), golang.ID())
store.EntityTaxonomyAssign(ctx, post.ID(), tags.ID(), tutorial.ID())

// Get all tags for post
postTags, _ := store.EntityTaxonomyList(ctx, post.ID())
```

### 4.3 Hierarchical Tree Navigation

```go
// Get full category tree
tree, _ := store.TaxonomyTermListTree(ctx, categories.ID(), "")
// Returns: Electronics -> Phones -> Smartphones

// Get children of Electronics only
phoneCategories, _ := store.TaxonomyTermListTree(ctx, categories.ID(), electronics.ID())
```

---

## 5. Implementation Phases ⏳ PENDING

**Note:** Implementation is pending stabilization of the dataobject pattern for entities and attributes.

### Phase 1: Taxonomy Core (1 day)

1. Create `taxonomy_implementation.go` with dataobject pattern
2. Create `taxonomy_implementation_test.go`
3. Add constants for taxonomy columns in `consts.go`
4. Add `TaxonomyInterface` to `interfaces.go`

### Phase 2: Query & SQL (1 day)

1. Create `taxonomy_query.go` with query builder
2. Create `taxonomy_query_interface.go`
3. Create `taxonomy_query_test.go`
4. Create `taxonomy_table_create_sql.go` with sb builder pattern

### Phase 3: Store Methods (1 day)

1. Create `store_taxonomies.go` with CRUD operations
2. Create `store_taxonomies_test.go`
3. Update `store_implementation.go` to include taxonomy table in AutoMigrate

### Phase 4: Taxonomy Terms (1.5 days)

Repeat Phases 1-3 for TaxonomyTerm (8 files)
- Hierarchical tree queries (`TaxonomyTermListTree`)

### Phase 5: Entity-Taxonomy Integration (1 day)

Repeat Phases 1-3 for EntityTaxonomy (8 files)
- Entity-taxonomy assignment operations

### Phase 6: Trash Support (1 day)

1. Create `taxonomy_trash_implementation.go` and test
2. Create `taxonomy_term_trash_implementation.go` and test
3. Create `store_taxonomies_trash.go` with Trash/Restore/ListTrash methods
4. Create `store_taxonomy_terms_trash.go` with Trash/Restore/ListTrash methods
5. Add trash table SQL to `SqlCreateTable`

### Phase 7: Documentation (0.5 day)

1. Update README.md with taxonomy examples
2. Add usage guide

**Total: ~7 days** (after dataobject pattern is stable)

---

## 6. Testing Strategy

### 6.1 Unit Tests (dataobject Pattern)

| Test File | Coverage |
|-----------|----------|
| `taxonomy_implementation_test.go` | Type definition, getters, setters |
| `taxonomy_query_test.go` | Query builder |
| `store_taxonomies_test.go` | Store CRUD operations |
| `taxonomy_term_implementation_test.go` | Term type definition |
| `taxonomy_term_query_test.go` | Term query builder |
| `store_taxonomy_terms_test.go` | Term store CRUD |
| `entity_taxonomy_implementation_test.go` | EntityTaxonomy type |
| `taxonomy_trash_implementation_test.go` | Trash type definition |
| `store_taxonomies_trash_test.go` | Trash/restore operations |
| `taxonomy_term_trash_implementation_test.go` | Term trash type |
| `store_taxonomy_terms_trash_test.go` | Term trash/restore operations |

### 6.2 Integration Tests

```go
func TestTaxonomyLifecycle(t *testing.T) {
    // Create taxonomy
    // Create terms with hierarchy
    // Assign entities
    // Trash taxonomy (soft delete)
    // Verify taxonomy in trash
    // Restore taxonomy
    // Verify taxonomy restored
    // Hard delete taxonomy
    // Verify taxonomy gone
}

func TestTaxonomyTrashCascade(t *testing.T) {
    // Create Electronics taxonomy
    // Create Phones term
    // Assign products to Phones
    // Trash Electronics taxonomy
    // Verify Phones term also trashed
    // Verify entity assignments preserved (or cascade removed)
    // Restore Electronics
    // Verify Phones also restored
}

func TestTaxonomyHierarchy(t *testing.T) {
    // Create Electronics > Phones > Smartphones
    // Verify tree structure
    // Move Phones to different parent
    // Verify reorganization
}
```

---

## 7. Backward Compatibility

### 7.1 Database Migration

New tables are additive:
- `entities_taxonomies` - new table
- `entities_taxonomy_terms` - new table
- `entities_entity_taxonomies` - new table

Existing code continues to work unchanged.

### 7.2 API Compatibility

- All existing methods unchanged
- New methods added to `StoreInterface`
- Existing implementations don't break

### 7.3 Migration Path

1. Update `entitystore` package
2. Run `store.AutoMigrate(ctx)` to create new tables
3. Optionally migrate existing attribute-based taxonomies

---

## 8. Risks and Mitigations

| Risk | Mitigation |
|------|------------|
| Taxonomy tree performance | Cache tree structure, lazy load children |
| Duplicate term slugs | Unique constraint on (taxonomy_id, slug) |
| Orphaned entity associations | FK constraints or cleanup job |
| Many-to-many complexity | Clear documentation with examples |
| Feature bloat for simple use cases | Optional flag - only create tables if enabled |

---

## 9. Future Considerations

### 9.1 Potential Enhancements

1. **Taxonomy permissions** - Which roles can assign which taxonomies
2. **Taxonomy synonyms** - Alternative names for terms
3. **Term metadata** - Store additional data on terms (icon, color, etc.)
4. **Auto-suggest** - Typeahead for term assignment
5. **Bulk operations** - Assign taxonomy to multiple entities

### 9.2 Query Optimization

Future PR could add:
```go
// Filter entities by taxonomy
store.EntityList(ctx, entitystore.EntityQueryOptions{
    EntityType: "product",
    TaxonomyFilter: map[string][]string{
        "categories": ["electronics", "phones"],
        "tags": ["featured"],
    },
})
```

---

## 10. Conclusion

### Status

⏳ **PENDING** - This proposal is on hold until the dataobject pattern implementation for entities and attributes is fully stabilized.

### Recommendation

**Proceed with implementation AFTER:**
1. ✅ dataobject pattern for entities is stable
2. ✅ dataobject pattern for attributes is stable
3. ✅ Trash table implementation is stable
4. ✅ All 4 core entities (Entity, Attribute, EntityTrash, AttributeTrash) are complete

Then implement taxonomy following the same dataobject pattern:
- Taxonomy (8 files)
- TaxonomyTerm (8 files)
- EntityTaxonomy (8 files)

### Benefits

1. **Completeness** - Entity store now supports full content management patterns
2. **Performance** - Native taxonomy queries vs. attribute parsing
3. **Hierarchy** - Built-in support for nested categories
4. **Reusability** - All entitystore consumers get taxonomy for free
5. **Consistency** - Uses same dataobject pattern as entities and attributes

### Next Actions

1. ⏳ Wait for dataobject pattern stabilization
2. ⏳ Create feature branch for taxonomy
3. ⏳ Implement following 8-file pattern
4. ⏳ PR and review

---

**End of Proposal - Updated March 28, 2026**
**Status: PENDING (waiting for dataobject pattern stabilization)**
