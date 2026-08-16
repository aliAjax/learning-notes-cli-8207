package model

import "testing"

func TestNormalizeTagsReturnsEmptyNonNil(t *testing.T) {
	tags := NormalizeTags(nil)
	if tags == nil {
		t.Fatal("NormalizeTags(nil) returned nil, want empty slice")
	}
	if len(tags) != 0 {
		t.Fatalf("NormalizeTags(nil) length = %d, want 0", len(tags))
	}
}
