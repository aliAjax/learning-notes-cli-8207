package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"learning-notes-cli/internal/model"
)

const frontMatterDelimiter = "---"

var (
	// ErrNotFound is returned when a requested note does not exist.
	ErrNotFound = errors.New("note not found")
	// ErrInvalidNoteID is returned when an ID cannot be used as a filename.
	ErrInvalidNoteID = errors.New("invalid note id")
)

// Store is the persistence boundary used by commands.
type Store interface {
	Save(ctx context.Context, note model.Note) error
	List(ctx context.Context) ([]model.Note, error)
	Get(ctx context.Context, id string) (model.Note, error)
	Delete(ctx context.Context, id string) error
}

// MarkdownStore stores each note as a Markdown file with JSON front matter.
type MarkdownStore struct {
	dir string
}

// NewMarkdownStore creates a MarkdownStore rooted at dir.
func NewMarkdownStore(dir string) *MarkdownStore {
	return &MarkdownStore{dir: dir}
}

// Save writes a note to <data-dir>/<id>.md.
func (s *MarkdownStore) Save(ctx context.Context, note model.Note) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := validateID(note.ID); err != nil {
		return err
	}
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return fmt.Errorf("create data directory %q: %w", s.dir, err)
	}

	data, err := marshalNote(note)
	if err != nil {
		return err
	}
	path := s.path(note.ID)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write note %q: %v", path, err)
	}
	return nil
}

// List returns all notes ordered by update time descending.
func (s *MarkdownStore) List(ctx context.Context) ([]model.Note, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(s.dir)
	if errors.Is(err, os.ErrNotExist) {
		return []model.Note{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read data directory %q: %v", s.dir, err)
	}

	notes := make([]model.Note, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		note, err := s.readFile(entry.Name())
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	sort.Slice(notes, func(i, j int) bool {
		if notes[i].UpdatedAt.Equal(notes[j].UpdatedAt) {
			return notes[i].ID > notes[j].ID
		}
		return notes[i].UpdatedAt.After(notes[j].UpdatedAt)
	})
	return notes, nil
}

// Get returns a note by exact ID.
func (s *MarkdownStore) Get(ctx context.Context, id string) (model.Note, error) {
	if err := ctx.Err(); err != nil {
		return model.Note{}, err
	}
	if err := validateID(id); err != nil {
		return model.Note{}, err
	}

	note, err := s.readFile(id + ".md")
	if errors.Is(err, os.ErrNotExist) {
		return model.Note{}, fmt.Errorf("note %q: %w", id, ErrNotFound)
	}
	if err != nil {
		return model.Note{}, err
	}
	return note, nil
}

// Delete removes a note by exact ID.
func (s *MarkdownStore) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := validateID(id); err != nil {
		return err
	}

	path := s.path(id)
	if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("note %q: %w", id, ErrNotFound)
	} else if err != nil {
		return fmt.Errorf("delete note %q: %w", path, err)
	}
	return nil
}

func (s *MarkdownStore) readFile(filename string) (model.Note, error) {
	path := filepath.Join(s.dir, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return model.Note{}, err
	}
	note, err := unmarshalNote(data)
	if err != nil {
		return model.Note{}, fmt.Errorf("parse note %q: %w", path, err)
	}
	if note.ID+".md" != filename {
		return model.Note{}, fmt.Errorf("parse note %q: filename does not match metadata id", path)
	}
	return note, nil
}

func (s *MarkdownStore) path(id string) string {
	return filepath.Join(s.dir, id+".md")
}

func validateID(id string) error {
	if id == "" || strings.ContainsAny(id, `/\`) || id == "." || id == ".." {
		return ErrInvalidNoteID
	}
	return nil
}

type noteMetadata struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func marshalNote(note model.Note) ([]byte, error) {
	meta := noteMetadata{
		ID:        note.ID,
		Title:     note.Title,
		Tags:      note.Tags,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
	metaJSON, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode note metadata: %w", err)
	}

	var builder strings.Builder
	builder.WriteString(frontMatterDelimiter)
	builder.WriteByte('\n')
	builder.Write(metaJSON)
	builder.WriteByte('\n')
	builder.WriteString(frontMatterDelimiter)
	builder.WriteByte('\n')
	builder.WriteByte('\n')
	builder.WriteString(note.Content)
	return []byte(builder.String()), nil
}

func unmarshalNote(data []byte) (model.Note, error) {
	text := string(data)
	if !strings.HasPrefix(text, frontMatterDelimiter+"\n") {
		return model.Note{}, fmt.Errorf("missing front matter delimiter")
	}

	body := text[len(frontMatterDelimiter)+1:]
	end := strings.Index(body, "\n"+frontMatterDelimiter)
	if end < 0 {
		return model.Note{}, fmt.Errorf("missing closing front matter delimiter")
	}

	metaJSON := body[:end]
	contentStart := end + len("\n"+frontMatterDelimiter)
	content := strings.TrimPrefix(body[contentStart:], "\n")
	content = strings.TrimPrefix(content, "\n")

	var meta noteMetadata
	if err := json.Unmarshal([]byte(metaJSON), &meta); err != nil {
		return model.Note{}, fmt.Errorf("decode front matter: %w", err)
	}
	if meta.ID == "" || meta.Title == "" {
		return model.Note{}, fmt.Errorf("front matter requires id and title")
	}

	return model.Note{
		ID:        meta.ID,
		Title:     meta.Title,
		Tags:      model.NormalizeTags(meta.Tags),
		CreatedAt: meta.CreatedAt,
		UpdatedAt: meta.UpdatedAt,
		Content:   content,
	}, nil
}
