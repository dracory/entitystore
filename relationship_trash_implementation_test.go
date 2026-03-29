package entitystore

import (
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestNewRelationshipTrash(t *testing.T) {
	trash := NewRelationshipTrash()

	if trash == nil {
		t.Fatal("NewRelationshipTrash() returned nil")
	}

	// Should have deleted_at set
	if trash.GetDeletedAt() == "" {
		t.Error("NewRelationshipTrash() should set DeletedAt")
	}
}

func TestRelationshipTrashSettersAndGetters(t *testing.T) {
	trash := NewRelationshipTrash()

	// Test ID
	trash.SetID("trash_id_123")
	if trash.ID() != "trash_id_123" {
		t.Errorf("SetID/GetID failed: expected 'trash_id_123', got '%s'", trash.ID())
	}

	// Test EntityID
	trash.SetEntityID("entity_123")
	if trash.GetEntityID() != "entity_123" {
		t.Errorf("SetEntityID/GetEntityID failed: expected 'entity_123', got '%s'", trash.GetEntityID())
	}

	// Test RelatedEntityID
	trash.SetRelatedEntityID("related_456")
	if trash.GetRelatedEntityID() != "related_456" {
		t.Errorf("SetRelatedEntityID/GetRelatedEntityID failed: expected 'related_456', got '%s'", trash.GetRelatedEntityID())
	}

	// Test RelationshipType
	trash.SetRelationshipType(RELATIONSHIP_TYPE_HAS_MANY)
	if trash.GetRelationshipType() != RELATIONSHIP_TYPE_HAS_MANY {
		t.Errorf("SetRelationshipType/GetRelationshipType failed: expected '%s', got '%s'", RELATIONSHIP_TYPE_HAS_MANY, trash.GetRelationshipType())
	}

	// Test ParentID
	trash.SetParentID("parent_789")
	if trash.GetParentID() != "parent_789" {
		t.Errorf("SetParentID/GetParentID failed: expected 'parent_789', got '%s'", trash.GetParentID())
	}

	// Test Sequence
	trash.SetSequence(100)
	if trash.GetSequence() != 100 {
		t.Errorf("SetSequence/GetSequence failed: expected 100, got %d", trash.GetSequence())
	}

	// Test Metadata
	trash.SetMetadata("{\"deleted\": true}")
	if trash.GetMetadata() != "{\"deleted\": true}" {
		t.Errorf("SetMetadata/GetMetadata failed: expected '{\"deleted\": true}', got '%s'", trash.GetMetadata())
	}

	// Test CreatedAt
	testTime := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	trash.SetCreatedAt(testTime)
	if trash.GetCreatedAt() != testTime {
		t.Errorf("SetCreatedAt/GetCreatedAt failed: expected '%s', got '%s'", testTime, trash.GetCreatedAt())
	}

	// Test DeletedAt
	deletedTime := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	trash.SetDeletedAt(deletedTime)
	if trash.GetDeletedAt() != deletedTime {
		t.Errorf("SetDeletedAt/GetDeletedAt failed: expected '%s', got '%s'", deletedTime, trash.GetDeletedAt())
	}

	// Test DeletedBy
	trash.SetDeletedBy("user_123")
	if trash.GetDeletedBy() != "user_123" {
		t.Errorf("SetDeletedBy/GetDeletedBy failed: expected 'user_123', got '%s'", trash.GetDeletedBy())
	}
}

func TestRelationshipTrashCreatedAtCarbon(t *testing.T) {
	trash := NewRelationshipTrash()

	testTime := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	trash.SetCreatedAt(testTime)

	carbonTime := trash.GetCreatedAtCarbon()
	if carbonTime == nil {
		t.Error("CreatedAtCarbon() returned nil")
	}

	if carbonTime.ToDateTimeString(carbon.UTC) != testTime {
		t.Errorf("CreatedAtCarbon() returned wrong time: expected '%s', got '%s'", testTime, carbonTime.ToDateTimeString(carbon.UTC))
	}
}

func TestRelationshipTrashDeletedAtCarbon(t *testing.T) {
	trash := NewRelationshipTrash()

	testTime := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	trash.SetDeletedAt(testTime)

	carbonTime := trash.GetDeletedAtCarbon()
	if carbonTime == nil {
		t.Error("DeletedAtCarbon() returned nil")
	}

	if carbonTime.ToDateTimeString(carbon.UTC) != testTime {
		t.Errorf("DeletedAtCarbon() returned wrong time: expected '%s', got '%s'", testTime, carbonTime.ToDateTimeString(carbon.UTC))
	}
}

func TestNewRelationshipTrashFromExistingData(t *testing.T) {
	now := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	data := map[string]string{
		COLUMN_ID:                "trash_123",
		COLUMN_ENTITY_ID:         "entity_456",
		COLUMN_RELATED_ENTITY_ID: "entity_789",
		COLUMN_RELATIONSHIP_TYPE: RELATIONSHIP_TYPE_BELONGS_TO,
		COLUMN_PARENT_ID:         "parent_000",
		COLUMN_SEQUENCE:          "3",
		COLUMN_METADATA:          "{\"reason\": \"test\"}",
		COLUMN_CREATED_AT:        now,
		COLUMN_DELETED_AT:        now,
		COLUMN_DELETED_BY:        "admin_001",
	}

	trash := NewRelationshipTrashFromExistingData(data)

	if trash.ID() != "trash_123" {
		t.Errorf("ID from existing data failed: expected 'trash_123', got '%s'", trash.ID())
	}

	if trash.GetEntityID() != "entity_456" {
		t.Errorf("EntityID from existing data failed: expected 'entity_456', got '%s'", trash.GetEntityID())
	}

	if trash.GetRelatedEntityID() != "entity_789" {
		t.Errorf("RelatedEntityID from existing data failed: expected 'entity_789', got '%s'", trash.GetRelatedEntityID())
	}

	if trash.GetRelationshipType() != RELATIONSHIP_TYPE_BELONGS_TO {
		t.Errorf("RelationshipType from existing data failed: expected '%s', got '%s'", RELATIONSHIP_TYPE_BELONGS_TO, trash.GetRelationshipType())
	}

	if trash.GetParentID() != "parent_000" {
		t.Errorf("ParentID from existing data failed: expected 'parent_000', got '%s'", trash.GetParentID())
	}

	if trash.GetSequence() != 3 {
		t.Errorf("Sequence from existing data failed: expected 3, got %d", trash.GetSequence())
	}

	if trash.GetMetadata() != "{\"reason\": \"test\"}" {
		t.Errorf("Metadata from existing data failed: expected '{\"reason\": \"test\"}', got '%s'", trash.GetMetadata())
	}

	if trash.GetDeletedBy() != "admin_001" {
		t.Errorf("DeletedBy from existing data failed: expected 'admin_001', got '%s'", trash.GetDeletedBy())
	}
}

func TestRelationshipTrashInterfaceCompliance(t *testing.T) {
	var _ RelationshipTrashInterface = (*relationshipTrashImplementation)(nil)
}
