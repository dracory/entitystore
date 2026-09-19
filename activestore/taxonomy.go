package activestore

import (
	"context"
	"errors"

	"github.com/dracory/entitystore"
)

// ActiveTaxonomyInterface wraps a TaxonomyInterface with navigation and
// lifecycle helpers.
type ActiveTaxonomyInterface interface {
	GetTaxonomy() entitystore.TaxonomyInterface
	Terms() ([]ActiveTaxonomyTermInterface, error) // terms in this taxonomy
	Entities() ([]ActiveEntityInterface, error)    // assigned entities
	Trash(deletedBy string) (bool, error)
	Restore() (bool, error)
	Delete() (bool, error)
}

// ActiveTaxonomyTermInterface wraps a TaxonomyTermInterface with
// navigation and lifecycle helpers.
type ActiveTaxonomyTermInterface interface {
	GetTerm() entitystore.TaxonomyTermInterface
	GetTaxonomy() (ActiveTaxonomyInterface, error)
	Parent() (ActiveTaxonomyTermInterface, error) // hierarchical nav
	Children() ([]ActiveTaxonomyTermInterface, error)
	Entities() ([]ActiveEntityInterface, error) // entities tagged with this term
	Trash(deletedBy string) (bool, error)
	Restore() (bool, error)
	Delete() (bool, error)
}

// ---------------------------------------------------------------------------
// Taxonomy wrapper
// ---------------------------------------------------------------------------

type activeTaxonomyImplementation struct {
	ctx      context.Context
	store    entitystore.StoreInterface
	taxonomy entitystore.TaxonomyInterface
}

func newActiveTaxonomy(ctx context.Context, store entitystore.StoreInterface, taxonomy entitystore.TaxonomyInterface) (ActiveTaxonomyInterface, error) {
	if taxonomy == nil {
		return nil, errors.New("taxonomy cannot be nil")
	}
	return &activeTaxonomyImplementation{ctx: ctx, store: store, taxonomy: taxonomy}, nil
}

func (t *activeTaxonomyImplementation) GetTaxonomy() entitystore.TaxonomyInterface {
	return t.taxonomy
}

func (t *activeTaxonomyImplementation) Terms() ([]ActiveTaxonomyTermInterface, error) {
	terms, err := t.store.TaxonomyTermList(t.ctx, entitystore.TaxonomyTermQuery().WithTaxonomyID(t.taxonomy.GetID()))
	if err != nil {
		return nil, err
	}
	active := make([]ActiveTaxonomyTermInterface, 0, len(terms))
	for _, term := range terms {
		wrapped, err := newActiveTerm(t.ctx, t.store, term)
		if err != nil {
			return nil, err
		}
		active = append(active, wrapped)
	}
	return active, nil
}

func (t *activeTaxonomyImplementation) Entities() ([]ActiveEntityInterface, error) {
	assignments, err := t.store.EntityTaxonomyList(t.ctx, entitystore.EntityTaxonomyQuery().WithTaxonomyID(t.taxonomy.GetID()))
	if err != nil {
		return nil, err
	}
	return hydrateEntities(t.ctx, t.store, assignments)
}

func (t *activeTaxonomyImplementation) Trash(deletedBy string) (bool, error) {
	return t.store.TaxonomyTrash(t.ctx, t.taxonomy.GetID(), deletedBy)
}

func (t *activeTaxonomyImplementation) Restore() (bool, error) {
	return t.store.TaxonomyRestore(t.ctx, t.taxonomy.GetID())
}

func (t *activeTaxonomyImplementation) Delete() (bool, error) {
	return t.store.TaxonomyDelete(t.ctx, t.taxonomy.GetID())
}

// ---------------------------------------------------------------------------
// Term wrapper
// ---------------------------------------------------------------------------

type activeTermImplementation struct {
	ctx   context.Context
	store entitystore.StoreInterface
	term  entitystore.TaxonomyTermInterface
}

func newActiveTerm(ctx context.Context, store entitystore.StoreInterface, term entitystore.TaxonomyTermInterface) (ActiveTaxonomyTermInterface, error) {
	if term == nil {
		return nil, errors.New("term cannot be nil")
	}
	return &activeTermImplementation{ctx: ctx, store: store, term: term}, nil
}

