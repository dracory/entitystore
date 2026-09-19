package entitystore

import (
	"context"
	"errors"
	"log"

	"github.com/dracory/neat/contracts/database/orm"
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

	if err := validateIDLength(taxonomy.ID(), COLUMN_ID); err != nil {
		return err
	}
	if err := validateIDLength(taxonomy.GetParentID(), COLUMN_PARENT_ID); err != nil {
		return err
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

	termsCount, err := st.TaxonomyTermCount(ctx, TaxonomyTermQuery().
		WithTaxonomyID(taxonomyID))
	if err != nil {
		return false, err
	}
	if termsCount > 0 {
		return false, errors.New("cannot delete taxonomy: it has associated terms")
	}

	assignmentsCount, err := st.EntityTaxonomyCount(ctx, EntityTaxonomyQuery().
		WithTaxonomyID(taxonomyID))
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

	list, err := st.TaxonomyList(ctx, TaxonomyQuery().
		WithID(taxonomyID).
		WithLimit(1))

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

	list, err := st.TaxonomyList(ctx, TaxonomyQuery().
		WithSlug(slug).
		WithLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// applyTaxonomyFilters applies the common filter clauses from a validated
// taxonomy query to a query. Shared by TaxonomyList, TaxonomyCount and
// TaxonomyTrashList to prevent filter drift.
func (st *storeImplementation) applyTaxonomyFilters(q orm.Query, query TaxonomyQueryInterface) orm.Query {
	if query.GetID() != "" {
		q = q.Where(COLUMN_ID+" = ?", query.GetID())
	}

	if len(query.GetIDs()) > 0 {
		ids := make([]any, len(query.GetIDs()))
		for i, id := range query.GetIDs() {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_ID, ids)
	}

	if query.GetSlug() != "" {
		q = q.Where(COLUMN_SLUG+" = ?", query.GetSlug())
	}

	if query.GetParentID() != "" {
		q = q.Where(COLUMN_PARENT_ID+" = ?", query.GetParentID())
	}

	q = applyTimeRange(q, COLUMN_CREATED_AT, query.GetCreatedAtGte(), query.GetCreatedAtLte())
	q = applyTimeRange(q, COLUMN_UPDATED_AT, query.GetUpdatedAtGte(), query.GetUpdatedAtLte())

	return q
}

// TaxonomyList lists taxonomies matching the given fluent query
func (st *storeImplementation) TaxonomyList(ctx context.Context, query TaxonomyQueryInterface) ([]TaxonomyInterface, error) {
	if !st.taxonomiesEnabled {
		return nil, errors.New("taxonomies are not enabled")
	}

	if query == nil {
		return nil, errors.New("taxonomy query cannot be nil")
	}

	if err := query.Validate(); err != nil {
		return nil, err
	}

	q := st.applyTaxonomyFilters(st.db.Query().Table(st.taxonomyTableName), query)

	sortByColumn := COLUMN_NAME
	sortOrder := "asc"

	if query.GetSortOrder() != "" {
		sortOrder = query.GetSortOrder()
	}

	if query.GetSortBy() != "" {
		sortByColumn = query.GetSortBy()
	}

	q = q.OrderBy(sortByColumn, sortOrder)

	if query.GetOffset() > 0 {
		q = q.Offset(int(query.GetOffset()))
	}

	if query.GetLimit() > 0 {
		q = q.Limit(int(query.GetLimit()))
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

// TaxonomyCount counts taxonomies matching the given fluent query
func (st *storeImplementation) TaxonomyCount(ctx context.Context, query TaxonomyQueryInterface) (int64, error) {
	if !st.taxonomiesEnabled {
		return 0, errors.New("taxonomies are not enabled")
	}

	if query == nil {
		return 0, errors.New("taxonomy query cannot be nil")
	}

	if err := query.Validate(); err != nil {
		return 0, err
	}

	q := st.applyTaxonomyFilters(st.db.Query().Table(st.taxonomyTableName), query)

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

	if err := validateIDRequired(taxonomy.ID(), COLUMN_ID); err != nil {
		return err
	}
	if err := validateIDLength(taxonomy.ID(), COLUMN_ID); err != nil {
		return err
	}
	if err := validateIDLength(taxonomy.GetParentID(), COLUMN_PARENT_ID); err != nil {
		return err
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
