package entitystore

import (
	"github.com/dracory/sb"
)

// relationshipTableCreateSql returns a SQL string for creating the relationships table
func (st *storeImplementation) relationshipTableCreateSql() (string, error) {
	sql, err := sb.NewBuilder(st.dbDriverName).
		Table(st.relationshipTableName).
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
			Name:   COLUMN_RELATED_ENTITY_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 9,
		}).
		Column(sb.Column{
			Name:   COLUMN_RELATIONSHIP_TYPE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 50,
		}).
		Column(sb.Column{
			Name:     COLUMN_PARENT_ID,
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   9,
			Nullable: true,
		}).
		Column(sb.Column{
			Name:    COLUMN_SEQUENCE,
			Type:    sb.COLUMN_TYPE_INTEGER,
			Default: "0",
		}).
		Column(sb.Column{
			Name:     COLUMN_METADATA,
			Type:     sb.COLUMN_TYPE_TEXT,
			Nullable: true,
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
