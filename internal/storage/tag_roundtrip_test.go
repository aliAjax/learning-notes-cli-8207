package storage

import (
	"context"
	"testing"
	"time"

	"learning-notes-cli/internal/model"
)

func TestMarkdownStoreNormalizesTagsOnRead(t *testing.T) {
	now := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)
	note := model.Note{
		ID:        "note-1",
		Title:     "Tag round trip",
		Tags:      []string{"Go", "go", "CLI"},
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
	want := []string{"CLI", "Go"}
	if len(got.Tags) != len(want) {
		t.Fatalf("tags = %#v, want %#v", got.Tags, want)
	}
	for i := range want {
		if got.Tags[i] != want[i] {
			t.Fatalf("tags = %#v, want %#v", got.Tags, want)
		}
	}
}
