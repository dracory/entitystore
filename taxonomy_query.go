package entitystore

import "errors"

// taxonomy_query.go defines TaxonomyQueryInterface — the fluent query type
// accepted by StoreInterface.TaxonomyList, TaxonomyCount and
// TaxonomyTrashList — and its implementation. The implementation uses
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
//	store.TaxonomyList(ctx, TaxonomyQuery().
//		WithSlug("categories").
//		WithEntityTypes([]string{"product"}).
//		WithLimit(100))

// TaxonomyOptions provides the options for creating a new taxonomy
type TaxonomyOptions struct {
	Name        string   // Display name of the taxonomy
	Slug        string   // URL-friendly identifier
	Description string   // Optional description
	ParentID    string   // Parent taxonomy ID for hierarchical taxonomies
	EntityTypes []string // Entity types that can use this taxonomy
}

// TaxonomyTermOptions provides the options for creating a new taxonomy term
type TaxonomyTermOptions struct {
	TaxonomyID string // Parent taxonomy ID
	Name       string // Display name of the term
	Slug       string // URL-friendly identifier
	ParentID   string // Parent term ID for hierarchical terms
	SortOrder  int    // Sort order for the term
}

// == INTERFACE ==============================================================

// TaxonomyQueryInterface is a fluent taxonomy query. Every With* method
// marks the field as set and returns the builder for chaining. The store
// calls Validate() before executing and returns an error for invalid queries.
type TaxonomyQueryInterface interface {
	// Validate checks the set parameters for invalid values (empty strings
	// or slices where a value was set, invalid sort order).
	Validate() error

	// HasID reports whether the ID filter was set.
	HasID() bool
	// GetID returns the taxonomy ID filter.
	GetID() string
	// WithID filters by a specific taxonomy ID.
	WithID(id string) TaxonomyQueryInterface

	// HasIDs reports whether the IDs filter was set.
	HasIDs() bool
	// GetIDs returns the taxonomy IDs filter.
	GetIDs() []string
	// WithIDs filters by multiple taxonomy IDs.
	WithIDs(ids []string) TaxonomyQueryInterface

	// HasSlug reports whether the Slug filter was set.
	HasSlug() bool
	// GetSlug returns the taxonomy slug filter.
	GetSlug() string
	// WithSlug filters by taxonomy slug.
	WithSlug(slug string) TaxonomyQueryInterface

	// HasParentID reports whether the ParentID filter was set.
	HasParentID() bool
	// GetParentID returns the parent taxonomy ID filter.
	GetParentID() string
	// WithParentID filters by parent taxonomy ID.
	WithParentID(parentID string) TaxonomyQueryInterface

	// HasEntityTypes reports whether the EntityTypes filter was set.
	HasEntityTypes() bool
	// GetEntityTypes returns the entity types filter.
	GetEntityTypes() []string
	// WithEntityTypes filters by allowed entity types.
	WithEntityTypes(entityTypes []string) TaxonomyQueryInterface

	// HasLimit reports whether Limit was set.
	HasLimit() bool
	// GetLimit returns the result limit.
	GetLimit() uint64
	// WithLimit sets the maximum number of results.
	WithLimit(limit uint64) TaxonomyQueryInterface

	// HasOffset reports whether Offset was set.
	HasOffset() bool
	// GetOffset returns the result offset.
	GetOffset() uint64
	// WithOffset sets the number of results to skip.
	WithOffset(offset uint64) TaxonomyQueryInterface

	// HasSortBy reports whether SortBy was set.
	HasSortBy() bool
	// GetSortBy returns the sort column.
	GetSortBy() string
	// WithSortBy sets the column to sort by.
	WithSortBy(sortBy string) TaxonomyQueryInterface

	// HasSortOrder reports whether SortOrder was set.
	HasSortOrder() bool
	// GetSortOrder returns the sort direction.
	GetSortOrder() string
	// WithSortOrder sets the sort direction ("asc" or "desc").
	WithSortOrder(sortOrder string) TaxonomyQueryInterface

	// HasCountOnly reports whether CountOnly was set.
	HasCountOnly() bool
	// GetCountOnly returns whether only a count is requested.
	GetCountOnly() bool
	// WithCountOnly returns only a count, not results.
	WithCountOnly(countOnly bool) TaxonomyQueryInterface
}

// == CONSTRUCTOR ============================================================

// TaxonomyQuery returns a new fluent taxonomy query builder.
func TaxonomyQuery() TaxonomyQueryInterface {
	return &taxonomyQueryImplementation{}
}

// == TYPE ===================================================================

type taxonomyQueryImplementation struct {
	id             string
	hasID          bool
	ids            []string
	hasIDs         bool
	slug           string
	hasSlug        bool
	parentID       string
	hasParentID    bool
	entityTypes    []string
	hasEntityTypes bool
	limit          uint64
	hasLimit       bool
	offset         uint64
	hasOffset      bool
	sortBy         string
	hasSortBy      bool
	sortOrder      string
	hasSortOrder   bool
	countOnly      bool
	hasCountOnly   bool
}

