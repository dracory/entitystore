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
	nameAttr, err := store.AttributeFind(ctx, person.ID(), "name")
	if err != nil {
		log.Fatalf("Failed to find name attribute: %v", err)
	}
	if nameAttr == nil {
		log.Fatal("Expected name attribute to be found")
	}

	ageAttr, err := store.AttributeFind(ctx, person.ID(), "age")
	if err != nil {
		log.Fatalf("Failed to find age attribute: %v", err)
	}
	if ageAttr == nil {
		log.Fatal("Expected age attribute to be found")
	}

	fmt.Printf("   Name: %s\n", nameAttr.GetValue())
	fmt.Printf("   Age: %s\n", ageAttr.GetValue())

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
	prodNameAttr, err := store.AttributeFind(ctx, product.ID(), "name")
	if err != nil {
		log.Fatalf("Failed to find name attribute: %v", err)
	}
	if prodNameAttr == nil {
		log.Fatal("Expected name attribute to be found")
	}

	prodPriceAttr, err := store.AttributeFind(ctx, product.ID(), "price")
	if err != nil {
		log.Fatalf("Failed to find price attribute: %v", err)
	}
	if prodPriceAttr == nil {
		log.Fatal("Expected price attribute to be found")
	}

	fmt.Printf("   Name: %s\n", prodNameAttr.GetValue())
	fmt.Printf("   Price: %s\n", prodPriceAttr.GetValue())

	// List all entities
	fmt.Println("\n3. Listing all entities...")
	entities, err := store.EntityList(ctx, entitystore.EntityQueryOptions{})
	if err != nil {
		log.Fatalf("Failed to list entities: %v", err)
	}
	fmt.Printf("   Found %d entities:\n", len(entities))
	for _, e := range entities {
		fmt.Printf("   - [%s] %s (handle: %s)\n", e.GetType(), e.ID(), e.GetHandle())
	}

	// Find entity by ID
	fmt.Println("\n4. Finding entity by ID...")
	found, err := store.EntityFindByID(ctx, person.ID())
	if err != nil {
		log.Fatalf("Failed to find entity: %v", err)
	}
	if found == nil {
		log.Fatal("Expected entity to be found")
	}
	foundNameAttr, err := store.AttributeFind(ctx, found.ID(), "name")
	if err != nil {
		log.Fatalf("Failed to find attribute: %v", err)
	}
	if foundNameAttr == nil {
		log.Fatal("Expected name attribute to be found")
	}
	fmt.Printf("   Found: %s (type: %s)\n", foundNameAttr.GetValue(), found.GetType())

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
		fmt.Printf("   - %s: %s\n", attr.GetKey(), attr.GetValue())
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
