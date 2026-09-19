package entitystore

import "errors"

// relationship_query.go defines RelationshipQueryInterface — the fluent
// query type accepted by StoreInterface.RelationshipList, RelationshipCount
// and RelationshipTrashList — and its implementation. The implementation
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
//	store.RelationshipList(ctx, RelationshipQuery().
//		WithEntityID("entity-1").
//		WithRelationshipType("belongs_to").
//		WithLimit(100))

// == INTERFACE ==============================================================

// RelationshipQueryInterface is a fluent relationship query. Every With*
// method marks the field as set and returns the builder for chaining. The
// store calls Validate() before executing and returns an error for invalid
// queries.
type RelationshipQueryInterface interface {
	// Validate checks the set parameters for invalid values (empty strings
	// or slices where a value was set, invalid sort order).
	Validate() error

	// HasID reports whether the ID filter was set.
	HasID() bool
	// GetID returns the relationship ID filter.
	GetID() string
	// WithID filters by a specific relationship ID.
	WithID(id string) RelationshipQueryInterface

	// HasIDs reports whether the IDs filter was set.
	HasIDs() bool
	// GetIDs returns the relationship IDs filter.
	GetIDs() []string
	// WithIDs filters by multiple relationship IDs.
	WithIDs(ids []string) RelationshipQueryInterface

	// HasEntityID reports whether the EntityID filter was set.
	HasEntityID() bool
	// GetEntityID returns the source entity ID filter.
	GetEntityID() string
	// WithEntityID filters by source entity ID.
	WithEntityID(entityID string) RelationshipQueryInterface

	// HasEntityIDs reports whether the EntityIDs filter was set.
	HasEntityIDs() bool
	// GetEntityIDs returns the source entity IDs filter.
	GetEntityIDs() []string
	// WithEntityIDs filters by multiple source entity IDs.
	WithEntityIDs(entityIDs []string) RelationshipQueryInterface

	// HasRelatedEntityID reports whether the RelatedEntityID filter was set.
	HasRelatedEntityID() bool
	// GetRelatedEntityID returns the target entity ID filter.
	GetRelatedEntityID() string
	// WithRelatedEntityID filters by target entity ID.
	WithRelatedEntityID(relatedEntityID string) RelationshipQueryInterface

	// HasRelatedEntityIDs reports whether the RelatedEntityIDs filter was set.
	HasRelatedEntityIDs() bool
	// GetRelatedEntityIDs returns the target entity IDs filter.
	GetRelatedEntityIDs() []string
	// WithRelatedEntityIDs filters by multiple target entity IDs.
	WithRelatedEntityIDs(relatedEntityIDs []string) RelationshipQueryInterface

	// HasRelationshipType reports whether the RelationshipType filter was set.
	HasRelationshipType() bool
	// GetRelationshipType returns the relationship type filter.
	GetRelationshipType() string
	// WithRelationshipType filters by relationship type.
	WithRelationshipType(relationshipType string) RelationshipQueryInterface

	// HasParentID reports whether the ParentID filter was set.
	HasParentID() bool
	// GetParentID returns the parent relationship ID filter.
	GetParentID() string
	// WithParentID filters by parent relationship ID.
	WithParentID(parentID string) RelationshipQueryInterface

	// HasLimit reports whether Limit was set.
	HasLimit() bool
	// GetLimit returns the result limit.
	GetLimit() uint64
	// WithLimit sets the maximum number of results.
	WithLimit(limit uint64) RelationshipQueryInterface

	// HasOffset reports whether Offset was set.
	HasOffset() bool
	// GetOffset returns the result offset.
	GetOffset() uint64
	// WithOffset sets the number of results to skip.
	WithOffset(offset uint64) RelationshipQueryInterface

	// HasSortBy reports whether SortBy was set.
	HasSortBy() bool
	// GetSortBy returns the sort column.
	GetSortBy() string
	// WithSortBy sets the column to sort by.
	WithSortBy(sortBy string) RelationshipQueryInterface

	// HasSortOrder reports whether SortOrder was set.
	HasSortOrder() bool
	// GetSortOrder returns the sort direction.
	GetSortOrder() string
	// WithSortOrder sets the sort direction ("asc" or "desc").
	WithSortOrder(sortOrder string) RelationshipQueryInterface

	// HasCountOnly reports whether CountOnly was set.
	HasCountOnly() bool
	// GetCountOnly returns whether only a count is requested.
	GetCountOnly() bool
	// WithCountOnly returns only a count, not results.
	WithCountOnly(countOnly bool) RelationshipQueryInterface
}

