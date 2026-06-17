package entitystore

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
)

// relationshipTrashRow is used for scanning relationship trash query results
type relationshipTrashRow struct {
	ID               string `db:"id"`
	EntityID         string `db:"entity_id"`
	RelatedEntityID  string `db:"related_entity_id"`
	RelationshipType string `db:"relationship_type"`
	ParentID         string `db:"parent_id"`
	Sequence         int    `db:"sequence"`
	Metadata         string `db:"metadata"`
	CreatedAt        string `db:"created_at"`
	DeletedAt        string `db:"deleted_at"`
	DeletedBy        string `db:"deleted_by"`
}

// RelationshipTrash soft-deletes a relationship by moving it to trash
func (st *storeImplementation) RelationshipTrash(ctx context.Context, relationshipID string, deletedBy string) (bool, error) {
	if relationshipID == "" {
		return false, errors.New("relationship ID cannot be empty")
	}

	rel, err := st.RelationshipFind(ctx, relationshipID)
	if err != nil {
		return false, err
	}

	if rel == nil {
		return false, nil
	}

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

	return true, st.db.Query().Transaction(func(tx orm.Query) error {
		trashRow := map[string]any{}
		for k, v := range trash.Data() {
			trashRow[k] = v
		}

		if st.GetDebug() {
			log.Println("RelationshipTrash insert:", trashRow)
		}

		if err := tx.Table(st.relationshipTrashTableName).Create(trashRow); err != nil {
			return err
		}

		_, err := tx.Table(st.relationshipTableName).Where(COLUMN_ID+" = ?", relationshipID).Delete()
		return err
	})
}

// RelationshipRestore restores a relationship from trash
func (st *storeImplementation) RelationshipRestore(ctx context.Context, relationshipID string) (bool, error) {
	if relationshipID == "" {
		return false, errors.New("relationship ID cannot be empty")
	}

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

	rel := NewRelationship()
	rel.SetID(trash.ID())
	rel.SetEntityID(trash.GetEntityID())
	rel.SetRelatedEntityID(trash.GetRelatedEntityID())
	rel.SetRelationshipType(trash.GetRelationshipType())
	rel.SetParentID(trash.GetParentID())
	rel.SetSequence(trash.GetSequence())
	rel.SetMetadata(trash.GetMetadata())
	rel.SetCreatedAt(trash.GetCreatedAt())

	return true, st.db.Query().Transaction(func(tx orm.Query) error {
		relRow := map[string]any{}
		for k, v := range rel.Data() {
			relRow[k] = v
		}

		if st.GetDebug() {
			log.Println("RelationshipRestore insert:", relRow)
		}

		if err := tx.Table(st.relationshipTableName).Create(relRow); err != nil {
			return err
		}

		_, err := tx.Table(st.relationshipTrashTableName).Where(COLUMN_ID+" = ?", relationshipID).Delete()
		return err
	})
}

// RelationshipTrashList lists deleted relationships in trash
func (st *storeImplementation) RelationshipTrashList(ctx context.Context, options RelationshipQueryOptions) ([]RelationshipTrashInterface, error) {
	q := st.db.Query().Table(st.relationshipTrashTableName)

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

	if options.RelatedEntityID != "" {
		q = q.Where(COLUMN_RELATED_ENTITY_ID+" = ?", options.RelatedEntityID)
	}

	if options.RelationshipType != "" {
		q = q.Where(COLUMN_RELATIONSHIP_TYPE+" = ?", options.RelationshipType)
	}

	sortByColumn := COLUMN_DELETED_AT
	sortOrder := "desc"

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

	var rows []relationshipTrashRow
	if err := q.Get(&rows); err != nil {
		return nil, err
	}

	var list []RelationshipTrashInterface
	for _, r := range rows {
		list = append(list, NewRelationshipTrashFromExistingData(map[string]string{
			COLUMN_ID:                r.ID,
			COLUMN_ENTITY_ID:         r.EntityID,
			COLUMN_RELATED_ENTITY_ID: r.RelatedEntityID,
			COLUMN_RELATIONSHIP_TYPE: r.RelationshipType,
			COLUMN_PARENT_ID:         r.ParentID,
			COLUMN_SEQUENCE:          strconv.Itoa(r.Sequence),
			COLUMN_METADATA:          r.Metadata,
			COLUMN_CREATED_AT:        r.CreatedAt,
			COLUMN_DELETED_AT:        r.DeletedAt,
			COLUMN_DELETED_BY:        r.DeletedBy,
		}))
	}

	return list, nil
}
