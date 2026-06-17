package entitystore

import (
	"context"
	"errors"
	"log"

	"github.com/dromara/carbon/v2"
)

// taxonomyTrashRow is used for scanning taxonomy trash query results
type taxonomyTrashRow struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	Slug        string `db:"slug"`
	Description string `db:"description"`
	ParentID    string `db:"parent_id"`
	EntityTypes string `db:"entity_types"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
	DeletedAt   string `db:"deleted_at"`
	DeletedBy   string `db:"deleted_by"`
}

// TaxonomyTrash soft-deletes a taxonomy by moving it to trash
func (st *storeImplementation) TaxonomyTrash(ctx context.Context, taxonomyID string, deletedBy string) (bool, error) {
	if !st.taxonomiesEnabled {
		return false, errors.New("taxonomies are not enabled")
	}

	taxonomy, err := st.TaxonomyFind(ctx, taxonomyID)
	if err != nil {
		return false, err
	}
	if taxonomy == nil {
		return false, errors.New("taxonomy not found")
	}

	termsCount, err := st.TaxonomyTermCount(ctx, TaxonomyTermQueryOptions{
		TaxonomyID: taxonomyID,
	})
	if err != nil {
		return false, err
	}
	if termsCount > 0 {
		return false, errors.New("cannot trash taxonomy: it has associated terms")
	}

	assignmentsCount, err := st.EntityTaxonomyCount(ctx, EntityTaxonomyQueryOptions{
		TaxonomyID: taxonomyID,
	})
	if err != nil {
		return false, err
	}
	if assignmentsCount > 0 {
		return false, errors.New("cannot trash taxonomy: it has associated entity assignments")
	}

	trash := NewTaxonomyTrash()
	trash.SetID(taxonomy.ID())
	trash.SetName(taxonomy.GetName())
	trash.SetSlug(taxonomy.GetSlug())
	trash.SetDescription(taxonomy.GetDescription())
	trash.SetParentID(taxonomy.GetParentID())
	trash.SetEntityTypes(taxonomy.GetEntityTypes())
	trash.SetCreatedAt(taxonomy.GetCreatedAt())
	trash.SetUpdatedAt(taxonomy.GetUpdatedAt())
	trash.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	trash.SetDeletedBy(deletedBy)

	row := map[string]any{}
	for k, v := range trash.Data() {
		row[k] = v
	}

	if st.GetDebug() {
		log.Println("TaxonomyTrash insert:", row)
	}

	if err := st.db.Query().Table(st.taxonomyTrashTableName).Create(row); err != nil {
		return false, err
	}

	return st.TaxonomyDelete(ctx, taxonomyID)
}

// TaxonomyRestore restores a taxonomy from trash
func (st *storeImplementation) TaxonomyRestore(ctx context.Context, taxonomyID string) (bool, error) {
	if !st.taxonomiesEnabled {
		return false, errors.New("taxonomies are not enabled")
	}

	list, err := st.TaxonomyTrashList(ctx, TaxonomyQueryOptions{
		ID:    taxonomyID,
		Limit: 1,
	})
	if err != nil {
		return false, err
	}
	if len(list) == 0 {
		return false, errors.New("taxonomy not found in trash")
	}

	trash := list[0]

	taxonomy := NewTaxonomy()
	taxonomy.SetID(trash.ID())
	taxonomy.SetName(trash.GetName())
	taxonomy.SetSlug(trash.GetSlug())
	taxonomy.SetDescription(trash.GetDescription())
	taxonomy.SetParentID(trash.GetParentID())
	taxonomy.SetEntityTypes(trash.GetEntityTypes())
	taxonomy.SetCreatedAt(trash.GetCreatedAt())
	taxonomy.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	if err := st.TaxonomyCreate(ctx, taxonomy); err != nil {
		return false, err
	}

	result, err := st.db.Query().Table(st.taxonomyTrashTableName).Where(COLUMN_ID+" = ?", taxonomyID).Delete()
	if err != nil {
		return false, err
	}

	return result.RowsAffected > 0, nil
}

// TaxonomyTrashList lists trashed taxonomies
func (st *storeImplementation) TaxonomyTrashList(ctx context.Context, options TaxonomyQueryOptions) ([]TaxonomyTrashInterface, error) {
	if !st.taxonomiesEnabled {
		return nil, errors.New("taxonomies are not enabled")
	}

	q := st.db.Query().Table(st.taxonomyTrashTableName)

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

	var rows []taxonomyTrashRow
	if err := q.Get(&rows); err != nil {
		return nil, err
	}

	var list []TaxonomyTrashInterface
	for _, r := range rows {
		list = append(list, NewTaxonomyTrashFromExistingData(map[string]string{
			COLUMN_ID:           r.ID,
			COLUMN_NAME:         r.Name,
			COLUMN_SLUG:         r.Slug,
			COLUMN_DESCRIPTION:  r.Description,
			COLUMN_PARENT_ID:    r.ParentID,
			COLUMN_ENTITY_TYPES: r.EntityTypes,
			COLUMN_CREATED_AT:   r.CreatedAt,
			COLUMN_UPDATED_AT:   r.UpdatedAt,
			COLUMN_DELETED_AT:   r.DeletedAt,
			COLUMN_DELETED_BY:   r.DeletedBy,
		}))
	}

	return list, nil
}
