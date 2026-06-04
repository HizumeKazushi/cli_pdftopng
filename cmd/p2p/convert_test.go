package main

import "testing"

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
