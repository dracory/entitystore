package entitystore

import (
	"context"
	"database/sql"
	"log"

	"github.com/doug-martin/goqu/v9"
)

// EntityFindByAttribute finds an entity by attribute
func (st *storeImplementation) EntityFindByAttribute(ctx context.Context, entityType string, attributeKey string, attributeValue string) (*Entity, error) {
	q := goqu.Dialect(st.dbDriverName).From(st.attributeTableName)
	q = q.LeftJoin(goqu.I(st.entityTableName), goqu.On(goqu.Ex{st.attributeTableName + "." + COLUMN_ENTITY_ID: goqu.I(st.entityTableName + "." + COLUMN_ID)}))
	q = q.Where(goqu.C(COLUMN_ENTITY_TYPE).Eq(entityType))
	q = q.Where(goqu.And(goqu.C(COLUMN_ATTRIBUTE_KEY).Eq(attributeKey), goqu.C(COLUMN_ATTRIBUTE_VALUE).Eq(attributeValue)))
	q = q.Select(COLUMN_ENTITY_ID)

	sqlStr, _, _ := q.ToSQL()
	if st.GetDebug() {
		log.Println(sqlStr)
	}

	var entityID string
	err := st.database.DB().QueryRow(sqlStr).Scan(&entityID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		if st.GetDebug() {
			log.Println(err)
		}
		return nil, err
	}

	return st.EntityFindByID(ctx, entityID)
}
