package activestore

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dracory/entitystore"
)

func TestTaxonomyAndTerms(t *testing.T) {
	ctx := context.Background()
	active, err := New(ctx, initFullStore(t, "activestore_tax.db"))
	if err != nil {
		t.Fatal(err)
	}

	cat, err := active.TaxonomyCreate(entitystore.TaxonomyOptions{
		Name: "Categories",
		Slug: "categories",
	})
	if err != nil {
		t.Fatal(err)
	}

	laptops, err := active.TermCreate(entitystore.TaxonomyTermOptions{
		TaxonomyID: cat.GetTaxonomy().GetID(),
		Name:       "Laptops",
		Slug:       "laptops",
	})
	if err != nil {
		t.Fatal(err)
	}

	// taxonomy navigation
	terms, err := cat.Terms()
	if err != nil || len(terms) != 1 {
		t.Fatalf("expected 1 term, got %v err=%v", len(terms), err)
	}
	if terms[0].GetTerm().GetSlug() != "laptops" {
		t.Fatal("expected slug laptops")
	}

	// find by slug round-trips through the wrapper
	found, err := active.TermFindBySlug(cat.GetTaxonomy().GetID(), "laptops")
	if err != nil || found.GetTerm().GetID() != laptops.GetTerm().GetID() {
		t.Fatalf("TermFindBySlug mismatch: %v", err)
	}
	tax, err := found.GetTaxonomy()
	if err != nil || tax.GetTaxonomy().GetID() != cat.GetTaxonomy().GetID() {
		t.Fatalf("GetTaxonomy mismatch: %v", err)
	}
}

func TestAssignTermAndEntities(t *testing.T) {
	ctx := context.Background()
	active, err := New(ctx, initFullStore(t, "activestore_assign.db"))
	if err != nil {
		t.Fatal(err)
	}

	cat, err := active.TaxonomyCreate(entitystore.TaxonomyOptions{Name: "Categories", Slug: "categories"})
	if err != nil {
		t.Fatal(err)
	}
	laptops, err := active.TermCreate(entitystore.TaxonomyTermOptions{
		TaxonomyID: cat.GetTaxonomy().GetID(),
		Name:       "Laptops",
		Slug:       "laptops",
	})
	if err != nil {
		t.Fatal(err)
	}

	product := active.EntityCreate("product")
	if _, err := product.Save(); err != nil {
		t.Fatal(err)
	}
	if err := product.AssignTerm(cat, laptops); err != nil {
		t.Fatal(err)
	}

	terms, err := product.Terms(cat)
	if err != nil || len(terms) != 1 {
		t.Fatalf("expected 1 term, got %v err=%v", len(terms), err)
	}

	entities, err := laptops.Entities()
	if err != nil || len(entities) != 1 {
		t.Fatalf("expected 1 entity, got %v err=%v", len(entities), err)
	}
	if entities[0].GetEntity().ID() != product.GetEntity().ID() {
		t.Fatal("expected the product entity")
	}

	if err := product.RemoveTerm(cat, laptops); err != nil {
		t.Fatal(err)
	}
	terms, _ = product.Terms(cat)
	if len(terms) != 0 {
		t.Fatalf("expected 0 terms after remove, got %v", len(terms))
	}
}

func TestTermHierarchy(t *testing.T) {
	ctx := context.Background()
	active, err := New(ctx, initFullStore(t, "activestore_hier.db"))
	if err != nil {
		t.Fatal(err)
	}

	cat, err := active.TaxonomyCreate(entitystore.TaxonomyOptions{Name: "Categories", Slug: "categories"})
	if err != nil {
		t.Fatal(err)
	}
	parent, err := active.TermCreate(entitystore.TaxonomyTermOptions{
		TaxonomyID: cat.GetTaxonomy().GetID(), Name: "Computers", Slug: "computers",
	})
	if err != nil {
		t.Fatal(err)
	}
	child, err := active.TermCreate(entitystore.TaxonomyTermOptions{
		TaxonomyID: cat.GetTaxonomy().GetID(), Name: "Laptops", Slug: "laptops",
		ParentID: parent.GetTerm().GetID(),
	})
	if err != nil {
		t.Fatal(err)
	}

	gotParent, err := child.Parent()
	if err != nil || gotParent == nil {
		t.Fatalf("expected parent, got %v err=%v", gotParent, err)
	}
	if gotParent.GetTerm().GetSlug() != "computers" {
		t.Fatal("expected parent slug computers")
	}

	children, err := parent.Children()
	if err != nil || len(children) != 1 {
		t.Fatalf("expected 1 child, got %v err=%v", len(children), err)
	}

	orphan, err := active.TermCreate(entitystore.TaxonomyTermOptions{
		TaxonomyID: cat.GetTaxonomy().GetID(), Name: "Phones", Slug: "phones",
	})
	if err != nil {
		t.Fatal(err)
	}
	noParent, err := orphan.Parent()
	if err != nil || noParent != nil {
		t.Fatalf("expected nil parent, got %v err=%v", noParent, err)
	}
}

