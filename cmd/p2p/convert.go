package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func convertPDF(cfg config, stdout, stderr io.Writer) error {
	if err := validateConfig(&cfg); err != nil {
		return err
	}
	if err := ensurePopplerTools(); err != nil {
		return err
	}
	if err := os.MkdirAll(cfg.outDir, 0o755); err != nil {
		return fmt.Errorf("出力ディレクトリを作成できません: %w", err)
	}

	totalPages, err := getPageCount(cfg.input)
	if err != nil {
		return err
	}
	pages, err := selectedPages(cfg, totalPages)
	if err != nil {
		return err
	}

	outPrefix := filepath.Join(cfg.outDir, cfg.prefix)
	fmt.Fprintf(stdout, "0%%\n")
	for i, page := range pages {
		if err := convertPage(cfg, outPrefix, page, stdout, stderr); err != nil {
			return err
		}

		percent := int(float64(i+1) / float64(len(pages)) * 100)
		fmt.Fprintf(stdout, "%d%%\n", percent)
	}

	fmt.Fprintf(stdout, "Output: %s\nDone\n", cfg.outDir)
	return nil
}
