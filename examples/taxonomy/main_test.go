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
		DB:                 db,
		EntityTableName:    "test_entities",
		AttributeTableName: "test_attributes",
		TaxonomiesEnabled:  true,
		AutomigrateEnabled: true,
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

	// Test GetID() method consistency
	if tax.GetID() != tax.ID() {
		t.Errorf("expected GetID() '%s' to match ID() '%s'", tax.GetID(), tax.ID())
	}

	if tax.GetName() != "Categories" {
		t.Errorf("Expected name 'Categories', got '%s'", tax.GetName())
	}

	if tax.GetSlug() != "categories" {
		t.Errorf("Expected slug 'categories', got '%s'", tax.GetSlug())
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

	if found.GetName() != "Categories" {
		t.Errorf("Expected name 'Categories', got '%s'", found.GetName())
	}
}

func TestTaxonomyList(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	_, err := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})
	if err != nil {
		t.Fatalf("Failed to create taxonomy: %v", err)
	}

	tax, err := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Tags",
		Slug: "tags",
	})
	if err != nil {
		t.Fatalf("Failed to create taxonomy: %v", err)
	}
	if tax == nil {
		t.Fatal("Expected taxonomy to be created")
	}

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

	tax, err := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{Name: "A", Slug: "a"})
	if err != nil {
		t.Fatalf("Failed to create taxonomy: %v", err)
	}
	if tax == nil {
		t.Fatal("Expected taxonomy to be created")
	}

	_, err = store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{Name: "B", Slug: "b"})
	if err != nil {
		t.Fatalf("Failed to create taxonomy: %v", err)
	}

	count, _ := store.TaxonomyCount(ctx, entitystore.TaxonomyQueryOptions{})
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}
}

func TestTaxonomyTermCreate(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	tax, err := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})
	if err != nil {
		t.Fatalf("Failed to create taxonomy: %v", err)
	}
	if tax == nil {
		t.Fatal("Expected taxonomy to be created")
	}

	term, err := store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: tax.ID(),
		Name:       "Electronics",
		Slug:       "electronics",
	})
	if err != nil {
		t.Fatalf("Failed to create term: %v", err)
	}

	if term.GetTaxonomyID() != tax.ID() {
		t.Errorf("Expected TaxonomyID %s, got %s", tax.ID(), term.GetTaxonomyID())
	}

	// Test GetID() method consistency
	if term.GetID() != term.ID() {
		t.Errorf("expected GetID() '%s' to match ID() '%s'", term.GetID(), term.ID())
	}

	if term.GetName() != "Electronics" {
		t.Errorf("Expected name 'Electronics', got '%s'", term.GetName())
	}
}

func TestTaxonomyTermHierarchy(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	tax, err := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})
	if err != nil {
		t.Fatalf("Failed to create taxonomy: %v", err)
	}
	if tax == nil {
		t.Fatal("Expected taxonomy to be created")
	}

	// Create parent term
	parent, err := store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: tax.ID(),
		Name:       "Electronics",
		Slug:       "electronics",
	})
	if err != nil {
		t.Fatalf("Failed to create parent term: %v", err)
	}
	if parent == nil {
		t.Fatal("Expected parent term to be created")
	}

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

	if child.GetParentID() != parent.ID() {
		t.Errorf("Expected ParentID %s, got %s", parent.ID(), child.GetParentID())
	}
}

