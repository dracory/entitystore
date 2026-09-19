package activestore

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dracory/entitystore"
	_ "modernc.org/sqlite"
)

func initFullStore(t *testing.T, name string) entitystore.StoreInterface {
	t.Helper()
	db, err := sql.Open("sqlite", "file:"+name+"?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	store, err := entitystore.NewStore(entitystore.NewStoreOptions{
		DB:                   db,
		EntityTableName:      "entity",
		AttributeTableName:   "attribute",
		AutomigrateEnabled:   true,
		RelationshipsEnabled: true,
		TaxonomiesEnabled:    true,
	})
	if err != nil {
		t.Fatal(err)
	}
	return store
}

func TestRelateToAndRelated(t *testing.T) {
	ctx := context.Background()
	active, err := New(ctx, initFullStore(t, "activestore_relate.db"))
	if err != nil {
		t.Fatal(err)
	}

	post := active.EntityCreate("post")
	author := active.EntityCreate("author")
	tag := active.EntityCreate("tag")
	for _, e := range []ActiveEntityInterface{post, author, tag} {
		if _, err := e.Save(); err != nil {
			t.Fatal(err)
		}
	}

	if err := post.RelateTo(author.GetEntity().ID(), "written_by"); err != nil {
		t.Fatal(err)
	}
	if err := post.RelateToOrdered(tag.GetEntity().ID(), "has_tag", 3); err != nil {
		t.Fatal(err)
	}

	authors, err := post.Related("written_by")
	if err != nil || len(authors) != 1 {
		t.Fatalf("expected 1 author, got %v err=%v", len(authors), err)
	}
	if authors[0].GetEntity().ID() != author.GetEntity().ID() {
		t.Fatal("related entity ID mismatch")
	}

	if err := post.Unrelate(author.GetEntity().ID(), "written_by"); err != nil {
		t.Fatal(err)
	}
	authors, _ = post.Related("written_by")
	if len(authors) != 0 {
		t.Fatalf("expected 0 authors after unrelate, got %v", len(authors))
	}
}

func TestRelationshipCRUD(t *testing.T) {
	ctx := context.Background()
	active, err := New(ctx, initFullStore(t, "activestore_relcrud.db"))
	if err != nil {
		t.Fatal(err)
	}

	a := active.EntityCreate("post")
	b := active.EntityCreate("author")
	for _, e := range []ActiveEntityInterface{a, b} {
		if _, err := e.Save(); err != nil {
			t.Fatal(err)
		}
	}

	rel, err := active.RelationshipCreate(entitystore.RelationshipOptions{
		EntityID:         a.GetEntity().ID(),
		RelatedEntityID:  b.GetEntity().ID(),
		RelationshipType: "written_by",
	})
	if err != nil {
		t.Fatal(err)
	}

	related, err := rel.GetRelatedEntity()
	if err != nil || related.GetEntity().ID() != b.GetEntity().ID() {
		t.Fatalf("GetRelatedEntity mismatch: %v", err)
	}

	found, err := active.RelationshipFindByID(rel.GetRelationship().ID())
	if err != nil || found.GetRelationship().ID() != rel.GetRelationship().ID() {
		t.Fatalf("RelationshipFindByID mismatch: %v", err)
	}

	count, err := active.RelationshipCount(entitystore.RelationshipQuery().WithEntityID(a.GetEntity().ID()))
	if err != nil || count != 1 {
		t.Fatalf("expected count=1, got %v err=%v", count, err)
	}

	trashed, err := rel.Trash("tester")
	if err != nil || !trashed {
		t.Fatalf("expected trashed=true, got %v err=%v", trashed, err)
	}
	restored, err := rel.Restore()
	if err != nil || !restored {
		t.Fatalf("expected restored=true, got %v err=%v", restored, err)
	}
	deleted, err := rel.Delete()
	if err != nil || !deleted {
		t.Fatalf("expected deleted=true, got %v err=%v", deleted, err)
	}
}
