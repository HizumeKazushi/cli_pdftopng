package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

func ensurePopplerTools() error {
	if _, err := exec.LookPath("pdftoppm"); err != nil {
		return errors.New("pdftoppmが見つかりません。macOSなら `brew install poppler` でインストールしてください")
	}
	if _, err := exec.LookPath("pdfinfo"); err != nil {
		return errors.New("pdfinfoが見つかりません。macOSなら `brew install poppler` でインストールしてください")
	}
	return nil
}

func getPageCount(input string) (int, error) {
	output, err := exec.Command("pdfinfo", input).Output()
	if err != nil {
		return 0, fmt.Errorf("PDFのページ数を取得できません: %w", err)
	}

	pages, err := parsePageCount(string(output))
	if err != nil {
		return 0, err
	}
	return pages, nil
}

func parsePageCount(output string) (int, error) {
	for _, line := range strings.Split(output, "\n") {
		key, value, ok := strings.Cut(line, ":")
		if !ok || strings.TrimSpace(key) != "Pages" {
			continue
		}

		pages, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil || pages <= 0 {
			return 0, errors.New("PDFのページ数が不正です")
		}
		return pages, nil
	}
	return 0, errors.New("PDFのページ数を取得できません")
}

func convertPage(ctx context.Context, cfg config, outPrefix string, page int, stdout, stderr io.Writer) error {
	cmdArgs := buildPdftoppmArgs(cfg, outPrefix, page)
	cmd := exec.CommandContext(ctx, "pdftoppm", cmdArgs...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%dページ目の変換に失敗しました: %w", page, err)
	}
	return nil
}

func buildPdftoppmArgs(cfg config, outPrefix string, page int) []string {
	return []string{
		"-png",
		"-r", fmt.Sprint(cfg.dpi),
		"-f", fmt.Sprint(page),
		"-l", fmt.Sprint(page),
		cfg.input,
		outPrefix,
	}
}
