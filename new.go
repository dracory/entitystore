package entitystore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dracory/sb"
)

// NewStoreOptions define the options for creating a new entity store
type NewStoreOptions struct {
	EntityTableName            string               // Name of the entities table (required)
	AttributeTableName         string               // Name of the attributes table (required)
	EntityTrashTableName       string               // Name of the trashed entities table (default: entities_trash)
	AttributeTrashTableName    string               // Name of the trashed attributes table (default: attributes_trash)
	RelationshipTableName      string               // Name of the relationships table (only if RelationshipsEnabled)
	RelationshipTrashTableName string               // Name of the trashed relationships table (only if RelationshipsEnabled)
	RelationshipsEnabled       bool                 // Enable relationship features (optional)
	TaxonomyTableName          string               // Name of the taxonomies table (only if TaxonomiesEnabled)
	TaxonomyTrashTableName     string               // Name of the trashed taxonomies table (only if TaxonomiesEnabled)
	TaxonomyTermTableName      string               // Name of the taxonomy terms table (only if TaxonomiesEnabled)
	TaxonomyTermTrashTableName string               // Name of the trashed taxonomy terms table (only if TaxonomiesEnabled)
	EntityTaxonomyTableName    string               // Name of the entity-taxonomy assignments table (only if TaxonomiesEnabled)
	TaxonomiesEnabled          bool                 // Enable taxonomy features (optional)
	DB                         *sql.DB              // Database connection (required if Database not set)
	Database                   sb.DatabaseInterface // Database interface (required if DB not set)
	DbDriverName               string               // Database driver name (auto-detected if empty)
	AutomigrateEnabled         bool                 // Automatically create/update tables on startup
	DebugEnabled               bool                 // Enable debug logging
}

// NewStore creates a new entity store with the provided options
// Automatically creates table names if not provided and runs automigration if enabled
// Returns an error if required options are missing or if automigration fails
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
