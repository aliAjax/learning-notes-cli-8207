package search

import (
	"testing"
	"time"

	"learning-notes-cli/internal/model"
)

func testNotes() []model.Note {
	now := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)
	return []model.Note{
		{ID: "1", Title: "Go goroutines", Tags: []string{"go", "concurrency"}, CreatedAt: now, UpdatedAt: now},
		{ID: "2", Title: "Learning log", Tags: []string{"daily", "go"}, CreatedAt: now, UpdatedAt: now},
		{ID: "3", Title: "REST design", Tags: []string{"api"}, CreatedAt: now, UpdatedAt: now},
	}
}

func TestFilterFreeTextMatchesTitleOrTag(t *testing.T) {
	got := Filter(testNotes(), Query{Text: "go"})
	if len(got) != 2 {
		t.Fatalf("free text matches = %d, want 2", len(got))
	}
}

func TestFilterCombinesCriteria(t *testing.T) {
	got := Filter(testNotes(), Query{Title: "log", Tag: "daily"})
	if len(got) != 1 || got[0].ID != "2" {
		t.Fatalf("combined filter = %#v, want note 2", got)
	}
}

func TestFilterIsCaseInsensitive(t *testing.T) {
	got := Filter(testNotes(), Query{Tag: "DAILY"})
	if len(got) != 1 || got[0].ID != "2" {
		t.Fatalf("case-insensitive tag filter = %#v, want note 2", got)
	}
}
