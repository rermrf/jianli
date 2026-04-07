package store

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

func TestResumeStoreSeedsOnFirstRead(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "resume.db"))
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db.DB() returned error: %v", err)
	}
	defer sqlDB.Close()

	store := NewResumeStore(db)
	resume, err := store.Get()
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(resume, &payload); err != nil {
		t.Fatalf("expected empty resume JSON, got error: %v", err)
	}

	profile, ok := payload["profile"].(map[string]any)
	if !ok || profile["name"] != "" {
		t.Fatalf("expected empty profile name, got %#v", payload["profile"])
	}

	var count int64
	if err := db.Raw("SELECT COUNT(*) FROM resume").Scan(&count).Error; err != nil {
		t.Fatalf("raw count query returned error: %v", err)
	}

	if count != 0 {
		t.Fatalf("expected Get() on empty db not to persist a default resume, got %d stored rows", count)
	}
}

func TestResumeStoreSaveReplacesExistingRecord(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "resume.db"))
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db.DB() returned error: %v", err)
	}
	defer sqlDB.Close()

	store := NewResumeStore(db)
	updated := json.RawMessage(`{"profile":{"name":"测试用户"}}`)
	if err := store.Save(updated); err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	resume, err := store.Get()
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}

	if string(resume) != string(updated) {
		t.Fatalf("expected saved resume to match update, got %s", string(resume))
	}
}
