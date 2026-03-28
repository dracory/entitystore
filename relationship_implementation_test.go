package entitystore

import (
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestNewRelationship(t *testing.T) {
	rel := NewRelationship()

	if rel == nil {
		t.Fatal("NewRelationship() returned nil")
	}
}

func TestRelationshipSettersAndGetters(t *testing.T) {
	rel := NewRelationship()

	// Test ID
	rel.SetID("test_id_123")
	if rel.ID() != "test_id_123" {
		t.Errorf("SetID/GetID failed: expected 'test_id_123', got '%s'", rel.ID())
	}

	// Test EntityID
	rel.SetEntityID("entity_123")
	if rel.EntityID() != "entity_123" {
		t.Errorf("SetEntityID/GetEntityID failed: expected 'entity_123', got '%s'", rel.EntityID())
	}

	// Test RelatedEntityID
	rel.SetRelatedEntityID("related_456")
	if rel.RelatedEntityID() != "related_456" {
		t.Errorf("SetRelatedEntityID/GetRelatedEntityID failed: expected 'related_456', got '%s'", rel.RelatedEntityID())
	}

	// Test RelationshipType
	rel.SetRelationshipType(RELATIONSHIP_TYPE_BELONGS_TO)
	if rel.RelationshipType() != RELATIONSHIP_TYPE_BELONGS_TO {
		t.Errorf("SetRelationshipType/GetRelationshipType failed: expected '%s', got '%s'", RELATIONSHIP_TYPE_BELONGS_TO, rel.RelationshipType())
	}

	// Test ParentID
	rel.SetParentID("parent_789")
	if rel.ParentID() != "parent_789" {
		t.Errorf("SetParentID/GetParentID failed: expected 'parent_789', got '%s'", rel.ParentID())
	}

	// Test Sequence
	rel.SetSequence(42)
	if rel.Sequence() != 42 {
		t.Errorf("SetSequence/GetSequence failed: expected 42, got %d", rel.Sequence())
	}

	// Test Metadata
	rel.SetMetadata("{\"key\": \"value\"}")
	if rel.Metadata() != "{\"key\": \"value\"}" {
		t.Errorf("SetMetadata/GetMetadata failed: expected '{\"key\": \"value\"}', got '%s'", rel.Metadata())
	}

	// Test CreatedAt
	testTime := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	rel.SetCreatedAt(testTime)
	if rel.CreatedAt() != testTime {
		t.Errorf("SetCreatedAt/GetCreatedAt failed: expected '%s', got '%s'", testTime, rel.CreatedAt())
	}
}

func TestRelationshipCreatedAtCarbon(t *testing.T) {
	rel := NewRelationship()

	testTime := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	rel.SetCreatedAt(testTime)

	carbonTime := rel.CreatedAtCarbon()
	if carbonTime == nil {
		t.Error("CreatedAtCarbon() returned nil")
	}

	if carbonTime.ToDateTimeString(carbon.UTC) != testTime {
		t.Errorf("CreatedAtCarbon() returned wrong time: expected '%s', got '%s'", testTime, carbonTime.ToDateTimeString(carbon.UTC))
	}
}

func TestNewRelationshipFromExistingData(t *testing.T) {
	data := map[string]string{
		COLUMN_ID:                "rel_123",
		COLUMN_ENTITY_ID:         "entity_456",
		COLUMN_RELATED_ENTITY_ID: "entity_789",
		COLUMN_RELATIONSHIP_TYPE: RELATIONSHIP_TYPE_MANY_MANY,
		COLUMN_PARENT_ID:         "parent_000",
		COLUMN_SEQUENCE:          "5",
		COLUMN_METADATA:          "{\"status\": \"active\"}",
		COLUMN_CREATED_AT:      carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC),
	}

	rel := NewRelationshipFromExistingData(data)

	if rel.ID() != "rel_123" {
		t.Errorf("ID from existing data failed: expected 'rel_123', got '%s'", rel.ID())
	}

	if rel.EntityID() != "entity_456" {
		t.Errorf("EntityID from existing data failed: expected 'entity_456', got '%s'", rel.EntityID())
	}

	if rel.RelatedEntityID() != "entity_789" {
		t.Errorf("RelatedEntityID from existing data failed: expected 'entity_789', got '%s'", rel.RelatedEntityID())
	}

	if rel.RelationshipType() != RELATIONSHIP_TYPE_MANY_MANY {
		t.Errorf("RelationshipType from existing data failed: expected '%s', got '%s'", RELATIONSHIP_TYPE_MANY_MANY, rel.RelationshipType())
	}

	if rel.ParentID() != "parent_000" {
		t.Errorf("ParentID from existing data failed: expected 'parent_000', got '%s'", rel.ParentID())
	}

	if rel.Sequence() != 5 {
		t.Errorf("Sequence from existing data failed: expected 5, got %d", rel.Sequence())
	}

	if rel.Metadata() != "{\"status\": \"active\"}" {
		t.Errorf("Metadata from existing data failed: expected '{\"status\": \"active\"}', got '%s'", rel.Metadata())
	}
}

func TestRelationshipSequenceAsString(t *testing.T) {
	rel := NewRelationship()

	// Test that sequence is stored as string but retrieved as int
	rel.SetSequence(0)
	if rel.Sequence() != 0 {
		t.Errorf("Sequence 0 failed: expected 0, got %d", rel.Sequence())
	}

	rel.SetSequence(999)
	if rel.Sequence() != 999 {
		t.Errorf("Sequence 999 failed: expected 999, got %d", rel.Sequence())
	}

	// Test with negative (should handle gracefully)
	rel.SetSequence(-1)
	if rel.Sequence() != -1 {
		t.Errorf("Sequence -1 failed: expected -1, got %d", rel.Sequence())
	}
}

func TestRelationshipInterfaceCompliance(t *testing.T) {
	var _ RelationshipInterface = (*relationshipImplementation)(nil)
}
