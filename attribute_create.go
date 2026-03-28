package entitystore

import (
	"context"
	"errors"
	"log"

	"github.com/doug-martin/goqu/v9"
	"github.com/dromara/carbon/v2"
)

// AttributeCreate persists a new attribute record
func (st *storeImplementation) AttributeCreate(ctx context.Context, attr AttributeInterface) error {
	if attr == nil {
		return errors.New("attribute is required")
	}

	if attr.AttributeKey() == "" {
		return errors.New("attribute key is required field")
	}

	if attr.ID() == "" {
		attr.SetID(GenerateShortID())
	}

	if attr.CreatedAt() == "" {
		attr.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	if attr.UpdatedAt() == "" {
		attr.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	record := goqu.Record{}
	for k, v := range attr.Data() {
		record[k] = v
	}

	q := goqu.Dialect(st.dbDriverName).Insert(st.attributeTableName).Rows(record)
	sqlStr, _, _ := q.ToSQL()

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := st.database.Exec(ctx, sqlStr)
	if err != nil && st.GetDebug() {
		log.Println(err)
	}

	return err
}
