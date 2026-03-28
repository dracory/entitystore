package entitystore

// SqlCreateAttributeTrashTable returns the SQL for creating the attributes_trash table
func (st *storeImplementation) SqlCreateAttributeTrashTable() string {
	sql := `
		CREATE TABLE IF NOT EXISTS ` + st.attributeTrashTableName + ` (
			id varchar(9) NOT NULL PRIMARY KEY,
			entity_id varchar(9) NOT NULL,
			attribute_key varchar(255) NOT NULL,
			attribute_value text,
			created_at datetime NOT NULL,
			updated_at datetime NOT NULL,
			deleted_at datetime NOT NULL,
			deleted_by varchar(9) DEFAULT ''
		);
	`
	return sql
}
