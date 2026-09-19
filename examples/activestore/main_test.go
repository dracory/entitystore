package main

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dracory/entitystore"
	"github.com/dracory/entitystore/activestore"
	_ "modernc.org/sqlite"
)

func setupTestActiveStore(t *testing.T) (activestore.ActiveStoreInterface, func()) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	store, err := entitystore.NewStore(entitystore.NewStoreOptions{
		DB:                 db,
		EntityTableName:    "test_entities",
		AttributeTableName: "test_attributes",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	active, err := activestore.New(context.Background(), store)
	if err != nil {
		t.Fatalf("Failed to create active store: %v", err)
	}

	return active, func() { db.Close() } //nolint:errcheck
}

func TestFluentCreate(t *testing.T) {
	active, cleanup := setupTestActiveStore(t)
	defer cleanup()

	product := active.New("product").
		SetString("name", "Laptop").
		SetFloat("price", 1299.99).
		SetInt("stock", 50)

	if err := product.Save(); err != nil {
		t.Fatalf("Failed to save product: %v", err)
	}

	name, exists, err := product.GetString("name")
	if err != nil || !exists {
		t.Fatalf("Failed to read name: %v", err)
	}
	if name != "Laptop" {
		t.Errorf("Expected name 'Laptop', got '%s'", name)
	}
}

func TestFindAndUpdate(t *testing.T) {
	active, cleanup := setupTestActiveStore(t)
	defer cleanup()

	created := active.New("product").SetString("name", "Phone")
	if err := created.Save(); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	found, err := active.FindByID(created.GetEntity().ID())
	if err != nil {
		t.Fatalf("Failed to find: %v", err)
	}

	found.SetString("name", "Tablet")
	if err := found.Err(); err != nil {
		t.Fatalf("Failed to update: %v", err)
	}

	name, _, _ := found.GetString("name")
	if name != "Tablet" {
		t.Errorf("Expected name 'Tablet', got '%s'", name)
	}
}

func TestListCountTrash(t *testing.T) {
	active, cleanup := setupTestActiveStore(t)
	defer cleanup()

	for range 2 {
		if err := active.New("tag").Save(); err != nil {
			t.Fatalf("Failed to save: %v", err)
		}
	}

	query := entitystore.EntityQuery().WithEntityType("tag")
	list, err := active.List(query)
	if err != nil || len(list) != 2 {
		t.Fatalf("Expected 2 tags, got %d err=%v", len(list), err)
	}

	trashed, err := list[0].Trash()
	if err != nil || !trashed {
		t.Fatalf("Expected trashed=true, got %v err=%v", trashed, err)
	}

	count, _ := active.Count(query)
	if count != 1 {
		t.Errorf("Expected count 1 after trash, got %d", count)
	}
}
