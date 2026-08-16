package storage

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"learning-notes-cli/internal/model"
)

func TestMarkdownStorePreservesEmptyTags(t *testing.T) {
	now := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)
	note := model.Note{
		ID:        "note-1",
		Title:     "No tags",
		Tags:      nil,
		CreatedAt: now,
		UpdatedAt: now,
		Content:   "body\n",
	}
	store := NewMarkdownStore(t.TempDir())
	if err := store.Save(context.Background(), note); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	got, err := store.Get(context.Background(), note.ID)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got.Tags == nil {
		t.Fatal("Tags after round trip are nil, want empty slice")
	}
	data, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("json.Marshal returned error: %v", err)
	}
	if strings.Contains(string(data), `"tags":null`) {
		t.Fatalf("JSON contains null tags: %s", data)
	}
}
