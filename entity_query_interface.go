package entitystore

// EntityQueryInterface defines the contract for entity query builders
type EntityQueryInterface interface {
	// ID filters
	ID(id string) EntityQueryInterface
	IDIn(ids ...string) EntityQueryInterface

	// EntityType filters
	EntityType(entityType string) EntityQueryInterface
	EntityTypeIn(entityTypes ...string) EntityQueryInterface

	// EntityHandle filters
	EntityHandle(handle string) EntityQueryInterface
	EntityHandleLike(handle string) EntityQueryInterface

	// Time-based filters
	CreatedAtGte(createdAt string) EntityQueryInterface
	CreatedAtLte(createdAt string) EntityQueryInterface
	UpdatedAtGte(updatedAt string) EntityQueryInterface
	UpdatedAtLte(updatedAt string) EntityQueryInterface

	// Sorting
	OrderByCreatedAt(direction string) EntityQueryInterface
	OrderByUpdatedAt(direction string) EntityQueryInterface

	// Pagination
	Limit(limit int) EntityQueryInterface
	Offset(offset int) EntityQueryInterface

	// Execution
	Count() (int64, error)
	Execute() ([]EntityInterface, error)
	First() (EntityInterface, error)
}
