package entitystore

import (
	"context"
	"testing"

	_ "modernc.org/sqlite"
)

func TestRelationshipCreate(t *testing.T) {
	db := InitDB("relationship_create")
	defer db.Close()

	ctx := context.Background()
	store, err := NewStore(NewStoreOptions{
		DB:                    db,
		EntityTableName:       "entities",
		AttributeTableName:    "attributes",
		RelationshipsEnabled:  true,
		RelationshipTableName: "relationships",
		AutomigrateEnabled:    true,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Create two entities
	entity1, err := store.EntityCreateWithType(ctx, "book")
	if err != nil {
		t.Fatalf("Failed to create entity1: %v", err)
	}

	entity2, err := store.EntityCreateWithType(ctx, "author")
	if err != nil {
		t.Fatalf("Failed to create entity2: %v", err)
	}

	// Create relationship
	rel := NewRelationship()
	rel.SetEntityID(entity1.ID())
	rel.SetRelatedEntityID(entity2.ID())
	rel.SetRelationshipType(RELATIONSHIP_TYPE_BELONGS_TO)
	rel.SetSequence(1)
	rel.SetMetadata("{\"order\": 1}")

	err = store.RelationshipCreate(ctx, rel)
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}

	if rel.ID() == "" {
		t.Error("Relationship ID should be set after creation")
	}

	if rel.GetCreatedAt() == "" {
		t.Error("Relationship CreatedAt should be set after creation")
	}
}

func TestRelationshipCreateByOptions(t *testing.T) {
	db := InitDB("relationship_options")
	defer db.Close()

	ctx := context.Background()
	store, err := NewStore(NewStoreOptions{
		DB:                    db,
		EntityTableName:       "entities",
		AttributeTableName:    "attributes",
		RelationshipsEnabled:  true,
		RelationshipTableName: "relationships",
		AutomigrateEnabled:    true,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Create two entities
	entity1, err := store.EntityCreateWithType(ctx, "product")
	if err != nil {
		t.Fatalf("Failed to create entity1: %v", err)
	}

	entity2, err := store.EntityCreateWithType(ctx, "category")
	if err != nil {
		t.Fatalf("Failed to create entity2: %v", err)
	}

	// Create relationship using options
	rel, err := store.RelationshipCreateByOptions(ctx, RelationshipOptions{
		EntityID:         entity1.ID(),
		RelatedEntityID:  entity2.ID(),
		RelationshipType: RELATIONSHIP_TYPE_BELONGS_TO,
		Sequence:         5,
		Metadata:         "{\"priority\": \"high\"}",
	})
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}

	if rel.ID() == "" {
		t.Error("Relationship ID should be set")
	}

	if rel.GetEntityID() != entity1.ID() {
		t.Errorf("EntityID mismatch: expected %s, got %s", entity1.ID(), rel.GetEntityID())
	}

	if rel.GetRelatedEntityID() != entity2.ID() {
		t.Errorf("RelatedEntityID mismatch: expected %s, got %s", entity2.ID(), rel.GetRelatedEntityID())
	}

	if rel.GetSequence() != 5 {
		t.Errorf("Sequence mismatch: expected 5, got %d", rel.GetSequence())
	}
}

func TestRelationshipFind(t *testing.T) {
	db := InitDB("relationship_find")
	defer db.Close()

	ctx := context.Background()
	store, err := NewStore(NewStoreOptions{
		DB:                    db,
		EntityTableName:       "entities",
		AttributeTableName:    "attributes",
		RelationshipsEnabled:  true,
		RelationshipTableName: "relationships",
		AutomigrateEnabled:    true,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Create entities
	entity1, err := store.EntityCreateWithType(ctx, "post")
	if err != nil {
		t.Fatalf("Failed to create entity1: %v", err)
	}
	if entity1 == nil {
		t.Fatal("Expected entity1 to be created")
	}

	entity2, err := store.EntityCreateWithType(ctx, "user")
	if err != nil {
		t.Fatalf("Failed to create entity2: %v", err)
	}
	if entity2 == nil {
		t.Fatal("Expected entity2 to be created")
	}

	// Create relationship
	rel, err := store.RelationshipCreateByOptions(ctx, RelationshipOptions{
		EntityID:         entity1.ID(),
		RelatedEntityID:  entity2.ID(),
		RelationshipType: RELATIONSHIP_TYPE_BELONGS_TO,
	})
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}
	if rel == nil {
		t.Fatal("Expected relationship to be created")
	}

	// Find by ID
	found, err := store.RelationshipFind(ctx, rel.ID())
	if err != nil {
		t.Fatalf("Failed to find relationship: %v", err)
	}

	if found == nil {
		t.Fatal("Should have found the relationship")
	}

	if found.ID() != rel.ID() {
		t.Errorf("ID mismatch: expected %s, got %s", rel.ID(), found.ID())
	}

	// Find non-existent
	notFound, err := store.RelationshipFind(ctx, "non_existent_id")
	if err != nil {
		t.Fatalf("Error finding non-existent: %v", err)
	}

	if notFound != nil {
		t.Error("Should not find non-existent relationship")
	}
}

func TestRelationshipFindByEntities(t *testing.T) {
	db := InitDB("relationship_find_by_entities")
	defer db.Close()

	ctx := context.Background()
	store, err := NewStore(NewStoreOptions{
		DB:                    db,
		EntityTableName:       "entities",
		AttributeTableName:    "attributes",
		RelationshipsEnabled:  true,
		RelationshipTableName: "relationships",
		AutomigrateEnabled:    true,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Create entities
	entity1, err := store.EntityCreateWithType(ctx, "child")
	if err != nil {
		t.Fatalf("Failed to create entity1: %v", err)
	}
	if entity1 == nil {
		t.Fatal("Expected entity1 to be created")
	}

	entity2, err := store.EntityCreateWithType(ctx, "parent")
	if err != nil {
		t.Fatalf("Failed to create entity2: %v", err)
	}
	if entity2 == nil {
		t.Fatal("Expected entity2 to be created")
	}

	// Create relationship
	_, err = store.RelationshipCreateByOptions(ctx, RelationshipOptions{
		EntityID:         entity1.ID(),
		RelatedEntityID:  entity2.ID(),
		RelationshipType: RELATIONSHIP_TYPE_BELONGS_TO,
	})
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}

	// Find by entities
	found, err := store.RelationshipFindByEntities(ctx, entity1.ID(), entity2.ID(), RELATIONSHIP_TYPE_BELONGS_TO)
	if err != nil {
		t.Fatalf("Failed to find relationship: %v", err)
	}
	if found == nil {
		t.Fatal("Should have found the relationship")
	}

	// Try different type
	notFound, err := store.RelationshipFindByEntities(ctx, entity1.ID(), entity2.ID(), RELATIONSHIP_TYPE_HAS_MANY)
	if err != nil {
		t.Fatalf("Failed to find relationship: %v", err)
	}
	if notFound != nil {
		t.Error("Should not find relationship with wrong type")
	}
}

func TestRelationshipList(t *testing.T) {
	db := InitDB("relationship_list")
	defer db.Close()

	ctx := context.Background()
	store, err := NewStore(NewStoreOptions{
		DB:                    db,
		EntityTableName:       "entities",
		AttributeTableName:    "attributes",
		RelationshipsEnabled:  true,
		RelationshipTableName: "relationships",
		AutomigrateEnabled:    true,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Create entities
	author, err := store.EntityCreateWithType(ctx, "author")
	if err != nil {
		t.Fatalf("Failed to create author: %v", err)
	}
	if author == nil {
		t.Fatal("Expected author to be created")
	}

	book1, err := store.EntityCreateWithType(ctx, "book")
	if err != nil {
		t.Fatalf("Failed to create book1: %v", err)
	}
	if book1 == nil {
		t.Fatal("Expected book1 to be created")
	}

	book2, err := store.EntityCreateWithType(ctx, "book")
	if err != nil {
		t.Fatalf("Failed to create book2: %v", err)
	}
	if book2 == nil {
		t.Fatal("Expected book2 to be created")
	}

	book3, err := store.EntityCreateWithType(ctx, "book")
	if err != nil {
		t.Fatalf("Failed to create book3: %v", err)
	}
	if book3 == nil {
		t.Fatal("Expected book3 to be created")
	}

	// Create relationships
	_, err = store.RelationshipCreateByOptions(ctx, RelationshipOptions{
		EntityID:         book1.ID(),
		RelatedEntityID:  author.ID(),
		RelationshipType: RELATIONSHIP_TYPE_BELONGS_TO,
		Sequence:         1,
	})
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}

	_, err = store.RelationshipCreateByOptions(ctx, RelationshipOptions{
		EntityID:         book2.ID(),
		RelatedEntityID:  author.ID(),
		RelationshipType: RELATIONSHIP_TYPE_BELONGS_TO,
		Sequence:         2,
	})
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}

	_, err = store.RelationshipCreateByOptions(ctx, RelationshipOptions{
		EntityID:         book3.ID(),
		RelatedEntityID:  author.ID(),
		RelationshipType: RELATIONSHIP_TYPE_BELONGS_TO,
		Sequence:         3,
	})
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}

	// List all for author (as related entity)
	rels, err := store.RelationshipListRelated(ctx, author.ID(), RELATIONSHIP_TYPE_BELONGS_TO)
	if err != nil {
		t.Fatalf("Failed to list relationships: %v", err)
	}

	if len(rels) != 3 {
		t.Errorf("Expected 3 relationships, got %d", len(rels))
	}

	// List with pagination
	paginated, err := store.RelationshipList(ctx, RelationshipQueryOptions{
		RelatedEntityID:  author.ID(),
		RelationshipType: RELATIONSHIP_TYPE_BELONGS_TO,
		Limit:            2,
	})
	if err != nil {
		t.Fatalf("Failed to list paginated: %v", err)
	}

	if len(paginated) != 2 {
		t.Errorf("Expected 2 paginated relationships, got %d", len(paginated))
	}
}

func TestRelationshipDelete(t *testing.T) {
	db := InitDB("relationship_delete")
	defer db.Close()

	ctx := context.Background()
	store, err := NewStore(NewStoreOptions{
		DB:                    db,
		EntityTableName:       "entities",
		AttributeTableName:    "attributes",
		RelationshipsEnabled:  true,
		RelationshipTableName: "relationships",
		AutomigrateEnabled:    true,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Create entities and relationship
	entity1, err := store.EntityCreateWithType(ctx, "item")
	if err != nil {
		t.Fatalf("Failed to create entity1: %v", err)
	}
	if entity1 == nil {
		t.Fatal("Expected entity1 to be created")
	}

	entity2, err := store.EntityCreateWithType(ctx, "category")
	if err != nil {
		t.Fatalf("Failed to create entity2: %v", err)
	}
	if entity2 == nil {
		t.Fatal("Expected entity2 to be created")
	}

	rel, err := store.RelationshipCreateByOptions(ctx, RelationshipOptions{
		EntityID:         entity1.ID(),
		RelatedEntityID:  entity2.ID(),
		RelationshipType: RELATIONSHIP_TYPE_BELONGS_TO,
	})
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}
	if rel == nil {
		t.Fatal("Expected relationship to be created")
	}

	// Delete
	deleted, err := store.RelationshipDelete(ctx, rel.ID())
	if err != nil {
		t.Fatalf("Failed to delete relationship: %v", err)
	}

	if !deleted {
		t.Error("Expected deleted to be true")
	}

	// Verify it's gone
	found, err := store.RelationshipFind(ctx, rel.ID())
	if err != nil {
		t.Fatalf("Failed to find relationship: %v", err)
	}
	if found != nil {
		t.Error("Relationship should be deleted")
	}

	// Delete non-existent
	notDeleted, _ := store.RelationshipDelete(ctx, "non_existent")
	if notDeleted {
		t.Error("Deleting non-existent should return false")
	}
}

func TestRelationshipDeleteAll(t *testing.T) {
	db := InitDB("relationship_delete_all")
	defer db.Close()

	ctx := context.Background()
	store, err := NewStore(NewStoreOptions{
		DB:                    db,
		EntityTableName:       "entities",
		AttributeTableName:    "attributes",
		RelationshipsEnabled:  true,
		RelationshipTableName: "relationships",
		AutomigrateEnabled:    true,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Create entity with multiple relationships
	entity1, err := store.EntityCreateWithType(ctx, "parent")
	if err != nil {
		t.Fatalf("Failed to create entity1: %v", err)
	}
	if entity1 == nil {
		t.Fatal("Expected entity1 to be created")
	}

	child1, err := store.EntityCreateWithType(ctx, "child")
	if err != nil {
		t.Fatalf("Failed to create child1: %v", err)
	}
	if child1 == nil {
		t.Fatal("Expected child1 to be created")
	}

	child2, err := store.EntityCreateWithType(ctx, "child")
	if err != nil {
		t.Fatalf("Failed to create child2: %v", err)
	}
	if child2 == nil {
		t.Fatal("Expected child2 to be created")
	}

	child3, err := store.EntityCreateWithType(ctx, "child")
	if err != nil {
		t.Fatalf("Failed to create child3: %v", err)
	}
	if child3 == nil {
		t.Fatal("Expected child3 to be created")
	}

	_, err = store.RelationshipCreateByOptions(ctx, RelationshipOptions{
		EntityID:         child1.ID(),
		RelatedEntityID:  entity1.ID(),
		RelationshipType: RELATIONSHIP_TYPE_BELONGS_TO,
	})
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}

	_, err = store.RelationshipCreateByOptions(ctx, RelationshipOptions{
		EntityID:         child2.ID(),
		RelatedEntityID:  entity1.ID(),
		RelationshipType: RELATIONSHIP_TYPE_BELONGS_TO,
	})
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}

	_, err = store.RelationshipCreateByOptions(ctx, RelationshipOptions{
		EntityID:         child3.ID(),
		RelatedEntityID:  entity1.ID(),
		RelationshipType: RELATIONSHIP_TYPE_BELONGS_TO,
	})
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}

	// Count before delete
	before, err := store.RelationshipList(ctx, RelationshipQueryOptions{EntityID: child1.ID()})
	if err != nil {
		t.Fatalf("Failed to list relationships: %v", err)
	}
	if len(before) == 0 {
		t.Fatal("Should have relationships before delete")
	}

	// Delete all for entity1
	err = store.RelationshipDeleteAll(ctx, entity1.ID())
	if err != nil {
		t.Fatalf("Failed to delete all: %v", err)
	}

	// Verify they're gone
	after, err := store.RelationshipListRelated(ctx, entity1.ID(), RELATIONSHIP_TYPE_BELONGS_TO)
	if err != nil {
		t.Fatalf("Failed to list relationships: %v", err)
	}
	if len(after) != 0 {
		t.Errorf("Expected 0 relationships after delete all, got %d", len(after))
	}
}

func TestRelationshipTrashAndRestore(t *testing.T) {
	db := InitDB("relationship_trash_restore")
	defer db.Close()

	ctx := context.Background()
	store, err := NewStore(NewStoreOptions{
		DB:                         db,
		EntityTableName:            "entities",
		AttributeTableName:         "attributes",
		RelationshipsEnabled:       true,
		RelationshipTableName:      "relationships",
		RelationshipTrashTableName: "relationships_trash",
		AutomigrateEnabled:         true,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Create entities and relationship
	entity1, err := store.EntityCreateWithType(ctx, "post")
	if err != nil {
		t.Fatalf("Failed to create entity1: %v", err)
	}
	if entity1 == nil {
		t.Fatal("Expected entity1 to be created")
	}

	entity2, err := store.EntityCreateWithType(ctx, "user")
	if err != nil {
		t.Fatalf("Failed to create entity2: %v", err)
	}
	if entity2 == nil {
		t.Fatal("Expected entity2 to be created")
	}

	rel, err := store.RelationshipCreateByOptions(ctx, RelationshipOptions{
		EntityID:         entity1.ID(),
		RelatedEntityID:  entity2.ID(),
		RelationshipType: RELATIONSHIP_TYPE_BELONGS_TO,
		Metadata:         "{\"test\": \"data\"}",
	})
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}
	if rel == nil {
		t.Fatal("Expected relationship to be created")
	}

	// Trash it
	trashed, err := store.RelationshipTrash(ctx, rel.ID(), "user_123")
	if err != nil {
		t.Fatalf("Failed to trash relationship: %v", err)
	}

	if !trashed {
		t.Error("Expected trashed to be true")
	}

	// Verify it's in main table
	found, err := store.RelationshipFind(ctx, rel.ID())
	if err != nil {
		t.Fatalf("Failed to find relationship: %v", err)
	}
	if found != nil {
		t.Error("Relationship should not be in main table after trash")
	}

	// Verify it's in trash
	trashItems, err := store.RelationshipTrashList(ctx, RelationshipQueryOptions{
		ID:    rel.ID(),
		Limit: 1,
	})
	if err != nil {
		t.Fatalf("Failed to list trash: %v", err)
	}

	if len(trashItems) != 1 {
		t.Errorf("Expected 1 item in trash, got %d", len(trashItems))
	}

	if len(trashItems) > 0 {
		if trashItems[0].GetDeletedBy() != "user_123" {
			t.Errorf("Expected deleted_by to be 'user_123', got '%s'", trashItems[0].GetDeletedBy())
		}

		if trashItems[0].GetMetadata() != "{\"test\": \"data\"}" {
			t.Errorf("Metadata should be preserved in trash, got '%s'", trashItems[0].GetMetadata())
		}
	}

	// Restore it
	restored, err := store.RelationshipRestore(ctx, rel.ID())
	if err != nil {
		t.Fatalf("Failed to restore relationship: %v", err)
	}

	if !restored {
		t.Error("Expected restored to be true")
	}

	// Verify it's back in main table
	found2, err := store.RelationshipFind(ctx, rel.ID())
	if err != nil {
		t.Fatalf("Failed to find relationship: %v", err)
	}
	if found2 == nil {
		t.Fatal("Relationship should be restored to main table")
	}

	if found2.GetMetadata() != "{\"test\": \"data\"}" {
		t.Errorf("Metadata should be preserved after restore, got '%s'", found2.GetMetadata())
	}

	// Verify it's gone from trash
	trashItems2, err := store.RelationshipTrashList(ctx, RelationshipQueryOptions{
		ID:    rel.ID(),
		Limit: 1,
	})
	if err != nil {
		t.Fatalf("Failed to list trash: %v", err)
	}

	if len(trashItems2) != 0 {
		t.Errorf("Expected 0 items in trash after restore, got %d", len(trashItems2))
	}
}

func TestRelationshipDuplicatePrevention(t *testing.T) {
	db := InitDB("relationship_duplicate_prevention")
	defer db.Close()

	ctx := context.Background()
	store, err := NewStore(NewStoreOptions{
		DB:                    db,
		EntityTableName:       "entities",
		AttributeTableName:    "attributes",
		RelationshipsEnabled:  true,
		RelationshipTableName: "relationships",
		AutomigrateEnabled:    true,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Create entities
	entity1, err := store.EntityCreateWithType(ctx, "book")
	if err != nil {
		t.Fatalf("Failed to create entity1: %v", err)
	}
	if entity1 == nil {
		t.Fatal("Expected entity1 to be created")
	}

	entity2, err := store.EntityCreateWithType(ctx, "author")
	if err != nil {
		t.Fatalf("Failed to create entity2: %v", err)
	}
	if entity2 == nil {
		t.Fatal("Expected entity2 to be created")
	}

	// Create first relationship
	rel1, err := store.RelationshipCreateByOptions(ctx, RelationshipOptions{
		EntityID:         entity1.ID(),
		RelatedEntityID:  entity2.ID(),
		RelationshipType: RELATIONSHIP_TYPE_BELONGS_TO,
	})
	if err != nil {
		t.Fatalf("Failed to create first relationship: %v", err)
	}

	if rel1 == nil {
		t.Fatal("First relationship should be created")
	}

	// Try to create duplicate relationship
	rel2, err := store.RelationshipCreateByOptions(ctx, RelationshipOptions{
		EntityID:         entity1.ID(),
		RelatedEntityID:  entity2.ID(),
		RelationshipType: RELATIONSHIP_TYPE_BELONGS_TO,
	})

	if err == nil {
		t.Error("Expected error when creating duplicate relationship")
	}

	if rel2 != nil {
		t.Error("Duplicate relationship should not be created")
	}

	if err != nil && err.Error() != "relationship already exists" {
		t.Errorf("Expected 'relationship already exists' error, got: %v", err)
	}

	// Verify only one relationship exists
	count, _ := store.RelationshipCount(ctx, RelationshipQueryOptions{
		EntityID:         entity1.ID(),
		RelatedEntityID:  entity2.ID(),
		RelationshipType: RELATIONSHIP_TYPE_BELONGS_TO,
	})

	if count != 1 {
		t.Errorf("Expected 1 relationship, got %d", count)
	}
}

func TestRelationshipValidation(t *testing.T) {
	db := InitDB("relationship_validation")
	defer db.Close()

	ctx := context.Background()
	store, err := NewStore(NewStoreOptions{
		DB:                    db,
		EntityTableName:       "entities",
		AttributeTableName:    "attributes",
		RelationshipsEnabled:  true,
		RelationshipTableName: "relationships",
		AutomigrateEnabled:    true,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Test nil relationship
	err = store.RelationshipCreate(ctx, nil)
	if err == nil || err.Error() != "relationship cannot be nil" {
		t.Errorf("Expected 'relationship cannot be nil' error, got: %v", err)
	}

	// Test missing entity_id
	rel := NewRelationship()
	rel.SetRelatedEntityID("related_123")
	rel.SetRelationshipType(RELATIONSHIP_TYPE_BELONGS_TO)
	err = store.RelationshipCreate(ctx, rel)
	if err == nil || err.Error() != "entity_id is required" {
		t.Errorf("Expected 'entity_id is required' error, got: %v", err)
	}

	// Test missing related_entity_id
	rel2 := NewRelationship()
	rel2.SetEntityID("entity_123")
	rel2.SetRelationshipType(RELATIONSHIP_TYPE_BELONGS_TO)
	err = store.RelationshipCreate(ctx, rel2)
	if err == nil || err.Error() != "related_entity_id is required" {
		t.Errorf("Expected 'related_entity_id is required' error, got: %v", err)
	}

	// Test missing relationship_type
	rel3 := NewRelationship()
	rel3.SetEntityID("entity_123")
	rel3.SetRelatedEntityID("related_123")
	err = store.RelationshipCreate(ctx, rel3)
	if err == nil || err.Error() != "relationship_type is required" {
		t.Errorf("Expected 'relationship_type is required' error, got: %v", err)
	}

	// Test self-referencing belongs_to (should fail)
	rel4 := NewRelationship()
	rel4.SetEntityID("same_id")
	rel4.SetRelatedEntityID("same_id")
	rel4.SetRelationshipType(RELATIONSHIP_TYPE_BELONGS_TO)
	err = store.RelationshipCreate(ctx, rel4)
	if err == nil {
		t.Error("Expected error for self-referencing belongs_to relationship")
	}

	// Test self-referencing many_to_many (should succeed)
	rel5 := NewRelationship()
	rel5.SetEntityID("same_id2")
	rel5.SetRelatedEntityID("same_id2")
	rel5.SetRelationshipType(RELATIONSHIP_TYPE_MANY_MANY)
	err = store.RelationshipCreate(ctx, rel5)
	if err != nil {
		t.Errorf("Self-referencing many_to_many should be allowed, got error: %v", err)
	}
}
