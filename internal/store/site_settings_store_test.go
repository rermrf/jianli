package store

import (
	"path/filepath"
	"testing"
)

func TestSiteSettingsStoreReturnsDefaultRecord(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "resume.db"))
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db.DB() returned error: %v", err)
	}
	defer sqlDB.Close()

	store := NewSiteSettingsStore(db)
	settings, err := store.Get()
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}

	if !settings.AllowPDFExport {
		t.Fatal("expected allowPdfExport to default to true")
	}
}

func TestSiteSettingsStoreSavePersistsUpdatedValue(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "resume.db"))
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db.DB() returned error: %v", err)
	}
	defer sqlDB.Close()

	store := NewSiteSettingsStore(db)
	if err := store.Save(false); err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	settings, err := store.Get()
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}

	if settings.AllowPDFExport {
		t.Fatal("expected allowPdfExport to persist false")
	}
}
