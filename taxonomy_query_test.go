package entitystore

import (
	"context"
	"testing"
)

// TestTaxonomyQueryFluentOptions verifies that the fluent builder records
// the explicitly set fields and tracks presence via has* flags.
func TestTaxonomyQueryFluentOptions(t *testing.T) {
	query := TaxonomyQuery().
		WithSlug("categories").
		WithEntityTypes([]string{"product", "post"}).
		WithLimit(100).
		WithSortOrder("desc")

	if query.GetSlug() != "categories" {
		t.Errorf("Slug: expected 'categories', got %q", query.GetSlug())
	}
	types := query.GetEntityTypes()
	if len(types) != 2 || types[0] != "product" || types[1] != "post" {
		t.Errorf("EntityTypes mismatch: %v", types)
	}
	if query.GetLimit() != 100 {
		t.Errorf("Limit: expected 100, got %d", query.GetLimit())
	}
	if query.GetSortOrder() != "desc" {
		t.Errorf("SortOrder: expected 'desc', got %q", query.GetSortOrder())
	}

	// Presence flags: set fields report true, unset fields report false.
	if !query.HasSlug() || !query.HasEntityTypes() || !query.HasLimit() || !query.HasSortOrder() {
		t.Error("Set fields should report Has*() = true")
	}
	if query.HasID() || query.HasIDs() || query.HasParentID() ||
		query.HasOffset() || query.HasSortBy() || query.HasCountOnly() {
		t.Error("Unset fields should report Has*() = false")
	}
}

// TestTaxonomyQueryValidate verifies Validate() catches set-but-invalid
// values and ignores unset fields.
func TestTaxonomyQueryValidate(t *testing.T) {
	// Empty builder is valid
	if err := TaxonomyQuery().Validate(); err != nil {
		t.Error("Empty query should be valid:", err)
	}

	cases := []struct {
		name  string
		query TaxonomyQueryInterface
	}{
		{"empty id", TaxonomyQuery().WithID("")},
		{"empty ids", TaxonomyQuery().WithIDs([]string{})},
		{"empty slug", TaxonomyQuery().WithSlug("")},
		{"empty parent id", TaxonomyQuery().WithParentID("")},
		{"empty entity types", TaxonomyQuery().WithEntityTypes(nil)},
		{"empty sort by", TaxonomyQuery().WithSortBy("")},
		{"bad sort order", TaxonomyQuery().WithSortOrder("sideways")},
	}
	for _, c := range cases {
		if err := c.query.Validate(); err == nil {
			t.Errorf("%s: expected validation error, got nil", c.name)
		}
	}
}

// TestTaxonomyQueryIntegration runs the fluent query against a real store.
func TestTaxonomyQueryIntegration(t *testing.T) {
	db := InitDB("taxonomy_fluent_test")

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
	for _, slug := range []string{"categories", "tags"} {
		_, err := store.TaxonomyCreateByOptions(ctx, TaxonomyOptions{
			Name: slug,
			Slug: slug,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	query := TaxonomyQuery().WithSlug("categories")
	if err := query.Validate(); err != nil {
		t.Fatal("Validate failed:", err)
	}

	list, err := store.TaxonomyList(ctx, query)
	if err != nil {
		t.Fatal("TaxonomyList failed:", err)
	}
	if len(list) != 1 {
		t.Fatalf("Expected 1 taxonomy, got %d", len(list))
	}

	count, err := store.TaxonomyCount(ctx, TaxonomyQuery())
	if err != nil {
		t.Fatal("TaxonomyCount failed:", err)
	}
	if count != 2 {
		t.Fatalf("Expected count 2, got %d", count)
	}
}
