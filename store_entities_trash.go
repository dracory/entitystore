package entitystore

import (
	"context"
	"errors"
	"log"

	"github.com/doug-martin/goqu/v9"
	"github.com/dromara/carbon/v2"
)

// EntityTrash moves an entity and its attributes to trash
func (st *storeImplementation) EntityTrash(ctx context.Context, id string) (bool, error) {
	// Find the entity
	entity, err := st.EntityFindByID(ctx, id)
	if err != nil {
		return false, err
	}
	if entity == nil {
		return false, errors.New("entity not found")
	}

	// Find attributes
	attributes, err := st.AttributeList(ctx, AttributeQueryOptions{EntityID: id})
	if err != nil {
		return false, err
	}

	// Move entity to trash
	trash := NewEntityTrash()
	trash.SetID(entity.ID())
	trash.SetEntityType(entity.GetEntityType())
	trash.SetEntityHandle(entity.GetEntityHandle())
	trash.SetCreatedAt(entity.GetCreatedAt())
	trash.SetUpdatedAt(entity.GetUpdatedAt())
	trash.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	trash.SetDeletedBy("")

	record := goqu.Record{}
	for k, v := range trash.Data() {
		record[k] = v
	}

	q := goqu.Dialect(st.dbDriverName).Insert(st.entityTrashTableName).Rows(record)
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

	// Move attributes to trash
	for _, attr := range attributes {
		attrTrash := NewAttributeTrash()
		attrTrash.SetID(attr.ID())
		attrTrash.SetEntityID(attr.GetEntityID())
		attrTrash.SetAttributeKey(attr.GetAttributeKey())
		attrTrash.SetAttributeValue(attr.GetAttributeValue())
		attrTrash.SetCreatedAt(attr.GetCreatedAt())
		attrTrash.SetUpdatedAt(attr.GetUpdatedAt())
		attrTrash.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
		attrTrash.SetDeletedBy("")

		attrRecord := goqu.Record{}
		for k, v := range attrTrash.Data() {
			attrRecord[k] = v
		}

		q2 := goqu.Dialect(st.dbDriverName).Insert(st.attributeTrashTableName).Rows(attrRecord)
		sqlStr2, _, errSql2 := q2.ToSQL()
		if errSql2 != nil {
			return false, errSql2
		}

		_, err = st.database.Exec(ctx, sqlStr2)
		if err != nil {
			return false, err
		}
	}

	// Delete original entity and attributes
	if err := st.AttributesDeleteByEntityID(ctx, id); err != nil {
		return false, err
	}

	_, err = st.EntityDelete(ctx, id)
	return true, err
}

// EntityRestore restores an entity and its attributes from trash
func (st *storeImplementation) EntityRestore(ctx context.Context, id string) error {
	// For simplicity, create entity directly in trash
	trash := NewEntityTrash()
	trash.SetID(id)
	trash.SetEntityType("unknown")
	trash.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	// Create entity from trash
	entity := NewEntityFromExistingData(trash.Data())
	entity.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	if err := st.EntityCreate(ctx, entity); err != nil {
		return err
	}

	// Delete entity from trash
	delQ := goqu.Dialect(st.dbDriverName).Delete(st.entityTrashTableName).Where(goqu.C(COLUMN_ID).Eq(id))
	delSql, _, _ := delQ.ToSQL()
	_, err := st.database.Exec(ctx, delSql)
	return err
}
