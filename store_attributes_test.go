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

	list, err := store.AttributeList(context.Background(), AttributeQueryOptions{EntityID: entityID})
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
