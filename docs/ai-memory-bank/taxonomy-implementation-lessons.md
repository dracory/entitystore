# Taxonomy Implementation Lessons Learned

**Date:** 2026-03-28
**Status:** Completed and Documented

## Implementation Summary

Successfully implemented entity taxonomy support for entitystore following the dataobject pattern. All tests passing.

## Critical Bugs Fixed During Code Review

### 1. EntityTypes Array Handling (Data Loss Prevention)

**Problem:** Using comma-separated strings to store array data caused data loss when arrays contained empty strings.

**Solution:** Changed to JSON serialization in both `taxonomy_implementation.go` and `taxonomy_trash_implementation.go`:

```go
// Before (data loss issue)
func (o *taxonomyImplementation) EntityTypes() []string {
    typesStr := o.Get(COLUMN_ENTITY_TYPES)
    if typesStr == "" {
        return []string{}
    }
    return strings.Split(typesStr, ",")
}

// After (JSON serialization)
func (o *taxonomyImplementation) EntityTypes() []string {
    typesStr := o.Get(COLUMN_ENTITY_TYPES)
    if typesStr == "" {
        return []string{}
    }
    var types []string
    if err := json.Unmarshal([]byte(typesStr), &types); err != nil {
        return []string{}
    }
    return types
}
```

**Lesson:** Always use JSON serialization for array/slice data in dataobject pattern, never comma-separated strings.

### 2. Referential Integrity Validation

**Problem:** `EntityTaxonomyAssign` didn't validate that referenced entities, taxonomies, and terms actually exist.

**Solution:** Added comprehensive validation:

```go
// Validate entity exists
entity, err := st.EntityFindByID(ctx, entityID)
if err != nil {
    return err
}
if entity == nil {
    return errors.New("entity not found")
}

// Validate taxonomy exists
taxonomy, err := st.TaxonomyFind(ctx, taxonomyID)
if err != nil {
    return err
}
if taxonomy == nil {
    return errors.New("taxonomy not found")
}

// Validate term exists and belongs to the taxonomy
term, err := st.TaxonomyTermFind(ctx, termID)
if err != nil {
    return err
}
if term == nil {
    return errors.New("taxonomy term not found")
}
if term.TaxonomyID() != taxonomyID {
    return errors.New("taxonomy term does not belong to the specified taxonomy")
}
```

**Lesson:** Always validate referential integrity before creating relationships/assignments. Don't rely solely on database constraints.

### 3. Cascade Delete Prevention

**Problem:** Deleting taxonomies or terms with dependencies would leave orphaned records.

**Solution:** Added dependency checks before deletion:

```go
// Check for dependent taxonomy terms
termsCount, err := st.TaxonomyTermCount(ctx, TaxonomyTermQueryOptions{
    TaxonomyID: taxonomyID,
})
if err != nil {
    return false, err
}
if termsCount > 0 {
    return false, errors.New("cannot delete taxonomy: it has associated terms")
}

// Check for entity assignments
assignmentsCount, err := st.EntityTaxonomyCount(ctx, EntityTaxonomyQueryOptions{
    TaxonomyID: taxonomyID,
})
if err != nil {
    return false, err
}
if assignmentsCount > 0 {
    return false, errors.New("cannot delete taxonomy: it has associated entity assignments")
}
```

**Lesson:** Implement cascade delete prevention by checking for dependencies before deletion. Provide clear error messages about what's blocking the deletion.

### 4. Slug Conflict Validation in Updates

**Problem:** `TaxonomyUpdate` and `TaxonomyTermUpdate` didn't check for slug conflicts, leading to database constraint violations.

**Solution:** Added duplicate slug checks:

```go
// Check for slug conflicts with other taxonomies
if taxonomy.Slug() != "" {
    existing, err := st.TaxonomyFindBySlug(ctx, taxonomy.Slug())
    if err != nil {
        return err
    }
    if existing != nil && existing.ID() != taxonomy.ID() {
        return errors.New("taxonomy with this slug already exists")
    }
}
```

