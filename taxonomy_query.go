package entitystore

// TaxonomyQueryOptions defines the query options for taxonomy queries
type TaxonomyQueryOptions struct {
	ID          string
	IDs         []string
	Slug        string
	ParentID    string
	EntityTypes []string
	Limit       uint64
	Offset      uint64
	SortBy      string
	SortOrder   string // asc / desc
	CountOnly   bool
}

// TaxonomyOptions defines the options for creating a taxonomy
type TaxonomyOptions struct {
	Name        string
	Slug        string
	Description string
	ParentID    string
	EntityTypes []string
}

// TaxonomyTermQueryOptions defines the query options for taxonomy term queries
type TaxonomyTermQueryOptions struct {
	ID         string
	IDs        []string
	TaxonomyID string
	Slug       string
	ParentID   string
	Limit      uint64
	Offset     uint64
	SortBy     string
	SortOrder  string // asc / desc
	CountOnly  bool
}

// TaxonomyTermOptions defines the options for creating a taxonomy term
type TaxonomyTermOptions struct {
	TaxonomyID string
	Name       string
	Slug       string
	ParentID   string
	SortOrder  int
}

// EntityTaxonomyQueryOptions defines the query options for entity-taxonomy queries
type EntityTaxonomyQueryOptions struct {
	ID         string
	EntityID   string
	EntityIDs  []string
	TaxonomyID string
	TermID     string
	TermIDs    []string
	Limit      uint64
	Offset     uint64
	SortBy     string
	SortOrder  string // asc / desc
	CountOnly  bool
}
