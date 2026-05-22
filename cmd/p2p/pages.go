package main

import "fmt"

func selectedPages(cfg config, totalPages int) ([]int, error) {
	first := cfg.first
	if first == 0 {
		first = 1
	}
	last := cfg.last
	if last == 0 {
		last = totalPages
	}
	if first > totalPages {
		return nil, fmt.Errorf("--firstがPDFの総ページ数を超えています: %d", totalPages)
	}
	if last > totalPages {
		last = totalPages
	}

	pages := make([]int, 0, last-first+1)
	for page := first; page <= last; page++ {
		pages = append(pages, page)
	}
	return pages, nil
}
