# Entity Taxonomy Example

This example demonstrates EntityStore's taxonomy system for categorizing entities with hierarchical terms.

## What This Example Shows

### 1. Enabling Taxonomies
- Creating a store with `TaxonomiesEnabled: true`
- This creates `entities_taxonomies`, `entities_taxonomy_terms`, and `entities_entity_taxonomies` tables

### 2. Taxonomy Structure
- **Taxonomy**: A classification system (e.g., "Product Categories")
- **Terms**: Hierarchical categories within a taxonomy (e.g., "Electronics" → "Computers" → "Laptops")
- **Assignments**: Links between entities and terms

### 3. Creating Taxonomies
- Taxonomy with name, slug, and description
- Restricting taxonomies to specific entity types
- Slug-based lookups for SEO-friendly URLs

### 4. Managing Terms
- Creating hierarchical terms (parent-child relationships)
- Sort order for term ordering
- Finding terms by slug within a taxonomy

### 5. Entity Categorization
- Assigning entities to taxonomy terms
- Querying entities by category
- Counting assignments
- Removing entity from categories

### 6. Soft Deletes
- Trashing taxonomies and terms
- Restore functionality

## Running the Example

```bash
go run examples/taxonomy/main.go
```

## Running Tests

```bash
go test ./examples/taxonomy/... -v
```

## Key Concepts

**Taxonomy Tables**:
- `entities_taxonomies` - Taxonomy definitions (name, slug, entity_types as JSON)
- `entities_taxonomy_terms` - Hierarchical terms within taxonomies
- `entities_entity_taxonomies` - Entity-to-term assignments (many-to-many)

**Entity Type Restrictions**: Taxonomies can be restricted to specific entity types via the `EntityTypes` array.

**Hierarchical Terms**: Terms can have parents, creating nested category structures (e.g., Electronics > Computers > Laptops).

## Code Highlights

```go
// Enable taxonomies in store
store, _ := entitystore.NewStore(entitystore.NewStoreOptions{
    TaxonomiesEnabled: true,
    // ... other options
})

// Create taxonomy
categories, _ := store.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
    Name:        "Product Categories",
    Slug:        "product_categories",
    EntityTypes: []string{"product"},
})

// Create hierarchical term
computers, _ := store.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
    TaxonomyID: categories.ID(),
    Name:       "Computers",
    Slug:       "computers",
    ParentID:   electronics.ID(), // Parent term
})

// Assign entity to term
store.EntityTaxonomyAssign(ctx, product.ID(), categories.ID(), laptops.ID())

// Find all products in a category
assignments, _ := store.EntityTaxonomyList(ctx, entitystore.EntityTaxonomyQueryOptions{
    TaxonomyID: categories.ID(),
    TermID:     laptops.ID(),
})
```

## Use Cases

- **E-commerce**: Product categories, brands, tags
- **Content Management**: Article categories, tags, topics
- **Project Management**: Task priorities, statuses, labels
