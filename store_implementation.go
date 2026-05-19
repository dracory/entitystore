package entitystore

import (
	"context"
	"database/sql"

	"github.com/dracory/sb"
)

// storeImplementation implements StoreInterface
type storeImplementation struct {
	entityTableName            string
	attributeTableName         string
	entityTrashTableName       string
	attributeTrashTableName    string
	relationshipTableName      string
	relationshipTrashTableName string
	relationshipsEnabled       bool
	taxonomyTableName          string
	taxonomyTrashTableName     string
	taxonomyTermTableName      string
	taxonomyTermTrashTableName string
	entityTaxonomyTableName    string
	taxonomiesEnabled          bool
	database                   sb.DatabaseInterface
	dbDriverName               string
	automigrateEnabled         bool
	debugEnabled               bool
}

// StoreOption options for the vault store
type StoreOption func(*storeImplementation)

// MigrateUp creates the entity store tables
func (st *storeImplementation) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	sqlArray, err := st.SqlCreateTable()

	if err != nil {
		return err
	}

	for _, sql := range sqlArray {
		var errExec error
		if len(tx) > 0 && tx[0] != nil {
			_, errExec = tx[0].ExecContext(ctx, sql)
		} else {
			_, errExec = st.database.Exec(ctx, sql)
		}
		if errExec != nil {
			return errExec
		}
	}

	return nil
}