func TestTaxonomyTermListByTaxonomy(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	tax, err := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})
	if err != nil {
		t.Fatalf("Failed to create taxonomy: %v", err)
	}
	if tax == nil {
		t.Fatal("Expected taxonomy to be created")
	}

	_, err = store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: tax.ID(),
		Name:       "Electronics",
		Slug:       "electronics",
	})
	if err != nil {
		t.Fatalf("Failed to create term: %v", err)
	}

	_, err = store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: tax.ID(),
		Name:       "Books",
		Slug:       "books",
	})
	if err != nil {
		t.Fatalf("Failed to create term: %v", err)
	}

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

	tax, err := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})
	if err != nil {
		t.Fatalf("Failed to create taxonomy: %v", err)
	}
	if tax == nil {
		t.Fatal("Expected taxonomy to be created")
	}

	term, err := store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: tax.ID(),
		Name:       "Electronics",
		Slug:       "electronics",
	})
	if err != nil {
		t.Fatalf("Failed to create term: %v", err)
	}
	if term == nil {
		t.Fatal("Expected term to be created")
	}

	entity, err := store.EntityCreateWithTypeAndAttributes(ctx, "product", map[string]string{
		"name": "Laptop",
	})
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}
	if entity == nil {
		t.Fatal("Expected entity to be created")
	}

	// Assign entity to term
	err = store.EntityTaxonomyAssign(ctx, entity.ID(), tax.ID(), term.ID())
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

	tax, err := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})
	if err != nil {
		t.Fatalf("Failed to create taxonomy: %v", err)
	}
	if tax == nil {
		t.Fatal("Expected taxonomy to be created")
	}

	term, err := store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: tax.ID(),
		Name:       "Electronics",
		Slug:       "electronics",
	})
	if err != nil {
		t.Fatalf("Failed to create term: %v", err)
	}
	if term == nil {
		t.Fatal("Expected term to be created")
	}

	entity, err := store.EntityCreateWithTypeAndAttributes(ctx, "product", nil)
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}
	if entity == nil {
		t.Fatal("Expected entity to be created")
	}

	// Assign and then remove
	store.EntityTaxonomyAssign(ctx, entity.ID(), tax.ID(), term.ID())
	err = store.EntityTaxonomyRemove(ctx, entity.ID(), tax.ID(), term.ID())
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

	tax, err := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})
	if err != nil {
		t.Fatalf("Failed to create taxonomy: %v", err)
	}
	if tax == nil {
		t.Fatal("Expected taxonomy to be created")
	}

	term, err := store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: tax.ID(),
		Name:       "Electronics",
		Slug:       "electronics",
	})
	if err != nil {
		t.Fatalf("Failed to create term: %v", err)
	}
	if term == nil {
		t.Fatal("Expected term to be created")
	}

	entity1, err := store.EntityCreateWithTypeAndAttributes(ctx, "product", nil)
	if err != nil {
		t.Fatalf("Failed to create entity1: %v", err)
	}
	if entity1 == nil {
		t.Fatal("Expected entity1 to be created")
	}

	entity2, err := store.EntityCreateWithTypeAndAttributes(ctx, "product", nil)
	if err != nil {
		t.Fatalf("Failed to create entity2: %v", err)
	}
	if entity2 == nil {
		t.Fatal("Expected entity2 to be created")
	}

	err = store.EntityTaxonomyAssign(ctx, entity1.ID(), tax.ID(), term.ID())
	if err != nil {
		t.Fatalf("Failed to assign taxonomy to entity1: %v", err)
	}

	err = store.EntityTaxonomyAssign(ctx, entity2.ID(), tax.ID(), term.ID())
	if err != nil {
		t.Fatalf("Failed to assign taxonomy to entity2: %v", err)
	}

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

	tax, err := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})
	if err != nil {
		t.Fatalf("Failed to create taxonomy: %v", err)
	}
	if tax == nil {
		t.Fatal("Expected taxonomy to be created")
	}

	// Update
	tax.SetDescription("Updated description")
	err = store.TaxonomyUpdate(ctx, tax)
	if err != nil {
		t.Fatalf("Failed to update taxonomy: %v", err)
	}

	// Verify
	found, err := store.TaxonomyFind(ctx, tax.ID())
	if err != nil {
		t.Fatalf("Failed to find taxonomy: %v", err)
	}
	if found == nil {
		t.Fatal("Expected to find taxonomy")
	}
	if found.GetDescription() != "Updated description" {
		t.Errorf("Expected description 'Updated description', got '%s'", found.GetDescription())
	}
}

func TestTaxonomyTermUpdate(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	tax, err := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})
	if err != nil {
		t.Fatalf("Failed to create taxonomy: %v", err)
	}
	if tax == nil {
		t.Fatal("Expected taxonomy to be created")
	}

	term, err := store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: tax.ID(),
		Name:       "Electronics",
		Slug:       "electronics",
		SortOrder:  1,
	})
	if err != nil {
		t.Fatalf("Failed to create term: %v", err)
	}
	if term == nil {
		t.Fatal("Expected term to be created")
	}

	// Update
	term.SetSortOrder(5)
	err = store.TaxonomyTermUpdate(ctx, term)
	if err != nil {
		t.Fatalf("Failed to update term: %v", err)
	}

	// Verify
	found, err := store.TaxonomyTermFind(ctx, term.ID())
	if err != nil {
		t.Fatalf("Failed to find term: %v", err)
	}
	if found == nil {
		t.Fatal("Expected to find term")
	}
	if found.GetSortOrder() != 5 {
		t.Errorf("Expected sort order 5, got %d", found.GetSortOrder())
	}
}

func TestTaxonomyTrash(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	defer cleanup()

	ctx := context.Background()

	tax, err := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})
	if err != nil {
		t.Fatalf("Failed to create taxonomy: %v", err)
	}
	if tax == nil {
		t.Fatal("Expected taxonomy to be created")
	}

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

	tax, err := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})
	if err != nil {
		t.Fatalf("Failed to create taxonomy: %v", err)
	}
	if tax == nil {
		t.Fatal("Expected taxonomy to be created")
	}

	term, err := store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: tax.ID(),
		Name:       "Electronics",
		Slug:       "electronics",
	})
	if err != nil {
		t.Fatalf("Failed to create term: %v", err)
	}
	if term == nil {
		t.Fatal("Expected term to be created")
	}

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
