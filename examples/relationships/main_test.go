package main

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dracory/entitystore"
	_ "modernc.org/sqlite"
)

func setupTestStore(t *testing.T) (entitystore.StoreInterface, *sql.DB, func()) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	store, err := entitystore.NewStore(entitystore.NewStoreOptions{
		DB:                   db,
		EntityTableName:      "test_entities",
		AttributeTableName:   "test_attributes",
		RelationshipsEnabled: true,
		AutomigrateEnabled:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return store, db, cleanup
}

func TestRelationshipCreate(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	// Create two entities
	parent, _ := store.EntityCreateWithTypeAndAttributes(ctx, "parent", map[string]string{"name": "Parent"})
	child, _ := store.EntityCreateWithTypeAndAttributes(ctx, "child", map[string]string{"name": "Child"})

	// Create relationship
	rel, err := store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
		EntityID:         child.ID(),
		RelatedEntityID:  parent.ID(),
		RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
	})
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}

	if rel.ID() == "" {
		t.Error("Relationship ID should not be empty")
	}

	if rel.GetEntityID() != child.ID() {
		t.Errorf("Expected EntityID %s, got %s", child.ID(), rel.GetEntityID())
	}

	if rel.GetRelatedEntityID() != parent.ID() {
		t.Errorf("Expected RelatedEntityID %s, got %s", parent.ID(), rel.GetRelatedEntityID())
	}
}

func TestRelationshipFindByEntities(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	parent, _ := store.EntityCreateWithTypeAndAttributes(ctx, "parent", nil)
	child, _ := store.EntityCreateWithTypeAndAttributes(ctx, "child", nil)

	// Create relationship
	store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
		EntityID:         child.ID(),
		RelatedEntityID:  parent.ID(),
		RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
	})

	// Find relationship
	found, err := store.RelationshipFindByEntities(ctx, child.ID(), parent.ID(), entitystore.RELATIONSHIP_TYPE_BELONGS_TO)
	if err != nil {
		t.Fatalf("Failed to find relationship: %v", err)
	}

	if found.GetEntityID() != child.ID() {
		t.Errorf("Expected EntityID %s, got %s", child.ID(), found.GetEntityID())
	}
}

func TestRelationshipList(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	author, _ := store.EntityCreateWithTypeAndAttributes(ctx, "author", nil)
	book1, _ := store.EntityCreateWithTypeAndAttributes(ctx, "book", nil)
	book2, _ := store.EntityCreateWithTypeAndAttributes(ctx, "book", nil)

	// Create relationships
	store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
		EntityID:         book1.ID(),
		RelatedEntityID:  author.ID(),
		RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
	})
	store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
		EntityID:         book2.ID(),
		RelatedEntityID:  author.ID(),
		RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
	})

	// List relationships where author is the related entity
	rels, err := store.RelationshipList(ctx, entitystore.RelationshipQueryOptions{
		RelatedEntityID:  author.ID(),
		RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
	})
	if err != nil {
		t.Fatalf("Failed to list relationships: %v", err)
	}

	if len(rels) != 2 {
		t.Errorf("Expected 2 relationships, got %d", len(rels))
	}
}

func TestRelationshipCount(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	author, _ := store.EntityCreateWithTypeAndAttributes(ctx, "author", nil)
	book1, _ := store.EntityCreateWithTypeAndAttributes(ctx, "book", nil)
	book2, _ := store.EntityCreateWithTypeAndAttributes(ctx, "book", nil)

	store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
		EntityID:         book1.ID(),
		RelatedEntityID:  author.ID(),
		RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
	})
	store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
		EntityID:         book2.ID(),
		RelatedEntityID:  author.ID(),
		RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
	})

	// Count all
	total, _ := store.RelationshipCount(ctx, entitystore.RelationshipQueryOptions{})
	if total != 2 {
		t.Errorf("Expected count 2, got %d", total)
	}

	// Count by type
	belongsCount, _ := store.RelationshipCount(ctx, entitystore.RelationshipQueryOptions{
		RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
	})
	if belongsCount != 2 {
		t.Errorf("Expected belongs_to count 2, got %d", belongsCount)
	}
}

func TestRelationshipTypes(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	entity1, _ := store.EntityCreateWithTypeAndAttributes(ctx, "test", nil)
	entity2, _ := store.EntityCreateWithTypeAndAttributes(ctx, "test", nil)

	tests := []struct {
		relType string
	}{
		{entitystore.RELATIONSHIP_TYPE_BELONGS_TO},
		{entitystore.RELATIONSHIP_TYPE_HAS_MANY},
		{entitystore.RELATIONSHIP_TYPE_MANY_MANY},
	}

	for _, tt := range tests {
		rel, err := store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
			EntityID:         entity1.ID(),
			RelatedEntityID:  entity2.ID(),
			RelationshipType: tt.relType,
		})
		if err != nil {
			t.Errorf("Failed to create %s relationship: %v", tt.relType, err)
			continue
		}
		if rel.GetRelationshipType() != tt.relType {
			t.Errorf("Expected type %s, got %s", tt.relType, rel.GetRelationshipType())
		}
	}
}

func TestRelationshipTrash(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	parent, _ := store.EntityCreateWithTypeAndAttributes(ctx, "parent", nil)
	child, _ := store.EntityCreateWithTypeAndAttributes(ctx, "child", nil)

	rel, _ := store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
		EntityID:         child.ID(),
		RelatedEntityID:  parent.ID(),
		RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
	})

	// Trash relationship
	deleted, err := store.RelationshipTrash(ctx, rel.ID(), "test_user")
	if err != nil {
		t.Fatalf("Failed to trash relationship: %v", err)
	}
	if !deleted {
		t.Error("Expected relationship to be deleted")
	}

	// Count should be 0
	count, _ := store.RelationshipCount(ctx, entitystore.RelationshipQueryOptions{})
	if count != 0 {
		t.Errorf("Expected count 0 after trash, got %d", count)
	}
}

func TestRelationshipDeleteAll(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	parent, _ := store.EntityCreateWithTypeAndAttributes(ctx, "parent", nil)
	child1, _ := store.EntityCreateWithTypeAndAttributes(ctx, "child", nil)
	child2, _ := store.EntityCreateWithTypeAndAttributes(ctx, "child", nil)

	store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
		EntityID:         child1.ID(),
		RelatedEntityID:  parent.ID(),
		RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
	})
	store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
		EntityID:         child2.ID(),
		RelatedEntityID:  parent.ID(),
		RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
	})

	// Delete all relationships for parent
	err := store.RelationshipDeleteAll(ctx, parent.ID())
	if err != nil {
		t.Fatalf("Failed to delete all relationships: %v", err)
	}

	// Verify
	count, _ := store.RelationshipCount(ctx, entitystore.RelationshipQueryOptions{})
	if count != 0 {
		t.Errorf("Expected count 0 after delete all, got %d", count)
	}
}
