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

	// Create the store with taxonomies enabled
	store, err := entitystore.NewStore(entitystore.NewStoreOptions{
		DB:                 db,
		EntityTableName:    "entities",
		AttributeTableName: "attributes",
		TaxonomiesEnabled:  true,
		AutomigrateEnabled: true,
	})
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	fmt.Println("=== EntityStore Taxonomy Example ===")

	// Create "Product Categories" taxonomy
	fmt.Println("1. Creating 'Product Categories' taxonomy...")
	categoriesTax, err := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name:        "Product Categories",
		Slug:        "product_categories",
		Description: "Categories for organizing products",
		EntityTypes: []string{"product"},
	})
	if err != nil {
		log.Fatalf("Failed to create taxonomy: %v", err)
	}
	fmt.Printf("   Created taxonomy: %s (ID: %s)\n", categoriesTax.GetName(), categoriesTax.ID())
	fmt.Printf("   Slug: %s, Entity Types: %v\n", categoriesTax.GetSlug(), categoriesTax.GetEntityTypes())

	// Create taxonomy terms (hierarchical categories)
	fmt.Println("\n2. Creating taxonomy terms (hierarchical categories)...")

	// Top level: Electronics
	electronics, _ := store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: categoriesTax.ID(),
		Name:       "Electronics",
		Slug:       "electronics",
		SortOrder:  1,
	})
	fmt.Printf("   - Electronics (ID: %s)\n", electronics.ID())

	// Sub-category: Computers (parent = Electronics)
	computers, _ := store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: categoriesTax.ID(),
		Name:       "Computers",
		Slug:       "computers",
		ParentID:   electronics.ID(),
		SortOrder:  1,
	})
	fmt.Printf("   - Computers (ID: %s, Parent: Electronics)\n", computers.ID())

	// Sub-category: Laptops (parent = Computers)
	laptops, _ := store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: categoriesTax.ID(),
		Name:       "Laptops",
		Slug:       "laptops",
		ParentID:   computers.ID(),
		SortOrder:  1,
	})
	fmt.Printf("   - Laptops (ID: %s, Parent: Computers)\n", laptops.ID())

	// Another top level: Books
	books, _ := store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: categoriesTax.ID(),
		Name:       "Books",
		Slug:       "books",
		SortOrder:  2,
	})
	fmt.Printf("   - Books (ID: %s)\n", books.ID())

	// Create product entities
	fmt.Println("\n3. Creating product entities...")
	macbook, _ := store.EntityCreateWithTypeAndAttributes(ctx, "product", map[string]string{
		"name":  "MacBook Pro",
		"price": "1999.99",
		"sku":   "MBP-16-001",
	})
	hpLaptop, _ := store.EntityCreateWithTypeAndAttributes(ctx, "product", map[string]string{
		"name":  "HP Pavilion",
		"price": "699.99",
		"sku":   "HP-PAV-001",
	})
	theHobbit, _ := store.EntityCreateWithTypeAndAttributes(ctx, "product", map[string]string{
		"name":  "The Hobbit",
		"price": "14.99",
		"sku":   "BK-HOBBIT-001",
	})
	fmt.Printf("   - %s (ID: %s)\n", macbook.GetTemp("name"), macbook.ID())
	fmt.Printf("   - %s (ID: %s)\n", hpLaptop.GetTemp("name"), hpLaptop.ID())
	fmt.Printf("   - %s (ID: %s)\n", theHobbit.GetTemp("name"), theHobbit.ID())

	// Assign products to taxonomy terms
	fmt.Println("\n4. Assigning products to taxonomy terms...")
	err = store.EntityTaxonomyAssign(ctx, macbook.ID(), categoriesTax.ID(), laptops.ID())
	if err != nil {
		log.Fatalf("Failed to assign taxonomy: %v", err)
	}
	fmt.Printf("   ✓ %s assigned to Laptops\n", macbook.GetTemp("name"))

	store.EntityTaxonomyAssign(ctx, hpLaptop.ID(), categoriesTax.ID(), laptops.ID())
	fmt.Printf("   ✓ %s assigned to Laptops\n", hpLaptop.GetTemp("name"))

	store.EntityTaxonomyAssign(ctx, theHobbit.ID(), categoriesTax.ID(), books.ID())
	fmt.Printf("   ✓ %s assigned to Books\n", theHobbit.GetTemp("name"))

	// Query taxonomy assignments
	fmt.Println("\n5. Finding products in 'Laptops' category...")
	assignments, err := store.EntityTaxonomyList(ctx, entitystore.EntityTaxonomyQueryOptions{
		TaxonomyID: categoriesTax.ID(),
		TermID:     laptops.ID(),
	})
	if err != nil {
		log.Fatalf("Failed to list assignments: %v", err)
	}
	fmt.Printf("   Found %d products in Laptops:\n", len(assignments))
	for _, assignment := range assignments {
		product, _ := store.EntityFindByID(ctx, assignment.GetEntityID())
		fmt.Printf("   - %s ($%s)\n", product.GetTemp("name"), product.GetTemp("price"))
	}

	// List all terms in taxonomy
	fmt.Println("\n6. Listing all terms in 'Product Categories' taxonomy...")
	terms, _ := store.TaxonomyTermList(ctx, entitystore.TaxonomyTermQueryOptions{
		TaxonomyID: categoriesTax.ID(),
	})
	fmt.Printf("   Found %d terms:\n", len(terms))
	for _, term := range terms {
		parentInfo := ""
		if term.GetParentID() != "" {
			parent, _ := store.TaxonomyTermFind(ctx, term.GetParentID())
			if parent != nil {
				parentInfo = fmt.Sprintf(" (parent: %s)", parent.GetName())
			}
		}
		fmt.Printf("   - %s%s\n", term.GetName(), parentInfo)
	}

	// Find taxonomy by slug
	fmt.Println("\n7. Finding taxonomy by slug...")
	foundTax, _ := store.TaxonomyFindBySlug(ctx, "product_categories")
	if foundTax != nil {
		fmt.Printf("   Found: %s\n", foundTax.GetName())
	}

	// Find term by slug
	fmt.Println("\n8. Finding term by slug within taxonomy...")
	foundTerm, _ := store.TaxonomyTermFindBySlug(ctx, categoriesTax.ID(), "laptops")
	if foundTerm != nil {
		fmt.Printf("   Found: %s\n", foundTerm.GetName())
	}

	// Count terms in taxonomy
	fmt.Println("\n9. Counting taxonomy terms...")
	termCount, _ := store.TaxonomyTermCount(ctx, entitystore.TaxonomyTermQueryOptions{
		TaxonomyID: categoriesTax.ID(),
	})
	fmt.Printf("   Total terms in Product Categories: %d\n", termCount)

	// Count entity assignments
	fmt.Println("\n10. Counting entity-taxonomy assignments...")
	assignmentCount, _ := store.EntityTaxonomyCount(ctx, entitystore.EntityTaxonomyQueryOptions{
		TaxonomyID: categoriesTax.ID(),
	})
	fmt.Printf("    Total product assignments: %d\n", assignmentCount)

	// Remove assignment
	fmt.Println("\n11. Removing a product assignment...")
	err = store.EntityTaxonomyRemove(ctx, theHobbit.ID(), categoriesTax.ID(), books.ID())
	if err != nil {
		log.Fatalf("Failed to remove assignment: %v", err)
	}
	fmt.Println("    ✓ Removed The Hobbit from Books category")

	// Verify removal
	assignmentCount, _ = store.EntityTaxonomyCount(ctx, entitystore.EntityTaxonomyQueryOptions{})
	fmt.Printf("    Total assignments after removal: %d\n", assignmentCount)

	// Update taxonomy
	fmt.Println("\n12. Updating taxonomy...")
	categoriesTax.SetDescription("Organized product categories for e-commerce")
	store.TaxonomyUpdate(ctx, categoriesTax)
	fmt.Println("    ✓ Updated taxonomy description")

	fmt.Println("\n=== Example completed successfully! ===")
}
