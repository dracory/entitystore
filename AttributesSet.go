package entitystore

import (
	"context"
	"log"
)

// AttributesSet upserts an entity attribute
func (st *storeImplementation) AttributesSet(ctx context.Context, entityID string, attributes map[string]string) error {
	// err := st.database.BeginTransaction()

	// if err != nil {
	// 	if st.GetDebug() {
	// 		log.Println(err)
	// 	}
	// 	return err
	// }

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		err = st.database.RollbackTransaction()
	// 		if st.GetDebug() {
	// 			log.Println(err)
	// 		}
	// 	}
	// }()

	for k, v := range attributes {
		err := st.AttributeSetString(ctx, entityID, k, v)

		if err != nil {
			if st.GetDebug() {
				log.Println(err)
			}

			// err = st.database.RollbackTransaction()

			// if st.GetDebug() {
			// 	log.Println(err)
			// }

			return err
		}
	}

	// err = st.database.CommitTransaction()

	// if err != nil {
	// 	err = st.database.RollbackTransaction()

	// 	if st.GetDebug() {
	// 		log.Println(err)
	// 	}

	// 	return err
	// }

	return nil
}

// // AttributesSet upserts an entity attribute
// func (st *Store) AttributesSet(entityID string, attributes map[string]string) error {
// 	err := st.database.BeginTransaction()

// 	if err != nil {
// 		if st.GetDebug() {
// 			log.Println(err)
// 		}
// 		return err
// 	}

// 	defer func() {
// 		if r := recover(); r != nil {
// 			err = st.database.RollbackTransaction()
// 			if st.GetDebug() {
// 				log.Println(err)
// 			}
// 		}
// 	}()

// 	for k, v := range attributes {
// 		attr, err := st.AttributeFind(entityID, k)

// 		if err != nil {
// 			if st.GetDebug() {
// 				log.Println(err)
// 			}

// 			err = st.database.RollbackTransaction()

// 			if st.GetDebug() {
// 				log.Println(err)
// 			}

// 			return err
// 		}

// 		if attr == nil {
// 			attr = st.NewAttribute(NewAttributeOptions{ID: uid.HumanUid(), EntityID: entityID, AttributeKey: k, CreatedAt: time.Now(), UpdatedAt: time.Now()})
// 			attr.SetString(v)

// 			q := goqu.Dialect(st.dbDriverName).Insert(st.attributeTableName)
// 			q = q.Rows(attr.ToMap())
// 			sqlStr, _, err := q.ToSQL()

// 			if err != nil {
// 				if st.GetDebug() {
// 					log.Println(err)
// 				}

// 				err = st.database.RollbackTransaction()

// 				if st.GetDebug() {
// 					log.Println(err)
// 				}

// 				return err
// 			}

// 			if st.GetDebug() {
// 				log.Println(sqlStr)
// 			}

// 			_, err = st.database.Exec(sqlStr)

// 			if err != nil {
// 				log.Println(err)
// 				err = st.database.RollbackTransaction()

// 				if st.GetDebug() {
// 					log.Println(err)
// 				}

// 				return err
// 			}

// 		}

// 		attr.SetString(v)
// 		attr.SetUpdatedAt(time.Now())

// 		q := goqu.Dialect(st.dbDriverName).Update(st.attributeTableName)
// 		q = q.Where(goqu.C("id").Eq(attr.ID()))
// 		q = q.Set(attr.ToMap())

// 		sqlStr, _, err := q.ToSQL()

// 		if err != nil {
// 			if st.GetDebug() {
// 				log.Println(err)
// 			}

// 			err = st.database.RollbackTransaction()

// 			if st.GetDebug() {
// 				log.Println(err)
// 			}

// 			return err
// 		}

// 		if st.GetDebug() {
// 			log.Println(sqlStr)
// 		}

// 		_, err = st.database.Exec(sqlStr)

// 		if err != nil {
// 			if st.GetDebug() {
// 				log.Println(err)
// 			}

// 			err = st.database.RollbackTransaction()

// 			if st.GetDebug() {
// 				log.Println(err)
// 			}

// 			return err
// 		}
// 	}

// 	err = st.database.CommitTransaction()

// 	if err != nil {
// 		err = st.database.RollbackTransaction()

// 		if st.GetDebug() {
// 			log.Println(err)
// 		}

// 		return err
// 	}

// 	return nil
// }
