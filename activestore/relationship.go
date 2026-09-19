package activestore

import (
	"context"
	"errors"

	"github.com/dracory/entitystore"
)

// ActiveRelationshipInterface wraps a RelationshipInterface with
// navigation and lifecycle helpers.
type ActiveRelationshipInterface interface {
	GetRelationship() entitystore.RelationshipInterface
	GetEntity() (ActiveEntityInterface, error)        // source entity, hydrated
	GetRelatedEntity() (ActiveEntityInterface, error) // target entity, hydrated
	Trash(deletedBy string) (bool, error)
	Restore() (bool, error)
	Delete() (bool, error)
}

// activeRelationshipImplementation is the concrete relationship wrapper.
type activeRelationshipImplementation struct {
	ctx          context.Context
	store        entitystore.StoreInterface
	relationship entitystore.RelationshipInterface
}

func newActiveRelationship(ctx context.Context, store entitystore.StoreInterface, relationship entitystore.RelationshipInterface) (ActiveRelationshipInterface, error) {
	if relationship == nil {
		return nil, errors.New("relationship cannot be nil")
	}
	return &activeRelationshipImplementation{ctx: ctx, store: store, relationship: relationship}, nil
}

func (r *activeRelationshipImplementation) GetRelationship() entitystore.RelationshipInterface {
	return r.relationship
}

func (r *activeRelationshipImplementation) GetEntity() (ActiveEntityInterface, error) {
	ent, err := r.store.EntityFindByID(r.ctx, r.relationship.GetEntityID())
	if err != nil {
		return nil, err
	}
	return newActiveEntity(r.ctx, r.store, ent)
}

func (r *activeRelationshipImplementation) GetRelatedEntity() (ActiveEntityInterface, error) {
	ent, err := r.store.EntityFindByID(r.ctx, r.relationship.GetRelatedEntityID())
	if err != nil {
		return nil, err
	}
	return newActiveEntity(r.ctx, r.store, ent)
}

func (r *activeRelationshipImplementation) Trash(deletedBy string) (bool, error) {
	return r.store.RelationshipTrash(r.ctx, r.relationship.ID(), deletedBy)
}

func (r *activeRelationshipImplementation) Restore() (bool, error) {
	return r.store.RelationshipRestore(r.ctx, r.relationship.ID())
}

func (r *activeRelationshipImplementation) Delete() (bool, error) {
	return r.store.RelationshipDelete(r.ctx, r.relationship.ID())
}

// ---------------------------------------------------------------------------
// Store-level methods
// ---------------------------------------------------------------------------

func (s *activeStoreImplementation) RelationshipCreate(options entitystore.RelationshipOptions) (ActiveRelationshipInterface, error) {
	rel, err := s.store.RelationshipCreateByOptions(s.ctx, options)
	if err != nil {
		return nil, err
	}
	return newActiveRelationship(s.ctx, s.store, rel)
}

func (s *activeStoreImplementation) RelationshipFindByID(relationshipID string) (ActiveRelationshipInterface, error) {
	rel, err := s.store.RelationshipFind(s.ctx, relationshipID)
	if err != nil {
		return nil, err
	}
	return newActiveRelationship(s.ctx, s.store, rel)
}

func (s *activeStoreImplementation) RelationshipList(query entitystore.RelationshipQueryInterface) ([]ActiveRelationshipInterface, error) {
	rels, err := s.store.RelationshipList(s.ctx, query)
	if err != nil {
		return nil, err
	}
	active := make([]ActiveRelationshipInterface, 0, len(rels))
	for _, rel := range rels {
		wrapped, err := newActiveRelationship(s.ctx, s.store, rel)
		if err != nil {
			return nil, err
		}
		active = append(active, wrapped)
	}
	return active, nil
}

func (s *activeStoreImplementation) RelationshipCount(query entitystore.RelationshipQueryInterface) (int64, error) {
	return s.store.RelationshipCount(s.ctx, query)
}

func (s *activeStoreImplementation) RelationshipTrash(relationshipID, deletedBy string) (bool, error) {
	return s.store.RelationshipTrash(s.ctx, relationshipID, deletedBy)
}

func (s *activeStoreImplementation) RelationshipRestore(relationshipID string) (bool, error) {
	return s.store.RelationshipRestore(s.ctx, relationshipID)
}

func (s *activeStoreImplementation) RelationshipDelete(relationshipID string) (bool, error) {
	return s.store.RelationshipDelete(s.ctx, relationshipID)
}

func (s *activeStoreImplementation) RelationshipDeleteAll(entityID string) error {
	return s.store.RelationshipDeleteAll(s.ctx, entityID)
}
