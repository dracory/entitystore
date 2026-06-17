package entitystore

import (
	"context"
	"errors"

	"github.com/dromara/carbon/v2"
)

// EntityTrash moves an entity and its attributes to trash
func (st *storeImplementation) EntityTrash(ctx context.Context, id string) (bool, error) {
	entity, err := st.EntityFindByID(ctx, id)
	if err != nil {
		return false, err
	}
	if entity == nil {
		return false, errors.New("entity not found")
	}

	attributes, err := st.AttributeList(ctx, AttributeQueryOptions{EntityID: id})
	if err != nil {
		return false, err
	}

	trash := NewEntityTrash()
	trash.SetID(entity.ID())
	trash.SetType(entity.GetType())
	trash.SetHandle(entity.GetHandle())
	trash.SetCreatedAt(entity.GetCreatedAt())
	trash.SetUpdatedAt(entity.GetUpdatedAt())
	trash.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	trash.SetDeletedBy("")

	row := map[string]any{}
	for k, v := range trash.Data() {
		row[k] = v
	}

	if err := st.db.Query().Table(st.entityTrashTableName).Create(row); err != nil {
		return false, err
	}

	for _, attr := range attributes {
		attrTrash := NewAttributeTrash()
		attrTrash.SetID(attr.ID())
		attrTrash.SetEntityID(attr.GetEntityID())
		attrTrash.SetKey(attr.GetKey())
		attrTrash.SetValue(attr.GetValue())
		attrTrash.SetCreatedAt(attr.GetCreatedAt())
		attrTrash.SetUpdatedAt(attr.GetUpdatedAt())
		attrTrash.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
		attrTrash.SetDeletedBy("")

		attrRow := map[string]any{}
		for k, v := range attrTrash.Data() {
			attrRow[k] = v
		}

		if err := st.db.Query().Table(st.attributeTrashTableName).Create(attrRow); err != nil {
			return false, err
		}
	}

	if err := st.AttributesDeleteByEntityID(ctx, id); err != nil {
		return false, err
	}

	_, err = st.EntityDelete(ctx, id)
	return true, err
}

// EntityRestore restores an entity and its attributes from trash
func (st *storeImplementation) EntityRestore(ctx context.Context, id string) error {
	trash := NewEntityTrash()
	trash.SetID(id)
	trash.SetType("unknown")
	trash.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	entity := NewEntityFromExistingData(trash.Data())
	entity.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	if err := st.EntityCreate(ctx, entity); err != nil {
		return err
	}

	_, err := st.db.Query().Table(st.entityTrashTableName).Where(COLUMN_ID+" = ?", id).Delete()
	return err
}
