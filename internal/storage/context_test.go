package storage

import (
	"context"
	"errors"
	"testing"
	"time"

	"learning-notes-cli/internal/model"
)

func TestStoreRespectsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	store := NewMarkdownStore(t.TempDir())
	now := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)
	note := model.Note{
		ID:        "note-1",
		Title:     "Context",
		CreatedAt: now,
		UpdatedAt: now,
		Content:   "body\n",
	}

	if err := store.Save(ctx, note); !errors.Is(err, context.Canceled) {
		t.Fatalf("Save canceled error = %v, want context.Canceled", err)
	}
	if _, err := store.Get(ctx, note.ID); !errors.Is(err, context.Canceled) {
		t.Fatalf("Get canceled error = %v, want context.Canceled", err)
	}
	if _, err := store.List(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("List canceled error = %v, want context.Canceled", err)
	}
	if err := store.Delete(ctx, note.ID); !errors.Is(err, context.Canceled) {
		t.Fatalf("Delete canceled error = %v, want context.Canceled", err)
	}
}
