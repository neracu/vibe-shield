package analyzer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const contextRadius = 15

type CodeContextOpts struct {
	RunStartedAt    time.Time
	BaselineModTime time.Time // zero = baseline stat was unavailable
}

func FileModTime(path string) (time.Time, error) {
	absPath, err := resolveAbsPath(path)
	if err != nil {
		return time.Time{}, err
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

func ExtractCodeContext(de *DetectedError, opts CodeContextOpts) (snippet string, stale bool, err error) {
	path, err := resolveAbsPath(de.FilePath)
	if err != nil {
		return "", false, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return "", false, fmt.Errorf("stat file %q: %w", path, err)
	}
	stale = isStaleModTime(info.ModTime(), opts)

	f, err := os.Open(path)
	if err != nil {
		return "", stale, fmt.Errorf("open file %q: %w", path, err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", stale, fmt.Errorf("read file %q: %w", path, err)
	}

	if de.LineNumber < 1 || de.LineNumber > len(lines) {
		return "", stale, fmt.Errorf("line %d out of range (file has %d lines)", de.LineNumber, len(lines))
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
	return b.String(), stale, nil
}

func resolveAbsPath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("resolve working directory: %w", err)
	}
	return filepath.Join(wd, path), nil
}

func isStaleModTime(modTime time.Time, opts CodeContextOpts) bool {
	if !opts.RunStartedAt.IsZero() && !modTime.Before(opts.RunStartedAt) {
		return true
	}
	if !opts.BaselineModTime.IsZero() && !modTime.Equal(opts.BaselineModTime) {
		return true
	}
	return false
}
