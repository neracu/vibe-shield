package analyzer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const contextRadius = 15

func ExtractCodeContext(de *DetectedError) (string, error) {
	path := de.FilePath
	if !filepath.IsAbs(path) {
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("resolve working directory: %w", err)
		}
		path = filepath.Join(wd, path)
	}

	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open file %q: %w", path, err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("read file %q: %w", path, err)
	}

	if de.LineNumber < 1 || de.LineNumber > len(lines) {
		return "", fmt.Errorf("line %d out of range (file has %d lines)", de.LineNumber, len(lines))
	}

	start := de.LineNumber - contextRadius
	if start < 1 {
		start = 1
	}
	end := de.LineNumber + contextRadius
	if end > len(lines) {
		end = len(lines)
	}

	var b strings.Builder
	for i := start; i <= end; i++ {
		if i > start {
			b.WriteByte('\n')
		}
		if i == de.LineNumber {
			fmt.Fprintf(&b, "%d >> %s", i, lines[i-1])
		} else {
			fmt.Fprintf(&b, "%d | %s", i, lines[i-1])
		}
	}
	return b.String(), nil
}
