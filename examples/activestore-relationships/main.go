package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/dracory/entitystore"
	"github.com/dracory/entitystore/activestore"
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
		DB:                   db,
		EntityTableName:      "entities",
		AttributeTableName:   "attributes",
		AutomigrateEnabled:   true,
		RelationshipsEnabled: true, // required for relationships
	})
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	active, err := activestore.New(ctx, store)
	if err != nil {
		log.Fatalf("Failed to create active store: %v", err)
	}

	fmt.Println("=== ActiveStore Relationships Example ===")

	// Create two entities and link them — the entity's own ID is implicit
	fmt.Println("1. Creating and linking entities...")
	post := active.EntityCreate("post").SetString("title", "Hello World")
	author := active.EntityCreate("author").SetString("name", "Ada")
	for _, e := range []activestore.ActiveEntityInterface{post, author} {
		if _, err := e.Save(); err != nil {
			log.Fatalf("Failed to save: %v", err)
		}
	}

	if err := post.RelateTo(author, "written_by"); err != nil {
		log.Fatalf("Failed to relate: %v", err)
	}
	fmt.Println("   Linked post --written_by--> author")

	// Related() hydrates the linked entities in one call — no manual
	// RelationshipList + EntityFindByID loop
	fmt.Println("\n2. Navigating the relationship...")
	authors, err := post.Related("written_by")
	if err != nil {
		log.Fatalf("Failed to list related: %v", err)
	}
	for _, a := range authors {
		name, _, _ := a.GetString("name")
		fmt.Printf("   written_by: %s\n", name)
	}

	// Remove the link
	fmt.Println("\n3. Unrelating...")
	if err := post.Unrelate(author, "written_by"); err != nil {
		log.Fatalf("Failed to unrelate: %v", err)
	}
	authors, _ = post.Related("written_by")
	fmt.Printf("   written_by count after unrelate: %d\n", len(authors))

	fmt.Println("\n=== Example completed successfully! ===")
}
