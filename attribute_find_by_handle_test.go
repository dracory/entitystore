package entitystore

// func TestAttributeFindByHandle(t *testing.T) {
// 	db := InitDB("test_attribute_find_by_handle.db")

// 	store, err := NewStore(NewStoreOptions{
// 		DB:                 db,
// 		EntityTableName:    "cms_entity",
// 		AttributeTableName: "cms_attribute",
// 		AutomigrateEnabled: true,
// 	})

// 	for _, entityHandle := range []string{"entityHandle1", "entityHandle2", "entityHandle3", "entityHandle4", "entityHandle5", "entityHandle6", "entityHandle7", "entityHandle8"} {
// 		entity, errCreate := store.EntityCreate("default_type")
// 		if errCreate != nil {
// 			t.Fatalf(errCreate.Error())
// 		}
// 		if entity == nil {
// 			t.Fatal("Error MUST NOT BE nil:", entity)
// 		}
// 		entity.SetHandle(entityHandle)
// 		errUpdate := store.EntityUpdate(*entity)
// 		if errUpdate != nil {
// 			t.Fatalf(errUpdate.Error())
// 		}
// 		errSet1 := store.AttributeSetString(entity.ID(), "attr1", "val1")
// 		if errSet1 != nil {
// 			t.Fatalf(errSet1.Error())
// 		}
// 		errSet2 := store.AttributeSetString(entity.ID(), "attr2", "val2")
// 		if errSet2 != nil {
// 			t.Fatalf(errSet2.Error())
// 		}
// 		errSet3 := store.AttributeSetString(entity.ID(), "attr3", "val3")
// 		if errSet3 != nil {
// 			t.Fatalf(errSet3.Error())
// 		}
// 	}

// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}
// 	store.debugEnabled = true

// 	list, _ := store.EntityList(EntityQueryOptions{EntityType: "default_type", EntityHandle: "entityHandle3"})

// 	log.Println(list)

// 	attrs, _ := store.AttributeList(AttributeQueryOptions{EntityID: list[0].ID()})
// 	log.Println(attrs)

// 	attr, errFind := store.AttributeFindByHandle("default_type", "entityHandle3", "attr2")

// 	if errFind != nil {
// 		t.Fatal("Error MUST BE nil:", errFind.Error())
// 	}

// 	if attr == nil {
// 		t.Fatal("Attribute could not be found:", attr)
// 	}

// 	if attr.AttributeValue() != "val2" {
// 		t.Fatal("Attribute value MUST BE val2:", attr.AttributeValue())
// 	}
// }
