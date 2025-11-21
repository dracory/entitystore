package entitystore

import (
	"context"
	"errors"
	"log"

	"github.com/doug-martin/goqu/v9"
)

// EntityDelete deletes an entity and all attributes
func (st *storeImplementation) EntityDelete(ctx context.Context, entityID string) (bool, error) {
	if entityID == "" {
		if st.GetDebug() {
			log.Println("in EntityDelete entity ID cannot be empty")
		}
		return false, errors.New("in EntityDelete entity ID cannot be empty")
	}

	// Note the use of tx as the database handle once you are within a transaction
	err := st.database.BeginTransaction()

	if err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		return false, err
	}

	defer func() {
		if r := recover(); r != nil {
			txErr := st.database.RollbackTransaction()
			if txErr != nil && st.GetDebug() {
				log.Println(txErr)
			}
		}
	}()

	sqlStr1, _, _ := goqu.Dialect(st.dbDriverName).From(st.attributeTableName).Where(goqu.C("entity_id").Eq(entityID)).Delete().ToSQL()

	if _, err := st.database.Exec(ctx, sqlStr1); err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		txErr := st.database.RollbackTransaction()
		if txErr != nil && st.GetDebug() {
			log.Println(txErr)
		}
		return false, err
	}

	sqlStr2, _, _ := goqu.Dialect(st.dbDriverName).From(st.entityTableName).Where(goqu.C("id").Eq(entityID)).Delete().ToSQL()

	if _, err := st.database.Exec(ctx, sqlStr2); err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		err := st.database.RollbackTransaction()
		if err != nil && st.GetDebug() {
			log.Println(err)
		}
		return false, err
	}

	err = st.database.CommitTransaction()

	if err != nil {
		if st.GetDebug() {
			log.Println(err)
		}

		err := st.database.RollbackTransaction()
		if err != nil && st.GetDebug() {
			log.Println(err)
		}

		return false, err
	}

	return true, nil
}
