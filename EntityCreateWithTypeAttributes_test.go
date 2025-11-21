package entitystore

import (
	"context"
	"testing"
)

func TestEntityCreateWithAttributes(t *testing.T) {
	db := InitDB("test_entity_create_with_type_and_attributes.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Must be NIL:", err.Error())
	}

	entity, err := store.EntityCreateWithTypeAndAttributes(context.Background(), "post", map[string]string{
		"name": "Hello world",
	})

	if err != nil {
		t.Fatal("Entity could not be created:", err.Error())
	}

	if entity == nil {
		t.Fatalf("Entity could not be created")
	}

	val, _ := entity.GetString("name", "")
	if val != "Hello world" {
		t.Fatalf("Entity attribute mismatch")
	}
}
