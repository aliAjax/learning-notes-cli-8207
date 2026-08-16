package search

import (
	"testing"
	"time"

	"learning-notes-cli/internal/model"
)

func TestFilterTagIsExactCaseInsensitive(t *testing.T) {
	now := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)
	notes := []model.Note{
		{ID: "1", Title: "Go notes", Tags: []string{"go"}, CreatedAt: now, UpdatedAt: now},
		{ID: "2", Title: "Golang tips", Tags: []string{"golang"}, CreatedAt: now, UpdatedAt: now},
	}
	got := Filter(notes, Query{Tag: "GO"})
	if len(got) != 1 || got[0].ID != "1" {
		t.Fatalf("tag filter = %#v, want only note 1", got)
	}
}
