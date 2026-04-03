package store

import (
	"path/filepath"
	"testing"
)

func TestOpenInitializesSchema(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "resume.db")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}
	defer db.Close()

	for _, tableName := range []string{"resume", "visitors"} {
		var existingTable string
		if err := db.QueryRow("SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?", tableName).Scan(&existingTable); err != nil {
			t.Fatalf("expected table %s to exist: %v", tableName, err)
		}
	}
}