// == CONSTRUCTOR ============================================================

// RelationshipQuery returns a new fluent relationship query builder.
func RelationshipQuery() RelationshipQueryInterface {
	return &relationshipQueryImplementation{}
}

// == TYPE ===================================================================

type relationshipQueryImplementation struct {
	id                  string
	hasID               bool
	ids                 []string
	hasIDs              bool
	entityID            string
	hasEntityID         bool
	entityIDs           []string
	hasEntityIDs        bool
	relatedEntityID     string
	hasRelatedEntityID  bool
	relatedEntityIDs    []string
	hasRelatedEntityIDs bool
	relationshipType    string
	hasRelationshipType bool
	parentID            string
	hasParentID         bool
	limit               uint64
	hasLimit            bool
	offset              uint64
	hasOffset           bool
	sortBy              string
	hasSortBy           bool
	sortOrder           string
	hasSortOrder        bool
	countOnly           bool
	hasCountOnly        bool
}

// == INTERFACE VERIFICATION =================================================

var _ RelationshipQueryInterface = (*relationshipQueryImplementation)(nil)

// == INTERFACE IMPLEMENTATION ===============================================

// Validate checks the set parameters for invalid values.
func (q *relationshipQueryImplementation) Validate() error {
	if q.hasID && q.id == "" {
		return errors.New("relationship query: id cannot be empty")
	}
	if q.hasIDs && len(q.ids) < 1 {
		return errors.New("relationship query: ids cannot be empty array")
	}
	if q.hasEntityID && q.entityID == "" {
		return errors.New("relationship query: entity_id cannot be empty")
	}
	if q.hasEntityIDs && len(q.entityIDs) < 1 {
		return errors.New("relationship query: entity_ids cannot be empty array")
	}
	if q.hasRelatedEntityID && q.relatedEntityID == "" {
		return errors.New("relationship query: related_entity_id cannot be empty")
	}
	if q.hasRelatedEntityIDs && len(q.relatedEntityIDs) < 1 {
		return errors.New("relationship query: related_entity_ids cannot be empty array")
	}
	if q.hasRelationshipType && q.relationshipType == "" {
		return errors.New("relationship query: relationship_type cannot be empty")
	}
	if q.hasParentID && q.parentID == "" {
		return errors.New("relationship query: parent_id cannot be empty")
	}
	if q.hasSortBy && q.sortBy == "" {
		return errors.New("relationship query: sort_by cannot be empty")
	}
	if q.hasSortOrder && q.sortOrder != "asc" && q.sortOrder != "desc" {
		return errors.New("relationship query: sort_order must be \"asc\" or \"desc\"")
	}
	return nil
}

func (q *relationshipQueryImplementation) HasID() bool   { return q.hasID }
func (q *relationshipQueryImplementation) GetID() string { return q.id }
func (q *relationshipQueryImplementation) WithID(id string) RelationshipQueryInterface {
	q.id, q.hasID = id, true
	return q
}

func (q *relationshipQueryImplementation) HasIDs() bool     { return q.hasIDs }
func (q *relationshipQueryImplementation) GetIDs() []string { return q.ids }
func (q *relationshipQueryImplementation) WithIDs(ids []string) RelationshipQueryInterface {
	q.ids, q.hasIDs = ids, true
	return q
}

func (q *relationshipQueryImplementation) HasEntityID() bool   { return q.hasEntityID }
func (q *relationshipQueryImplementation) GetEntityID() string { return q.entityID }
func (q *relationshipQueryImplementation) WithEntityID(entityID string) RelationshipQueryInterface {
	q.entityID, q.hasEntityID = entityID, true
	return q
}

