package entitystore

import "github.com/doug-martin/goqu/v9"

// AttributeQueryOptions provides filtering and pagination options for attribute queries
type AttributeQueryOptions struct {
	ID           string   // Filter by specific attribute ID
	IDs          []string // Filter by multiple attribute IDs
	EntityID     string   // Filter by associated entity ID
	EntityType   string   // Filter by entity type (requires EntityHandle or join)
	EntityHandle string   // Filter by entity handle (requires EntityType or join)
	AttributeKey string   // Filter by attribute key/name
	Limit        uint64   // Maximum number of results to return
	Offset       uint64   // Number of results to skip
	SortBy       string   // Column to sort by (default: id)
	SortOrder    string   // Sort direction: "asc" or "desc"
	CountOnly    bool     // Return only count, not results
}

// AttributeQuery builds a goqu query for attributes based on the provided options
// Returns a SelectDataset that can be further customized or executed
func (st *storeImplementation) AttributeQuery(options AttributeQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(st.dbDriverName).From(st.attributeTableName)

	if options.EntityType != "" && options.EntityHandle != "" {
		q = q.LeftJoin(goqu.I(st.entityTableName), goqu.On(goqu.Ex{st.attributeTableName + ".entity_id": goqu.I(st.entityTableName + ".id")}))
	}

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

	if options.EntityID != "" {
		q = q.Where(goqu.C(COLUMN_ENTITY_ID).Eq(options.EntityID))
	}

	if options.EntityType != "" && options.EntityHandle != "" {
		q = q.Where(goqu.C(COLUMN_ENTITY_TYPE).Eq(options.EntityType))
		q = q.Where(goqu.C(COLUMN_ENTITY_HANDLE).Eq(options.EntityHandle))
	}

	if options.AttributeKey != "" {
		q = q.Where(goqu.C(COLUMN_ATTRIBUTE_KEY).Eq(options.AttributeKey))
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
