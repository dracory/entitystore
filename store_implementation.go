package entitystore

import (
	"context"
	"database/sql"

	"github.com/dracory/sb"
)

// storeImplementation implements StoreInterface
type storeImplementation struct {
	entityTableName         string
	attributeTableName      string
	entityTrashTableName    string
	attributeTrashTableName string
	database                sb.DatabaseInterface
	dbDriverName            string
	automigrateEnabled      bool
	debugEnabled            bool
}

// StoreOption options for the vault store
type StoreOption func(*storeImplementation)

// AutoMigrate auto migrate
func (st *storeImplementation) AutoMigrate(ctx context.Context) error {
	sqlArray, err := st.SqlCreateTable()

	if err != nil {
		return err
	}

	for _, sql := range sqlArray {
		_, err := st.database.Exec(ctx, sql)
		if err != nil {
			return err
		}
	}

	return nil
}

// EnableDebug - enables the debug option
func (st *storeImplementation) EnableDebug(debug bool) {
	st.debugEnabled = debug
}

func (st *storeImplementation) GetAttributeTableName() string {
	return st.attributeTableName
}

func (st *storeImplementation) GetAttributeTrashTableName() string {
	return st.attributeTrashTableName
}

func (st *storeImplementation) GetDB() *sql.DB {
	return st.database.DB()
}

func (st *storeImplementation) GetDebug() bool {
	return st.debugEnabled
}

func (st *storeImplementation) GetEntityTableName() string {
	return st.entityTableName
}

func (st *storeImplementation) GetEntityTrashTableName() string {
	return st.entityTrashTableName
}

func (st *storeImplementation) SqlCreateTable() ([]string, error) {
	sqls := []string{}

	// Create entities table
	sql1, err := st.entityTableCreateSql()
	if err != nil {
		return nil, err
	}
	sqls = append(sqls, sql1)

	// Create attributes table
	sql2, err := st.attributeTableCreateSql()
	if err != nil {
		return nil, err
	}
	sqls = append(sqls, sql2)

	// Create entities_trash table
	sql3, err := st.entityTrashTableCreateSql()
	if err != nil {
		return nil, err
	}
	sqls = append(sqls, sql3)

	// Create attributes_trash table
	sql4, err := st.attributeTrashTableCreateSql()
	if err != nil {
		return nil, err
	}
	sqls = append(sqls, sql4)

	return sqls, nil
}
