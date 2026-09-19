package main

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dracory/entitystore"
	"github.com/dracory/entitystore/activestore"
	_ "modernc.org/sqlite"
)

func TestRelateAndUnrelate(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close() //nolint:errcheck

	store, err := entitystore.NewStore(entitystore.NewStoreOptions{
		DB:                   db,
		EntityTableName:      "test_entities",
		AttributeTableName:   "test_attributes",
		AutomigrateEnabled:   true,
		RelationshipsEnabled: true,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	active, err := activestore.New(context.Background(), store)
	if err != nil {
		t.Fatalf("Failed to create active store: %v", err)
	}

	post := active.EntityCreate("post")
	author := active.EntityCreate("author")
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

	if err := post.Unrelate(author.GetEntity().ID(), "written_by"); err != nil {
		t.Fatalf("Failed to unrelate: %v", err)
	}

	authors, _ = post.Related("written_by")
	if len(authors) != 0 {
		t.Errorf("Expected 0 authors after unrelate, got %d", len(authors))
	}
}
