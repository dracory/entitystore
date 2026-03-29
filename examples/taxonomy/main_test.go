package main

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dracory/entitystore"
	_ "modernc.org/sqlite"
)

func setupTestStore(t *testing.T) (entitystore.StoreInterface, *sql.DB, func()) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	store, err := entitystore.NewStore(entitystore.NewStoreOptions{
		DB:                     db,
		EntityTableName:        "test_entities",
		AttributeTableName:     "test_attributes",
		TaxonomiesEnabled:      true,
		AutomigrateEnabled:     true,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return store, db, cleanup
}

func TestTaxonomyCreate(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	tax, err := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name:        "Categories",
		Slug:        "categories",
		Description: "Product categories",
		EntityTypes: []string{"product"},
	})
	if err != nil {
		t.Fatalf("Failed to create taxonomy: %v", err)
	}

	if tax.ID() == "" {
		t.Error("Taxonomy ID should not be empty")
	}

	if tax.Name() != "Categories" {
		t.Errorf("Expected name 'Categories', got '%s'", tax.Name())
	}

	if tax.Slug() != "categories" {
		t.Errorf("Expected slug 'categories', got '%s'", tax.Slug())
	}
}

func TestTaxonomyFindBySlug(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})

	found, err := store.TaxonomyFindBySlug(ctx, "categories")
	if err != nil {
		t.Fatalf("Failed to find taxonomy: %v", err)
	}

	if found == nil {
		t.Fatal("Taxonomy should be found")
	}

	if found.Name() != "Categories" {
		t.Errorf("Expected name 'Categories', got '%s'", found.Name())
	}
}

func TestTaxonomyList(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})
	store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Tags",
		Slug: "tags",
	})

	taxonomies, err := store.TaxonomyList(ctx, entitystore.TaxonomyQueryOptions{})
	if err != nil {
		t.Fatalf("Failed to list taxonomies: %v", err)
	}

	if len(taxonomies) != 2 {
		t.Errorf("Expected 2 taxonomies, got %d", len(taxonomies))
	}
}

func TestTaxonomyCount(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{Name: "A", Slug: "a"})
	store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{Name: "B", Slug: "b"})

	count, _ := store.TaxonomyCount(ctx, entitystore.TaxonomyQueryOptions{})
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}
}

func TestTaxonomyTermCreate(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	tax, _ := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})

	term, err := store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: tax.ID(),
		Name:       "Electronics",
		Slug:       "electronics",
	})
	if err != nil {
		t.Fatalf("Failed to create term: %v", err)
	}

	if term.TaxonomyID() != tax.ID() {
		t.Errorf("Expected TaxonomyID %s, got %s", tax.ID(), term.TaxonomyID())
	}

	if term.Name() != "Electronics" {
		t.Errorf("Expected name 'Electronics', got '%s'", term.Name())
	}
}

func TestTaxonomyTermHierarchy(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	tax, _ := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})

	// Create parent term
	parent, _ := store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: tax.ID(),
		Name:       "Electronics",
		Slug:       "electronics",
	})

	// Create child term
	child, err := store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: tax.ID(),
		Name:       "Computers",
		Slug:       "computers",
		ParentID:   parent.ID(),
	})
	if err != nil {
		t.Fatalf("Failed to create child term: %v", err)
	}

	if child.ParentID() != parent.ID() {
		t.Errorf("Expected ParentID %s, got %s", parent.ID(), child.ParentID())
	}
}

func TestTaxonomyTermListByTaxonomy(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	tax, _ := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})

	store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: tax.ID(),
		Name:       "Electronics",
		Slug:       "electronics",
	})
	store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: tax.ID(),
		Name:       "Books",
		Slug:       "books",
	})

	terms, err := store.TaxonomyTermList(ctx, entitystore.TaxonomyTermQueryOptions{
		TaxonomyID: tax.ID(),
	})
	if err != nil {
		t.Fatalf("Failed to list terms: %v", err)
	}

	if len(terms) != 2 {
		t.Errorf("Expected 2 terms, got %d", len(terms))
	}
}

func TestEntityTaxonomyAssign(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	tax, _ := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})

	term, _ := store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: tax.ID(),
		Name:       "Electronics",
		Slug:       "electronics",
	})

	entity, _ := store.EntityCreateWithTypeAndAttributes(ctx, "product", map[string]string{
		"name": "Laptop",
	})

	// Assign entity to term
	err := store.EntityTaxonomyAssign(ctx, entity.ID(), tax.ID(), term.ID())
	if err != nil {
		t.Fatalf("Failed to assign taxonomy: %v", err)
	}

	// Verify assignment exists
	assignments, _ := store.EntityTaxonomyList(ctx, entitystore.EntityTaxonomyQueryOptions{
		EntityID: entity.ID(),
	})

	if len(assignments) != 1 {
		t.Errorf("Expected 1 assignment, got %d", len(assignments))
	}
}

