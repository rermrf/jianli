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
		t.Fatalf("expected seeded resume JSON, got error: %v", err)
	}

	profile, ok := payload["profile"].(map[string]any)
	if !ok || profile["name"] != "温庆京" {
		t.Fatalf("expected seeded profile name 温庆京, got %#v", payload["profile"])
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
