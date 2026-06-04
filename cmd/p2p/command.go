package main

import (
	"io"
	"runtime"

	"github.com/spf13/cobra"
)

func run(args []string, stdout, stderr io.Writer) error {
	cmd := newRootCommand(stdout, stderr)
	cmd.SetArgs(args)
	return cmd.Execute()
}

func newRootCommand(stdout, stderr io.Writer) *cobra.Command {
	cfg := config{}
	rootCmd := &cobra.Command{
		Use:           appName + " [flags] input.pdf [input2.pdf...]",
		Short:         "Convert PDF pages to PNG files",
		Args:          cobra.MinimumNArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return convertPDFs(cfg, args, stdout, stderr)
		},
	}
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	addConvertFlags(rootCmd, &cfg)

	convertCfg := config{}
	convertCmd := &cobra.Command{
		Use:           "convert [flags] input.pdf [input2.pdf...]",
		Short:         "Convert PDF pages to PNG files",
		Args:          cobra.MinimumNArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return convertPDFs(convertCfg, args, stdout, stderr)
		},
	}
	addConvertFlags(convertCmd, &convertCfg)
	rootCmd.AddCommand(convertCmd)

	return rootCmd
}

func addConvertFlags(cmd *cobra.Command, cfg *config) {
	cmd.Flags().StringVarP(&cfg.outDir, "out", "o", "", "output directory (default: input file name without .pdf)")
	cmd.Flags().StringVar(&cfg.prefix, "prefix", "", "output file prefix")
	cmd.Flags().IntVar(&cfg.dpi, "dpi", 200, "rendering DPI")
	cmd.Flags().IntVar(&cfg.first, "first", 0, "first page to convert, 1-based")
	cmd.Flags().IntVar(&cfg.last, "last", 0, "last page to convert, 1-based")
	cmd.Flags().IntVarP(&cfg.jobs, "jobs", "j", runtime.NumCPU(), "number of pages to convert in parallel")
}
