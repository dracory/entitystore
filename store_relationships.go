package entitystore

import (
	"context"
	"errors"
	"log"
	"strconv"

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

	list, err := st.RelationshipList(ctx, RelationshipQueryOptions{
		ID:    relationshipID,
		Limit: 1,
	})

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

	list, err := st.RelationshipList(ctx, RelationshipQueryOptions{
		EntityID:         entityID,
		RelatedEntityID:  relatedEntityID,
		RelationshipType: relationshipType,
		Limit:            1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// RelationshipList lists relationships matching the given query options
func (st *storeImplementation) RelationshipList(ctx context.Context, options RelationshipQueryOptions) ([]RelationshipInterface, error) {
	q := st.db.Query().Table(st.relationshipTableName)

	if options.ID != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID)
	}

	if len(options.IDs) > 0 {
		ids := make([]any, len(options.IDs))
		for i, id := range options.IDs {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_ID, ids)
	}

	if options.EntityID != "" {
		q = q.Where(COLUMN_ENTITY_ID+" = ?", options.EntityID)
	}

	if len(options.EntityIDs) > 0 {
		ids := make([]any, len(options.EntityIDs))
		for i, id := range options.EntityIDs {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_ENTITY_ID, ids)
	}

	if options.RelatedEntityID != "" {
		q = q.Where(COLUMN_RELATED_ENTITY_ID+" = ?", options.RelatedEntityID)
	}

	if len(options.RelatedEntityIDs) > 0 {
		ids := make([]any, len(options.RelatedEntityIDs))
		for i, id := range options.RelatedEntityIDs {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_RELATED_ENTITY_ID, ids)
	}

	if options.RelationshipType != "" {
		q = q.Where(COLUMN_RELATIONSHIP_TYPE+" = ?", options.RelationshipType)
	}

	if options.ParentID != "" {
		q = q.Where(COLUMN_PARENT_ID+" = ?", options.ParentID)
	}

	sortByColumn := COLUMN_CREATED_AT
	sortOrder := "asc"

	if options.SortOrder != "" {
		sortOrder = options.SortOrder
	}

	if options.SortBy != "" {
		sortByColumn = options.SortBy
	}

	q = q.OrderBy(sortByColumn, sortOrder)

	if options.Offset > 0 {
		q = q.Offset(int(options.Offset))
	}

	if options.Limit > 0 {
		q = q.Limit(int(options.Limit))
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
	return st.RelationshipList(ctx, RelationshipQueryOptions{
		RelatedEntityID:  relatedEntityID,
		RelationshipType: relationshipType,
	})
}

// RelationshipCount counts relationships matching the given options
func (st *storeImplementation) RelationshipCount(ctx context.Context, options RelationshipQueryOptions) (int64, error) {
	q := st.db.Query().Table(st.relationshipTableName)

	if options.ID != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID)
	}

	if len(options.IDs) > 0 {
		ids := make([]any, len(options.IDs))
		for i, id := range options.IDs {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_ID, ids)
	}

	if options.EntityID != "" {
		q = q.Where(COLUMN_ENTITY_ID+" = ?", options.EntityID)
	}

	if len(options.EntityIDs) > 0 {
		ids := make([]any, len(options.EntityIDs))
		for i, id := range options.EntityIDs {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_ENTITY_ID, ids)
	}

	if options.RelatedEntityID != "" {
		q = q.Where(COLUMN_RELATED_ENTITY_ID+" = ?", options.RelatedEntityID)
	}

	if len(options.RelatedEntityIDs) > 0 {
		ids := make([]any, len(options.RelatedEntityIDs))
		for i, id := range options.RelatedEntityIDs {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_RELATED_ENTITY_ID, ids)
	}

	if options.RelationshipType != "" {
		q = q.Where(COLUMN_RELATIONSHIP_TYPE+" = ?", options.RelationshipType)
	}

	if options.ParentID != "" {
		q = q.Where(COLUMN_PARENT_ID+" = ?", options.ParentID)
	}

	var count int64
	if err := q.Count(&count); err != nil {
		return 0, err
	}

	return count, nil
}
