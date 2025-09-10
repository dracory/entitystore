package entitystore

import (
	"testing"
)

func TestAttributeString(t *testing.T) {
	db := InitDB("test_attribute_string.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	errSetString := store.AttributeSetString("default", "hello", "world")

	if errSetString != nil {
		t.Fatal("Attribute could not be created:", errSetString)
	}

	// store.EnableDebug(true)

	attr, err := store.AttributeFind("default", "hello")

	if err != nil {
		t.Fatal("Attribute could not be retrieved:", err)
	}

	if attr == nil {
		t.Fatal("Attribute could not be retrieved")
	}

	if attr.GetString() != "world" {
		t.Fatal("Attribute value incorrect", "must be 'world'", "found", attr.GetString())
	}
}
