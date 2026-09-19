package entitystore

import (
	"context"
	"testing"
)

func TestStoreEntityCreate(t *testing.T) {
	db := InitDB("store_entity_create_test")
	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "entity_create_entity",
		AttributeTableName: "entity_create_attribute",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	entity := NewEntity()
	entity.SetType("product")
	entity.SetHandle("iphone-15")

	err = store.EntityCreate(context.Background(), entity)
	if err != nil {
		t.Fatal("EntityCreate failed:", err)
	}

	if entity.ID() == "" {
		t.Fatal("Entity ID should be generated")
	}
}

func TestStoreEntityFindByID(t *testing.T) {
	db := InitDB("store_entity_find_test")
	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "entity_find_entity",
		AttributeTableName: "entity_find_attribute",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	entity := NewEntity()
	entity.SetType("product")
	err = store.EntityCreate(context.Background(), entity)
	if err != nil {
		t.Fatal("EntityCreate failed:", err)
	}

	found, err := store.EntityFindByID(context.Background(), entity.ID())
	if err != nil {
		t.Fatal("EntityFindByID failed:", err)
	}
	if found == nil {
		t.Fatal("Entity should be found")
	}
}

func TestStoreEntityList(t *testing.T) {
	db := InitDB("store_entity_list_test")
	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "entity_list_entity",
		AttributeTableName: "entity_list_attribute",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		entity := NewEntity()
		entity.SetType("product")
		err = store.EntityCreate(context.Background(), entity)
		if err != nil {
			t.Fatal("EntityCreate failed:", err)
		}
	}

	list, err := store.EntityList(context.Background(), EntityQuery())
	if err != nil {
		t.Fatal("EntityList failed:", err)
	}
	if len(list) != 3 {
		t.Fatal("Expected 3 entities, got:", len(list))
	}
}

func TestStoreEntityListByAttributeErrorPropagation(t *testing.T) {
	db := InitDB("entity_list_by_attr_error_test")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "err_prop_entity",
		AttributeTableName: "err_prop_attribute",
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

	err = store.AttributeSetString(ctx, entity.ID(), "color", "red")
	if err != nil {
		t.Fatal("AttributeSetString failed:", err)
	}

	_, err = db.Exec("DROP TABLE err_prop_attribute")
	if err != nil {
		t.Fatal("DROP TABLE failed:", err)
	}

	results, err := store.EntityListByAttribute(ctx, "product", "color", "red")
	if err == nil {
		t.Fatal("Expected error when attributes table is missing, got nil. " +
			"EntityListByAttribute silently swallows database errors from AttributeFind.")
	}
	if results != nil {
		t.Fatalf("Expected nil results when error occurs, got %d entities", len(results))
	}
}

func TestStoreEntityListByAttributeJoin(t *testing.T) {
	db := InitDB("entity_list_by_attr_join_test")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "join_entity",
		AttributeTableName: "join_attribute",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("NewStore failed:", err)
	}

	ctx := context.Background()

	// Create 5 products: 3 with color=red, 2 with color=blue
	for i := 0; i < 5; i++ {
		entity, err := store.EntityCreateWithType(ctx, "product")
		if err != nil {
			t.Fatal("EntityCreateWithType failed:", err)
		}
		color := "blue"
		if i < 3 {
			color = "red"
		}
		err = store.AttributeSetString(ctx, entity.ID(), "color", color)
		if err != nil {
			t.Fatal("AttributeSetString failed:", err)
		}
	}

	// Create an order with color=red to verify type filtering
	orderEntity, err := store.EntityCreateWithType(ctx, "order")
	if err != nil {
		t.Fatal("EntityCreateWithType failed:", err)
	}
	err = store.AttributeSetString(ctx, orderEntity.ID(), "color", "red")
	if err != nil {
		t.Fatal("AttributeSetString failed:", err)
	}

	// Query: products with color=red → should return exactly 3
	results, err := store.EntityListByAttribute(ctx, "product", "color", "red")
	if err != nil {
		t.Fatal("EntityListByAttribute failed:", err)
	}
	if len(results) != 3 {
		t.Fatalf("Expected 3 products with color=red, got %d", len(results))
	}

	// Verify all results are type "product"
	for _, e := range results {
		if e.GetType() != "product" {
			t.Fatalf("Expected all results to be type 'product', got '%s'", e.GetType())
		}
	}

	// Verify no duplicate entity IDs
	seen := make(map[string]bool)
	for _, e := range results {
		if seen[e.ID()] {
			t.Fatalf("Duplicate entity ID %s in results", e.ID())
		}
		seen[e.ID()] = true
	}

	// Query: products with color=blue → should return exactly 2
	resultsBlue, err := store.EntityListByAttribute(ctx, "product", "color", "blue")
	if err != nil {
		t.Fatal("EntityListByAttribute failed:", err)
	}
	if len(resultsBlue) != 2 {
		t.Fatalf("Expected 2 products with color=blue, got %d", len(resultsBlue))
	}

	// Query: orders with color=red → should return exactly 1
	resultsOrder, err := store.EntityListByAttribute(ctx, "order", "color", "red")
	if err != nil {
		t.Fatal("EntityListByAttribute failed:", err)
	}
	if len(resultsOrder) != 1 {
		t.Fatalf("Expected 1 order with color=red, got %d", len(resultsOrder))
	}

	// Query: products with non-existent value → should return 0
	resultsNone, err := store.EntityListByAttribute(ctx, "product", "color", "green")
	if err != nil {
		t.Fatal("EntityListByAttribute failed:", err)
	}
	if len(resultsNone) != 0 {
		t.Fatalf("Expected 0 products with color=green, got %d", len(resultsNone))
	}
}
