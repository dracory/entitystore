package entitystore

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/doug-martin/goqu/v9"
	"github.com/dromara/carbon/v2"
)

// AttributeCreate persists a new attribute record
func (st *storeImplementation) AttributeCreate(ctx context.Context, attribute AttributeInterface) error {
	if attribute == nil {
		return errors.New("attribute cannot be nil")
	}

	if attribute.ID() == "" {
		attribute.SetID(GenerateShortID())
	}

	if attribute.GetCreatedAt() == "" {
		attribute.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	if attribute.GetUpdatedAt() == "" {
		attribute.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	record := goqu.Record{}
	for k, v := range attribute.Data() {
		record[k] = v
	}

	q := goqu.Dialect(st.dbDriverName).Insert(st.attributeTableName).Rows(record)

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

// AttributeUpdate persists changes to an existing attribute record
func (st *storeImplementation) AttributeUpdate(ctx context.Context, attribute AttributeInterface) error {
	attribute.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	record := goqu.Record{}
	for k, v := range attribute.Data() {
		record[k] = v
	}

	q := goqu.Dialect(st.dbDriverName).
		Update(st.attributeTableName).
		Where(goqu.C(COLUMN_ID).Eq(attribute.ID())).
		Set(record)

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

// AttributeDelete removes an attribute record by ID
func (st *storeImplementation) AttributeDelete(ctx context.Context, id string) error {
	q := goqu.Dialect(st.dbDriverName).
		Delete(st.attributeTableName).
		Where(goqu.C(COLUMN_ID).Eq(id))

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

// AttributesDeleteByEntityID removes all attributes for an entity
func (st *storeImplementation) AttributesDeleteByEntityID(ctx context.Context, entityID string) error {
	q := goqu.Dialect(st.dbDriverName).
		Delete(st.attributeTableName).
		Where(goqu.C(COLUMN_ENTITY_ID).Eq(entityID))

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

// AttributeFind finds a single attribute by entity ID and attribute key
func (st *storeImplementation) AttributeFind(ctx context.Context, entityID string, attributeKey string) (AttributeInterface, error) {
	if entityID == "" {
		return nil, errors.New("entity id cannot be empty")
	}

	if attributeKey == "" {
		return nil, errors.New("attribute key cannot be empty")
	}

	list, err := st.AttributeList(ctx, AttributeQueryOptions{
		EntityID:     entityID,
		AttributeKey: attributeKey,
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

// AttributeFindByHandle finds a single attribute by entity type, handle, and attribute key
func (st *storeImplementation) AttributeFindByHandle(ctx context.Context, entityType string, entityHandle string, attributeKey string) (AttributeInterface, error) {
	if entityType == "" {
		return nil, errors.New("entity type cannot be empty")
	}

	if entityHandle == "" {
		return nil, errors.New("entity handle cannot be empty")
	}

	if attributeKey == "" {
		return nil, errors.New("attribute key cannot be empty")
	}

	list, err := st.AttributeList(ctx, AttributeQueryOptions{
		EntityType:   entityType,
		EntityHandle: entityHandle,
		AttributeKey: attributeKey,
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

// AttributeList lists attributes matching the given query options
func (st *storeImplementation) AttributeList(ctx context.Context, options AttributeQueryOptions) ([]AttributeInterface, error) {
	q := st.AttributeQuery(options)

	sqlStr, _, errSql := q.ToSQL()
	if errSql != nil {
		return nil, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	attributeMaps, err := st.database.SelectToMapString(ctx, sqlStr)
	if err != nil {
		return nil, err
	}

	var list []AttributeInterface
	for _, m := range attributeMaps {
		list = append(list, NewAttributeFromExistingData(m))
	}

	return list, nil
}

// AttributeSetString upserts a string attribute value for an entity
func (st *storeImplementation) AttributeSetString(ctx context.Context, entityID string, attributeKey string, attributeValue string) error {
	attr, err := st.AttributeFind(ctx, entityID, attributeKey)
	if err != nil {
		return err
	}

	if attr == nil {
		_, err := st.AttributeCreateWithKeyAndValue(ctx, entityID, attributeKey, attributeValue)
		return err
	}

	attr.SetValue(attributeValue)
	return st.AttributeUpdate(ctx, attr)
}

// AttributeSetInt creates a new attribute or updates existing with int value
func (st *storeImplementation) AttributeSetInt(ctx context.Context, entityID string, attributeKey string, attributeValue int64) error {
	attributeValueAsString := strconv.FormatInt(attributeValue, 10)
	return st.AttributeSetString(ctx, entityID, attributeKey, attributeValueAsString)
}

// AttributeSetFloat creates a new attribute or updates existing with float value
func (st *storeImplementation) AttributeSetFloat(ctx context.Context, entityID string, attributeKey string, attributeValue float64) error {
	attributeValueAsString := strconv.FormatFloat(attributeValue, 'f', 30, 64)
	return st.AttributeSetString(ctx, entityID, attributeKey, attributeValueAsString)
}

// AttributeCreateWithKeyAndValue creates a new attribute with the given key and value
func (st *storeImplementation) AttributeCreateWithKeyAndValue(ctx context.Context, entityID string, attributeKey string, attributeValue string) (AttributeInterface, error) {
	attr := NewAttribute()
	attr.SetEntityID(entityID)
	attr.SetKey(attributeKey)
	attr.SetValue(attributeValue)

	if err := st.AttributeCreate(ctx, attr); err != nil {
		return nil, err
	}

	return attr, nil
}

// AttributesSet upserts multiple entity attributes
func (st *storeImplementation) AttributesSet(ctx context.Context, entityID string, attributes map[string]string) error {
	for k, v := range attributes {
		err := st.AttributeSetString(ctx, entityID, k, v)
		if err != nil {
			if st.GetDebug() {
				log.Println(err)
			}
			return err
		}
	}
	return nil
}
