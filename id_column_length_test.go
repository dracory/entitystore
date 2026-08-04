package entitystore

import (
	"context"
	"os"
	"strings"
	"testing"
)

// TestGenerateShortID_FitsColumnWidth is half 1 of the proof:
// GenerateShortID() must produce IDs that fit within ID_COLUMN_LENGTH.
// If this fails, the generated IDs are longer than the declared column width.
func TestGenerateShortID_FitsColumnWidth(t *testing.T) {
	for i := 0; i < 1000; i++ {
		id := GenerateShortID()
		if len(id) > ID_COLUMN_LENGTH {
			t.Errorf("GenerateShortID() returned %q (length %d) which exceeds ID_COLUMN_LENGTH (%d) — "+
				"inserts will fail on MySQL/MariaDB/PostgreSQL with 'Data too long for column'",
				id, len(id), ID_COLUMN_LENGTH)
			return
		}
	}
}

// TestStoreSchema_UsesIDColumnLength is half 2 of the proof:
// store.go must use ID_COLUMN_LENGTH for all id-type column definitions.
// This verifies the schema side — the columns are declared with the same
// constant that TestGenerateShortID_FitsColumnWidth checks against.
func TestStoreSchema_UsesIDColumnLength(t *testing.T) {
	src, err := os.ReadFile("store.go")
	if err != nil {
		t.Fatal("Could not read store.go:", err)
	}

	source := string(src)

	// Count how many times ID_COLUMN_LENGTH is used in table.String() calls
	count := strings.Count(source, "ID_COLUMN_LENGTH)")

	t.Logf("store.go uses ID_COLUMN_LENGTH %d times in schema definitions", count)

	if count == 0 {
		t.Error("store.go does not use ID_COLUMN_LENGTH — id-type columns are not linked to the constant")
	}
}

// TestEntityCreate_RejectsOversizedID verifies that EntityCreate returns an error
// when the entity ID exceeds ID_COLUMN_LENGTH.
func TestEntityCreate_RejectsOversizedID(t *testing.T) {
	db := InitDB("test_oversized_id.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "test_entity",
		AttributeTableName: "test_attribute",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("Store could not be created:", err)
	}

	entity := NewEntity()
	entity.SetType("test_type")
	entity.SetID(strings.Repeat("a", ID_COLUMN_LENGTH+1))

	err = store.EntityCreate(context.Background(), entity)
	if err == nil {
		t.Fatal("EntityCreate should have returned an error for oversized ID")
	}

	t.Logf("EntityCreate correctly rejected oversized ID: %v", err)
}

// TestEntityUpdate_RejectsOversizedID verifies that EntityUpdate returns an error
// when the entity ID exceeds ID_COLUMN_LENGTH.
func TestEntityUpdate_RejectsOversizedID(t *testing.T) {
	db := InitDB("test_oversized_id_update.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "test_entity",
		AttributeTableName: "test_attribute",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("Store could not be created:", err)
	}

	entity := NewEntity()
	entity.SetType("test_type")
	entity.SetID(strings.Repeat("a", ID_COLUMN_LENGTH+1))

	err = store.EntityUpdate(context.Background(), entity)
	if err == nil {
		t.Fatal("EntityUpdate should have returned an error for oversized ID")
	}

	t.Logf("EntityUpdate correctly rejected oversized ID: %v", err)
}

// TestAttributeCreate_RejectsOversizedID verifies that AttributeCreate returns an error
// when the attribute ID or entity_id exceeds ID_COLUMN_LENGTH.
func TestAttributeCreate_RejectsOversizedID(t *testing.T) {
	db := InitDB("test_oversized_attr_id.db")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		EntityTableName:    "test_entity",
		AttributeTableName: "test_attribute",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("Store could not be created:", err)
	}

	// Test oversized attribute ID
	attr := NewAttribute()
	attr.SetID(strings.Repeat("a", ID_COLUMN_LENGTH+1))
	attr.SetEntityID("valid_id")

	err = store.AttributeCreate(context.Background(), attr)
	if err == nil {
		t.Fatal("AttributeCreate should have returned an error for oversized attribute ID")
	}
	t.Logf("AttributeCreate correctly rejected oversized attribute ID: %v", err)

	// Test oversized entity_id
	attr2 := NewAttribute()
	attr2.SetID("valid_id")
	attr2.SetEntityID(strings.Repeat("b", ID_COLUMN_LENGTH+1))

	err = store.AttributeCreate(context.Background(), attr2)
	if err == nil {
		t.Fatal("AttributeCreate should have returned an error for oversized entity_id")
	}
	t.Logf("AttributeCreate correctly rejected oversized entity_id: %v", err)
}
