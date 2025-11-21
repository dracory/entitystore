package entitystore

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
)

// EntityTrash moves an entity and all attributes to the trash bin
func (st *storeImplementation) EntityTrash(ctx context.Context, entityID string) (bool, error) {
	if entityID == "" {
		return false, errors.New("entity ID cannot be empty")
	}

	_ = ctx
	// Note the use of tx as the database handle once you are within a transaction
	err := st.database.BeginTransaction()

	defer func() {
		if r := recover(); r != nil {
			err = st.database.RollbackTransaction()
			if err != nil {
				log.Println(err)
			}
		}
	}()

	if err != nil {
		return false, err
	}

	ent, err := st.EntityFindByID(ctx, entityID)

	if err != nil {
		_ = st.database.RollbackTransaction()
		return false, err
	}

	if ent == nil {
		_ = st.database.RollbackTransaction()
		return false, err
	}

	entTrash := EntityTrash{
		ID:        ent.ID(),
		Type:      ent.Type(),
		CreatedAt: ent.CreatedAt(),
		UpdatedAt: ent.UpdatedAt(),
		DeletedAt: time.Now(),
	}

	q := goqu.Dialect(st.dbDriverName).Insert(st.entityTrashTableName)
	q = q.Rows(entTrash)
	sqlStr, _, _ := q.ToSQL()

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	if _, err := st.database.Exec(ctx, sqlStr); err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		_ = st.database.RollbackTransaction()
		return false, err
	}

	attrs, err := st.EntityAttributeList(ctx, entityID)

	if err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		_ = st.database.RollbackTransaction()
		return false, err
	}

	for _, attr := range attrs {
		attrTrash := AttributeTrash{
			ID:             attr.ID(),
			EntityID:       attr.EntityID(),
			AttributeKey:   attr.AttributeKey(),
			AttributeValue: attr.AttributeValue(),
			CreatedAt:      attr.CreatedAt(),
			UpdatedAt:      attr.UpdatedAt(),
			DeletedAt:      time.Now(),
		}

		q := goqu.Dialect(st.dbDriverName).Insert(st.attributeTrashTableName)
		q = q.Rows(attrTrash)
		sqlStrAttr, _, _ := q.ToSQL()

		if st.GetDebug() {
			log.Println(sqlStrAttr)
		}

		if _, err := st.database.Exec(ctx, sqlStrAttr); err != nil {
			if st.GetDebug() {
				log.Println(err)
			}
			_ = st.database.RollbackTransaction()
			return false, err
		}
	}

	q1 := goqu.Dialect(st.dbDriverName).From(st.attributeTableName).Where(goqu.C(COLUMN_ENTITY_ID).Eq(entityID)).Delete()
	sqlStr1, _, _ := q1.ToSQL()

	if _, err := st.database.Exec(ctx, sqlStr1); err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		_ = st.database.RollbackTransaction()
		return false, err
	}

	q2 := goqu.Dialect(st.dbDriverName).From(st.entityTableName).Where(goqu.C(COLUMN_ID).Eq(entityID)).Delete()
	sqlStr2, _, _ := q2.ToSQL()

	if _, err := st.database.Exec(ctx, sqlStr2); err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		_ = st.database.RollbackTransaction()
		return false, err
	}

	err = st.database.CommitTransaction()

	if err == nil {
		return true, nil
	}

	if st.GetDebug() {
		log.Println(err)
	}

	err = st.database.RollbackTransaction()
	if err != nil && st.GetDebug() {
		log.Println(err)
	}

	return false, err
}
