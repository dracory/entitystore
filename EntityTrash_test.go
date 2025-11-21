package entitystore

import (
	"context"
	"testing"
)

func TestEntityTrash(t *testing.T) {
	db := InitDB("test_entity_trash.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Must be NIL:", err)
	}

	entity, err := store.EntityCreateWithTypeAndAttributes(context.Background(), "post", map[string]string{
		"title": "Test Post Title",
		"text":  "Test Post Text",
	})

	if err != nil {
		t.Fatal("Entity could not be created:", err)
	}

	if entity == nil {
		t.Fatal("Entity could not be created")
	}

	attr, err := store.AttributeFind(context.Background(), entity.ID(), "title")

	if err != nil {
		t.Fatal("Attribute could not be found:", err)
	}

	if attr == nil {
		t.Fatal("Attribute should not be nil")
	}

	isDeleted, err := store.EntityTrash(context.Background(), entity.ID())

	if err != nil {
		t.Fatal("Entity could not be deleted:", err)
	}

	if isDeleted == false {
		t.Fatal("Entity could not be soft deleted")
	}

	val, err := store.EntityFindByID(context.Background(), entity.ID())

	if err != nil {
		t.Fatal(err)
	}

	if val != nil {
		t.Fatal("Entity should no longer be present")
	}

	attr, err = store.AttributeFind(context.Background(), entity.ID(), "title")

	if err != nil {
		t.Fatal("Attribute could not be found:", err)
	}

	if attr != nil {
		t.Fatal("Attribute should be nil")
	}
}
