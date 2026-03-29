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
	defer db.Close()

	// Create the store
	store, err := entitystore.NewStore(entitystore.NewStoreOptions{
		DB:                 db,
		EntityTableName:    "entities",
		AttributeTableName: "attributes",
		AutomigrateEnabled: true,
	})
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	fmt.Println("=== EntityStore Basic Example ===")

	// Create a person entity
	fmt.Println("1. Creating a person entity...")
	person, err := store.EntityCreateWithTypeAndAttributes(ctx, "person", map[string]string{
		"name":  "John Doe",
		"age":   "30",
		"email": "john@example.com",
	})
	if err != nil {
		log.Fatalf("Failed to create person: %v", err)
	}
	fmt.Printf("   Created person with ID: %s\n", person.ID())

	// Retrieve attributes from store
	nameAttr, _ := store.AttributeFind(ctx, person.ID(), "name")
	ageAttr, _ := store.AttributeFind(ctx, person.ID(), "age")

	fmt.Printf("   Name: %s\n", nameAttr.GetAttributeValue())
	fmt.Printf("   Age: %s\n", ageAttr.GetAttributeValue())

	// Create a product entity
	fmt.Println("\n2. Creating a product entity...")
	product, err := store.EntityCreateWithTypeAndAttributes(ctx, "product", map[string]string{
		"name":  "Laptop",
		"price": "999.99",
		"sku":   "LAPTOP-001",
	})
	if err != nil {
		log.Fatalf("Failed to create product: %v", err)
	}
	fmt.Printf("   Created product with ID: %s\n", product.ID())

	// Retrieve attributes from store
	prodNameAttr, _ := store.AttributeFind(ctx, product.ID(), "name")
	prodPriceAttr, _ := store.AttributeFind(ctx, product.ID(), "price")

	fmt.Printf("   Name: %s\n", prodNameAttr.GetAttributeValue())
	fmt.Printf("   Price: %s\n", prodPriceAttr.GetAttributeValue())

	// List all entities
	fmt.Println("\n3. Listing all entities...")
	entities, err := store.EntityList(ctx, entitystore.EntityQueryOptions{})
	if err != nil {
		log.Fatalf("Failed to list entities: %v", err)
	}
	fmt.Printf("   Found %d entities:\n", len(entities))
	for _, e := range entities {
		fmt.Printf("   - [%s] %s (handle: %s)\n", e.GetEntityType(), e.ID(), e.GetEntityHandle())
	}

	// Find entity by ID
	fmt.Println("\n4. Finding entity by ID...")
	found, err := store.EntityFindByID(ctx, person.ID())
	if err != nil {
		log.Fatalf("Failed to find entity: %v", err)
	}
	foundNameAttr, _ := store.AttributeFind(ctx, found.ID(), "name")
	fmt.Printf("   Found: %s (type: %s)\n", foundNameAttr.GetAttributeValue(), found.GetEntityType())

	// Update entity attributes via store
	fmt.Println("\n5. Updating entity attributes...")
	err = store.AttributeSetString(ctx, person.ID(), "age", "31")
	if err != nil {
		log.Fatalf("Failed to update age: %v", err)
	}
	err = store.AttributeSetString(ctx, person.ID(), "city", "New York")
	if err != nil {
		log.Fatalf("Failed to add city: %v", err)
	}
	fmt.Println("   Updated person age to 31 and added city: New York")

	// List attributes
	fmt.Println("\n6. Listing all attributes for person...")
	attrs, err := store.EntityAttributeList(ctx, person.ID())
	if err != nil {
		log.Fatalf("Failed to list attributes: %v", err)
	}
	fmt.Printf("   Found %d attributes:\n", len(attrs))
	for _, attr := range attrs {
		fmt.Printf("   - %s: %s\n", attr.GetAttributeKey(), attr.GetAttributeValue())
	}

	// Count entities
	fmt.Println("\n7. Counting entities...")
	count, err := store.EntityCount(ctx, entitystore.EntityQueryOptions{})
	if err != nil {
		log.Fatalf("Failed to count entities: %v", err)
	}
	fmt.Printf("   Total entities: %d\n", count)

	// Count by type
	personCount, _ := store.EntityCount(ctx, entitystore.EntityQueryOptions{EntityType: "person"})
	productCount, _ := store.EntityCount(ctx, entitystore.EntityQueryOptions{EntityType: "product"})
	fmt.Printf("   Persons: %d, Products: %d\n", personCount, productCount)

	// Soft delete (trash)
	fmt.Println("\n8. Soft deleting product...")
	deleted, err := store.EntityTrash(ctx, product.ID())
	if err != nil {
		log.Fatalf("Failed to trash entity: %v", err)
	}
	fmt.Printf("   Deleted: %v\n", deleted)

	// Count after delete
	count, _ = store.EntityCount(ctx, entitystore.EntityQueryOptions{})
	fmt.Printf("   Total entities after delete: %d\n", count)

	fmt.Println("\n=== Example completed successfully! ===")
}
