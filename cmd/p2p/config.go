package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	input  string
	outDir string
	prefix string
	dpi    int
	first  int
	last   int
	batch  bool
}

func validateConfig(cfg *config) error {
	if filepath.Ext(strings.ToLower(cfg.input)) != ".pdf" {
		return errors.New("入力ファイルは.pdfである必要があります")
	}
	if _, err := os.Stat(cfg.input); err != nil {
		return fmt.Errorf("入力ファイルを読めません: %w", err)
	}
	if cfg.dpi <= 0 {
		return errors.New("--dpiは1以上を指定してください")
	}
	if cfg.first < 0 || cfg.last < 0 {
		return errors.New("--firstと--lastは0以上を指定してください")
	}
	if cfg.first > 0 && cfg.last > 0 && cfg.first > cfg.last {
		return errors.New("--firstは--last以下にしてください")
	}
	resolveOutputDefaults(cfg)
	return nil
}

func resolveOutputDefaults(cfg *config) {
	base := filepath.Base(cfg.input)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	if cfg.outDir == "" {
		cfg.outDir = name
	} else if cfg.batch {
		cfg.outDir = filepath.Join(cfg.outDir, name)
	}
	if cfg.prefix == "" {
		cfg.prefix = name
	}
}