func TestTaxonomyAndTermFind_NotFound(t *testing.T) {
	active, err := New(context.Background(), initFullStore(t, "activestore_taxnotfound.db"))
	if err != nil {
		t.Fatal(err)
	}
	tax, err := active.TaxonomyFindByID("no-such-id")
	if err != nil || tax != nil {
		t.Fatalf("expected (nil, nil), got tax=%v err=%v", tax, err)
	}
	tax, err = active.TaxonomyFindBySlug("no-such-slug")
	if err != nil || tax != nil {
		t.Fatalf("expected (nil, nil), got tax=%v err=%v", tax, err)
	}
	term, err := active.TermFindByID("no-such-id")
	if err != nil || term != nil {
		t.Fatalf("expected (nil, nil), got term=%v err=%v", term, err)
	}
	term, err = active.TermFindBySlug("no-such-taxonomy", "no-such-slug")
	if err != nil || term != nil {
		t.Fatalf("expected (nil, nil), got term=%v err=%v", term, err)
	}
}

func TestEntities_SkipsDanglingEntity(t *testing.T) {
	ctx := context.Background()
	store := initFullStore(t, "activestore_taxdangle.db")
	active, err := New(ctx, store)
	if err != nil {
		t.Fatal(err)
	}

	cat, err := active.TaxonomyCreate(entitystore.TaxonomyOptions{Name: "Categories", Slug: "categories"})
	if err != nil {
		t.Fatal(err)
	}
	laptops, err := active.TermCreate(entitystore.TaxonomyTermOptions{
		TaxonomyID: cat.GetTaxonomy().GetID(), Name: "Laptops", Slug: "laptops",
	})
	if err != nil {
		t.Fatal(err)
	}

	product := active.EntityCreate("product")
	if _, err := product.Save(); err != nil {
		t.Fatal(err)
	}
	if err := product.AssignTerm(cat, laptops); err != nil {
		t.Fatal(err)
	}

	// hard-delete the entity, leaving a dangling assignment
	if _, err := store.EntityDelete(ctx, product.GetEntity().ID()); err != nil {
		t.Fatal(err)
	}

	entities, err := laptops.Entities()
	if err != nil {
		t.Fatalf("expected nil error for dangling assignment, got %v", err)
	}
	if len(entities) != 0 {
		t.Fatalf("expected dangling assignment to be skipped, got %v entities", len(entities))
	}
}

func TestTerms_SkipsDanglingTerm(t *testing.T) {
	ctx := context.Background()
	initFullStore(t, "activestore_termdangle.db")
	active, err := New(ctx, initFullStore(t, "activestore_termdangle.db"))
	if err != nil {
		t.Fatal(err)
	}

	cat, err := active.TaxonomyCreate(entitystore.TaxonomyOptions{Name: "Categories", Slug: "categories"})
	if err != nil {
		t.Fatal(err)
	}
	laptops, err := active.TermCreate(entitystore.TaxonomyTermOptions{
		TaxonomyID: cat.GetTaxonomy().GetID(), Name: "Laptops", Slug: "laptops",
	})
	if err != nil {
		t.Fatal(err)
	}

	product := active.EntityCreate("product")
	if _, err := product.Save(); err != nil {
		t.Fatal(err)
	}
	if err := product.AssignTerm(cat, laptops); err != nil {
		t.Fatal(err)
	}

	// delete the term row directly (TaxonomyTermDelete refuses while
	// assignments exist), leaving a dangling assignment
	db, err := sql.Open("sqlite", "file:activestore_termdangle.db?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("DELETE FROM entities_taxonomy_terms WHERE id = ?", laptops.GetTerm().GetID()); err != nil {
		t.Fatal(err)
	}

	terms, err := product.Terms(cat)
	if err != nil {
		t.Fatalf("expected nil error for dangling assignment, got %v", err)
	}
	if len(terms) != 0 {
		t.Fatalf("expected dangling term to be skipped, got %v terms", len(terms))
	}
}
