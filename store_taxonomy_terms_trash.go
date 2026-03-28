package entitystore

import (
	"context"
	"errors"
	"log"

	"github.com/doug-martin/goqu/v9"
	"github.com/dromara/carbon/v2"
)

// ==========================================
// TaxonomyTerm Trash/Restore
// ==========================================

// TaxonomyTermTrash soft-deletes a taxonomy term by moving it to trash
func (st *storeImplementation) TaxonomyTermTrash(ctx context.Context, termID string, deletedBy string) (bool, error) {
	if !st.taxonomiesEnabled {
		return false, errors.New("taxonomies are not enabled")
	}

	// Find the term
	term, err := st.TaxonomyTermFind(ctx, termID)
	if err != nil {
		return false, err
	}
	if term == nil {
		return false, errors.New("taxonomy term not found")
	}

	// Check for entity assignments
	assignmentsCount, err := st.EntityTaxonomyCount(ctx, EntityTaxonomyQueryOptions{
		TermID: termID,
	})
	if err != nil {
		return false, err
	}
	if assignmentsCount > 0 {
		return false, errors.New("cannot trash taxonomy term: it has associated entity assignments")
	}

	// Create trash record
	trash := NewTaxonomyTermTrash()
	trash.SetID(term.ID())
	trash.SetTaxonomyID(term.TaxonomyID())
	trash.SetName(term.Name())
	trash.SetSlug(term.Slug())
	trash.SetParentID(term.ParentID())
	trash.SetSortOrder(term.SortOrder())
	trash.SetCreatedAt(term.CreatedAt())
	trash.SetUpdatedAt(term.UpdatedAt())
	trash.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	trash.SetDeletedBy(deletedBy)

	// Insert into trash
	record := goqu.Record{}
	for k, v := range trash.Data() {
		record[k] = v
	}

	q := goqu.Dialect(st.dbDriverName).Insert(st.taxonomyTermTrashTableName).Rows(record)
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
	return st.TaxonomyTermDelete(ctx, termID)
}

// TaxonomyTermRestore restores a taxonomy term from trash
func (st *storeImplementation) TaxonomyTermRestore(ctx context.Context, termID string) (bool, error) {
	if !st.taxonomiesEnabled {
		return false, errors.New("taxonomies are not enabled")
	}

	// Find in trash
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

	// Recreate term
	term := NewTaxonomyTerm()
	term.SetID(trash.ID())
	term.SetTaxonomyID(trash.TaxonomyID())
	term.SetName(trash.Name())
	term.SetSlug(trash.Slug())
	term.SetParentID(trash.ParentID())
	term.SetSortOrder(trash.SortOrder())
	term.SetCreatedAt(trash.CreatedAt())
	term.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	if err := st.TaxonomyTermCreate(ctx, term); err != nil {
		return false, err
	}

	// Delete from trash
	q := goqu.Dialect(st.dbDriverName).
		Delete(st.taxonomyTermTrashTableName).
		Where(goqu.C(COLUMN_ID).Eq(termID))

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

// TaxonomyTermTrashList lists trashed taxonomy terms
func (st *storeImplementation) TaxonomyTermTrashList(ctx context.Context, options TaxonomyTermQueryOptions) ([]TaxonomyTermTrashInterface, error) {
	if !st.taxonomiesEnabled {
		return nil, errors.New("taxonomies are not enabled")
	}

	q := goqu.Dialect(st.dbDriverName).From(st.taxonomyTermTrashTableName)

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

	var list []TaxonomyTermTrashInterface
	for _, m := range trashMaps {
		list = append(list, NewTaxonomyTermTrashFromExistingData(m))
	}

	return list, nil
}
