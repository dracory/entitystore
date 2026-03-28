package entitystore

import (
	"context"
	"log"
)

// EntityList lists entities matching the given query options
func (st *storeImplementation) EntityList(ctx context.Context, options EntityQueryOptions) ([]EntityInterface, error) {
	q := st.EntityQuery(options)

	sqlStr, _, errSql := q.ToSQL()
	if errSql != nil {
		return nil, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	entityMaps, err := st.database.SelectToMapString(ctx, sqlStr)
	if err != nil {
		log.Println("EntityList error:", err)
		return nil, err
	}

	var list []EntityInterface
	for _, m := range entityMaps {
		list = append(list, NewEntityFromExistingData(m))
	}

	return list, nil
}
