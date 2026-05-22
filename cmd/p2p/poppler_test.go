package main

import (
	"reflect"
	"testing"
)

func TestParsePageCount(t *testing.T) {
	got, err := parsePageCount("Title: sample\nPages:          12\n")
	if err != nil {
		t.Fatalf("parsePageCount returned error: %v", err)
	}
	if got != 12 {
		t.Fatalf("pages = %d, want 12", got)
	}
}

func TestBuildPdftoppmArgsForSinglePage(t *testing.T) {
	cfg := config{
		input: "doc.pdf",
		dpi:   300,
	}

	got := buildPdftoppmArgs(cfg, "out/doc", 4)
	want := []string{"-png", "-r", "300", "-f", "4", "-l", "4", "doc.pdf", "out/doc"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("args = %#v, want %#v", got, want)
	}
}
