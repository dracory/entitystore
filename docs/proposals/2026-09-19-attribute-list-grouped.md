# Proposal: AttributeListGrouped — EAV Pivot Helper

**Date:** 2026-09-19
**Status:** PROPOSED
**Author:** AI Assistant

## Summary

Add a convenience method that runs an attribute query and returns the
results pre-grouped by entity — the in-memory equivalent of an EAV
"PIVOT", where attribute keys behave as the logical columns of an entity
row.

```go
// map[entityID]map[attributeKey]attributeValue
grouped, err := store.AttributeListGrouped(ctx, AttributeQuery().
    WithEntityType("character").
    WithAttributeKeys([]string{"name", "title", "avatar"}))
```

## Motivation

In the EAV model, `attribute_key` plays the role of a column and
`attribute_value` the role of that column's value for a given entity.
SQL cannot pivot rows into columns dynamically without hardcoding the
key list, so the standard pattern is:

1. Fetch the attribute rows you need (`WHERE attribute_key IN (...)` —
   already supported via `WithAttributeKeys`)
2. Group them by `entity_id` in memory

Every consumer that wants entity-shaped data currently re-implements
step 2 by hand:

```go
// Repeated in every consumer today
byEntity := map[string][]entitystore.AttributeInterface{}
for _, a := range attrs {
    byEntity[a.GetEntityID()] = append(byEntity[a.GetEntityID()], a)
}
```

This is boilerplate, easy to get subtly wrong (e.g. losing attributes
when an entity has no rows), and — critically — is the *desired result
shape* for most real call sites, not the flat row list.

## Proposed API

```go
// AttributeListGrouped retrieves attributes matching the query and
// groups them by entity ID: map[entityID]map[attributeKey]attributeValue.
// Entities with no matching attributes are absent from the result.
// Returns an error if the query is nil or fails validation.
AttributeListGrouped(ctx context.Context, query AttributeQueryInterface) (map[string]map[string]string, error)
```

Implementation is a thin wrapper over `AttributeList` — same query,
same validation, same single round trip:

```go
func (st *storeImplementation) AttributeListGrouped(ctx context.Context, query AttributeQueryInterface) (map[string]map[string]string, error) {
    attrs, err := st.AttributeList(ctx, query)
    if err != nil {
        return nil, err
    }
    grouped := map[string]map[string]string{}
    for _, a := range attrs {
        if grouped[a.GetEntityID()] == nil {
            grouped[a.GetEntityID()] = map[string]string{}
        }
        grouped[a.GetEntityID()][a.GetKey()] = a.GetValue()
    }
    return grouped, nil
}
```

## Design Notes

- **Return `map[string]map[string]string`, not objects.** The value is
  the convenience; callers that need IDs/timestamps per attribute can
  still use `AttributeList`. A richer shape
  (`map[entityID][]AttributeInterface`) keeps the nested-loop
  boilerplate the proposal exists to remove.
- **Duplicate keys:** last write wins, matching the read order of the
  underlying query. (Attribute keys are effectively unique per entity in
  practice — `AttributeSet*` upserts — so collisions are an edge case.)
- **Entities with zero matching attributes** are simply absent;
  callers check `grouped[id]` for nil.
- **Validation:** inherits `AttributeList`'s `Validate()`/`nil` checks
  for free.

## Alternatives Considered

1. **`Columns`/`Select` parameter on the query** (cmsstore
   `SetColumns`-style) — rejected: it returns hollow `AttributeInterface`
   objects with silently empty fields, breaking the "fully populated
   object" invariant, and the payload being skipped is the value column
   callers actually need.
2. **True SQL PIVOT** — impossible dynamically; key names would have to
   be hardcoded in the query, defeating the EAV purpose.
3. **Consumer-side grouping (status quo)** — works, but every caller
   duplicates the same five lines; a canonical helper removes drift.

## Cost

One method, ~15 lines, plus a test. No schema changes, no breaking
changes — purely additive to `StoreInterface`.
