# Entity Taxonomy Support Proposal

**Date:** 2026-03-28
**Status:** Draft
**Author:** AI Assistant
**Repository:** github.com/dracory/entitystore

---

## 1. Executive Summary

**Problem:** `entitystore` provides excellent EAV storage but lacks native support for categorizing/classifying entities (taxonomy).

**Current Workaround:** Store taxonomy IDs as string attributes, query separately.

```go
// Current workaround - manual attribute storage
product.SetString("category_id", "cat_123")
product.SetString("tag_ids", `["tag_1", "tag_2"]`)

// Manual filtering in application code
categoryID := product.GetString("category_id", "")
// Query all products, filter in code - inefficient
```

**Solution:** Add native taxonomy support to `entitystore`.

**Impact:**
- All projects using `entitystore` get categorization for free
- Enables taxonomy-based queries ("find all products in Electronics")
- Supports hierarchical categories, tags, labels, etc.

---

## 2. Current State (As-Is)

### 2.1 Entitystore Architecture

```
entitystore/
├── entity.go                 # Entity struct with Get/Set methods
├── entity_*.go              # Entity operations (CRUD, List, Trash)
├── attribute.go             # Attribute struct
├── attribute_*.go           # Attribute operations
├── store_implementation.go  # Core store with table management
├── interfaces.go            # StoreInterface definition
└── new.go                   # Store initialization
```

### 2.2 Current Database Schema

```sql
-- Entities table (exists)
CREATE TABLE entities (
    id varchar(40) PRIMARY KEY,
    entity_type varchar(40) NOT NULL,
    entity_handle varchar(60),
    created_at datetime NOT NULL,
    updated_at datetime NOT NULL
);

-- Attributes table (exists)
CREATE TABLE attributes (
    id varchar(40) PRIMARY KEY,
    entity_id varchar(40) NOT NULL,
    attribute_key varchar(255) NOT NULL,
    attribute_value text,
    created_at datetime NOT NULL,
    updated_at datetime NOT NULL
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

## 3. Proposed Solution (To-Be)

### 3.1 New Types

**File:** `taxonomy.go`

```go
package entitystore

import "time"

// Taxonomy is a categorization system (categories, tags, etc.)
// ID is 9-char short ID (e.g., "86ccrtsgx") for space efficiency
type Taxonomy struct {
	id          string // 9-char short ID
	name        string
	slug        string
	description string
	parentID    string // 9-char short ID
	entityTypes []string
	createdAt   time.Time
	updatedAt   time.Time
	st          *storeImplementation
}

// TaxonomyTerm is a single term/category within a taxonomy
// ID is 9-char short ID (e.g., "86ccrtsgx") for space efficiency
type TaxonomyTerm struct {
	id         string // 9-char short ID
	taxonomyID string // 9-char short ID
	name       string
	slug       string
	parentID   string // 9-char short ID
	sortOrder  int
	createdAt  time.Time
	updatedAt  time.Time
	st         *storeImplementation
}

// EntityTaxonomy links an entity to a taxonomy term
// Uses 9-char short IDs for all ID fields
type EntityTaxonomy struct {
	id         string // 9-char short ID
	entityID   string // 9-char short ID
	taxonomyID string // 9-char short ID
	termID     string // 9-char short ID
	createdAt  time.Time
}

// ==========================================
// Taxonomy Getters
// ==========================================

func (t *Taxonomy) ID() string {
	return t.id
}

func (t *Taxonomy) Name() string {
	return t.name
}

func (t *Taxonomy) Slug() string {
	return t.slug
}

func (t *Taxonomy) Description() string {
	return t.description
}

func (t *Taxonomy) ParentID() string {
	return t.parentID
}

func (t *Taxonomy) EntityTypes() []string {
	return t.entityTypes
}

func (t *Taxonomy) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Taxonomy) UpdatedAt() time.Time {
	return t.updatedAt
}

// ==========================================
// Taxonomy Setters (fluent interface)
// ==========================================

func (t *Taxonomy) SetID(id string) *Taxonomy {
	t.id = id
	return t
}

func (t *Taxonomy) SetName(name string) *Taxonomy {
	t.name = name
	return t
}

func (t *Taxonomy) SetSlug(slug string) *Taxonomy {
	t.slug = slug
	return t
}

func (t *Taxonomy) SetDescription(desc string) *Taxonomy {
	t.description = desc
	return t
}

func (t *Taxonomy) SetParentID(parentID string) *Taxonomy {
	t.parentID = parentID
	return t
}

func (t *Taxonomy) SetEntityTypes(types []string) *Taxonomy {
	t.entityTypes = types
	return t
}

func (t *Taxonomy) SetCreatedAt(createdAt time.Time) *Taxonomy {
	t.createdAt = createdAt
	return t
}

