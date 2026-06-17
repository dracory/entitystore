package entitystore

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/dracory/neat"
	contractsschema "github.com/dracory/neat/contracts/database/schema"
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
	db                         *neat.Database
	automigrateEnabled         bool
	debugEnabled               bool
	logger                     *slog.Logger
}

// StoreOption options for the vault store
type StoreOption func(*storeImplementation)

// MigrateUp creates the entity store tables
func (st *storeImplementation) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	if err := st.entityTableCreate(); err != nil {
		return err
	}
	if err := st.attributeTableCreate(); err != nil {
		return err
	}
	if err := st.entityTrashTableCreate(); err != nil {
		return err
	}
	if err := st.attributeTrashTableCreate(); err != nil {
		return err
	}
	if st.relationshipsEnabled {
		if err := st.relationshipTableCreate(); err != nil {
			return err
		}
		if err := st.relationshipTrashTableCreate(); err != nil {
			return err
		}
	}
	if st.taxonomiesEnabled {
		if err := st.taxonomyTableCreate(); err != nil {
			return err
		}
		if err := st.taxonomyTermTableCreate(); err != nil {
			return err
		}
		if err := st.entityTaxonomyTableCreate(); err != nil {
			return err
		}
		if err := st.taxonomyTrashTableCreate(); err != nil {
			return err
		}
		if err := st.taxonomyTermTrashTableCreate(); err != nil {
			return err
		}
	}
	return nil
}

// MigrateDown drops the entity store tables
func (st *storeImplementation) MigrateDown(ctx context.Context, tx ...*sql.Tx) error {
	if st.taxonomiesEnabled {
		_ = st.db.Schema().DropIfExists(st.entityTaxonomyTableName)
		_ = st.db.Schema().DropIfExists(st.taxonomyTermTableName)
		_ = st.db.Schema().DropIfExists(st.taxonomyTermTrashTableName)
		_ = st.db.Schema().DropIfExists(st.taxonomyTableName)
		_ = st.db.Schema().DropIfExists(st.taxonomyTrashTableName)
	}
	if st.relationshipsEnabled {
		_ = st.db.Schema().DropIfExists(st.relationshipTableName)
		_ = st.db.Schema().DropIfExists(st.relationshipTrashTableName)
	}
	_ = st.db.Schema().DropIfExists(st.attributeTrashTableName)
	_ = st.db.Schema().DropIfExists(st.attributeTableName)
	_ = st.db.Schema().DropIfExists(st.entityTrashTableName)
	_ = st.db.Schema().DropIfExists(st.entityTableName)
	return nil
}

// EnableDebug - enables the debug option
func (st *storeImplementation) EnableDebug(debug bool) {
	st.debugEnabled = debug
	if debug {
		st.db.EnableDebug()
		st.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		st.db.DisableDebug()
		st.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
}

func (st *storeImplementation) GetAttributeTableName() string {
	return st.attributeTableName
}

func (st *storeImplementation) GetAttributeTrashTableName() string {
	return st.attributeTrashTableName
}

func (st *storeImplementation) GetDB() *sql.DB {
	db, _ := st.db.DB()
	return db
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

func (st *storeImplementation) entityTableCreate() error {
	if st.db.Schema().HasTable(st.entityTableName) {
		return nil
	}
	return st.db.Schema().Create(st.entityTableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 9)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_ENTITY_TYPE, 40)
		table.String(COLUMN_ENTITY_HANDLE, 60)
		table.DateTime(COLUMN_CREATED_AT)
		table.DateTime(COLUMN_UPDATED_AT)
	})
}

func (st *storeImplementation) attributeTableCreate() error {
	if st.db.Schema().HasTable(st.attributeTableName) {
		return nil
	}
	return st.db.Schema().Create(st.attributeTableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 9)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_ENTITY_ID, 9)
		table.String(COLUMN_ATTRIBUTE_KEY, 255)
		table.Text(COLUMN_ATTRIBUTE_VALUE)
		table.DateTime(COLUMN_CREATED_AT)
		table.DateTime(COLUMN_UPDATED_AT)
	})
}

func (st *storeImplementation) entityTrashTableCreate() error {
	if st.db.Schema().HasTable(st.entityTrashTableName) {
		return nil
	}
	return st.db.Schema().Create(st.entityTrashTableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 9)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_ENTITY_TYPE, 40)
		table.String(COLUMN_ENTITY_HANDLE, 60)
		table.DateTime(COLUMN_CREATED_AT)
		table.DateTime(COLUMN_UPDATED_AT)
		table.DateTime(COLUMN_DELETED_AT)
		table.String(COLUMN_DELETED_BY, 9)
	})
}

func (st *storeImplementation) attributeTrashTableCreate() error {
	if st.db.Schema().HasTable(st.attributeTrashTableName) {
		return nil
	}
	return st.db.Schema().Create(st.attributeTrashTableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 9)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_ENTITY_ID, 9)
		table.String(COLUMN_ATTRIBUTE_KEY, 255)
		table.Text(COLUMN_ATTRIBUTE_VALUE)
		table.DateTime(COLUMN_CREATED_AT)
		table.DateTime(COLUMN_UPDATED_AT)
		table.DateTime(COLUMN_DELETED_AT)
		table.String(COLUMN_DELETED_BY, 9)
	})
}

