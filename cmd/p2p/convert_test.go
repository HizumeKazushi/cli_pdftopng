package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestWorkerCountCapsAtPageCount(t *testing.T) {
	got := workerCount(8, 3)
	if got != 3 {
		t.Fatalf("workerCount = %d, want 3", got)
	}
}

func TestWorkerCountDefaultsInvalidJobsToOne(t *testing.T) {
	got := workerCount(0, 3)
	if got != 1 {
		t.Fatalf("workerCount = %d, want 1", got)
	}
}

func TestWorkerCountHandlesNoPages(t *testing.T) {
	got := workerCount(4, 0)
	if got != 0 {
		t.Fatalf("workerCount = %d, want 0", got)
	}
}

func TestConvertPagesCancelsPendingPagesAfterError(t *testing.T) {
	var converted []int
	convert := func(ctx context.Context, cfg config, outPrefix string, page int, stdout, stderr io.Writer) error {
		converted = append(converted, page)
		if page == 2 {
			return errors.New("page 2 failed")
		}
		return nil
	}

	err := convertPagesWithConverter(context.Background(), config{jobs: 1}, "out/doc", []int{1, 2, 3, 4}, &bytes.Buffer{}, io.Discard, convert)
	if err == nil {
		t.Fatal("convertPagesWithConverter returned nil error")
	}

	want := []int{1, 2}
	if !reflect.DeepEqual(converted, want) {
		t.Fatalf("converted pages = %#v, want %#v", converted, want)
	}
}

func TestConvertPagesJoinsConcurrentErrors(t *testing.T) {
	started := make(chan int, 2)
	release := make(chan struct{})
	convert := func(ctx context.Context, cfg config, outPrefix string, page int, stdout, stderr io.Writer) error {
		started <- page
		<-release
		return fmt.Errorf("page %d failed", page)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- convertPagesWithConverter(context.Background(), config{jobs: 2}, "out/doc", []int{1, 2, 3}, &bytes.Buffer{}, io.Discard, convert)
	}()

	<-started
	<-started
	close(release)

	err := <-errCh
	if err == nil {
		t.Fatal("convertPagesWithConverter returned nil error")
	}
	message := err.Error()
	if !strings.Contains(message, "page 1 failed") || !strings.Contains(message, "page 2 failed") {
		t.Fatalf("error = %q, want both page errors", message)
	}
}
