package entitystore

import (
	"context"
	"log"
)

// EntityList lists entities
func (st *storeImplementation) EntityList(ctx context.Context, options EntityQueryOptions) (entityList []Entity, err error) {
	q := st.EntityQuery(options)

	sqlStr, _, errSql := q.ToSQL()

	if errSql != nil {
		return entityList, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	entityMaps, errSelect := st.database.SelectToMapString(ctx, sqlStr)
	// errScan := sqlscan.Select(context.Background(), st.db, &entityMaps, sqlStr)
	// if errScan != nil {
	// 	if errScan == sql.ErrNoRows {
	// 		// sqlscan does not use this anymore
	// 		return nil, errScan
	// 	}

	// 	if sqlscan.NotFound(errScan) {
	// 		return nil, nil
	// 	}

	// 	log.Println("EntityList. Error: ", errScan.Error())
	// 	return nil, err
	// }

	if errSelect != nil {
		log.Println("EntityList. Error: ", errSelect.Error())
		return nil, errSelect
	}

	for i := 0; i < len(entityMaps); i++ {
		entity := st.NewEntityFromMap(entityMaps[i])
		entityList = append(entityList, entity)
	}

	return entityList, nil
}
