package main

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dracory/entitystore"
	"github.com/dracory/entitystore/activestore"
	_ "modernc.org/sqlite"
)

func TestAssignAndNavigate(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer func() { _ = db.Close() }()

	store, err := entitystore.NewStore(entitystore.NewStoreOptions{
		DB:                 db,
		EntityTableName:    "test_entities",
		AttributeTableName: "test_attributes",
		AutomigrateEnabled: true,
		TaxonomiesEnabled:  true,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	active, err := activestore.New(context.Background(), store)
	if err != nil {
		t.Fatalf("Failed to create active store: %v", err)
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

	post := active.EntityCreate("post")
	if _, err := post.Save(); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}
	if err := post.AssignTerm(cat, term); err != nil {
		t.Fatalf("Failed to assign term: %v", err)
	}

	terms, err := post.Terms(cat)
	if err != nil || len(terms) != 1 {
		t.Fatalf("Expected 1 term, got %d err=%v", len(terms), err)
	}

	tagged, err := term.Entities()
	if err != nil || len(tagged) != 1 {
		t.Fatalf("Expected 1 tagged entity, got %d err=%v", len(tagged), err)
	}
}
