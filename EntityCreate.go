package entitystore

import (
	"errors"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/uid"
)

// EntityCreate creates a new entity
func (st *storeImplementation) EntityCreate(entity *Entity) error {
	if entity == nil {
		return errors.New("entity cannot be nil")
	}

	if entity.ID() == "" {
		entity.SetID(uid.HumanUid())
	}

	if entity.CreatedAt().IsZero() {
		entity.SetCreatedAt(time.Now())
	}

	if entity.UpdatedAt().IsZero() {
		entity.SetUpdatedAt(time.Now())
	}

	q := goqu.Dialect(st.dbDriverName).Insert(st.entityTableName)
	q = q.Rows(entity.ToMap())

	sqlStr, _, errSql := q.ToSQL()

	if errSql != nil {
		return errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := st.database.Exec(sqlStr)

	if err != nil {
		return err
	}

	return nil
}
