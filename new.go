package entitystore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dracory/sb"
)

// NewStoreOptions define the options for creating a new entity store
type NewStoreOptions struct {
	EntityTableName            string
	AttributeTableName         string
	EntityTrashTableName       string
	AttributeTrashTableName    string
	RelationshipTableName      string
	RelationshipTrashTableName string
	RelationshipsEnabled       bool
	TaxonomyTableName          string
	TaxonomyTrashTableName     string
	TaxonomyTermTableName      string
	TaxonomyTermTrashTableName string
	EntityTaxonomyTableName    string
	TaxonomiesEnabled          bool
	DB                         *sql.DB
	Database                   sb.DatabaseInterface
	DbDriverName               string
	AutomigrateEnabled         bool
	DebugEnabled               bool
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
		entityTableName:            opts.EntityTableName,
		attributeTableName:         opts.AttributeTableName,
		entityTrashTableName:       opts.EntityTrashTableName,
		attributeTrashTableName:    opts.AttributeTrashTableName,
		relationshipTableName:      opts.RelationshipTableName,
		relationshipTrashTableName: opts.RelationshipTrashTableName,
		relationshipsEnabled:       opts.RelationshipsEnabled,
		taxonomyTableName:          opts.TaxonomyTableName,
		taxonomyTrashTableName:     opts.TaxonomyTrashTableName,
		taxonomyTermTableName:      opts.TaxonomyTermTableName,
		taxonomyTermTrashTableName: opts.TaxonomyTermTrashTableName,
		entityTaxonomyTableName:    opts.EntityTaxonomyTableName,
		taxonomiesEnabled:          opts.TaxonomiesEnabled,
		automigrateEnabled:         opts.AutomigrateEnabled,
		database:                   opts.Database,
		dbDriverName:               opts.DbDriverName,
		debugEnabled:               opts.DebugEnabled,
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

	// Set default relationship table names if relationships are enabled
	if store.relationshipsEnabled {
		if store.relationshipTableName == "" {
			store.relationshipTableName = DEFAULT_RELATIONSHIP_TABLE_NAME
		}
		if store.relationshipTrashTableName == "" {
			store.relationshipTrashTableName = DEFAULT_RELATIONSHIP_TRASH_TABLE_NAME
		}
	}

	// Set default taxonomy table names if taxonomies are enabled
	if store.taxonomiesEnabled {
		if store.taxonomyTableName == "" {
			store.taxonomyTableName = DEFAULT_TAXONOMY_TABLE_NAME
		}
		if store.taxonomyTrashTableName == "" {
			store.taxonomyTrashTableName = DEFAULT_TAXONOMY_TRASH_TABLE_NAME
		}
		if store.taxonomyTermTableName == "" {
			store.taxonomyTermTableName = DEFAULT_TAXONOMY_TERM_TABLE_NAME
		}
		if store.taxonomyTermTrashTableName == "" {
			store.taxonomyTermTrashTableName = DEFAULT_TAXONOMY_TERM_TRASH_TABLE_NAME
		}
		if store.entityTaxonomyTableName == "" {
			store.entityTaxonomyTableName = DEFAULT_ENTITY_TAXONOMY_TABLE_NAME
		}
	}

	if store.automigrateEnabled {
		if err := store.AutoMigrate(context.Background()); err != nil {
			return nil, err
		}
	}

	return store, nil
}
