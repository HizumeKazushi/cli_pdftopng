package main

import "testing"

func TestFormatProgress(t *testing.T) {
	got := formatProgress(2, 4)
	want := "[###############---------------]  50% (2/4)"

	if got != want {
		t.Fatalf("progress = %q, want %q", got, want)
	}
}

func TestFormatProgressClampsDone(t *testing.T) {
	got := formatProgress(5, 4)
	want := "[##############################] 100% (4/4)"

	if got != want {
		t.Fatalf("progress = %q, want %q", got, want)
	}
}
