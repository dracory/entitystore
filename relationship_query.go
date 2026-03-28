package entitystore

// RelationshipQueryOptions defines the query options for relationship queries
type RelationshipQueryOptions struct {
	ID               string
	IDs              []string
	EntityID         string
	EntityIDs        []string
	RelatedEntityID  string
	RelatedEntityIDs []string
	RelationshipType string
	ParentID         string
	Limit            uint64
	Offset           uint64
	SortBy           string
	SortOrder        string // asc / desc
	CountOnly        bool
}

// RelationshipOptions defines the options for creating a relationship
type RelationshipOptions struct {
	EntityID         string
	RelatedEntityID  string
	RelationshipType string
	ParentID         string
	Sequence         int
	Metadata         string
}
