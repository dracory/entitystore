package entitystore

import (
	"context"
	"testing"
)

func TestStoreAttributeFind(t *testing.T) {
	db := InitDB("store_attr_find_test")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "attr_find_entity",
		AttributeTableName: "attr_find_attribute",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	for _, entityID := range []string{"entity1", "entity2", "entity3"} {
		errSet1 := store.AttributeSetString(context.Background(), entityID, "attr1", "val1")
		if errSet1 != nil {
			t.Fatal(errSet1)
		}
		errSet2 := store.AttributeSetString(context.Background(), entityID, "attr2", "val2")
		if errSet2 != nil {
			t.Fatal(errSet2)
		}
	}

	attr, errFind := store.AttributeFind(context.Background(), "entity2", "attr1")

	if errFind != nil {
		t.Fatal("AttributeFind failed:", errFind)
	}

	if attr == nil {
		t.Fatal("Attribute should be found")
	}

	if attr.GetValue() != "val1" {
		t.Fatal("Attribute value mismatch")
	}
}

func TestStoreAttributeList(t *testing.T) {
	db := InitDB("store_attr_list_test")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "attr_list_entity",
		AttributeTableName: "attr_list_attribute",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	entityID := "test-entity"
	err = store.AttributeSetString(context.Background(), entityID, "attr1", "val1")
	if err != nil {
		t.Fatal("AttributeSetString failed:", err)
	}
	err = store.AttributeSetString(context.Background(), entityID, "attr2", "val2")
	if err != nil {
		t.Fatal("AttributeSetString failed:", err)
	}

	list, err := store.AttributeList(context.Background(), AttributeQuery().WithEntityID(entityID))
	if err != nil {
		t.Fatal("AttributeList failed:", err)
	}

	if len(list) != 2 {
		t.Fatal("Expected 2 attributes, got:", len(list))
	}
}

func TestStoreAttributeGetString(t *testing.T) {
	db := InitDB("store_attr_get_str_test")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "attr_get_str_entity",
		AttributeTableName: "attr_get_str_attribute",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	entityID := "test-entity"

	// Test getting non-existent attribute
	value, exists, err := store.AttributeGetString(ctx, entityID, "missing")
	if err != nil {
		t.Fatal("AttributeGetString error for missing attr:", err)
	}
	if exists {
		t.Fatal("Expected exists=false for missing attribute")
	}
	if value != "" {
		t.Fatal("Expected empty value for missing attribute")
	}

	// Set and get
	err = store.AttributeSetString(ctx, entityID, "name", "John")
	if err != nil {
		t.Fatal("AttributeSetString failed:", err)
	}

	value, exists, err = store.AttributeGetString(ctx, entityID, "name")
	if err != nil {
		t.Fatal("AttributeGetString failed:", err)
	}
	if !exists {
		t.Fatal("Expected exists=true")
	}
	if value != "John" {
		t.Fatalf("Expected 'John', got '%s'", value)
	}
}

func TestStoreAttributeGetInt(t *testing.T) {
	db := InitDB("store_attr_get_int_test")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "attr_get_int_entity",
		AttributeTableName: "attr_get_int_attribute",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	entityID := "test-entity"

	// Test getting non-existent attribute
	value, exists, err := store.AttributeGetInt(ctx, entityID, "missing")
	if err != nil {
		t.Fatal("AttributeGetInt error for missing attr:", err)
	}
	if exists {
		t.Fatal("Expected exists=false for missing attribute")
	}
	if value != 0 {
		t.Fatal("Expected zero value for missing attribute")
	}

	// Set and get
	err = store.AttributeSetInt(ctx, entityID, "age", 42)
	if err != nil {
		t.Fatal("AttributeSetInt failed:", err)
	}

	value, exists, err = store.AttributeGetInt(ctx, entityID, "age")
	if err != nil {
		t.Fatal("AttributeGetInt failed:", err)
	}
	if !exists {
		t.Fatal("Expected exists=true")
	}
	if value != 42 {
		t.Fatalf("Expected 42, got %d", value)
	}

	// Test negative number
	err = store.AttributeSetInt(ctx, entityID, "negative", -100)
	if err != nil {
		t.Fatal("AttributeSetInt failed:", err)
	}

	value, exists, err = store.AttributeGetInt(ctx, entityID, "negative")
	if err != nil {
		t.Fatal("AttributeGetInt failed:", err)
	}
	if !exists {
		t.Fatal("Expected exists=true")
	}
	if value != -100 {
		t.Fatalf("Expected -100, got %d", value)
	}
}