func TestEntityTaxonomyRemove(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	tax, _ := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})

	term, _ := store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: tax.ID(),
		Name:       "Electronics",
		Slug:       "electronics",
	})

	entity, _ := store.EntityCreateWithTypeAndAttributes(ctx, "product", nil)

	// Assign and then remove
	store.EntityTaxonomyAssign(ctx, entity.ID(), tax.ID(), term.ID())
	err := store.EntityTaxonomyRemove(ctx, entity.ID(), tax.ID(), term.ID())
	if err != nil {
		t.Fatalf("Failed to remove assignment: %v", err)
	}

	// Verify removal
	assignments, _ := store.EntityTaxonomyList(ctx, entitystore.EntityTaxonomyQueryOptions{
		EntityID: entity.ID(),
	})

	if len(assignments) != 0 {
		t.Errorf("Expected 0 assignments after removal, got %d", len(assignments))
	}
}

func TestEntityTaxonomyCount(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	tax, _ := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})

	term, _ := store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: tax.ID(),
		Name:       "Electronics",
		Slug:       "electronics",
	})

	entity1, _ := store.EntityCreateWithTypeAndAttributes(ctx, "product", nil)
	entity2, _ := store.EntityCreateWithTypeAndAttributes(ctx, "product", nil)

	store.EntityTaxonomyAssign(ctx, entity1.ID(), tax.ID(), term.ID())
	store.EntityTaxonomyAssign(ctx, entity2.ID(), tax.ID(), term.ID())

	count, _ := store.EntityTaxonomyCount(ctx, entitystore.EntityTaxonomyQueryOptions{
		TaxonomyID: tax.ID(),
		TermID:     term.ID(),
	})

	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}
}

func TestTaxonomyUpdate(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	tax, _ := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})

	// Update
	tax.SetDescription("Updated description")
	err := store.TaxonomyUpdate(ctx, tax)
	if err != nil {
		t.Fatalf("Failed to update taxonomy: %v", err)
	}

	// Verify
	found, _ := store.TaxonomyFind(ctx, tax.ID())
	if found.Description() != "Updated description" {
		t.Errorf("Expected description 'Updated description', got '%s'", found.Description())
	}
}

func TestTaxonomyTermUpdate(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	tax, _ := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})

	term, _ := store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: tax.ID(),
		Name:       "Electronics",
		Slug:       "electronics",
		SortOrder:  1,
	})

	// Update
	term.SetSortOrder(5)
	err := store.TaxonomyTermUpdate(ctx, term)
	if err != nil {
		t.Fatalf("Failed to update term: %v", err)
	}

	// Verify
	found, _ := store.TaxonomyTermFind(ctx, term.ID())
	if found.SortOrder() != 5 {
		t.Errorf("Expected sort order 5, got %d", found.SortOrder())
	}
}

func TestTaxonomyTrash(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	tax, _ := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})

	// Trash
	deleted, err := store.TaxonomyTrash(ctx, tax.ID(), "test_user")
	if err != nil {
		t.Fatalf("Failed to trash taxonomy: %v", err)
	}
	if !deleted {
		t.Error("Expected taxonomy to be deleted")
	}

	// Count should be 0
	count, _ := store.TaxonomyCount(ctx, entitystore.TaxonomyQueryOptions{})
	if count != 0 {
		t.Errorf("Expected count 0 after trash, got %d", count)
	}
}

func TestTaxonomyTermTrash(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	tax, _ := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})

	term, _ := store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: tax.ID(),
		Name:       "Electronics",
		Slug:       "electronics",
	})

	// Trash
	deleted, err := store.TaxonomyTermTrash(ctx, term.ID(), "test_user")
	if err != nil {
		t.Fatalf("Failed to trash term: %v", err)
	}
	if !deleted {
		t.Error("Expected term to be deleted")
	}

	// Count should be 0
	count, _ := store.TaxonomyTermCount(ctx, entitystore.TaxonomyTermQueryOptions{
		TaxonomyID: tax.ID(),
	})
	if count != 0 {
		t.Errorf("Expected count 0 after trash, got %d", count)
	}
}
