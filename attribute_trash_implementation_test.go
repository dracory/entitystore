package entitystore

import (
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestAttributeTrashImplementation(t *testing.T) {
	trash := NewAttributeTrash()

	// Test ID generation
	if trash.ID() == "" {
		t.Error("expected ID to be generated")
	}
	if len(trash.ID()) < 9 {
		t.Errorf("expected ID length >= 9, got %d", len(trash.ID()))
	}

	// Test EntityID getter/setter
	trash.SetEntityID("entity123")
	if trash.EntityID() != "entity123" {
		t.Errorf("expected EntityID 'entity123', got '%s'", trash.EntityID())
	}

	// Test AttributeKey getter/setter
	trash.SetAttributeKey("name")
	if trash.AttributeKey() != "name" {
		t.Errorf("expected AttributeKey 'name', got '%s'", trash.AttributeKey())
	}

	// Test AttributeValue getter/setter
	trash.SetAttributeValue("value")
	if trash.AttributeValue() != "value" {
		t.Errorf("expected AttributeValue 'value', got '%s'", trash.AttributeValue())
	}

	// Test DeletedAt
	if trash.DeletedAt() == "" {
		t.Error("expected DeletedAt to be set")
	}

	// Test DeletedBy getter/setter
	trash.SetDeletedBy("user456")
	if trash.DeletedBy() != "user456" {
		t.Errorf("expected DeletedBy 'user456', got '%s'", trash.DeletedBy())
	}
}

func TestAttributeTrashFromExistingData(t *testing.T) {
	data := map[string]string{
		COLUMN_ID:              "attrtrash123",
		COLUMN_ENTITY_ID:       "entity789",
		COLUMN_ATTRIBUTE_KEY:   "color",
		COLUMN_ATTRIBUTE_VALUE: "blue",
		COLUMN_CREATED_AT:      carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC),
		COLUMN_UPDATED_AT:      carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC),
		COLUMN_DELETED_AT:      carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC),
		COLUMN_DELETED_BY:      "admin",
	}

	trash := NewAttributeTrashFromExistingData(data)

	if trash.ID() != "attrtrash123" {
		t.Errorf("expected ID 'attrtrash123', got '%s'", trash.ID())
	}

	if trash.EntityID() != "entity789" {
		t.Errorf("expected EntityID 'entity789', got '%s'", trash.EntityID())
	}

	if trash.DeletedBy() != "admin" {
		t.Errorf("expected DeletedBy 'admin', got '%s'", trash.DeletedBy())
	}
}