func TestStoreAttributeGetFloat(t *testing.T) {
	db := InitDB("store_attr_get_float_test")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "attr_get_float_entity",
		AttributeTableName: "attr_get_float_attribute",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	entityID := "test-entity"

	// Test getting non-existent attribute
	value, exists, err := store.AttributeGetFloat(ctx, entityID, "missing")
	if err != nil {
		t.Fatal("AttributeGetFloat error for missing attr:", err)
	}
	if exists {
		t.Fatal("Expected exists=false for missing attribute")
	}
	if value != 0 {
		t.Fatal("Expected zero value for missing attribute")
	}

	// Set and get
	err = store.AttributeSetFloat(ctx, entityID, "price", 19.99)
	if err != nil {
		t.Fatal("AttributeSetFloat failed:", err)
	}

	value, exists, err = store.AttributeGetFloat(ctx, entityID, "price")
	if err != nil {
		t.Fatal("AttributeGetFloat failed:", err)
	}
	if !exists {
		t.Fatal("Expected exists=true")
	}
	if value != 19.99 {
		t.Fatalf("Expected 19.99, got %f", value)
	}

	// Test negative number
	err = store.AttributeSetFloat(ctx, entityID, "temperature", -273.15)
	if err != nil {
		t.Fatal("AttributeSetFloat failed:", err)
	}

	value, exists, err = store.AttributeGetFloat(ctx, entityID, "temperature")
	if err != nil {
		t.Fatal("AttributeGetFloat failed:", err)
	}
	if !exists {
		t.Fatal("Expected exists=true")
	}
	if value != -273.15 {
		t.Fatalf("Expected -273.15, got %f", value)
	}
}

// TestStoreAttributeListWithEntityTypeJoin tests the JOIN-based filtering
// by entity type (WP posts/postmeta pattern). Attributes from entities of
// the requested type should be returned; attributes from other entity types
// should be excluded.
func TestStoreAttributeListWithEntityTypeJoin(t *testing.T) {
	db := InitDB("store_attr_join_test")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "attr_join_entity",
		AttributeTableName: "attr_join_attribute",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// Create entities of two types
	entity1, err := store.EntityCreateWithType(ctx, "page_view")
	if err != nil {
		t.Fatal("EntityCreateWithType failed:", err)
	}
	entity2, err := store.EntityCreateWithType(ctx, "page_view")
	if err != nil {
		t.Fatal("EntityCreateWithType failed:", err)
	}
	entity3, err := store.EntityCreateWithType(ctx, "order")
	if err != nil {
		t.Fatal("EntityCreateWithType failed:", err)
	}

	// Set utm_source on all three entities
	if err := store.AttributeSetString(ctx, entity1.GetID(), "utm_source", "reddit"); err != nil {
		t.Fatal("AttributeSetString failed:", err)
	}
	if err := store.AttributeSetString(ctx, entity2.GetID(), "utm_source", "linkedin"); err != nil {
		t.Fatal("AttributeSetString failed:", err)
	}
	if err := store.AttributeSetString(ctx, entity3.GetID(), "utm_source", "direct"); err != nil {
		t.Fatal("AttributeSetString failed:", err)
	}

	// Query: all utm_source attributes for page_view entities only
	list, err := store.AttributeList(ctx, AttributeQuery().
		WithEntityType("page_view").
		WithAttributeKey("utm_source").
		WithLimit(100))
	if err != nil {
		t.Fatal("AttributeList with EntityType join failed:", err)
	}

	if len(list) != 2 {
		t.Fatalf("Expected 2 attributes for page_view entities, got %d", len(list))
	}

	// Verify both page_view attributes are present (reddit + linkedin), order entity excluded
	values := map[string]bool{}
	for _, attr := range list {
		values[attr.GetValue()] = true
	}
	if !values["reddit"] {
		t.Error("Expected 'reddit' in results")
	}
	if !values["linkedin"] {
		t.Error("Expected 'linkedin' in results")
	}
	if values["direct"] {
		t.Error("'direct' (order entity) should be excluded by EntityType filter")
	}
}

