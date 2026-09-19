package entitystore

import "errors"

// entity_taxonomy_query.go defines EntityTaxonomyQueryInterface — the
// fluent query type accepted by StoreInterface.EntityTaxonomyList and
// EntityTaxonomyCount — and its implementation. The implementation uses
// typed fields with has* presence flags instead of a map[string]any, so
// every setter is compile-time type safe. Presence flags let the store read
// only the fields the caller explicitly set, and let Validate() distinguish
// "set to an invalid value" from "not set".
//
// The store reads the query through the Get* accessors and calls Validate()
// before executing, returning an error for invalid queries.
//
// Usage:
//
//	store.EntityTaxonomyList(ctx, EntityTaxonomyQuery().
//		WithEntityID("entity-1").
//		WithTaxonomyID("taxonomy-1").
//		WithLimit(100))

// == INTERFACE ==============================================================

// EntityTaxonomyQueryInterface is a fluent entity-taxonomy query. Every
// With* method marks the field as set and returns the builder for chaining.
// The store calls Validate() before executing and returns an error for
// invalid queries.
type EntityTaxonomyQueryInterface interface {
	// Validate checks the set parameters for invalid values (empty strings
	// or slices where a value was set, invalid sort order).
	Validate() error

	// HasID reports whether the ID filter was set.
	HasID() bool
	// GetID returns the assignment ID filter.
	GetID() string
	// WithID filters by a specific assignment ID.
	WithID(id string) EntityTaxonomyQueryInterface

	// HasEntityID reports whether the EntityID filter was set.
	HasEntityID() bool
	// GetEntityID returns the entity ID filter.
	GetEntityID() string
	// WithEntityID filters by entity ID.
	WithEntityID(entityID string) EntityTaxonomyQueryInterface

	// HasEntityIDs reports whether the EntityIDs filter was set.
	HasEntityIDs() bool
	// GetEntityIDs returns the entity IDs filter.
	GetEntityIDs() []string
	// WithEntityIDs filters by multiple entity IDs.
	WithEntityIDs(entityIDs []string) EntityTaxonomyQueryInterface

	// HasTaxonomyID reports whether the TaxonomyID filter was set.
	HasTaxonomyID() bool
	// GetTaxonomyID returns the taxonomy ID filter.
	GetTaxonomyID() string
	// WithTaxonomyID filters by taxonomy ID.
	WithTaxonomyID(taxonomyID string) EntityTaxonomyQueryInterface

	// HasTermID reports whether the TermID filter was set.
	HasTermID() bool
	// GetTermID returns the term ID filter.
	GetTermID() string
	// WithTermID filters by term ID.
	WithTermID(termID string) EntityTaxonomyQueryInterface

	// HasTermIDs reports whether the TermIDs filter was set.
	HasTermIDs() bool
	// GetTermIDs returns the term IDs filter.
	GetTermIDs() []string
	// WithTermIDs filters by multiple term IDs.
	WithTermIDs(termIDs []string) EntityTaxonomyQueryInterface

	// HasLimit reports whether Limit was set.
	HasLimit() bool
	// GetLimit returns the result limit.
	GetLimit() uint64
	// WithLimit sets the maximum number of results.
	WithLimit(limit uint64) EntityTaxonomyQueryInterface

	// HasOffset reports whether Offset was set.
	HasOffset() bool
	// GetOffset returns the result offset.
	GetOffset() uint64
	// WithOffset sets the number of results to skip.
	WithOffset(offset uint64) EntityTaxonomyQueryInterface

	// HasSortBy reports whether SortBy was set.
	HasSortBy() bool
	// GetSortBy returns the sort column.
	GetSortBy() string
	// WithSortBy sets the column to sort by.
	WithSortBy(sortBy string) EntityTaxonomyQueryInterface

	// HasSortOrder reports whether SortOrder was set.
	HasSortOrder() bool
	// GetSortOrder returns the sort direction.
	GetSortOrder() string
	// WithSortOrder sets the sort direction ("asc" or "desc").
	WithSortOrder(sortOrder string) EntityTaxonomyQueryInterface

	// HasCountOnly reports whether CountOnly was set.
	HasCountOnly() bool
	// GetCountOnly returns whether only a count is requested.
	GetCountOnly() bool
	// WithCountOnly returns only a count, not results.
	WithCountOnly(countOnly bool) EntityTaxonomyQueryInterface
}

// == CONSTRUCTOR ============================================================

// EntityTaxonomyQuery returns a new fluent entity-taxonomy query builder.
func EntityTaxonomyQuery() EntityTaxonomyQueryInterface {
	return &entityTaxonomyQueryImplementation{}
}

// == TYPE ===================================================================

type entityTaxonomyQueryImplementation struct {
	id            string
	hasID         bool
	entityID      string
	hasEntityID   bool
	entityIDs     []string
	hasEntityIDs  bool
	taxonomyID    string
	hasTaxonomyID bool
	termID        string
	hasTermID     bool
	termIDs       []string
	hasTermIDs    bool
	limit         uint64
	hasLimit      bool
	offset        uint64
	hasOffset     bool
	sortBy        string
	hasSortBy     bool
	sortOrder     string
	hasSortOrder  bool
	countOnly     bool
	hasCountOnly  bool
}

// == INTERFACE VERIFICATION =================================================

var _ EntityTaxonomyQueryInterface = (*entityTaxonomyQueryImplementation)(nil)

// == INTERFACE IMPLEMENTATION ===============================================

