package entitystore

import "testing"

func TestAttributeFind(t *testing.T) {
	db := InitDB("test_attribute_find.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	for _, entityID := range []string{"entity1", "entity2", "entity3", "entity4", "entity5", "entity6", "entity7", "entity8"} {
		errSet1 := store.AttributeSetString(entityID, "attr1", "val1")
		if errSet1 != nil {
			t.Fatal(errSet1)
		}
		errSet2 := store.AttributeSetString(entityID, "attr2", "val2")
		if errSet2 != nil {
			t.Fatal(errSet2)
		}
		errSet3 := store.AttributeSetString(entityID, "attr3", "val3")
		if errSet3 != nil {
			t.Fatal(errSet3)
		}
	}

	attr, errFind := store.AttributeFind("entity3", "attr2")

	if errFind != nil {
		t.Fatal("Error MUST BE nil:", errFind)
	}

	if attr == nil {
		t.Fatal("Attribute could not be found:", attr)
	}

	if attr.AttributeValue() != "val2" {
		t.Fatal("Attribute value MUST BE val2:", attr.AttributeValue())
	}
}
