package entitystore

// SqlCreateEntityTable returns the SQL for creating the entities table
func (st *storeImplementation) SqlCreateEntityTable() string {
	sql := `
		CREATE TABLE IF NOT EXISTS ` + st.entityTableName + ` (
			id varchar(9) NOT NULL PRIMARY KEY,
			entity_type varchar(40) NOT NULL,
			entity_handle varchar(60) DEFAULT '',
			created_at datetime NOT NULL,
			updated_at datetime NOT NULL
		);
	`
	return sql
}
