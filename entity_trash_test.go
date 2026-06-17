package entitystore

import (
	"context"
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestEntityTrash(t *testing.T) {
	db := InitDB("test_entity_trash.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Must be NIL:", err)
	}

	entity, err := store.EntityCreateWithTypeAndAttributes(context.Background(), "post", map[string]string{
		"title": "Test Post Title",
		"text":  "Test Post Text",
	})

	if err != nil {
		t.Fatal("Entity could not be created:", err)
	}

	if entity == nil {
		t.Fatal("Entity could not be created")
	}

	attr, err := store.AttributeFind(context.Background(), entity.ID(), "title")

	if err != nil {
		t.Fatal("Attribute could not be found:", err)
	}

	if attr == nil {
		t.Fatal("Attribute should not be nil")
	}

	isDeleted, err := store.EntityTrash(context.Background(), entity.ID())

	if err != nil {
		t.Fatal("Entity could not be deleted:", err)
	}

	if isDeleted == false {
		t.Fatal("Entity could not be soft deleted")
	}

	val, err := store.EntityFindByID(context.Background(), entity.ID())

	if err != nil {
		t.Fatal(err)
	}

	if val != nil {
		t.Fatal("Entity should no longer be present")
	}

	attr, err = store.AttributeFind(context.Background(), entity.ID(), "title")

	if err != nil {
		t.Fatal("Attribute could not be found:", err)
	}

	if attr != nil {
		t.Fatal("Attribute should be nil")
	}
}

func TestEntityTrashImplementation(t *testing.T) {
	trash := NewEntityTrash()

	// Test ID generation
	if trash.ID() == "" {
		t.Error("expected ID to be generated")
	}
	if len(trash.ID()) < 9 {
		t.Errorf("expected ID length >= 9, got %d", len(trash.ID()))
	}

	// Test GetID() method consistency
	if trash.GetID() != trash.ID() {
		t.Errorf("expected GetID() '%s' to match ID() '%s'", trash.GetID(), trash.ID())
	}

	// Test EntityType getter/setter
	trash.SetType("product")
	if trash.GetType() != "product" {
		t.Errorf("expected EntityType 'product', got '%s'", trash.GetType())
	}

	// Test EntityHandle getter/setter
	trash.SetHandle("my-handle")
	if trash.GetHandle() != "my-handle" {
		t.Errorf("expected EntityHandle 'my-handle', got '%s'", trash.GetHandle())
	}

	// Test DeletedAt
	if trash.GetDeletedAt() == "" {
		t.Error("expected DeletedAt to be set")
	}

	// Test DeletedBy getter/setter
	trash.SetDeletedBy("user123")
	if trash.GetDeletedBy() != "user123" {
		t.Errorf("expected DeletedBy 'user123', got '%s'", trash.GetDeletedBy())
	}
}

func TestEntityTrashFromExistingData(t *testing.T) {
	data := map[string]string{
		COLUMN_ID:            "trash123",
		COLUMN_ENTITY_TYPE:   "product",
		COLUMN_ENTITY_HANDLE: "deleted-product",
		COLUMN_CREATED_AT:    carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC),
		COLUMN_UPDATED_AT:    carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC),
		COLUMN_DELETED_AT:    carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC),
		COLUMN_DELETED_BY:    "admin",
	}

	trash := NewEntityTrashFromExistingData(data)

	if trash.ID() != "trash123" {
		t.Errorf("expected ID 'trash123', got '%s'", trash.ID())
	}

	if trash.GetType() != "product" {
		t.Errorf("expected EntityType 'product', got '%s'", trash.GetType())
	}

	if trash.GetDeletedBy() != "admin" {
		t.Errorf("expected DeletedBy 'admin', got '%s'", trash.GetDeletedBy())
	}
}
