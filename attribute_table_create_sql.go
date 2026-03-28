package entitystore

// SqlCreateAttributeTable returns the SQL for creating the attributes table
func (st *storeImplementation) SqlCreateAttributeTable() string {
	sql := `
		CREATE TABLE IF NOT EXISTS ` + st.attributeTableName + ` (
			id varchar(9) NOT NULL PRIMARY KEY,
			entity_id varchar(9) NOT NULL,
			attribute_key varchar(255) NOT NULL,
			attribute_value text,
			created_at datetime NOT NULL,
			updated_at datetime NOT NULL
		);
	`
	return sql
}
