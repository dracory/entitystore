package entitystore

import (
	"github.com/dracory/sb"
)

// taxonomyTermTableCreateSql returns a SQL string for creating the taxonomy terms table
func (st *storeImplementation) taxonomyTermTableCreateSql() (string, error) {
	sql, err := sb.NewBuilder(st.dbDriverName).
		Table(st.taxonomyTermTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			PrimaryKey: true,
			Length:     9,
		}).
		Column(sb.Column{
			Name:   COLUMN_TAXONOMY_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 9,
		}).
		Column(sb.Column{
			Name:   COLUMN_NAME,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name:   COLUMN_SLUG,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name:   COLUMN_PARENT_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 9,
		}).
		Column(sb.Column{
			Name:    COLUMN_SORT_ORDER,
			Type:    sb.COLUMN_TYPE_INTEGER,
			Default: "0",
		}).
		Column(sb.Column{
			Name: COLUMN_CREATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_UPDATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		CreateIfNotExists()
	if err != nil {
		return "", err
	}
	return sql, nil
}

// taxonomyTermIndexesCreateSql returns SQL for creating taxonomy term indexes
func (st *storeImplementation) taxonomyTermIndexesCreateSql() ([]string, error) {
	sqls := []string{}

	// Index on taxonomy_id for queries filtering by taxonomy
	sql1, err := sb.NewBuilder(st.dbDriverName).
		Table(st.taxonomyTermTableName).
		CreateIndex("idx_"+st.taxonomyTermTableName+"_taxonomy", COLUMN_TAXONOMY_ID)
	if err != nil {
		return nil, err
	}
	sqls = append(sqls, sql1)

	// Index on parent_id for hierarchical queries
	sql2, err := sb.NewBuilder(st.dbDriverName).
		Table(st.taxonomyTermTableName).
		CreateIndex("idx_"+st.taxonomyTermTableName+"_parent", COLUMN_PARENT_ID)
	if err != nil {
		return nil, err
	}
	sqls = append(sqls, sql2)

	// Unique index on taxonomy_id + slug to prevent duplicate slugs within a taxonomy
	sql3, err := sb.NewBuilder(st.dbDriverName).
		Table(st.taxonomyTermTableName).
		CreateUniqueIndex(
			"idx_"+st.taxonomyTermTableName+"_taxonomy_slug",
			COLUMN_TAXONOMY_ID,
			COLUMN_SLUG,
		)
	if err != nil {
		return nil, err
	}
	sqls = append(sqls, sql3)

	return sqls, nil
}
