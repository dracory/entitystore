package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/dracory/entitystore"
	_ "modernc.org/sqlite"
)

func main() {
	ctx := context.Background()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer func() { _ = db.Close() }()

	store, err := entitystore.NewStore(entitystore.NewStoreOptions{
		DB:                 db,
		EntityTableName:    "entities",
		AttributeTableName: "attributes",
		AutomigrateEnabled: true,
	})
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	fmt.Println("=== EntityStore Aggregate Example ===")

	// Seed blogposts with a publishedDay attribute
	fmt.Println("1. Creating blogposts...")
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
			log.Fatalf("Failed to create blogpost: %v", err)
		}
	}
	fmt.Printf("   Created %d blogposts\n", len(days))

	// Count posts per published day — the YesSql BlogPostByDay index
	// expressed as a live GROUP BY instead of a maintained table
	fmt.Println("\n2. Posts per day (COUNT grouped by publishedDay)...")
	byDay, err := store.AttributeGroupBy(ctx, entitystore.AttributeQuery().
		WithAttributeKey("publishedDay").
		WithEntityType("blogpost").
		WithAggregate(entitystore.AGGREGATE_COUNT))
	if err != nil {
		log.Fatalf("Failed to group by day: %v", err)
	}
	for day, count := range byDay {
		fmt.Printf("   %s: %v posts\n", day, count)
	}

	// Count posts per status
	fmt.Println("\n3. Posts per status (COUNT grouped by status)...")
	byStatus, err := store.AttributeGroupBy(ctx, entitystore.AttributeQuery().
		WithAttributeKey("status").
		WithEntityType("blogpost").
		WithAggregate(entitystore.AGGREGATE_COUNT))
	if err != nil {
		log.Fatalf("Failed to group by status: %v", err)
	}
	for status, count := range byStatus {
		fmt.Printf("   %s: %v posts\n", status, count)
	}

	// Distinct published days
	fmt.Println("\n4. Days with posts (DISTINCT publishedDay)...")
	distinctDays, err := store.AttributeGroupBy(ctx, entitystore.AttributeQuery().
		WithAttributeKey("publishedDay").
		WithEntityType("blogpost").
		WithAggregate(entitystore.AGGREGATE_DISTINCT))
	if err != nil {
		log.Fatalf("Failed to list distinct days: %v", err)
	}
	for day := range distinctDays {
		fmt.Printf("   %s\n", day)
	}

	fmt.Println("\n=== Example completed successfully! ===")
}
