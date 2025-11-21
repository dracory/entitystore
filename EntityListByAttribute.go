package entitystore

import (
	"context"
	"log"

	"github.com/doug-martin/goqu/v9"
)

// EntityListByAttribute finds an entity by attribute
func (st *storeImplementation) EntityListByAttribute(ctx context.Context, entityType string, attributeKey string, attributeValue string) (entityList []Entity, err error) {
	var entityIDs []string

	q := goqu.Dialect(st.dbDriverName).From(st.attributeTableName).
		LeftJoin(goqu.I(st.entityTableName), goqu.On(goqu.Ex{st.attributeTableName + "." + COLUMN_ENTITY_ID: goqu.I(st.entityTableName + "." + COLUMN_ID)})).
		Where(goqu.C(COLUMN_ENTITY_TYPE).Eq(entityType)).
		Where(goqu.And(goqu.C(COLUMN_ATTRIBUTE_KEY).Eq(attributeKey), goqu.C(COLUMN_ATTRIBUTE_VALUE).Eq(attributeValue))).
		Select(COLUMN_ENTITY_ID)

	sqlStr, _, err := q.ToSQL()

	if err != nil {
		if st.GetDebug() {
			log.Println(err.Error())
		}
		return nil, err
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	rows, err := st.database.Query(ctx, sqlStr)

	if err != nil {
		return []Entity{}, err
	}

	for rows.Next() {
		var entityID string
		err := rows.Scan(&entityID)
		if err != nil {
			return []Entity{}, err
		}
		entityIDs = append(entityIDs, entityID)
	}

	if len(entityIDs) < 1 {
		return entityList, nil
	}

	return st.EntityList(ctx, EntityQueryOptions{
		EntityType: entityType,
		IDs:        entityIDs,
		SortBy:     COLUMN_ID,
	})
}
