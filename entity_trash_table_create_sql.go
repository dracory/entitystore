package entitystore

// SqlCreateEntityTrashTable returns the SQL for creating the entities_trash table
func (st *storeImplementation) SqlCreateEntityTrashTable() string {
	sql := `
		CREATE TABLE IF NOT EXISTS ` + st.entityTrashTableName + ` (
			id varchar(9) NOT NULL PRIMARY KEY,
			entity_type varchar(40) NOT NULL,
			entity_handle varchar(60) DEFAULT '',
			created_at datetime NOT NULL,
			updated_at datetime NOT NULL,
			deleted_at datetime NOT NULL,
			deleted_by varchar(9) DEFAULT ''
		);
	`
	return sql
}