func (t *Taxonomy) SetUpdatedAt(updatedAt time.Time) *Taxonomy {
	t.updatedAt = updatedAt
	return t
}

func (t *Taxonomy) ToMap() map[string]any {
	return map[string]any{
		"id":           t.ID(),
		"name":         t.Name(),
		"slug":         t.Slug(),
		"description":  t.Description(),
		"parent_id":    t.ParentID(),
		"entity_types": t.EntityTypes(),
		"created_at":   t.CreatedAt(),
		"updated_at":   t.UpdatedAt(),
	}
}

// ==========================================
// TaxonomyTerm Getters
// ==========================================

func (tt *TaxonomyTerm) ID() string {
	return tt.id
}

func (tt *TaxonomyTerm) TaxonomyID() string {
	return tt.taxonomyID
}

func (tt *TaxonomyTerm) Name() string {
	return tt.name
}

func (tt *TaxonomyTerm) Slug() string {
	return tt.slug
}

func (tt *TaxonomyTerm) ParentID() string {
	return tt.parentID
}

func (tt *TaxonomyTerm) SortOrder() int {
	return tt.sortOrder
}

func (tt *TaxonomyTerm) CreatedAt() time.Time {
	return tt.createdAt
}

func (tt *TaxonomyTerm) UpdatedAt() time.Time {
	return tt.updatedAt
}

// ==========================================
// TaxonomyTerm Setters (fluent interface)
// ==========================================

func (tt *TaxonomyTerm) SetID(id string) *TaxonomyTerm {
	tt.id = id
	return tt
}

func (tt *TaxonomyTerm) SetTaxonomyID(taxID string) *TaxonomyTerm {
	tt.taxonomyID = taxID
	return tt
}

func (tt *TaxonomyTerm) SetName(name string) *TaxonomyTerm {
	tt.name = name
	return tt
}

func (tt *TaxonomyTerm) SetSlug(slug string) *TaxonomyTerm {
	tt.slug = slug
	return tt
}

func (tt *TaxonomyTerm) SetParentID(parentID string) *TaxonomyTerm {
	tt.parentID = parentID
	return tt
}

func (tt *TaxonomyTerm) SetSortOrder(order int) *TaxonomyTerm {
	tt.sortOrder = order
	return tt
}

func (tt *TaxonomyTerm) SetCreatedAt(createdAt time.Time) *TaxonomyTerm {
	tt.createdAt = createdAt
	return tt
}

func (tt *TaxonomyTerm) SetUpdatedAt(updatedAt time.Time) *TaxonomyTerm {
	tt.updatedAt = updatedAt
	return tt
}

func (tt *TaxonomyTerm) ToMap() map[string]any {
	return map[string]any{
		"id":          tt.ID(),
		"taxonomy_id": tt.TaxonomyID(),
		"name":        tt.Name(),
		"slug":        tt.Slug(),
		"parent_id":   tt.ParentID(),
		"sort_order":  tt.SortOrder(),
		"created_at":  tt.CreatedAt(),
		"updated_at":  tt.UpdatedAt(),
	}
}

// ==========================================
// EntityTaxonomy Getters
// ==========================================

func (et *EntityTaxonomy) ID() string {
	return et.id
}

func (et *EntityTaxonomy) EntityID() string {
	return et.entityID
}

func (et *EntityTaxonomy) TaxonomyID() string {
	return et.taxonomyID
}

func (et *EntityTaxonomy) TermID() string {
	return et.termID
}

func (et *EntityTaxonomy) CreatedAt() time.Time {
	return et.createdAt
}

// ==========================================
// EntityTaxonomy Setters (fluent interface)
// ==========================================

func (et *EntityTaxonomy) SetID(id string) *EntityTaxonomy {
	et.id = id
	return et
}

func (et *EntityTaxonomy) SetEntityID(entityID string) *EntityTaxonomy {
	et.entityID = entityID
	return et
}

func (et *EntityTaxonomy) SetTaxonomyID(taxID string) *EntityTaxonomy {
	et.taxonomyID = taxID
	return et
}

func (et *EntityTaxonomy) SetTermID(termID string) *EntityTaxonomy {
	et.termID = termID
	return et
}

func (et *EntityTaxonomy) SetCreatedAt(createdAt time.Time) *EntityTaxonomy {
	et.createdAt = createdAt
	return et
}

func (et *EntityTaxonomy) ToMap() map[string]any {
	return map[string]any{
		"id":          et.ID(),
		"entity_id":   et.EntityID(),
		"taxonomy_id": et.TaxonomyID(),
		"term_id":     et.TermID(),
		"created_at":  et.CreatedAt(),
	}
}
```

**Options structs:**

```go
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

### 3.3 Database Schema Additions

Uses **9-char short IDs** (varchar(9)) for space efficiency, matching cmsstore pattern:

