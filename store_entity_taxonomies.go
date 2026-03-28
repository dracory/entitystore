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
// EntityTaxonomy CRUD (Entity-Term Assignments)
// ==========================================

// EntityTaxonomyAssign assigns an entity to a taxonomy term
func (st *storeImplementation) EntityTaxonomyAssign(ctx context.Context, entityID string, taxonomyID string, termID string) error {
	if !st.taxonomiesEnabled {
		return errors.New("taxonomies are not enabled")
	}

	if entityID == "" || taxonomyID == "" || termID == "" {
		return errors.New("entity ID, taxonomy ID, and term ID are all required")
	}

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

	// Check if assignment already exists
	existing, err := st.EntityTaxonomyList(ctx, EntityTaxonomyQueryOptions{
		EntityID:   entityID,
		TaxonomyID: taxonomyID,
		TermID:     termID,
		Limit:      1,
	})
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return errors.New("entity is already assigned to this taxonomy term")
	}

	// Create new assignment
	assignment := NewEntityTaxonomy()
	assignment.SetEntityID(entityID)
	assignment.SetTaxonomyID(taxonomyID)
	assignment.SetTermID(termID)
	assignment.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	record := goqu.Record{}
	for k, v := range assignment.Data() {
		record[k] = v
	}

	q := goqu.Dialect(st.dbDriverName).Insert(st.entityTaxonomyTableName).Rows(record)

	sqlStr, _, errSql := q.ToSQL()
	if errSql != nil {
		return errSql
	}

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err = st.database.Exec(ctx, sqlStr)
	return err
}

// EntityTaxonomyRemove removes an entity from a taxonomy term
func (st *storeImplementation) EntityTaxonomyRemove(ctx context.Context, entityID string, taxonomyID string, termID string) error {
	if !st.taxonomiesEnabled {
		return errors.New("taxonomies are not enabled")
	}

	if entityID == "" || taxonomyID == "" || termID == "" {
		return errors.New("entity ID, taxonomy ID, and term ID are all required")
	}

	q := goqu.Dialect(st.dbDriverName).
		Delete(st.entityTaxonomyTableName).
		Where(
			goqu.C(COLUMN_ENTITY_ID).Eq(entityID),
			goqu.C(COLUMN_TAXONOMY_ID).Eq(taxonomyID),
			goqu.C(COLUMN_TERM_ID).Eq(termID),
		)

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

// EntityTaxonomyList lists entity-taxonomy assignments matching the given query options
func (st *storeImplementation) EntityTaxonomyList(ctx context.Context, options EntityTaxonomyQueryOptions) ([]EntityTaxonomyInterface, error) {
	if !st.taxonomiesEnabled {
		return nil, errors.New("taxonomies are not enabled")
	}

	q := goqu.Dialect(st.dbDriverName).From(st.entityTaxonomyTableName)

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if options.EntityID != "" {
		q = q.Where(goqu.C(COLUMN_ENTITY_ID).Eq(options.EntityID))
	}

	if len(options.EntityIDs) > 0 {
		q = q.Where(goqu.C(COLUMN_ENTITY_ID).In(options.EntityIDs))
	}

	if options.TaxonomyID != "" {
		q = q.Where(goqu.C(COLUMN_TAXONOMY_ID).Eq(options.TaxonomyID))
	}

	if options.TermID != "" {
		q = q.Where(goqu.C(COLUMN_TERM_ID).Eq(options.TermID))
	}

	if len(options.TermIDs) > 0 {
		q = q.Where(goqu.C(COLUMN_TERM_ID).In(options.TermIDs))
	}

	sortByColumn := COLUMN_CREATED_AT
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

	assignmentMaps, err := st.database.SelectToMapString(ctx, sqlStr)
	if err != nil {
		return nil, err
	}

	var list []EntityTaxonomyInterface
	for _, m := range assignmentMaps {
		list = append(list, NewEntityTaxonomyFromExistingData(m))
	}

	return list, nil
}

// EntityTaxonomyCount counts entity-taxonomy assignments matching the given options
func (st *storeImplementation) EntityTaxonomyCount(ctx context.Context, options EntityTaxonomyQueryOptions) (int64, error) {
	if !st.taxonomiesEnabled {
		return 0, errors.New("taxonomies are not enabled")
	}

	q := goqu.Dialect(st.dbDriverName).From(st.entityTaxonomyTableName)

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if options.EntityID != "" {
		q = q.Where(goqu.C(COLUMN_ENTITY_ID).Eq(options.EntityID))
	}

	if len(options.EntityIDs) > 0 {
		q = q.Where(goqu.C(COLUMN_ENTITY_ID).In(options.EntityIDs))
	}

	if options.TaxonomyID != "" {
		q = q.Where(goqu.C(COLUMN_TAXONOMY_ID).Eq(options.TaxonomyID))
	}

	if options.TermID != "" {
		q = q.Where(goqu.C(COLUMN_TERM_ID).Eq(options.TermID))
	}

	if len(options.TermIDs) > 0 {
		q = q.Where(goqu.C(COLUMN_TERM_ID).In(options.TermIDs))
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
