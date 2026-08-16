package storage

import (
	"context"
	"testing"
	"time"

	"learning-notes-cli/internal/model"
)

func TestMarkdownStorePreservesLeadingBlankLines(t *testing.T) {
	now := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)
	note, err := model.NewNote("Leading blanks", nil, "\n\nFirst line\n", now)
	if err != nil {
		t.Fatalf("NewNote returned error: %v", err)
	}
	store := NewMarkdownStore(t.TempDir())
	if err := store.Save(context.Background(), note); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	got, err := store.Get(context.Background(), note.ID)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	want := "\n\nFirst line\n"
	if got.Content != want {
		t.Fatalf("content = %q, want %q", got.Content, want)
	}
}
