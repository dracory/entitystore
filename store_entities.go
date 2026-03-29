package entitystore

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/doug-martin/goqu/v9"
	"github.com/dromara/carbon/v2"
)

// EntityCreate persists a new entity record
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

	record := goqu.Record{}
	for k, v := range entity.Data() {
		record[k] = v
	}

	q := goqu.Dialect(st.dbDriverName).Insert(st.entityTableName).Rows(record)

	sqlStr, _, errSql := q.ToSQL()
	if errSql != nil {
		return errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := st.database.Exec(ctx, sqlStr)
	return err
}

// EntityUpdate persists changes to an existing entity record
func (st *storeImplementation) EntityUpdate(ctx context.Context, entity EntityInterface) error {
	entity.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	record := goqu.Record{}
	for k, v := range entity.Data() {
		record[k] = v
	}

	q := goqu.Dialect(st.dbDriverName).
		Update(st.entityTableName).
		Where(goqu.C(COLUMN_ID).Eq(entity.ID())).
		Set(record)

	sqlStr, _, errSql := q.ToSQL()
	if errSql != nil {
		return errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := st.database.Exec(ctx, sqlStr)
	if err != nil && st.GetDebug() {
		log.Println(err)
	}

	return err
}

// EntityDelete removes an entity record by ID
func (st *storeImplementation) EntityDelete(ctx context.Context, id string) (bool, error) {
	q := goqu.Dialect(st.dbDriverName).
		Delete(st.entityTableName).
		Where(goqu.C(COLUMN_ID).Eq(id))

	sqlStr, _, errSql := q.ToSQL()
	if errSql != nil {
		return false, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	result, err := st.database.Exec(ctx, sqlStr)
	if err != nil {
		return false, err
	}

	affected, _ := result.RowsAffected()
	return affected > 0, nil
}

// EntityFindByID finds an entity by its ID
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

// EntityList lists entities matching the given query options
func (st *storeImplementation) EntityList(ctx context.Context, options EntityQueryOptions) ([]EntityInterface, error) {
	q := st.EntityQuery(options)

	sqlStr, _, errSql := q.ToSQL()
	if errSql != nil {
		return nil, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	entityMaps, err := st.database.SelectToMapString(ctx, sqlStr)
	if err != nil {
		log.Println("EntityList error:", err)
		return nil, err
	}

	var list []EntityInterface
	for _, m := range entityMaps {
		list = append(list, NewEntityFromExistingData(m))
	}

	return list, nil
}

// EntityCount counts entities
func (st *storeImplementation) EntityCount(ctx context.Context, options EntityQueryOptions) (int64, error) {
	q := st.EntityQuery(options)
	q = q.Limit(1).Select(goqu.COUNT(goqu.Star()).As("count"))
	sqlStr, _, errSql := q.ToSQL()
	if errSql != nil {
		return 0, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	maps, err := st.database.SelectToMapString(ctx, sqlStr)
	if err != nil {
		return 0, err
	}

	if len(maps) == 0 {
		return 0, nil
	}

	count, _ := strconv.ParseInt(maps[0]["count"], 10, 64)
	return count, nil
}

// EntityAttributeList lists all attributes for an entity
func (st *storeImplementation) EntityAttributeList(ctx context.Context, entityID string) ([]AttributeInterface, error) {
	return st.AttributeList(ctx, AttributeQueryOptions{EntityID: entityID})
}

// EntityFindByAttribute finds an entity by type and attribute key/value
func (st *storeImplementation) EntityFindByAttribute(ctx context.Context, entityType string, attributeKey string, attributeValue string) (EntityInterface, error) {
	// Find by attribute first
	attrs, err := st.AttributeList(ctx, AttributeQueryOptions{
		AttributeKey: attributeKey,
		EntityType:   entityType,
	})
	if err != nil {
		return nil, err
	}

	for _, attr := range attrs {
		if attr.GetAttributeValue() == attributeValue {
			return st.EntityFindByID(ctx, attr.GetEntityID())
		}
	}

	trash := NewEntity()
	trash.SetType("unknown")
	return trash, nil
}

// EntityListByAttribute finds entities by type and attribute key/value
func (st *storeImplementation) EntityListByAttribute(ctx context.Context, entityType string, attributeKey string, attributeValue string) ([]EntityInterface, error) {
	// Get all entities of this type
	entities, err := st.EntityList(ctx, EntityQueryOptions{EntityType: entityType})
	if err != nil {
		return nil, err
	}

	var results []EntityInterface
	for _, entity := range entities {
		attr, err := st.AttributeFind(ctx, entity.ID(), attributeKey)
		if err == nil && attr != nil && attr.GetAttributeValue() == attributeValue {
			results = append(results, entity)
		}
	}

	return results, nil
}

// EntityCreateWithType is a shortcut to create an entity by providing only the type
func (st *storeImplementation) EntityCreateWithType(ctx context.Context, entityType string) (EntityInterface, error) {
	entity := NewEntity()
	entity.SetType(entityType)
	if err := st.EntityCreate(ctx, entity); err != nil {
		return entity, err
	}
	return entity, nil
}

// EntityCreateWithTypeAndAttributes creates an entity with attributes
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