func (t *activeTermImplementation) GetTerm() entitystore.TaxonomyTermInterface {
	return t.term
}

func (t *activeTermImplementation) GetTaxonomy() (ActiveTaxonomyInterface, error) {
	taxonomy, err := t.store.TaxonomyFind(t.ctx, t.term.GetTaxonomyID())
	if err != nil || taxonomy == nil {
		return nil, err
	}
	return newActiveTaxonomy(t.ctx, t.store, taxonomy)
}

func (t *activeTermImplementation) Parent() (ActiveTaxonomyTermInterface, error) {
	parentID := t.term.GetParentID()
	if parentID == "" {
		return nil, nil
	}
	parent, err := t.store.TaxonomyTermFind(t.ctx, parentID)
	if err != nil || parent == nil {
		return nil, err
	}
	return newActiveTerm(t.ctx, t.store, parent)
}

func (t *activeTermImplementation) Children() ([]ActiveTaxonomyTermInterface, error) {
	children, err := t.store.TaxonomyTermList(t.ctx, entitystore.TaxonomyTermQuery().
		WithTaxonomyID(t.term.GetTaxonomyID()).
		WithParentID(t.term.GetID()))
	if err != nil {
		return nil, err
	}
	active := make([]ActiveTaxonomyTermInterface, 0, len(children))
	for _, child := range children {
		wrapped, err := newActiveTerm(t.ctx, t.store, child)
		if err != nil {
			return nil, err
		}
		active = append(active, wrapped)
	}
	return active, nil
}

func (t *activeTermImplementation) Entities() ([]ActiveEntityInterface, error) {
	assignments, err := t.store.EntityTaxonomyList(t.ctx, entitystore.EntityTaxonomyQuery().WithTermID(t.term.GetID()))
	if err != nil {
		return nil, err
	}
	return hydrateEntities(t.ctx, t.store, assignments)
}

func (t *activeTermImplementation) Trash(deletedBy string) (bool, error) {
	return t.store.TaxonomyTermTrash(t.ctx, t.term.GetID(), deletedBy)
}

func (t *activeTermImplementation) Restore() (bool, error) {
	return t.store.TaxonomyTermRestore(t.ctx, t.term.GetID())
}

func (t *activeTermImplementation) Delete() (bool, error) {
	return t.store.TaxonomyTermDelete(t.ctx, t.term.GetID())
}

// hydrateEntities resolves entity-taxonomy assignments into wrapped entities.
func hydrateEntities(ctx context.Context, store entitystore.StoreInterface, assignments []entitystore.EntityTaxonomyInterface) ([]ActiveEntityInterface, error) {
	active := make([]ActiveEntityInterface, 0, len(assignments))
	for _, assignment := range assignments {
		ent, err := store.EntityFindByID(ctx, assignment.GetEntityID())
		if err != nil {
			return nil, err
		}
		if ent == nil {
			continue // dangling assignment: entity deleted or missing
		}
		wrapped, err := newActiveEntity(ctx, store, ent)
		if err != nil {
			return nil, err
		}
		active = append(active, wrapped)
	}
	return active, nil
}

// ---------------------------------------------------------------------------
// Store-level taxonomy methods
// ---------------------------------------------------------------------------

func (s *activeStoreImplementation) TaxonomyCreate(options entitystore.TaxonomyOptions) (ActiveTaxonomyInterface, error) {
	taxonomy, err := s.store.TaxonomyCreateByOptions(s.ctx, options)
	if err != nil {
		return nil, err
	}
	return newActiveTaxonomy(s.ctx, s.store, taxonomy)
}

func (s *activeStoreImplementation) TaxonomyFindByID(taxonomyID string) (ActiveTaxonomyInterface, error) {
	taxonomy, err := s.store.TaxonomyFind(s.ctx, taxonomyID)
	if err != nil || taxonomy == nil {
		return nil, err
	}
	return newActiveTaxonomy(s.ctx, s.store, taxonomy)
}

func (s *activeStoreImplementation) TaxonomyFindBySlug(slug string) (ActiveTaxonomyInterface, error) {
	taxonomy, err := s.store.TaxonomyFindBySlug(s.ctx, slug)
	if err != nil || taxonomy == nil {
		return nil, err
	}
	return newActiveTaxonomy(s.ctx, s.store, taxonomy)
}

