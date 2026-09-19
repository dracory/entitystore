package entitystore

import "errors"

// taxonomy_term_query.go defines TaxonomyTermQueryInterface — the fluent
// query type accepted by StoreInterface.TaxonomyTermList, TaxonomyTermCount
// and TaxonomyTermTrashList — and its implementation. The implementation
// uses typed fields with has* presence flags instead of a map[string]any,
// so every setter is compile-time type safe. Presence flags let the store
// read only the fields the caller explicitly set, and let Validate()
// distinguish "set to an invalid value" from "not set".
//
// The store reads the query through the Get* accessors and calls Validate()
// before executing, returning an error for invalid queries.
//
// Usage:
//
//	store.TaxonomyTermList(ctx, TaxonomyTermQuery().
//		WithTaxonomyID("taxonomy-1").
//		WithSlug("electronics").
//		WithLimit(100))

// == INTERFACE ==============================================================

// TaxonomyTermQueryInterface is a fluent taxonomy term query. Every With*
// method marks the field as set and returns the builder for chaining. The
// store calls Validate() before executing and returns an error for invalid
// queries.
type TaxonomyTermQueryInterface interface {
	// Validate checks the set parameters for invalid values (empty strings
	// or slices where a value was set, invalid sort order).
	Validate() error

	// HasID reports whether the ID filter was set.
	HasID() bool
	// GetID returns the term ID filter.
	GetID() string
	// WithID filters by a specific term ID.
	WithID(id string) TaxonomyTermQueryInterface

	// HasIDs reports whether the IDs filter was set.
	HasIDs() bool
	// GetIDs returns the term IDs filter.
	GetIDs() []string
	// WithIDs filters by multiple term IDs.
	WithIDs(ids []string) TaxonomyTermQueryInterface

	// HasTaxonomyID reports whether the TaxonomyID filter was set.
	HasTaxonomyID() bool
	// GetTaxonomyID returns the taxonomy ID filter.
	GetTaxonomyID() string
	// WithTaxonomyID filters by parent taxonomy ID.
	WithTaxonomyID(taxonomyID string) TaxonomyTermQueryInterface

	// HasSlug reports whether the Slug filter was set.
	HasSlug() bool
	// GetSlug returns the term slug filter.
	GetSlug() string
	// WithSlug filters by term slug.
	WithSlug(slug string) TaxonomyTermQueryInterface

	// HasParentID reports whether the ParentID filter was set.
	HasParentID() bool
	// GetParentID returns the parent term ID filter.
	GetParentID() string
	// WithParentID filters by parent term ID.
	WithParentID(parentID string) TaxonomyTermQueryInterface

	// HasLimit reports whether Limit was set.
	HasLimit() bool
	// GetLimit returns the result limit.
	GetLimit() uint64
	// WithLimit sets the maximum number of results.
	WithLimit(limit uint64) TaxonomyTermQueryInterface

	// HasOffset reports whether Offset was set.
	HasOffset() bool
	// GetOffset returns the result offset.
	GetOffset() uint64
	// WithOffset sets the number of results to skip.
	WithOffset(offset uint64) TaxonomyTermQueryInterface

	// HasSortBy reports whether SortBy was set.
	HasSortBy() bool
	// GetSortBy returns the sort column.
	GetSortBy() string
	// WithSortBy sets the column to sort by.
	WithSortBy(sortBy string) TaxonomyTermQueryInterface

	// HasSortOrder reports whether SortOrder was set.
	HasSortOrder() bool
	// GetSortOrder returns the sort direction.
	GetSortOrder() string
	// WithSortOrder sets the sort direction ("asc" or "desc").
	WithSortOrder(sortOrder string) TaxonomyTermQueryInterface

	// HasCountOnly reports whether CountOnly was set.
	HasCountOnly() bool
	// GetCountOnly returns whether only a count is requested.
	GetCountOnly() bool
	// WithCountOnly returns only a count, not results.
	WithCountOnly(countOnly bool) TaxonomyTermQueryInterface

	// HasCreatedAtGte reports whether the created_at lower bound was set.
	HasCreatedAtGte() bool
	// GetCreatedAtGte returns the created_at lower bound (inclusive).
	GetCreatedAtGte() string
	// WithCreatedAtGte filters to rows created at or after the given UTC
	// datetime ("YYYY-MM-DD HH:MM:SS").
	WithCreatedAtGte(createdAtGte string) TaxonomyTermQueryInterface

	// HasCreatedAtLte reports whether the created_at upper bound was set.
	HasCreatedAtLte() bool
	// GetCreatedAtLte returns the created_at upper bound (inclusive).
	GetCreatedAtLte() string
	// WithCreatedAtLte filters to rows created at or before the given UTC
	// datetime ("YYYY-MM-DD HH:MM:SS").
	WithCreatedAtLte(createdAtLte string) TaxonomyTermQueryInterface

	// HasUpdatedAtGte reports whether the updated_at lower bound was set.
	HasUpdatedAtGte() bool
	// GetUpdatedAtGte returns the updated_at lower bound (inclusive).
	GetUpdatedAtGte() string
	// WithUpdatedAtGte filters to rows updated at or after the given UTC
	// datetime ("YYYY-MM-DD HH:MM:SS").
	WithUpdatedAtGte(updatedAtGte string) TaxonomyTermQueryInterface

	// HasUpdatedAtLte reports whether the updated_at upper bound was set.
	HasUpdatedAtLte() bool
	// GetUpdatedAtLte returns the updated_at upper bound (inclusive).
	GetUpdatedAtLte() string
	// WithUpdatedAtLte filters to rows updated at or before the given UTC
	// datetime ("YYYY-MM-DD HH:MM:SS").
	WithUpdatedAtLte(updatedAtLte string) TaxonomyTermQueryInterface
}

