package entitystore

import (
	"github.com/dracory/sb"
)

// entityTaxonomyTableCreateSql returns a SQL string for creating the entity_taxonomies table
func (st *storeImplementation) entityTaxonomyTableCreateSql() (string, error) {
	sql, err := sb.NewBuilder(st.dbDriverName).
		Table(st.entityTaxonomyTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			PrimaryKey: true,
			Length:     9,
		}).
		Column(sb.Column{
			Name:   COLUMN_ENTITY_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 9,
		}).
		Column(sb.Column{
			Name:   COLUMN_TAXONOMY_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 9,
		}).
		Column(sb.Column{
			Name:   COLUMN_TERM_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 9,
		}).
		Column(sb.Column{
			Name: COLUMN_CREATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		CreateIfNotExists()
	if err != nil {
		return "", err
	}
	return sql, nil
}

// entityTaxonomyIndexesCreateSql returns SQL for creating entity_taxonomy indexes
func (st *storeImplementation) entityTaxonomyIndexesCreateSql() ([]string, error) {
	sqls := []string{}

	// Index on entity_id for queries filtering by entity
	sql1, err := sb.NewBuilder(st.dbDriverName).
		Table(st.entityTaxonomyTableName).
		CreateIndex("idx_"+st.entityTaxonomyTableName+"_entity", COLUMN_ENTITY_ID)
	if err != nil {
		return nil, err
	}
	sqls = append(sqls, sql1)

	// Index on taxonomy_id for queries filtering by taxonomy
	sql2, err := sb.NewBuilder(st.dbDriverName).
		Table(st.entityTaxonomyTableName).
		CreateIndex("idx_"+st.entityTaxonomyTableName+"_taxonomy", COLUMN_TAXONOMY_ID)
	if err != nil {
		return nil, err
	}
	sqls = append(sqls, sql2)

	// Index on term_id for queries filtering by term
	sql3, err := sb.NewBuilder(st.dbDriverName).
		Table(st.entityTaxonomyTableName).
		CreateIndex("idx_"+st.entityTaxonomyTableName+"_term", COLUMN_TERM_ID)
	if err != nil {
		return nil, err
	}
	sqls = append(sqls, sql3)

	// Unique index on entity_id + taxonomy_id + term_id to prevent duplicate assignments
	sql4, err := sb.NewBuilder(st.dbDriverName).
		Table(st.entityTaxonomyTableName).
		CreateUniqueIndex(
			"idx_"+st.entityTaxonomyTableName+"_entity_term",
			COLUMN_ENTITY_ID,
			COLUMN_TAXONOMY_ID,
			COLUMN_TERM_ID,
		)
	if err != nil {
		return nil, err
	}
	sqls = append(sqls, sql4)

	return sqls, nil
}
