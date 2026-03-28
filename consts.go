package entitystore

// Column constants
const (
	COLUMN_ID                = "id"
	COLUMN_ENTITY_TYPE       = "entity_type"
	COLUMN_ENTITY_HANDLE     = "entity_handle"
	COLUMN_ENTITY_ID         = "entity_id"
	COLUMN_ATTRIBUTE_KEY     = "attribute_key"
	COLUMN_ATTRIBUTE_VALUE   = "attribute_value"
	COLUMN_CREATED_AT        = "created_at"
	COLUMN_UPDATED_AT        = "updated_at"
	COLUMN_DELETED_AT        = "deleted_at"
	COLUMN_DELETED_BY        = "deleted_by"
	COLUMN_RELATED_ENTITY_ID = "related_entity_id"
	COLUMN_RELATIONSHIP_TYPE = "relationship_type"
	COLUMN_PARENT_ID         = "parent_id"
	COLUMN_SEQUENCE          = "sequence"
	COLUMN_METADATA          = "metadata"

	// Taxonomy columns
	COLUMN_NAME         = "name"
	COLUMN_SLUG         = "slug"
	COLUMN_DESCRIPTION  = "description"
	COLUMN_ENTITY_TYPES = "entity_types"
	COLUMN_TAXONOMY_ID  = "taxonomy_id"
	COLUMN_TERM_ID      = "term_id"
	COLUMN_SORT_ORDER   = "sort_order"
)

// Relationship types
const (
	RELATIONSHIP_TYPE_BELONGS_TO = "belongs_to"   // Entity belongs to one parent
	RELATIONSHIP_TYPE_HAS_MANY   = "has_many"     // Entity has many children
	RELATIONSHIP_TYPE_MANY_MANY  = "many_to_many" // Entities linked bidirectionally
)

// Default table names
const (
	DEFAULT_RELATIONSHIP_TABLE_NAME       = "entities_relationships"
	DEFAULT_RELATIONSHIP_TRASH_TABLE_NAME = "entities_relationships_trash"

	// Taxonomy table names
	DEFAULT_TAXONOMY_TABLE_NAME            = "entities_taxonomies"
	DEFAULT_TAXONOMY_TERM_TABLE_NAME       = "entities_taxonomy_terms"
	DEFAULT_ENTITY_TAXONOMY_TABLE_NAME     = "entities_entity_taxonomies"
	DEFAULT_TAXONOMY_TRASH_TABLE_NAME      = "entities_taxonomies_trash"
	DEFAULT_TAXONOMY_TERM_TRASH_TABLE_NAME = "entities_taxonomy_terms_trash"
)
