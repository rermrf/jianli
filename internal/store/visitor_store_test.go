package store

import (
	"path/filepath"
	"testing"
	"time"

	"jianli/internal/model"
)

func TestVisitorStoreDedupesRecentIP(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "resume.db"))
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}
	defer db.Close()

	store := NewVisitorStore(db)
	visitTime := time.Now().Add(-5 * time.Minute)
	input := model.VisitorRecord{IP: "112.17.45.201", City: "杭州", Device: "Windows", Browser: "Chrome", OS: "Windows", VisitTime: visitTime}

	firstID, err := store.RecordVisit(input)
	if err != nil {
		t.Fatalf("RecordVisit() returned error: %v", err)
	}

	secondID, err := store.RecordVisit(input)
	if err != nil {
		t.Fatalf("RecordVisit() returned error: %v", err)
	}

	if firstID != secondID {
		t.Fatalf("expected duplicate IP within 10 minutes to reuse row, got %d and %d", firstID, secondID)
	}
}

func TestVisitorStoreStatsAndListRespectDaysFilter(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "resume.db"))
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}
	defer db.Close()

	store := NewVisitorStore(db)
	records := []model.VisitorRecord{
		{IP: "112.17.45.201", City: "杭州", Device: "Windows", Browser: "Chrome", OS: "Windows", VisitTime: time.Now().Add(-2 * time.Hour), Duration: 120},
		{IP: "183.62.84.19", City: "深圳", Device: "iPhone", Browser: "Safari", OS: "iOS", VisitTime: time.Now().Add(-26 * time.Hour), Duration: 60},
		{IP: "59.56.11.8", City: "厦门", Device: "Mac", Browser: "Firefox", OS: "macOS", VisitTime: time.Now().Add(-4 * 24 * time.Hour), Duration: 180},
	}

	for _, record := range records {
		if _, err := store.RecordVisit(record); err != nil {
			t.Fatalf("RecordVisit() returned error: %v", err)
		}
	}

	list, err := store.List(3, 1, 20)
	if err != nil {
		t.Fatalf("List() returned error: %v", err)
	}

	if len(list) != 2 {
		t.Fatalf("expected 2 records in 3-day window, got %d", len(list))
	}

	stats, err := store.Stats(3)
	if err != nil {
		t.Fatalf("Stats() returned error: %v", err)
	}

	if stats.TotalVisits != 2 {
		t.Fatalf("expected total visits 2, got %d", stats.TotalVisits)
	}

	if stats.UniqueVisitors != 2 {
		t.Fatalf("expected unique visitors 2, got %d", stats.UniqueVisitors)
	}
}
