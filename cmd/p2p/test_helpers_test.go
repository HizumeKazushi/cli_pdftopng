package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempPDF(t *testing.T) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "input.pdf")
	if err := os.WriteFile(path, []byte("%PDF-1.4\n"), 0o644); err != nil {
		t.Fatalf("write temp pdf: %v", err)
	}
	return path
}
