package entitystore

import (
	"github.com/dracory/sb"
)

// taxonomyTableCreateSql returns a SQL string for creating the taxonomies table
func (st *storeImplementation) taxonomyTableCreateSql() (string, error) {
	sql, err := sb.NewBuilder(st.dbDriverName).
		Table(st.taxonomyTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			PrimaryKey: true,
			Length:     9,
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
			Unique: true,
		}).
		Column(sb.Column{
			Name: COLUMN_DESCRIPTION,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name:   COLUMN_PARENT_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 9,
		}).
		Column(sb.Column{
			Name: COLUMN_ENTITY_TYPES,
			Type: sb.COLUMN_TYPE_TEXT,
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