// MigrateDown drops the entity store tables
func (st *storeImplementation) MigrateDown(ctx context.Context, tx ...*sql.Tx) error {
	sqlArray, err := st.SqlDropTable()

	if err != nil {
		return err
	}

	for _, sql := range sqlArray {
		var errExec error
		if len(tx) > 0 && tx[0] != nil {
			_, errExec = tx[0].ExecContext(ctx, sql)
		} else {
			_, errExec = st.database.Exec(ctx, sql)
		}
		if errExec != nil {
			return errExec
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

func (st *storeImplementation) GetRelationshipTableName() string {
	return st.relationshipTableName
}

func (st *storeImplementation) GetRelationshipTrashTableName() string {
	return st.relationshipTrashTableName
}

func (st *storeImplementation) GetTaxonomyTableName() string {
	return st.taxonomyTableName
}

func (st *storeImplementation) GetTaxonomyTrashTableName() string {
	return st.taxonomyTrashTableName
}

func (st *storeImplementation) GetTaxonomyTermTableName() string {
	return st.taxonomyTermTableName
}

func (st *storeImplementation) GetTaxonomyTermTrashTableName() string {
	return st.taxonomyTermTrashTableName
}

func (st *storeImplementation) GetEntityTaxonomyTableName() string {
	return st.entityTaxonomyTableName
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

	// Create relationship tables if enabled
	if st.relationshipsEnabled {
		// Create relationships table
		sql5, err := st.relationshipTableCreateSql()
		if err != nil {
			return nil, err
		}
		sqls = append(sqls, sql5)

		// Create relationships_trash table
		sql6, err := st.relationshipTrashTableCreateSql()
		if err != nil {
			return nil, err
		}
		sqls = append(sqls, sql6)

		// Create indexes for relationships table
		sql7, err := st.relationshipIndexesCreateSql()
		if err != nil {
			return nil, err
		}
		sqls = append(sqls, sql7...)

		// Create indexes for relationships_trash table
		sql8, err := st.relationshipTrashIndexesCreateSql()
		if err != nil {
			return nil, err
		}
		sqls = append(sqls, sql8...)
	}

	// Create taxonomy tables if enabled
	if st.taxonomiesEnabled {
		// Create taxonomies table
		sql9, err := st.taxonomyTableCreateSql()
		if err != nil {
			return nil, err
		}
		sqls = append(sqls, sql9)

		// Create taxonomy_terms table
		sql10, err := st.taxonomyTermTableCreateSql()
		if err != nil {
			return nil, err
		}
		sqls = append(sqls, sql10)

		// Create entity_taxonomies table
		sql11, err := st.entityTaxonomyTableCreateSql()
		if err != nil {
			return nil, err
		}
		sqls = append(sqls, sql11)

		// Create taxonomies_trash table
		sql12, err := st.taxonomyTrashTableCreateSql()
		if err != nil {
			return nil, err
		}
		sqls = append(sqls, sql12)

		// Create taxonomy_terms_trash table
		sql13, err := st.taxonomyTermTrashTableCreateSql()
		if err != nil {
			return nil, err
		}
		sqls = append(sqls, sql13)

		// Create indexes for taxonomy_terms table
		sql14, err := st.taxonomyTermIndexesCreateSql()
		if err != nil {
			return nil, err
		}
		sqls = append(sqls, sql14...)

		// Create indexes for entity_taxonomies table
		sql15, err := st.entityTaxonomyIndexesCreateSql()
		if err != nil {
			return nil, err
		}
		sqls = append(sqls, sql15...)
	}

	return sqls, nil
}

// SqlDropTable returns SQL statements for dropping all entity store tables
func (st *storeImplementation) SqlDropTable() ([]string, error) {
	sqls := []string{}

	// Drop in reverse order of creation

	// Drop taxonomy tables if enabled
	if st.taxonomiesEnabled {
		// Drop entity_taxonomies table
		sql15, err := sb.NewBuilder(st.dbDriverName).Table(st.entityTaxonomyTableName).Drop()
		if err != nil {
			return nil, err
		}
		sqls = append(sqls, sql15)

		// Drop taxonomy_terms table
		sql13, err := sb.NewBuilder(st.dbDriverName).Table(st.taxonomyTermTableName).Drop()
		if err != nil {
			return nil, err
		}
		sqls = append(sqls, sql13)

		// Drop taxonomies table
		sql12, err := sb.NewBuilder(st.dbDriverName).Table(st.taxonomyTableName).Drop()
		if err != nil {
			return nil, err
		}
		sqls = append(sqls, sql12)

		// Drop taxonomy_terms_trash table
		sql202, err := sb.NewBuilder(st.dbDriverName).Table(st.taxonomyTermTrashTableName).Drop()
		if err != nil {
			return nil, err
		}
		sqls = append(sqls, sql202)

		// Drop taxonomies_trash table
		sql195, err := sb.NewBuilder(st.dbDriverName).Table(st.taxonomyTrashTableName).Drop()
		if err != nil {
			return nil, err
		}
		sqls = append(sqls, sql195)
	}

	// Drop relationship tables if enabled
	if st.relationshipsEnabled {
		// Drop relationships table
		sql5, err := sb.NewBuilder(st.dbDriverName).Table(st.relationshipTableName).Drop()
		if err != nil {
			return nil, err
		}
		sqls = append(sqls, sql5)

		// Drop relationships_trash table
		sql6, err := sb.NewBuilder(st.dbDriverName).Table(st.relationshipTrashTableName).Drop()
		if err != nil {
			return nil, err
		}
		sqls = append(sqls, sql6)
	}

	// Drop attributes_trash table
	sql4, err := sb.NewBuilder(st.dbDriverName).Table(st.attributeTrashTableName).Drop()
	if err != nil {
		return nil, err
	}
	sqls = append(sqls, sql4)

	// Drop entities_trash table
	sql3, err := sb.NewBuilder(st.dbDriverName).Table(st.entityTrashTableName).Drop()
	if err != nil {
		return nil, err
	}
	sqls = append(sqls, sql3)

	// Drop attributes table
	sql2, err := sb.NewBuilder(st.dbDriverName).Table(st.attributeTableName).Drop()
	if err != nil {
		return nil, err
	}
	sqls = append(sqls, sql2)

	// Drop entities table
	sql1, err := sb.NewBuilder(st.dbDriverName).Table(st.entityTableName).Drop()
	if err != nil {
		return nil, err
	}
	sqls = append(sqls, sql1)

	return sqls, nil
}
