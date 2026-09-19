package main

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dracory/entitystore"
	_ "modernc.org/sqlite"
)

func setupTestStore(t *testing.T) (entitystore.StoreInterface, func()) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	store, err := entitystore.NewStore(entitystore.NewStoreOptions{
		DB:                 db,
		EntityTableName:    "test_entities",
		AttributeTableName: "test_attributes",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	return store, func() { _ = db.Close() }
}

func seedPosts(t *testing.T, store entitystore.StoreInterface) {
	ctx := context.Background()
	days := []string{"2026-09-17", "2026-09-18", "2026-09-18", "2026-09-19", "2026-09-19", "2026-09-19"}
	for i, day := range days {
		status := "published"
		if i == 0 {
			status = "draft"
		}
		_, err := store.EntityCreateWithTypeAndAttributes(ctx, "blogpost", map[string]string{
			"publishedDay": day,
			"status":       status,
		})
		if err != nil {
			t.Fatalf("Failed to create blogpost: %v", err)
		}
	}
}

func TestGroupByCount(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()
	seedPosts(t, store)

	byDay, err := store.AttributeGroupBy(context.Background(), entitystore.AttributeQuery().
		WithAttributeKey("publishedDay").
		WithEntityType("blogpost").
		WithAggregate(entitystore.AGGREGATE_COUNT))
	if err != nil {
		t.Fatalf("AttributeGroupBy failed: %v", err)
	}

	expected := map[string]int64{"2026-09-17": 1, "2026-09-18": 2, "2026-09-19": 3}
	if len(byDay) != len(expected) {
		t.Fatalf("Expected %d groups, got %d", len(expected), len(byDay))
	}
	for day, want := range expected {
		if byDay[day] != want {
			t.Errorf("Day %s: expected %d, got %v", day, want, byDay[day])
		}
	}
}

func TestGroupByCountWithoutEntityType(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()
	seedPosts(t, store)

	byStatus, err := store.AttributeGroupBy(context.Background(), entitystore.AttributeQuery().
		WithAttributeKey("status").
		WithAggregate(entitystore.AGGREGATE_COUNT))
	if err != nil {
		t.Fatalf("AttributeGroupBy failed: %v", err)
	}

	if byStatus["published"] != int64(5) {
		t.Errorf("Expected 5 published, got %v", byStatus["published"])
	}
	if byStatus["draft"] != int64(1) {
		t.Errorf("Expected 1 draft, got %v", byStatus["draft"])
	}
}

func TestGroupByDistinct(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()
	seedPosts(t, store)

	days, err := store.AttributeGroupBy(context.Background(), entitystore.AttributeQuery().
		WithAttributeKey("publishedDay").
		WithAggregate(entitystore.AGGREGATE_DISTINCT))
	if err != nil {
		t.Fatalf("AttributeGroupBy failed: %v", err)
	}

	if len(days) != 3 {
		t.Fatalf("Expected 3 distinct days, got %d", len(days))
	}
}

func TestGroupByInvalidAggregate(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	_, err := store.AttributeGroupBy(context.Background(), entitystore.AttributeQuery().
		WithAttributeKey("publishedDay").
		WithAggregate(entitystore.AggregateFunc("bogus")))
	if err == nil {
		t.Fatal("Expected validation error for unknown aggregate")
	}
}
