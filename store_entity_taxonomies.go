package entitystore

import (
	"context"
	"errors"
	"log"

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

// EntityTaxonomyList lists entity-taxonomy assignments matching the given query options
func (st *storeImplementation) EntityTaxonomyList(ctx context.Context, options EntityTaxonomyQueryOptions) ([]EntityTaxonomyInterface, error) {
	if !st.taxonomiesEnabled {
		return nil, errors.New("taxonomies are not enabled")
	}

	q := st.db.Query().Table(st.entityTaxonomyTableName)

	if options.ID != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID)
	}

	if options.EntityID != "" {
		q = q.Where(COLUMN_ENTITY_ID+" = ?", options.EntityID)
	}

	if len(options.EntityIDs) > 0 {
		ids := make([]any, len(options.EntityIDs))
		for i, id := range options.EntityIDs {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_ENTITY_ID, ids)
	}

	if options.TaxonomyID != "" {
		q = q.Where(COLUMN_TAXONOMY_ID+" = ?", options.TaxonomyID)
	}

	if options.TermID != "" {
		q = q.Where(COLUMN_TERM_ID+" = ?", options.TermID)
	}

	if len(options.TermIDs) > 0 {
		ids := make([]any, len(options.TermIDs))
		for i, id := range options.TermIDs {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_TERM_ID, ids)
	}

	sortByColumn := COLUMN_CREATED_AT
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

// EntityTaxonomyCount counts entity-taxonomy assignments matching the given options
func (st *storeImplementation) EntityTaxonomyCount(ctx context.Context, options EntityTaxonomyQueryOptions) (int64, error) {
	if !st.taxonomiesEnabled {
		return 0, errors.New("taxonomies are not enabled")
	}

	q := st.db.Query().Table(st.entityTaxonomyTableName)

	if options.ID != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID)
	}

	if options.EntityID != "" {
		q = q.Where(COLUMN_ENTITY_ID+" = ?", options.EntityID)
	}

	if len(options.EntityIDs) > 0 {
		ids := make([]any, len(options.EntityIDs))
		for i, id := range options.EntityIDs {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_ENTITY_ID, ids)
	}

	if options.TaxonomyID != "" {
		q = q.Where(COLUMN_TAXONOMY_ID+" = ?", options.TaxonomyID)
	}

	if options.TermID != "" {
		q = q.Where(COLUMN_TERM_ID+" = ?", options.TermID)
	}

	if len(options.TermIDs) > 0 {
		ids := make([]any, len(options.TermIDs))
		for i, id := range options.TermIDs {
			ids[i] = id
		}
		q = q.WhereIn(COLUMN_TERM_ID, ids)
	}

	var count int64
	if err := q.Count(&count); err != nil {
		return 0, err
	}

	return count, nil
}
