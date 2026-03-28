# Entity Relationships Support Proposal

**Date:** 2026-03-28
**Status:** ⏳ PENDING (waiting for dataobject pattern stabilization)
**Author:** AI Assistant
**Repository:** github.com/dracory/entitystore

---

## 1. Executive Summary

**Problem:** `entitystore` provides excellent EAV storage but lacks native support for linking entities together (relationships).

**Current Workaround:** Store relationship IDs as string attributes, query separately.

```go
// Current workaround - manual attribute storage
entity.SetAttribute("author_id", "auth_123")
authorID := entity.GetAttribute("author_id")
author := store.EntityFindByID(ctx, authorID) // Separate query
```

**Solution:** Add native relationship support to `entitystore`.

**Impact:**
- All projects using `entitystore` get entity-to-entity linking for free
- Eliminates manual join logic in application code
- Enables relationship queries ("find all books by this author")

**Status:** ⏳ **PENDING** - Waiting for dataobject pattern implementation to stabilize before adding relationships

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
├── attribute_implementation_test.go  # Attribute tests
├── attribute_query.go                # Attribute query builder
├── attribute_query_interface.go      # AttributeQueryInterface
├── attribute_table_create_sql.go     # Attributes table SQL
├── entity_trash_implementation.go    # EntityTrash with dataobject
├── entity_trash_implementation_test.go
├── entity_trash_query_interface.go
├── entity_trash_table_create_sql.go
├── attribute_trash_implementation.go   # AttributeTrash with dataobject
├── attribute_trash_implementation_test.go
├── attribute_trash_query_interface.go
├── attribute_trash_table_create_sql.go
├── store_entities.go                  # Entity CRUD methods
├── store_entities_test.go             # Entity store tests
├── store_attributes.go                # Attribute CRUD methods
├── store_attributes_test.go           # Attribute store tests
├── store_entities_trash.go            # Entity trash/restore
├── store_attributes_trash.go          # Attribute trash/restore
├── interfaces.go                      # All entity interfaces
├── consts.go                          # Column constants
├── id_helpers.go                      # GenerateShortID()
└── new.go                             # Store initialization
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

**Problems with attribute-based relationships:**
- No referential integrity
- No cascade operations
- Complex queries require multiple round trips
- Trash/restore doesn't handle relationships
- No relationship metadata (e.g., "order" in belongs_to)

---

## 3. Proposed Solution (To-Be) ⏳ PENDING

**Note:** This proposal is pending stabilization of the dataobject pattern implementation. Once the core entitystore architecture (entities, attributes, trash tables with dataobject) is stable, relationships can be added following the same pattern.

### 3.1 New Types (dataobject Pattern)

Following the same pattern as entities and attributes:

**File:** `relationship_implementation.go`

