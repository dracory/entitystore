package entitystore

import (
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/gouniverse/uid"
)

// EntityCreate creates a new entity
func (st *Store) EntityCreate(entityType string) (*Entity, error) {
	entity := st.NewEntity(NewEntityOptions{
		ID:        uid.HumanUid(),
		Type:      entityType,
		Handle:    "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	q := goqu.Dialect(st.dbDriverName).Insert(st.entityTableName)
	q = q.Rows(entity.ToMap())
	sqlStr, _, errSql := q.ToSQL()

	if errSql != nil {
		return nil, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := st.database.Exec(sqlStr)

	if err != nil {
		return entity, err
	}

	return entity, nil
}
