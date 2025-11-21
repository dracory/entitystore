package entitystore

import (
	"context"
	"testing"
)

func TestAttributeInt(t *testing.T) {
	db := InitDB("test_attribute_int.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	errSet := store.AttributeSetInt(context.Background(), "default", "test_int", 12)

	if errSet != nil {
		t.Fatal("Attribute could not be created:", errSet)
	}

	attr, err := store.AttributeFind(context.Background(), "default", "test_int")

	if err != nil {
		t.Fatal("Attribute could not be retrieved:", err)
	}

	if attr == nil {
		t.Fatal("Attribute could not be retrieved")
	}

	v, _ := attr.GetInt()
	if v != 12 {
		t.Fatal("Attribute value incorrect")
	}
}