// == CONSTRUCTOR ============================================================

// TaxonomyTermQuery returns a new fluent taxonomy term query builder.
func TaxonomyTermQuery() TaxonomyTermQueryInterface {
	return &taxonomyTermQueryImplementation{}
}

// == TYPE ===================================================================

type taxonomyTermQueryImplementation struct {
	id              string
	hasID           bool
	ids             []string
	hasIDs          bool
	taxonomyID      string
	hasTaxonomyID   bool
	slug            string
	hasSlug         bool
	parentID        string
	hasParentID     bool
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
}

// == INTERFACE VERIFICATION =================================================

var _ TaxonomyTermQueryInterface = (*taxonomyTermQueryImplementation)(nil)

// == INTERFACE IMPLEMENTATION ===============================================

// Validate checks the set parameters for invalid values.
func (q *taxonomyTermQueryImplementation) Validate() error {
	if q.hasID && q.id == "" {
		return errors.New("taxonomy term query: id cannot be empty")
	}
	if q.hasIDs && len(q.ids) < 1 {
		return errors.New("taxonomy term query: ids cannot be empty array")
	}
	if q.hasTaxonomyID && q.taxonomyID == "" {
		return errors.New("taxonomy term query: taxonomy_id cannot be empty")
	}
	if q.hasSlug && q.slug == "" {
		return errors.New("taxonomy term query: slug cannot be empty")
	}
	if q.hasParentID && q.parentID == "" {
		return errors.New("taxonomy term query: parent_id cannot be empty")
	}
	if q.hasSortBy && q.sortBy == "" {
		return errors.New("taxonomy term query: sort_by cannot be empty")
	}
	if q.hasSortOrder && q.sortOrder != "asc" && q.sortOrder != "desc" {
		return errors.New("taxonomy term query: sort_order must be \"asc\" or \"desc\"")
	}
	if q.hasCreatedAtGte && q.createdAtGte == "" {
		return errors.New("taxonomy term query: created_at_gte cannot be empty")
	}
	if q.hasCreatedAtLte && q.createdAtLte == "" {
		return errors.New("taxonomy term query: created_at_lte cannot be empty")
	}
	if q.hasUpdatedAtGte && q.updatedAtGte == "" {
		return errors.New("taxonomy term query: updated_at_gte cannot be empty")
	}
	if q.hasUpdatedAtLte && q.updatedAtLte == "" {
		return errors.New("taxonomy term query: updated_at_lte cannot be empty")
	}
	return nil
}

func (q *taxonomyTermQueryImplementation) HasID() bool   { return q.hasID }
func (q *taxonomyTermQueryImplementation) GetID() string { return q.id }
func (q *taxonomyTermQueryImplementation) WithID(id string) TaxonomyTermQueryInterface {
	q.id, q.hasID = id, true
	return q
}

func (q *taxonomyTermQueryImplementation) HasIDs() bool     { return q.hasIDs }
func (q *taxonomyTermQueryImplementation) GetIDs() []string { return q.ids }
func (q *taxonomyTermQueryImplementation) WithIDs(ids []string) TaxonomyTermQueryInterface {
	q.ids, q.hasIDs = ids, true
	return q
}

func (q *taxonomyTermQueryImplementation) HasTaxonomyID() bool   { return q.hasTaxonomyID }
func (q *taxonomyTermQueryImplementation) GetTaxonomyID() string { return q.taxonomyID }
func (q *taxonomyTermQueryImplementation) WithTaxonomyID(taxonomyID string) TaxonomyTermQueryInterface {
	q.taxonomyID, q.hasTaxonomyID = taxonomyID, true
	return q
}

