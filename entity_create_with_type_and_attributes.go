package entitystore

import (
	"context"
	"log"
)

// EntityCreateWithTypeAndAttributes is a shortcut to create an entity with a type and
// a map of initial attributes. IDs are auto-assigned.
func (st *storeImplementation) EntityCreateWithTypeAndAttributes(ctx context.Context, entityType string, attributes map[string]string) (EntityInterface, error) {
	err := st.database.BeginTransaction()
	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			if rbErr := st.database.RollbackTransaction(); rbErr != nil {
				log.Println(rbErr)
			}
		}
	}()

	entity, err := st.EntityCreateWithType(ctx, entityType)
	if err != nil {
		_ = st.database.RollbackTransaction()
		return nil, err
	}

	for k, v := range attributes {
		if _, err := st.AttributeCreateWithKeyAndValue(ctx, entity.ID(), k, v); err != nil {
			_ = st.database.RollbackTransaction()
			return nil, err
		}
	}

	if err = st.database.CommitTransaction(); err != nil {
		_ = st.database.RollbackTransaction()
		return nil, err
	}

	return entity, nil
}
