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
		DB:                 db,
		EntityTableName:    "test_entities",
		AttributeTableName: "test_attributes",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return store, db, cleanup
}

func TestEntityCreateWithTypeAndAttributes(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	entity, err := store.EntityCreateWithTypeAndAttributes(ctx, "test", map[string]string{
		"name":  "Test Entity",
		"value": "123",
	})
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	if entity.ID() == "" {
		t.Error("Entity ID should not be empty")
	}

	if entity.GetType() != "test" {
		t.Errorf("Expected entity type 'test', got '%s'", entity.GetType())
	}

	// Verify attributes via store
	attr, err := store.AttributeFind(ctx, entity.ID(), "name")
	if err != nil {
		t.Fatalf("Failed to find attribute: %v", err)
	}
	if attr == nil {
		t.Fatal("Expected attribute to be found")
	}

	if attr.GetValue() != "Test Entity" {
		t.Errorf("Expected name 'Test Entity', got '%s'", attr.GetValue())
	}
}

func TestEntityFindByID(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	// Create entity
	created, err := store.EntityCreateWithTypeAndAttributes(ctx, "person", map[string]string{
		"name": "Alice",
		"age":  "25",
	})
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}
	if created == nil {
		t.Fatal("Expected entity to be created")
	}

	// Find by ID
	found, err := store.EntityFindByID(ctx, created.ID())
	if err != nil {
		t.Fatalf("Failed to find entity: %v", err)
	}
	if found == nil {
		t.Fatal("Expected entity to be found")
	}

	if found.ID() != created.ID() {
		t.Errorf("Expected ID %s, got %s", created.ID(), found.ID())
	}

	// Verify attributes via store
	attr, err := store.AttributeFind(ctx, found.ID(), "name")
	if err != nil {
		t.Fatalf("Failed to find attribute: %v", err)
	}
	if attr == nil {
		t.Fatal("Expected attribute to be found")
	}

	if attr.GetValue() != "Alice" {
		t.Errorf("Expected name 'Alice', got '%s'", attr.GetValue())
	}
}

func TestEntityList(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	// Create multiple entities
	_, err := store.EntityCreateWithTypeAndAttributes(ctx, "type_a", map[string]string{"key": "1"})
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	_, err = store.EntityCreateWithTypeAndAttributes(ctx, "type_a", map[string]string{"key": "2"})
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	_, err = store.EntityCreateWithTypeAndAttributes(ctx, "type_b", map[string]string{"key": "3"})
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	// List all
	all, err := store.EntityList(ctx, entitystore.EntityQueryOptions{})
	if err != nil {
		t.Fatalf("Failed to list entities: %v", err)
	}

	if len(all) != 3 {
		t.Errorf("Expected 3 entities, got %d", len(all))
	}

	// List by type
	typeA, _ := store.EntityList(ctx, entitystore.EntityQueryOptions{EntityType: "type_a"})
	if len(typeA) != 2 {
		t.Errorf("Expected 2 type_a entities, got %d", len(typeA))
	}
}

func TestEntityCount(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	// Create entities
	_, err := store.EntityCreateWithTypeAndAttributes(ctx, "person", map[string]string{})
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	_, err = store.EntityCreateWithTypeAndAttributes(ctx, "person", map[string]string{})
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	_, err = store.EntityCreateWithTypeAndAttributes(ctx, "product", map[string]string{})
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}

	// Count all
	total, _ := store.EntityCount(ctx, entitystore.EntityQueryOptions{})
	if total != 3 {
		t.Errorf("Expected count 3, got %d", total)
	}

	// Count by type
	personCount, _ := store.EntityCount(ctx, entitystore.EntityQueryOptions{EntityType: "person"})
	if personCount != 2 {
		t.Errorf("Expected person count 2, got %d", personCount)
	}
}

func TestEntityUpdate(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	// Create entity
	entity, err := store.EntityCreateWithTypeAndAttributes(ctx, "test", map[string]string{
		"name": "Original",
	})
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}
	if entity == nil {
		t.Fatal("Expected entity to be created")
	}

	// Update attribute via store
	err = store.AttributeSetString(ctx, entity.ID(), "name", "Updated")
	if err != nil {
		t.Fatalf("Failed to update attribute: %v", err)
	}

	// Verify
	attr, err := store.AttributeFind(ctx, entity.ID(), "name")
	if err != nil {
		t.Fatalf("Failed to find attribute: %v", err)
	}
	if attr == nil {
		t.Fatal("Expected attribute to be found")
	}
	if attr.GetValue() != "Updated" {
		t.Errorf("Expected name 'Updated', got '%s'", attr.GetValue())
	}
}

func TestEntityTrash(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	// Create and trash entity
	entity, err := store.EntityCreateWithTypeAndAttributes(ctx, "test", map[string]string{})
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}
	if entity == nil {
		t.Fatal("Expected entity to be created")
	}
	id := entity.ID()

	deleted, err := store.EntityTrash(ctx, id)
	if err != nil {
		t.Fatalf("Failed to trash entity: %v", err)
	}

	if !deleted {
		t.Error("Expected entity to be deleted")
	}

	// Verify count decreased
	count, _ := store.EntityCount(ctx, entitystore.EntityQueryOptions{})
	if count != 0 {
		t.Errorf("Expected count 0 after trash, got %d", count)
	}
}

func TestAttributeTypes(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	entity, err := store.EntityCreateWithTypeAndAttributes(ctx, "test", map[string]string{
		"string_val": "hello",
		"int_val":    "42",
		"float_val":  "3.14",
	})
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}
	if entity == nil {
		t.Fatal("Expected entity to be created")
	}

	// Test string attribute
	attrs, err := store.EntityAttributeList(ctx, entity.ID())
	if err != nil {
		t.Fatalf("Failed to list attributes: %v", err)
	}
	if len(attrs) != 3 {
		t.Fatalf("Expected 3 attributes, got %d", len(attrs))
	}

	// Find int attribute and verify conversion
	for _, attr := range attrs {
		if attr.GetKey() == "int_val" {
			intVal, err := attr.GetInt()
			if err != nil {
				t.Errorf("Failed to convert to int: %v", err)
			}
			if intVal != 42 {
				t.Errorf("Expected int 42, got %d", intVal)
			}
		}

		if attr.GetKey() == "float_val" {
			floatVal, err := attr.GetFloat()
			if err != nil {
				t.Errorf("Failed to convert to float: %v", err)
			}
			if floatVal != 3.14 {
				t.Errorf("Expected float 3.14, got %f", floatVal)
			}
		}
	}
}
