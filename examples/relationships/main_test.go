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
	parent, err := store.EntityCreateWithTypeAndAttributes(ctx, "parent", map[string]string{"name": "Parent"})
	if err != nil {
		t.Fatalf("Failed to create parent: %v", err)
	}
	if parent == nil {
		t.Fatal("Expected parent to be created")
	}

	child, err := store.EntityCreateWithTypeAndAttributes(ctx, "child", map[string]string{"name": "Child"})
	if err != nil {
		t.Fatalf("Failed to create child: %v", err)
	}
	if child == nil {
		t.Fatal("Expected child to be created")
	}

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

	parent, err := store.EntityCreateWithTypeAndAttributes(ctx, "parent", nil)
	if err != nil {
		t.Fatalf("Failed to create parent: %v", err)
	}
	if parent == nil {
		t.Fatal("Expected parent to be created")
	}

	child, err := store.EntityCreateWithTypeAndAttributes(ctx, "child", nil)
	if err != nil {
		t.Fatalf("Failed to create child: %v", err)
	}
	if child == nil {
		t.Fatal("Expected child to be created")
	}

	// Create relationship
	_, err = store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
		EntityID:         child.ID(),
		RelatedEntityID:  parent.ID(),
		RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
	})

	// Find relationship
	found, err := store.RelationshipFindByEntities(ctx, child.ID(), parent.ID(), entitystore.RELATIONSHIP_TYPE_BELONGS_TO)
	if err != nil {
		t.Fatalf("Failed to find relationship: %v", err)
	}
	if found == nil {
		t.Fatal("Expected relationship to be found")
	}

	if found.GetEntityID() != child.ID() {
		t.Errorf("Expected EntityID %s, got %s", child.ID(), found.GetEntityID())
	}
}

func TestRelationshipList(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	author, err := store.EntityCreateWithTypeAndAttributes(ctx, "author", nil)
	if err != nil {
		t.Fatalf("Failed to create author: %v", err)
	}
	if author == nil {
		t.Fatal("Expected author to be created")
	}

	book1, err := store.EntityCreateWithTypeAndAttributes(ctx, "book", nil)
	if err != nil {
		t.Fatalf("Failed to create book1: %v", err)
	}
	if book1 == nil {
		t.Fatal("Expected book1 to be created")
	}

	book2, err := store.EntityCreateWithTypeAndAttributes(ctx, "book", nil)
	if err != nil {
		t.Fatalf("Failed to create book2: %v", err)
	}
	if book2 == nil {
		t.Fatal("Expected book2 to be created")
	}

	// Create relationships
	_, err = store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
		EntityID:         book1.ID(),
		RelatedEntityID:  author.ID(),
		RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
	})
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}

	_, err = store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
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

	author, err := store.EntityCreateWithTypeAndAttributes(ctx, "author", nil)
	if err != nil {
		t.Fatalf("Failed to create author: %v", err)
	}
	if author == nil {
		t.Fatal("Expected author to be created")
	}

	book1, err := store.EntityCreateWithTypeAndAttributes(ctx, "book", nil)
	if err != nil {
		t.Fatalf("Failed to create book1: %v", err)
	}
	if book1 == nil {
		t.Fatal("Expected book1 to be created")
	}

	book2, err := store.EntityCreateWithTypeAndAttributes(ctx, "book", nil)
	if err != nil {
		t.Fatalf("Failed to create book2: %v", err)
	}
	if book2 == nil {
		t.Fatal("Expected book2 to be created")
	}

	_, err = store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
		EntityID:         book1.ID(),
		RelatedEntityID:  author.ID(),
		RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
	})
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}

	_, err = store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
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

	entity1, err := store.EntityCreateWithTypeAndAttributes(ctx, "test", nil)
	if err != nil {
		t.Fatalf("Failed to create entity1: %v", err)
	}
	if entity1 == nil {
		t.Fatal("Expected entity1 to be created")
	}

	entity2, err := store.EntityCreateWithTypeAndAttributes(ctx, "test", nil)
	if err != nil {
		t.Fatalf("Failed to create entity2: %v", err)
	}
	if entity2 == nil {
		t.Fatal("Expected entity2 to be created")
	}

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

	parent, err := store.EntityCreateWithTypeAndAttributes(ctx, "parent", nil)
	if err != nil {
		t.Fatalf("Failed to create parent: %v", err)
	}
	if parent == nil {
		t.Fatal("Expected parent to be created")
	}

	child, err := store.EntityCreateWithTypeAndAttributes(ctx, "child", nil)
	if err != nil {
		t.Fatalf("Failed to create child: %v", err)
	}
	if child == nil {
		t.Fatal("Expected child to be created")
	}

	rel, err := store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
		EntityID:         child.ID(),
		RelatedEntityID:  parent.ID(),
		RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
	})
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}
	if rel == nil {
		t.Fatal("Expected relationship to be created")
	}

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

	parent, err := store.EntityCreateWithTypeAndAttributes(ctx, "parent", nil)
	if err != nil {
		t.Fatalf("Failed to create parent: %v", err)
	}
	if parent == nil {
		t.Fatal("Expected parent to be created")
	}

	child1, err := store.EntityCreateWithTypeAndAttributes(ctx, "child", nil)
	if err != nil {
		t.Fatalf("Failed to create child1: %v", err)
	}
	if child1 == nil {
		t.Fatal("Expected child1 to be created")
	}

	child2, err := store.EntityCreateWithTypeAndAttributes(ctx, "child", nil)
	if err != nil {
		t.Fatalf("Failed to create child2: %v", err)
	}
	if child2 == nil {
		t.Fatal("Expected child2 to be created")
	}

	_, err = store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
		EntityID:         child1.ID(),
		RelatedEntityID:  parent.ID(),
		RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
	})
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}

	_, err = store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
		EntityID:         child2.ID(),
		RelatedEntityID:  parent.ID(),
		RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
	})
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}

	// Delete all relationships for parent
	err = store.RelationshipDeleteAll(ctx, parent.ID())
	if err != nil {
		t.Fatalf("Failed to delete all relationships: %v", err)
	}

	// Verify
	count, _ := store.RelationshipCount(ctx, entitystore.RelationshipQueryOptions{})
	if count != 0 {
		t.Errorf("Expected count 0 after delete all, got %d", count)
	}
}
