package entitystore

import "testing"

func TestEntityDelete(t *testing.T) {
	db := InitDB("test_entity_delete.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Must be NIL:", err)
	}

	entity, err := store.EntityCreateWithType("post")

	if err != nil {
		t.Fatal("Entity could not be created:", err)
	}

	if entity == nil {
		t.Fatal("Entity could not be created")
	}

	err = entity.SetString("title", "Hello world")

	if err != nil {
		t.Fatal("Entity title could not be created:", err)
	}

	isDeleted, err := store.EntityDelete(entity.ID())

	if err != nil {
		t.Fatal("Entity could not be soft deleted:", err)
	}

	if isDeleted == false {
		t.Fatal("Entity could not be soft deleted")
	}

	val, err := store.EntityFindByID(entity.ID())

	if err != nil {
		t.Fatal(err)
	}

	if val != nil {
		t.Fatal("Entity should no longer be present")
	}
}
