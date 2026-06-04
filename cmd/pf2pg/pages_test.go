package main

import (
	"reflect"
	"testing"
)

func TestSelectedPagesUsesFullDocumentByDefault(t *testing.T) {
	got, err := selectedPages(config{}, 3)
	if err != nil {
		t.Fatalf("selectedPages returned error: %v", err)
	}

	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("pages = %#v, want %#v", got, want)
	}
}

func TestSelectedPagesClampsLastPage(t *testing.T) {
	got, err := selectedPages(config{first: 2, last: 9}, 3)
	if err != nil {
		t.Fatalf("selectedPages returned error: %v", err)
	}

	want := []int{2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("pages = %#v, want %#v", got, want)
	}
}
