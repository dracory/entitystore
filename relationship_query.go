package entitystore

// RelationshipQueryOptions provides filtering and pagination options for relationship queries
type RelationshipQueryOptions struct {
	ID               string   // Filter by specific relationship ID
	IDs              []string // Filter by multiple relationship IDs
	EntityID         string   // Filter by source entity ID
	EntityIDs        []string // Filter by multiple source entity IDs
	RelatedEntityID  string   // Filter by target entity ID
	RelatedEntityIDs []string // Filter by multiple target entity IDs
	RelationshipType string   // Filter by relationship type
	ParentID         string   // Filter by parent relationship ID
	Limit            uint64   // Maximum number of results to return
	Offset           uint64   // Number of results to skip
	SortBy           string   // Column to sort by (default: id)
	SortOrder        string   // Sort direction: "asc" or "desc"
	CountOnly        bool     // Return only count, not results
}

// RelationshipOptions provides the options for creating a new relationship
type RelationshipOptions struct {
	EntityID         string // Source entity ID
	RelatedEntityID  string // Target/related entity ID
	RelationshipType string // Type of relationship (e.g., "belongs_to", "has_many")
	ParentID         string // Parent relationship ID for hierarchical relationships
	Sequence         int    // Sort order for the relationship
	Metadata         string // JSON metadata associated with the relationship
}
