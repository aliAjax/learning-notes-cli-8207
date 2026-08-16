package storage

import (
	"context"
	"errors"
	"testing"
)

func TestStoreErrorsPreserveSentinel(t *testing.T) {
	store := NewMarkdownStore(t.TempDir())
	if err := store.Delete(context.Background(), "missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Delete missing error = %v, want ErrNotFound", err)
	}
	if _, err := store.Get(context.Background(), "missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get missing error = %v, want ErrNotFound", err)
	}
}
