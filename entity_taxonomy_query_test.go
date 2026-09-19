package entitystore

import (
	"context"
	"testing"
)

// TestEntityTaxonomyQueryFluentOptions verifies that the fluent builder
// records the explicitly set fields and tracks presence via has* flags.
func TestEntityTaxonomyQueryFluentOptions(t *testing.T) {
	query := EntityTaxonomyQuery().
		WithEntityID("entity-1").
		WithTaxonomyID("taxonomy-1").
		WithTermIDs([]string{"term-1", "term-2"}).
		WithLimit(100).
		WithSortOrder("asc")

	if query.GetEntityID() != "entity-1" {
		t.Errorf("EntityID: expected 'entity-1', got %q", query.GetEntityID())
	}
	if query.GetTaxonomyID() != "taxonomy-1" {
		t.Errorf("TaxonomyID: expected 'taxonomy-1', got %q", query.GetTaxonomyID())
	}
	termIDs := query.GetTermIDs()
	if len(termIDs) != 2 || termIDs[0] != "term-1" || termIDs[1] != "term-2" {
		t.Errorf("TermIDs mismatch: %v", termIDs)
	}
	if query.GetLimit() != 100 {
		t.Errorf("Limit: expected 100, got %d", query.GetLimit())
	}
	if query.GetSortOrder() != "asc" {
		t.Errorf("SortOrder: expected 'asc', got %q", query.GetSortOrder())
	}

	// Presence flags: set fields report true, unset fields report false.
	if !query.HasEntityID() || !query.HasTaxonomyID() || !query.HasTermIDs() ||
		!query.HasLimit() || !query.HasSortOrder() {
		t.Error("Set fields should report Has*() = true")
	}
	if query.HasID() || query.HasEntityIDs() || query.HasTermID() ||
		query.HasOffset() || query.HasSortBy() || query.HasCountOnly() {
		t.Error("Unset fields should report Has*() = false")
	}
}

// TestEntityTaxonomyQueryValidate verifies Validate() catches set-but-invalid
// values and ignores unset fields.
func TestEntityTaxonomyQueryValidate(t *testing.T) {
	// Empty builder is valid
	if err := EntityTaxonomyQuery().Validate(); err != nil {
		t.Error("Empty query should be valid:", err)
	}

	cases := []struct {
		name  string
		query EntityTaxonomyQueryInterface
	}{
		{"empty id", EntityTaxonomyQuery().WithID("")},
		{"empty entity id", EntityTaxonomyQuery().WithEntityID("")},
		{"empty entity ids", EntityTaxonomyQuery().WithEntityIDs([]string{})},
		{"empty taxonomy id", EntityTaxonomyQuery().WithTaxonomyID("")},
		{"empty term id", EntityTaxonomyQuery().WithTermID("")},
		{"empty term ids", EntityTaxonomyQuery().WithTermIDs(nil)},
		{"empty sort by", EntityTaxonomyQuery().WithSortBy("")},
		{"bad sort order", EntityTaxonomyQuery().WithSortOrder("sideways")},
	}
	for _, c := range cases {
		if err := c.query.Validate(); err == nil {
			t.Errorf("%s: expected validation error, got nil", c.name)
		}
	}
}

// TestEntityTaxonomyQueryIntegration runs the fluent query against a real store.
func TestEntityTaxonomyQueryIntegration(t *testing.T) {
	db := InitDB("entity_taxonomy_fluent_test")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "fluent_entity",
		AttributeTableName: "fluent_attribute",
		TaxonomiesEnabled:  true,
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	taxonomy, err := store.TaxonomyCreateByOptions(ctx, TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})
	if err != nil {
		t.Fatal(err)
	}
	term, err := store.TaxonomyTermCreateByOptions(ctx, TaxonomyTermOptions{
		TaxonomyID: taxonomy.ID(),
		Name:       "Electronics",
		Slug:       "electronics",
	})
	if err != nil {
		t.Fatal(err)
	}

	entity1, err := store.EntityCreateWithType(ctx, "product")
	if err != nil {
		t.Fatal(err)
	}
	entity2, err := store.EntityCreateWithType(ctx, "product")
	if err != nil {
		t.Fatal(err)
	}

	for _, entity := range []EntityInterface{entity1, entity2} {
		if err := store.EntityTaxonomyAssign(ctx, entity.ID(), taxonomy.ID(), term.ID()); err != nil {
			t.Fatal(err)
		}
	}

	query := EntityTaxonomyQuery().WithEntityID(entity1.ID())
	if err := query.Validate(); err != nil {
		t.Fatal("Validate failed:", err)
	}

	list, err := store.EntityTaxonomyList(ctx, query)
	if err != nil {
		t.Fatal("EntityTaxonomyList failed:", err)
	}
	if len(list) != 1 {
		t.Fatalf("Expected 1 assignment, got %d", len(list))
	}

	count, err := store.EntityTaxonomyCount(ctx, EntityTaxonomyQuery().WithTermID(term.ID()))
	if err != nil {
		t.Fatal("EntityTaxonomyCount failed:", err)
	}
	if count != 2 {
		t.Fatalf("Expected count 2, got %d", count)
	}
}
