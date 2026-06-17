package entitystore

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/dromara/carbon/v2"
)

// taxonomyTermRow is used for scanning taxonomy term query results
type taxonomyTermRow struct {
	ID         string `db:"id"`
	TaxonomyID string `db:"taxonomy_id"`
	Name       string `db:"name"`
	Slug       string `db:"slug"`
	ParentID   string `db:"parent_id"`
	SortOrder  int    `db:"sort_order"`
	CreatedAt  string `db:"created_at"`
	UpdatedAt  string `db:"updated_at"`
}

// TaxonomyTermCreate persists a new taxonomy term record
func (st *storeImplementation) TaxonomyTermCreate(ctx context.Context, term TaxonomyTermInterface) error {
	if term == nil {
		return errors.New("taxonomy term cannot be nil")
	}

	if !st.taxonomiesEnabled {
		return errors.New("taxonomies are not enabled")
	}

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

	row := map[string]any{}
	for k, v := range term.Data() {
		row[k] = v
	}

	if st.GetDebug() {
		log.Println("TaxonomyTermCreate:", row)
	}

	return st.db.Query().Table(st.taxonomyTermTableName).Create(row)
}

// TaxonomyTermCreateByOptions creates a taxonomy term using the provided options
func (st *storeImplementation) TaxonomyTermCreateByOptions(ctx context.Context, opts TaxonomyTermOptions) (TaxonomyTermInterface, error) {
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

	assignmentsCount, err := st.EntityTaxonomyCount(ctx, EntityTaxonomyQueryOptions{
		TermID: termID,
	})
	if err != nil {
		return false, err
	}
	if assignmentsCount > 0 {
		return false, errors.New("cannot delete taxonomy term: it has associated entity assignments")
	}

	result, err := st.db.Query().Table(st.taxonomyTermTableName).Where(COLUMN_ID+" = ?", termID).Delete()
	if err != nil {
		return false, err
	}

	return result.RowsAffected > 0, nil
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

	q := st.db.Query().Table(st.taxonomyTermTableName)

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

	if options.TaxonomyID != "" {
		q = q.Where(COLUMN_TAXONOMY_ID+" = ?", options.TaxonomyID)
	}

	if options.Slug != "" {
		q = q.Where(COLUMN_SLUG+" = ?", options.Slug)
	}

	if options.ParentID != "" {
		q = q.Where(COLUMN_PARENT_ID+" = ?", options.ParentID)
	}

	sortByColumn := COLUMN_SORT_ORDER
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

	var rows []taxonomyTermRow
	if err := q.Get(&rows); err != nil {
		return nil, err
	}

	var list []TaxonomyTermInterface
	for _, r := range rows {
		list = append(list, NewTaxonomyTermFromExistingData(map[string]string{
			COLUMN_ID:          r.ID,
			COLUMN_TAXONOMY_ID: r.TaxonomyID,
			COLUMN_NAME:        r.Name,
			COLUMN_SLUG:        r.Slug,
			COLUMN_PARENT_ID:   r.ParentID,
			COLUMN_SORT_ORDER:  strconv.Itoa(r.SortOrder),
			COLUMN_CREATED_AT:  r.CreatedAt,
			COLUMN_UPDATED_AT:  r.UpdatedAt,
		}))
	}

	return list, nil
}

// TaxonomyTermCount counts taxonomy terms matching the given options
func (st *storeImplementation) TaxonomyTermCount(ctx context.Context, options TaxonomyTermQueryOptions) (int64, error) {
	if !st.taxonomiesEnabled {
		return 0, errors.New("taxonomies are not enabled")
	}

	q := st.db.Query().Table(st.taxonomyTermTableName)

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

	if options.TaxonomyID != "" {
		q = q.Where(COLUMN_TAXONOMY_ID+" = ?", options.TaxonomyID)
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

	row := map[string]any{}
	for k, v := range term.Data() {
		row[k] = v
	}

	if st.GetDebug() {
		log.Println("TaxonomyTermUpdate:", row)
	}

	_, err := st.db.Query().Table(st.taxonomyTermTableName).Where(COLUMN_ID+" = ?", term.ID()).Update(row)
	return err
}
