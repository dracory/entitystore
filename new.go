package entitystore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dracory/sb"
)

// NewStoreOptions define the options for creating a new entity store
type NewStoreOptions struct {
	EntityTableName         string
	AttributeTableName      string
	EntityTrashTableName    string
	AttributeTrashTableName string
	DB                      *sql.DB
	Database                sb.DatabaseInterface
	DbDriverName            string
	AutomigrateEnabled      bool
	DebugEnabled            bool
}

// NewStore creates a new entity store
func NewStore(opts NewStoreOptions) (StoreInterface, error) {
	if opts.DB == nil && opts.Database == nil {
		return nil, errors.New("entity store: DB or Database is required")
	}

	if opts.DbDriverName == "" {
		if opts.DB != nil {
			opts.DbDriverName = sb.DatabaseDriverName(opts.DB)
		}
		if opts.Database != nil {
			opts.DbDriverName = sb.DatabaseDriverName(opts.Database.DB())
		}
	}

	if opts.Database == nil {
		opts.Database = sb.NewDatabase(opts.DB, opts.DbDriverName)
	}

	store := &storeImplementation{
		entityTableName:         opts.EntityTableName,
		attributeTableName:      opts.AttributeTableName,
		entityTrashTableName:    opts.EntityTrashTableName,
		attributeTrashTableName: opts.AttributeTrashTableName,
		automigrateEnabled:      opts.AutomigrateEnabled,
		database:                opts.Database,
		dbDriverName:            opts.DbDriverName,
		debugEnabled:            opts.DebugEnabled,
	}

	if store.entityTableName == "" {
		return nil, errors.New("entity store: entityTableName is required")
	}

	if store.attributeTableName == "" {
		return nil, errors.New("entity store: attributeTableName is required")
	}

	if store.entityTrashTableName == "" {
		store.entityTrashTableName = store.entityTableName + "_trash"
	}

	if store.attributeTrashTableName == "" {
		store.attributeTrashTableName = store.attributeTableName + "_trash"
	}

	if store.automigrateEnabled {
		if err := store.AutoMigrate(context.Background()); err != nil {
			return nil, err
		}
	}

	return store, nil
}
