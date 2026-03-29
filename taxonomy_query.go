package entitystore

// TaxonomyQueryOptions provides filtering and pagination options for taxonomy queries
type TaxonomyQueryOptions struct {
	ID          string   // Filter by specific taxonomy ID
	IDs         []string // Filter by multiple taxonomy IDs
	Slug        string   // Filter by taxonomy slug
	ParentID    string   // Filter by parent taxonomy ID
	EntityTypes []string // Filter by allowed entity types
	Limit       uint64   // Maximum number of results to return
	Offset      uint64   // Number of results to skip
	SortBy      string   // Column to sort by (default: id)
	SortOrder   string   // Sort direction: "asc" or "desc"
	CountOnly   bool     // Return only count, not results
}

// TaxonomyOptions provides the options for creating a new taxonomy
type TaxonomyOptions struct {
	Name        string   // Display name of the taxonomy
	Slug        string   // URL-friendly identifier
	Description string   // Optional description
	ParentID    string   // Parent taxonomy ID for hierarchical taxonomies
	EntityTypes []string // Entity types that can use this taxonomy
}

// TaxonomyTermQueryOptions provides filtering and pagination options for taxonomy term queries
type TaxonomyTermQueryOptions struct {
	ID         string   // Filter by specific term ID
	IDs        []string // Filter by multiple term IDs
	TaxonomyID string   // Filter by parent taxonomy ID
	Slug       string   // Filter by term slug
	ParentID   string   // Filter by parent term ID
	Limit      uint64   // Maximum number of results to return
	Offset     uint64   // Number of results to skip
	SortBy     string   // Column to sort by (default: id)
	SortOrder  string   // Sort direction: "asc" or "desc"
	CountOnly  bool     // Return only count, not results
}

// TaxonomyTermOptions provides the options for creating a new taxonomy term
type TaxonomyTermOptions struct {
	TaxonomyID string // Parent taxonomy ID
	Name       string // Display name of the term
	Slug       string // URL-friendly identifier
	ParentID   string // Parent term ID for hierarchical terms
	SortOrder  int    // Sort order for the term
}

// EntityTaxonomyQueryOptions provides filtering and pagination options for entity-taxonomy queries
type EntityTaxonomyQueryOptions struct {
	ID         string   // Filter by specific assignment ID
	EntityID   string   // Filter by entity ID
	EntityIDs  []string // Filter by multiple entity IDs
	TaxonomyID string   // Filter by taxonomy ID
	TermID     string   // Filter by term ID
	TermIDs    []string // Filter by multiple term IDs
	Limit      uint64   // Maximum number of results to return
	Offset     uint64   // Number of results to skip
	SortBy     string   // Column to sort by (default: id)
	SortOrder  string   // Sort direction: "asc" or "desc"
	CountOnly  bool     // Return only count, not results
}