func (q *taxonomyTermQueryImplementation) HasSlug() bool   { return q.hasSlug }
func (q *taxonomyTermQueryImplementation) GetSlug() string { return q.slug }
func (q *taxonomyTermQueryImplementation) WithSlug(slug string) TaxonomyTermQueryInterface {
	q.slug, q.hasSlug = slug, true
	return q
}

func (q *taxonomyTermQueryImplementation) HasParentID() bool   { return q.hasParentID }
func (q *taxonomyTermQueryImplementation) GetParentID() string { return q.parentID }
func (q *taxonomyTermQueryImplementation) WithParentID(parentID string) TaxonomyTermQueryInterface {
	q.parentID, q.hasParentID = parentID, true
	return q
}

func (q *taxonomyTermQueryImplementation) HasLimit() bool   { return q.hasLimit }
func (q *taxonomyTermQueryImplementation) GetLimit() uint64 { return q.limit }
func (q *taxonomyTermQueryImplementation) WithLimit(limit uint64) TaxonomyTermQueryInterface {
	q.limit, q.hasLimit = limit, true
	return q
}

func (q *taxonomyTermQueryImplementation) HasOffset() bool   { return q.hasOffset }
func (q *taxonomyTermQueryImplementation) GetOffset() uint64 { return q.offset }
func (q *taxonomyTermQueryImplementation) WithOffset(offset uint64) TaxonomyTermQueryInterface {
	q.offset, q.hasOffset = offset, true
	return q
}

func (q *taxonomyTermQueryImplementation) HasSortBy() bool   { return q.hasSortBy }
func (q *taxonomyTermQueryImplementation) GetSortBy() string { return q.sortBy }
func (q *taxonomyTermQueryImplementation) WithSortBy(sortBy string) TaxonomyTermQueryInterface {
	q.sortBy, q.hasSortBy = sortBy, true
	return q
}

func (q *taxonomyTermQueryImplementation) HasSortOrder() bool   { return q.hasSortOrder }
func (q *taxonomyTermQueryImplementation) GetSortOrder() string { return q.sortOrder }
func (q *taxonomyTermQueryImplementation) WithSortOrder(sortOrder string) TaxonomyTermQueryInterface {
	q.sortOrder, q.hasSortOrder = sortOrder, true
	return q
}

func (q *taxonomyTermQueryImplementation) HasCountOnly() bool { return q.hasCountOnly }
func (q *taxonomyTermQueryImplementation) GetCountOnly() bool { return q.countOnly }
func (q *taxonomyTermQueryImplementation) WithCountOnly(countOnly bool) TaxonomyTermQueryInterface {
	q.countOnly, q.hasCountOnly = countOnly, true
	return q
}

func (q *taxonomyTermQueryImplementation) HasCreatedAtGte() bool   { return q.hasCreatedAtGte }
func (q *taxonomyTermQueryImplementation) GetCreatedAtGte() string { return q.createdAtGte }
func (q *taxonomyTermQueryImplementation) WithCreatedAtGte(createdAtGte string) TaxonomyTermQueryInterface {
	q.createdAtGte, q.hasCreatedAtGte = createdAtGte, true
	return q
}

func (q *taxonomyTermQueryImplementation) HasCreatedAtLte() bool   { return q.hasCreatedAtLte }
func (q *taxonomyTermQueryImplementation) GetCreatedAtLte() string { return q.createdAtLte }
func (q *taxonomyTermQueryImplementation) WithCreatedAtLte(createdAtLte string) TaxonomyTermQueryInterface {
	q.createdAtLte, q.hasCreatedAtLte = createdAtLte, true
	return q
}

func (q *taxonomyTermQueryImplementation) HasUpdatedAtGte() bool   { return q.hasUpdatedAtGte }
func (q *taxonomyTermQueryImplementation) GetUpdatedAtGte() string { return q.updatedAtGte }
func (q *taxonomyTermQueryImplementation) WithUpdatedAtGte(updatedAtGte string) TaxonomyTermQueryInterface {
	q.updatedAtGte, q.hasUpdatedAtGte = updatedAtGte, true
	return q
}

func (q *taxonomyTermQueryImplementation) HasUpdatedAtLte() bool   { return q.hasUpdatedAtLte }
func (q *taxonomyTermQueryImplementation) GetUpdatedAtLte() string { return q.updatedAtLte }
func (q *taxonomyTermQueryImplementation) WithUpdatedAtLte(updatedAtLte string) TaxonomyTermQueryInterface {
	q.updatedAtLte, q.hasUpdatedAtLte = updatedAtLte, true
	return q
}
