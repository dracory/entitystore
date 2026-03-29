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
// TaxonomyTerm CRUD
// ==========================================

// TaxonomyTermCreate persists a new taxonomy term record
func (st *storeImplementation) TaxonomyTermCreate(ctx context.Context, term TaxonomyTermInterface) error {
	if term == nil {
		return errors.New("taxonomy term cannot be nil")
	}

	if !st.taxonomiesEnabled {
		return errors.New("taxonomies are not enabled")
	}

	// Validate required fields
	if term.GetName() == "" {
		return errors.New("taxonomy term name is required")
	}
	if term.GetSlug() == "" {
		return errors.New("taxonomy term slug is required")
	}
	if term.GetTaxonomyID() == "" {
		return errors.New("taxonomy ID is required")
	}

	if term.ID() == "" {
		term.SetID(GenerateShortID())
	}

	if term.GetCreatedAt() == "" {
		term.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}
	if term.GetUpdatedAt() == "" {
		term.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	record := goqu.Record{}
	for k, v := range term.Data() {
		record[k] = v
	}

	q := goqu.Dialect(st.dbDriverName).Insert(st.taxonomyTermTableName).Rows(record)

	sqlStr, params, errSql := q.Prepared(true).ToSQL()
	if errSql != nil {
		return errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := st.database.Exec(ctx, sqlStr, params...)
	return err
}

// TaxonomyTermCreateByOptions creates a taxonomy term using the provided options
func (st *storeImplementation) TaxonomyTermCreateByOptions(ctx context.Context, opts TaxonomyTermOptions) (TaxonomyTermInterface, error) {
	// Check for duplicate slug within taxonomy
	existing, err := st.TaxonomyTermFindBySlug(ctx, opts.TaxonomyID, opts.Slug)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("taxonomy term with this slug already exists in this taxonomy")
	}

	term := NewTaxonomyTerm()
	term.SetTaxonomyID(opts.TaxonomyID)
	term.SetName(opts.Name)
	term.SetSlug(opts.Slug)
	term.SetParentID(opts.ParentID)
	term.SetSortOrder(opts.SortOrder)

	if err := st.TaxonomyTermCreate(ctx, term); err != nil {
		return nil, err
	}

	return term, nil
}

// TaxonomyTermDelete removes a taxonomy term record by ID (hard delete)
func (st *storeImplementation) TaxonomyTermDelete(ctx context.Context, termID string) (bool, error) {
	if !st.taxonomiesEnabled {
		return false, errors.New("taxonomies are not enabled")
	}

	// Check for entity assignments
	assignmentsCount, err := st.EntityTaxonomyCount(ctx, EntityTaxonomyQueryOptions{
		TermID: termID,
	})
	if err != nil {
		return false, err
	}
	if assignmentsCount > 0 {
		return false, errors.New("cannot delete taxonomy term: it has associated entity assignments")
	}

	q := goqu.Dialect(st.dbDriverName).
		Delete(st.taxonomyTermTableName).
		Where(goqu.C(COLUMN_ID).Eq(termID))

	sqlStr, params, errSql := q.Prepared(true).ToSQL()
	if errSql != nil {
		return false, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	result, err := st.database.Exec(ctx, sqlStr, params...)
	if err != nil {
		return false, err
	}

	affected, _ := result.RowsAffected()
	return affected > 0, nil
}

// TaxonomyTermFind finds a taxonomy term by its ID
func (st *storeImplementation) TaxonomyTermFind(ctx context.Context, termID string) (TaxonomyTermInterface, error) {
	if !st.taxonomiesEnabled {
		return nil, errors.New("taxonomies are not enabled")
	}

	if termID == "" {
		return nil, errors.New("taxonomy term ID cannot be empty")
	}

	list, err := st.TaxonomyTermList(ctx, TaxonomyTermQueryOptions{
		ID:    termID,
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

// TaxonomyTermFindBySlug finds a taxonomy term by its slug within a taxonomy
func (st *storeImplementation) TaxonomyTermFindBySlug(ctx context.Context, taxonomyID string, slug string) (TaxonomyTermInterface, error) {
	if !st.taxonomiesEnabled {
		return nil, errors.New("taxonomies are not enabled")
	}

	if taxonomyID == "" || slug == "" {
		return nil, errors.New("taxonomy ID and slug are required")
	}

	list, err := st.TaxonomyTermList(ctx, TaxonomyTermQueryOptions{
		TaxonomyID: taxonomyID,
		Slug:       slug,
		Limit:      1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// TaxonomyTermList lists taxonomy terms matching the given query options
func (st *storeImplementation) TaxonomyTermList(ctx context.Context, options TaxonomyTermQueryOptions) ([]TaxonomyTermInterface, error) {
	if !st.taxonomiesEnabled {
		return nil, errors.New("taxonomies are not enabled")
	}

	q := goqu.Dialect(st.dbDriverName).From(st.taxonomyTermTableName)

	if len(options.IDs) > 0 {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDs))
	}

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if options.TaxonomyID != "" {
		q = q.Where(goqu.C(COLUMN_TAXONOMY_ID).Eq(options.TaxonomyID))
	}

	if options.Slug != "" {
		q = q.Where(goqu.C(COLUMN_SLUG).Eq(options.Slug))
	}

	if options.ParentID != "" {
		q = q.Where(goqu.C(COLUMN_PARENT_ID).Eq(options.ParentID))
	}

	sortByColumn := COLUMN_SORT_ORDER
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

	sqlStr, params, errSql := q.Prepared(true).Select().ToSQL()
	if errSql != nil {
		return nil, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	termMaps, err := st.database.SelectToMapString(ctx, sqlStr, params...)
	if err != nil {
		return nil, err
	}

	var list []TaxonomyTermInterface
	for _, m := range termMaps {
		list = append(list, NewTaxonomyTermFromExistingData(m))
	}

	return list, nil
}

// TaxonomyTermCount counts taxonomy terms matching the given options
func (st *storeImplementation) TaxonomyTermCount(ctx context.Context, options TaxonomyTermQueryOptions) (int64, error) {
	if !st.taxonomiesEnabled {
		return 0, errors.New("taxonomies are not enabled")
	}

	q := goqu.Dialect(st.dbDriverName).From(st.taxonomyTermTableName)

	if len(options.IDs) > 0 {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDs))
	}

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if options.TaxonomyID != "" {
		q = q.Where(goqu.C(COLUMN_TAXONOMY_ID).Eq(options.TaxonomyID))
	}

	if options.Slug != "" {
		q = q.Where(goqu.C(COLUMN_SLUG).Eq(options.Slug))
	}

	if options.ParentID != "" {
		q = q.Where(goqu.C(COLUMN_PARENT_ID).Eq(options.ParentID))
	}

	sqlStr, params, errSql := q.Prepared(true).Select(goqu.COUNT(goqu.Star()).As("count")).ToSQL()
	if errSql != nil {
		return 0, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	maps, err := st.database.SelectToMapString(ctx, sqlStr, params...)
	if err != nil {
		return 0, err
	}

	if len(maps) == 0 {
		return 0, nil
	}

	count := cast.ToInt64(maps[0]["count"])
	return count, nil
}

// TaxonomyTermUpdate updates a taxonomy term record
func (st *storeImplementation) TaxonomyTermUpdate(ctx context.Context, term TaxonomyTermInterface) error {
	if !st.taxonomiesEnabled {
		return errors.New("taxonomies are not enabled")
	}

	if term == nil {
		return errors.New("taxonomy term cannot be nil")
	}

	if term.ID() == "" {
		return errors.New("taxonomy term ID is required")
	}

	// Check for slug conflicts within the same taxonomy
	if term.GetSlug() != "" && term.GetTaxonomyID() != "" {
		existing, err := st.TaxonomyTermFindBySlug(ctx, term.GetTaxonomyID(), term.GetSlug())
		if err != nil {
			return err
		}
		if existing != nil && existing.ID() != term.ID() {
			return errors.New("taxonomy term with this slug already exists in this taxonomy")
		}
	}

	term.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	record := goqu.Record{}
	for k, v := range term.Data() {
		record[k] = v
	}

	q := goqu.Dialect(st.dbDriverName).
		Update(st.taxonomyTermTableName).
		Set(record).
		Where(goqu.C(COLUMN_ID).Eq(term.ID()))

	sqlStr, params, errSql := q.Prepared(true).ToSQL()
	if errSql != nil {
		return errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := st.database.Exec(ctx, sqlStr, params...)
	return err
}
