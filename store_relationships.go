package entitystore

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/doug-martin/goqu/v9"
	"github.com/dromara/carbon/v2"
)

// RelationshipCreate persists a new relationship record
func (st *storeImplementation) RelationshipCreate(ctx context.Context, relationship RelationshipInterface) error {
	if relationship == nil {
		return errors.New("relationship cannot be nil")
	}

	// Validate required fields
	if relationship.EntityID() == "" {
		return errors.New("entity_id is required")
	}
	if relationship.RelatedEntityID() == "" {
		return errors.New("related_entity_id is required")
	}
	if relationship.RelationshipType() == "" {
		return errors.New("relationship_type is required")
	}

	// Prevent self-referencing relationships for belongs_to and has_many types
	if relationship.EntityID() == relationship.RelatedEntityID() {
		if relationship.RelationshipType() == RELATIONSHIP_TYPE_BELONGS_TO || relationship.RelationshipType() == RELATIONSHIP_TYPE_HAS_MANY {
			return errors.New("self-referencing relationships not allowed for belongs_to and has_many types")
		}
	}

	if relationship.ID() == "" {
		relationship.SetID(GenerateShortID())
	}

	if relationship.CreatedAt() == "" {
		relationship.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	record := goqu.Record{}
	for k, v := range relationship.Data() {
		record[k] = v
	}

	q := goqu.Dialect(st.dbDriverName).Insert(st.relationshipTableName).Rows(record)

	sqlStr, _, errSql := q.ToSQL()
	if errSql != nil {
		return errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := st.database.Exec(ctx, sqlStr)
	return err
}

// RelationshipCreateByOptions creates a relationship using the provided options
func (st *storeImplementation) RelationshipCreateByOptions(ctx context.Context, opts RelationshipOptions) (RelationshipInterface, error) {
	// Check for duplicate relationship
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
	q := goqu.Dialect(st.dbDriverName).
		Delete(st.relationshipTableName).
		Where(goqu.C(COLUMN_ID).Eq(relationshipID))

	sqlStr, _, errSql := q.ToSQL()
	if errSql != nil {
		return false, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	result, err := st.database.Exec(ctx, sqlStr)
	if err != nil {
		return false, err
	}

	affected, _ := result.RowsAffected()
	return affected > 0, nil
}

// RelationshipDeleteAll removes all relationships for an entity (both as source and target)
func (st *storeImplementation) RelationshipDeleteAll(ctx context.Context, entityID string) error {
	// Delete where entity is the source
	q1 := goqu.Dialect(st.dbDriverName).
		Delete(st.relationshipTableName).
		Where(goqu.C(COLUMN_ENTITY_ID).Eq(entityID))

	sqlStr1, _, errSql := q1.ToSQL()
	if errSql != nil {
		return errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr1)
	}

	_, err := st.database.Exec(ctx, sqlStr1)
	if err != nil {
		return err
	}

	// Delete where entity is the target
	q2 := goqu.Dialect(st.dbDriverName).
		Delete(st.relationshipTableName).
		Where(goqu.C(COLUMN_RELATED_ENTITY_ID).Eq(entityID))

	sqlStr2, _, errSql := q2.ToSQL()
	if errSql != nil {
		return errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr2)
	}

	_, err = st.database.Exec(ctx, sqlStr2)
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

// RelationshipList lists relationships matching the given query options.
// Supports filtering by entity IDs, relationship type, parent ID, and pagination.
// Default sort order is ascending by created_at.
func (st *storeImplementation) RelationshipList(ctx context.Context, options RelationshipQueryOptions) ([]RelationshipInterface, error) {
	q := goqu.Dialect(st.dbDriverName).From(st.relationshipTableName)

	if len(options.IDs) > 0 {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDs))
	}

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if options.EntityID != "" {
		q = q.Where(goqu.C(COLUMN_ENTITY_ID).Eq(options.EntityID))
	}

	if len(options.EntityIDs) > 0 {
		q = q.Where(goqu.C(COLUMN_ENTITY_ID).In(options.EntityIDs))
	}

	if options.RelatedEntityID != "" {
		q = q.Where(goqu.C(COLUMN_RELATED_ENTITY_ID).Eq(options.RelatedEntityID))
	}

	if len(options.RelatedEntityIDs) > 0 {
		q = q.Where(goqu.C(COLUMN_RELATED_ENTITY_ID).In(options.RelatedEntityIDs))
	}

	if options.RelationshipType != "" {
		q = q.Where(goqu.C(COLUMN_RELATIONSHIP_TYPE).Eq(options.RelationshipType))
	}

	if options.ParentID != "" {
		q = q.Where(goqu.C(COLUMN_PARENT_ID).Eq(options.ParentID))
	}

	sortByColumn := COLUMN_CREATED_AT
	sortOrder := "asc"

	if options.SortOrder != "" {
		sortOrder = options.SortOrder
	}

	if options.SortBy != "" {
		sortByColumn = options.SortBy
	}

	if sortOrder == "asc" {
		q = q.Order(goqu.I(sortByColumn).Asc())
	} else {
		q = q.Order(goqu.I(sortByColumn).Desc())
	}

	if options.Offset > 0 {
		q = q.Offset(uint(options.Offset))
	}

	if options.Limit > 0 {
		q = q.Limit(uint(options.Limit))
	}

	sqlStr, _, errSql := q.Select().ToSQL()
	if errSql != nil {
		return nil, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	relationshipMaps, err := st.database.SelectToMapString(ctx, sqlStr)
	if err != nil {
		return nil, err
	}

	var list []RelationshipInterface
	for _, m := range relationshipMaps {
		list = append(list, NewRelationshipFromExistingData(m))
	}

	return list, nil
}

// RelationshipListRelated lists all relationships where the given entity is the related (target) entity.
// This is useful for finding all entities that reference the given entity.
func (st *storeImplementation) RelationshipListRelated(ctx context.Context, relatedEntityID string, relationshipType string) ([]RelationshipInterface, error) {
	return st.RelationshipList(ctx, RelationshipQueryOptions{
		RelatedEntityID:  relatedEntityID,
		RelationshipType: relationshipType,
	})
}

// RelationshipCount counts relationships matching the given options.
// Returns the total number of relationships that match the query criteria.
func (st *storeImplementation) RelationshipCount(ctx context.Context, options RelationshipQueryOptions) (int64, error) {
	options.CountOnly = true
	q := goqu.Dialect(st.dbDriverName).From(st.relationshipTableName)

	if options.EntityID != "" {
		q = q.Where(goqu.C(COLUMN_ENTITY_ID).Eq(options.EntityID))
	}

	if options.RelatedEntityID != "" {
		q = q.Where(goqu.C(COLUMN_RELATED_ENTITY_ID).Eq(options.RelatedEntityID))
	}

	if options.RelationshipType != "" {
		q = q.Where(goqu.C(COLUMN_RELATIONSHIP_TYPE).Eq(options.RelationshipType))
	}

	sqlStr, _, errSql := q.Select(goqu.COUNT(goqu.Star()).As("count")).ToSQL()
	if errSql != nil {
		return 0, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	maps, err := st.database.SelectToMapString(ctx, sqlStr)
	if err != nil {
		return 0, err
	}

	if len(maps) == 0 {
		return 0, nil
	}

	count, _ := strconv.ParseInt(maps[0]["count"], 10, 64)
	return count, nil
}
