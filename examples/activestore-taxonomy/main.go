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
	defer db.Close() //nolint:errcheck

	store, err := entitystore.NewStore(entitystore.NewStoreOptions{
		DB:                 db,
		EntityTableName:    "entities",
		AttributeTableName: "attributes",
		AutomigrateEnabled: true,
		TaxonomiesEnabled:  true, // required for taxonomies
	})
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	active, err := activestore.New(ctx, store)
	if err != nil {
		log.Fatalf("Failed to create active store: %v", err)
	}

	fmt.Println("=== ActiveStore Taxonomy Example ===")

	// Create a taxonomy and a term — both return active wrappers
	fmt.Println("1. Creating taxonomy and term...")
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
	fmt.Println("   Created categories/laptops")

	// Assign a product to the term — the entity's ID is implicit
	fmt.Println("\n2. Assigning the term to a product...")
	product := active.EntityCreate("product").SetString("name", "Laptop")
	if _, err := product.Save(); err != nil {
		log.Fatalf("Failed to save product: %v", err)
	}
	if err := product.AssignTerm(cat.GetTaxonomy().GetID(), laptops.GetTerm().GetID()); err != nil {
		log.Fatalf("Failed to assign term: %v", err)
	}
	fmt.Println("   Assigned")

	// Navigate both directions without writing queries:
	// entity -> its terms, term -> its entities
	fmt.Println("\n3. Navigating...")
	terms, err := product.Terms(cat.GetTaxonomy().GetID())
	if err != nil {
		log.Fatalf("Failed to list terms: %v", err)
	}
	for _, t := range terms {
		fmt.Printf("   product has term: %s\n", t.GetTerm().GetName())
	}

	tagged, err := laptops.Entities()
	if err != nil {
		log.Fatalf("Failed to list entities: %v", err)
	}
	for _, e := range tagged {
		name, _, _ := e.GetString("name")
		fmt.Printf("   'laptops' tags: %s\n", name)
	}

	fmt.Println("\n=== Example completed successfully! ===")
}
