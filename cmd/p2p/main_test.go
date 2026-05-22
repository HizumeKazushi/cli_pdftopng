package main

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestValidateConfigDefaults(t *testing.T) {
	pdf := writeTempPDF(t)
	cfg := config{input: pdf, outDir: ".", dpi: 200}

	if err := validateConfig(&cfg); err != nil {
		t.Fatalf("validateConfig returned error: %v", err)
	}

	if cfg.input != pdf {
		t.Fatalf("input = %q, want %q", cfg.input, pdf)
	}
	if cfg.outDir != "." {
		t.Fatalf("outDir = %q, want .", cfg.outDir)
	}
	if cfg.prefix != "input" {
		t.Fatalf("prefix = %q, want input", cfg.prefix)
	}
	if cfg.dpi != 200 {
		t.Fatalf("dpi = %d, want 200", cfg.dpi)
	}
}

func TestValidateConfigRejectsInvalidRange(t *testing.T) {
	pdf := writeTempPDF(t)
	cfg := config{input: pdf, outDir: ".", dpi: 200, first: 3, last: 2}

	if err := validateConfig(&cfg); err == nil {
		t.Fatal("validateConfig returned nil error")
	}
}

func TestRootCommandHelpIncludesConvertSubcommand(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := newRootCommand(&stdout, &stderr)
	cmd.SetArgs([]string{"--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if !bytes.Contains(stdout.Bytes(), []byte("convert")) {
		t.Fatalf("help output = %q, want convert subcommand", stdout.String())
	}
	if !bytes.Contains(stdout.Bytes(), []byte("p2p [flags] input.pdf")) {
		t.Fatalf("help output = %q, want p2p usage", stdout.String())
	}
}

func TestBuildPdftoppmArgs(t *testing.T) {
	cfg := config{
		input: "doc.pdf",
		dpi:   300,
		first: 2,
		last:  4,
	}

	got := buildPdftoppmArgs(cfg, "out/doc")
	want := []string{"-png", "-r", "300", "-f", "2", "-l", "4", "doc.pdf", "out/doc"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("args = %#v, want %#v", got, want)
	}
}

func writeTempPDF(t *testing.T) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "input.pdf")
	if err := os.WriteFile(path, []byte("%PDF-1.4\n"), 0o644); err != nil {
		t.Fatalf("write temp pdf: %v", err)
	}
	return path
}
