package entitystore

import "errors"

// entity_query.go defines EntityQueryInterface — the fluent query type
// accepted by StoreInterface.EntityList and EntityCount — and its
// implementation. The implementation uses typed fields with has* presence
// flags instead of a map[string]any, so every setter is compile-time type
// safe. Presence flags let the store read only the fields the caller
// explicitly set, and let Validate() distinguish "set to an invalid value"
// from "not set".
//
// The store reads the query through the Get* accessors and calls Validate()
// before executing, returning an error for invalid queries.
//
// Usage:
//
//	store.EntityList(ctx, EntityQuery().
//		WithEntityType("character").
//		WithLimit(100))

// == INTERFACE ==============================================================

// EntityQueryInterface is a fluent entity query. Every With* method marks
// the field as set and returns the builder for chaining. The store calls
// Validate() before executing and returns an error for invalid queries.
type EntityQueryInterface interface {
	// Validate checks the set parameters for invalid values (empty strings
	// or slices where a value was set, invalid sort order).
	Validate() error

	// HasID reports whether the ID filter was set.
	HasID() bool
	// GetID returns the entity ID filter.
	GetID() string
	// WithID filters by a specific entity ID.
	WithID(id string) EntityQueryInterface

	// HasIDs reports whether the IDs filter was set.
	HasIDs() bool
	// GetIDs returns the entity IDs filter.
	GetIDs() []string
	// WithIDs filters by multiple entity IDs.
	WithIDs(ids []string) EntityQueryInterface

	// HasEntityType reports whether the EntityType filter was set.
	HasEntityType() bool
	// GetEntityType returns the entity type filter.
	GetEntityType() string
	// WithEntityType filters by entity type.
	WithEntityType(entityType string) EntityQueryInterface

	// HasEntityHandle reports whether the EntityHandle filter was set.
	HasEntityHandle() bool
	// GetEntityHandle returns the entity handle filter.
	GetEntityHandle() string
	// WithEntityHandle filters by entity handle.
	WithEntityHandle(entityHandle string) EntityQueryInterface

	// HasSearch reports whether the Search filter was set.
	HasSearch() bool
	// GetSearch returns the text search filter.
	GetSearch() string
	// WithSearch sets the text search filter.
	WithSearch(search string) EntityQueryInterface

	// HasLimit reports whether Limit was set.
	HasLimit() bool
	// GetLimit returns the result limit.
	GetLimit() uint64
	// WithLimit sets the maximum number of results.
	WithLimit(limit uint64) EntityQueryInterface

	// HasOffset reports whether Offset was set.
	HasOffset() bool
	// GetOffset returns the result offset.
	GetOffset() uint64
	// WithOffset sets the number of results to skip.
	WithOffset(offset uint64) EntityQueryInterface

	// HasSortBy reports whether SortBy was set.
	HasSortBy() bool
	// GetSortBy returns the sort column.
	GetSortBy() string
	// WithSortBy sets the column to sort by.
	WithSortBy(sortBy string) EntityQueryInterface

	// HasSortOrder reports whether SortOrder was set.
	HasSortOrder() bool
	// GetSortOrder returns the sort direction.
	GetSortOrder() string
	// WithSortOrder sets the sort direction ("asc" or "desc").
	WithSortOrder(sortOrder string) EntityQueryInterface

	// HasCountOnly reports whether CountOnly was set.
	HasCountOnly() bool
	// GetCountOnly returns whether only a count is requested.
	GetCountOnly() bool
	// WithCountOnly returns only a count, not results.
	WithCountOnly(countOnly bool) EntityQueryInterface

	// HasCreatedAtGte reports whether the created_at lower bound was set.
	HasCreatedAtGte() bool
	// GetCreatedAtGte returns the created_at lower bound (inclusive).
	GetCreatedAtGte() string
	// WithCreatedAtGte filters to rows created at or after the given UTC
	// datetime ("YYYY-MM-DD HH:MM:SS").
	WithCreatedAtGte(createdAtGte string) EntityQueryInterface

	// HasCreatedAtLte reports whether the created_at upper bound was set.
	HasCreatedAtLte() bool
	// GetCreatedAtLte returns the created_at upper bound (inclusive).
	GetCreatedAtLte() string
	// WithCreatedAtLte filters to rows created at or before the given UTC
	// datetime ("YYYY-MM-DD HH:MM:SS").
	WithCreatedAtLte(createdAtLte string) EntityQueryInterface

	// HasUpdatedAtGte reports whether the updated_at lower bound was set.
	HasUpdatedAtGte() bool
	// GetUpdatedAtGte returns the updated_at lower bound (inclusive).
	GetUpdatedAtGte() string
	// WithUpdatedAtGte filters to rows updated at or after the given UTC
	// datetime ("YYYY-MM-DD HH:MM:SS").
	WithUpdatedAtGte(updatedAtGte string) EntityQueryInterface

	// HasUpdatedAtLte reports whether the updated_at upper bound was set.
	HasUpdatedAtLte() bool
	// GetUpdatedAtLte returns the updated_at upper bound (inclusive).
	GetUpdatedAtLte() string
	// WithUpdatedAtLte filters to rows updated at or before the given UTC
	// datetime ("YYYY-MM-DD HH:MM:SS").
	WithUpdatedAtLte(updatedAtLte string) EntityQueryInterface

	// HasPrefetchAttributes reports whether attribute prefetching was set.
	HasPrefetchAttributes() bool
	// GetPrefetchAttributes returns the attribute keys to prefetch. An empty
	// slice means "all attributes" when HasPrefetchAttributes() is true.
	GetPrefetchAttributes() []string
	// WithPrefetchAttributes eagerly loads attributes for all returned
	// entities in a single batch query, populating each entity's in-memory
	// attributes (readable via GetTempKey). Calling it with specific keys
	// prefetches only those attributes; calling it with no arguments
	// prefetches all attributes. Use it to avoid an N+1 query pattern when
	// you know which attributes you will read — or want them all.
	WithPrefetchAttributes(attributeKeys ...string) EntityQueryInterface
}

