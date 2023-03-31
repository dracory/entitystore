package entitystore

import (
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
)

// AttributeUpdate updates an attribute
func (st *Store) AttributeUpdate(attr Attribute) error {
	attr.SetUpdatedAt(time.Now())

	q := goqu.Dialect(st.dbDriverName).Update(st.attributeTableName)
	q = q.Where(goqu.C("id").Eq(attr.ID()))
	q = q.Set(attr.ToMap())

	sqlStr, _, errSql := q.ToSQL()

	if errSql != nil {
		return errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := st.database.Exec(sqlStr)

	if err != nil {
		if st.GetDebug() {
			log.Println(err)
		}

		return err
	}

	return nil
}
