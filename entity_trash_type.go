package entitystore

import (
	"time"
)

// EntityTrash type
type EntityTrash struct {
	ID        string    `db:"id"`
	Type      string    `db:"entity_type"`
	Handle    string    `db:"entity_handle"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
	DeletedBy string    `db:"deleted_by"`
}
