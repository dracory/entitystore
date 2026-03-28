package entitystore

// AttributeTrashQueryInterface defines the contract for attribute trash query builders
type AttributeTrashQueryInterface interface {
	// ID filters
	ID(id string) AttributeTrashQueryInterface
	IDIn(ids ...string) AttributeTrashQueryInterface

	// EntityID filters
	EntityID(entityID string) AttributeTrashQueryInterface
	EntityIDIn(entityIDs ...string) AttributeTrashQueryInterface

	// AttributeKey filters
	AttributeKey(key string) AttributeTrashQueryInterface
	AttributeKeyIn(keys ...string) AttributeTrashQueryInterface

	// DeletedAt filters
	DeletedAtGte(deletedAt string) AttributeTrashQueryInterface
	DeletedAtLte(deletedAt string) AttributeTrashQueryInterface

	// DeletedBy filters
	DeletedBy(deletedBy string) AttributeTrashQueryInterface

	// Sorting
	OrderByDeletedAt(direction string) AttributeTrashQueryInterface
	OrderByCreatedAt(direction string) AttributeTrashQueryInterface

	// Pagination
	Limit(limit int) AttributeTrashQueryInterface
	Offset(offset int) AttributeTrashQueryInterface

	// Execution
	Count() (int64, error)
	Execute() ([]AttributeTrashInterface, error)
	First() (AttributeTrashInterface, error)
}
