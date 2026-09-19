package entitystore

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
)

// relationshipRow is used for scanning relationship query results
type relationshipRow struct {
	ID               string `db:"id"`
	EntityID         string `db:"entity_id"`
	RelatedEntityID  string `db:"related_entity_id"`
	RelationshipType string `db:"relationship_type"`
	ParentID         string `db:"parent_id"`
	Sequence         int    `db:"sequence"`
	Metadata         string `db:"metadata"`
	CreatedAt        string `db:"created_at"`
}

// RelationshipCreate persists a new relationship record
func (st *storeImplementation) RelationshipCreate(ctx context.Context, relationship RelationshipInterface) error {
	if relationship == nil {
		return errors.New("relationship cannot be nil")
	}

	if relationship.GetEntityID() == "" {
		return errors.New("entity_id is required")
	}
	if relationship.GetRelatedEntityID() == "" {
		return errors.New("related_entity_id is required")
	}
	if relationship.GetRelationshipType() == "" {
		return errors.New("relationship_type is required")
	}

	if relationship.GetEntityID() == relationship.GetRelatedEntityID() {
		if relationship.GetRelationshipType() == RELATIONSHIP_TYPE_BELONGS_TO || relationship.GetRelationshipType() == RELATIONSHIP_TYPE_HAS_MANY {
			return errors.New("self-referencing relationships not allowed for belongs_to and has_many types")
		}
	}

	if relationship.ID() == "" {
		relationship.SetID(GenerateShortID())
	}

	if err := validateIDLength(relationship.ID(), COLUMN_ID); err != nil {
		return err
	}
	if err := validateIDLength(relationship.GetEntityID(), COLUMN_ENTITY_ID); err != nil {
		return err
	}
	if err := validateIDLength(relationship.GetRelatedEntityID(), COLUMN_RELATED_ENTITY_ID); err != nil {
		return err
	}
	if err := validateIDLength(relationship.GetParentID(), COLUMN_PARENT_ID); err != nil {
		return err
	}

	if relationship.GetCreatedAt() == "" {
		relationship.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	row := map[string]any{}
	for k, v := range relationship.Data() {
		row[k] = v
	}

	if st.GetDebug() {
		log.Println("RelationshipCreate:", row)
	}

	return st.db.Query().Table(st.relationshipTableName).Create(row)
}

// RelationshipCreateByOptions creates a relationship using the provided options
func (st *storeImplementation) RelationshipCreateByOptions(ctx context.Context, opts RelationshipOptions) (RelationshipInterface, error) {
	existing, err := st.RelationshipFindByEntities(ctx, opts.EntityID, opts.RelatedEntityID, opts.RelationshipType)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("relationship already exists")
	}

	rel := NewRelationship()
	rel.SetEntityID(opts.EntityID)
	rel.SetRelatedEntityID(opts.RelatedEntityID)
	rel.SetRelationshipType(opts.RelationshipType)
	rel.SetParentID(opts.ParentID)
	rel.SetSequence(opts.Sequence)
	rel.SetMetadata(opts.Metadata)

	if err := st.RelationshipCreate(ctx, rel); err != nil {
		return nil, err
	}

	return rel, nil
}

// RelationshipDelete removes a relationship record by ID
func (st *storeImplementation) RelationshipDelete(ctx context.Context, relationshipID string) (bool, error) {
	result, err := st.db.Query().Table(st.relationshipTableName).Where(COLUMN_ID+" = ?", relationshipID).Delete()
	if err != nil {
		return false, err
	}

	return result.RowsAffected > 0, nil
}

// RelationshipDeleteAll removes all relationships for an entity (both as source and target)
func (st *storeImplementation) RelationshipDeleteAll(ctx context.Context, entityID string) error {
	_, err := st.db.Query().Table(st.relationshipTableName).Where(COLUMN_ENTITY_ID+" = ?", entityID).Delete()
	if err != nil {
		return err
	}

	_, err = st.db.Query().Table(st.relationshipTableName).Where(COLUMN_RELATED_ENTITY_ID+" = ?", entityID).Delete()
	return err
}

