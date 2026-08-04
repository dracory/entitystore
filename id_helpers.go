package entitystore

import (
	"fmt"

	neatuid "github.com/dracory/neat/support/uid"
)

// GenerateShortID delegates to neat/support/uid which produces 11-character
// lowercase Crockford Base32 IDs from a microsecond timestamp + 4-bit counter.
// Thread-safe via internal mutex.
func GenerateShortID() string {
	return neatuid.GenerateShortID()
}

// validateIDLength checks that an ID does not exceed ID_COLUMN_LENGTH.
// Returns an error if the ID is too long, nil otherwise.
// len() in Go is O(1) — the length is stored in the string header.
func validateIDLength(id, columnName string) error {
	if len(id) > ID_COLUMN_LENGTH {
		return fmt.Errorf("column %q value %q exceeds maximum length %d (got %d)", columnName, id, ID_COLUMN_LENGTH, len(id))
	}
	return nil
}

// validateIDRequired checks that a primary key ID is not empty.
// Use for COLUMN_ID fields that must always have a value.
func validateIDRequired(id, columnName string) error {
	if id == "" {
		return fmt.Errorf("column %q is required and cannot be empty", columnName)
	}
	return nil
}
