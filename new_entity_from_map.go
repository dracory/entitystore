package entitystore

import "github.com/dromara/carbon/v2"

func (st *storeImplementation) NewEntityFromMap(entityMap map[string]string) Entity {
	opts := NewEntityOptions{}
	if id, exists := entityMap[COLUMN_ID]; exists {
		opts.ID = id
	}
	if entityType, exists := entityMap[COLUMN_ENTITY_TYPE]; exists {
		opts.Type = entityType
	}
	if entityHandle, exists := entityMap[COLUMN_ENTITY_HANDLE]; exists {
		opts.Handle = entityHandle
	}
	if createdAt, exists := entityMap[COLUMN_CREATED_AT]; exists {
		opts.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	}
	if updatedAt, exists := entityMap[COLUMN_CREATED_AT]; exists {
		opts.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).StdTime()
	}

	return st.NewEntity(opts)
}
