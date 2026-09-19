# Proposal: AttributeGroupBy — SQL-Level Aggregation over Attributes

**Date:** 2026-09-19
**Status:** IMPLEMENTED (`WithAggregate` on `AttributeQuery` +
`StoreInterface.AttributeGroupBy`; runnable demo in
`examples/aggregate/`)
**Author:** AI Assistant
**Supersedes:** parts of `2026-09-19-map-reduce-indexes.md` (kept as
the Go-side fallback for multi-attribute aggregation)

## Summary

Expose SQL `GROUP BY` over the attributes table. Because entitystore
is EAV — attribute keys/values are real rows, not opaque document
fields — the YesSql map/reduce problem mostly dissolves into a single
query the fluent API currently cannot express:

```sql
SELECT attribute_value, COUNT(*)
FROM attributes
WHERE attribute_key = 'publishedDay'
GROUP BY attribute_value
```

```go
// map[value]count — the BlogPostByDay equivalent, live and always correct
byDay, err := store.AttributeGroupBy(ctx, entitystore.AttributeQuery().
    WithAttributeKey("publishedDay").
    WithEntityType("blogpost").
    Aggregate(entitystore.AGGREGATE_COUNT))
// map[string]int{"2026-09-18": 3, "2026-09-19": 7}
```

## Before / After — What Do We Save?

### Today

```go
// Load everything, aggregate in Go
attrs, _ := store.AttributeList(ctx, entitystore.AttributeQuery().
    WithAttributeKey("publishedDay"))

byDay := map[string]int{}
for _, a := range attrs {
    byDay[a.GetValue()]++
}
```

Fine at small scale, but it pulls every matching row over the wire
to compute a number SQL can return directly.

### After

```go
byDay, _ := store.AttributeGroupBy(ctx, entitystore.AttributeQuery().
    WithAttributeKey("publishedDay").
    Aggregate(entitystore.AGGREGATE_COUNT))
```

### What is actually saved

| | Today | After |
|---|---|---|
| Rows transferred | N attribute rows | K distinct values |
| Lines per aggregate | fetch + hand-rolled loop | 1 call |
| Correctness | recomputed live (always correct) | recomputed live (always correct) |
| Write paths | untouched | untouched — nothing is persisted or maintained |

## Worked Example

```go
// Seed some blogposts
for _, d := range []string{"2026-09-17", "2026-09-18", "2026-09-18", "2026-09-19", "2026-09-19", "2026-09-19"} {
    post, _ := store.EntityCreate(ctx, "blogpost")
    post.AttributeSetString(ctx, "publishedDay", d)
}

// "How many posts per day?" — one call, one GROUP BY
byDay, _ := store.AttributeGroupBy(ctx, entitystore.AttributeQuery().
    WithAttributeKey("publishedDay").
    WithEntityType("blogpost").
    WithAggregate(entitystore.AGGREGATE_COUNT))
// map[string]any{"2026-09-17": 1, "2026-09-18": 2, "2026-09-19": 3}

// "How many posts in each status?"
byStatus, _ := store.AttributeGroupBy(ctx, entitystore.AttributeQuery().
    WithAttributeKey("status").
    WithEntityType("blogpost").
    WithAggregate(entitystore.AGGREGATE_COUNT))
// map[string]any{"draft": 2, "published": 4}

// "Which days have any posts at all?"
days, _ := store.AttributeGroupBy(ctx, entitystore.AttributeQuery().
    WithAttributeKey("publishedDay").
    WithAggregate(entitystore.AGGREGATE_DISTINCT))
// map keys only: {"2026-09-17", "2026-09-18", "2026-09-19"}

// Scope it the same way any attribute query is scoped
recentByDay, _ := store.AttributeGroupBy(ctx, entitystore.AttributeQuery().
    WithAttributeKey("publishedDay").
    WithCreatedAtGte("2026-09-01 00:00:00").
    WithAggregate(entitystore.AGGREGATE_COUNT))
```

SQL emitted (SQLite):

```sql
SELECT attribute_value, COUNT(*)
FROM attributes
WHERE attribute_key = 'publishedDay'
  AND entity_id IN (SELECT id FROM entities WHERE entity_type = 'blogpost')
GROUP BY attribute_value
```

## Proposed API

`AttributeQuery` gains one optional field:

```go
type AggregateFunc string

const (
    AGGREGATE_COUNT   AggregateFunc = "count"
    AGGREGATE_SUM     AggregateFunc = "sum"   // numeric cast per dialect
    AGGREGATE_MIN     AggregateFunc = "min"
    AGGREGATE_MAX     AggregateFunc = "max"
    AGGREGATE_DISTINCT AggregateFunc = "distinct" // just the keys, no value
)

// WithAggregate groups rows by attribute_value and applies fn.
// Result keys are the distinct attribute_values.
WithAggregate(fn AggregateFunc)
```

And one store method:

```go
// AttributeGroupBy runs the attribute query with GROUP BY on
// attribute_value and returns map[attribute_value]result.
// Result type depends on the aggregate: int for count/sum,
// string for min/max, presence-only for distinct.
AttributeGroupBy(ctx context.Context, query AttributeQueryInterface) (map[string]any, error)
```

(`map[string]any` is the honest shape for mixed aggregates; if only
`COUNT` ships in v1, `map[string]int` is nicer — open detail.)

### Why this covers the YesSql sample

| YesSql | This proposal |
|---|---|
| `BlogPostByAuthor` index (lookup by author) | `EntityListByAttribute("blogpost","author",x)` — already exists |
| `BlogPostByDay` (count per day) | `AttributeGroupBy` + `AGGREGATE_COUNT` |
| `Delete` (count hits zero → row gone) | n/a — computed live, absent keys simply don't appear |

## Design Notes

- **Fits existing conventions:** `WithAggregate` is a query field
  like any other — `Get`/`Has`/`Validate`, applied uniformly; empty
  query still valid. Grouping happens in
  `applyAttributeFilters`-adjacent code, single place.
- **Dialect parity:** `COUNT`/`MIN`/`MAX` are portable via goqu;
  `SUM` needs a `CAST(attribute_value AS numeric)` which differs
  slightly across SQLite/MySQL/Postgres — same problem class as
  `escapeLike`, solved once in `query_filters.go`.
- **The EAV limit remains:** grouping by one attribute while
  aggregating *another* ("sum `amount` by `sku`") still needs a
  self-join or Go-side reduce — that's what the other proposal's
  `ReduceBy` is for. The two compose: SQL `GroupBy` for the common
  case, `ReduceBy` for the rest.
- **No persistence, no drift:** like the lo-style version, results
  are computed live — but in the database, not in memory.

## Alternatives Considered

1. **Persisted projections (YesSql model)** — registration machinery
   and hidden write-path behaviour to maintain tables that this one
   query makes mostly unnecessary at entitystore's scale.
2. **Go-side `ReduceBy` only** — works but transfers N rows to
   compute a number; SQL `GROUP BY` is strictly better when the
   grouping key is a single attribute value.
3. **Raw SQL in consumers** — already possible, but bypasses the
   query abstraction and re-solves dialect escaping per app.
4. **`dracory/indexstore`** — a place to *persist* aggregates if a
   hot query ever needs it; orthogonal, composes with either
   proposal.

## Cost

- One `AggregateFunc` type + `WithAggregate` on `AttributeQuery`
  (fluent + `Has`/`Get`/`Validate`).
- One `AttributeGroupBy` store method (~40 lines incl. dialect cast).
- Tests: count/min/max/distinct grouping, entity-type scoping,
  validation of unknown aggregate names.
