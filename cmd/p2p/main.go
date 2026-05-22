package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const appName = "p2p"

type config struct {
	input  string
	outDir string
	prefix string
	dpi    int
	first  int
	last   int
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	cmd := newRootCommand(stdout, stderr)
	cmd.SetArgs(args)
	return cmd.Execute()
}

func newRootCommand(stdout, stderr io.Writer) *cobra.Command {
	cfg := config{}
	rootCmd := &cobra.Command{
		Use:           appName + " [flags] input.pdf",
		Short:         "Convert PDF pages to PNG files",
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg.input = args[0]
			return convertPDF(cfg, stdout, stderr)
		},
	}
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	addConvertFlags(rootCmd, &cfg)

	convertCfg := config{}
	convertCmd := &cobra.Command{
		Use:           "convert [flags] input.pdf",
		Short:         "Convert PDF pages to PNG files",
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			convertCfg.input = args[0]
			return convertPDF(convertCfg, stdout, stderr)
		},
	}
	addConvertFlags(convertCmd, &convertCfg)
	rootCmd.AddCommand(convertCmd)

	return rootCmd
}

func addConvertFlags(cmd *cobra.Command, cfg *config) {
	cmd.Flags().StringVarP(&cfg.outDir, "out", "o", ".", "output directory")
	cmd.Flags().StringVar(&cfg.prefix, "prefix", "", "output file prefix")
	cmd.Flags().IntVar(&cfg.dpi, "dpi", 200, "rendering DPI")
	cmd.Flags().IntVar(&cfg.first, "first", 0, "first page to convert, 1-based")
	cmd.Flags().IntVar(&cfg.last, "last", 0, "last page to convert, 1-based")
}

func convertPDF(cfg config, stdout, stderr io.Writer) error {
	if err := validateConfig(&cfg); err != nil {
		return err
	}
	if _, err := exec.LookPath("pdftoppm"); err != nil {
		return errors.New("pdftoppmが見つかりません。macOSなら `brew install poppler` でインストールしてください")
	}

	if err := os.MkdirAll(cfg.outDir, 0o755); err != nil {
		return fmt.Errorf("出力ディレクトリを作成できません: %w", err)
	}

	outPrefix := filepath.Join(cfg.outDir, cfg.prefix)
	cmdArgs := buildPdftoppmArgs(cfg, outPrefix)
	cmd := exec.Command("pdftoppm", cmdArgs...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("変換に失敗しました: %w", err)
	}

	fmt.Fprintf(stdout, "PNGを書き出しました: %s\n", cfg.outDir)
	return nil
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
	if cfg.prefix == "" {
		base := filepath.Base(cfg.input)
		cfg.prefix = strings.TrimSuffix(base, filepath.Ext(base))
	}

	return nil
}

func buildPdftoppmArgs(cfg config, outPrefix string) []string {
	args := []string{"-png", "-r", fmt.Sprint(cfg.dpi)}
	if cfg.first > 0 {
		args = append(args, "-f", fmt.Sprint(cfg.first))
	}
	if cfg.last > 0 {
		args = append(args, "-l", fmt.Sprint(cfg.last))
	}
	args = append(args, cfg.input, outPrefix)
	return args
}
