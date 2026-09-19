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
		DB:                   db,
		EntityTableName:      "test_entities",
		AttributeTableName:   "test_attributes",
		AutomigrateEnabled:   true,
		RelationshipsEnabled: true,
		TaxonomiesEnabled:    true,
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

	product := active.EntityCreate("product").
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

	created := active.EntityCreate("product").SetString("name", "Phone")
	if err := created.Save(); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	found, err := active.EntityFindByID(created.GetEntity().ID())
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
		if err := active.EntityCreate("tag").Save(); err != nil {
			t.Fatalf("Failed to save: %v", err)
		}
	}

	query := entitystore.EntityQuery().WithEntityType("tag")
	list, err := active.EntityList(query)
	if err != nil || len(list) != 2 {
		t.Fatalf("Expected 2 tags, got %d err=%v", len(list), err)
	}

	trashed, err := list[0].Trash()
	if err != nil || !trashed {
		t.Fatalf("Expected trashed=true, got %v err=%v", trashed, err)
	}

	count, _ := active.EntityCount(query)
	if count != 1 {
		t.Errorf("Expected count 1 after trash, got %d", count)
	}
}

func TestRelationshipsAndTaxonomy(t *testing.T) {
	active, cleanup := setupTestActiveStore(t)
	defer cleanup()

	post := active.EntityCreate("post")
	author := active.EntityCreate("author").SetString("name", "Ada")
	for _, e := range []activestore.ActiveEntityInterface{post, author} {
		if err := e.Save(); err != nil {
			t.Fatalf("Failed to save: %v", err)
		}
	}

	if err := post.RelateTo(author.GetEntity().ID(), "written_by"); err != nil {
		t.Fatalf("Failed to relate: %v", err)
	}
	authors, err := post.Related("written_by")
	if err != nil || len(authors) != 1 {
		t.Fatalf("Expected 1 author, got %d err=%v", len(authors), err)
	}

	cat, err := active.TaxonomyCreate(entitystore.TaxonomyOptions{Name: "Categories", Slug: "categories"})
	if err != nil {
		t.Fatalf("Failed to create taxonomy: %v", err)
	}
	term, err := active.TermCreate(entitystore.TaxonomyTermOptions{
		TaxonomyID: cat.GetTaxonomy().GetID(), Name: "Go", Slug: "go",
	})
	if err != nil {
		t.Fatalf("Failed to create term: %v", err)
	}
	if err := post.AssignTerm(cat.GetTaxonomy().GetID(), term.GetTerm().GetID()); err != nil {
		t.Fatalf("Failed to assign term: %v", err)
	}
	tagged, err := term.Entities()
	if err != nil || len(tagged) != 1 {
		t.Fatalf("Expected 1 tagged entity, got %d err=%v", len(tagged), err)
	}
}
