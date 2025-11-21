package entitystore

import "github.com/doug-martin/goqu/v9"

type EntityQueryOptions struct {
	ID           string
	IDs          []string
	EntityType   string
	EntityHandle string
	Limit        uint64
	Offset       uint64
	Search       string
	SortBy       string
	SortOrder    string // asc / dec
	CountOnly    bool
}

func (st *storeImplementation) EntityQuery(options EntityQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(st.dbDriverName).From(st.entityTableName)

	if len(options.IDs) > 0 {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDs))
	}

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	sortByColumn := COLUMN_ID
	sortOrder := "asc"

	if options.SortOrder != "" {
		sortOrder = options.SortOrder
	}

	if options.SortBy != "" {
		sortByColumn = options.SortBy
	}

	if sortOrder == "asc" {
		q = q.Order(goqu.I(sortByColumn).Asc())
	} else {
		q = q.Order(goqu.I(sortByColumn).Desc())
	}

	if options.EntityType != "" {
		q = q.Where(goqu.C(COLUMN_ENTITY_TYPE).Eq(options.EntityType))
	}

	if options.EntityHandle != "" {
		q = q.Where(goqu.C(COLUMN_ENTITY_HANDLE).Eq(options.EntityHandle))
	}

	q = q.Offset(uint(options.Offset))

	if options.Limit != 0 {
		q = q.Limit(uint(options.Limit))
	}

	if !options.CountOnly {
		if options.Limit > 0 {
			q = q.Limit(uint(options.Limit))
		}

		if options.Offset > 0 {
			q = q.Offset(uint(options.Offset))
		}
	}

	return q.Select()
}
