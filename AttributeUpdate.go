package entitystore

import (
	"context"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
)

// AttributeUpdate updates an attribute
func (st *storeImplementation) AttributeUpdate(attr Attribute) error {
	attr.SetUpdatedAt(time.Now())

	q := goqu.Dialect(st.dbDriverName).Update(st.attributeTableName)
	q = q.Where(goqu.C(COLUMN_ID).Eq(attr.ID()))
	q = q.Set(attr.ToMap())

	sqlStr, _, errSql := q.ToSQL()

	if errSql != nil {
		return errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := st.database.Exec(context.Background(), sqlStr)

	if err != nil {
		if st.GetDebug() {
			log.Println(err)
		}

		return err
	}

	return nil
}
