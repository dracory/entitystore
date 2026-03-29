package entitystore

import (
	"context"
	"testing"
)

func TestEntityCreateWithType(t *testing.T) {
	db := InitDB("test_entity_create_with_type.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	entity, err := store.EntityCreateWithType(context.Background(), "post")
	if entity == nil {
		t.Fatal("Entity could not be created")
	}

	if err != nil {
		t.Fatal("Entity could not be created:", err)
	}

	if len(entity.ID()) < 9 || len(entity.ID()) > 15 {
		t.Fatal("Entity ID is not a short ID (expected 9-15 chars):", entity.ID(), "len:", len(entity.ID()))
	}

	// CreatedAt() is now a string — use CreatedAtCarbon() for time comparisons
	createdAt := entity.GetCreatedAtCarbon()
	if createdAt.IsZero() {
		t.Fatal("Entity CreatedAt is empty")
	}

	updatedAt := entity.GetUpdatedAtCarbon()
	if updatedAt.IsZero() {
		t.Fatal("Entity UpdatedAt is empty")
	}
}
