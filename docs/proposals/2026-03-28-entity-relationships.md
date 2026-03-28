# Entity Relationships Support Proposal

**Date:** 2026-03-28
**Status:** Draft
**Author:** AI Assistant
**Repository:** github.com/dracory/entitystore

---

## 1. Executive Summary

**Problem:** `entitystore` provides excellent EAV storage but lacks native support for linking entities together (relationships).

**Current Workaround:** Store relationship IDs as string attributes, query separately.

```go
// Current workaround - manual attribute storage
product.SetString("author_id", "auth_123")
authorID := product.GetString("author_id", "")
author := store.EntityFindByID(ctx, authorID) // Separate query
```

**Solution:** Add native relationship support to `entitystore`.

**Impact:**
- All projects using `entitystore` get entity-to-entity linking for free
- Eliminates manual join logic in application code
- Enables relationship queries ("find all books by this author")

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

**Problems with attribute-based relationships:**
- No referential integrity
- No cascade operations
- Complex queries require multiple round trips
- Trash/restore doesn't handle relationships
- No relationship metadata (e.g., "order" in belongs_to)

---

## 3. Proposed Solution (To-Be)

### 3.1 New Types

**File:** `relationship.go`

```go
package entitystore

import "time"

// Relationship types
const (
	RELATIONSHIP_TYPE_BELONGS_TO = "belongs_to"  // Entity belongs to one parent
	RELATIONSHIP_TYPE_HAS_MANY   = "has_many"    // Entity has many children
	RELATIONSHIP_TYPE_MANY_MANY  = "many_to_many" // Entities linked bidirectionally
)

// Relationship links two entities
// ID is 9-char short ID (e.g., "86ccrtsgx") for space efficiency
type Relationship struct {
	id               string // 9-char short ID
	entityID         string // 9-char short ID
	relatedEntityID  string // 9-char short ID
	relationshipType string
	metadata         map[string]string
	createdAt        time.Time
	st               *storeImplementation
}

// Getters
func (r *Relationship) ID() string {
	return r.id
}

func (r *Relationship) EntityID() string {
	return r.entityID
}

func (r *Relationship) RelatedEntityID() string {
	return r.relatedEntityID
}

func (r *Relationship) RelationshipType() string {
	return r.relationshipType
}

func (r *Relationship) Metadata() map[string]string {
	return r.metadata
}

func (r *Relationship) CreatedAt() time.Time {
	return r.createdAt
}

// Setters (fluent interface)
func (r *Relationship) SetID(id string) *Relationship {
	r.id = id
	return r
}

func (r *Relationship) SetEntityID(entityID string) *Relationship {
	r.entityID = entityID
	return r
}

func (r *Relationship) SetRelatedEntityID(relatedID string) *Relationship {
	r.relatedEntityID = relatedID
	return r
}

func (r *Relationship) SetRelationshipType(relType string) *Relationship {
	r.relationshipType = relType
	return r
}

func (r *Relationship) SetMetadata(metadata map[string]string) *Relationship {
	r.metadata = metadata
	return r
}

func (r *Relationship) SetCreatedAt(createdAt time.Time) *Relationship {
	r.createdAt = createdAt
	return r
}

func (r *Relationship) ToMap() map[string]any {
	return map[string]any{
		"id":                r.ID(),
		"entity_id":         r.EntityID(),
		"related_entity_id": r.RelatedEntityID(),
		"relationship_type": r.RelationshipType(),
		"metadata":          r.Metadata(),
		"created_at":        r.CreatedAt(),
	}
}
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
```

**Options struct:**

```go
type RelationshipOptions struct {
    EntityID         string
    RelatedEntityID  string
    RelationshipType string
    Metadata         map[string]string
}
```

### 3.3 Database Schema Addition

Uses **9-char short IDs** (varchar(9)) for space efficiency, matching cmsstore pattern:

**File:** `store_implementation.go` (extend SqlCreateTable)