```go
package entitystore

import "github.com/dracory/dataobject"

// Relationship types
const (
	RELATIONSHIP_TYPE_BELONGS_TO = "belongs_to"  // Entity belongs to one parent
	RELATIONSHIP_TYPE_HAS_MANY   = "has_many"    // Entity has many children
	RELATIONSHIP_TYPE_MANY_MANY  = "many_to_many" // Entities linked bidirectionally
)

// relationshipImplementation implements RelationshipInterface
type relationshipImplementation struct {
	*dataobject.DataObject
}

// Column constants for relationships
const (
	COLUMN_PARENT_ID   = "parent_id"
	COLUMN_SEQUENCE    = "sequence"
)

// NewRelationship creates a new relationship instance
func NewRelationship() RelationshipInterface {
	return &relationshipImplementation{
		DataObject: dataobject.NewDataObject(),
	}
}

// ID returns the relationship ID
func (o *relationshipImplementation) ID() string {
	return o.Get(COLUMN_ID)
}

// SetID sets the relationship ID (fluent)
func (o *relationshipImplementation) SetID(id string) RelationshipInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// EntityID returns the source entity ID
func (o *relationshipImplementation) EntityID() string {
	return o.Get(COLUMN_ENTITY_ID)
}

// SetEntityID sets the source entity ID (fluent)
func (o *relationshipImplementation) SetEntityID(entityID string) RelationshipInterface {
	o.Set(COLUMN_ENTITY_ID, entityID)
	return o
}

// RelatedEntityID returns the target entity ID
func (o *relationshipImplementation) RelatedEntityID() string {
	return o.Get(COLUMN_RELATED_ENTITY_ID)
}

// SetRelatedEntityID sets the target entity ID (fluent)
func (o *relationshipImplementation) SetRelatedEntityID(relatedID string) RelationshipInterface {
	o.Set(COLUMN_RELATED_ENTITY_ID, relatedID)
	return o
}

// RelationshipType returns the relationship type
func (o *relationshipImplementation) RelationshipType() string {
	return o.Get(COLUMN_RELATIONSHIP_TYPE)
}

// SetRelationshipType sets the relationship type (fluent)
func (o *relationshipImplementation) SetRelationshipType(relType string) RelationshipInterface {
	o.Set(COLUMN_RELATIONSHIP_TYPE, relType)
	return o
}

// ParentID returns the parent relationship ID (for hierarchical relationships)
func (o *relationshipImplementation) ParentID() string {
	return o.Get(COLUMN_PARENT_ID)
}

// SetParentID sets the parent relationship ID (fluent)
func (o *relationshipImplementation) SetParentID(parentID string) RelationshipInterface {
	o.Set(COLUMN_PARENT_ID, parentID)
	return o
}

// Sequence returns the sequence/order number
func (o *relationshipImplementation) Sequence() int {
	return o.GetInt(COLUMN_SEQUENCE)
}

// SetSequence sets the sequence/order number (fluent)
func (o *relationshipImplementation) SetSequence(sequence int) RelationshipInterface {
	o.Set(COLUMN_SEQUENCE, sequence)
	return o
}

// Metadata returns relationship metadata
func (o *relationshipImplementation) Metadata() map[string]string {
	return o.GetAllAttributes()
}

// SetMetadata sets relationship metadata (fluent)
func (o *relationshipImplementation) SetMetadata(metadata map[string]string) RelationshipInterface {
	for k, v := range metadata {
		o.SetAttribute(k, v)
	}
	return o
}

// CreatedAt returns creation timestamp
func (o *relationshipImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

// SetCreatedAt sets creation timestamp (fluent)
func (o *relationshipImplementation) SetCreatedAt(createdAt string) RelationshipInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}
```
```

### 3.2 Store Interface Extensions

**File:** `interfaces.go` (append to StoreInterface)

```go
type StoreInterface interface {
    // ... existing methods ...
    
    // ==========================================
    // Relationships
    // ==========================================
    
    // RelationshipCreate creates a new relationship between entities
    RelationshipCreate(ctx context.Context, opts RelationshipOptions) (*Relationship, error)
    
    // RelationshipFind finds a specific relationship
    RelationshipFind(ctx context.Context, entityID, relatedID, relType string) (*Relationship, error)
    
    // RelationshipList lists all relationships for an entity
    RelationshipList(ctx context.Context, entityID string, relType string) ([]Relationship, error)
    
    // RelationshipListRelated finds all entities related to the given entity
    RelationshipListRelated(ctx context.Context, relatedID string, relType string) ([]Relationship, error)
    
    // RelationshipDelete removes a relationship
    RelationshipDelete(ctx context.Context, entityID, relatedID, relType string) error
    
    // RelationshipDeleteAll removes all relationships for an entity
    RelationshipDeleteAll(ctx context.Context, entityID string) error
}

