package entitystore

import "time"

type NewEntityOptions struct {
	ID        string
	Type      string
	Handle    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (st *storeImplementation) NewEntity(opts NewEntityOptions) Entity {
	entity := Entity{}
	entity.SetID(opts.ID)
	entity.SetType(opts.Type)
	entity.SetHandle(opts.Handle)
	entity.SetCreatedAt(opts.CreatedAt)
	entity.SetUpdatedAt(opts.UpdatedAt)
	entity.st = st
	return entity
}
