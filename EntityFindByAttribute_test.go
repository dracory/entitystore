package entitystore

import "testing"

func TestEntityFindByAttribute(t *testing.T) {
	db := InitDB("test_entity_find_by_attribute.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Must be NIL:", err)
	}

	entity, err := store.EntityCreateWithTypeAndAttributes("post", map[string]string{
		"path": "/",
	})

	if err != nil {
		t.Fatal("Entity could not be created:", err)
	}

	val, _ := entity.GetString("path", "")
	if val != "/" {
		t.Fatal("Entity attribute mismatch")
	}

	// store.SetDebug(true)

	homePage, err := store.EntityFindByAttribute("post", "path", "/")

	if err != nil {
		t.Fatal("Entity find by attribute failed:", err)
	}

	if homePage == nil {
		t.Fatal("Entity could not be found")
	}
}