```go
func (st *storeImplementation) SqlCreateTable() ([]string, error) {
    // ... existing entity and attribute tables ...
    
    sqlRelationship := `
    CREATE TABLE IF NOT EXISTS ` + st.entityTableName + `_relationships (
        id varchar(9) NOT NULL PRIMARY KEY,
        entity_id varchar(9) NOT NULL,
        related_entity_id varchar(9) NOT NULL,
        relationship_type varchar(50) NOT NULL,
        metadata text,
        created_at datetime NOT NULL,
        INDEX idx_entity (entity_id),
        INDEX idx_related (related_entity_id),
        INDEX idx_type (relationship_type),
        INDEX idx_entity_type (entity_id, relationship_type)
    );`
    
    sqlArray := append(existingSQL, sqlRelationship)
    return sqlArray, nil
}
```

### 3.4 Implementation Files

| File | Purpose | Lines (Est) |
|------|---------|-------------|
| `relationship.go` | Relationship type definition | 90 |
| `relationship_create.go` | Create relationship | 50 |
| `relationship_find.go` | Find relationship | 50 |
| `relationship_list.go` | List relationships | 60 |
| `relationship_delete.go` | Delete relationship | 50 |
| `id_helpers.go` | Short ID generation | 50 |
| `interfaces.go` | Extend StoreInterface | +30 |
| `store_implementation.go` | Update SqlCreateTable | +20 |
| **Total** | | **~400** |

---

**File:** `id_helpers.go` (new file)

Add short ID generation matching cmsstore pattern:

```go
package entitystore

import (
	"strings"
	"sync"
	"time"

	"github.com/dracory/uid"
)

var (
	idMutex    sync.Mutex
	lastIDTime int64
	idSequence int
)

// GenerateShortID generates a new shortened ID using TimestampMicro + Crockford Base32 (lowercase)
// Returns a 9-character lowercase ID (e.g., "86ccrtsgx")
// Thread-safe: Uses mutex to prevent duplicate IDs when called concurrently
func GenerateShortID() string {
	idMutex.Lock()
	defer idMutex.Unlock()

	// Get current microsecond timestamp
	now := time.Now().UnixMicro()

	// If same microsecond as last ID, add sequence number to ensure uniqueness
	if now == lastIDTime {
		idSequence++
		now += int64(idSequence)
	} else {
		lastIDTime = now
		idSequence = 0
	}

	timestampID := uid.TimestampMicro()
	shortened, _ := uid.ShortenCrockford(timestampID)
	return strings.ToLower(shortened)
}

// NormalizeID normalizes an ID to lowercase for consistent lookups
func NormalizeID(id string) string {
	return strings.ToLower(strings.TrimSpace(id))
}
```

### 4.1 Basic Relationship

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
    fmt.Println(book.GetString("title", ""))
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
    fmt.Printf("Friend: %s, Status: %s\n", friend.GetString("name", ""), status)
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

---

## 5. Implementation Phases

### Phase 1: Core Types (1 day)

1. Create `relationship.go` with type definition
2. Create `new_relationship.go` with constructor
3. Add constants for relationship types
4. Write basic type tests

### Phase 2: CRUD Operations (2 days)

1. Create `relationship_create.go`
2. Create `relationship_find.go`
3. Create `relationship_list.go`
4. Create `relationship_delete.go`
5. Write CRUD tests

### Phase 3: Integration (1 day)

1. Update `interfaces.go` with new methods
2. Update `store_implementation.go` with relationship table
3. Add table to AutoMigrate
4. Write integration tests

### Phase 4: Documentation (1 day)

1. Update README.md
2. Create examples/relationships.go
3. Write usage guide

**Total: ~5 days**

---

## 6. Testing Strategy

### 6.1 Unit Tests

| Test File | Coverage |
|-----------|----------|
| `relationship_test.go` | Type definition, getters, setters |
| `relationship_create_test.go` | Create relationships, duplicate handling |
| `relationship_find_test.go` | Find by entities |
| `relationship_list_test.go` | List by entity, list by related |
| `relationship_delete_test.go` | Delete single, delete all |

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

### Recommendation

**Proceed with implementation.**

Relationships are a fundamental data modeling concept. Adding native support to `entitystore` transforms it from a simple EAV store into a complete entity management system.

### Benefits

1. **Simplicity** - No more manual attribute-based linking
2. **Performance** - Single-query relationship loading
3. **Integrity** - Clear relationship semantics
4. **Reusability** - All entitystore consumers benefit

### Next Actions

1. Create feature branch
2. Implement Phase 1 (Core Types)
3. PR and review
4. Proceed to Phase 2-4

---

**End of Proposal**
