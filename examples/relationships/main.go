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

	// Open SQLite database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close() //nolint:errcheck

	// Create the store with relationships enabled
	store, err := entitystore.NewStore(entitystore.NewStoreOptions{
		DB:                   db,
		EntityTableName:      "entities",
		AttributeTableName:   "attributes",
		RelationshipsEnabled: true,
		AutomigrateEnabled:   true,
	})
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	fmt.Println("=== EntityStore Relationships Example ===")

	// Create Author entity
	fmt.Println("1. Creating Author entity...")
	author, err := store.EntityCreateWithTypeAndAttributes(ctx, "author", map[string]string{
		"name": "Jane Smith",
		"bio":  "Bestselling author of fiction novels",
	})
	if err != nil {
		log.Fatalf("Failed to create author: %v", err)
	}
	fmt.Printf("   Created author: %s (ID: %s)\n", "Jane Smith", author.ID())

	// Create Book entities
	fmt.Println("\n2. Creating Book entities...")
	book1, err := store.EntityCreateWithTypeAndAttributes(ctx, "book", map[string]string{
		"title": "The Mystery",
		"isbn":  "978-1234567890",
	})
	if err != nil {
		log.Fatalf("Failed to create book: %v", err)
	}
	book2, err := store.EntityCreateWithTypeAndAttributes(ctx, "book", map[string]string{
		"title": "The Adventure",
		"isbn":  "978-0987654321",
	})
	if err != nil {
		log.Fatalf("Failed to create book: %v", err)
	}
	fmt.Printf("   Book 1: %s (ID: %s)\n", "The Mystery", book1.ID())
	fmt.Printf("   Book 2: %s (ID: %s)\n", "The Adventure", book2.ID())

	// Create BELONGS_TO relationship (book belongs to author)
	fmt.Println("\n3. Creating BELONGS_TO relationships (books → author)...")
	rel1, err := store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
		EntityID:         book1.ID(),
		RelatedEntityID:  author.ID(),
		RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
	})
	if err != nil {
		log.Fatalf("Failed to create relationship: %v", err)
	}
	fmt.Printf("   Created relationship: %s\n", rel1.ID())

	rel2, err := store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
		EntityID:         book2.ID(),
		RelatedEntityID:  author.ID(),
		RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
	})
	if err != nil {
		log.Fatalf("Failed to create relationship: %v", err)
	}
	fmt.Printf("   Created relationship: %s\n", rel2.ID())

	// Create Category entity
	fmt.Println("\n4. Creating Category entity...")
	fictionCategory, err := store.EntityCreateWithTypeAndAttributes(ctx, "category", map[string]string{
		"name": "Fiction",
		"slug": "fiction",
	})
	if err != nil {
		log.Fatalf("Failed to create category: %v", err)
	}
	fmt.Printf("   Category: %s (ID: %s)\n", fictionCategory.GetTempKey("name"), fictionCategory.ID())

	// Create MANY_TO_MANY relationship (books <-> categories)
	fmt.Println("\n5. Creating MANY_TO_MANY relationships (books ↔ category)...")
	rel3, err := store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
		EntityID:         book1.ID(),
		RelatedEntityID:  fictionCategory.ID(),
		RelationshipType: entitystore.RELATIONSHIP_TYPE_MANY_MANY,
	})
	if err != nil {
		log.Fatalf("Failed to create relationship: %v", err)
	}
	rel4, err := store.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
		EntityID:         book2.ID(),
		RelatedEntityID:  fictionCategory.ID(),
		RelationshipType: entitystore.RELATIONSHIP_TYPE_MANY_MANY,
	})
	if err != nil {
		log.Fatalf("Failed to create relationship: %v", err)
	}
	fmt.Printf("   Book 1 → Category: %s\n", rel3.ID())
	fmt.Printf("   Book 2 → Category: %s\n", rel4.ID())

	// Query relationships
	fmt.Println("\n6. Querying relationships for author (finding their books)...")
	relationships, err := store.RelationshipList(ctx, entitystore.RelationshipQueryOptions{
		RelatedEntityID:  author.ID(),
		RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
	})
	if err != nil {
		log.Fatalf("Failed to list relationships: %v", err)
	}
	fmt.Printf("   Found %d books belonging to author:\n", len(relationships))
	for _, rel := range relationships {
		book, err := store.EntityFindByID(ctx, rel.GetEntityID())
		if err != nil {
			log.Fatalf("Failed to find book: %v", err)
		}
		if book == nil {
			log.Printf("Book not found for entity ID: %s", rel.GetEntityID())
			continue
		}
		fmt.Printf("   - %s\n", book.GetTempKey("title"))
	}

	// Query reverse relationships
	fmt.Println("\n7. Querying books and their categories...")
	book1Categories, err := store.RelationshipList(ctx, entitystore.RelationshipQueryOptions{
		EntityID:         book1.ID(),
		RelationshipType: entitystore.RELATIONSHIP_TYPE_MANY_MANY,
	})
	if err != nil {
		log.Fatalf("Failed to list relationships: %v", err)
	}
	fmt.Printf("   Book '%s' is in %d categories\n",
		book1.GetTempKey("title"), len(book1Categories))

	// Count relationships
	fmt.Println("\n8. Counting relationships...")
	count, err := store.RelationshipCount(ctx, entitystore.RelationshipQueryOptions{
		RelationshipType: entitystore.RELATIONSHIP_TYPE_BELONGS_TO,
	})
	if err != nil {
		log.Fatalf("Failed to count relationships: %v", err)
	}
	fmt.Printf("   Total BELONGS_TO relationships: %d\n", count)

	// Find specific relationship
	fmt.Println("\n9. Finding specific relationship between book and author...")
	foundRel, err := store.RelationshipFindByEntities(ctx,
		book1.ID(),
		author.ID(),
		entitystore.RELATIONSHIP_TYPE_BELONGS_TO)
	if err != nil {
		fmt.Printf("   Relationship not found: %v\n", err)
	} else if foundRel == nil {
		fmt.Println("   Relationship not found (nil)")
	} else {
		fmt.Printf("   Found relationship: %s\n", foundRel.ID())
	}

	// Soft delete relationship (trash)
	fmt.Println("\n10. Soft deleting a relationship...")
	trashed, err := store.RelationshipTrash(ctx, rel4.ID(), "admin")
	if err != nil {
		log.Fatalf("Failed to trash relationship: %v", err)
	}
	fmt.Printf("   Deleted: %v\n", trashed)

	// Count after delete
	count, err = store.RelationshipCount(ctx, entitystore.RelationshipQueryOptions{})
	if err != nil {
		log.Fatalf("Failed to count relationships: %v", err)
	}
	fmt.Printf("    Total relationships after delete: %d\n", count)

	fmt.Println("\n=== Example completed successfully! ===")
}
