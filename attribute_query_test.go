package entitystore

import (
	"context"
	"testing"
)

// TestAttributeQueryFluentOptions verifies that the fluent builder records
// the explicitly set fields and tracks presence via has* flags.
func TestAttributeQueryFluentOptions(t *testing.T) {
	query := AttributeQuery().
		WithEntityType("character").
		WithAttributeKeys([]string{"name", "title"}).
		WithLimit(100).
		WithSortOrder("desc")

	if query.GetEntityType() != "character" {
		t.Errorf("EntityType: expected 'character', got %q", query.GetEntityType())
	}
	keys := query.GetAttributeKeys()
	if len(keys) != 2 || keys[0] != "name" || keys[1] != "title" {
		t.Errorf("AttributeKeys mismatch: %v", keys)
	}
	if query.GetLimit() != 100 {
		t.Errorf("Limit: expected 100, got %d", query.GetLimit())
	}
	if query.GetSortOrder() != "desc" {
		t.Errorf("SortOrder: expected 'desc', got %q", query.GetSortOrder())
	}

	// Presence flags: set fields report true, unset fields report false.
	if !query.HasEntityType() || !query.HasAttributeKeys() || !query.HasLimit() || !query.HasSortOrder() {
		t.Error("Set fields should report Has*() = true")
	}
	if query.HasID() || query.HasEntityID() || query.HasAttributeKey() ||
		query.HasOffset() || query.HasSortBy() || query.HasCountOnly() {
		t.Error("Unset fields should report Has*() = false")
	}
}

// TestAttributeQueryValidate verifies Validate() catches set-but-invalid
// values and ignores unset fields.
func TestAttributeQueryValidate(t *testing.T) {
	// Empty builder is valid
	if err := AttributeQuery().Validate(); err != nil {
		t.Error("Empty query should be valid:", err)
	}

	cases := []struct {
		name  string
		query AttributeQueryInterface
	}{
		{"empty id", AttributeQuery().WithID("")},
		{"empty ids", AttributeQuery().WithIDs([]string{})},
		{"empty entity type", AttributeQuery().WithEntityType("")},
		{"empty key", AttributeQuery().WithAttributeKey("")},
		{"empty keys", AttributeQuery().WithAttributeKeys(nil)},
		{"bad sort order", AttributeQuery().WithSortOrder("sideways")},
	}
	for _, c := range cases {
		if err := c.query.Validate(); err == nil {
			t.Errorf("%s: expected validation error, got nil", c.name)
		}
	}
}

// TestAttributeQueryIntegration runs the fluent query against a real store.
func TestAttributeQueryIntegration(t *testing.T) {
	db := InitDB("attr_fluent_test")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "fluent_entity",
		AttributeTableName: "fluent_attribute",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	entity, err := store.EntityCreateWithType(ctx, "character")
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range map[string]string{"name": "Ada", "title": "Countess", "secret": "x"} {
		if err := store.AttributeSetString(ctx, entity.ID(), k, v); err != nil {
			t.Fatal(err)
		}
	}

	query := AttributeQuery().
		WithEntityType("character").
		WithAttributeKeys([]string{"name", "title"})
	if err := query.Validate(); err != nil {
		t.Fatal("Validate failed:", err)
	}

	list, err := store.AttributeList(ctx, query)
	if err != nil {
		t.Fatal("AttributeList failed:", err)
	}
	if len(list) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(list))
	}
}

// TestAttributeQueryTimeFilters verifies created_at/updated_at range
// filtering executes against a real store.
func TestAttributeQueryTimeFilters(t *testing.T) {
	db := InitDB("attr_time_test")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "time_entity",
		AttributeTableName: "time_attribute",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	entity, err := store.EntityCreateWithType(ctx, "character")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.AttributeSetString(ctx, entity.ID(), "name", "Ada"); err != nil {
		t.Fatal(err)
	}

	// Future lower bound excludes everything
	list, err := store.AttributeList(ctx, AttributeQuery().WithCreatedAtGte("2999-01-01 00:00:00"))
	if err != nil {
		t.Fatal("AttributeList failed:", err)
	}
	if len(list) != 0 {
		t.Fatalf("Expected 0 attributes after future created_at_gte, got %d", len(list))
	}

	// Past bounds include the row; combined range works too
	list, err = store.AttributeList(ctx, AttributeQuery().
		WithCreatedAtGte("2000-01-01 00:00:00").
		WithCreatedAtLte("2999-01-01 00:00:00").
		WithUpdatedAtGte("2000-01-01 00:00:00"))
	if err != nil {
		t.Fatal("AttributeList failed:", err)
	}
	if len(list) != 1 {
		t.Fatalf("Expected 1 attribute in range, got %d", len(list))
	}

	// Invalid: set-but-empty bound rejected by Validate
	if _, err := store.AttributeList(ctx, AttributeQuery().WithCreatedAtGte("")); err == nil {
		t.Error("Expected validation error for empty created_at_gte")
	}
}
