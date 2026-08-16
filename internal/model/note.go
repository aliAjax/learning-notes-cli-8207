package model

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Note is the application representation of a Markdown learning note.
type Note struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Content   string    `json:"-"`
}

// NewNote creates a note with a unique, filesystem-safe identifier.
func NewNote(title string, tags []string, content string, now time.Time) (Note, error) {
	id, err := NewID(now)
	if err != nil {
		return Note{}, err
	}

	title = strings.TrimSpace(title)
	if title == "" {
		return Note{}, fmt.Errorf("title cannot be empty")
	}

	tags = NormalizeTags(tags)
	now = now.UTC()
	return Note{
		ID:        id,
		Title:     title,
		Tags:      tags,
		CreatedAt: now,
		UpdatedAt: now,
		Content:   normalizeContent(content),
	}, nil
}

// NewID builds a timestamped random identifier.
func NewID(now time.Time) (string, error) {
	random := make([]byte, 4)
	if _, err := rand.Read(random); err != nil {
		return "", fmt.Errorf("generate note id: %w", err)
	}

	return now.UTC().Format("20060102-150405") + "-" + hex.EncodeToString(random), nil
}

// NormalizeTags trims, de-duplicates, and sorts tags while preserving case.
// It always returns a non-nil slice, so notes without tags serialize as
// `"tags":[]` rather than `"tags":null`.
func NormalizeTags(tags []string) []string {
	seen := make(map[string]struct{}, len(tags))
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		key := strings.ToLower(tag)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, tag)
	}
	sort.Slice(result, func(i, j int) bool {
		return strings.ToLower(result[i]) < strings.ToLower(result[j])
	})
	return result
}

// ParseTags splits a comma or semicolon separated tag list.
func ParseTags(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []string{}
	}
	return NormalizeTags(strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ';'
	}))
}

// TagString returns a compact display string for note tags.
func (n Note) TagString() string {
	if len(n.Tags) == 0 {
		return ""
	}
	return strings.Join(n.Tags, ", ")
}

func normalizeContent(content string) string {
	content = strings.TrimPrefix(content, "\ufeff")
	content = strings.TrimRight(content, " \t")
	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return content
}
