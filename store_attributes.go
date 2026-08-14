package entitystore

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
)

// attributeRow is used for scanning attribute query results
type attributeRow struct {
	ID             string `db:"id"`
	EntityID       string `db:"entity_id"`
	AttributeKey   string `db:"attribute_key"`
	AttributeValue string `db:"attribute_value"`
	CreatedAt      string `db:"created_at"`
	UpdatedAt      string `db:"updated_at"`
}

// AttributeCreate persists a new attribute record to the database
func (st *storeImplementation) AttributeCreate(ctx context.Context, attribute AttributeInterface) error {
	if attribute == nil {
		return errors.New("attribute cannot be nil")
	}

	if attribute.ID() == "" {
		attribute.SetID(GenerateShortID())
	}

	if err := validateIDLength(attribute.ID(), COLUMN_ID); err != nil {
		return err
	}
	if err := validateIDLength(attribute.GetEntityID(), COLUMN_ENTITY_ID); err != nil {
		return err
	}

	if attribute.GetCreatedAt() == "" {
		attribute.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	if attribute.GetUpdatedAt() == "" {
		attribute.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	row := map[string]any{}
	for k, v := range attribute.Data() {
		row[k] = v
	}

	if st.GetDebug() {
		log.Println("AttributeCreate:", row)
	}

	return st.db.Query().Table(st.attributeTableName).Create(row)
}

// AttributeUpdate persists changes to an existing attribute record
func (st *storeImplementation) AttributeUpdate(ctx context.Context, attribute AttributeInterface) error {
	if attribute == nil {
		return errors.New("attribute cannot be nil")
	}

	if err := validateIDRequired(attribute.ID(), COLUMN_ID); err != nil {
		return err
	}
	if err := validateIDLength(attribute.ID(), COLUMN_ID); err != nil {
		return err
	}
	if err := validateIDLength(attribute.GetEntityID(), COLUMN_ENTITY_ID); err != nil {
		return err
	}

	attribute.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	row := map[string]any{}
	for k, v := range attribute.Data() {
		row[k] = v
	}

	if st.GetDebug() {
		log.Println("AttributeUpdate:", row)
	}

	_, err := st.db.Query().Table(st.attributeTableName).Where(COLUMN_ID+" = ?", attribute.ID()).Update(row)
	return err
}

// AttributeDelete permanently removes an attribute record by ID
func (st *storeImplementation) AttributeDelete(ctx context.Context, id string) error {
	_, err := st.db.Query().Table(st.attributeTableName).Where(COLUMN_ID+" = ?", id).Delete()
	return err
}

// AttributesDeleteByEntityID permanently removes all attributes for a given entity
func (st *storeImplementation) AttributesDeleteByEntityID(ctx context.Context, entityID string) error {
	_, err := st.db.Query().Table(st.attributeTableName).Where(COLUMN_ENTITY_ID+" = ?", entityID).Delete()
	return err
}

// AttributeFind retrieves a single attribute by entity ID and attribute key
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

// applyAttributeFilters applies the common filter clauses from AttributeQueryOptions
// to a query. Shared by AttributeList and AttributeCount to prevent filter drift.
// When hasJoin is true, column references are table-qualified to avoid ambiguity.
func (st *storeImplementation) applyAttributeFilters(q orm.Query, options AttributeQueryOptions, hasJoin bool) orm.Query {
	if options.ID != "" {
		q = q.Where(st.attributeTableName+"."+COLUMN_ID+" = ?", options.ID)
	}

	if len(options.IDs) > 0 {
		ids := make([]any, len(options.IDs))
		for i, id := range options.IDs {
			ids[i] = id
		}
		q = q.WhereIn(st.attributeTableName+"."+COLUMN_ID, ids)
	}

	if options.EntityID != "" {
		q = q.Where(st.attributeTableName+"."+COLUMN_ENTITY_ID+" = ?", options.EntityID)
	}

	if options.AttributeKey != "" {
		q = q.Where(st.attributeTableName+"."+COLUMN_ATTRIBUTE_KEY+" = ?", options.AttributeKey)
	}

	if options.EntityType != "" {
		q = q.Where(st.entityTableName+"."+COLUMN_ENTITY_TYPE+" = ?", options.EntityType)
	}

	if options.EntityHandle != "" {
		q = q.Where(st.entityTableName+"."+COLUMN_ENTITY_HANDLE+" = ?", options.EntityHandle)
	}

	return q
}

// AttributeList retrieves attributes matching the given query options.
// When EntityType or EntityHandle is set, the query joins the entities table
// (WP posts/postmeta pattern) to filter attributes by the parent entity's
// type or handle. This avoids N+1 query patterns when loading attributes for
// a specific entity type.
func (st *storeImplementation) AttributeList(ctx context.Context, options AttributeQueryOptions) ([]AttributeInterface, error) {
	q := st.db.Query().Table(st.attributeTableName)

	// Join entities table when filtering by EntityType or EntityHandle
	// (WP wp_postmeta JOIN wp_posts ON post_id = id pattern)
	hasJoin := options.EntityType != "" || options.EntityHandle != ""

	if hasJoin {
		q = q.Join(st.entityTableName + " ON " + st.attributeTableName + "." + COLUMN_ENTITY_ID + " = " + st.entityTableName + "." + COLUMN_ID)

		// Explicit column selection to avoid ambiguous column names (id, created_at,
		// updated_at exist in both tables). Without this, SELECT * on a JOIN fails
		// on MySQL/PostgreSQL with "ambiguous column name" — SQLite silently picks
		// the first table's columns, masking the bug in tests.
		q = q.Select(
			st.attributeTableName + "." + COLUMN_ID + ", " +
				st.attributeTableName + "." + COLUMN_ENTITY_ID + ", " +
				st.attributeTableName + "." + COLUMN_ATTRIBUTE_KEY + ", " +
				st.attributeTableName + "." + COLUMN_ATTRIBUTE_VALUE + ", " +
				st.attributeTableName + "." + COLUMN_CREATED_AT + ", " +
				st.attributeTableName + "." + COLUMN_UPDATED_AT,
		)
	}

	q = st.applyAttributeFilters(q, options, hasJoin)

	sortByColumn := COLUMN_ID
	sortOrder := "asc"

	if options.SortOrder != "" {
		sortOrder = options.SortOrder
	}

	if options.SortBy != "" {
		sortByColumn = options.SortBy
	}

	// Use Order (not OrderBy) when JOIN is active to allow table-qualified column
	// names. OrderBy rejects dotted identifiers via isSimpleIdentifier, which
	// would silently drop the ORDER BY clause. Order uses isValidColumnReference
	// which accepts "table.column" format.
	if hasJoin {
		q = q.Order(st.attributeTableName + "." + sortByColumn + " " + sortOrder)
	} else {
		q = q.OrderBy(sortByColumn, sortOrder)
	}

	if options.Offset > 0 {
		q = q.Offset(int(options.Offset))
	}

	if options.Limit > 0 {
		q = q.Limit(int(options.Limit))
	}

	var rows []attributeRow
	if err := q.Get(&rows); err != nil {
		return nil, err
	}

	var list []AttributeInterface
	for _, r := range rows {
		list = append(list, NewAttributeFromExistingData(map[string]string{
			COLUMN_ID:              r.ID,
			COLUMN_ENTITY_ID:       r.EntityID,
			COLUMN_ATTRIBUTE_KEY:   r.AttributeKey,
			COLUMN_ATTRIBUTE_VALUE: r.AttributeValue,
			COLUMN_CREATED_AT:      r.CreatedAt,
			COLUMN_UPDATED_AT:      r.UpdatedAt,
		}))
	}

	return list, nil
}

// AttributeCount counts attributes matching the given query options.
// Applies the same filters as AttributeList but returns a count instead of rows.
func (st *storeImplementation) AttributeCount(ctx context.Context, options AttributeQueryOptions) (int64, error) {
	q := st.db.Query().Table(st.attributeTableName)

	hasJoin := options.EntityType != "" || options.EntityHandle != ""

	if hasJoin {
		q = q.Join(st.entityTableName + " ON " + st.attributeTableName + "." + COLUMN_ENTITY_ID + " = " + st.entityTableName + "." + COLUMN_ID)
	}

	q = st.applyAttributeFilters(q, options, hasJoin)

	var count int64
	if err := q.Count(&count); err != nil {
		return 0, err
	}

	return count, nil
}

// AttributeSetString creates or updates a string attribute value for an entity
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
func (st *storeImplementation) AttributeSetInt(ctx context.Context, entityID string, attributeKey string, attributeValue int64) error {
	attributeValueAsString := strconv.FormatInt(attributeValue, 10)
	return st.AttributeSetString(ctx, entityID, attributeKey, attributeValueAsString)
}

// AttributeSetFloat creates or updates a float attribute value for an entity
func (st *storeImplementation) AttributeSetFloat(ctx context.Context, entityID string, attributeKey string, attributeValue float64) error {
	attributeValueAsString := strconv.FormatFloat(attributeValue, 'f', 30, 64)
	return st.AttributeSetString(ctx, entityID, attributeKey, attributeValueAsString)
}

// AttributeGetString retrieves a string attribute value for an entity
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

// AttributeCreateWithKeyAndValue creates an attribute with the given key and value for an entity
func (st *storeImplementation) AttributeCreateWithKeyAndValue(ctx context.Context, entityID string, attributeKey string, attributeValue string) (AttributeInterface, error) {
	attr := NewAttribute()
	attr.SetEntityID(entityID)
	attr.SetKey(attributeKey)
	attr.SetValue(attributeValue)
	if err := st.AttributeCreate(ctx, attr); err != nil {
		return attr, err
	}
	return attr, nil
}

// AttributesSet creates or updates multiple attributes for an entity atomically.
// All writes are wrapped in a single database transaction — if any attribute
// fails, the entire operation is rolled back and no partial writes persist.
func (st *storeImplementation) AttributesSet(ctx context.Context, entityID string, attributes map[string]string) error {
	return st.db.Transaction(func(tx orm.Query) error {
		for key, value := range attributes {
			if key == "" {
				return errors.New("attribute key cannot be empty")
			}

			var rows []attributeRow
			err := tx.Table(st.attributeTableName).
				Where(COLUMN_ENTITY_ID+" = ?", entityID).
				Where(COLUMN_ATTRIBUTE_KEY+" = ?", key).
				Limit(1).
				Get(&rows)
			if err != nil {
				return err
			}

			now := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)

			if len(rows) > 0 {
				_, err := tx.Table(st.attributeTableName).
					Where(COLUMN_ID+" = ?", rows[0].ID).
					Update(map[string]any{
						COLUMN_ATTRIBUTE_VALUE: value,
						COLUMN_UPDATED_AT:      now,
					})
				if err != nil {
					return err
				}
			} else {
				err := tx.Table(st.attributeTableName).Create(map[string]any{
					COLUMN_ID:              GenerateShortID(),
					COLUMN_ENTITY_ID:       entityID,
					COLUMN_ATTRIBUTE_KEY:   key,
					COLUMN_ATTRIBUTE_VALUE: value,
					COLUMN_CREATED_AT:      now,
					COLUMN_UPDATED_AT:      now,
				})
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}