**Lesson:** Always validate unique constraints in application code before database operations for better error messages.

### 5. Query Filter Consistency

**Problem:** Count methods didn't apply the same filters as List methods, causing inconsistent behavior.

**Solution:** Added missing filters to all Count methods to match List methods:

```go
// TaxonomyCount now matches TaxonomyList filters
if len(options.IDs) > 0 {
    q = q.Where(goqu.C(COLUMN_ID).In(options.IDs))
}

if options.ID != "" {
    q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
}
```

**Lesson:** Count and List methods must support the same filters for consistent API behavior.

## Best Practices Established

### 1. Optional Features Pattern

Taxonomies are disabled by default and only enabled via flag:

```go
type NewStoreOptions struct {
    TaxonomiesEnabled bool
    TaxonomyTableName string
    // ... other options
}
```

**Benefit:** Backward compatibility - existing code continues to work unchanged.

### 2. Dataobject Pattern Consistency

All taxonomy types follow the same pattern as entities and attributes:
- `*_implementation.go` - Type with dataobject.DataObject
- `*_table_create_sql.go` - SQL schema using sb builder
- `store_*.go` - CRUD operations
- `*_trash_implementation.go` - Trash type for soft delete

**Benefit:** Consistent codebase, easier to maintain and extend.

### 3. Comprehensive Validation

Every operation validates:
- Required fields are present
- Referenced entities exist
- No conflicts with existing data
- Dependencies before deletion

**Benefit:** Data integrity and clear error messages.

## Files Created (20 files)

1. `taxonomy_implementation.go`
2. `taxonomy_term_implementation.go`
3. `entity_taxonomy_implementation.go`
4. `taxonomy_trash_implementation.go`
5. `taxonomy_term_trash_implementation.go`
6. `taxonomy_table_create_sql.go`
7. `taxonomy_term_table_create_sql.go`
8. `entity_taxonomy_table_create_sql.go`
9. `taxonomy_trash_table_create_sql.go`
10. `taxonomy_term_trash_table_create_sql.go`
11. `taxonomy_query.go`
12. `store_taxonomies.go`
13. `store_taxonomy_terms.go`
14. `store_entity_taxonomies.go`
15. `store_taxonomies_trash.go`
16. `store_taxonomy_terms_trash.go`

## Files Modified (4 files)

1. `consts.go` - Added taxonomy column constants
2. `interfaces.go` - Added taxonomy interfaces
3. `store_implementation.go` - Added taxonomy tables and flags
4. `new.go` - Added taxonomy options

## Testing Results

✅ All tests passing (`task test` - exit code 0)

## Known Limitations (Not Fixed)

1. **Transaction Support** - Trash/restore operations are not atomic (requires transaction support in store interface)
2. **Slug Collision on Restore** - If slug is taken while item is in trash, restore will fail
3. **Race Condition** - Small window between duplicate check and insert (mitigated by unique index)

## Future Considerations

1. Add transaction support for atomic trash/restore operations
2. Implement slug conflict resolution on restore (e.g., append suffix)
3. Add bulk operations for assigning multiple entities to taxonomies
4. Consider caching for hierarchical tree queries

## Documentation Updated

1. ✅ `docs/proposals/2026-03-28-entity-taxonomy.md` - Updated status to IMPLEMENTED with full implementation details
2. ✅ `README.md` - Added taxonomy feature, setup instructions, usage examples, and method references
3. ✅ Created this lessons learned document

## Key Takeaways

1. **JSON > CSV** - Always use JSON for array serialization in dataobject pattern
2. **Validate Early** - Check referential integrity before creating relationships
3. **Prevent Orphans** - Check dependencies before deletion
4. **Consistent APIs** - Count and List methods must support same filters
5. **User-Friendly Errors** - Validate constraints in code for better error messages
6. **Optional Features** - Use feature flags for backward compatibility
7. **Follow Patterns** - Consistency across codebase makes it maintainable
