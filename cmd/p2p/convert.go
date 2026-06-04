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
	printProgress(stdout, 0, len(pages))
	for i, page := range pages {
		if err := convertPage(cfg, outPrefix, page, stdout, stderr); err != nil {
			return err
		}

		printProgress(stdout, i+1, len(pages))
	}

	fmt.Fprintf(stdout, "\nOutput: %s\nDone\n", cfg.outDir)
	return nil
}

func convertPDFs(cfg config, inputs []string, stdout, stderr io.Writer) error {
	if len(inputs) == 1 {
		cfg.input = inputs[0]
		return convertPDF(cfg, stdout, stderr)
	}

	for i, input := range inputs {
		pdfCfg := cfg
		pdfCfg.input = input
		pdfCfg.batch = true

		fmt.Fprintf(stdout, "Converting %s (%d/%d)\n", input, i+1, len(inputs))
		if err := convertPDF(pdfCfg, stdout, stderr); err != nil {
			return err
		}
	}
	return nil
}
