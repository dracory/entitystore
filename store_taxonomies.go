package entitystore

import (
	"context"
	"errors"
	"log"

	"github.com/doug-martin/goqu/v9"
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
)

// ==========================================
// Taxonomy CRUD
// ==========================================

// TaxonomyCreate persists a new taxonomy record
func (st *storeImplementation) TaxonomyCreate(ctx context.Context, taxonomy TaxonomyInterface) error {
	if taxonomy == nil {
		return errors.New("taxonomy cannot be nil")
	}

	if !st.taxonomiesEnabled {
		return errors.New("taxonomies are not enabled")
	}

	// Validate required fields
	if taxonomy.Name() == "" {
		return errors.New("taxonomy name is required")
	}
	if taxonomy.Slug() == "" {
		return errors.New("taxonomy slug is required")
	}

	if taxonomy.ID() == "" {
		taxonomy.SetID(GenerateShortID())
	}

	if taxonomy.CreatedAt() == "" {
		taxonomy.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}
	if taxonomy.UpdatedAt() == "" {
		taxonomy.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	record := goqu.Record{}
	for k, v := range taxonomy.Data() {
		record[k] = v
	}

	q := goqu.Dialect(st.dbDriverName).Insert(st.taxonomyTableName).Rows(record)

	sqlStr, _, errSql := q.ToSQL()
	if errSql != nil {
		return errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := st.database.Exec(ctx, sqlStr)
	return err
}

// TaxonomyCreateByOptions creates a taxonomy using the provided options
func (st *storeImplementation) TaxonomyCreateByOptions(ctx context.Context, opts TaxonomyOptions) (TaxonomyInterface, error) {
	// Check for duplicate slug
	existing, err := st.TaxonomyFindBySlug(ctx, opts.Slug)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("taxonomy with this slug already exists")
	}

	tax := NewTaxonomy()
	tax.SetName(opts.Name)
	tax.SetSlug(opts.Slug)
	tax.SetDescription(opts.Description)
	tax.SetParentID(opts.ParentID)
	tax.SetEntityTypes(opts.EntityTypes)

	if err := st.TaxonomyCreate(ctx, tax); err != nil {
		return nil, err
	}

	return tax, nil
}

// TaxonomyDelete removes a taxonomy record by ID (hard delete)
func (st *storeImplementation) TaxonomyDelete(ctx context.Context, taxonomyID string) (bool, error) {
	if !st.taxonomiesEnabled {
		return false, errors.New("taxonomies are not enabled")
	}

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

	q := goqu.Dialect(st.dbDriverName).
		Delete(st.taxonomyTableName).
		Where(goqu.C(COLUMN_ID).Eq(taxonomyID))

	sqlStr, _, errSql := q.ToSQL()
	if errSql != nil {
		return false, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	result, err := st.database.Exec(ctx, sqlStr)
	if err != nil {
		return false, err
	}

	affected, _ := result.RowsAffected()
	return affected > 0, nil
}

// TaxonomyFind finds a taxonomy by its ID
func (st *storeImplementation) TaxonomyFind(ctx context.Context, taxonomyID string) (TaxonomyInterface, error) {
	if !st.taxonomiesEnabled {
		return nil, errors.New("taxonomies are not enabled")
	}

	if taxonomyID == "" {
		return nil, errors.New("taxonomy ID cannot be empty")
	}

	list, err := st.TaxonomyList(ctx, TaxonomyQueryOptions{
		ID:    taxonomyID,
		Limit: 1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// TaxonomyFindBySlug finds a taxonomy by its slug
func (st *storeImplementation) TaxonomyFindBySlug(ctx context.Context, slug string) (TaxonomyInterface, error) {
	if !st.taxonomiesEnabled {
		return nil, errors.New("taxonomies are not enabled")
	}

	if slug == "" {
		return nil, errors.New("slug cannot be empty")
	}

	list, err := st.TaxonomyList(ctx, TaxonomyQueryOptions{
		Slug:  slug,
		Limit: 1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// TaxonomyList lists taxonomies matching the given query options
func (st *storeImplementation) TaxonomyList(ctx context.Context, options TaxonomyQueryOptions) ([]TaxonomyInterface, error) {
	if !st.taxonomiesEnabled {
		return nil, errors.New("taxonomies are not enabled")
	}

	q := goqu.Dialect(st.dbDriverName).From(st.taxonomyTableName)

	if len(options.IDs) > 0 {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDs))
	}

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if options.Slug != "" {
		q = q.Where(goqu.C(COLUMN_SLUG).Eq(options.Slug))
	}

	if options.ParentID != "" {
		q = q.Where(goqu.C(COLUMN_PARENT_ID).Eq(options.ParentID))
	}

	sortByColumn := COLUMN_NAME
	sortOrder := "asc"

	if options.SortOrder != "" {
		sortOrder = options.SortOrder
	}

	if options.SortBy != "" {
		sortByColumn = options.SortBy
	}

	if sortOrder == "asc" {
		q = q.Order(goqu.I(sortByColumn).Asc())
	} else {
		q = q.Order(goqu.I(sortByColumn).Desc())
	}

	if options.Offset > 0 {
		q = q.Offset(uint(options.Offset))
	}

	if options.Limit > 0 {
		q = q.Limit(uint(options.Limit))
	}

	sqlStr, _, errSql := q.Select().ToSQL()
	if errSql != nil {
		return nil, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	taxonomyMaps, err := st.database.SelectToMapString(ctx, sqlStr)
	if err != nil {
		return nil, err
	}

	var list []TaxonomyInterface
	for _, m := range taxonomyMaps {
		list = append(list, NewTaxonomyFromExistingData(m))
	}

	return list, nil
}

// TaxonomyCount counts taxonomies matching the given options
func (st *storeImplementation) TaxonomyCount(ctx context.Context, options TaxonomyQueryOptions) (int64, error) {
	if !st.taxonomiesEnabled {
		return 0, errors.New("taxonomies are not enabled")
	}

	q := goqu.Dialect(st.dbDriverName).From(st.taxonomyTableName)

	if len(options.IDs) > 0 {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDs))
	}

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if options.Slug != "" {
		q = q.Where(goqu.C(COLUMN_SLUG).Eq(options.Slug))
	}

	if options.ParentID != "" {
		q = q.Where(goqu.C(COLUMN_PARENT_ID).Eq(options.ParentID))
	}

	sqlStr, _, errSql := q.Select(goqu.COUNT(goqu.Star()).As("count")).ToSQL()
	if errSql != nil {
		return 0, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	maps, err := st.database.SelectToMapString(ctx, sqlStr)
	if err != nil {
		return 0, err
	}

	if len(maps) == 0 {
		return 0, nil
	}

	count := cast.ToInt64(maps[0]["count"])
	return count, nil
}

// TaxonomyUpdate updates a taxonomy record
func (st *storeImplementation) TaxonomyUpdate(ctx context.Context, taxonomy TaxonomyInterface) error {
	if !st.taxonomiesEnabled {
		return errors.New("taxonomies are not enabled")
	}

	if taxonomy == nil {
		return errors.New("taxonomy cannot be nil")
	}

	if taxonomy.ID() == "" {
		return errors.New("taxonomy ID is required")
	}

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

	taxonomy.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	record := goqu.Record{}
	for k, v := range taxonomy.Data() {
		record[k] = v
	}

	q := goqu.Dialect(st.dbDriverName).
		Update(st.taxonomyTableName).
		Set(record).
		Where(goqu.C(COLUMN_ID).Eq(taxonomy.ID()))

	sqlStr, _, errSql := q.ToSQL()
	if errSql != nil {
		return errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := st.database.Exec(ctx, sqlStr)
	return err
}
