package activestore

import (
	"context"
	"errors"
	"strconv"

	"github.com/dracory/entitystore"
)

// ActiveEntityInterface defines the rich, store-aware API for entities.
// Mutating methods return the ActiveEntityInterface for chaining; errors are
// accumulated internally and surfaced by Save() or Err().
type ActiveEntityInterface interface {
	SetString(key, value string) ActiveEntityInterface
	SetInt(key string, value int64) ActiveEntityInterface
	SetFloat(key string, value float64) ActiveEntityInterface
	GetString(key string) (string, bool, error)
	GetAttributes() ([]entitystore.AttributeInterface, error)
	Prefetch() error
	Save() error
	Trash() (bool, error)
	Delete() (bool, error)
	Err() error
	GetEntity() entitystore.EntityInterface
}

// pendingOp is a deferred store write applied when Save() persists the
// entity for the first time.
type pendingOp func(ctx context.Context, store entitystore.StoreInterface, entityID string) error

// activeEntityImplementation is the concrete wrapper.
type activeEntityImplementation struct {
	ctx        context.Context
	entity     entitystore.EntityInterface
	store      entitystore.StoreInterface
	persisted  bool
	pendingOps []pendingOp
	cache      map[string]string // populated by Prefetch(); nil = not prefetched
	err        error             // first accumulated error, deferred to Save()/Err()
}

// newActiveEntity binds an existing (persisted) entity to a store context.
// Entities are only created through ActiveStoreInterface, so a store is
// always present.
func newActiveEntity(ctx context.Context, store entitystore.StoreInterface, entity entitystore.EntityInterface) (ActiveEntityInterface, error) {
	if entity == nil {
		return nil, errors.New("entity cannot be nil")
	}
	return &activeEntityImplementation{
		ctx:       ctx,
		entity:    entity,
		store:     store,
		persisted: true,
	}, nil
}

// apply runs op immediately when the entity is persisted, or stages it
// until Save() otherwise. The first error is accumulated into e.err.
func (e *activeEntityImplementation) apply(op pendingOp) ActiveEntityInterface {
	if e.err != nil {
		return e
	}
	if !e.persisted {
		e.pendingOps = append(e.pendingOps, op)
		return e
	}
	e.err = op(e.ctx, e.store, e.entity.ID())
	return e
}

// setCached records a successful write in the prefetch cache, if active.
func (e *activeEntityImplementation) setCached(key, value string) {
	if e.cache != nil {
		e.cache[key] = value
	}
}

func (e *activeEntityImplementation) SetString(key, value string) ActiveEntityInterface {
	return e.apply(func(ctx context.Context, store entitystore.StoreInterface, entityID string) error {
		if err := store.AttributeSetString(ctx, entityID, key, value); err != nil {
			return err
		}
		e.setCached(key, value)
		return nil
	})
}

func (e *activeEntityImplementation) SetInt(key string, value int64) ActiveEntityInterface {
	return e.apply(func(ctx context.Context, store entitystore.StoreInterface, entityID string) error {
		if err := store.AttributeSetInt(ctx, entityID, key, value); err != nil {
			return err
		}
		e.setCached(key, strconv.FormatInt(value, 10))
		return nil
	})
}

func (e *activeEntityImplementation) SetFloat(key string, value float64) ActiveEntityInterface {
	return e.apply(func(ctx context.Context, store entitystore.StoreInterface, entityID string) error {
		if err := store.AttributeSetFloat(ctx, entityID, key, value); err != nil {
			return err
		}
		e.setCached(key, strconv.FormatFloat(value, 'f', -1, 64))
		return nil
	})
}

// Prefetch loads all of the entity's attributes into memory in a single
// query. Afterwards GetString reads from the cache without hitting the
// store, and setters keep the cache in sync.
func (e *activeEntityImplementation) Prefetch() error {
	attrs, err := e.store.EntityAttributeList(e.ctx, e.entity.ID())
	if err != nil {
		return err
	}
	e.cache = make(map[string]string, len(attrs))
	for _, attr := range attrs {
		e.cache[attr.GetKey()] = attr.GetValue()
	}
	return nil
}

// GetString gets a string attribute — from the prefetch cache when
// Prefetch() was called, otherwise from the store.
func (e *activeEntityImplementation) GetString(key string) (string, bool, error) {
	if e.cache != nil {
		value, exists := e.cache[key]
		return value, exists, nil
	}
	return e.store.AttributeGetString(e.ctx, e.entity.ID(), key)
}

// GetAttributes returns all attributes of the entity via the store.
func (e *activeEntityImplementation) GetAttributes() ([]entitystore.AttributeInterface, error) {
	return e.store.EntityAttributeList(e.ctx, e.entity.ID())
}

// Err returns the first error accumulated by the fluent setters, if any.
func (e *activeEntityImplementation) Err() error {
	return e.err
}

// Save persists the underlying entity to the store, then flushes any
// attribute writes that were staged while the entity was unpersisted.
func (e *activeEntityImplementation) Save() error {
	if e.err != nil {
		return e.err
	}
	if e.persisted {
		return e.store.EntityUpdate(e.ctx, e.entity)
	}
	if err := e.store.EntityCreate(e.ctx, e.entity); err != nil {
		return err
	}
	e.persisted = true
	for _, op := range e.pendingOps {
		if err := op(e.ctx, e.store, e.entity.ID()); err != nil {
			e.err = err
			return err
		}
	}
	e.pendingOps = nil
	return nil
}

// Trash moves the entity to the trash table.
func (e *activeEntityImplementation) Trash() (bool, error) {
	return e.store.EntityTrash(e.ctx, e.entity.ID())
}

// Delete permanently removes the entity.
func (e *activeEntityImplementation) Delete() (bool, error) {
	return e.store.EntityDelete(e.ctx, e.entity.ID())
}

// GetEntity returns the underlying pure data entity if the developer
// needs to pass it to a core method that expects EntityInterface.
func (e *activeEntityImplementation) GetEntity() entitystore.EntityInterface {
	return e.entity
}
