package entitystore

import (
	"context"
	"testing"
)

// TestTaxonomyTermQueryFluentOptions verifies that the fluent builder
// records the explicitly set fields and tracks presence via has* flags.
func TestTaxonomyTermQueryFluentOptions(t *testing.T) {
	query := TaxonomyTermQuery().
		WithTaxonomyID("taxonomy-1").
		WithSlug("electronics").
		WithParentID("parent-1").
		WithLimit(100).
		WithSortOrder("desc")

	if query.GetTaxonomyID() != "taxonomy-1" {
		t.Errorf("TaxonomyID: expected 'taxonomy-1', got %q", query.GetTaxonomyID())
	}
	if query.GetSlug() != "electronics" {
		t.Errorf("Slug: expected 'electronics', got %q", query.GetSlug())
	}
	if query.GetParentID() != "parent-1" {
		t.Errorf("ParentID: expected 'parent-1', got %q", query.GetParentID())
	}
	if query.GetLimit() != 100 {
		t.Errorf("Limit: expected 100, got %d", query.GetLimit())
	}
	if query.GetSortOrder() != "desc" {
		t.Errorf("SortOrder: expected 'desc', got %q", query.GetSortOrder())
	}

	// Presence flags: set fields report true, unset fields report false.
	if !query.HasTaxonomyID() || !query.HasSlug() || !query.HasParentID() ||
		!query.HasLimit() || !query.HasSortOrder() {
		t.Error("Set fields should report Has*() = true")
	}
	if query.HasID() || query.HasIDs() ||
		query.HasOffset() || query.HasSortBy() || query.HasCountOnly() {
		t.Error("Unset fields should report Has*() = false")
	}
}

// TestTaxonomyTermQueryValidate verifies Validate() catches set-but-invalid
// values and ignores unset fields.
func TestTaxonomyTermQueryValidate(t *testing.T) {
	// Empty builder is valid
	if err := TaxonomyTermQuery().Validate(); err != nil {
		t.Error("Empty query should be valid:", err)
	}

	cases := []struct {
		name  string
		query TaxonomyTermQueryInterface
	}{
		{"empty id", TaxonomyTermQuery().WithID("")},
		{"empty ids", TaxonomyTermQuery().WithIDs([]string{})},
		{"empty taxonomy id", TaxonomyTermQuery().WithTaxonomyID("")},
		{"empty slug", TaxonomyTermQuery().WithSlug("")},
		{"empty parent id", TaxonomyTermQuery().WithParentID("")},
		{"empty sort by", TaxonomyTermQuery().WithSortBy("")},
		{"bad sort order", TaxonomyTermQuery().WithSortOrder("sideways")},
	}
	for _, c := range cases {
		if err := c.query.Validate(); err == nil {
			t.Errorf("%s: expected validation error, got nil", c.name)
		}
	}
}

// TestTaxonomyTermQueryIntegration runs the fluent query against a real store.
func TestTaxonomyTermQueryIntegration(t *testing.T) {
	db := InitDB("taxonomy_term_fluent_test")

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

	for _, slug := range []string{"electronics", "books", "toys"} {
		_, err := store.TaxonomyTermCreateByOptions(ctx, TaxonomyTermOptions{
			TaxonomyID: taxonomy.ID(),
			Name:       slug,
			Slug:       slug,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	query := TaxonomyTermQuery().
		WithTaxonomyID(taxonomy.ID()).
		WithSlug("electronics")
	if err := query.Validate(); err != nil {
		t.Fatal("Validate failed:", err)
	}

	list, err := store.TaxonomyTermList(ctx, query)
	if err != nil {
		t.Fatal("TaxonomyTermList failed:", err)
	}
	if len(list) != 1 {
		t.Fatalf("Expected 1 term, got %d", len(list))
	}

	count, err := store.TaxonomyTermCount(ctx, TaxonomyTermQuery().WithTaxonomyID(taxonomy.ID()))
	if err != nil {
		t.Fatal("TaxonomyTermCount failed:", err)
	}
	if count != 3 {
		t.Fatalf("Expected count 3, got %d", count)
	}
}