// TestStoreAttributeFindByHandleWithJoin tests AttributeFindByHandle which
// relies on the EntityType + EntityHandle join path in AttributeList.
func TestStoreAttributeFindByHandleWithJoin(t *testing.T) {
	db := InitDB("store_attr_handle_join_test")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "attr_handle_join_entity",
		AttributeTableName: "attr_handle_join_attribute",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// Create entity with handle
	entity, err := store.EntityCreateWithType(ctx, "product")
	if err != nil {
		t.Fatal("EntityCreateWithType failed:", err)
	}
	entity.SetHandle("my-product")
	if err := store.EntityUpdate(ctx, entity); err != nil {
		t.Fatal("EntityUpdate failed:", err)
	}

	// Set attribute
	if err := store.AttributeSetString(ctx, entity.GetID(), "price", "29.99"); err != nil {
		t.Fatal("AttributeSetString failed:", err)
	}

	// Find by handle — uses EntityType + EntityHandle join
	attr, err := store.AttributeFindByHandle(ctx, "product", "my-product", "price")
	if err != nil {
		t.Fatal("AttributeFindByHandle failed:", err)
	}
	if attr == nil {
		t.Fatal("Expected attribute to be found by handle")
	}
	if attr.GetValue() != "29.99" {
		t.Fatalf("Expected '29.99', got '%s'", attr.GetValue())
	}

	// Wrong handle should return nil
	attrWrong, err := store.AttributeFindByHandle(ctx, "product", "wrong-handle", "price")
	if err != nil {
		t.Fatal("AttributeFindByHandle failed:", err)
	}
	if attrWrong != nil {
		t.Error("Expected nil for wrong handle")
	}

	// Wrong type should return nil
	attrWrongType, err := store.AttributeFindByHandle(ctx, "order", "my-product", "price")
	if err != nil {
		t.Fatal("AttributeFindByHandle failed:", err)
	}
	if attrWrongType != nil {
		t.Error("Expected nil for wrong entity type")
	}
}

func TestStoreAttributesSetAtomicity(t *testing.T) {
	db := InitDB("attr_set_atomicity_test")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "atomic_entity",
		AttributeTableName: "atomic_attribute",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("NewStore failed:", err)
	}

	ctx := context.Background()

	entity, err := store.EntityCreateWithType(ctx, "product")
	if err != nil {
		t.Fatal("EntityCreateWithType failed:", err)
	}

	// Call AttributesSet with one valid attribute and one with an empty key.
	// The empty key causes AttributeFind to return "attribute key cannot be empty".
	// Map iteration order is non-deterministic, so:
	// - If empty key is iterated first: error returns immediately, nothing committed.
	// - If valid key is iterated first: "color" commits, then empty key fails.
	//   In this case, "color" should be rolled back if the operation is atomic.
	//
	// We run the test with -count=10 to hit both orderings with high probability.
	// The test FAILS when the valid key commits first and is not rolled back.
	attrs := map[string]string{
		"color": "red",
		"":      "should-fail",
	}

	err = store.AttributesSet(ctx, entity.ID(), attrs)
	if err == nil {
		t.Fatal("Expected error from AttributesSet with empty key, got nil")
	}

	// If atomic, "color" should NOT exist after the error.
	// If not atomic, "color" may exist if it was iterated before the empty key.
	_, exists, err := store.AttributeGetString(ctx, entity.ID(), "color")
	if err != nil {
		t.Fatal("AttributeGetString failed:", err)
	}
	if exists {
		t.Fatal("Expected color attribute to NOT exist after AttributesSet failure. " +
			"Partial write detected — AttributesSet is not atomic: the valid attribute " +
			"was committed before the empty-key error and was not rolled back.")
	}
}

