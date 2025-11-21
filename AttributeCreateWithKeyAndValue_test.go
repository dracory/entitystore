package entitystore

import (
	"context"
	"testing"
)

func TestAttributeCreateWithKeyAndValue(t *testing.T) {
	db := InitDB("test_attribute_create_with_key_and_value.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	errSet := store.AttributeSetString(context.Background(), "default", "hello", "world")

	if errSet != nil {
		t.Fatal("Attribute could not be created:", errSet.Error())
	}
}
