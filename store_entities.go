package entitystore

import (
	"context"
	"errors"
	"log"

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

	if entity.CreatedAt() == "" {
		entity.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	if entity.UpdatedAt() == "" {
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
