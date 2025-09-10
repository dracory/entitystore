package entitystore

import (
	//"log"
	// "log"
	"database/sql"
	"os"
	"testing"

	//"database/sql"
	_ "github.com/mattn/go-sqlite3"
	// _ "modernc.org/sqlite"
)

func InitDB(filepath string) *sql.DB {
	_ = os.Remove(filepath) // remove database

	dsn := filepath
	db, err := sql.Open("sqlite3", dsn)

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

	errAutomigrate := store.AutoMigrate()

	if errAutomigrate != nil {
		t.Fatal("Automigrate failed:", errAutomigrate)
	}
}
