package entitystore

import (
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestEntityTrashImplementation(t *testing.T) {
	trash := NewEntityTrash()

	// Test ID generation
	if trash.ID() == "" {
		t.Error("expected ID to be generated")
	}
	if len(trash.ID()) < 9 {
		t.Errorf("expected ID length >= 9, got %d", len(trash.ID()))
	}

	// Test EntityType getter/setter
	trash.SetEntityType("product")
	if trash.EntityType() != "product" {
		t.Errorf("expected EntityType 'product', got '%s'", trash.EntityType())
	}

	// Test EntityHandle getter/setter
	trash.SetEntityHandle("my-handle")
	if trash.EntityHandle() != "my-handle" {
		t.Errorf("expected EntityHandle 'my-handle', got '%s'", trash.EntityHandle())
	}

	// Test DeletedAt
	if trash.DeletedAt() == "" {
		t.Error("expected DeletedAt to be set")
	}

	// Test DeletedBy getter/setter
	trash.SetDeletedBy("user123")
	if trash.DeletedBy() != "user123" {
		t.Errorf("expected DeletedBy 'user123', got '%s'", trash.DeletedBy())
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

	if trash.EntityType() != "product" {
		t.Errorf("expected EntityType 'product', got '%s'", trash.EntityType())
	}

	if trash.DeletedBy() != "admin" {
		t.Errorf("expected DeletedBy 'admin', got '%s'", trash.DeletedBy())
	}
}