// == CONSTRUCTOR ============================================================

// EntityQuery returns a new fluent entity query builder.
func EntityQuery() EntityQueryInterface {
	return &entityQueryImplementation{}
}

// == TYPE ===================================================================

type entityQueryImplementation struct {
	id              string
	hasID           bool
	ids             []string
	hasIDs          bool
	entityType      string
	hasEntityType   bool
	entityHandle    string
	hasEntityHandle bool
	search          string
	hasSearch       bool
	limit           uint64
	hasLimit        bool
	offset          uint64
	hasOffset       bool
	sortBy          string
	hasSortBy       bool
	sortOrder       string
	hasSortOrder    bool
	countOnly       bool
	hasCountOnly    bool
	createdAtGte    string
	hasCreatedAtGte bool
	createdAtLte    string
	hasCreatedAtLte bool
	updatedAtGte    string
	hasUpdatedAtGte bool
	updatedAtLte    string
	hasUpdatedAtLte bool

	prefetchAttributes    []string
	hasPrefetchAttributes bool
}

// == INTERFACE VERIFICATION =================================================

var _ EntityQueryInterface = (*entityQueryImplementation)(nil)

// == INTERFACE IMPLEMENTATION ===============================================

// Validate checks the set parameters for invalid values.
func (q *entityQueryImplementation) Validate() error {
	if q.hasID && q.id == "" {
		return errors.New("entity query: id cannot be empty")
	}
	if q.hasIDs && len(q.ids) < 1 {
		return errors.New("entity query: ids cannot be empty array")
	}
	if q.hasEntityType && q.entityType == "" {
		return errors.New("entity query: entity_type cannot be empty")
	}
	if q.hasEntityHandle && q.entityHandle == "" {
		return errors.New("entity query: entity_handle cannot be empty")
	}
	if q.hasSortBy && q.sortBy == "" {
		return errors.New("entity query: sort_by cannot be empty")
	}
	if q.hasSortOrder && q.sortOrder != "asc" && q.sortOrder != "desc" {
		return errors.New("entity query: sort_order must be \"asc\" or \"desc\"")
	}
	if q.hasCreatedAtGte && q.createdAtGte == "" {
		return errors.New("entity query: created_at_gte cannot be empty")
	}
	if q.hasCreatedAtLte && q.createdAtLte == "" {
		return errors.New("entity query: created_at_lte cannot be empty")
	}
	if q.hasUpdatedAtGte && q.updatedAtGte == "" {
		return errors.New("entity query: updated_at_gte cannot be empty")
	}
	if q.hasUpdatedAtLte && q.updatedAtLte == "" {
		return errors.New("entity query: updated_at_lte cannot be empty")
	}
	return nil
}

func (q *entityQueryImplementation) HasID() bool   { return q.hasID }
func (q *entityQueryImplementation) GetID() string { return q.id }
func (q *entityQueryImplementation) WithID(id string) EntityQueryInterface {
	q.id, q.hasID = id, true
	return q
}

func (q *entityQueryImplementation) HasIDs() bool     { return q.hasIDs }
func (q *entityQueryImplementation) GetIDs() []string { return q.ids }
func (q *entityQueryImplementation) WithIDs(ids []string) EntityQueryInterface {
	q.ids, q.hasIDs = ids, true
	return q
}

