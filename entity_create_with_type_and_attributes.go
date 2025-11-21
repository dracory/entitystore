package entitystore

import (
	"context"
	"log"
)

// EntityCreateWithTypeAndAttributes quick shortcut method
// to create an entity by providing only the type as string
// and the attributes as map
// NB. The IDs will be auto-assigned
func (st *storeImplementation) EntityCreateWithTypeAndAttributes(ctx context.Context, entityType string, attributes map[string]string) (*Entity, error) {
	err := st.database.BeginTransaction()

	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			err = st.database.RollbackTransaction()

			if err != nil {
				log.Println(err)
				return
			}
		}
	}()

	entity, err := st.EntityCreateWithType(ctx, entityType)

	if err != nil {
		_ = st.database.RollbackTransaction()
		return nil, err
	}

	for k, v := range attributes {
		_, err := st.AttributeCreateWithKeyAndValue(ctx, entity.ID(), k, v)

		if err != nil {
			_ = st.database.RollbackTransaction()
			return nil, err
		}
	}

	err = st.database.CommitTransaction()

	if err != nil {
		_ = st.database.RollbackTransaction()
		return nil, err
	}

	return entity, nil
}
