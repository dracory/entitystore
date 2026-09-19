package activestore

import (
	"context"
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

	cat, _ := active.TaxonomyCreate(entitystore.TaxonomyOptions{Name: "Categories", Slug: "categories"})
	laptops, _ := active.TermCreate(entitystore.TaxonomyTermOptions{
		TaxonomyID: cat.GetTaxonomy().GetID(),
		Name:       "Laptops",
		Slug:       "laptops",
	})

	product := active.EntityCreate("product")
	if _, err := product.Save(); err != nil {
		t.Fatal(err)
	}
	if err := product.AssignTerm(cat.GetTaxonomy().GetID(), laptops.GetTerm().GetID()); err != nil {
		t.Fatal(err)
	}

	terms, err := product.Terms(cat.GetTaxonomy().GetID())
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

	if err := product.RemoveTerm(cat.GetTaxonomy().GetID(), laptops.GetTerm().GetID()); err != nil {
		t.Fatal(err)
	}
	terms, _ = product.Terms(cat.GetTaxonomy().GetID())
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

	cat, _ := active.TaxonomyCreate(entitystore.TaxonomyOptions{Name: "Categories", Slug: "categories"})
	parent, _ := active.TermCreate(entitystore.TaxonomyTermOptions{
		TaxonomyID: cat.GetTaxonomy().GetID(), Name: "Computers", Slug: "computers",
	})
	child, _ := active.TermCreate(entitystore.TaxonomyTermOptions{
		TaxonomyID: cat.GetTaxonomy().GetID(), Name: "Laptops", Slug: "laptops",
		ParentID: parent.GetTerm().GetID(),
	})

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
	noParent, err := orphan.Parent()
	if err != nil || noParent != nil {
		t.Fatalf("expected nil parent, got %v err=%v", noParent, err)
	}
}
