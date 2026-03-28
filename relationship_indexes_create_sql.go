package entitystore

import (
	"github.com/dracory/sb"
)

// relationshipIndexesCreateSql returns SQL strings for creating indexes on the relationships table
func (st *storeImplementation) relationshipIndexesCreateSql() ([]string, error) {
	sqls := []string{}

	// Index on entity_id for queries filtering by source entity
	sql1, err := sb.NewBuilder(st.dbDriverName).
		Table(st.relationshipTableName).
		CreateIndex("idx_"+st.relationshipTableName+"_entity_id", COLUMN_ENTITY_ID)
	if err != nil {
		return nil, err
	}
	sqls = append(sqls, sql1)

	// Index on related_entity_id for queries filtering by target entity
	sql2, err := sb.NewBuilder(st.dbDriverName).
		Table(st.relationshipTableName).
		CreateIndex("idx_"+st.relationshipTableName+"_related_entity_id", COLUMN_RELATED_ENTITY_ID)
	if err != nil {
		return nil, err
	}
	sqls = append(sqls, sql2)

	// Index on relationship_type for queries filtering by type
	sql3, err := sb.NewBuilder(st.dbDriverName).
		Table(st.relationshipTableName).
		CreateIndex("idx_"+st.relationshipTableName+"_type", COLUMN_RELATIONSHIP_TYPE)
	if err != nil {
		return nil, err
	}
	sqls = append(sqls, sql3)

	// Unique composite index to prevent duplicate relationships
	sql4, err := sb.NewBuilder(st.dbDriverName).
		Table(st.relationshipTableName).
		CreateUniqueIndex(
			"idx_"+st.relationshipTableName+"_unique_composite",
			COLUMN_ENTITY_ID,
			COLUMN_RELATED_ENTITY_ID,
			COLUMN_RELATIONSHIP_TYPE,
		)
	if err != nil {
		return nil, err
	}
	sqls = append(sqls, sql4)

	return sqls, nil
}
