# Proposal: Generic Map/Reduce Helpers over Entity Queries

**Date:** 2026-09-19
**Status:** PROPOSED (Go-side fallback — see
`2026-09-19-attribute-group-by.md` for the preferred SQL-level
approach covering the common single-attribute case)
**Author:** AI Assistant

## Summary

Add small generic helpers — in the style of `samber/lo` — that map,
group, and reduce entity query results. Read-side computation only:
no registration, no write-path hooks, no hidden machinery. What you
call is what happens.

```go
posts, _ := store.EntityList(ctx, entitystore.EntityQuery().
    WithEntityType("blogpost"))

// Map — transform entities to anything
authors := entitystore.Map(posts, func(e entitystore.EntityInterface) string {
    a, _ := e.AttributeGetString(ctx, "author")
    return a
})

// GroupBy — bucket entities by a key, like lo.GroupBy
byDay := entitystore.GroupBy(posts, func(e entitystore.EntityInterface) string {
    d, _ := e.AttributeGetString(ctx, "publishedDay")
    return d
})
// map[string][]EntityInterface — {"2026-09-18": [post, post, post], ...}

// ReduceBy — group + fold in one step, the YesSql BlogPostByDay equivalent
counts := entitystore.ReduceBy(posts,
    func(e entitystore.EntityInterface) string { // key
        d, _ := e.AttributeGetString(ctx, "publishedDay")
        return d
    },
    func(count int, _ entitystore.EntityInterface) int { // fold
        return count + 1
    },
    0) // initial value per key
// map[string]int — {"2026-09-18": 3, "2026-09-19": 7, ...}
```

## Before / After — What Do We Save?

### Today

```go
posts, _ := store.EntityList(ctx, entitystore.EntityQuery().
    WithEntityType("blogpost"))

byDay := map[string]int{}
for _, p := range posts {
    d, _ := p.AttributeGetString(ctx, "publishedDay")
    byDay[d]++
}
```

Works, but every consumer re-derives the group-then-fold loop, and
the two-step shape (key extraction vs. aggregation) is invisible in
the boilerplate.

### After

```go
byDay := entitystore.ReduceBy(posts, keyFunc, foldFunc, 0)
```

### What is actually saved

| | Today | After |
|---|---|---|
| Lines per aggregate | 5–8 lines of hand-rolled loop | 1 expression |
| Correctness | caller must get bucketing right | library function |
| Write paths touched | none | none — identical |
| Persistence | none — computed per call | none — computed per call |

This proposal deliberately does **not** persist results. The
aggregate is computed from live data on every call, so it can never
drift, never needs reindexing, and never touches the write path.

## Proposed API

Four generic functions, ~10 lines each:

```go
func Map[T any](entities []EntityInterface, fn func(EntityInterface) T) []T

func GroupBy[K comparable](entities []EntityInterface, key func(EntityInterface) K) map[K][]EntityInterface

func Reduce[T any](entities []EntityInterface, fold func(T, EntityInterface) T, initial T) T

func ReduceBy[K comparable, T any](entities []EntityInterface,
    key func(EntityInterface) K,
    fold func(T, EntityInterface) T,
    initial T) map[K]T
```

`ReduceBy` is the workhorse — it is literally `GroupBy` followed by a
fold per bucket:

```go
func ReduceBy[K comparable, T any](entities []EntityInterface,
    key func(EntityInterface) K,
    fold func(T, EntityInterface) T,
    initial T) map[K]T {
    out := map[K]T{}
    for _, e := range entities {
        k := key(e)
        out[k] = fold(out[k], e) // zero value of T starts the fold
    }
    return out
}
```

(When `out[k]` doesn't exist yet, `fold` receives the zero value —
the `initial` parameter is applied to absent keys, or dropped in
favour of zero-value semantics; open detail.)

## How it maps to YesSql

| YesSql | This proposal |
|---|---|
| `Map(blogPost => ...)` | `Map(entities, fn)` or the `key` func |
| `.Group(blogPost => blogPost.Day)` | `GroupBy` / the `key` func |
| `.Reduce(group => new BlogPostByDay{Count = group.Sum(...)})` | the `fold` func |
| `.Delete(...)` | not needed — nothing is persisted, so nothing is maintained |
| index row type (`BlogPostByDay` struct) | your Go type `T` — real types, not `map[string]string` |
| index table, updated on save | a `map[K]T` returned to the caller |

## Design Notes

- **Pure functions over `[]EntityInterface`** — they don't even need
  the store, so they work on results of `EntityList`,
  `EntityListByAttribute`, or any filtered slice.
- **One extra query cost:** `AttributeGetString` inside a `Map`/`key`
  func may hit the store per entity. An `AttributeListGrouped`-style
  prefetch (separate proposal) removes the N+1 — the two compose
  naturally.
- **Scale ceiling:** recompute-on-read is correct and simple up to
  entitystore's documented target (<100K entities). If a hot
  aggregate ever outgrows it, a persisted variant can be layered on
  later without changing these signatures.

## Alternatives Considered

1. **Persisted write-path projections (YesSql model)** — index rows
   maintained inside `EntityCreate`/`Update`/`Trash` transactions.
   Always-consistent reads, but adds registration machinery,
   `_index:` entity types, reindex tooling, and hidden behaviour on
   writes — a lot of mechanism for a store whose consumers are small.
2. **`samber/lo` directly** — `lo.GroupBy`/`lo.Map` work on
   `[]EntityInterface` today, zero new code. But they'd live outside
   the package, the `ctx`-ful attribute access is still hand-written,
   and a first-party helper documents the intended pattern.
3. **Consumer-side loops (status quo)** — five lines each, everywhere.
4. **SQL `GROUP BY`** — can't express it: attribute keys are data,
   not columns; folding requires one attribute join per key.

## Cost

Four generic functions (~40 lines total) plus tests. No schema
changes, no new interfaces, no store changes, no write-path hooks.
Purely additive.
