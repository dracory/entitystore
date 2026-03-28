# Pros and Cons

When to use Entity Store and when to choose something else.

## When to Use Entity Store

### Good Fit

| Scenario | Why Entity Store Works |
|----------|------------------------|
| **Rapid prototyping** | No migrations needed, add fields on the fly |
| **User-defined fields** | End users can create custom attributes |
| **CMS/CMDB systems** | Flexible content types with relationships |
| **E-commerce catalogs** | Products with varying specifications |
| **Multi-tenant apps** | Different tenants need different schemas |
| **Audit trails** | Soft deletes built into all operations |
| **SQL expertise exists** | Team knows SQL, wants schemaless flexibility |

### Specific Use Cases

```go
// E-commerce product with custom attributes
product := store.EntityCreateWithType("product")
product.SetString("name", "Laptop")
product.SetFloat("price", 999.99)
product.SetInterface("specs", map[string]string{
    "cpu": "Intel i7",
    "ram": "16GB",
    "ssd": "512GB",
})

// Later add new attribute without migration
product.SetString("warranty", "2 years")
```

```go
// CMS content with taxonomy
article := store.EntityCreateWithType("article")
article.SetString("title", "Go Tips")
store.EntityTaxonomyAssign(ctx, article.ID(), blogCategories.ID(), golangTag.ID())
```

## When NOT to Use Entity Store

### Bad Fit

| Scenario | Better Alternative |
|----------|-------------------|
| **Simple CRUD apps** | Regular SQL tables + GORM/sqlx |
| **High-write analytics** | ClickHouse, TimescaleDB |
| **Full-text search** | Elasticsearch, Meilisearch |
| **Graph data** | Neo4j, Dgraph |
| **Document storage** | MongoDB (native JSON) |
| **Key-value cache** | Redis |
| **Time-series data** | InfluxDB, TimescaleDB |

### Anti-Patterns

```go
// DON'T: Use for simple user table with fixed schema
// Regular SQL is better:
// CREATE TABLE users (id INT, email VARCHAR, name VARCHAR);

// DON'T: Store large blobs/files
// Use S3/file storage instead:
// file.SetString("s3_url", "s3://bucket/file.pdf")

// DON'T: Complex reporting without indexes
// EAV requires joins; materialized views or separate tables for reporting
```

## Comparison with Alternatives

### vs Raw SQL Tables

| Aspect | Entity Store | Raw SQL |
|--------|--------------|---------|
| **Schema changes** | No migrations | Migrations required |
| **Query complexity** | Higher (EAV joins) | Lower (direct tables) |
| **Type safety** | Runtime | Compile-time with structs |
| **Performance** | Slightly slower | Faster (native columns) |
| **Flexibility** | High | Low |

**Verdict:** Use raw SQL for stable schemas, Entity Store for evolving data.

### vs MongoDB

| Aspect | Entity Store | MongoDB |
|--------|--------------|---------|
| **Transactions** | Full ACID | Limited (single doc) |
| **Query language** | SQL | Mongo query language |
| **Horizontal scale** | Hard (SQL) | Easy (sharding) |
| **Embedded documents** | Manual (JSON strings) | Native |
| **Referential integrity** | Enforced | Manual |
| **Hosting** | Any SQL provider | MongoDB Atlas or self-hosted |

**Verdict:** MongoDB for massive scale or document-centric apps; Entity Store for transactional integrity with flexibility.

### vs GORM/ORM

| Aspect | Entity Store | GORM |
|--------|--------------|------|
| **Migrations** | Not needed | Required |
| **Model definitions** | None | Struct tags |
| **Query builder** | Manual (goqu) | Rich API |
| **Hooks/callbacks** | Manual | Built-in |
| **Community** | Smaller | Large ecosystem |

**Verdict:** GORM for standard apps, Entity Store when schema varies by tenant/user.

### vs Key-Value Stores (Redis, DynamoDB)

| Aspect | Entity Store | Key-Value |
|--------|--------------|-----------|
| **Querying** | SQL queries | Key lookup only |
| **Relationships** | Supported | Manual |
| **Persistence** | Durable | Varies (Redis = volatile) |
| **Complex data** | Good | Poor (flattening needed) |

**Verdict:** Key-value for caching/session storage; Entity Store for primary data with queries.

## Performance Considerations

### Where It's Fast

- **Entity lookups by ID** - Indexed primary key
- **Attribute writes** - Simple INSERT/UPDATE
- **Small datasets (< 100K entities)** - Negligible overhead

### Where It's Slow

- **Complex filters** - Multiple attribute joins
- **Full-text search** - No built-in text indexing
- **Aggregation** - GROUP BY requires creative SQL

### Optimization Tips

```go
// GOOD: Filter by entity_type first
store.EntityList(ctx, EntityQueryOptions{
    EntityType: "product", // Uses index
    Limit:      20,
})

// BAD: Load all then filter in code
all, _ := store.EntityList(ctx, EntityQueryOptions{})
products := filterByType(all, "product") // Slow!

// GOOD: Use IDs for batch operations
store.EntityList(ctx, EntityQueryOptions{
    IDs: []string{"id1", "id2", "id3"},
})
```

## Complexity Trade-offs

### Added Complexity

- **Learning EAV** - 30 minutes to understand pattern
- **Query construction** - More verbose than SQL tables
- **Type conversions** - Manual int/float/interface handling

### Removed Complexity

- **No migrations** - Zero schema management
- **No model files** - No struct definitions
- **Flexible attributes** - Add fields without deployment

## Team Fit

### Ideal Team Profile

- Comfortable with SQL
- Building flexible/CMS-like systems
- Want ACID without schema rigidity
- Go expertise exists

### Poor Team Profile

- Need strict compile-time type safety
- Heavy reporting/analytics needs
- Want MongoDB-style document embedding
- No SQL knowledge

## Summary Decision Tree

```
Need schemaless storage?
├── No → Use regular SQL tables
└── Yes → Need ACID transactions?
    ├── No → Consider MongoDB
    └── Yes → Need complex relationships?
        ├── No → Key-value store might work
        └── Yes → Entity Store is a good fit
```