type RelationshipInterface interface {
	dataobject.DataObjectInterface

	// Core fields
	ID() string
	SetID(id string) RelationshipInterface
	EntityID() string
	SetEntityID(entityID string) RelationshipInterface
	RelatedEntityID() string
	SetRelatedEntityID(relatedID string) RelationshipInterface
	RelationshipType() string
	SetRelationshipType(relType string) RelationshipInterface

	// Hierarchical support
	ParentID() string
	SetParentID(parentID string) RelationshipInterface
	Sequence() int
	SetSequence(sequence int) RelationshipInterface

	// Metadata
	Metadata() map[string]string
	SetMetadata(metadata map[string]string) RelationshipInterface

	// Timestamps
	CreatedAt() string
	SetCreatedAt(createdAt string) RelationshipInterface
}

// Options struct:

**Store Options (enable/disable feature):**

```go
type NewStoreOptions struct {
    DB                  *sql.DB
    EntityTableName     string
    AttributeTableName  string
    
    // Feature flags - relationships are optional
    RelationshipsEnabled bool   // Enable relationship support
    RelationshipTableName string // Default: "entities_relationships"
}

type RelationshipOptions struct {
    EntityID         string
    RelatedEntityID  string
    RelationshipType string
    ParentID         string            // For hierarchical relationships
    Sequence         int               // For ordering
    Metadata         map[string]string
}
```

### 3.3 Database Schema Addition

Uses **9-char short IDs** (varchar(9)) for space efficiency, matching cmsstore pattern:

**File:** `store_implementation.go` (extend SqlCreateTable)

```go
func (st *storeImplementation) SqlCreateTable() ([]string, error) {
    // ... existing entity and attribute tables ...
    
    sqlArray := existingSQL
    
    // Only create relationships table if enabled
    if st.relationshipsEnabled {
        sqlRelationship := `
        CREATE TABLE IF NOT EXISTS ` + st.relationshipTableName + ` (
            id varchar(9) NOT NULL PRIMARY KEY,
            entity_id varchar(9) NOT NULL,
            related_entity_id varchar(9) NOT NULL,
            relationship_type varchar(50) NOT NULL,
            parent_id varchar(9) DEFAULT NULL,
            sequence int DEFAULT 0,
            metadata text,
            created_at datetime NOT NULL,
            INDEX idx_entity (entity_id),
            INDEX idx_related (related_entity_id),
            INDEX idx_type (relationship_type),
            INDEX idx_entity_type (entity_id, relationship_type),
            INDEX idx_parent (parent_id),
            INDEX idx_sequence (entity_id, sequence)
        );`
        sqlArray = append(sqlArray, sqlRelationship)
    }
    
    return sqlArray, nil
}
```

### 3.4 Implementation Files (dataobject Pattern)

Following the same 8-file pattern as entities and attributes:

| File | Purpose | Lines (Est) |
|------|---------|-------------|
| `relationship_implementation.go` | dataobject-based struct with getters/setters | 130 |
| `relationship_implementation_test.go` | Tests for relationship implementation | 100 |
| `relationship_query.go` | Query builder for relationships | 120 |
| `relationship_query_interface.go` | Query interface definitions | 50 |
| `relationship_query_test.go` | Query tests | 80 |
| `relationship_table_create_sql.go` | SQL schema using sb builder | 40 |
| `store_relationships.go` | Store CRUD methods | 150 |
| `store_relationships_test.go` | Store method tests | 100 |
| **Relationship Subtotal** | | **~770** |

Plus column constants in `consts.go` and interface in `interfaces.go`.

**Prerequisite:** This should only be implemented AFTER the dataobject pattern for entities and attributes is stable.

---

## 4. Usage Examples

### 4.1 Enabling Relationships (Optional Feature)

```go
// Create store WITH relationship support
store, _ := entitystore.NewStore(entitystore.NewStoreOptions{
    DB:                   db,
    EntityTableName:      "entities",
    AttributeTableName:   "attributes",
    RelationshipsEnabled: true,                    // Enable relationships
    RelationshipTableName: "my_relationships",     // Optional: custom table name
})

// Create store WITHOUT relationships (default behavior)
storeMinimal, _ := entitystore.NewStore(entitystore.NewStoreOptions{
    DB:                  db,
    EntityTableName:     "entities",
    AttributeTableName:  "attributes",
    // RelationshipsEnabled defaults to false
})
```

