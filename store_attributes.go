package entitystore

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/doug-martin/goqu/v9"
	"github.com/dromara/carbon/v2"
)

// AttributeCreate persists a new attribute record to the database
// Automatically generates an ID if not set and timestamps if empty
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

	sqlStr, params, errSql := q.Prepared(true).ToSQL()
	if errSql != nil {
		return errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := st.database.Exec(ctx, sqlStr, params...)
	return err
}

// AttributeUpdate persists changes to an existing attribute record
// Automatically updates the updated_at timestamp
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

	sqlStr, params, errSql := q.Prepared(true).ToSQL()
	if errSql != nil {
		return errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := st.database.Exec(ctx, sqlStr, params...)
	return err
}

// AttributeDelete permanently removes an attribute record by ID
func (st *storeImplementation) AttributeDelete(ctx context.Context, id string) error {
	q := goqu.Dialect(st.dbDriverName).
		Delete(st.attributeTableName).
		Where(goqu.C(COLUMN_ID).Eq(id))

	sqlStr, params, errSql := q.Prepared(true).ToSQL()
	if errSql != nil {
		return errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := st.database.Exec(ctx, sqlStr, params...)
	return err
}

// AttributesDeleteByEntityID permanently removes all attributes for a given entity
func (st *storeImplementation) AttributesDeleteByEntityID(ctx context.Context, entityID string) error {
	q := goqu.Dialect(st.dbDriverName).
		Delete(st.attributeTableName).
		Where(goqu.C(COLUMN_ENTITY_ID).Eq(entityID))

	sqlStr, params, errSql := q.Prepared(true).ToSQL()
	if errSql != nil {
		return errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := st.database.Exec(ctx, sqlStr, params...)
	return err
}

// AttributeFind retrieves a single attribute by entity ID and attribute key
// Returns nil if not found
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

// AttributeFindByHandle retrieves an attribute by entity type, handle, and attribute key
// Returns nil if not found
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

// AttributeList retrieves attributes matching the given query options
// Supports filtering by entity ID, entity type, handle, and attribute key
func (st *storeImplementation) AttributeList(ctx context.Context, options AttributeQueryOptions) ([]AttributeInterface, error) {
	q := st.AttributeQuery(options)

	sqlStr, params, errSql := q.Prepared(true).ToSQL()
	if errSql != nil {
		return nil, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	attributeMaps, err := st.database.SelectToMapString(ctx, sqlStr, params...)
	if err != nil {
		return nil, err
	}

	var list []AttributeInterface
	for _, m := range attributeMaps {
		list = append(list, NewAttributeFromExistingData(m))
	}

	return list, nil
}

// AttributeSetString creates or updates a string attribute value for an entity
// If the attribute doesn't exist, it will be created
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

// AttributeSetInt creates or updates an integer attribute value for an entity
// Converts the int64 to a string for storage
func (st *storeImplementation) AttributeSetInt(ctx context.Context, entityID string, attributeKey string, attributeValue int64) error {
	attributeValueAsString := strconv.FormatInt(attributeValue, 10)
	return st.AttributeSetString(ctx, entityID, attributeKey, attributeValueAsString)
}

// AttributeSetFloat creates or updates a float attribute value for an entity
// Converts the float64 to a string for storage
func (st *storeImplementation) AttributeSetFloat(ctx context.Context, entityID string, attributeKey string, attributeValue float64) error {
	attributeValueAsString := strconv.FormatFloat(attributeValue, 'f', 30, 64)
	return st.AttributeSetString(ctx, entityID, attributeKey, attributeValueAsString)
}

// AttributeGetString retrieves a string attribute value for an entity
// Returns exists=false if the attribute doesn't exist
func (st *storeImplementation) AttributeGetString(ctx context.Context, entityID string, attributeKey string) (value string, exists bool, err error) {
	attr, err := st.AttributeFind(ctx, entityID, attributeKey)
	if err != nil {
		return "", false, err
	}
	if attr == nil {
		return "", false, nil
	}
	return attr.GetValue(), true, nil
}

// AttributeGetInt retrieves an int64 attribute value for an entity
// Returns exists=false if the attribute doesn't exist
func (st *storeImplementation) AttributeGetInt(ctx context.Context, entityID string, attributeKey string) (value int64, exists bool, err error) {
	valueStr, exists, err := st.AttributeGetString(ctx, entityID, attributeKey)
	if err != nil || !exists {
		return 0, exists, err
	}
	value, parseErr := strconv.ParseInt(valueStr, 10, 64)
	if parseErr != nil {
		return 0, false, parseErr
	}
	return value, true, nil
}

// AttributeGetFloat retrieves a float64 attribute value for an entity
// Returns exists=false if the attribute doesn't exist
func (st *storeImplementation) AttributeGetFloat(ctx context.Context, entityID string, attributeKey string) (value float64, exists bool, err error) {
	valueStr, exists, err := st.AttributeGetString(ctx, entityID, attributeKey)
	if err != nil || !exists {
		return 0, exists, err
	}
	value, parseErr := strconv.ParseFloat(valueStr, 64)
	if parseErr != nil {
		return 0, false, parseErr
	}
	return value, true, nil
}

// AttributeCreateWithKeyAndValue creates a new attribute with the given key and value
// Convenience method that creates and persists the attribute in one call
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

// AttributesSet creates or updates multiple entity attributes in a batch
// If any attribute fails, the error is returned immediately
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
