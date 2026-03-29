package entitystore

import (
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestAttributeImplementation(t *testing.T) {
	attr := NewAttribute()

	// Test ID generation
	if attr.ID() == "" {
		t.Error("expected ID to be generated")
	}
	if len(attr.ID()) < 9 {
		t.Errorf("expected ID length >= 9, got %d", len(attr.ID()))
	}

	// Test GetID() method consistency
	if attr.GetID() != attr.ID() {
		t.Errorf("expected GetID() '%s' to match ID() '%s'", attr.GetID(), attr.ID())
	}

	// Test EntityID getter/setter
	attr.SetEntityID("entity123")
	if attr.GetEntityID() != "entity123" {
		t.Errorf("expected EntityID 'entity123', got '%s'", attr.GetEntityID())
	}

	// Test AttributeKey getter/setter
	attr.SetKey("name")
	if attr.GetKey() != "name" {
		t.Errorf("expected AttributeKey 'name', got '%s'", attr.GetKey())
	}

	// Test AttributeValue getter/setter
	attr.SetValue("iPhone")
	if attr.GetValue() != "iPhone" {
		t.Errorf("expected AttributeValue 'iPhone', got '%s'", attr.GetValue())
	}

	// Test CreatedAt
	if attr.GetCreatedAt() == "" {
		t.Error("expected CreatedAt to be set")
	}

	// Test UpdatedAt
	if attr.GetUpdatedAt() == "" {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestAttributeFromExistingData(t *testing.T) {
	data := map[string]string{
		COLUMN_ID:              "attr123",
		COLUMN_ENTITY_ID:       "entity456",
		COLUMN_ATTRIBUTE_KEY:   "price",
		COLUMN_ATTRIBUTE_VALUE: "999",
		COLUMN_CREATED_AT:      carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC),
		COLUMN_UPDATED_AT:      carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC),
	}

	attr := NewAttributeFromExistingData(data)

	if attr.ID() != "attr123" {
		t.Errorf("expected ID 'attr123', got '%s'", attr.ID())
	}

	if attr.GetEntityID() != "entity456" {
		t.Errorf("expected EntityID 'entity456', got '%s'", attr.GetEntityID())
	}

	if attr.GetKey() != "price" {
		t.Errorf("expected AttributeKey 'price', got '%s'", attr.GetKey())
	}

	if attr.GetValue() != "999" {
		t.Errorf("expected AttributeValue '999', got '%s'", attr.GetValue())
	}
}

func TestAttributeTypeConversions(t *testing.T) {
	// Test SetInt/GetInt
	attr := NewAttribute()
	attr.SetInt(42)
	val, err := attr.GetInt()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if val != 42 {
		t.Errorf("expected int 42, got %d", val)
	}

	// Test SetFloat/GetFloat
	attr2 := NewAttribute()
	attr2.SetFloat(3.14)
	fval, err := attr2.GetFloat()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if fval != 3.14 {
		t.Errorf("expected float 3.14, got %f", fval)
	}
}

func TestAttributeFluentInterface(t *testing.T) {
	attr := NewAttribute().
		SetEntityID("entity789").
		SetKey("color").
		SetValue("red")

	if attr.GetEntityID() != "entity789" {
		t.Errorf("expected EntityID 'entity789', got '%s'", attr.GetEntityID())
	}

	if attr.GetKey() != "color" {
		t.Errorf("expected AttributeKey 'color', got '%s'", attr.GetKey())
	}

	if attr.GetValue() != "red" {
		t.Errorf("expected AttributeValue 'red', got '%s'", attr.GetValue())
	}
}
