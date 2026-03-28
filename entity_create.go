package entitystore

import (
	"context"
	"errors"
	"log"

	"github.com/doug-martin/goqu/v9"
)

// EntityCreate creates a new entity
func (st *storeImplementation) EntityCreate(ctx context.Context, entity *Entity) error {
	if entity == nil {
		return errors.New("entity cannot be nil")
	}

	if entity.ID() == "" {
		entity.SetID(GenerateShortID())
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

	_, err := st.database.Exec(ctx, sqlStr)

	if err != nil {
		return err
	}

	return nil
}