// RelationshipFind finds a relationship by its ID
func (st *storeImplementation) RelationshipFind(ctx context.Context, relationshipID string) (RelationshipInterface, error) {
	if relationshipID == "" {
		return nil, errors.New("relationship ID cannot be empty")
	}

	list, err := st.RelationshipList(ctx, RelationshipQuery().
		WithID(relationshipID).
		WithLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// RelationshipFindByEntities finds a relationship by entity IDs and type
func (st *storeImplementation) RelationshipFindByEntities(ctx context.Context, entityID string, relatedEntityID string, relationshipType string) (RelationshipInterface, error) {
	if entityID == "" || relatedEntityID == "" || relationshipType == "" {
		return nil, errors.New("entityID, relatedEntityID, and relationshipType are required")
	}

	list, err := st.RelationshipList(ctx, RelationshipQuery().
		WithEntityID(entityID).
		WithRelatedEntityID(relatedEntityID).
		WithRelationshipType(relationshipType).
		WithLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// applyRelationshipFilters applies the common filter clauses from a
// validated relationship query to a query. Shared by RelationshipList,
// RelationshipCount and RelationshipTrashList to prevent filter drift.
func (st *storeImplementation) applyRelationshipFilters(q orm.Query, query RelationshipQueryInterface) orm.Query {
	if query.GetID() != "" {
		q = q.Where(COLUMN_ID+" = ?", query.GetID())
	}

	if len(query.GetIDs()) > 0 {
		ids := make([]any, len(query.GetIDs()))
		for i, id := range query.GetIDs() {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_ID, ids)
	}

	if query.GetEntityID() != "" {
		q = q.Where(COLUMN_ENTITY_ID+" = ?", query.GetEntityID())
	}

	if len(query.GetEntityIDs()) > 0 {
		ids := make([]any, len(query.GetEntityIDs()))
		for i, id := range query.GetEntityIDs() {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_ENTITY_ID, ids)
	}

	if query.GetRelatedEntityID() != "" {
		q = q.Where(COLUMN_RELATED_ENTITY_ID+" = ?", query.GetRelatedEntityID())
	}

	if len(query.GetRelatedEntityIDs()) > 0 {
		ids := make([]any, len(query.GetRelatedEntityIDs()))
		for i, id := range query.GetRelatedEntityIDs() {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_RELATED_ENTITY_ID, ids)
	}

	if query.GetRelationshipType() != "" {
		q = q.Where(COLUMN_RELATIONSHIP_TYPE+" = ?", query.GetRelationshipType())
	}

	if query.GetParentID() != "" {
		q = q.Where(COLUMN_PARENT_ID+" = ?", query.GetParentID())
	}

	q = applyTimeRange(q, COLUMN_CREATED_AT, query.GetCreatedAtGte(), query.GetCreatedAtLte())

	return q
}

// RelationshipList lists relationships matching the given fluent query
func (st *storeImplementation) RelationshipList(ctx context.Context, query RelationshipQueryInterface) ([]RelationshipInterface, error) {
	if query == nil {
		return nil, errors.New("relationship query cannot be nil")
	}

	if err := query.Validate(); err != nil {
		return nil, err
	}

	q := st.applyRelationshipFilters(st.db.Query().Table(st.relationshipTableName), query)

	sortByColumn := COLUMN_CREATED_AT
	sortOrder := "asc"

	if query.GetSortOrder() != "" {
		sortOrder = query.GetSortOrder()
	}

	if query.GetSortBy() != "" {
		sortByColumn = query.GetSortBy()
	}

	q = q.OrderBy(sortByColumn, sortOrder)

	if query.GetOffset() > 0 {
		q = q.Offset(int(query.GetOffset()))
	}

	if query.GetLimit() > 0 {
		q = q.Limit(int(query.GetLimit()))
	}

	var rows []relationshipRow
	if err := q.Get(&rows); err != nil {
		return nil, err
	}

	var list []RelationshipInterface
	for _, r := range rows {
		list = append(list, NewRelationshipFromExistingData(map[string]string{
			COLUMN_ID:                r.ID,
			COLUMN_ENTITY_ID:         r.EntityID,
			COLUMN_RELATED_ENTITY_ID: r.RelatedEntityID,
			COLUMN_RELATIONSHIP_TYPE: r.RelationshipType,
			COLUMN_PARENT_ID:         r.ParentID,
			COLUMN_SEQUENCE:          strconv.Itoa(r.Sequence),
			COLUMN_METADATA:          r.Metadata,
			COLUMN_CREATED_AT:        r.CreatedAt,
		}))
	}

	return list, nil
}

// RelationshipListRelated lists all relationships where the given entity is the related (target) entity
func (st *storeImplementation) RelationshipListRelated(ctx context.Context, relatedEntityID string, relationshipType string) ([]RelationshipInterface, error) {
	return st.RelationshipList(ctx, RelationshipQuery().
		WithRelatedEntityID(relatedEntityID).
		WithRelationshipType(relationshipType))
}

// RelationshipCount counts relationships matching the given fluent query
func (st *storeImplementation) RelationshipCount(ctx context.Context, query RelationshipQueryInterface) (int64, error) {
	if query == nil {
		return 0, errors.New("relationship query cannot be nil")
	}

	if err := query.Validate(); err != nil {
		return 0, err
	}

	q := st.applyRelationshipFilters(st.db.Query().Table(st.relationshipTableName), query)

	var count int64
	if err := q.Count(&count); err != nil {
		return 0, err
	}

	return count, nil
}
