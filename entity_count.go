package entitystore

import (
	"context"
	"database/sql"
	"log"

	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/sqlscan"
)

// EntityCount counts the entities of a specified type
// EntityCount counts entities
func (st *storeImplementation) EntityCount(ctx context.Context, options EntityQueryOptions) (int64, error) {
	options.CountOnly = true

	q := st.EntityQuery(options)
	sqlStr, _, errSql := q.Limit(1).Select(goqu.COUNT(goqu.Star()).As("count")).ToSQL()

	if errSql != nil {
		return 0, errSql
	}

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	type countResult struct {
		Count int64 `db:"count"`
	}

	var result countResult
	err := sqlscan.Get(ctx, st.database.DB(), &result, sqlStr)
	if err != nil {
		if err == sql.ErrNoRows {
			// sqlscan does not use this anymore
			return 0, err
		}

		if sqlscan.NotFound(err) {
			return 0, nil
		}

		return 0, err
	}

	return result.Count, nil
}
