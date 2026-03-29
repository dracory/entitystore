# Entity Relationships Example

This example demonstrates how to create and manage relationships between entities using EntityStore's relationship system.

## What This Example Shows

### 1. Enabling Relationships
- Creating a store with `RelationshipsEnabled: true`
- This creates `entities_relationships` and `entities_relationships_trash` tables

### 2. Relationship Types
- **BELONGS_TO**: One entity belongs to another (e.g., Book belongs to Author)
- **HAS_MANY**: One entity has many related entities
- **MANY_TO_MANY**: Bidirectional linking (e.g., Book has multiple Categories)

### 3. Creating Relationships
- Linking two entities with a relationship type
- Optional metadata, sequence, and parent relationships

### 4. Querying Relationships
- Finding relationships by entity ID
- Finding relationships by related entity ID (reverse lookup)
- Finding specific relationships between two entities
- Counting relationships by type

### 5. Managing Relationships
- Soft deleting relationships (trash)
- Deleting all relationships for an entity

## Running the Example

```bash
go run examples/relationships/main.go
```

## Running Tests

```bash
go test ./examples/relationships/... -v
```

## Key Concepts

**Relationship Storage**: Relationships are stored in their own table with:
- `entity_id` - The source entity
- `related_entity_id` - The target entity
- `relationship_type` - BELONGS_TO, HAS_MANY, or MANY_MANY
- `parent_id` - For hierarchical relationships
- `sequence` - For ordering
- `metadata` - JSON text for additional data

## Code Highlights

```go
// Enable relationships in store
store, _ := entitystore.NewStore(entitystore.NewStoreOptions{
    RelationshipsEnabled: true,
    // ... other options
})

// Create BELONGS_TO relationship (book belongs to author)
rel, _ := store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
    EntityID:         book.ID(),
    RelatedEntityID:  author.ID(),
    RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
})

// Find all books belonging to an author
relationships, _ := store.RelationshipList(ctx, entitystore.RelationshipQueryOptions{
    RelatedEntityID:  author.ID(),
    RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
})

// Find specific relationship
found, _ := store.RelationshipFindByEntities(ctx, 
    book.ID(), 
    author.ID(), 
    entitystore.RELATIONSHIP_TYPE_BELONGS_TO)
```
