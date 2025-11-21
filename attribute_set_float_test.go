package entitystore

import (
	"context"
	"testing"
)

func TestAttributeSetFloat(t *testing.T) {
	db := InitDB("test_attribute_float.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Must be NIL:", err.Error())
	}

	errSetFloat := store.AttributeSetFloat(context.Background(), "default", "test_float", 12.123456789123456789123456789)

	if errSetFloat != nil {
		t.Fatal("Attribute could not be created:", errSetFloat.Error())
	}

	attr, err := store.AttributeFind(context.Background(), "default", "test_float")

	if err != nil {
		t.Fatal("Attribute could not be retrieved:", err.Error())
	}

	if attr == nil {
		t.Fatal("Attribute could not be retrieved")
	}

	v, _ := attr.GetFloat()
	if v != 12.123456789123456789123456789 {
		t.Fatal("Attribute value incorrect")
	}
}
