# ActiveStore Relationships Example

Demonstrates entity-centric relationship sugar: `RelateTo`, `Related`, `RelateToOrdered`, `Unrelate`.

Requires `RelationshipsEnabled: true` in `NewStoreOptions`.

## Running

```bash
go run examples/activestore/relationships/main.go
go test ./examples/activestore/relationships/... -v
```

## Code Highlights

```go
_ = post.RelateTo(author.GetEntity().ID(), "written_by")

// Hydrates linked entities in one call — no RelationshipList +
// EntityFindByID loop
authors, _ := post.Related("written_by")

_ = post.Unrelate(author.GetEntity().ID(), "written_by")
```
