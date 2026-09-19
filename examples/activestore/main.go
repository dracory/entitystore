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

	// Open SQLite database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer func() { _ = db.Close() }()

	// Create the core store
	store, err := entitystore.NewStore(entitystore.NewStoreOptions{
		DB:                 db,
		EntityTableName:    "entities",
		AttributeTableName: "attributes",
		AutomigrateEnabled: true,
	})
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	// Wrap it in the active store — the single entry point
	active, err := activestore.New(ctx, store)
	if err != nil {
		log.Fatalf("Failed to create active store: %v", err)
	}

	fmt.Println("=== ActiveStore Example ===")

	// Create a product entity with the fluent API.
	// Attributes are staged until Save() creates the row.
	fmt.Println("1. Creating a product entity...")
	product, err := active.EntityCreate("product").
		SetString("name", "Laptop").
		SetFloat("price", 1299.99).
		SetInt("stock", 50).
		Save()
	if err != nil {
		log.Fatalf("Failed to save product: %v", err)
	}
	fmt.Printf("   Created product with ID: %s\n", product.GetEntity().ID())

	// Find an existing entity — writes on persisted entities
	// go to the store immediately.
	fmt.Println("\n2. Finding and updating the product...")
	found, err := active.EntityFindByID(product.GetEntity().ID())
	if err != nil {
		log.Fatalf("Failed to find product: %v", err)
	}
	found.SetInt("stock", 45)
	if err := found.Err(); err != nil {
		log.Fatalf("Failed to update stock: %v", err)
	}
	fmt.Println("   Updated stock to 45")

	// List entities by type — results are wrapped automatically
	fmt.Println("\n3. Listing products...")
	products, err := active.EntityList(entitystore.EntityQuery().WithEntityType("product"))
	if err != nil {
		log.Fatalf("Failed to list products: %v", err)
	}
	for _, p := range products {
		name, _, _ := p.GetString("name")
		fmt.Printf("   - %s\n", name)
	}

	// Soft delete via the entity itself
	fmt.Println("\n4. Trashing the product...")
	trashed, err := product.Trash()
	if err != nil {
		log.Fatalf("Failed to trash product: %v", err)
	}
	fmt.Printf("   Trashed: %v\n", trashed)

	fmt.Println("\n=== Example completed successfully! ===")
}
