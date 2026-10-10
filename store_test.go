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

	errAutomigrate := store.MigrateUp(context.Background())

	if errAutomigrate != nil {
		t.Fatal("Automigrate failed:", errAutomigrate)
	}
}

func TestStoreTableIndexes(t *testing.T) {
	db := InitDB("test_store_table_indexes.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store creation failed:", err)
	}

	tables := []string{
		store.GetEntityTableName(),
		store.GetEntityTrashTableName(),
		store.GetAttributeTableName(),
		store.GetAttributeTrashTableName(),
	}

	for _, tbl := range tables {
		rows, err := db.Query("PRAGMA index_list(" + tbl + ")")
		if err != nil {
			t.Fatalf("Failed to get index list for %s: %v", tbl, err)
		}

		indexCount := 0
		for rows.Next() {
			var seq int
			var name string
			var unique int
			var origin string
			var partial int
			if err := rows.Scan(&seq, &name, &unique, &origin, &partial); err != nil {
				rows.Close()
				t.Fatalf("Failed to scan index list row for %s: %v", tbl, err)
			}
			indexCount++
		}
		rows.Close()

		if indexCount < 2 {
			t.Errorf("Expected at least 2 indexes on table %s, got %d", tbl, indexCount)
		}
	}
}
