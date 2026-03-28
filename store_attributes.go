package entitystore

import (
	"context"
	"errors"
	"log"

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

	if attribute.CreatedAt() == "" {
		attribute.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	if attribute.UpdatedAt() == "" {
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
