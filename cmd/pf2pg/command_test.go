package main

import (
	"bytes"
	"testing"
)

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
	if !bytes.Contains(stdout.Bytes(), []byte("pf2pg [flags] input.pdf [input2.pdf...]")) {
		t.Fatalf("help output = %q, want pf2pg usage", stdout.String())
	}
}
