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
	defer db.Close() //nolint:errcheck

	// Create the core store
	store, err := entitystore.NewStore(entitystore.NewStoreOptions{
		DB:                   db,
		EntityTableName:      "entities",
		AttributeTableName:   "attributes",
		AutomigrateEnabled:   true,
		RelationshipsEnabled: true,
		TaxonomiesEnabled:    true,
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
	product := active.EntityCreate("product").
		SetString("name", "Laptop").
		SetFloat("price", 1299.99).
		SetInt("stock", 50)
	if err := product.Save(); err != nil {
		log.Fatalf("Failed to save product: %v", err)
	}
	fmt.Printf("   Created product with ID: %s\n", product.GetEntity().ID())

	// Read an attribute back
	name, exists, err := product.GetString("name")
	if err != nil || !exists {
		log.Fatalf("Failed to read name: %v", err)
	}
	fmt.Printf("   Name: %s\n", name)

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

	// List and count entities by type
	fmt.Println("\n3. Listing products...")
	products, err := active.EntityList(entitystore.EntityQuery().WithEntityType("product"))
	if err != nil {
		log.Fatalf("Failed to list products: %v", err)
	}
	count, _ := active.EntityCount(entitystore.EntityQuery().WithEntityType("product"))
	fmt.Printf("   Found %d product(s):\n", count)
	for _, p := range products {
		n, _, _ := p.GetString("name")
		fmt.Printf("   - %s (ID: %s)\n", n, p.GetEntity().ID())
	}

	// Relationships — entity-centric, no IDs to juggle
	fmt.Println("\n4. Creating relationships...")
	author := active.EntityCreate("author").SetString("name", "Ada")
	if err := author.Save(); err != nil {
		log.Fatalf("Failed to save author: %v", err)
	}
	if err := product.RelateTo(author.GetEntity().ID(), "invented_by"); err != nil {
		log.Fatalf("Failed to relate: %v", err)
	}
	inventors, err := product.Related("invented_by")
	if err != nil {
		log.Fatalf("Failed to list related: %v", err)
	}
	for _, inventor := range inventors {
		n, _, _ := inventor.GetString("name")
		fmt.Printf("   %s invented_by %s\n", name, n)
	}

	// Taxonomies — assign a term, navigate back to tagged entities
	fmt.Println("\n5. Assigning a taxonomy term...")
	cat, err := active.TaxonomyCreate(entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})
	if err != nil {
		log.Fatalf("Failed to create taxonomy: %v", err)
	}
	laptops, err := active.TermCreate(entitystore.TaxonomyTermOptions{
		TaxonomyID: cat.GetTaxonomy().GetID(),
		Name:       "Laptops",
		Slug:       "laptops",
	})
	if err != nil {
		log.Fatalf("Failed to create term: %v", err)
	}
	if err := product.AssignTerm(cat.GetTaxonomy().GetID(), laptops.GetTerm().GetID()); err != nil {
		log.Fatalf("Failed to assign term: %v", err)
	}
	tagged, err := laptops.Entities()
	if err != nil {
		log.Fatalf("Failed to list tagged entities: %v", err)
	}
	fmt.Printf("   %d entit(ies) tagged 'laptops':\n", len(tagged))
	for _, e := range tagged {
		n, _, _ := e.GetString("name")
		fmt.Printf("   - %s\n", n)
	}

	// Soft delete via the entity itself
	fmt.Println("\n6. Trashing the product...")
	trashed, err := product.Trash()
	if err != nil {
		log.Fatalf("Failed to trash product: %v", err)
	}
	fmt.Printf("   Trashed: %v\n", trashed)

	fmt.Println("\n=== Example completed successfully! ===")
}
