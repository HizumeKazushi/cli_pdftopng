package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
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
	if err := convertPages(cfg, outPrefix, pages, stdout, stderr); err != nil {
		return err
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

func convertPages(cfg config, outPrefix string, pages []int, stdout, stderr io.Writer) error {
	jobs := workerCount(cfg.jobs, len(pages))
	pageJobs := make(chan int)
	results := make(chan error, len(pages))

	var wg sync.WaitGroup
	for range jobs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for page := range pageJobs {
				results <- convertPage(cfg, outPrefix, page, io.Discard, stderr)
			}
		}()
	}

	go func() {
		for _, page := range pages {
			pageJobs <- page
		}
		close(pageJobs)
		wg.Wait()
		close(results)
	}()

	done := 0
	var firstErr error
	for err := range results {
		done++
		if err != nil && firstErr == nil {
			firstErr = err
		}
		printProgress(stdout, done, len(pages))
	}
	return firstErr
}

func workerCount(jobs, pages int) int {
	if pages <= 0 {
		return 0
	}
	if jobs <= 0 {
		return 1
	}
	if jobs > pages {
		return pages
	}
	return jobs
}
