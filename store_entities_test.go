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
	store.EntityCreate(context.Background(), entity)

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
		store.EntityCreate(context.Background(), entity)
	}

	list, err := store.EntityList(context.Background(), EntityQueryOptions{})
	if err != nil {
		t.Fatal("EntityList failed:", err)
	}
	if len(list) != 3 {
		t.Fatal("Expected 3 entities, got:", len(list))
	}
}
