package store

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"testing"

	"gorm.io/gorm"
)

func openResumeDraftTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := Open(filepath.Join(t.TempDir(), "resume.db"))
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db.DB() returned error: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	return db
}

func TestResumeDraftStoreCreateAndList(t *testing.T) {
	db := openResumeDraftTestDB(t)
	store := NewResumeDraftStore(db, NewResumeStore(db))

	firstID, err := store.Create(
		"初版草稿",
		"最早的草稿",
		json.RawMessage(`{"profile":{"name":"初版"}}`),
	)
	if err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	secondID, err := store.Create(
		"终版草稿",
		"最近一次修改",
		json.RawMessage(`{"profile":{"name":"终版"}}`),
	)
	if err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	drafts, err := store.List()
	if err != nil {
		t.Fatalf("List() returned error: %v", err)
	}

	if len(drafts) != 2 {
		t.Fatalf("expected 2 drafts, got %d", len(drafts))
	}

	if drafts[0].ID != secondID {
		t.Fatalf("expected most recent draft %d first, got %d", secondID, drafts[0].ID)
	}

	if drafts[1].ID != firstID {
		t.Fatalf("expected oldest draft %d second, got %d", firstID, drafts[1].ID)
	}

	if drafts[0].Name != "终版草稿" || drafts[0].Note != "最近一次修改" {
		t.Fatalf("expected latest draft metadata to round-trip, got %#v", drafts[0])
	}
}

func TestResumeDraftStoreGetByID(t *testing.T) {
	db := openResumeDraftTestDB(t)
	store := NewResumeDraftStore(db, NewResumeStore(db))

	draftID, err := store.Create(
		"面试版",
		"补充项目亮点",
		json.RawMessage(`{"profile":{"name":"草稿版"}}`),
	)
	if err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	draft, err := store.Get(draftID)
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}

	if draft.ID != draftID {
		t.Fatalf("expected draft id %d, got %d", draftID, draft.ID)
	}

	if draft.Name != "面试版" || draft.Note != "补充项目亮点" {
		t.Fatalf("expected draft metadata to match, got %#v", draft)
	}

	if string(draft.Data) != `{"profile":{"name":"草稿版"}}` {
		t.Fatalf("expected stored draft json to match, got %s", string(draft.Data))
	}
}

func TestResumeDraftStorePublishReplacesMainResume(t *testing.T) {
	db := openResumeDraftTestDB(t)
	store := NewResumeDraftStore(db, NewResumeStore(db))

	draftID, err := store.Create(
		"发布版",
		"准备用于线上",
		json.RawMessage(`{"profile":{"name":"已发布简历"}}`),
	)
	if err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	if err := store.Publish(draftID); err != nil {
		t.Fatalf("Publish() returned error: %v", err)
	}

	mainResume, err := NewResumeStore(db).Get()
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}

	if string(mainResume) != `{"profile":{"name":"已发布简历"}}` {
		t.Fatalf("expected main resume to be replaced, got %s", string(mainResume))
	}
}

func TestResumeDraftStoreDelete(t *testing.T) {
	db := openResumeDraftTestDB(t)
	store := NewResumeDraftStore(db, NewResumeStore(db))

	draftID, err := store.Create(
		"待删除草稿",
		"删除后不应再读取到",
		json.RawMessage(`{"profile":{"name":"待删除"}}`),
	)
	if err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	if err := store.Delete(draftID); err != nil {
		t.Fatalf("Delete() returned error: %v", err)
	}

	drafts, err := store.List()
	if err != nil {
		t.Fatalf("List() returned error: %v", err)
	}

	if len(drafts) != 0 {
		t.Fatalf("expected no drafts after delete, got %d", len(drafts))
	}

	_, err = store.Get(draftID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected Get() after delete to return gorm.ErrRecordNotFound, got %v", err)
	}
}
