package entitystore

import (
	"context"
	"testing"
)

func TestAttributesSet(t *testing.T) {
	db := InitDB("test_attributes_set.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	attributes := map[string]string{
		"attribute 1": "value 1",
		"attribute 2": "value 2",
		"attribute 3": "value 3",
	}

	errSet := store.AttributesSet(context.Background(), "ENTITY_ID", attributes)

	if errSet != nil {
		t.Fatal(errSet)
	}

	for key, value := range attributes {
		attr, err := store.AttributeFind(context.Background(), "ENTITY_ID", key)

		if err != nil {
			t.Fatal(err)
		}

		if attr == nil {
			t.Fatal("Attribute could not be nil")
		}

		if attr.GetString() != value {
			t.Fatal("Attribute value mismatch")
		}

	}
}
