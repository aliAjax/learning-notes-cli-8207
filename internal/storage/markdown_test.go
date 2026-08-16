package storage

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"learning-notes-cli/internal/model"
)

func TestMarkdownStoreRoundTrip(t *testing.T) {
	store := NewMarkdownStore(t.TempDir())
	now := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)
	note, err := model.NewNote("Go notes", []string{"go", "cli"}, "first line\nsecond line\n", now)
	if err != nil {
		t.Fatalf("NewNote returned error: %v", err)
	}

	if err := store.Save(context.Background(), note); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	got, err := store.Get(context.Background(), note.ID)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got.Title != note.Title || got.ID != note.ID {
		t.Fatalf("Get = %#v, want matching note", got)
	}
	if !strings.HasPrefix(got.Content, "first line") {
		t.Fatalf("content has unexpected leading blank line: %q", got.Content)
	}

	notes, err := store.List(context.Background())
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("List length = %d, want 1", len(notes))
	}

	if err := store.Delete(context.Background(), note.ID); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if _, err := store.Get(context.Background(), note.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get after Delete error = %v, want ErrNotFound", err)
	}
}

func TestMarkdownStoreListMissingDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "does-not-exist")
	notes, err := NewMarkdownStore(dir).List(context.Background())
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(notes) != 0 {
		t.Fatalf("List length = %d, want 0", len(notes))
	}
	if _, err := os.Stat(dir); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("List should not create missing directory")
	}
}
