package entitystore

import (
	"context"
	"testing"
	"time"
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

	if entity.CreatedAt().Before(time.Now().Add(-1 * time.Minute)) {
		t.Fatal("Entity CreatedAt is not recent (before 1 min):", entity.CreatedAt())
	}

	if entity.CreatedAt().After(time.Now().Add(1 * time.Minute)) {
		t.Fatal("Entity CreatedAt is not recent (after 1 min):", entity.CreatedAt())
	}

	if entity.UpdatedAt().Before(time.Now().Add(-1 * time.Minute)) {
		t.Fatal("Entity UpdatedAt is not recent (before 1 min):", entity.UpdatedAt())
	}

	if entity.UpdatedAt().After(time.Now().Add(1 * time.Minute)) {
		t.Fatal("Entity UpdatedAt is not recent (after 1 min):", entity.UpdatedAt())
	}
}