### 4.2 Basic Relationship

```go
store, _ := entitystore.NewStore(entitystore.NewStoreOptions{
    DB: db,
    EntityTableName: "entities",
    AttributeTableName: "attributes",
})

// Create entities (uses 9-char short IDs)
author := store.EntityCreateWithType(ctx, "author")
author.SetString("name", "John Doe")
// author.ID() = "86ccrtsgx" (9 chars)

book := store.EntityCreateWithType(ctx, "book")
book.SetString("title", "Go Programming")
// book.ID() = "86ccrtsgy" (9 chars)

// Create relationship (uses short IDs internally)
store.RelationshipCreate(ctx, entitystore.RelationshipOptions{
    EntityID: book.ID(),           // "86ccrtsgy"
    RelatedEntityID: author.ID(),  // "86ccrtsgx"
    RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
})
```

// Find all books by author
relationships, _ := store.RelationshipListRelated(ctx, author.ID(), entitystore.RELATIONSHIP_TYPE_BELONGS_TO)
for _, rel := range relationships {
    book, _ := store.EntityFindByID(ctx, rel.EntityID())
    fmt.Println(book.GetAttribute("title"))
}
```

### 4.2 Many-to-Many Relationship

```go
// User friends with other users
user1 := store.EntityCreateWithType(ctx, "user")
user2 := store.EntityCreateWithType(ctx, "user")
user3 := store.EntityCreateWithType(ctx, "user")

// Create mutual friendships
store.RelationshipCreate(ctx, entitystore.RelationshipOptions{
    EntityID: user1.ID(),
    RelatedEntityID: user2.ID(),
    RelationshipType: entitystore.RELATIONSHIP_TYPE_MANY_MANY,
    Metadata: map[string]string{"status": "confirmed"},
})

store.RelationshipCreate(ctx, entitystore.RelationshipOptions{
    EntityID: user2.ID(),
    RelatedEntityID: user1.ID(),
    RelationshipType: entitystore.RELATIONSHIP_TYPE_MANY_MANY,
    Metadata: map[string]string{"status": "confirmed"},
})

// Find all friends of user1
friendships, _ := store.RelationshipList(ctx, user1.ID(), entitystore.RELATIONSHIP_TYPE_MANY_MANY)
for _, friendship := range friendships {
    friend, _ := store.EntityFindByID(ctx, friendship.RelatedEntityID())
    status := friendship.Metadata()["status"]
    fmt.Printf("Friend: %s, Status: %s\n", friend.GetAttribute("name"), status)
}
```

### 4.3 Relationship with Metadata

```go
// Store order information on belongs_to relationship
store.RelationshipCreate(ctx, entitystore.RelationshipOptions{
    EntityID: lineItem.ID(),
    RelatedEntityID: order.ID(),
    RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
    Metadata: map[string]string{
        "line_number": "1",
        "quantity": "3",
    },
})
```

### 4.4 Hierarchical Relationships (Trees/Nesting)

```go
// Build a nested category tree using parent_id
electronics := store.EntityCreateWithType(ctx, "category")
electronics.SetString("name", "Electronics")

phones := store.EntityCreateWithType(ctx, "category")
phones.SetString("name", "Phones")

laptops := store.EntityCreateWithType(ctx, "category")
laptops.SetString("name", "Laptops")

// Create root relationship
relElectronics := store.RelationshipCreate(ctx, entitystore.RelationshipOptions{
    EntityID:         electronics.ID(),
    RelatedEntityID:  rootCategory.ID(),
    RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
    Sequence:         1,
})

// Create child relationships with parent_id
relPhones := store.RelationshipCreate(ctx, entitystore.RelationshipOptions{
    EntityID:         phones.ID(),
    RelatedEntityID:  electronics.ID(),
    RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
    ParentID:         relElectronics.ID(),  // Child of electronics
    Sequence:         1,
})

