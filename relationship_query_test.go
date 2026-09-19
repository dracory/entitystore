package entitystore

import (
	"context"
	"testing"
)

// TestRelationshipQueryFluentOptions verifies that the fluent builder
// records the explicitly set fields and tracks presence via has* flags.
func TestRelationshipQueryFluentOptions(t *testing.T) {
	query := RelationshipQuery().
		WithEntityID("entity-1").
		WithRelatedEntityID("entity-2").
		WithRelationshipType(RELATIONSHIP_TYPE_BELONGS_TO).
		WithLimit(100).
		WithSortOrder("desc")

	if query.GetEntityID() != "entity-1" {
		t.Errorf("EntityID: expected 'entity-1', got %q", query.GetEntityID())
	}
	if query.GetRelatedEntityID() != "entity-2" {
		t.Errorf("RelatedEntityID: expected 'entity-2', got %q", query.GetRelatedEntityID())
	}
	if query.GetRelationshipType() != RELATIONSHIP_TYPE_BELONGS_TO {
		t.Errorf("RelationshipType: expected %q, got %q", RELATIONSHIP_TYPE_BELONGS_TO, query.GetRelationshipType())
	}
	if query.GetLimit() != 100 {
		t.Errorf("Limit: expected 100, got %d", query.GetLimit())
	}
	if query.GetSortOrder() != "desc" {
		t.Errorf("SortOrder: expected 'desc', got %q", query.GetSortOrder())
	}

	// Presence flags: set fields report true, unset fields report false.
	if !query.HasEntityID() || !query.HasRelatedEntityID() || !query.HasRelationshipType() ||
		!query.HasLimit() || !query.HasSortOrder() {
		t.Error("Set fields should report Has*() = true")
	}
	if query.HasID() || query.HasIDs() || query.HasEntityIDs() || query.HasRelatedEntityIDs() ||
		query.HasParentID() || query.HasOffset() || query.HasSortBy() || query.HasCountOnly() {
		t.Error("Unset fields should report Has*() = false")
	}
}

// TestRelationshipQueryValidate verifies Validate() catches set-but-invalid
// values and ignores unset fields.
func TestRelationshipQueryValidate(t *testing.T) {
	// Empty builder is valid
	if err := RelationshipQuery().Validate(); err != nil {
		t.Error("Empty query should be valid:", err)
	}

	cases := []struct {
		name  string
		query RelationshipQueryInterface
	}{
		{"empty id", RelationshipQuery().WithID("")},
		{"empty ids", RelationshipQuery().WithIDs([]string{})},
		{"empty entity id", RelationshipQuery().WithEntityID("")},
		{"empty entity ids", RelationshipQuery().WithEntityIDs([]string{})},
		{"empty related entity id", RelationshipQuery().WithRelatedEntityID("")},
		{"empty related entity ids", RelationshipQuery().WithRelatedEntityIDs(nil)},
		{"empty relationship type", RelationshipQuery().WithRelationshipType("")},
		{"empty parent id", RelationshipQuery().WithParentID("")},
		{"empty sort by", RelationshipQuery().WithSortBy("")},
		{"bad sort order", RelationshipQuery().WithSortOrder("sideways")},
	}
	for _, c := range cases {
		if err := c.query.Validate(); err == nil {
			t.Errorf("%s: expected validation error, got nil", c.name)
		}
	}
}

// TestRelationshipQueryIntegration runs the fluent query against a real store.
func TestRelationshipQueryIntegration(t *testing.T) {
	db := InitDB("relationship_fluent_test")

	store, err := NewStore(NewStoreOptions{
		DB:                    db,
		EntityTableName:       "fluent_entity",
		AttributeTableName:    "fluent_attribute",
		RelationshipsEnabled:  true,
		RelationshipTableName: "fluent_relationship",
		AutomigrateEnabled:    true,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	author, err := store.EntityCreateWithType(ctx, "author")
	if err != nil {
		t.Fatal(err)
	}
	book1, err := store.EntityCreateWithType(ctx, "book")
	if err != nil {
		t.Fatal(err)
	}
	book2, err := store.EntityCreateWithType(ctx, "book")
	if err != nil {
		t.Fatal(err)
	}

	for _, book := range []EntityInterface{book1, book2} {
		_, err := store.RelationshipCreateByOptions(ctx, RelationshipOptions{
			EntityID:         book.ID(),
			RelatedEntityID:  author.ID(),
			RelationshipType: RELATIONSHIP_TYPE_BELONGS_TO,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	query := RelationshipQuery().
		WithRelatedEntityID(author.ID()).
		WithRelationshipType(RELATIONSHIP_TYPE_BELONGS_TO)
	if err := query.Validate(); err != nil {
		t.Fatal("Validate failed:", err)
	}

	list, err := store.RelationshipList(ctx, query)
	if err != nil {
		t.Fatal("RelationshipList failed:", err)
	}
	if len(list) != 2 {
		t.Fatalf("Expected 2 relationships, got %d", len(list))
	}

	count, err := store.RelationshipCount(ctx, query)
	if err != nil {
		t.Fatal("RelationshipCount failed:", err)
	}
	if count != 2 {
		t.Fatalf("Expected count 2, got %d", count)
	}
}
