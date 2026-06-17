package entitystore

import (
	"context"

	"github.com/dromara/carbon/v2"
)

// AttributeTrash moves an attribute to trash
func (st *storeImplementation) AttributeTrash(ctx context.Context, id string, deletedBy string) error {
	trash := NewAttributeTrash()
	trash.SetID(id)
	trash.SetEntityID("")
	trash.SetKey("unknown")
	trash.SetValue("")
	trash.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	trash.SetDeletedBy(deletedBy)

	row := map[string]any{}
	for k, v := range trash.Data() {
		row[k] = v
	}

	if err := st.db.Query().Table(st.attributeTrashTableName).Create(row); err != nil {
		return err
	}

	return st.AttributeDelete(ctx, id)
}

// AttributeRestore restores an attribute from trash
func (st *storeImplementation) AttributeRestore(ctx context.Context, id string) error {
	trash := NewAttributeTrash()
	trash.SetID(id)

	attr := NewAttributeFromExistingData(trash.Data())
	attr.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	if err := st.AttributeCreate(ctx, attr); err != nil {
		return err
	}

	_, err := st.db.Query().Table(st.attributeTrashTableName).Where(COLUMN_ID+" = ?", id).Delete()
	return err
}