func (q *entityQueryImplementation) HasEntityType() bool   { return q.hasEntityType }
func (q *entityQueryImplementation) GetEntityType() string { return q.entityType }
func (q *entityQueryImplementation) WithEntityType(entityType string) EntityQueryInterface {
	q.entityType, q.hasEntityType = entityType, true
	return q
}

func (q *entityQueryImplementation) HasEntityHandle() bool   { return q.hasEntityHandle }
func (q *entityQueryImplementation) GetEntityHandle() string { return q.entityHandle }
func (q *entityQueryImplementation) WithEntityHandle(entityHandle string) EntityQueryInterface {
	q.entityHandle, q.hasEntityHandle = entityHandle, true
	return q
}

func (q *entityQueryImplementation) HasSearch() bool   { return q.hasSearch }
func (q *entityQueryImplementation) GetSearch() string { return q.search }
func (q *entityQueryImplementation) WithSearch(search string) EntityQueryInterface {
	q.search, q.hasSearch = search, true
	return q
}

func (q *entityQueryImplementation) HasLimit() bool   { return q.hasLimit }
func (q *entityQueryImplementation) GetLimit() uint64 { return q.limit }
func (q *entityQueryImplementation) WithLimit(limit uint64) EntityQueryInterface {
	q.limit, q.hasLimit = limit, true
	return q
}

func (q *entityQueryImplementation) HasOffset() bool   { return q.hasOffset }
func (q *entityQueryImplementation) GetOffset() uint64 { return q.offset }
func (q *entityQueryImplementation) WithOffset(offset uint64) EntityQueryInterface {
	q.offset, q.hasOffset = offset, true
	return q
}

func (q *entityQueryImplementation) HasSortBy() bool   { return q.hasSortBy }
func (q *entityQueryImplementation) GetSortBy() string { return q.sortBy }
func (q *entityQueryImplementation) WithSortBy(sortBy string) EntityQueryInterface {
	q.sortBy, q.hasSortBy = sortBy, true
	return q
}

func (q *entityQueryImplementation) HasSortOrder() bool   { return q.hasSortOrder }
func (q *entityQueryImplementation) GetSortOrder() string { return q.sortOrder }
func (q *entityQueryImplementation) WithSortOrder(sortOrder string) EntityQueryInterface {
	q.sortOrder, q.hasSortOrder = sortOrder, true
	return q
}

func (q *entityQueryImplementation) HasCountOnly() bool { return q.hasCountOnly }
func (q *entityQueryImplementation) GetCountOnly() bool { return q.countOnly }
func (q *entityQueryImplementation) WithCountOnly(countOnly bool) EntityQueryInterface {
	q.countOnly, q.hasCountOnly = countOnly, true
	return q
}

func (q *entityQueryImplementation) HasCreatedAtGte() bool   { return q.hasCreatedAtGte }
func (q *entityQueryImplementation) GetCreatedAtGte() string { return q.createdAtGte }
func (q *entityQueryImplementation) WithCreatedAtGte(createdAtGte string) EntityQueryInterface {
	q.createdAtGte, q.hasCreatedAtGte = createdAtGte, true
	return q
}

func (q *entityQueryImplementation) HasCreatedAtLte() bool   { return q.hasCreatedAtLte }
func (q *entityQueryImplementation) GetCreatedAtLte() string { return q.createdAtLte }
func (q *entityQueryImplementation) WithCreatedAtLte(createdAtLte string) EntityQueryInterface {
	q.createdAtLte, q.hasCreatedAtLte = createdAtLte, true
	return q
}

func (q *entityQueryImplementation) HasUpdatedAtGte() bool   { return q.hasUpdatedAtGte }
func (q *entityQueryImplementation) GetUpdatedAtGte() string { return q.updatedAtGte }
func (q *entityQueryImplementation) WithUpdatedAtGte(updatedAtGte string) EntityQueryInterface {
	q.updatedAtGte, q.hasUpdatedAtGte = updatedAtGte, true
	return q
}

func (q *entityQueryImplementation) HasUpdatedAtLte() bool   { return q.hasUpdatedAtLte }
func (q *entityQueryImplementation) GetUpdatedAtLte() string { return q.updatedAtLte }
func (q *entityQueryImplementation) WithUpdatedAtLte(updatedAtLte string) EntityQueryInterface {
	q.updatedAtLte, q.hasUpdatedAtLte = updatedAtLte, true
	return q
}

func (q *entityQueryImplementation) HasPrefetchAttributes() bool {
	return q.hasPrefetchAttributes
}
func (q *entityQueryImplementation) GetPrefetchAttributes() []string {
	return q.prefetchAttributes
}
func (q *entityQueryImplementation) WithPrefetchAttributes(attributeKeys ...string) EntityQueryInterface {
	q.prefetchAttributes, q.hasPrefetchAttributes = attributeKeys, true
	return q
}
