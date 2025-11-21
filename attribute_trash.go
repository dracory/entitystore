package entitystore

import (
	"time"
)

// AttributeTrash type
type AttributeTrash struct {
	ID             string    `db:"id"`
	EntityID       string    `db:"entity_id"`
	AttributeKey   string    `db:"attribute_key"`
	AttributeValue string    `db:"attribute_value"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	DeletedAt      time.Time `db:"deleted_at"`
	DeletedBy      string    `db:"deleted_by"`
}
