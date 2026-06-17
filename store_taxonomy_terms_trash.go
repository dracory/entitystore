package entitystore

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/dromara/carbon/v2"
)

// taxonomyTermTrashRow is used for scanning taxonomy term trash query results
type taxonomyTermTrashRow struct {
	ID         string `db:"id"`
	TaxonomyID string `db:"taxonomy_id"`
	Name       string `db:"name"`
	Slug       string `db:"slug"`
	ParentID   string `db:"parent_id"`
	SortOrder  int    `db:"sort_order"`
	CreatedAt  string `db:"created_at"`
	UpdatedAt  string `db:"updated_at"`
	DeletedAt  string `db:"deleted_at"`
	DeletedBy  string `db:"deleted_by"`
}

// TaxonomyTermTrash soft-deletes a taxonomy term by moving it to trash
func (st *storeImplementation) TaxonomyTermTrash(ctx context.Context, termID string, deletedBy string) (bool, error) {
	if !st.taxonomiesEnabled {
		return false, errors.New("taxonomies are not enabled")
	}

	term, err := st.TaxonomyTermFind(ctx, termID)
	if err != nil {
		return false, err
	}
	if term == nil {
		return false, errors.New("taxonomy term not found")
	}

	assignmentsCount, err := st.EntityTaxonomyCount(ctx, EntityTaxonomyQueryOptions{
		TermID: termID,
	})
	if err != nil {
		return false, err
	}
	if assignmentsCount > 0 {
		return false, errors.New("cannot trash taxonomy term: it has associated entity assignments")
	}

	trash := NewTaxonomyTermTrash()
	trash.SetID(term.ID())
	trash.SetTaxonomyID(term.GetTaxonomyID())
	trash.SetName(term.GetName())
	trash.SetSlug(term.GetSlug())
	trash.SetParentID(term.GetParentID())
	trash.SetSortOrder(term.GetSortOrder())
	trash.SetCreatedAt(term.GetCreatedAt())
	trash.SetUpdatedAt(term.GetUpdatedAt())
	trash.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	trash.SetDeletedBy(deletedBy)

	row := map[string]any{}
	for k, v := range trash.Data() {
		row[k] = v
	}

	if st.GetDebug() {
		log.Println("TaxonomyTermTrash insert:", row)
	}

	if err := st.db.Query().Table(st.taxonomyTermTrashTableName).Create(row); err != nil {
		return false, err
	}

	return st.TaxonomyTermDelete(ctx, termID)
}

// TaxonomyTermRestore restores a taxonomy term from trash
func (st *storeImplementation) TaxonomyTermRestore(ctx context.Context, termID string) (bool, error) {
	if !st.taxonomiesEnabled {
		return false, errors.New("taxonomies are not enabled")
	}

	list, err := st.TaxonomyTermTrashList(ctx, TaxonomyTermQueryOptions{
		ID:    termID,
		Limit: 1,
	})
	if err != nil {
		return false, err
	}
	if len(list) == 0 {
		return false, errors.New("taxonomy term not found in trash")
	}

	trash := list[0]

	term := NewTaxonomyTerm()
	term.SetID(trash.ID())
	term.SetTaxonomyID(trash.GetTaxonomyID())
	term.SetName(trash.GetName())
	term.SetSlug(trash.GetSlug())
	term.SetParentID(trash.GetParentID())
	term.SetSortOrder(trash.GetSortOrder())
	term.SetCreatedAt(trash.GetCreatedAt())
	term.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	if err := st.TaxonomyTermCreate(ctx, term); err != nil {
		return false, err
	}

	result, err := st.db.Query().Table(st.taxonomyTermTrashTableName).Where(COLUMN_ID+" = ?", termID).Delete()
	if err != nil {
		return false, err
	}

	return result.RowsAffected > 0, nil
}

// TaxonomyTermTrashList lists trashed taxonomy terms
func (st *storeImplementation) TaxonomyTermTrashList(ctx context.Context, options TaxonomyTermQueryOptions) ([]TaxonomyTermTrashInterface, error) {
	if !st.taxonomiesEnabled {
		return nil, errors.New("taxonomies are not enabled")
	}

	q := st.db.Query().Table(st.taxonomyTermTrashTableName)

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

	sortByColumn := COLUMN_DELETED_AT
	sortOrder := "desc"

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

	var rows []taxonomyTermTrashRow
	if err := q.Get(&rows); err != nil {
		return nil, err
	}

	var list []TaxonomyTermTrashInterface
	for _, r := range rows {
		list = append(list, NewTaxonomyTermTrashFromExistingData(map[string]string{
			COLUMN_ID:          r.ID,
			COLUMN_TAXONOMY_ID: r.TaxonomyID,
			COLUMN_NAME:        r.Name,
			COLUMN_SLUG:        r.Slug,
			COLUMN_PARENT_ID:   r.ParentID,
			COLUMN_SORT_ORDER:  strconv.Itoa(r.SortOrder),
			COLUMN_CREATED_AT:  r.CreatedAt,
			COLUMN_UPDATED_AT:  r.UpdatedAt,
			COLUMN_DELETED_AT:  r.DeletedAt,
			COLUMN_DELETED_BY:  r.DeletedBy,
		}))
	}

	return list, nil
}
