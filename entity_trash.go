package entitystore

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
)

// EntityTrash moves an entity and all its attributes to the trash tables, then deletes the originals
func (st *storeImplementation) EntityTrash(ctx context.Context, entityID string) (bool, error) {
	if entityID == "" {
		return false, errors.New("entity ID cannot be empty")
	}

	err := st.database.BeginTransaction()
	if err != nil {
		return false, err
	}

	defer func() {
		if r := recover(); r != nil {
			if rbErr := st.database.RollbackTransaction(); rbErr != nil {
				log.Println(rbErr)
			}
		}
	}()

	ent, err := st.EntityFindByID(ctx, entityID)
	if err != nil {
		_ = st.database.RollbackTransaction()
		return false, err
	}

	if ent == nil {
		_ = st.database.RollbackTransaction()
		return false, nil
	}

	// Insert into entity trash table
	entTrash := EntityTrash{
		ID:        ent.ID(),
		Type:      ent.EntityType(),
		Handle:    ent.EntityHandle(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: time.Now(),
	}

	q := goqu.Dialect(st.dbDriverName).Insert(st.entityTrashTableName).Rows(entTrash)
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

	// Move each attribute to trash
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
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			DeletedAt:      time.Now(),
		}

		q := goqu.Dialect(st.dbDriverName).Insert(st.attributeTrashTableName).Rows(attrTrash)
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

	// Delete attributes then entity
	sqlDelAttrs, _, _ := goqu.Dialect(st.dbDriverName).From(st.attributeTableName).Where(goqu.C(COLUMN_ENTITY_ID).Eq(entityID)).Delete().ToSQL()
	if _, err := st.database.Exec(ctx, sqlDelAttrs); err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		_ = st.database.RollbackTransaction()
		return false, err
	}

	sqlDelEnt, _, _ := goqu.Dialect(st.dbDriverName).From(st.entityTableName).Where(goqu.C(COLUMN_ID).Eq(entityID)).Delete().ToSQL()
	if _, err := st.database.Exec(ctx, sqlDelEnt); err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		_ = st.database.RollbackTransaction()
		return false, err
	}

	if err = st.database.CommitTransaction(); err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		_ = st.database.RollbackTransaction()
		return false, err
	}

	return true, nil
}
