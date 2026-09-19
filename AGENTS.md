# AGENTS.md

Entity Store (`github.com/dracory/entitystore`) — a schemaless entity-attribute-value store for Go, backed by SQL (SQLite, MySQL, PostgreSQL) via `goqu`.

## Commands

```bash
go build ./...     # build
go vet ./...       # vet
go test ./...      # all tests (root package + examples + testutils)
gofmt -l .         # formatting check
```

Requires Go 1.27+ (see `go.mod`). Tests use in-memory SQLite (`InitDB` helper, `modernc.org/sqlite`).

## Layout

- `interfaces.go` — all public interfaces (`EntityInterface`, `AttributeInterface`, `RelationshipInterface`, `TaxonomyInterface`, `TaxonomyTermInterface`, `EntityTaxonomyInterface`, query interfaces, `StoreInterface`)
- `*_query.go` — fluent query builders (one per domain)
- `query_filters.go` — shared SQL filter helpers
- `store_*.go` — store method implementations (one file per domain; `*_trash.go` for soft-delete ops)
- `*_query_test.go` — query builder tests (fluent usage, presence flags, `Validate()` table tests, real-store integration)
- `examples/` — runnable examples (basic, relationships, taxonomy) each with `main_test.go`
- `docs/` — user-facing documentation; `docs/proposals/` are historical design proposals
- `docs/ai-memory-bank/` — implementation lessons; read `lessons.md` before non-trivial changes and append dated entries when you learn something durable (see its `README.md` for format)

## Conventions

### Fluent queries (breaking-change pattern)

List/count/trash-list store methods accept query interfaces directly — never options structs:

```go
entities, err := store.EntityList(ctx, entitystore.EntityQuery().
    WithEntityType("person").
    WithLimit(10))
```

Every query interface exposes:

- `With<Field>(value)` — fluent setter returning the query for chaining
- `Get<Field>()` — accessor (empty value when unset)
- `Has<Field>() bool` — presence flag; distinguishes "unset" from "explicitly set to an invalid value"
- `Validate() error` — called by store methods, which also reject nil queries

Rules:

- An empty query is valid; explicitly set invalid values are not (`WithID("")`, `WithIDs([]string{})`, unknown `SortOrder` — only `asc`/`desc`, empty time bounds).
- All queries support inclusive UTC time bounds in `"YYYY-MM-DD HH:MM:SS"` format: `WithCreatedAtGte/Lte`, `WithUpdatedAtGte/Lte`.
- Filters must be applied uniformly to List/Count/TrashList via the shared `apply<Domain>Filters` helper in the store file — do not duplicate filter logic per method.
- Common helpers in `query_filters.go`: `applyTimeRange` (inclusive `>=`/`<=`), `escapeLike`/`applyLike` (literal LIKE matching with explicit `ESCAPE '\'` for SQLite/MySQL/PostgreSQL parity).
- `AttributeQuery` also has raw-pattern `WithAttributeKeyLike` (wildcards active) vs literal `WithAttributeKeyStartsWith/EndsWith/Contains` (wildcards escaped).

### Options structs

`*Options` structs are only for creation (`RelationshipOptions`, `TaxonomyOptions`, `TaxonomyTermOptions`, `NewStoreOptions`) — not for queries.

### Getters/setters

Domain objects use `Get*`/`Set*` (`GetID`, `GetType`, `GetEntityID`, `GetKey`, `GetValue`, `GetSlug`, `GetName`, `GetParentID`, ...). `ID()` alone is available via the embedded `dataobject.DataObjectInterface`. Entities do NOT have `GetString`/`SetString` attribute helpers — persisted attributes are read/written via `AttributeGetString`/`AttributeSetString`/etc.

### Trash

Trash tables mirror main tables with `deleted_at`/`deleted_by`. Entity/attribute trash methods live on the store implementation (`EntityTrash`, `EntityRestore`, `AttributeTrash`, `AttributeRestore`); relationship/taxonomy trash lists accept the same fluent queries as their live counterparts.

### Testing

- Query tests mirror `attribute_query_test.go`: fluent getter/presence test, `Validate()` table test, and a real-store integration test using `InitDB`.
- Enable optional features via `NewStoreOptions` (`RelationshipsEnabled`, `TaxonomiesEnabled`, `AutomigrateEnabled`).

### Documentation

- `docs/api-reference.md` is the canonical API surface — keep signatures in sync with `interfaces.go`.
- `docs/proposals/*` are dated, often-declined historical documents; mark superseded APIs as historical rather than rewriting them.
- Keep code examples compilable against the current API (fluent queries, `ctx`-first signatures, `Get*`/`Set*` names).
