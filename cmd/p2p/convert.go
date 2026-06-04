package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type pageConverter func(context.Context, config, string, int, io.Writer, io.Writer) error
type pdfConverter func(config, io.Writer, io.Writer) error
type pageResult struct {
	page int
	err  error
}

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
	return convertPDFsWithConverter(cfg, inputs, stdout, stderr, convertPDF)
}

func convertPDFsWithConverter(cfg config, inputs []string, stdout, stderr io.Writer, convert pdfConverter) error {
	if len(inputs) == 1 {
		cfg.input = inputs[0]
		return convert(cfg, stdout, stderr)
	}

	var errs []error
	for i, input := range inputs {
		pdfCfg := cfg
		pdfCfg.input = input
		pdfCfg.batch = true

		fmt.Fprintf(stdout, "Converting %s (%d/%d)\n", input, i+1, len(inputs))
		if err := convert(pdfCfg, stdout, stderr); err != nil {
			err = fmt.Errorf("%s: %w", input, err)
			if !cfg.continueOnError {
				return err
			}
			fmt.Fprintf(stderr, "%v\n", err)
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func convertPages(cfg config, outPrefix string, pages []int, stdout, stderr io.Writer) error {
	return convertPagesWithConverter(context.Background(), cfg, outPrefix, pages, stdout, stderr, convertPage)
}

func convertPagesWithConverter(parent context.Context, cfg config, outPrefix string, pages []int, stdout, stderr io.Writer, convert pageConverter) error {
	jobs := workerCount(cfg.jobs, len(pages))
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	pageJobs := make(chan int)
	results := make(chan pageResult, len(pages))

	var wg sync.WaitGroup
	for range jobs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case page, ok := <-pageJobs:
					if !ok {
						return
					}
					err := convert(ctx, cfg, outPrefix, page, io.Discard, stderr)
					if err != nil {
						cancel()
					}
					results <- pageResult{page: page, err: err}
				}
			}
		}()
	}

	go func() {
		defer close(pageJobs)
		for _, page := range pages {
			select {
			case <-ctx.Done():
				return
			case pageJobs <- page:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	done := 0
	var errs []error
	for result := range results {
		done++
		if result.err != nil {
			errs = append(errs, result.err)
		}
		printProgress(stdout, done, len(pages))
	}
	return errors.Join(errs...)
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
