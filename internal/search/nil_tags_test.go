package search

import (
	"testing"
	"time"

	"learning-notes-cli/internal/model"
)

func TestFilterMatchesTitleWhenTagsNil(t *testing.T) {
	now := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)
	notes := []model.Note{
		{ID: "1", Title: "Go notes", Tags: nil, CreatedAt: now, UpdatedAt: now},
	}
	got := Filter(notes, Query{Text: "go"})
	if len(got) != 1 || got[0].ID != "1" {
		t.Fatalf("filter = %#v, want note 1", got)
	}
}
