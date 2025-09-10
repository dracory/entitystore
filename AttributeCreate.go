package entitystore

import (
	"errors"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/uid"
)

// AttributeCreate creates a new attribute
func (st *storeImplementation) AttributeCreate(attr *Attribute) error {
	if attr == nil {
		return errors.New("attribute is required")
	}

	if attr.AttributeKey() == "" {
		return errors.New("attribute key is required field")
	}

	if attr.ID() == "" {
		attr.SetID(uid.HumanUid())
	}

	if attr.CreatedAt().IsZero() {
		attr.SetCreatedAt(time.Now())
	}

	if attr.UpdatedAt().IsZero() {
		attr.SetUpdatedAt(time.Now())
	}

	q := goqu.Dialect(st.dbDriverName).Insert(st.attributeTableName)
	q = q.Rows(attr.ToMap())
	sqlStr, _, _ := q.ToSQL()

	if st.GetDebug() {
		log.Println(sqlStr)
	}

	_, err := st.database.Exec(sqlStr)

	if err != nil {
		if st.GetDebug() {
			log.Println(err)
		}
		return err
	}

	return nil
}
