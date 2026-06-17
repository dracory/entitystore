package entitystore

import (
	"context"
	"errors"
	"log"

	"github.com/dromara/carbon/v2"
)

// entityRow is used for scanning entity query results
type entityRow struct {
	ID           string `db:"id"`
	EntityType   string `db:"entity_type"`
	EntityHandle string `db:"entity_handle"`
	CreatedAt    string `db:"created_at"`
	UpdatedAt    string `db:"updated_at"`
}

// EntityCreate persists a new entity record to the database
func (st *storeImplementation) EntityCreate(ctx context.Context, entity EntityInterface) error {
	if entity == nil {
		return errors.New("entity cannot be nil")
	}

	if entity.ID() == "" {
		entity.SetID(GenerateShortID())
	}

	if entity.GetCreatedAt() == "" {
		entity.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	if entity.GetUpdatedAt() == "" {
		entity.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	row := map[string]any{}
	for k, v := range entity.Data() {
		row[k] = v
	}

	if st.GetDebug() {
		log.Println("EntityCreate:", row)
	}

	return st.db.Query().Table(st.entityTableName).Create(row)
}

// EntityUpdate persists changes to an existing entity record
func (st *storeImplementation) EntityUpdate(ctx context.Context, entity EntityInterface) error {
	entity.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	row := map[string]any{}
	for k, v := range entity.Data() {
		row[k] = v
	}

	if st.GetDebug() {
		log.Println("EntityUpdate:", row)
	}

	_, err := st.db.Query().Table(st.entityTableName).Where(COLUMN_ID+" = ?", entity.ID()).Update(row)
	return err
}

// EntityDelete permanently removes an entity record by ID
func (st *storeImplementation) EntityDelete(ctx context.Context, id string) (bool, error) {
	result, err := st.db.Query().Table(st.entityTableName).Where(COLUMN_ID+" = ?", id).Delete()
	if err != nil {
		return false, err
	}

	return result.RowsAffected > 0, nil
}

// EntityFindByID finds an entity by its unique ID
func (st *storeImplementation) EntityFindByID(ctx context.Context, entityID string) (EntityInterface, error) {
	if entityID == "" {
		return nil, errors.New("entity ID cannot be empty")
	}

	list, err := st.EntityList(ctx, EntityQueryOptions{
		ID:    entityID,
		Limit: 1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// EntityFindByHandle finds an entity by its type and handle
func (st *storeImplementation) EntityFindByHandle(ctx context.Context, entityType string, entityHandle string) (EntityInterface, error) {
	if entityType == "" {
		return nil, errors.New("entity type cannot be empty")
	}

	if entityHandle == "" {
		return nil, errors.New("entity handle cannot be empty")
	}

	list, err := st.EntityList(ctx, EntityQueryOptions{
		EntityType:   entityType,
		EntityHandle: entityHandle,
		Limit:        1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// EntityList retrieves entities matching the given query options
func (st *storeImplementation) EntityList(ctx context.Context, options EntityQueryOptions) ([]EntityInterface, error) {
	q := st.db.Query().Table(st.entityTableName)

	if options.ID != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID)
	}

	if len(options.IDs) > 0 {
		ids := make([]any, len(options.IDs))
		for i, id := range options.IDs {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_ID, ids)
	}

	if options.EntityType != "" {
		q = q.Where(COLUMN_ENTITY_TYPE+" = ?", options.EntityType)
	}

	if options.EntityHandle != "" {
		q = q.Where(COLUMN_ENTITY_HANDLE+" = ?", options.EntityHandle)
	}

	sortByColumn := COLUMN_ID
	sortOrder := "asc"

	if options.SortOrder != "" {
		sortOrder = options.SortOrder
	}

	if options.SortBy != "" {
		sortByColumn = options.SortBy
	}

	q = q.OrderBy(sortByColumn, sortOrder)

	if options.Offset > 0 {
		q = q.Offset(int(options.Offset))
	}

	if options.Limit > 0 {
		q = q.Limit(int(options.Limit))
	}

	var rows []entityRow
	if err := q.Get(&rows); err != nil {
		return nil, err
	}

	var list []EntityInterface
	for _, r := range rows {
		list = append(list, NewEntityFromExistingData(map[string]string{
			COLUMN_ID:            r.ID,
			COLUMN_ENTITY_TYPE:   r.EntityType,
			COLUMN_ENTITY_HANDLE: r.EntityHandle,
			COLUMN_CREATED_AT:    r.CreatedAt,
			COLUMN_UPDATED_AT:    r.UpdatedAt,
		}))
	}

	return list, nil
}

// EntityCount counts entities matching the given query options
func (st *storeImplementation) EntityCount(ctx context.Context, options EntityQueryOptions) (int64, error) {
	q := st.db.Query().Table(st.entityTableName)

	if options.ID != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID)
	}

	if len(options.IDs) > 0 {
		ids := make([]any, len(options.IDs))
		for i, id := range options.IDs {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_ID, ids)
	}

	if options.EntityType != "" {
		q = q.Where(COLUMN_ENTITY_TYPE+" = ?", options.EntityType)
	}

	if options.EntityHandle != "" {
		q = q.Where(COLUMN_ENTITY_HANDLE+" = ?", options.EntityHandle)
	}

	var count int64
	if err := q.Count(&count); err != nil {
		return 0, err
	}

	return count, nil
}

// EntityAttributeList retrieves all attributes for a given entity
func (st *storeImplementation) EntityAttributeList(ctx context.Context, entityID string) ([]AttributeInterface, error) {
	return st.AttributeList(ctx, AttributeQueryOptions{EntityID: entityID})
}

// EntityFindByAttribute finds an entity by type and attribute key/value
func (st *storeImplementation) EntityFindByAttribute(ctx context.Context, entityType string, attributeKey string, attributeValue string) (EntityInterface, error) {
	results, err := st.EntityListByAttribute(ctx, entityType, attributeKey, attributeValue)
	if err != nil {
		return nil, err
	}
	if len(results) > 0 {
		return results[0], nil
	}
	return nil, nil
}

// EntityListByAttribute finds all entities of a type with a specific attribute value
func (st *storeImplementation) EntityListByAttribute(ctx context.Context, entityType string, attributeKey string, attributeValue string) ([]EntityInterface, error) {
	entities, err := st.EntityList(ctx, EntityQueryOptions{EntityType: entityType})
	if err != nil {
		return nil, err
	}

	var results []EntityInterface
	for _, entity := range entities {
		attr, err := st.AttributeFind(ctx, entity.ID(), attributeKey)
		if err == nil && attr != nil && attr.GetValue() == attributeValue {
			results = append(results, entity)
		}
	}

	return results, nil
}

// EntityCreateWithType creates a new entity with only the type specified
func (st *storeImplementation) EntityCreateWithType(ctx context.Context, entityType string) (EntityInterface, error) {
	entity := NewEntity()
	entity.SetType(entityType)
	if err := st.EntityCreate(ctx, entity); err != nil {
		return entity, err
	}
	return entity, nil
}

// EntityCreateWithTypeAndAttributes creates a new entity with type and attributes
func (st *storeImplementation) EntityCreateWithTypeAndAttributes(ctx context.Context, entityType string, attributes map[string]string) (EntityInterface, error) {
	entity, err := st.EntityCreateWithType(ctx, entityType)
	if err != nil {
		return entity, err
	}

	if len(attributes) > 0 {
		if err := st.AttributesSet(ctx, entity.ID(), attributes); err != nil {
			return entity, err
		}
	}

	return entity, nil
}
