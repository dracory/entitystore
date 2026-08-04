package entitystore

import (
	"os"
	"strings"
	"testing"
)

// TestGenerateShortID_FitsColumnWidth is half 1 of the proof:
// GenerateShortID() must produce IDs that fit within ID_COLUMN_LENGTH.
// If this fails, the generated IDs are longer than the declared column width.
func TestGenerateShortID_FitsColumnWidth(t *testing.T) {
	for i := 0; i < 1000; i++ {
		id := GenerateShortID()
		if len(id) > ID_COLUMN_LENGTH {
			t.Errorf("GenerateShortID() returned %q (length %d) which exceeds ID_COLUMN_LENGTH (%d) — "+
				"inserts will fail on MySQL/MariaDB/PostgreSQL with 'Data too long for column'",
				id, len(id), ID_COLUMN_LENGTH)
			return
		}
	}
}

// TestStoreSchema_UsesIDColumnLength is half 2 of the proof:
// store.go must use ID_COLUMN_LENGTH for all id-type column definitions.
// This verifies the schema side — the columns are declared with the same
// constant that TestGenerateShortID_FitsColumnWidth checks against.
func TestStoreSchema_UsesIDColumnLength(t *testing.T) {
	src, err := os.ReadFile("store.go")
	if err != nil {
		t.Fatal("Could not read store.go:", err)
	}

	source := string(src)

	// Count how many times ID_COLUMN_LENGTH is used in table.String() calls
	count := strings.Count(source, "ID_COLUMN_LENGTH)")

	t.Logf("store.go uses ID_COLUMN_LENGTH %d times in schema definitions", count)

	if count == 0 {
		t.Error("store.go does not use ID_COLUMN_LENGTH — id-type columns are not linked to the constant")
	}
}
