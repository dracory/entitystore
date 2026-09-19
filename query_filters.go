package entitystore

import (
	"strings"

	"github.com/dracory/neat/contracts/database/orm"
)

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

// escapeLike escapes LIKE wildcards in a literal so it matches verbatim.
// Paired with an explicit ESCAPE '\' clause — the default escape varies
// across drivers (MySQL backslash, PostgreSQL standard, SQLite none).
func escapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

// applyLike adds a LIKE clause with an explicit ESCAPE so escaped
// wildcards behave identically across SQLite, MySQL, and PostgreSQL.
// The pattern argument is already wildcard-processed by the caller
// (raw patterns pass through untouched; literals go through escapeLike).
func applyLike(q orm.Query, column string, pattern string) orm.Query {
	if pattern == "" {
		return q
	}
	return q.Where(column+` LIKE ? ESCAPE '\'`, pattern)
}