relLaptops := store.RelationshipCreate(ctx, entitystore.RelationshipOptions{
    EntityID:         laptops.ID(),
    RelatedEntityID:  electronics.ID(),
    RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
    ParentID:         relElectronics.ID(),  // Child of electronics
    Sequence:         2,  // After phones
})
```

---

## 5. Implementation Phases ⏳ PENDING

**Note:** Implementation is pending stabilization of the dataobject pattern for entities and attributes.

### Phase 1: Core Types (1 day)

1. Create `relationship_implementation.go` with dataobject pattern
2. Create `relationship_implementation_test.go`
3. Add constants for relationship types in `consts.go`
4. Add `RelationshipInterface` to `interfaces.go`

### Phase 2: Query & SQL (1 day)

1. Create `relationship_query.go` with query builder
2. Create `relationship_query_interface.go`
3. Create `relationship_query_test.go`
4. Create `relationship_table_create_sql.go` with sb builder pattern

### Phase 3: Store Methods (1 day)

1. Create `store_relationships.go` with CRUD operations
2. Create `store_relationships_test.go`
3. Update `store_implementation.go` to include relationship table in AutoMigrate

### Phase 4: Documentation (0.5 day)

1. Update README.md with relationship examples
2. Add usage guide

**Total: ~3.5 days** (after dataobject pattern is stable)

---

## 6. Testing Strategy

### 6.1 Unit Tests (dataobject Pattern)

| Test File | Coverage |
|-----------|----------|
| `relationship_implementation_test.go` | Type definition, getters, setters |
| `relationship_query_test.go` | Query builder |
| `store_relationships_test.go` | Store CRUD operations |
| `store_relationships_trash_test.go` | Trash/restore (if applicable) |

### 6.2 Integration Tests

```go
func TestRelationshipLifecycle(t *testing.T) {
    // Create parent and child
    // Link with relationship
    // Verify relationship exists
    // Delete relationship
    // Verify relationship gone
}

func TestRelationshipMetadata(t *testing.T) {
    // Create relationship with metadata
    // Retrieve and verify metadata preserved
}
```

---

## 7. Backward Compatibility

### 7.1 Database Migration

- `entities_relationships` is a new table
- Existing code continues to work unchanged
- Run `store.AutoMigrate(ctx)` to create new table

### 7.2 API Compatibility

- All existing methods unchanged
- New methods added to `StoreInterface`
- Existing implementations don't break

---

## 8. Risks and Mitigations

| Risk | Mitigation |
|------|------------|
| Relationship query performance | Proper indexing on entity_id, related_entity_id columns |
| Circular relationships | Document best practices, provide validation helpers |
| Cascade delete confusion | Document that relationships don't cascade (by design) |
| Feature bloat for simple use cases | Optional flag - only create tables if enabled |

---

## 9. Future Considerations

### 9.1 Potential Enhancements

1. **Relationship validation** - Ensure entity types match allowed types
2. **Cascade operations** - Optional cascade delete as separate feature
3. **Relationship constraints** - Prevent duplicate relationships
4. **Relationship querying** - Filter entities by relationship criteria

### 9.2 Query Optimization

Future PR could add:
```go
// Preload relationships in single query
store.EntityList(ctx, entitystore.EntityQueryOptions{
    EntityType: "book",
    PreloadRelationships: []string{"author", "publisher"},
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

Then implement relationships following the same dataobject pattern (8 files: implementation, test, query, query interface, query test, SQL schema, store methods, store test).

### Benefits

1. **Simplicity** - No more manual attribute-based linking
2. **Performance** - Single-query relationship loading
3. **Integrity** - Clear relationship semantics
4. **Reusability** - All entitystore consumers benefit
5. **Consistency** - Uses same dataobject pattern as entities and attributes

### Next Actions

1. ⏳ Wait for dataobject pattern stabilization
2. ⏳ Create feature branch for relationships
3. ⏳ Implement following 8-file pattern
4. ⏳ PR and review

---

**End of Proposal - Updated March 28, 2026**
**Status: PENDING (waiting for dataobject pattern stabilization)**