**File:** `store_implementation.go` (extend SqlCreateTable)

```go
func (st *storeImplementation) SqlCreateTable() ([]string, error) {
    // ... existing entity and attribute tables ...
    
    sqlTaxonomy := `
    CREATE TABLE IF NOT EXISTS ` + st.entityTableName + `_taxonomies (
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
    CREATE TABLE IF NOT EXISTS ` + st.entityTableName + `_taxonomy_terms (
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
    CREATE TABLE IF NOT EXISTS ` + st.entityTableName + `_entity_taxonomies (
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
    
    sqlArray := append(existingSQL, sqlTaxonomy, sqlTaxonomyTerm, sqlEntityTaxonomy)
    return sqlArray, nil
}
```

### 3.4 Implementation Files

| File | Purpose | Lines (Est) |
|------|---------|-------------|
| `taxonomy.go` | Taxonomy, Term, EntityTaxonomy types | 300 |
| `taxonomy_create.go` | Create taxonomy | 50 |
| `taxonomy_find.go` | Find taxonomy | 50 |
| `taxonomy_list.go` | List taxonomies | 50 |
| `taxonomy_update.go` | Update taxonomy | 50 |
| `taxonomy_delete.go` | Delete taxonomy | 50 |
| `taxonomy_term_*.go` | Term CRUD operations | 250 |
| `entity_taxonomy_*.go` | Entity-taxonomy association | 150 |
| `new_taxonomy.go` | Constructor functions | 80 |
| `interfaces.go` | Extend StoreInterface | +60 |
| `store_implementation.go` | Update SqlCreateTable | +40 |
| **Total** | | **~700** |

---

## 4. Usage Examples

### 4.1 Basic Taxonomy (Categories)

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
product := store.EntityCreateWithType(ctx, "product")
product.SetString("name", "iPhone 15")
store.EntityTaxonomyAssign(ctx, product.ID(), categories.ID(), phones.ID())

// Find all products in Electronics category
assignments, _ := store.EntityTaxonomyListByTaxonomy(ctx, categories.ID(), electronics.ID())
for _, assignment := range assignments {
    product, _ := store.EntityFindByID(ctx, assignment.EntityID())
    fmt.Println(product.GetString("name", ""))
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

## 5. Implementation Phases

### Phase 1: Taxonomy Core (2 days)

1. Create `taxonomy.go` with types
2. Create `new_taxonomy.go` with constructors
3. Implement taxonomy CRUD operations
4. Write tests

### Phase 2: Taxonomy Terms (2 days)

1. Implement term CRUD operations
2. Implement hierarchical tree queries (`TaxonomyTermListTree`)
3. Write tests

### Phase 3: Entity-Taxonomy Integration (2 days)

1. Implement `EntityTaxonomyAssign`, `EntityTaxonomyRemove`
2. Implement `EntityTaxonomyList`, `EntityTaxonomyListByTaxonomy`
3. Write tests

### Phase 4: Integration (1 day)

1. Update `interfaces.go` with new methods
2. Update `store_implementation.go` with taxonomy tables
3. Add tables to AutoMigrate
4. Write integration tests

### Phase 5: Documentation (1 day)

1. Update README.md
2. Create examples/taxonomy.go
3. Write usage guide

**Total: ~8 days**

---

## 6. Testing Strategy

### 6.1 Unit Tests

| Test File | Coverage |
|-----------|----------|
| `taxonomy_test.go` | Type definitions, getters, setters |
| `taxonomy_create_test.go` | Create taxonomy, slug uniqueness |
| `taxonomy_find_test.go` | Find by ID, find by slug |
| `taxonomy_list_test.go` | List taxonomies |
| `taxonomy_term_create_test.go` | Create terms, hierarchy |
| `taxonomy_term_tree_test.go` | Tree structure, parent/child |
| `entity_taxonomy_test.go` | Assign, remove, list associations |

### 6.2 Integration Tests

```go
func TestTaxonomyLifecycle(t *testing.T) {
    // Create taxonomy
    // Create terms with hierarchy
    // Assign entities
    // Delete taxonomy (cascade to terms and associations)
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

### Recommendation

**Proceed with implementation.**

Taxonomy is a fundamental CMS/e-commerce feature. Adding native support to `entitystore` enables categorization, filtering, and navigation patterns essential for complex applications.

### Benefits

1. **Completeness** - Entity store now supports full content management patterns
2. **Performance** - Native taxonomy queries vs. attribute parsing
3. **Hierarchy** - Built-in support for nested categories
4. **Reusability** - All entitystore consumers get taxonomy for free

### Next Actions

1. Create feature branch
2. Implement Phase 1 (Taxonomy Core)
3. PR and review
4. Proceed to Phase 2-5

---

**End of Proposal**