// Validate checks the set parameters for invalid values.
func (q *entityTaxonomyQueryImplementation) Validate() error {
	if q.hasID && q.id == "" {
		return errors.New("entity taxonomy query: id cannot be empty")
	}
	if q.hasEntityID && q.entityID == "" {
		return errors.New("entity taxonomy query: entity_id cannot be empty")
	}
	if q.hasEntityIDs && len(q.entityIDs) < 1 {
		return errors.New("entity taxonomy query: entity_ids cannot be empty array")
	}
	if q.hasTaxonomyID && q.taxonomyID == "" {
		return errors.New("entity taxonomy query: taxonomy_id cannot be empty")
	}
	if q.hasTermID && q.termID == "" {
		return errors.New("entity taxonomy query: term_id cannot be empty")
	}
	if q.hasTermIDs && len(q.termIDs) < 1 {
		return errors.New("entity taxonomy query: term_ids cannot be empty array")
	}
	if q.hasSortBy && q.sortBy == "" {
		return errors.New("entity taxonomy query: sort_by cannot be empty")
	}
	if q.hasSortOrder && q.sortOrder != "asc" && q.sortOrder != "desc" {
		return errors.New("entity taxonomy query: sort_order must be \"asc\" or \"desc\"")
	}
	return nil
}

func (q *entityTaxonomyQueryImplementation) HasID() bool   { return q.hasID }
func (q *entityTaxonomyQueryImplementation) GetID() string { return q.id }
func (q *entityTaxonomyQueryImplementation) WithID(id string) EntityTaxonomyQueryInterface {
	q.id, q.hasID = id, true
	return q
}

func (q *entityTaxonomyQueryImplementation) HasEntityID() bool   { return q.hasEntityID }
func (q *entityTaxonomyQueryImplementation) GetEntityID() string { return q.entityID }
func (q *entityTaxonomyQueryImplementation) WithEntityID(entityID string) EntityTaxonomyQueryInterface {
	q.entityID, q.hasEntityID = entityID, true
	return q
}

func (q *entityTaxonomyQueryImplementation) HasEntityIDs() bool     { return q.hasEntityIDs }
func (q *entityTaxonomyQueryImplementation) GetEntityIDs() []string { return q.entityIDs }
func (q *entityTaxonomyQueryImplementation) WithEntityIDs(entityIDs []string) EntityTaxonomyQueryInterface {
	q.entityIDs, q.hasEntityIDs = entityIDs, true
	return q
}

func (q *entityTaxonomyQueryImplementation) HasTaxonomyID() bool   { return q.hasTaxonomyID }
func (q *entityTaxonomyQueryImplementation) GetTaxonomyID() string { return q.taxonomyID }
func (q *entityTaxonomyQueryImplementation) WithTaxonomyID(taxonomyID string) EntityTaxonomyQueryInterface {
	q.taxonomyID, q.hasTaxonomyID = taxonomyID, true
	return q
}

func (q *entityTaxonomyQueryImplementation) HasTermID() bool   { return q.hasTermID }
func (q *entityTaxonomyQueryImplementation) GetTermID() string { return q.termID }
func (q *entityTaxonomyQueryImplementation) WithTermID(termID string) EntityTaxonomyQueryInterface {
	q.termID, q.hasTermID = termID, true
	return q
}

func (q *entityTaxonomyQueryImplementation) HasTermIDs() bool     { return q.hasTermIDs }
func (q *entityTaxonomyQueryImplementation) GetTermIDs() []string { return q.termIDs }
func (q *entityTaxonomyQueryImplementation) WithTermIDs(termIDs []string) EntityTaxonomyQueryInterface {
	q.termIDs, q.hasTermIDs = termIDs, true
	return q
}

func (q *entityTaxonomyQueryImplementation) HasLimit() bool   { return q.hasLimit }
func (q *entityTaxonomyQueryImplementation) GetLimit() uint64 { return q.limit }
func (q *entityTaxonomyQueryImplementation) WithLimit(limit uint64) EntityTaxonomyQueryInterface {
	q.limit, q.hasLimit = limit, true
	return q
}

func (q *entityTaxonomyQueryImplementation) HasOffset() bool   { return q.hasOffset }
func (q *entityTaxonomyQueryImplementation) GetOffset() uint64 { return q.offset }
func (q *entityTaxonomyQueryImplementation) WithOffset(offset uint64) EntityTaxonomyQueryInterface {
	q.offset, q.hasOffset = offset, true
	return q
}

func (q *entityTaxonomyQueryImplementation) HasSortBy() bool   { return q.hasSortBy }
func (q *entityTaxonomyQueryImplementation) GetSortBy() string { return q.sortBy }
func (q *entityTaxonomyQueryImplementation) WithSortBy(sortBy string) EntityTaxonomyQueryInterface {
	q.sortBy, q.hasSortBy = sortBy, true
	return q
}

func (q *entityTaxonomyQueryImplementation) HasSortOrder() bool   { return q.hasSortOrder }
func (q *entityTaxonomyQueryImplementation) GetSortOrder() string { return q.sortOrder }
func (q *entityTaxonomyQueryImplementation) WithSortOrder(sortOrder string) EntityTaxonomyQueryInterface {
	q.sortOrder, q.hasSortOrder = sortOrder, true
	return q
}

func (q *entityTaxonomyQueryImplementation) HasCountOnly() bool { return q.hasCountOnly }
func (q *entityTaxonomyQueryImplementation) GetCountOnly() bool { return q.countOnly }
func (q *entityTaxonomyQueryImplementation) WithCountOnly(countOnly bool) EntityTaxonomyQueryInterface {
	q.countOnly, q.hasCountOnly = countOnly, true
	return q
}