// TestStoreAttributeListByKeys tests the AttributeKeys filter —
// a single WHERE IN on attribute_key returning only the requested keys,
// combined with the EntityType join.
func TestStoreAttributeListByKeys(t *testing.T) {
	db := InitDB("store_attr_keys_test")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "attr_keys_entity",
		AttributeTableName: "attr_keys_attribute",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	entity, err := store.EntityCreateWithType(ctx, "character")
	if err != nil {
		t.Fatal("EntityCreateWithType failed:", err)
	}
	other, err := store.EntityCreateWithType(ctx, "order")
	if err != nil {
		t.Fatal("EntityCreateWithType failed:", err)
	}

	for k, v := range map[string]string{"name": "Ada", "title": "Countess", "system_prompt": "secret"} {
		if err := store.AttributeSetString(ctx, entity.ID(), k, v); err != nil {
			t.Fatal("AttributeSetString failed:", err)
		}
	}
	if err := store.AttributeSetString(ctx, other.ID(), "name", "NotAda"); err != nil {
		t.Fatal("AttributeSetString failed:", err)
	}

	// Only the two requested keys, only for the character entity type
	list, err := store.AttributeList(ctx, AttributeQuery().
		WithEntityType("character").
		WithAttributeKeys([]string{"name", "title"}))
	if err != nil {
		t.Fatal("AttributeList failed:", err)
	}
	if len(list) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(list))
	}
	keys := map[string]string{}
	for _, a := range list {
		keys[a.GetKey()] = a.GetValue()
	}
	if keys["name"] != "Ada" {
		t.Errorf("Expected name=Ada, got %q", keys["name"])
	}
	if keys["title"] != "Countess" {
		t.Errorf("Expected title=Countess, got %q", keys["title"])
	}
	if _, ok := keys["system_prompt"]; ok {
		t.Error("system_prompt should be excluded by AttributeKeys filter")
	}

	// AttributeKeys also applies to AttributeCount (shared filter builder)
	count, err := store.AttributeCount(ctx, AttributeQuery().
		WithEntityType("character").
		WithAttributeKeys([]string{"name", "title"}))
	if err != nil {
		t.Fatal("AttributeCount failed:", err)
	}
	if count != 2 {
		t.Fatalf("Expected count 2, got %d", count)
	}
}

func TestStoreAttributeCount(t *testing.T) {
	db := InitDB("attr_count_test")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "count_entity",
		AttributeTableName: "count_attribute",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("NewStore failed:", err)
	}

	ctx := context.Background()

	// Create 3 entities with attributes
	for i := 0; i < 3; i++ {
		entity, err := store.EntityCreateWithType(ctx, "product")
		if err != nil {
			t.Fatal("EntityCreateWithType failed:", err)
		}
		err = store.AttributeSetString(ctx, entity.ID(), "color", "red")
		if err != nil {
			t.Fatal("AttributeSetString failed:", err)
		}
		err = store.AttributeSetString(ctx, entity.ID(), "size", "large")
		if err != nil {
			t.Fatal("AttributeSetString failed:", err)
		}
	}

	// Count all attributes (6 total: 3 entities × 2 attrs)
	count, err := store.AttributeCount(ctx, AttributeQuery())
	if err != nil {
		t.Fatal("AttributeCount failed:", err)
	}
	if count != 6 {
		t.Fatalf("Expected count 6, got %d", count)
	}

	// Count by AttributeKey
	count, err = store.AttributeCount(ctx, AttributeQuery().WithAttributeKey("color"))
	if err != nil {
		t.Fatal("AttributeCount failed:", err)
	}
	if count != 3 {
		t.Fatalf("Expected count 3 for key=color, got %d", count)
	}

	// Count with JOIN path (EntityType filter)
	count, err = store.AttributeCount(ctx, AttributeQuery().WithEntityType("product"))
	if err != nil {
		t.Fatal("AttributeCount with JOIN failed:", err)
	}
	if count != 6 {
		t.Fatalf("Expected count 6 for EntityType=product, got %d", count)
	}

	// Verify count matches list length for the JOIN path
	list, err := store.AttributeList(ctx, AttributeQuery().WithEntityType("product"))
	if err != nil {
		t.Fatal("AttributeList failed:", err)
	}
	if int64(len(list)) != count {
		t.Fatalf("Count (%d) does not match list length (%d) for EntityType=product", count, len(list))
	}

	// Count by EntityID
	entity, err := store.EntityCreateWithType(ctx, "product")
	if err != nil {
		t.Fatal("EntityCreateWithType failed:", err)
	}
	err = store.AttributeSetString(ctx, entity.ID(), "weight", "100")
	if err != nil {
		t.Fatal("AttributeSetString failed:", err)
	}

	count, err = store.AttributeCount(ctx, AttributeQuery().WithEntityID(entity.ID()))
	if err != nil {
		t.Fatal("AttributeCount failed:", err)
	}
	if count != 1 {
		t.Fatalf("Expected count 1 for single entity, got %d", count)
	}
}
