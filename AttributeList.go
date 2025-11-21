package entitystore

import (
	"context"
	"log"
)

// AttributeList lists attributes
func (st *storeImplementation) AttributeList(ctx context.Context, options AttributeQueryOptions) (attributeList []Attribute, err error) {
	q := st.AttributeQuery(options)

	sqlStr, _, errSql := q.ToSQL()

	if errSql != nil {
		return attributeList, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	attributeMaps, errSelect := st.database.SelectToMapString(ctx, sqlStr)

	if errSelect != nil {
		return nil, err
	}

	// attributeMaps := []map[string]string{}
	// errScan := sqlscan.Select(context.Background(), st.db, &attributeMaps, sqlStr)
	// if errScan != nil {
	// 	if errScan == sql.ErrNoRows {
	// 		// sqlscan does not use this anymore
	// 		return nil, errScan
	// 	}

	// 	if sqlscan.NotFound(errScan) {
	// 		return nil, nil
	// 	}

	// 	return nil, err
	// }

	for i := 0; i < len(attributeMaps); i++ {
		attribute := st.NewAttributeFromMap(attributeMaps[i])
		attributeList = append(attributeList, attribute)
	}

	return attributeList, nil
}
