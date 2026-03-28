package entitystore

import (
	"github.com/dracory/sb"
)

// taxonomyTermTrashTableCreateSql returns a SQL string for creating the taxonomy_terms_trash table
func (st *storeImplementation) taxonomyTermTrashTableCreateSql() (string, error) {
	sql, err := sb.NewBuilder(st.dbDriverName).
		Table(st.taxonomyTermTrashTableName).
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
		Column(sb.Column{
			Name: COLUMN_DELETED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name:   COLUMN_DELETED_BY,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 9,
		}).
		CreateIfNotExists()
	if err != nil {
		return "", err
	}
	return sql, nil
}
