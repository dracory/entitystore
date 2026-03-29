package entitystore

import "github.com/doug-martin/goqu/v9"

// EntityQueryOptions provides filtering and pagination options for entity queries
type EntityQueryOptions struct {
	ID           string   // Filter by specific entity ID
	IDs          []string // Filter by multiple entity IDs
	EntityType   string   // Filter by entity type
	EntityHandle string   // Filter by entity handle
	Limit        uint64   // Maximum number of results to return
	Offset       uint64   // Number of results to skip
	Search       string   // Text search (not implemented yet)
	SortBy       string   // Column to sort by (default: id)
	SortOrder    string   // Sort direction: "asc" or "desc"
	CountOnly    bool     // Return only count, not results
}

// EntityQuery builds a goqu query for entities based on the provided options
// Returns a SelectDataset that can be further customized or executed
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
