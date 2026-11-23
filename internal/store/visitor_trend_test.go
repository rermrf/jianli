package store

import (
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"jianli/internal/model"
)

func TestVisitorAnalyticsUseCompleteBeijingDaysAndDistinctIPs(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "visitors.db"))
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	// Include midnight boundaries, mixed stored offsets, repeated IPs,
	// an empty IP, a zero-visit day and a future calendar day.
	for _, seed := range []struct{ ip, at string }{
		{"192.0.2.9", "2025-12-31T15:59:59Z"},
		{"192.0.2.1", "2025-12-31T16:00:00Z"},
		{"192.0.2.1", "2026-01-01T02:00:00Z"},
		{"192.0.2.3", "2026-01-01T08:00:00+08:00"},
		{"192.0.2.2", "2026-01-02T16:00:00Z"},
		{"192.0.2.1", "2026-01-03T01:00:00Z"},
		{"", "2026-01-03T02:00:00Z"},
		{"192.0.2.8", "2026-01-03T16:00:00Z"},
	} {
		at, err := time.Parse(time.RFC3339, seed.at)
		if err != nil {
			t.Fatal(err)
		}
		if err := db.Create(&model.VisitorRecord{IP: seed.ip, VisitTime: at, Duration: 60}).Error; err != nil {
			t.Fatal(err)
		}
	}
	visitorStore := NewVisitorStore(db)
	now, _ := time.Parse(time.RFC3339, "2026-01-03T03:00:00Z")
	points, err := visitorStore.trendAt(3, now)
	if err != nil {
		t.Fatal(err)
	}
	want := []model.VisitorTrendPoint{
		{Date: "2026-01-01", Count: 3, UniqueVisitors: 2},
		{Date: "2026-01-02", Count: 0, UniqueVisitors: 0},
		{Date: "2026-01-03", Count: 3, UniqueVisitors: 2},
	}
	if !reflect.DeepEqual(points, want) {
		t.Fatalf("daily points=%+v, want %+v", points, want)
	}
	stats, err := visitorStore.statsAt(3, now)
	if err != nil {
		t.Fatal(err)
	}
	if stats.TotalVisits != 6 || stats.TodayVisits != 3 || stats.UniqueVisitors != 3 || stats.AverageDurationSec != 60 {
		t.Fatalf("cards and daily buckets must use the same window: %+v", stats)
	}
	all, err := visitorStore.trendAt(0, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 4 || all[0].Date != "2025-12-31" || all[0].Count != 1 || !reflect.DeepEqual(all[1:], want) {
		t.Fatalf("all-history buckets must start at the earliest visit: %+v", all)
	}
}

func TestVisitorTrendRetainsZeroVisitDays(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "empty.db"))
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	visitorStore := NewVisitorStore(db)
	now, _ := time.Parse(time.RFC3339, "2026-03-01T00:00:00+08:00")
	points, err := visitorStore.trendAt(7, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(points) != 7 || points[0].Date != "2026-02-23" || points[6].Date != "2026-03-01" {
		t.Fatalf("unexpected empty range: %+v", points)
	}
	for _, point := range points {
		if point.Count != 0 || point.UniqueVisitors != 0 {
			t.Fatalf("empty days must retain zero values: %+v", point)
		}
	}
	all, err := visitorStore.trendAt(0, now)
	if err != nil || len(all) != 0 {
		t.Fatalf("empty all-history range=%+v err=%v", all, err)
	}
	if _, err := visitorStore.trendAt(-1, now); err == nil {
		t.Fatal("negative day ranges must be rejected")
	}
}
