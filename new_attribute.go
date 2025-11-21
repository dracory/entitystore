package entitystore

import "time"

type NewAttributeOptions struct {
	ID             string
	EntityID       string
	AttributeKey   string
	AttributeValue string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (st *storeImplementation) NewAttribute(opts NewAttributeOptions) Attribute {
	attribute := Attribute{}
	attribute.SetID(opts.ID)
	attribute.SetEntityID(opts.EntityID)
	attribute.SetAttributeKey(opts.AttributeKey)
	attribute.SetAttributeValue(opts.AttributeValue)
	attribute.SetCreatedAt(opts.CreatedAt)
	attribute.SetUpdatedAt(opts.UpdatedAt)
	attribute.st = st
	return attribute
}
