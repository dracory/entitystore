package entitystore

// EntityTrashQueryInterface defines the contract for entity trash query builders
type EntityTrashQueryInterface interface {
	// ID filters
	ID(id string) EntityTrashQueryInterface
	IDIn(ids ...string) EntityTrashQueryInterface

	// EntityType filters
	EntityType(entityType string) EntityTrashQueryInterface
	EntityTypeIn(entityTypes ...string) EntityTrashQueryInterface

	// DeletedAt filters
	DeletedAtGte(deletedAt string) EntityTrashQueryInterface
	DeletedAtLte(deletedAt string) EntityTrashQueryInterface

	// DeletedBy filters
	DeletedBy(deletedBy string) EntityTrashQueryInterface

	// Sorting
	OrderByDeletedAt(direction string) EntityTrashQueryInterface
	OrderByCreatedAt(direction string) EntityTrashQueryInterface

	// Pagination
	Limit(limit int) EntityTrashQueryInterface
	Offset(offset int) EntityTrashQueryInterface

	// Execution
	Count() (int64, error)
	Execute() ([]EntityTrashInterface, error)
	First() (EntityTrashInterface, error)
}
