package entitystore

import "errors"

// attribute_query.go defines AttributeQueryInterface — the fluent query type
// accepted by StoreInterface.AttributeList and AttributeCount — and its
// implementation. Unlike the cmsstore pageQuery pattern, the implementation
// uses typed fields with has* presence flags instead of a map[string]any, so
// every setter is compile-time type safe. Presence flags let the store read
// only the fields the caller explicitly set, and let Validate() distinguish
// "set to an invalid value" from "not set".
//
// The store reads the query through the Get* accessors and calls Validate()
// before executing, returning an error for invalid queries.
//
// Usage:
//
//	store.AttributeList(ctx, AttributeQuery().
//		WithEntityType("character").
//		WithAttributeKeys([]string{"name", "title"}).
//		WithLimit(100))

// == INTERFACE ==============================================================

// AttributeQueryInterface is a fluent attribute query. Every With* method
// marks the field as set and returns the builder for chaining. The store
// calls Validate() before executing and returns an error for invalid queries.
type AttributeQueryInterface interface {
	// Validate checks the set parameters for invalid values (empty strings
	// or slices where a value was set, invalid sort order).
	Validate() error

	// HasID reports whether the ID filter was set.
	HasID() bool
	// GetID returns the attribute ID filter.
	GetID() string
	// WithID filters by a specific attribute ID.
	WithID(id string) AttributeQueryInterface

	// HasIDs reports whether the IDs filter was set.
	HasIDs() bool
	// GetIDs returns the attribute IDs filter.
	GetIDs() []string
	// WithIDs filters by multiple attribute IDs.
	WithIDs(ids []string) AttributeQueryInterface

	// HasEntityID reports whether the EntityID filter was set.
	HasEntityID() bool
	// GetEntityID returns the entity ID filter.
	GetEntityID() string
	// WithEntityID filters by associated entity ID.
	WithEntityID(entityID string) AttributeQueryInterface

	// HasEntityType reports whether the EntityType filter was set.
	HasEntityType() bool
	// GetEntityType returns the entity type filter.
	GetEntityType() string
	// WithEntityType filters by entity type (joins the entities table).
	WithEntityType(entityType string) AttributeQueryInterface

	// HasEntityHandle reports whether the EntityHandle filter was set.
	HasEntityHandle() bool
	// GetEntityHandle returns the entity handle filter.
	GetEntityHandle() string
	// WithEntityHandle filters by entity handle (joins the entities table).
	WithEntityHandle(entityHandle string) AttributeQueryInterface

	// HasAttributeKey reports whether the AttributeKey filter was set.
	HasAttributeKey() bool
	// GetAttributeKey returns the attribute key filter.
	GetAttributeKey() string
	// WithAttributeKey filters by a single attribute key/name.
	WithAttributeKey(key string) AttributeQueryInterface

	// HasAttributeKeys reports whether the AttributeKeys filter was set.
	HasAttributeKeys() bool
	// GetAttributeKeys returns the attribute keys filter.
	GetAttributeKeys() []string
	// WithAttributeKeys filters by multiple attribute keys (WHERE IN).
	WithAttributeKeys(keys []string) AttributeQueryInterface

	// HasLimit reports whether Limit was set.
	HasLimit() bool
	// GetLimit returns the result limit.
	GetLimit() uint64
	// WithLimit sets the maximum number of results.
	WithLimit(limit uint64) AttributeQueryInterface

	// HasOffset reports whether Offset was set.
	HasOffset() bool
	// GetOffset returns the result offset.
	GetOffset() uint64
	// WithOffset sets the number of results to skip.
	WithOffset(offset uint64) AttributeQueryInterface

	// HasSortBy reports whether SortBy was set.
	HasSortBy() bool
	// GetSortBy returns the sort column.
	GetSortBy() string
	// WithSortBy sets the column to sort by.
	WithSortBy(sortBy string) AttributeQueryInterface

	// HasSortOrder reports whether SortOrder was set.
	HasSortOrder() bool
	// GetSortOrder returns the sort direction.
	GetSortOrder() string
	// WithSortOrder sets the sort direction ("asc" or "desc").
	WithSortOrder(sortOrder string) AttributeQueryInterface

	// HasCountOnly reports whether CountOnly was set.
	HasCountOnly() bool
	// GetCountOnly returns whether only a count is requested.
	GetCountOnly() bool
	// WithCountOnly returns only a count, not results.
	WithCountOnly(countOnly bool) AttributeQueryInterface
}

// == CONSTRUCTOR ============================================================

// AttributeQuery returns a new fluent attribute query builder.
func AttributeQuery() AttributeQueryInterface {
	return &attributeQueryImplementation{}
}

// == TYPE ===================================================================

type attributeQueryImplementation struct {
	id               string
	hasID            bool
	ids              []string
	hasIDs           bool
	entityID         string
	hasEntityID      bool
	entityType       string
	hasEntityType    bool
	entityHandle     string
	hasEntityHandle  bool
	attributeKey     string
	hasAttributeKey  bool
	attributeKeys    []string
	hasAttributeKeys bool
	limit            uint64
	hasLimit         bool
	offset           uint64
	hasOffset        bool
	sortBy           string
	hasSortBy        bool
	sortOrder        string
	hasSortOrder     bool
	countOnly        bool
	hasCountOnly     bool
}

// == INTERFACE VERIFICATION =================================================

var _ AttributeQueryInterface = (*attributeQueryImplementation)(nil)

// == INTERFACE IMPLEMENTATION ===============================================

