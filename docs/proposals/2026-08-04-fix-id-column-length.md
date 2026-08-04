# Bug Fix: ID Column Length Mismatch

**Date:** 2026-08-04
**Status:** FIXED
**Severity:** High — causes insert failures on MySQL/MariaDB

## Problem

`GenerateShortID()` produces IDs up to ~13 characters (Crockford Base32 encoding of a 19-digit nanosecond timestamp), but `MigrateUp` creates all `id` (and `entity_id`) columns as `VARCHAR(9)`.

On MySQL/MariaDB this causes:

```
Error 1406 (22001): Data too long for column 'id' at row 1
```

Every `EntityCreateWithType`, `EntityCreateWithTypeAndAttributes`, `AttributeCreate`, etc. fails because the generated ID exceeds the column width.

## Root Cause

`store.go` — all 33 occurrences of hardcoded `9` in id-type column definitions:
```go
table.String(COLUMN_ID, 9)
```

`id_helpers.go`:
```go
func GenerateShortID() string {
    timestampID := uid.TimestampNano()          // 19-digit number, e.g. "1725111153123456789"
    shortened, _ := uid.ShortenCrockford(timestampID)
    return strings.ToLower(shortened)           // ~13 chars, e.g. "1hj3m05wq7je0"
}
```

The docs (`architecture.md` line 152) said "9 characters by default" but the actual output is ~13 chars.

## Why It Wasn't Caught

All tests run on **SQLite**, which uses type affinity and does **not** enforce `VARCHAR(n)` length limits. A `VARCHAR(9)` column in SQLite silently accepts a 13-character string with no error. The bug only manifests on databases that enforce column constraints:

- **MySQL/MariaDB:** `Error 1406 (22001): Data too long for column`
- **PostgreSQL:** `value too long for type character varying(9)`

There were no integration tests against MySQL/MariaDB/PostgreSQL — only SQLite. Additionally, the column width (`9`) was hardcoded separately from the ID generation logic, with no constant or test linking the two, so the mismatch was invisible.

## Fix Applied

Introduced `ID_COLUMN_LENGTH` constant in `consts.go` and replaced all hardcoded `9` values in `store.go` schema definitions with this constant. Set to **40** to cover:

- `GenerateShortID()` output (13 chars) — the library's internal ID generator
- `ShortenID()` output from 32-char UUIDs (21 chars) — the library's UUID shortening utility
- Standard UUIDs that users may set directly (36 chars) — see below
- Small margin for future growth

### Why 40 and not 21

The library auto-generates IDs via `GenerateShortID()` when `entity.ID() == ""`, but nothing prevents a user from calling `entity.SetID("any-string")` before create. The `*Create` methods only auto-generate if the ID is empty — user-provided IDs are stored as-is with no validation or length check. This means the column width acts as an implicit constraint on user-provided IDs.

21 would cover the library's own generators but would break if a user passes a standard 36-char UUID (e.g. `550e8400-e29b-41d4-a716-446655440000`). 40 covers all three ID paths (internal generator, UUID shortener, user-provided UUIDs) with a small margin, without being wastefully wide like 255.

### Files changed

- `consts.go` — added `ID_COLUMN_LENGTH = 40` constant
- `store.go` — replaced all 33 occurrences of hardcoded `9` with `ID_COLUMN_LENGTH`
- `id_column_length_test.go` — regression tests proving the fix

### Migration for existing databases

Existing databases created with `VARCHAR(9)` need an ALTER TABLE migration. Downstream applications (e.g. CourseThread) have a temporary application-level migration (`2026_08_04_0002_entity_store_widen_id_columns`) that should be removed once this fix is published.

## Impact

- **MySQL/MariaDB:** Broken — inserts fail. This fix is required.
- **SQLite:** Not affected — SQLite uses type affinity and does not enforce VARCHAR length.
- **PostgreSQL:** Likely affected — enforces VARCHAR length.

## Workaround (for existing databases)

Downstream applications with existing `VARCHAR(9)` databases can ALTER the columns manually:

```sql
ALTER TABLE snv_entities_entity MODIFY COLUMN id VARCHAR(40) NOT NULL;
ALTER TABLE snv_entities_entity_trash MODIFY COLUMN id VARCHAR(40) NOT NULL;
ALTER TABLE snv_entities_attribute MODIFY COLUMN id VARCHAR(40) NOT NULL;
ALTER TABLE snv_entities_attribute_trash MODIFY COLUMN id VARCHAR(40) NOT NULL;
ALTER TABLE snv_entities_attribute MODIFY COLUMN entity_id VARCHAR(40) NOT NULL;
ALTER TABLE snv_entities_attribute_trash MODIFY COLUMN entity_id VARCHAR(40) NOT NULL;
```
