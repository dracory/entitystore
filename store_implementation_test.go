package entitystore

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func InitDB(filepath string) *sql.DB {
	// Use a shared in-memory SQLite database per logical filepath so
	// all connections see the same schema during tests.
	dsn := "file:" + filepath + "?mode=memory&cache=shared"
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		panic(err)
	}

	return db
}

func TestStoreCreate(t *testing.T) {
	db := InitDB("test_store_create.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created:", err)
	}

	if store == nil {
		t.Fatal("Store could not be created")
	}
}

func TestStoreAutomigrate(t *testing.T) {
	db := InitDB("test_entity_automigrate.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
	})

	if err != nil {
		t.Fatal("Store could not be created:", err)
	}

	errAutomigrate := store.AutoMigrate(context.Background())

	if errAutomigrate != nil {
		t.Fatal("Automigrate failed:", errAutomigrate)
	}
}
