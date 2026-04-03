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

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db.DB() returned error: %v", err)
	}
	defer sqlDB.Close()

	for _, tableName := range []string{"resume", "visitors"} {
		if !db.Migrator().HasTable(tableName) {
			t.Fatalf("expected table %s to exist", tableName)
		}
	}
}
