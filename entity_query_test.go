package entitystore

import (
	"context"
	"testing"
)

// TestEntityQueryFluentOptions verifies that the fluent builder records
// the explicitly set fields and tracks presence via has* flags.
func TestEntityQueryFluentOptions(t *testing.T) {
	query := EntityQuery().
		WithEntityType("character").
		WithEntityHandle("ada").
		WithLimit(100).
		WithSortOrder("desc")

	if query.GetEntityType() != "character" {
		t.Errorf("EntityType: expected 'character', got %q", query.GetEntityType())
	}
	if query.GetEntityHandle() != "ada" {
		t.Errorf("EntityHandle: expected 'ada', got %q", query.GetEntityHandle())
	}
	if query.GetLimit() != 100 {
		t.Errorf("Limit: expected 100, got %d", query.GetLimit())
	}
	if query.GetSortOrder() != "desc" {
		t.Errorf("SortOrder: expected 'desc', got %q", query.GetSortOrder())
	}

	// Presence flags: set fields report true, unset fields report false.
	if !query.HasEntityType() || !query.HasEntityHandle() || !query.HasLimit() || !query.HasSortOrder() {
		t.Error("Set fields should report Has*() = true")
	}
	if query.HasID() || query.HasIDs() || query.HasSearch() ||
		query.HasOffset() || query.HasSortBy() || query.HasCountOnly() {
		t.Error("Unset fields should report Has*() = false")
	}
}

// TestEntityQueryValidate verifies Validate() catches set-but-invalid
// values and ignores unset fields.
func TestEntityQueryValidate(t *testing.T) {
	// Empty builder is valid
	if err := EntityQuery().Validate(); err != nil {
		t.Error("Empty query should be valid:", err)
	}

	cases := []struct {
		name  string
		query EntityQueryInterface
	}{
		{"empty id", EntityQuery().WithID("")},
		{"empty ids", EntityQuery().WithIDs([]string{})},
		{"empty entity type", EntityQuery().WithEntityType("")},
		{"empty entity handle", EntityQuery().WithEntityHandle("")},
		{"empty sort by", EntityQuery().WithSortBy("")},
		{"bad sort order", EntityQuery().WithSortOrder("sideways")},
	}
	for _, c := range cases {
		if err := c.query.Validate(); err == nil {
			t.Errorf("%s: expected validation error, got nil", c.name)
		}
	}
}

// TestEntityQueryIntegration runs the fluent query against a real store.
func TestEntityQueryIntegration(t *testing.T) {
	db := InitDB("entity_fluent_test")

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
	if _, err := store.EntityCreateWithType(ctx, "character"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.EntityCreateWithType(ctx, "character"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.EntityCreateWithType(ctx, "product"); err != nil {
		t.Fatal(err)
	}

	query := EntityQuery().WithEntityType("character")
	if err := query.Validate(); err != nil {
		t.Fatal("Validate failed:", err)
	}

	list, err := store.EntityList(ctx, query)
	if err != nil {
		t.Fatal("EntityList failed:", err)
	}
	if len(list) != 2 {
		t.Fatalf("Expected 2 entities, got %d", len(list))
	}

	count, err := store.EntityCount(ctx, query)
	if err != nil {
		t.Fatal("EntityCount failed:", err)
	}
	if count != 2 {
		t.Fatalf("Expected count 2, got %d", count)
	}
}

func TestEntityQuery_PrefetchAttributes(t *testing.T) {
	db := InitDB("test_entity_query_prefetch.db")
	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "entity",
		AttributeTableName: "attribute",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	e1, err := store.EntityCreateWithType(ctx, "product")
	if err != nil {
		t.Fatal(err)
	}
	e2, err := store.EntityCreateWithType(ctx, "product")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.AttributeSetString(ctx, e1.ID(), "name", "Laptop"); err != nil {
		t.Fatal(err)
	}
	if err := store.AttributeSetString(ctx, e1.ID(), "sku", "LAP-1"); err != nil {
		t.Fatal(err)
	}
	if err := store.AttributeSetString(ctx, e2.ID(), "name", "Phone"); err != nil {
		t.Fatal(err)
	}

	list, err := store.EntityList(ctx, EntityQuery().
		WithEntityType("product").
		WithPrefetchAttributes("name"))
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Fatalf("Expected 2 entities, got %d", len(list))
	}

	byID := map[string]EntityInterface{list[0].ID(): list[0], list[1].ID(): list[1]}
	gotE1, ok := byID[e1.ID()]
	if !ok {
		t.Fatalf("Expected entity %s in result", e1.ID())
	}
	gotE2, ok := byID[e2.ID()]
	if !ok {
		t.Fatalf("Expected entity %s in result", e2.ID())
	}
	if got := gotE1.GetTempKey("name"); got != "Laptop" {
		t.Fatalf("Expected prefetched name Laptop, got %q", got)
	}
	if got := gotE2.GetTempKey("name"); got != "Phone" {
		t.Fatalf("Expected prefetched name Phone, got %q", got)
	}
	// sku was not in the prefetch list — must not be loaded
	if got := gotE1.GetTempKey("sku"); got != "" {
		t.Fatalf("Expected sku to not be prefetched, got %q", got)
	}

}

func TestEntityQuery_PrefetchAll(t *testing.T) {
	db := InitDB("test_entity_query_prefetchall.db")
	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "entity",
		AttributeTableName: "attribute",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	e1, err := store.EntityCreateWithType(ctx, "product")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.AttributeSetString(ctx, e1.ID(), "name", "Laptop"); err != nil {
		t.Fatal(err)
	}
	if err := store.AttributeSetString(ctx, e1.ID(), "sku", "LAP-1"); err != nil {
		t.Fatal(err)
	}

	// no keys = prefetch all attributes
	list, err := store.EntityList(ctx, EntityQuery().
		WithEntityType("product").
		WithPrefetchAttributes())
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("Expected 1 entity, got %d", len(list))
	}
	// all attributes must be loaded — both name and sku
	if got := list[0].GetTempKey("name"); got != "Laptop" {
		t.Fatalf("Expected prefetched name Laptop, got %q", got)
	}
	if got := list[0].GetTempKey("sku"); got != "LAP-1" {
		t.Fatalf("Expected prefetched sku LAP-1, got %q", got)
	}

	// not setting prefetch at all must not prefetch
	list2, err := store.EntityList(ctx, EntityQuery().WithEntityType("product"))
	if err != nil {
		t.Fatal(err)
	}
	if len(list2) != 1 {
		t.Fatalf("Expected 1 entity, got %d", len(list2))
	}
	if got := list2[0].GetTempKey("name"); got != "" {
		t.Fatalf("Expected no prefetch, got %q", got)
	}
}
