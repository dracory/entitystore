package entitystore

import (
	"context"
	"errors"
	"log"

	"github.com/dromara/carbon/v2"
)

// taxonomyRow is used for scanning taxonomy query results
type taxonomyRow struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	Slug        string `db:"slug"`
	Description string `db:"description"`
	ParentID    string `db:"parent_id"`
	EntityTypes string `db:"entity_types"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
}

// TaxonomyCreate persists a new taxonomy record
func (st *storeImplementation) TaxonomyCreate(ctx context.Context, taxonomy TaxonomyInterface) error {
	if taxonomy == nil {
		return errors.New("taxonomy cannot be nil")
	}

	if !st.taxonomiesEnabled {
		return errors.New("taxonomies are not enabled")
	}

	if taxonomy.GetName() == "" {
		return errors.New("taxonomy name is required")
	}
	if taxonomy.GetSlug() == "" {
		return errors.New("taxonomy slug is required")
	}

	if taxonomy.ID() == "" {
		taxonomy.SetID(GenerateShortID())
	}

	if taxonomy.GetCreatedAt() == "" {
		taxonomy.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}
	if taxonomy.GetUpdatedAt() == "" {
		taxonomy.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	row := map[string]any{}
	for k, v := range taxonomy.Data() {
		row[k] = v
	}

	if st.GetDebug() {
		log.Println("TaxonomyCreate:", row)
	}

	return st.db.Query().Table(st.taxonomyTableName).Create(row)
}

// TaxonomyCreateByOptions creates a taxonomy using the provided options
func (st *storeImplementation) TaxonomyCreateByOptions(ctx context.Context, opts TaxonomyOptions) (TaxonomyInterface, error) {
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

	termsCount, err := st.TaxonomyTermCount(ctx, TaxonomyTermQueryOptions{
		TaxonomyID: taxonomyID,
	})
	if err != nil {
		return false, err
	}
	if termsCount > 0 {
		return false, errors.New("cannot delete taxonomy: it has associated terms")
	}

	assignmentsCount, err := st.EntityTaxonomyCount(ctx, EntityTaxonomyQueryOptions{
		TaxonomyID: taxonomyID,
	})
	if err != nil {
		return false, err
	}
	if assignmentsCount > 0 {
		return false, errors.New("cannot delete taxonomy: it has associated entity assignments")
	}

	result, err := st.db.Query().Table(st.taxonomyTableName).Where(COLUMN_ID+" = ?", taxonomyID).Delete()
	if err != nil {
		return false, err
	}

	return result.RowsAffected > 0, nil
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

	q := st.db.Query().Table(st.taxonomyTableName)

	if options.ID != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID)
	}

	if len(options.IDs) > 0 {
		ids := make([]any, len(options.IDs))
		for i, id := range options.IDs {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_ID, ids)
	}

	if options.Slug != "" {
		q = q.Where(COLUMN_SLUG+" = ?", options.Slug)
	}

	if options.ParentID != "" {
		q = q.Where(COLUMN_PARENT_ID+" = ?", options.ParentID)
	}

	sortByColumn := COLUMN_NAME
	sortOrder := "asc"

	if options.SortOrder != "" {
		sortOrder = options.SortOrder
	}

	if options.SortBy != "" {
		sortByColumn = options.SortBy
	}

	q = q.OrderBy(sortByColumn, sortOrder)

	if options.Offset > 0 {
		q = q.Offset(int(options.Offset))
	}

	if options.Limit > 0 {
		q = q.Limit(int(options.Limit))
	}

	var rows []taxonomyRow
	if err := q.Get(&rows); err != nil {
		return nil, err
	}

	var list []TaxonomyInterface
	for _, r := range rows {
		list = append(list, NewTaxonomyFromExistingData(map[string]string{
			COLUMN_ID:           r.ID,
			COLUMN_NAME:         r.Name,
			COLUMN_SLUG:         r.Slug,
			COLUMN_DESCRIPTION:  r.Description,
			COLUMN_PARENT_ID:    r.ParentID,
			COLUMN_ENTITY_TYPES: r.EntityTypes,
			COLUMN_CREATED_AT:   r.CreatedAt,
			COLUMN_UPDATED_AT:   r.UpdatedAt,
		}))
	}

	return list, nil
}

// TaxonomyCount counts taxonomies matching the given options
func (st *storeImplementation) TaxonomyCount(ctx context.Context, options TaxonomyQueryOptions) (int64, error) {
	if !st.taxonomiesEnabled {
		return 0, errors.New("taxonomies are not enabled")
	}

	q := st.db.Query().Table(st.taxonomyTableName)

	if options.ID != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID)
	}

	if len(options.IDs) > 0 {
		ids := make([]any, len(options.IDs))
		for i, id := range options.IDs {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_ID, ids)
	}

	if options.Slug != "" {
		q = q.Where(COLUMN_SLUG+" = ?", options.Slug)
	}

	if options.ParentID != "" {
		q = q.Where(COLUMN_PARENT_ID+" = ?", options.ParentID)
	}

	var count int64
	if err := q.Count(&count); err != nil {
		return 0, err
	}

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

	if taxonomy.GetSlug() != "" {
		existing, err := st.TaxonomyFindBySlug(ctx, taxonomy.GetSlug())
		if err != nil {
			return err
		}
		if existing != nil && existing.ID() != taxonomy.ID() {
			return errors.New("taxonomy with this slug already exists")
		}
	}

	taxonomy.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	row := map[string]any{}
	for k, v := range taxonomy.Data() {
		row[k] = v
	}

	if st.GetDebug() {
		log.Println("TaxonomyUpdate:", row)
	}

	_, err := st.db.Query().Table(st.taxonomyTableName).Where(COLUMN_ID+" = ?", taxonomy.ID()).Update(row)
	return err
}
