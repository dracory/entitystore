// Package activestore provides an optional Active Record style wrapper
// around entitystore. It binds a StoreInterface to rich entity objects so
// that attribute writes, persistence, and trash/delete operations can be
// performed with a fluent, object-oriented API.
package activestore

import (
	"context"
	"errors"

	"github.com/dracory/entitystore"
)

// ActiveStoreInterface is a store-aware facade that returns
// ActiveEntityInterface results.
type ActiveStoreInterface interface {
	EntityCreate(entityType string) ActiveEntityInterface
	EntityFindByID(entityID string) (ActiveEntityInterface, error)
	EntityList(query entitystore.EntityQueryInterface) ([]ActiveEntityInterface, error)
	EntityCount(query entitystore.EntityQueryInterface) (int64, error)
	EntityTrash(entityID string) (bool, error)
	EntityDelete(entityID string) (bool, error)
	EntityWrap(entity entitystore.EntityInterface) (ActiveEntityInterface, error)
	GetStore() entitystore.StoreInterface

	// Relationships (requires RelationshipsEnabled)
	RelationshipCreate(options entitystore.RelationshipOptions) (ActiveRelationshipInterface, error)
	RelationshipFindByID(relationshipID string) (ActiveRelationshipInterface, error)
	RelationshipList(query entitystore.RelationshipQueryInterface) ([]ActiveRelationshipInterface, error)
	RelationshipCount(query entitystore.RelationshipQueryInterface) (int64, error)
	RelationshipTrash(relationshipID, deletedBy string) (bool, error)
	RelationshipRestore(relationshipID string) (bool, error)
	RelationshipDelete(relationshipID string) (bool, error)
	RelationshipDeleteAll(entityID string) error

	// Taxonomies (requires TaxonomiesEnabled)
	TaxonomyCreate(options entitystore.TaxonomyOptions) (ActiveTaxonomyInterface, error)
	TaxonomyFindByID(taxonomyID string) (ActiveTaxonomyInterface, error)
	TaxonomyFindBySlug(slug string) (ActiveTaxonomyInterface, error)
	TaxonomyList(query entitystore.TaxonomyQueryInterface) ([]ActiveTaxonomyInterface, error)
	TaxonomyCount(query entitystore.TaxonomyQueryInterface) (int64, error)
	TaxonomyUpdate(taxonomy entitystore.TaxonomyInterface) error
	TaxonomyTrash(taxonomyID, deletedBy string) (bool, error)
	TaxonomyRestore(taxonomyID string) (bool, error)
	TaxonomyDelete(taxonomyID string) (bool, error)

	// Taxonomy terms (requires TaxonomiesEnabled)
	TermCreate(options entitystore.TaxonomyTermOptions) (ActiveTaxonomyTermInterface, error)
	TermFindByID(termID string) (ActiveTaxonomyTermInterface, error)
	TermFindBySlug(taxonomyID, slug string) (ActiveTaxonomyTermInterface, error)
	TermList(query entitystore.TaxonomyTermQueryInterface) ([]ActiveTaxonomyTermInterface, error)
	TermCount(query entitystore.TaxonomyTermQueryInterface) (int64, error)
	TermUpdate(term entitystore.TaxonomyTermInterface) error
	TermTrash(termID, deletedBy string) (bool, error)
	TermRestore(termID string) (bool, error)
	TermDelete(termID string) (bool, error)
}

// activeStoreImplementation is the concrete store wrapper.
type activeStoreImplementation struct {
	ctx   context.Context
	store entitystore.StoreInterface
}

// New converts a standard store into an ActiveStoreInterface — the
// single entry point for the package.
func New(ctx context.Context, store entitystore.StoreInterface) (ActiveStoreInterface, error) {
	if store == nil {
		return nil, errors.New("store cannot be nil")
	}
	return &activeStoreImplementation{ctx: ctx, store: store}, nil
}

// EntityCreate returns a new staged entity — the row is not inserted
// until Save() is called on it (unlike store.EntityCreate which
// persists immediately).
func (s *activeStoreImplementation) EntityCreate(entityType string) ActiveEntityInterface {
	entity := entitystore.NewEntity()
	entity.SetType(entityType)
	return &activeEntityImplementation{
		ctx:        s.ctx,
		entity:     entity,
		store:      s.store,
		persisted:  false,
		pendingOps: []pendingOp{},
	}
}

func (s *activeStoreImplementation) EntityFindByID(entityID string) (ActiveEntityInterface, error) {
	ent, err := s.store.EntityFindByID(s.ctx, entityID)
	if err != nil || ent == nil {
		return nil, err
	}
	return newActiveEntity(s.ctx, s.store, ent)
}

func (s *activeStoreImplementation) EntityList(query entitystore.EntityQueryInterface) ([]ActiveEntityInterface, error) {
	entities, err := s.store.EntityList(s.ctx, query)
	if err != nil {
		return nil, err
	}
	active := make([]ActiveEntityInterface, 0, len(entities))
	for _, ent := range entities {
		wrapped, err := newActiveEntity(s.ctx, s.store, ent)
		if err != nil {
			return nil, err
		}
		active = append(active, wrapped)
	}
	return active, nil
}

func (s *activeStoreImplementation) EntityCount(query entitystore.EntityQueryInterface) (int64, error) {
	return s.store.EntityCount(s.ctx, query)
}

func (s *activeStoreImplementation) EntityTrash(entityID string) (bool, error) {
	return s.store.EntityTrash(s.ctx, entityID)
}

func (s *activeStoreImplementation) EntityDelete(entityID string) (bool, error) {
	return s.store.EntityDelete(s.ctx, entityID)
}

func (s *activeStoreImplementation) EntityWrap(entity entitystore.EntityInterface) (ActiveEntityInterface, error) {
	return newActiveEntity(s.ctx, s.store, entity)
}

func (s *activeStoreImplementation) GetStore() entitystore.StoreInterface {
	return s.store
}
