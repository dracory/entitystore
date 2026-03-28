package entitystore

import (
	"github.com/dracory/sb"
)

// relationshipTrashIndexesCreateSql returns SQL strings for creating indexes on the relationships_trash table
func (st *storeImplementation) relationshipTrashIndexesCreateSql() ([]string, error) {
	sqls := []string{}

	// Index on entity_id for queries filtering by source entity
	sql1, err := sb.NewBuilder(st.dbDriverName).
		Table(st.relationshipTrashTableName).
		CreateIndex("idx_"+st.relationshipTrashTableName+"_entity_id", COLUMN_ENTITY_ID)
	if err != nil {
		return nil, err
	}
	sqls = append(sqls, sql1)

	// Index on related_entity_id for queries filtering by target entity
	sql2, err := sb.NewBuilder(st.dbDriverName).
		Table(st.relationshipTrashTableName).
		CreateIndex("idx_"+st.relationshipTrashTableName+"_related_entity_id", COLUMN_RELATED_ENTITY_ID)
	if err != nil {
		return nil, err
	}
	sqls = append(sqls, sql2)

	// Index on deleted_at for sorting and filtering by deletion time
	sql3, err := sb.NewBuilder(st.dbDriverName).
		Table(st.relationshipTrashTableName).
		CreateIndex("idx_"+st.relationshipTrashTableName+"_deleted_at", COLUMN_DELETED_AT)
	if err != nil {
		return nil, err
	}
	sqls = append(sqls, sql3)

	return sqls, nil
}