func (st *storeImplementation) relationshipTableCreate() error {
	if st.db.Schema().HasTable(st.relationshipTableName) {
		return nil
	}
	return st.db.Schema().Create(st.relationshipTableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 9)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_ENTITY_ID, 9)
		table.String(COLUMN_RELATED_ENTITY_ID, 9)
		table.String(COLUMN_RELATIONSHIP_TYPE, 50)
		table.String(COLUMN_PARENT_ID, 9).Nullable()
		table.Integer(COLUMN_SEQUENCE)
		table.Text(COLUMN_METADATA).Nullable()
		table.DateTime(COLUMN_CREATED_AT)
		table.Index(COLUMN_ENTITY_ID)
		table.Index(COLUMN_RELATED_ENTITY_ID)
		table.Index(COLUMN_RELATIONSHIP_TYPE)
		table.Unique(COLUMN_ENTITY_ID, COLUMN_RELATED_ENTITY_ID, COLUMN_RELATIONSHIP_TYPE)
	})
}

func (st *storeImplementation) relationshipTrashTableCreate() error {
	if st.db.Schema().HasTable(st.relationshipTrashTableName) {
		return nil
	}
	return st.db.Schema().Create(st.relationshipTrashTableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 9)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_ENTITY_ID, 9)
		table.String(COLUMN_RELATED_ENTITY_ID, 9)
		table.String(COLUMN_RELATIONSHIP_TYPE, 50)
		table.String(COLUMN_PARENT_ID, 9).Nullable()
		table.Integer(COLUMN_SEQUENCE)
		table.Text(COLUMN_METADATA).Nullable()
		table.DateTime(COLUMN_CREATED_AT)
		table.DateTime(COLUMN_DELETED_AT)
		table.String(COLUMN_DELETED_BY, 9)
		table.Index(COLUMN_ENTITY_ID)
		table.Index(COLUMN_RELATED_ENTITY_ID)
		table.Index(COLUMN_RELATIONSHIP_TYPE)
		table.Unique(COLUMN_ENTITY_ID, COLUMN_RELATED_ENTITY_ID, COLUMN_RELATIONSHIP_TYPE)
	})
}

func (st *storeImplementation) taxonomyTableCreate() error {
	if st.db.Schema().HasTable(st.taxonomyTableName) {
		return nil
	}
	return st.db.Schema().Create(st.taxonomyTableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 9)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_NAME, 255)
		table.String(COLUMN_SLUG, 255)
		table.Unique(COLUMN_SLUG)
		table.Text(COLUMN_DESCRIPTION).Nullable()
		table.String(COLUMN_PARENT_ID, 9).Nullable()
		table.Text(COLUMN_ENTITY_TYPES).Nullable()
		table.DateTime(COLUMN_CREATED_AT)
		table.DateTime(COLUMN_UPDATED_AT)
	})
}

func (st *storeImplementation) taxonomyTrashTableCreate() error {
	if st.db.Schema().HasTable(st.taxonomyTrashTableName) {
		return nil
	}
	return st.db.Schema().Create(st.taxonomyTrashTableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 9)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_NAME, 255)
		table.String(COLUMN_SLUG, 255)
		table.Unique(COLUMN_SLUG)
		table.Text(COLUMN_DESCRIPTION).Nullable()
		table.String(COLUMN_PARENT_ID, 9).Nullable()
		table.Text(COLUMN_ENTITY_TYPES).Nullable()
		table.DateTime(COLUMN_CREATED_AT)
		table.DateTime(COLUMN_UPDATED_AT)
		table.DateTime(COLUMN_DELETED_AT)
		table.String(COLUMN_DELETED_BY, 9)
	})
}

func (st *storeImplementation) taxonomyTermTableCreate() error {
	if st.db.Schema().HasTable(st.taxonomyTermTableName) {
		return nil
	}
	return st.db.Schema().Create(st.taxonomyTermTableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 9)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_TAXONOMY_ID, 9)
		table.String(COLUMN_NAME, 255)
		table.String(COLUMN_SLUG, 255)
		table.String(COLUMN_PARENT_ID, 9).Nullable()
		table.Integer(COLUMN_SORT_ORDER)
		table.DateTime(COLUMN_CREATED_AT)
		table.DateTime(COLUMN_UPDATED_AT)
		table.Index(COLUMN_TAXONOMY_ID)
		table.Index(COLUMN_PARENT_ID)
		table.Unique(COLUMN_TAXONOMY_ID, COLUMN_SLUG)
	})
}

func (st *storeImplementation) taxonomyTermTrashTableCreate() error {
	if st.db.Schema().HasTable(st.taxonomyTermTrashTableName) {
		return nil
	}
	return st.db.Schema().Create(st.taxonomyTermTrashTableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 9)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_TAXONOMY_ID, 9)
		table.String(COLUMN_NAME, 255)
		table.String(COLUMN_SLUG, 255)
		table.String(COLUMN_PARENT_ID, 9).Nullable()
		table.Integer(COLUMN_SORT_ORDER)
		table.DateTime(COLUMN_CREATED_AT)
		table.DateTime(COLUMN_UPDATED_AT)
		table.DateTime(COLUMN_DELETED_AT)
		table.String(COLUMN_DELETED_BY, 9)
		table.Index(COLUMN_TAXONOMY_ID)
		table.Index(COLUMN_PARENT_ID)
		table.Unique(COLUMN_TAXONOMY_ID, COLUMN_SLUG)
	})
}

func (st *storeImplementation) entityTaxonomyTableCreate() error {
	if st.db.Schema().HasTable(st.entityTaxonomyTableName) {
		return nil
	}
	return st.db.Schema().Create(st.entityTaxonomyTableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 9)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_ENTITY_ID, 9)
		table.String(COLUMN_TAXONOMY_ID, 9)
		table.String(COLUMN_TERM_ID, 9)
		table.DateTime(COLUMN_CREATED_AT)
		table.Index(COLUMN_ENTITY_ID)
		table.Index(COLUMN_TAXONOMY_ID)
		table.Index(COLUMN_TERM_ID)
		table.Unique(COLUMN_ENTITY_ID, COLUMN_TAXONOMY_ID, COLUMN_TERM_ID)
	})
}
