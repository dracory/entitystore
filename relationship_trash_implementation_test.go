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
	if trash.DeletedAt() == "" {
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
	if trash.EntityID() != "entity_123" {
		t.Errorf("SetEntityID/GetEntityID failed: expected 'entity_123', got '%s'", trash.EntityID())
	}

	// Test RelatedEntityID
	trash.SetRelatedEntityID("related_456")
	if trash.RelatedEntityID() != "related_456" {
		t.Errorf("SetRelatedEntityID/GetRelatedEntityID failed: expected 'related_456', got '%s'", trash.RelatedEntityID())
	}

	// Test RelationshipType
	trash.SetRelationshipType(RELATIONSHIP_TYPE_HAS_MANY)
	if trash.RelationshipType() != RELATIONSHIP_TYPE_HAS_MANY {
		t.Errorf("SetRelationshipType/GetRelationshipType failed: expected '%s', got '%s'", RELATIONSHIP_TYPE_HAS_MANY, trash.RelationshipType())
	}

	// Test ParentID
	trash.SetParentID("parent_789")
	if trash.ParentID() != "parent_789" {
		t.Errorf("SetParentID/GetParentID failed: expected 'parent_789', got '%s'", trash.ParentID())
	}

	// Test Sequence
	trash.SetSequence(100)
	if trash.Sequence() != 100 {
		t.Errorf("SetSequence/GetSequence failed: expected 100, got %d", trash.Sequence())
	}

	// Test Metadata
	trash.SetMetadata("{\"deleted\": true}")
	if trash.Metadata() != "{\"deleted\": true}" {
		t.Errorf("SetMetadata/GetMetadata failed: expected '{\"deleted\": true}', got '%s'", trash.Metadata())
	}

	// Test CreatedAt
	testTime := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	trash.SetCreatedAt(testTime)
	if trash.CreatedAt() != testTime {
		t.Errorf("SetCreatedAt/GetCreatedAt failed: expected '%s', got '%s'", testTime, trash.CreatedAt())
	}

	// Test DeletedAt
	deletedTime := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	trash.SetDeletedAt(deletedTime)
	if trash.DeletedAt() != deletedTime {
		t.Errorf("SetDeletedAt/GetDeletedAt failed: expected '%s', got '%s'", deletedTime, trash.DeletedAt())
	}

	// Test DeletedBy
	trash.SetDeletedBy("user_123")
	if trash.DeletedBy() != "user_123" {
		t.Errorf("SetDeletedBy/GetDeletedBy failed: expected 'user_123', got '%s'", trash.DeletedBy())
	}
}

func TestRelationshipTrashCreatedAtCarbon(t *testing.T) {
	trash := NewRelationshipTrash()

	testTime := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	trash.SetCreatedAt(testTime)

	carbonTime := trash.CreatedAtCarbon()
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

	carbonTime := trash.DeletedAtCarbon()
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
		COLUMN_CREATED_AT:      now,
		COLUMN_DELETED_AT:      now,
		COLUMN_DELETED_BY:        "admin_001",
	}

	trash := NewRelationshipTrashFromExistingData(data)

	if trash.ID() != "trash_123" {
		t.Errorf("ID from existing data failed: expected 'trash_123', got '%s'", trash.ID())
	}

	if trash.EntityID() != "entity_456" {
		t.Errorf("EntityID from existing data failed: expected 'entity_456', got '%s'", trash.EntityID())
	}

	if trash.RelatedEntityID() != "entity_789" {
		t.Errorf("RelatedEntityID from existing data failed: expected 'entity_789', got '%s'", trash.RelatedEntityID())
	}

	if trash.RelationshipType() != RELATIONSHIP_TYPE_BELONGS_TO {
		t.Errorf("RelationshipType from existing data failed: expected '%s', got '%s'", RELATIONSHIP_TYPE_BELONGS_TO, trash.RelationshipType())
	}

	if trash.ParentID() != "parent_000" {
		t.Errorf("ParentID from existing data failed: expected 'parent_000', got '%s'", trash.ParentID())
	}

	if trash.Sequence() != 3 {
		t.Errorf("Sequence from existing data failed: expected 3, got %d", trash.Sequence())
	}

	if trash.Metadata() != "{\"reason\": \"test\"}" {
		t.Errorf("Metadata from existing data failed: expected '{\"reason\": \"test\"}', got '%s'", trash.Metadata())
	}

	if trash.DeletedBy() != "admin_001" {
		t.Errorf("DeletedBy from existing data failed: expected 'admin_001', got '%s'", trash.DeletedBy())
	}
}

func TestRelationshipTrashInterfaceCompliance(t *testing.T) {
	var _ RelationshipTrashInterface = (*relationshipTrashImplementation)(nil)
}
