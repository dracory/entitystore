package entitystore

import (
	"context"
	"errors"
	"log"

	"github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
)

// entityTaxonomyRow is used for scanning entity taxonomy query results
type entityTaxonomyRow struct {
	ID         string `db:"id"`
	EntityID   string `db:"entity_id"`
	TaxonomyID string `db:"taxonomy_id"`
	TermID     string `db:"term_id"`
	CreatedAt  string `db:"created_at"`
}

// EntityTaxonomyAssign assigns an entity to a taxonomy term
func (st *storeImplementation) EntityTaxonomyAssign(ctx context.Context, entityID string, taxonomyID string, termID string) error {
	if !st.taxonomiesEnabled {
		return errors.New("taxonomies are not enabled")
	}

	if entityID == "" || taxonomyID == "" || termID == "" {
		return errors.New("entity ID, taxonomy ID, and term ID are all required")
	}

	if err := validateIDLength(entityID, COLUMN_ENTITY_ID); err != nil {
		return err
	}
	if err := validateIDLength(taxonomyID, COLUMN_TAXONOMY_ID); err != nil {
		return err
	}
	if err := validateIDLength(termID, COLUMN_TERM_ID); err != nil {
		return err
	}

	entity, err := st.EntityFindByID(ctx, entityID)
	if err != nil {
		return err
	}
	if entity == nil {
		return errors.New("entity not found")
	}

	taxonomy, err := st.TaxonomyFind(ctx, taxonomyID)
	if err != nil {
		return err
	}
	if taxonomy == nil {
		return errors.New("taxonomy not found")
	}

	term, err := st.TaxonomyTermFind(ctx, termID)
	if err != nil {
		return err
	}
	if term == nil {
		return errors.New("taxonomy term not found")
	}
	if term.GetTaxonomyID() != taxonomyID {
		return errors.New("taxonomy term does not belong to the specified taxonomy")
	}

	existing, err := st.EntityTaxonomyList(ctx, EntityTaxonomyQuery().
		WithEntityID(entityID).
		WithTaxonomyID(taxonomyID).
		WithTermID(termID).
		WithLimit(1))
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return errors.New("entity is already assigned to this taxonomy term")
	}

	assignment := NewEntityTaxonomy()
	assignment.SetEntityID(entityID)
	assignment.SetTaxonomyID(taxonomyID)
	assignment.SetTermID(termID)
	assignment.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	row := map[string]any{}
	for k, v := range assignment.Data() {
		row[k] = v
	}

	if st.GetDebug() {
		log.Println("EntityTaxonomyAssign:", row)
	}

	return st.db.Query().Table(st.entityTaxonomyTableName).Create(row)
}

// EntityTaxonomyRemove removes an entity from a taxonomy term
func (st *storeImplementation) EntityTaxonomyRemove(ctx context.Context, entityID string, taxonomyID string, termID string) error {
	if !st.taxonomiesEnabled {
		return errors.New("taxonomies are not enabled")
	}

	if entityID == "" || taxonomyID == "" || termID == "" {
		return errors.New("entity ID, taxonomy ID, and term ID are all required")
	}

	_, err := st.db.Query().Table(st.entityTaxonomyTableName).Where(COLUMN_ENTITY_ID+" = ? AND "+COLUMN_TAXONOMY_ID+" = ? AND "+COLUMN_TERM_ID+" = ?", entityID, taxonomyID, termID).Delete()
	return err
}

// applyEntityTaxonomyFilters applies the common filter clauses from a
// validated entity-taxonomy query to a query. Shared by EntityTaxonomyList
// and EntityTaxonomyCount to prevent filter drift.
func (st *storeImplementation) applyEntityTaxonomyFilters(q orm.Query, query EntityTaxonomyQueryInterface) orm.Query {
	if query.GetID() != "" {
		q = q.Where(COLUMN_ID+" = ?", query.GetID())
	}

	if query.GetEntityID() != "" {
		q = q.Where(COLUMN_ENTITY_ID+" = ?", query.GetEntityID())
	}

	if len(query.GetEntityIDs()) > 0 {
		ids := make([]any, len(query.GetEntityIDs()))
		for i, id := range query.GetEntityIDs() {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_ENTITY_ID, ids)
	}

	if query.GetTaxonomyID() != "" {
		q = q.Where(COLUMN_TAXONOMY_ID+" = ?", query.GetTaxonomyID())
	}

	if query.GetTermID() != "" {
		q = q.Where(COLUMN_TERM_ID+" = ?", query.GetTermID())
	}

	if len(query.GetTermIDs()) > 0 {
		ids := make([]any, len(query.GetTermIDs()))
		for i, id := range query.GetTermIDs() {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_TERM_ID, ids)
	}

	q = applyTimeRange(q, COLUMN_CREATED_AT, query.GetCreatedAtGte(), query.GetCreatedAtLte())

	return q
}

// EntityTaxonomyList lists entity-taxonomy assignments matching the given fluent query
func (st *storeImplementation) EntityTaxonomyList(ctx context.Context, query EntityTaxonomyQueryInterface) ([]EntityTaxonomyInterface, error) {
	if !st.taxonomiesEnabled {
		return nil, errors.New("taxonomies are not enabled")
	}

	if query == nil {
		return nil, errors.New("entity taxonomy query cannot be nil")
	}

	if err := query.Validate(); err != nil {
		return nil, err
	}

	q := st.applyEntityTaxonomyFilters(st.db.Query().Table(st.entityTaxonomyTableName), query)

	sortByColumn := COLUMN_CREATED_AT
	sortOrder := "desc"

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

	var rows []entityTaxonomyRow
	if err := q.Get(&rows); err != nil {
		return nil, err
	}

	var list []EntityTaxonomyInterface
	for _, r := range rows {
		list = append(list, NewEntityTaxonomyFromExistingData(map[string]string{
			COLUMN_ID:          r.ID,
			COLUMN_ENTITY_ID:   r.EntityID,
			COLUMN_TAXONOMY_ID: r.TaxonomyID,
			COLUMN_TERM_ID:     r.TermID,
			COLUMN_CREATED_AT:  r.CreatedAt,
		}))
	}

	return list, nil
}

// EntityTaxonomyCount counts entity-taxonomy assignments matching the given fluent query
func (st *storeImplementation) EntityTaxonomyCount(ctx context.Context, query EntityTaxonomyQueryInterface) (int64, error) {
	if !st.taxonomiesEnabled {
		return 0, errors.New("taxonomies are not enabled")
	}

	if query == nil {
		return 0, errors.New("entity taxonomy query cannot be nil")
	}

	if err := query.Validate(); err != nil {
		return 0, err
	}

	q := st.applyEntityTaxonomyFilters(st.db.Query().Table(st.entityTaxonomyTableName), query)

	var count int64
	if err := q.Count(&count); err != nil {
		return 0, err
	}

	return count, nil
}
