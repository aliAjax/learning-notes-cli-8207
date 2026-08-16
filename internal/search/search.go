package search

import (
	"strings"

	"learning-notes-cli/internal/model"
)

// Query represents title, tag, and free-text search criteria.
type Query struct {
	Text  string
	Title string
	Tag   string
}

// IsEmpty reports whether no search criteria were supplied.
func (q Query) IsEmpty() bool {
	return strings.TrimSpace(q.Text) == "" &&
		strings.TrimSpace(q.Title) == "" &&
		strings.TrimSpace(q.Tag) == ""
}

// Filter returns notes matching all non-empty criteria.
func Filter(notes []model.Note, query Query) []model.Note {
	matches := make([]model.Note, 0)
	for _, note := range notes {
		if matchesText(note, query.Text) &&
			matchesTitle(note, query.Title) &&
			matchesTag(note, query.Tag) {
			matches = append(matches, note)
		}
	}
	return matches
}

func matchesText(note model.Note, text string) bool {
	text = strings.TrimSpace(text)
	if text == "" {
		return true
	}
	if containsFold(note.Title, text) {
		return true
	}
	for _, tag := range note.Tags {
		if containsFold(tag, text) {
			return true
		}
	}
	return false
}

func matchesTitle(note model.Note, title string) bool {
	title = strings.TrimSpace(title)
	return title == "" || containsFold(note.Title, title)
}

func matchesTag(note model.Note, tag string) bool {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return true
	}
	for _, existing := range note.Tags {
		if containsFold(existing, tag) {
			return true
		}
	}
	return false
}

func containsFold(value, fragment string) bool {
	return strings.Contains(strings.ToLower(value), strings.ToLower(fragment))
}
