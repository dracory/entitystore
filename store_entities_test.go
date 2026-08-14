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

	list, err := store.EntityList(context.Background(), EntityQueryOptions{})
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
