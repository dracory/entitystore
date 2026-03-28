package entitystore

import (
	"context"
	"log"
)

// AttributeList lists attributes matching the given query options
func (st *storeImplementation) AttributeList(ctx context.Context, options AttributeQueryOptions) ([]AttributeInterface, error) {
	q := st.AttributeQuery(options)

	sqlStr, _, errSql := q.ToSQL()
	if errSql != nil {
		return nil, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	attributeMaps, err := st.database.SelectToMapString(ctx, sqlStr)
	if err != nil {
		return nil, err
	}

	var list []AttributeInterface
	for _, m := range attributeMaps {
		list = append(list, NewAttributeFromExistingData(m))
	}

	return list, nil
}