func (q *relationshipQueryImplementation) HasEntityIDs() bool     { return q.hasEntityIDs }
func (q *relationshipQueryImplementation) GetEntityIDs() []string { return q.entityIDs }
func (q *relationshipQueryImplementation) WithEntityIDs(entityIDs []string) RelationshipQueryInterface {
	q.entityIDs, q.hasEntityIDs = entityIDs, true
	return q
}

func (q *relationshipQueryImplementation) HasRelatedEntityID() bool   { return q.hasRelatedEntityID }
func (q *relationshipQueryImplementation) GetRelatedEntityID() string { return q.relatedEntityID }
func (q *relationshipQueryImplementation) WithRelatedEntityID(relatedEntityID string) RelationshipQueryInterface {
	q.relatedEntityID, q.hasRelatedEntityID = relatedEntityID, true
	return q
}

func (q *relationshipQueryImplementation) HasRelatedEntityIDs() bool {
	return q.hasRelatedEntityIDs
}
func (q *relationshipQueryImplementation) GetRelatedEntityIDs() []string {
	return q.relatedEntityIDs
}
func (q *relationshipQueryImplementation) WithRelatedEntityIDs(relatedEntityIDs []string) RelationshipQueryInterface {
	q.relatedEntityIDs, q.hasRelatedEntityIDs = relatedEntityIDs, true
	return q
}

func (q *relationshipQueryImplementation) HasRelationshipType() bool {
	return q.hasRelationshipType
}
func (q *relationshipQueryImplementation) GetRelationshipType() string {
	return q.relationshipType
}
func (q *relationshipQueryImplementation) WithRelationshipType(relationshipType string) RelationshipQueryInterface {
	q.relationshipType, q.hasRelationshipType = relationshipType, true
	return q
}

func (q *relationshipQueryImplementation) HasParentID() bool   { return q.hasParentID }
func (q *relationshipQueryImplementation) GetParentID() string { return q.parentID }
func (q *relationshipQueryImplementation) WithParentID(parentID string) RelationshipQueryInterface {
	q.parentID, q.hasParentID = parentID, true
	return q
}

func (q *relationshipQueryImplementation) HasLimit() bool   { return q.hasLimit }
func (q *relationshipQueryImplementation) GetLimit() uint64 { return q.limit }
func (q *relationshipQueryImplementation) WithLimit(limit uint64) RelationshipQueryInterface {
	q.limit, q.hasLimit = limit, true
	return q
}

func (q *relationshipQueryImplementation) HasOffset() bool   { return q.hasOffset }
func (q *relationshipQueryImplementation) GetOffset() uint64 { return q.offset }
func (q *relationshipQueryImplementation) WithOffset(offset uint64) RelationshipQueryInterface {
	q.offset, q.hasOffset = offset, true
	return q
}

func (q *relationshipQueryImplementation) HasSortBy() bool   { return q.hasSortBy }
func (q *relationshipQueryImplementation) GetSortBy() string { return q.sortBy }
func (q *relationshipQueryImplementation) WithSortBy(sortBy string) RelationshipQueryInterface {
	q.sortBy, q.hasSortBy = sortBy, true
	return q
}

func (q *relationshipQueryImplementation) HasSortOrder() bool   { return q.hasSortOrder }
func (q *relationshipQueryImplementation) GetSortOrder() string { return q.sortOrder }
func (q *relationshipQueryImplementation) WithSortOrder(sortOrder string) RelationshipQueryInterface {
	q.sortOrder, q.hasSortOrder = sortOrder, true
	return q
}

func (q *relationshipQueryImplementation) HasCountOnly() bool { return q.hasCountOnly }
func (q *relationshipQueryImplementation) GetCountOnly() bool { return q.countOnly }
func (q *relationshipQueryImplementation) WithCountOnly(countOnly bool) RelationshipQueryInterface {
	q.countOnly, q.hasCountOnly = countOnly, true
	return q
}