func (s *activeStoreImplementation) TaxonomyList(query entitystore.TaxonomyQueryInterface) ([]ActiveTaxonomyInterface, error) {
	taxonomies, err := s.store.TaxonomyList(s.ctx, query)
	if err != nil {
		return nil, err
	}
	active := make([]ActiveTaxonomyInterface, 0, len(taxonomies))
	for _, taxonomy := range taxonomies {
		wrapped, err := newActiveTaxonomy(s.ctx, s.store, taxonomy)
		if err != nil {
			return nil, err
		}
		active = append(active, wrapped)
	}
	return active, nil
}

func (s *activeStoreImplementation) TaxonomyCount(query entitystore.TaxonomyQueryInterface) (int64, error) {
	return s.store.TaxonomyCount(s.ctx, query)
}

func (s *activeStoreImplementation) TaxonomyUpdate(taxonomy entitystore.TaxonomyInterface) error {
	return s.store.TaxonomyUpdate(s.ctx, taxonomy)
}

func (s *activeStoreImplementation) TaxonomyTrash(taxonomyID, deletedBy string) (bool, error) {
	return s.store.TaxonomyTrash(s.ctx, taxonomyID, deletedBy)
}

func (s *activeStoreImplementation) TaxonomyRestore(taxonomyID string) (bool, error) {
	return s.store.TaxonomyRestore(s.ctx, taxonomyID)
}

func (s *activeStoreImplementation) TaxonomyDelete(taxonomyID string) (bool, error) {
	return s.store.TaxonomyDelete(s.ctx, taxonomyID)
}

// ---------------------------------------------------------------------------
// Store-level term methods
// ---------------------------------------------------------------------------

func (s *activeStoreImplementation) TermCreate(options entitystore.TaxonomyTermOptions) (ActiveTaxonomyTermInterface, error) {
	term, err := s.store.TaxonomyTermCreateByOptions(s.ctx, options)
	if err != nil {
		return nil, err
	}
	return newActiveTerm(s.ctx, s.store, term)
}

func (s *activeStoreImplementation) TermFindByID(termID string) (ActiveTaxonomyTermInterface, error) {
	term, err := s.store.TaxonomyTermFind(s.ctx, termID)
	if err != nil || term == nil {
		return nil, err
	}
	return newActiveTerm(s.ctx, s.store, term)
}

func (s *activeStoreImplementation) TermFindBySlug(taxonomyID, slug string) (ActiveTaxonomyTermInterface, error) {
	term, err := s.store.TaxonomyTermFindBySlug(s.ctx, taxonomyID, slug)
	if err != nil || term == nil {
		return nil, err
	}
	return newActiveTerm(s.ctx, s.store, term)
}

func (s *activeStoreImplementation) TermList(query entitystore.TaxonomyTermQueryInterface) ([]ActiveTaxonomyTermInterface, error) {
	terms, err := s.store.TaxonomyTermList(s.ctx, query)
	if err != nil {
		return nil, err
	}
	active := make([]ActiveTaxonomyTermInterface, 0, len(terms))
	for _, term := range terms {
		wrapped, err := newActiveTerm(s.ctx, s.store, term)
		if err != nil {
			return nil, err
		}
		active = append(active, wrapped)
	}
	return active, nil
}

func (s *activeStoreImplementation) TermCount(query entitystore.TaxonomyTermQueryInterface) (int64, error) {
	return s.store.TaxonomyTermCount(s.ctx, query)
}

func (s *activeStoreImplementation) TermUpdate(term entitystore.TaxonomyTermInterface) error {
	return s.store.TaxonomyTermUpdate(s.ctx, term)
}

func (s *activeStoreImplementation) TermTrash(termID, deletedBy string) (bool, error) {
	return s.store.TaxonomyTermTrash(s.ctx, termID, deletedBy)
}

func (s *activeStoreImplementation) TermRestore(termID string) (bool, error) {
	return s.store.TaxonomyTermRestore(s.ctx, termID)
}

func (s *activeStoreImplementation) TermDelete(termID string) (bool, error) {
	return s.store.TaxonomyTermDelete(s.ctx, termID)
}
