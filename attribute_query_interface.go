package entitystore

// AttributeQueryInterface defines the contract for attribute query builders
type AttributeQueryInterface interface {
	// ID filters
	ID(id string) AttributeQueryInterface
	IDIn(ids ...string) AttributeQueryInterface

	// EntityID filters
	EntityID(entityID string) AttributeQueryInterface
	EntityIDIn(entityIDs ...string) AttributeQueryInterface

	// AttributeKey filters
	AttributeKey(key string) AttributeQueryInterface
	AttributeKeyIn(keys ...string) AttributeQueryInterface
	AttributeKeyLike(key string) AttributeQueryInterface

	// AttributeValue filters
	AttributeValue(value string) AttributeQueryInterface
	AttributeValueLike(value string) AttributeQueryInterface

	// Time-based filters
	CreatedAtGte(createdAt string) AttributeQueryInterface
	CreatedAtLte(createdAt string) AttributeQueryInterface
	UpdatedAtGte(updatedAt string) AttributeQueryInterface
	UpdatedAtLte(updatedAt string) AttributeQueryInterface

	// Sorting
	OrderByCreatedAt(direction string) AttributeQueryInterface
	OrderByUpdatedAt(direction string) AttributeQueryInterface

	// Pagination
	Limit(limit int) AttributeQueryInterface
	Offset(offset int) AttributeQueryInterface

	// Execution
	Count() (int64, error)
	Execute() ([]AttributeInterface, error)
	First() (AttributeInterface, error)
}
