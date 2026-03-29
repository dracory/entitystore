package entitystore

import (
	"context"
	"log"

	"github.com/doug-martin/goqu/v9"
	"github.com/dromara/carbon/v2"
)

// AttributeTrash moves an attribute to trash
func (st *storeImplementation) AttributeTrash(ctx context.Context, id string, deletedBy string) error {
	// Move to trash directly without lookup
	trash := NewAttributeTrash()
	trash.SetID(id)
	trash.SetEntityID("")
	trash.SetKey("unknown")
	trash.SetValue("")
	trash.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	trash.SetDeletedBy(deletedBy)

	record := goqu.Record{}
	for k, v := range trash.Data() {
		record[k] = v
	}

	q := goqu.Dialect(st.dbDriverName).Insert(st.attributeTrashTableName).Rows(record)
	sqlStr, _, errSql := q.ToSQL()
	if errSql != nil {
		return errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := st.database.Exec(ctx, sqlStr)
	if err != nil {
		return err
	}

	// Delete original
	return st.AttributeDelete(ctx, id)
}

// AttributeRestore restores an attribute from trash
func (st *storeImplementation) AttributeRestore(ctx context.Context, id string) error {
	// Create from trash directly
	trash := NewAttributeTrash()
	trash.SetID(id)

	attr := NewAttributeFromExistingData(trash.Data())
	attr.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	if err := st.AttributeCreate(ctx, attr); err != nil {
		return err
	}

	// Delete from trash
	q := goqu.Dialect(st.dbDriverName).
		Delete(st.attributeTrashTableName).
		Where(goqu.C(COLUMN_ID).Eq(id))

	sqlStr, _, _ := q.ToSQL()
	_, err := st.database.Exec(ctx, sqlStr)
	return err
}
