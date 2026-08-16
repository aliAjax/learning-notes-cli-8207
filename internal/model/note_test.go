package model

import (
	"strings"
	"testing"
	"time"
)

func TestNewNoteNormalizesFields(t *testing.T) {
	now := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)
	note, err := NewNote("  Go notes  ", []string{" go ", "cli", "GO"}, "content\n", now)
	if err != nil {
		t.Fatalf("NewNote returned error: %v", err)
	}

	if note.Title != "Go notes" {
		t.Fatalf("title = %q, want %q", note.Title, "Go notes")
	}
	if len(note.Tags) != 2 || note.Tags[0] != "cli" || note.Tags[1] != "go" {
		t.Fatalf("tags = %#v, want sorted de-duplicated tags", note.Tags)
	}
	if !strings.HasSuffix(note.Content, "\n") {
		t.Fatalf("content should end with newline")
	}
	if note.CreatedAt != now || note.UpdatedAt != now {
		t.Fatalf("created/updated times = %v/%v, want %v", note.CreatedAt, note.UpdatedAt, now)
	}
}

func TestParseTags(t *testing.T) {
	tags := ParseTags(" go;cli, go, ")
	if len(tags) != 2 || tags[0] != "cli" || tags[1] != "go" {
		t.Fatalf("tags = %#v, want [cli go]", tags)
	}
}
