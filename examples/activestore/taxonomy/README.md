# ActiveStore Taxonomy Example

Demonstrates entity-centric taxonomy sugar: `TaxonomyCreate`, `TermCreate`, `AssignTerm`, `Terms`, and navigation via `term.Entities()`.

Requires `TaxonomiesEnabled: true` in `NewStoreOptions`.

## Running

```bash
go run examples/activestore/taxonomy/main.go
go test ./examples/activestore/taxonomy/... -v
```

## Code Highlights

```go
cat, _ := active.TaxonomyCreate(entitystore.TaxonomyOptions{Slug: "categories"})
term, _ := active.TermCreate(entitystore.TaxonomyTermOptions{
    TaxonomyID: cat.GetTaxonomy().GetID(), Slug: "go",
})

_ = post.AssignTerm(cat.GetTaxonomy().GetID(), term.GetTerm().GetID())

// Navigate both directions without writing queries
terms, _ := post.Terms(cat.GetTaxonomy().GetID())
posts, _ := term.Entities()
```
