package testutils

import (
	"database/sql"
	"os"

	"github.com/dracory/entitystore"
)

// InitStore creates a new entitystore with in-memory SQLite for testing
func InitStore(filepath string) (entitystore.StoreInterface, error) {
	db := initDB(filepath)

	store, err := entitystore.NewStore(entitystore.NewStoreOptions{
		DB:                      db,
		EntityTableName:         "entities",
		AttributeTableName:        "attributes",
		EntityTrashTableName:    "entities_trash",
		AttributeTrashTableName: "attributes_trash",
		AutomigrateEnabled:      true,
	})

	if err != nil {
		return nil, err
	}

	return store, nil
}

func initDB(filepath string) *sql.DB {
	if filepath != ":memory:" && fileExists(filepath) {
		err := os.Remove(filepath)
		if err != nil {
			panic(err)
		}
	}

	dsn := filepath + "?parseTime=true"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		panic(err)
	}

	return db
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
