package command

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"learning-notes-cli/internal/model"
	"learning-notes-cli/internal/search"
	"learning-notes-cli/internal/storage"
)

type staticStore struct {
	notes []model.Note
}

func (s staticStore) Save(context.Context, model.Note) error { return nil }
func (s staticStore) List(context.Context) ([]model.Note, error) {
	return s.notes, nil
}
func (s staticStore) Get(context.Context, string) (model.Note, error) {
	return model.Note{}, storage.ErrNotFound
}
func (s staticStore) Delete(context.Context, string) error { return nil }

func TestRunSearchAppliesFilter(t *testing.T) {
	now := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)
	store := staticStore{notes: []model.Note{
		{ID: "1", Title: "Go notes", Tags: []string{"go"}, CreatedAt: now, UpdatedAt: now},
		{ID: "2", Title: "REST design", Tags: []string{"api"}, CreatedAt: now, UpdatedAt: now},
	}}
	var out, errOut bytes.Buffer
	app := NewApp(strings.NewReader(""), &out, &errOut)
	code := app.runSearch(context.Background(), []string{"go"}, store, search.Filter)
	if code != 0 {
		t.Fatalf("runSearch code = %d, want 0", code)
	}
	if !strings.Contains(out.String(), "Go notes") || strings.Contains(out.String(), "REST design") {
		t.Fatalf("search output = %q, want only Go notes", out.String())
	}
}
