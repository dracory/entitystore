package entitystore

import (
	"context"
	"errors"
	"log"

	"github.com/doug-martin/goqu/v9"
	"github.com/dromara/carbon/v2"
)

// RelationshipTrash soft-deletes a relationship by moving it to trash
// This operation is atomic - both insert to trash and delete from main table
// succeed together or fail together
func (st *storeImplementation) RelationshipTrash(ctx context.Context, relationshipID string, deletedBy string) (bool, error) {
	if relationshipID == "" {
		return false, errors.New("relationship ID cannot be empty")
	}

	// Find the relationship first
	rel, err := st.RelationshipFind(ctx, relationshipID)
	if err != nil {
		return false, err
	}

	if rel == nil {
		return false, nil
	}

	// Create trash record
	trash := NewRelationshipTrash()
	trash.SetID(rel.ID())
	trash.SetEntityID(rel.GetEntityID())
	trash.SetRelatedEntityID(rel.GetRelatedEntityID())
	trash.SetRelationshipType(rel.GetRelationshipType())
	trash.SetParentID(rel.GetParentID())
	trash.SetSequence(rel.GetSequence())
	trash.SetMetadata(rel.GetMetadata())
	trash.SetCreatedAt(rel.GetCreatedAt())
	trash.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	trash.SetDeletedBy(deletedBy)

	// Begin transaction for atomic operation
	tx, err := st.database.DB().BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	// Insert into trash
	record := goqu.Record{}
	for k, v := range trash.Data() {
		record[k] = v
	}

	q := goqu.Dialect(st.dbDriverName).Insert(st.relationshipTrashTableName).Rows(record)

	sqlStr, params, errSql := q.Prepared(true).ToSQL()
	if errSql != nil {
		return false, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err = tx.ExecContext(ctx, sqlStr, params...)
	if err != nil {
		return false, err
	}

	// Delete from main table
	q2 := goqu.Dialect(st.dbDriverName).
		Delete(st.relationshipTableName).
		Where(goqu.C(COLUMN_ID).Eq(relationshipID))

	sqlStr2, params2, errSql2 := q2.Prepared(true).ToSQL()
	if errSql2 != nil {
		return false, errSql2
	}

	if st.GetDebug() {
		log.Println(sqlStr2)
	}

	result, err := tx.ExecContext(ctx, sqlStr2, params2...)
	if err != nil {
		return false, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return false, err
	}

	affected, _ := result.RowsAffected()
	return affected > 0, nil
}

// RelationshipRestore restores a relationship from trash
// This operation is atomic - both insert to main table and delete from trash
// succeed together or fail together
func (st *storeImplementation) RelationshipRestore(ctx context.Context, relationshipID string) (bool, error) {
	if relationshipID == "" {
		return false, errors.New("relationship ID cannot be empty")
	}

	// Find in trash
	trashItems, err := st.RelationshipTrashList(ctx, RelationshipQueryOptions{
		ID:    relationshipID,
		Limit: 1,
	})

	if err != nil {
		return false, err
	}

	if len(trashItems) == 0 {
		return false, nil
	}

	trash := trashItems[0]

	// Create relationship from trash data
	rel := NewRelationship()
	rel.SetID(trash.ID())
	rel.SetEntityID(trash.GetEntityID())
	rel.SetRelatedEntityID(trash.GetRelatedEntityID())
	rel.SetRelationshipType(trash.GetRelationshipType())
	rel.SetParentID(trash.GetParentID())
	rel.SetSequence(trash.GetSequence())
	rel.SetMetadata(trash.GetMetadata())
	rel.SetCreatedAt(trash.GetCreatedAt())

	// Begin transaction for atomic operation
	tx, err := st.database.DB().BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	// Insert into main table
	record := goqu.Record{}
	for k, v := range rel.Data() {
		record[k] = v
	}

	q := goqu.Dialect(st.dbDriverName).Insert(st.relationshipTableName).Rows(record)

	sqlStr, params, errSql := q.Prepared(true).ToSQL()
	if errSql != nil {
		return false, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err = tx.ExecContext(ctx, sqlStr, params...)
	if err != nil {
		return false, err
	}

	// Delete from trash
	q2 := goqu.Dialect(st.dbDriverName).
		Delete(st.relationshipTrashTableName).
		Where(goqu.C(COLUMN_ID).Eq(relationshipID))

	sqlStr2, params2, errSql2 := q2.Prepared(true).ToSQL()
	if errSql2 != nil {
		return false, errSql2
	}

	if st.GetDebug() {
		log.Println(sqlStr2)
	}

	result, err := tx.ExecContext(ctx, sqlStr2, params2...)
	if err != nil {
		return false, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return false, err
	}

	affected, _ := result.RowsAffected()
	return affected > 0, nil
}

// RelationshipTrashList lists deleted relationships in trash.
// Default sort order is descending by deleted_at (most recent deletions first).
func (st *storeImplementation) RelationshipTrashList(ctx context.Context, options RelationshipQueryOptions) ([]RelationshipTrashInterface, error) {
	q := goqu.Dialect(st.dbDriverName).From(st.relationshipTrashTableName)

	if len(options.IDs) > 0 {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDs))
	}

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if options.EntityID != "" {
		q = q.Where(goqu.C(COLUMN_ENTITY_ID).Eq(options.EntityID))
	}

	if options.RelatedEntityID != "" {
		q = q.Where(goqu.C(COLUMN_RELATED_ENTITY_ID).Eq(options.RelatedEntityID))
	}

	if options.RelationshipType != "" {
		q = q.Where(goqu.C(COLUMN_RELATIONSHIP_TYPE).Eq(options.RelationshipType))
	}

	sortByColumn := COLUMN_DELETED_AT
	sortOrder := "desc"

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

	sqlStr, params, errSql := q.Prepared(true).Select().ToSQL()
	if errSql != nil {
		return nil, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	relationshipMaps, err := st.database.SelectToMapString(ctx, sqlStr, params...)
	if err != nil {
		return nil, err
	}

	var list []RelationshipTrashInterface
	for _, m := range relationshipMaps {
		list = append(list, NewRelationshipTrashFromExistingData(m))
	}

	return list, nil
}