// == INTERFACE VERIFICATION =================================================

var _ TaxonomyQueryInterface = (*taxonomyQueryImplementation)(nil)

// == INTERFACE IMPLEMENTATION ===============================================

// Validate checks the set parameters for invalid values.
func (q *taxonomyQueryImplementation) Validate() error {
	if q.hasID && q.id == "" {
		return errors.New("taxonomy query: id cannot be empty")
	}
	if q.hasIDs && len(q.ids) < 1 {
		return errors.New("taxonomy query: ids cannot be empty array")
	}
	if q.hasSlug && q.slug == "" {
		return errors.New("taxonomy query: slug cannot be empty")
	}
	if q.hasParentID && q.parentID == "" {
		return errors.New("taxonomy query: parent_id cannot be empty")
	}
	if q.hasEntityTypes && len(q.entityTypes) < 1 {
		return errors.New("taxonomy query: entity_types cannot be empty array")
	}
	if q.hasSortBy && q.sortBy == "" {
		return errors.New("taxonomy query: sort_by cannot be empty")
	}
	if q.hasSortOrder && q.sortOrder != "asc" && q.sortOrder != "desc" {
		return errors.New("taxonomy query: sort_order must be \"asc\" or \"desc\"")
	}
	return nil
}

func (q *taxonomyQueryImplementation) HasID() bool   { return q.hasID }
func (q *taxonomyQueryImplementation) GetID() string { return q.id }
func (q *taxonomyQueryImplementation) WithID(id string) TaxonomyQueryInterface {
	q.id, q.hasID = id, true
	return q
}

func (q *taxonomyQueryImplementation) HasIDs() bool     { return q.hasIDs }
func (q *taxonomyQueryImplementation) GetIDs() []string { return q.ids }
func (q *taxonomyQueryImplementation) WithIDs(ids []string) TaxonomyQueryInterface {
	q.ids, q.hasIDs = ids, true
	return q
}

func (q *taxonomyQueryImplementation) HasSlug() bool   { return q.hasSlug }
func (q *taxonomyQueryImplementation) GetSlug() string { return q.slug }
func (q *taxonomyQueryImplementation) WithSlug(slug string) TaxonomyQueryInterface {
	q.slug, q.hasSlug = slug, true
	return q
}

func (q *taxonomyQueryImplementation) HasParentID() bool   { return q.hasParentID }
func (q *taxonomyQueryImplementation) GetParentID() string { return q.parentID }
func (q *taxonomyQueryImplementation) WithParentID(parentID string) TaxonomyQueryInterface {
	q.parentID, q.hasParentID = parentID, true
	return q
}

func (q *taxonomyQueryImplementation) HasEntityTypes() bool     { return q.hasEntityTypes }
func (q *taxonomyQueryImplementation) GetEntityTypes() []string { return q.entityTypes }
func (q *taxonomyQueryImplementation) WithEntityTypes(entityTypes []string) TaxonomyQueryInterface {
	q.entityTypes, q.hasEntityTypes = entityTypes, true
	return q
}

func (q *taxonomyQueryImplementation) HasLimit() bool   { return q.hasLimit }
func (q *taxonomyQueryImplementation) GetLimit() uint64 { return q.limit }
func (q *taxonomyQueryImplementation) WithLimit(limit uint64) TaxonomyQueryInterface {
	q.limit, q.hasLimit = limit, true
	return q
}

func (q *taxonomyQueryImplementation) HasOffset() bool   { return q.hasOffset }
func (q *taxonomyQueryImplementation) GetOffset() uint64 { return q.offset }
func (q *taxonomyQueryImplementation) WithOffset(offset uint64) TaxonomyQueryInterface {
	q.offset, q.hasOffset = offset, true
	return q
}

func (q *taxonomyQueryImplementation) HasSortBy() bool   { return q.hasSortBy }
func (q *taxonomyQueryImplementation) GetSortBy() string { return q.sortBy }
func (q *taxonomyQueryImplementation) WithSortBy(sortBy string) TaxonomyQueryInterface {
	q.sortBy, q.hasSortBy = sortBy, true
	return q
}

func (q *taxonomyQueryImplementation) HasSortOrder() bool   { return q.hasSortOrder }
func (q *taxonomyQueryImplementation) GetSortOrder() string { return q.sortOrder }
func (q *taxonomyQueryImplementation) WithSortOrder(sortOrder string) TaxonomyQueryInterface {
	q.sortOrder, q.hasSortOrder = sortOrder, true
	return q
}

func (q *taxonomyQueryImplementation) HasCountOnly() bool { return q.hasCountOnly }
func (q *taxonomyQueryImplementation) GetCountOnly() bool { return q.countOnly }
func (q *taxonomyQueryImplementation) WithCountOnly(countOnly bool) TaxonomyQueryInterface {
	q.countOnly, q.hasCountOnly = countOnly, true
	return q
}
