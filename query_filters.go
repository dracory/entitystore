package entitystore

import "github.com/dracory/neat/contracts/database/orm"

// applyTimeRange adds inclusive lower/upper bound clauses on a datetime
// column ("YYYY-MM-DD HH:MM:SS" UTC strings compare lexicographically).
// Shared by all applyXFilters implementations to prevent filter drift.
func applyTimeRange(q orm.Query, column string, gte string, lte string) orm.Query {
	if gte != "" {
		q = q.Where(column+" >= ?", gte)
	}
	if lte != "" {
		q = q.Where(column+" <= ?", lte)
	}
	return q
}
