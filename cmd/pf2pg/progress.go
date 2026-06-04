package main

import (
	"fmt"
	"io"
	"strings"
)

const progressBarWidth = 30

func printProgress(w io.Writer, done, total int) {
	fmt.Fprint(w, "\r"+formatProgress(done, total))
}

func formatProgress(done, total int) string {
	if total <= 0 {
		return "[------------------------------]   0% (0/0)"
	}
	if done < 0 {
		done = 0
	}
	if done > total {
		done = total
	}

	percent := int(float64(done) / float64(total) * 100)
	filled := progressBarWidth * done / total
	empty := progressBarWidth - filled

	return fmt.Sprintf(
		"[%s%s] %3d%% (%d/%d)",
		strings.Repeat("#", filled),
		strings.Repeat("-", empty),
		percent,
		done,
		total,
	)
}
