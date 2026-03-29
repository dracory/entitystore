package entitystore

import (
	"context"
	"testing"
)

// TestEntityAttributesCreate tests creating an entity and setting attributes via the store
func TestEntityAttributesCreate(t *testing.T) {
	db := InitDB("test_attributes_create.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "cms_entity",
		AttributeTableName: "cms_attribute",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Must be NIL:", err)
	}

	entity, err := store.EntityCreateWithType(context.Background(), "post")
	if err != nil {
		t.Fatal("Entity could not be created:", err)
	}
	if entity == nil {
		t.Fatal("Entity could not be created")
	}

	// Set title via store (the new pattern — no store back-reference on entity)
	err = store.AttributeSetString(context.Background(), entity.ID(), "title", "Product 1")
	if err != nil {
		t.Fatal("Entity title could not be created:", err)
	}

	titleAttr, err := store.AttributeFind(context.Background(), entity.ID(), "title")
	if err != nil {
		t.Fatal("Entity title could not be retrieved:", err)
	}
	if titleAttr == nil || titleAttr.GetValue() != "Product 1" {
		t.Fatal("Title is incorrect:", titleAttr)
	}

	// Set price_float via store
	err = store.AttributeSetFloat(context.Background(), entity.ID(), "price_float", 12.35)
	if err != nil {
		t.Fatal("Entity price_float could not be created:", err)
	}

	priceFloatAttr, err := store.AttributeFind(context.Background(), entity.ID(), "price_float")
	if err != nil {
		t.Fatal("Entity price_float could not be retrieved:", err)
	}
	if priceFloatAttr == nil {
		t.Fatal("price_float attribute is nil")
	}
	priceFloat, err := priceFloatAttr.GetFloat()
	if err != nil {
		t.Fatal("price_float could not be parsed:", err)
	}
	if priceFloat != 12.35 {
		t.Fatal("Price float is incorrect:", priceFloat)
	}

	// Set price_int via store
	err = store.AttributeSetInt(context.Background(), entity.ID(), "price_int", 12)
	if err != nil {
		t.Fatal("Entity price_int could not be created:", err)
	}

	priceIntAttr, err := store.AttributeFind(context.Background(), entity.ID(), "price_int")
	if err != nil {
		t.Fatal("Entity price_int could not be retrieved:", err)
	}
	if priceIntAttr == nil {
		t.Fatal("price_int attribute is nil")
	}
	priceInt, err := priceIntAttr.GetInt()
	if err != nil {
		t.Fatal("price_int could not be parsed:", err)
	}
	if priceInt != 12 {
		t.Fatal("Price int is incorrect:", priceInt)
	}

	// Set description via store
	err = store.AttributeSetString(context.Background(), entity.ID(), "description", "Description text")
	if err != nil {
		t.Fatal("Entity description could not be created:", err)
	}

	descAttr, err := store.AttributeFind(context.Background(), entity.ID(), "description")
	if err != nil {
		t.Fatal("Entity description could not be retrieved:", err)
	}
	if descAttr == nil || descAttr.GetValue() != "Description text" {
		t.Fatal("Description is incorrect:", descAttr)
	}
}
