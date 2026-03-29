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
