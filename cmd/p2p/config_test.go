package main

import "testing"

func TestValidateConfigDefaults(t *testing.T) {
	pdf := writeTempPDF(t)
	cfg := config{input: pdf, dpi: 200}

	if err := validateConfig(&cfg); err != nil {
		t.Fatalf("validateConfig returned error: %v", err)
	}

	if cfg.input != pdf {
		t.Fatalf("input = %q, want %q", cfg.input, pdf)
	}
	if cfg.outDir != "input" {
		t.Fatalf("outDir = %q, want input", cfg.outDir)
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
	cfg := config{input: pdf, dpi: 200, first: 3, last: 2}

	if err := validateConfig(&cfg); err == nil {
		t.Fatal("validateConfig returned nil error")
	}
}

func TestValidateConfigUsesInputSubdirectoryForBatchOutput(t *testing.T) {
	pdf := writeTempPDF(t)
	cfg := config{input: pdf, outDir: "out", dpi: 200, batch: true}

	if err := validateConfig(&cfg); err != nil {
		t.Fatalf("validateConfig returned error: %v", err)
	}

	if cfg.outDir != "out/input" {
		t.Fatalf("outDir = %q, want out/input", cfg.outDir)
	}
	if cfg.prefix != "input" {
		t.Fatalf("prefix = %q, want input", cfg.prefix)
	}
}
