package model

import "testing"

func TestTagPipelineNormalizesAndSorts(t *testing.T) {
	tags := ParseTags(" Go;cli, go;CLI ")
	want := []string{"cli", "Go"}
	if len(tags) != len(want) {
		t.Fatalf("tags = %#v, want %#v", tags, want)
	}
	for i := range want {
		if tags[i] != want[i] {
			t.Fatalf("tags = %#v, want %#v", tags, want)
		}
	}
}
