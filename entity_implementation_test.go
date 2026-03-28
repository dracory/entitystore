package entitystore

import (
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestEntityImplementation(t *testing.T) {
	entity := NewEntity()

	// Test ID generation
	if entity.ID() == "" {
		t.Error("expected ID to be generated")
	}
	if len(entity.ID()) < 9 {
		t.Errorf("expected ID length >= 9, got %d", len(entity.ID()))
	}

	// Test EntityType getter/setter
	entity.SetEntityType("product")
	if entity.EntityType() != "product" {
		t.Errorf("expected EntityType 'product', got '%s'", entity.EntityType())
	}

	// Test EntityHandle getter/setter
	entity.SetEntityHandle("my-handle")
	if entity.EntityHandle() != "my-handle" {
		t.Errorf("expected EntityHandle 'my-handle', got '%s'", entity.EntityHandle())
	}

	// Test CreatedAt
	if entity.CreatedAt() == "" {
		t.Error("expected CreatedAt to be set")
	}

	// Test UpdatedAt
	if entity.UpdatedAt() == "" {
		t.Error("expected UpdatedAt to be set")
	}

	// Test Carbon helpers
	createdAtCarbon := entity.CreatedAtCarbon()
	if createdAtCarbon == nil {
		t.Error("expected CreatedAtCarbon to return valid carbon")
	}

	updatedAtCarbon := entity.UpdatedAtCarbon()
	if updatedAtCarbon == nil {
		t.Error("expected UpdatedAtCarbon to return valid carbon")
	}
}

func TestEntityFromExistingData(t *testing.T) {
	data := map[string]string{
		COLUMN_ID:            "abc123",
		COLUMN_ENTITY_TYPE:   "product",
		COLUMN_ENTITY_HANDLE: "my-handle",
		COLUMN_CREATED_AT:    carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC),
		COLUMN_UPDATED_AT:    carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC),
		"custom_field":       "custom_value",
	}

	entity := NewEntityFromExistingData(data)

	if entity.ID() != "abc123" {
		t.Errorf("expected ID 'abc123', got '%s'", entity.ID())
	}

	if entity.EntityType() != "product" {
		t.Errorf("expected EntityType 'product', got '%s'", entity.EntityType())
	}

	if entity.EntityHandle() != "my-handle" {
		t.Errorf("expected EntityHandle 'my-handle', got '%s'", entity.EntityHandle())
	}

	// Test dynamic attribute access
	if entity.GetAttribute("custom_field") != "custom_value" {
		t.Errorf("expected custom_field 'custom_value', got '%s'", entity.GetAttribute("custom_field"))
	}
}

func TestEntityDynamicAttributes(t *testing.T) {
	entity := NewEntity()

	// Test SetAttribute / GetAttribute
	entity.SetAttribute("name", "iPhone")
	if entity.GetAttribute("name") != "iPhone" {
		t.Errorf("expected attribute 'name' to be 'iPhone', got '%s'", entity.GetAttribute("name"))
	}

	// Test GetAllAttributes
	entity.SetAttribute("price", "999")
	allAttrs := entity.GetAllAttributes()

	if len(allAttrs) != 2 {
		t.Errorf("expected 2 dynamic attributes, got %d", len(allAttrs))
	}

	if allAttrs["name"] != "iPhone" {
		t.Errorf("expected allAttrs['name'] = 'iPhone', got '%s'", allAttrs["name"])
	}

	if allAttrs["price"] != "999" {
		t.Errorf("expected allAttrs['price'] = '999', got '%s'", allAttrs["price"])
	}

	// Test that system columns are excluded
	if _, exists := allAttrs[COLUMN_ID]; exists {
		t.Error("expected COLUMN_ID to be excluded from GetAllAttributes")
	}
}

func TestEntityFluentInterface(t *testing.T) {
	entity := NewEntity().
		SetEntityType("product").
		SetEntityHandle("iphone-15").
		SetAttribute("name", "iPhone 15").
		SetAttribute("price", "999")

	if entity.EntityType() != "product" {
		t.Errorf("expected EntityType 'product', got '%s'", entity.EntityType())
	}

	if entity.EntityHandle() != "iphone-15" {
		t.Errorf("expected EntityHandle 'iphone-15', got '%s'", entity.EntityHandle())
	}

	if entity.GetAttribute("name") != "iPhone 15" {
		t.Errorf("expected attribute 'name' = 'iPhone 15', got '%s'", entity.GetAttribute("name"))
	}
}

func TestEntityDataObject(t *testing.T) {
	entity := NewEntity()
	entity.SetEntityType("test")
	entity.SetAttribute("foo", "bar")

	// Test Data() returns underlying map
	data := entity.Data()
	if data[COLUMN_ENTITY_TYPE] != "test" {
		t.Errorf("expected Data() to contain entity_type='test'")
	}

	// Test Hydrate
	newData := map[string]string{
		COLUMN_ID:          "new-id",
		COLUMN_ENTITY_TYPE: "hydrated",
		"dynamic":          "value",
	}

	entity.Hydrate(newData)
	if entity.ID() != "new-id" {
		t.Errorf("expected ID 'new-id' after Hydrate, got '%s'", entity.ID())
	}
	if entity.EntityType() != "hydrated" {
		t.Errorf("expected EntityType 'hydrated' after Hydrate, got '%s'", entity.EntityType())
	}
}
