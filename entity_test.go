package entitystore

import (
	"context"
	"testing"

	"github.com/dromara/carbon/v2"
)

// TestEntityAttributesCreate tests creating an entity and setting attributes via the store
func TestEntityAttributesCreate(t *testing.T) {
	db := InitDB("test_attributes_create.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Must be NIL:", err)
	}

	entity, err := store.EntityCreateWithType(context.Background(), "post")
	if err != nil {
		t.Fatal("Entity could not be created:", err)
	}
	if entity == nil {
		t.Fatal("Entity could not be created")
	}

	// Set title via store (the new pattern — no store back-reference on entity)
	err = store.AttributeSetString(context.Background(), entity.ID(), "title", "Product 1")
	if err != nil {
		t.Fatal("Entity title could not be created:", err)
	}

	titleAttr, err := store.AttributeFind(context.Background(), entity.ID(), "title")
	if err != nil {
		t.Fatal("Entity title could not be retrieved:", err)
	}
	if titleAttr == nil || titleAttr.GetValue() != "Product 1" {
		t.Fatal("Title is incorrect:", titleAttr)
	}

	// Set price_float via store
	err = store.AttributeSetFloat(context.Background(), entity.ID(), "price_float", 12.35)
	if err != nil {
		t.Fatal("Entity price_float could not be created:", err)
	}

	priceFloatAttr, err := store.AttributeFind(context.Background(), entity.ID(), "price_float")
	if err != nil {
		t.Fatal("Entity price_float could not be retrieved:", err)
	}
	if priceFloatAttr == nil {
		t.Fatal("price_float attribute is nil")
	}
	priceFloat, err := priceFloatAttr.GetFloat()
	if err != nil {
		t.Fatal("price_float could not be parsed:", err)
	}
	if priceFloat != 12.35 {
		t.Fatal("Price float is incorrect:", priceFloat)
	}

	// Set price_int via store
	err = store.AttributeSetInt(context.Background(), entity.ID(), "price_int", 12)
	if err != nil {
		t.Fatal("Entity price_int could not be created:", err)
	}

	priceIntAttr, err := store.AttributeFind(context.Background(), entity.ID(), "price_int")
	if err != nil {
		t.Fatal("Entity price_int could not be retrieved:", err)
	}
	if priceIntAttr == nil {
		t.Fatal("price_int attribute is nil")
	}
	priceInt, err := priceIntAttr.GetInt()
	if err != nil {
		t.Fatal("price_int could not be parsed:", err)
	}
	if priceInt != 12 {
		t.Fatal("Price int is incorrect:", priceInt)
	}

	// Set description via store
	err = store.AttributeSetString(context.Background(), entity.ID(), "description", "Description text")
	if err != nil {
		t.Fatal("Entity description could not be created:", err)
	}

	descAttr, err := store.AttributeFind(context.Background(), entity.ID(), "description")
	if err != nil {
		t.Fatal("Entity description could not be retrieved:", err)
	}
	if descAttr == nil || descAttr.GetValue() != "Description text" {
		t.Fatal("Description is incorrect:", descAttr)
	}
}

func TestEntityImplementation(t *testing.T) {
	entity := NewEntity()

	// Test ID generation
	if entity.ID() == "" {
		t.Error("expected ID to be generated")
	}
	if len(entity.ID()) < 9 {
		t.Errorf("expected ID length >= 9, got %d", len(entity.ID()))
	}

	// Test GetID() method consistency
	if entity.GetID() != entity.ID() {
		t.Errorf("expected GetID() '%s' to match ID() '%s'", entity.GetID(), entity.ID())
	}

	// Test EntityType getter/setter
	entity.SetType("product")
	if entity.GetType() != "product" {
		t.Errorf("expected EntityType 'product', got '%s'", entity.GetType())
	}

	// Test EntityHandle getter/setter
	entity.SetHandle("my-handle")
	if entity.GetHandle() != "my-handle" {
		t.Errorf("expected EntityHandle 'my-handle', got '%s'", entity.GetHandle())
	}

	// Test CreatedAt
	if entity.GetCreatedAt() == "" {
		t.Error("expected CreatedAt to be set")
	}

	// Test UpdatedAt
	if entity.GetUpdatedAt() == "" {
		t.Error("expected UpdatedAt to be set")
	}

	// Test Carbon helpers
	createdAtCarbon := entity.GetCreatedAtCarbon()
	if createdAtCarbon == nil {
		t.Error("expected CreatedAtCarbon to return valid carbon")
	}

	updatedAtCarbon := entity.GetUpdatedAtCarbon()
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

	if entity.GetType() != "product" {
		t.Errorf("expected EntityType 'product', got '%s'", entity.GetType())
	}

	if entity.GetHandle() != "my-handle" {
		t.Errorf("expected EntityHandle 'my-handle', got '%s'", entity.GetHandle())
	}

	// Test dynamic attribute access
	if entity.GetTempKey("custom_field") != "custom_value" {
		t.Errorf("expected custom_field 'custom_value', got '%s'", entity.GetTempKey("custom_field"))
	}
}

func TestEntityDynamicAttributes(t *testing.T) {
	entity := NewEntity()

	// Test SetTempKey / GetTempKey
	entity.SetTempKey("name", "iPhone")
	if entity.GetTempKey("name") != "iPhone" {
		t.Errorf("expected attribute 'name' to be 'iPhone', got '%s'", entity.GetTempKey("name"))
	}

	// Test GetTempKeys
	entity.SetTempKey("price", "999")
	allAttrs := entity.GetTempKeys()

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
		t.Error("expected COLUMN_ID to be excluded from GetTempKeys")
	}
}

func TestEntityFluentInterface(t *testing.T) {
	entity := NewEntity().
		SetType("product").
		SetHandle("iphone-15").
		SetTempKey("name", "iPhone 15").
		SetTempKey("price", "999")

	if entity.GetType() != "product" {
		t.Errorf("expected EntityType 'product', got '%s'", entity.GetType())
	}

	if entity.GetHandle() != "iphone-15" {
		t.Errorf("expected EntityHandle 'iphone-15', got '%s'", entity.GetHandle())
	}

	if entity.GetTempKey("name") != "iPhone 15" {
		t.Errorf("expected attribute 'name' = 'iPhone 15', got '%s'", entity.GetTempKey("name"))
	}
}

func TestEntityDataObject(t *testing.T) {
	entity := NewEntity()
	entity.SetType("test")
	entity.SetTempKey("foo", "bar")

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
	if entity.GetType() != "hydrated" {
		t.Errorf("expected EntityType 'hydrated' after Hydrate, got '%s'", entity.GetType())
	}
}
