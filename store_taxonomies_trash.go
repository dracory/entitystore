package entitystore

import (
	"context"
	"errors"
	"log"

	"github.com/doug-martin/goqu/v9"
	"github.com/dromara/carbon/v2"
)

// ==========================================
// Taxonomy Trash/Restore
// ==========================================

// TaxonomyTrash soft-deletes a taxonomy by moving it to trash
func (st *storeImplementation) TaxonomyTrash(ctx context.Context, taxonomyID string, deletedBy string) (bool, error) {
	if !st.taxonomiesEnabled {
		return false, errors.New("taxonomies are not enabled")
	}

	// Find the taxonomy
	taxonomy, err := st.TaxonomyFind(ctx, taxonomyID)
	if err != nil {
		return false, err
	}
	if taxonomy == nil {
		return false, errors.New("taxonomy not found")
	}

	// Check for dependent taxonomy terms
	termsCount, err := st.TaxonomyTermCount(ctx, TaxonomyTermQueryOptions{
		TaxonomyID: taxonomyID,
	})
	if err != nil {
		return false, err
	}
	if termsCount > 0 {
		return false, errors.New("cannot trash taxonomy: it has associated terms")
	}

	// Check for entity assignments
	assignmentsCount, err := st.EntityTaxonomyCount(ctx, EntityTaxonomyQueryOptions{
		TaxonomyID: taxonomyID,
	})
	if err != nil {
		return false, err
	}
	if assignmentsCount > 0 {
		return false, errors.New("cannot trash taxonomy: it has associated entity assignments")
	}

	// Create trash record
	trash := NewTaxonomyTrash()
	trash.SetID(taxonomy.ID())
	trash.SetName(taxonomy.Name())
	trash.SetSlug(taxonomy.Slug())
	trash.SetDescription(taxonomy.Description())
	trash.SetParentID(taxonomy.ParentID())
	trash.SetEntityTypes(taxonomy.EntityTypes())
	trash.SetCreatedAt(taxonomy.CreatedAt())
	trash.SetUpdatedAt(taxonomy.UpdatedAt())
	trash.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	trash.SetDeletedBy(deletedBy)

	// Insert into trash
	record := goqu.Record{}
	for k, v := range trash.Data() {
		record[k] = v
	}

	q := goqu.Dialect(st.dbDriverName).Insert(st.taxonomyTrashTableName).Rows(record)
	sqlStr, _, errSql := q.ToSQL()
	if errSql != nil {
		return false, errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err = st.database.Exec(ctx, sqlStr)
	if err != nil {
		return false, err
	}

	// Delete from main table
	return st.TaxonomyDelete(ctx, taxonomyID)
}

// TaxonomyRestore restores a taxonomy from trash
func (st *storeImplementation) TaxonomyRestore(ctx context.Context, taxonomyID string) (bool, error) {
	if !st.taxonomiesEnabled {
		return false, errors.New("taxonomies are not enabled")
	}

	// Find in trash
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

	// Recreate taxonomy
	taxonomy := NewTaxonomy()
	taxonomy.SetID(trash.ID())
	taxonomy.SetName(trash.Name())
	taxonomy.SetSlug(trash.Slug())
	taxonomy.SetDescription(trash.Description())
	taxonomy.SetParentID(trash.ParentID())
	taxonomy.SetEntityTypes(trash.EntityTypes())
	taxonomy.SetCreatedAt(trash.CreatedAt())
	taxonomy.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	if err := st.TaxonomyCreate(ctx, taxonomy); err != nil {
		return false, err
	}

	// Delete from trash
	q := goqu.Dialect(st.dbDriverName).
		Delete(st.taxonomyTrashTableName).
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

// TaxonomyTrashList lists trashed taxonomies
func (st *storeImplementation) TaxonomyTrashList(ctx context.Context, options TaxonomyQueryOptions) ([]TaxonomyTrashInterface, error) {
	if !st.taxonomiesEnabled {
		return nil, errors.New("taxonomies are not enabled")
	}

	q := goqu.Dialect(st.dbDriverName).From(st.taxonomyTrashTableName)

	if len(options.IDs) > 0 {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDs))
	}

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if options.Slug != "" {
		q = q.Where(goqu.C(COLUMN_SLUG).Eq(options.Slug))
	}

	sortByColumn := COLUMN_DELETED_AT
	sortOrder := "desc"

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

	trashMaps, err := st.database.SelectToMapString(ctx, sqlStr)
	if err != nil {
		return nil, err
	}

	var list []TaxonomyTrashInterface
	for _, m := range trashMaps {
		list = append(list, NewTaxonomyTrashFromExistingData(m))
	}

	return list, nil
}
