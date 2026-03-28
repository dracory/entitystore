package entitystore

import (
	"context"
	"log"

	"github.com/doug-martin/goqu/v9"
	"github.com/dromara/carbon/v2"
)

// EntityUpdate persists changes to an existing entity record
func (st *storeImplementation) EntityUpdate(ctx context.Context, entity EntityInterface) error {
	entity.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	record := goqu.Record{}
	for k, v := range entity.Data() {
		record[k] = v
	}

	q := goqu.Dialect(st.dbDriverName).
		Update(st.GetEntityTableName()).
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
