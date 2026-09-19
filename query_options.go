package entitystore

// RelationshipOptions provides the options for creating a new relationship
type RelationshipOptions struct {
	EntityID         string // Source entity ID
	RelatedEntityID  string // Target/related entity ID
	RelationshipType string // Type of relationship (e.g., "belongs_to", "has_many")
	ParentID         string // Parent relationship ID for hierarchical relationships
	Sequence         int    // Sort order for the relationship
	Metadata         string // JSON metadata associated with the relationship
}