// Validate checks the set parameters for invalid values.
func (q *attributeQueryImplementation) Validate() error {
	if q.hasID && q.id == "" {
		return errors.New("attribute query: id cannot be empty")
	}
	if q.hasIDs && len(q.ids) < 1 {
		return errors.New("attribute query: ids cannot be empty array")
	}
	if q.hasEntityID && q.entityID == "" {
		return errors.New("attribute query: entity_id cannot be empty")
	}
	if q.hasEntityType && q.entityType == "" {
		return errors.New("attribute query: entity_type cannot be empty")
	}
	if q.hasEntityHandle && q.entityHandle == "" {
		return errors.New("attribute query: entity_handle cannot be empty")
	}
	if q.hasAttributeKey && q.attributeKey == "" {
		return errors.New("attribute query: attribute_key cannot be empty")
	}
	if q.hasAttributeKeys && len(q.attributeKeys) < 1 {
		return errors.New("attribute query: attribute_keys cannot be empty array")
	}
	if q.hasSortBy && q.sortBy == "" {
		return errors.New("attribute query: sort_by cannot be empty")
	}
	if q.hasSortOrder && q.sortOrder != "asc" && q.sortOrder != "desc" {
		return errors.New("attribute query: sort_order must be \"asc\" or \"desc\"")
	}
	return nil
}

func (q *attributeQueryImplementation) HasID() bool   { return q.hasID }
func (q *attributeQueryImplementation) GetID() string { return q.id }
func (q *attributeQueryImplementation) WithID(id string) AttributeQueryInterface {
	q.id, q.hasID = id, true
	return q
}

func (q *attributeQueryImplementation) HasIDs() bool     { return q.hasIDs }
func (q *attributeQueryImplementation) GetIDs() []string { return q.ids }
func (q *attributeQueryImplementation) WithIDs(ids []string) AttributeQueryInterface {
	q.ids, q.hasIDs = ids, true
	return q
}

func (q *attributeQueryImplementation) HasEntityID() bool   { return q.hasEntityID }
func (q *attributeQueryImplementation) GetEntityID() string { return q.entityID }
func (q *attributeQueryImplementation) WithEntityID(entityID string) AttributeQueryInterface {
	q.entityID, q.hasEntityID = entityID, true
	return q
}

func (q *attributeQueryImplementation) HasEntityType() bool   { return q.hasEntityType }
func (q *attributeQueryImplementation) GetEntityType() string { return q.entityType }
func (q *attributeQueryImplementation) WithEntityType(entityType string) AttributeQueryInterface {
	q.entityType, q.hasEntityType = entityType, true
	return q
}

func (q *attributeQueryImplementation) HasEntityHandle() bool   { return q.hasEntityHandle }
func (q *attributeQueryImplementation) GetEntityHandle() string { return q.entityHandle }
func (q *attributeQueryImplementation) WithEntityHandle(entityHandle string) AttributeQueryInterface {
	q.entityHandle, q.hasEntityHandle = entityHandle, true
	return q
}

func (q *attributeQueryImplementation) HasAttributeKey() bool   { return q.hasAttributeKey }
func (q *attributeQueryImplementation) GetAttributeKey() string { return q.attributeKey }
func (q *attributeQueryImplementation) WithAttributeKey(key string) AttributeQueryInterface {
	q.attributeKey, q.hasAttributeKey = key, true
	return q
}

func (q *attributeQueryImplementation) HasAttributeKeys() bool     { return q.hasAttributeKeys }
func (q *attributeQueryImplementation) GetAttributeKeys() []string { return q.attributeKeys }
func (q *attributeQueryImplementation) WithAttributeKeys(keys []string) AttributeQueryInterface {
	q.attributeKeys, q.hasAttributeKeys = keys, true
	return q
}

func (q *attributeQueryImplementation) HasLimit() bool   { return q.hasLimit }
func (q *attributeQueryImplementation) GetLimit() uint64 { return q.limit }
func (q *attributeQueryImplementation) WithLimit(limit uint64) AttributeQueryInterface {
	q.limit, q.hasLimit = limit, true
	return q
}

func (q *attributeQueryImplementation) HasOffset() bool   { return q.hasOffset }
func (q *attributeQueryImplementation) GetOffset() uint64 { return q.offset }
func (q *attributeQueryImplementation) WithOffset(offset uint64) AttributeQueryInterface {
	q.offset, q.hasOffset = offset, true
	return q
}

func (q *attributeQueryImplementation) HasSortBy() bool   { return q.hasSortBy }
func (q *attributeQueryImplementation) GetSortBy() string { return q.sortBy }
func (q *attributeQueryImplementation) WithSortBy(sortBy string) AttributeQueryInterface {
	q.sortBy, q.hasSortBy = sortBy, true
	return q
}

func (q *attributeQueryImplementation) HasSortOrder() bool   { return q.hasSortOrder }
func (q *attributeQueryImplementation) GetSortOrder() string { return q.sortOrder }
func (q *attributeQueryImplementation) WithSortOrder(sortOrder string) AttributeQueryInterface {
	q.sortOrder, q.hasSortOrder = sortOrder, true
	return q
}

func (q *attributeQueryImplementation) HasCountOnly() bool { return q.hasCountOnly }
func (q *attributeQueryImplementation) GetCountOnly() bool { return q.countOnly }
func (q *attributeQueryImplementation) WithCountOnly(countOnly bool) AttributeQueryInterface {
	q.countOnly, q.hasCountOnly = countOnly, true
	return q
}
