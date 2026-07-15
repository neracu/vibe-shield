package analyzer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func writeLinesFile(t *testing.T, path string, count int) {
	t.Helper()
	lines := make([]string, count)
	for i := range lines {
		lines[i] = strings.Repeat("x", i+1)
	}
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
		t.Fatalf("write test file: %v", err)
	}
}

func freshOpts() CodeContextOpts {
	return CodeContextOpts{RunStartedAt: time.Now().Add(time.Minute)}
}

func TestExtractCodeContext(t *testing.T) {
	t.Run("middle line", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "sample.go")
		writeLinesFile(t, path, 40)

		got, stale, err := ExtractCodeContext(&DetectedError{
			FilePath:   path,
			LineNumber: 20,
		}, freshOpts())
		if err != nil {
			t.Fatalf("ExtractCodeContext() error = %v", err)
		}
		if stale {
			t.Fatal("expected non-stale snippet")
		}

		lines := strings.Split(got, "\n")
		if len(lines) != 31 {
			t.Fatalf("got %d lines, want 31", len(lines))
		}
		if !strings.HasPrefix(lines[0], "5 | ") {
			t.Errorf("first line = %q, want prefix %q", lines[0], "5 | ")
		}
		if lines[15] != "20 >> xxxxxxxxxxxxxxxxxxxx" {
			t.Errorf("error line = %q, want %q", lines[15], "20 >> xxxxxxxxxxxxxxxxxxxx")
		}
		if !strings.HasPrefix(lines[30], "35 | ") {
			t.Errorf("last line = %q, want prefix %q", lines[30], "35 | ")
		}
	})

	t.Run("first lines", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "top.go")
		writeLinesFile(t, path, 10)

		got, _, err := ExtractCodeContext(&DetectedError{
			FilePath:   path,
			LineNumber: 1,
		}, freshOpts())
		if err != nil {
			t.Fatalf("ExtractCodeContext() error = %v", err)
		}

		lines := strings.Split(got, "\n")
		if len(lines) != 10 {
			t.Fatalf("got %d lines, want 10", len(lines))
		}
		if lines[0] != "1 >> x" {
			t.Errorf("first line = %q, want %q", lines[0], "1 >> x")
		}
	})

	t.Run("last lines", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "bottom.go")
		writeLinesFile(t, path, 10)

		got, _, err := ExtractCodeContext(&DetectedError{
			FilePath:   path,
			LineNumber: 10,
		}, freshOpts())
		if err != nil {
			t.Fatalf("ExtractCodeContext() error = %v", err)
		}

		lines := strings.Split(got, "\n")
		if len(lines) != 10 {
			t.Fatalf("got %d lines, want 10", len(lines))
		}
		if lines[9] != "10 >> xxxxxxxxxx" {
			t.Errorf("last line = %q, want %q", lines[9], "10 >> xxxxxxxxxx")
		}
	})

	t.Run("relative path", func(t *testing.T) {
		dir := t.TempDir()
		relName := "rel.go"
		writeLinesFile(t, filepath.Join(dir, relName), 3)

		wd, err := os.Getwd()
		if err != nil {
			t.Fatalf("Getwd() error = %v", err)
		}
		if err := os.Chdir(dir); err != nil {
			t.Fatalf("Chdir() error = %v", err)
		}
		defer func() {
			if err := os.Chdir(wd); err != nil {
				t.Errorf("restore wd: %v", err)
			}
		}()

		got, _, err := ExtractCodeContext(&DetectedError{
			FilePath:   relName,
			LineNumber: 2,
		}, freshOpts())
		if err != nil {
			t.Fatalf("ExtractCodeContext() error = %v", err)
		}
		if got != "1 | x\n2 >> xx\n3 | xxx" {
			t.Errorf("got %q, want %q", got, "1 | x\n2 >> xx\n3 | xxx")
		}
	})

	t.Run("missing file", func(t *testing.T) {
		got, stale, err := ExtractCodeContext(&DetectedError{
			FilePath:   filepath.Join(t.TempDir(), "missing.go"),
			LineNumber: 1,
		}, freshOpts())
		if err == nil {
			t.Fatal("expected error for missing file")
		}
		if got != "" {
			t.Errorf("got %q, want empty string", got)
		}
		if stale {
			t.Error("expected non-stale for missing file")
		}
	})

	t.Run("invalid line zero", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "short.go")
		writeLinesFile(t, path, 3)

		got, _, err := ExtractCodeContext(&DetectedError{
			FilePath:   path,
			LineNumber: 0,
		}, freshOpts())
		if err == nil {
			t.Fatal("expected error for line 0")
		}
		if got != "" {
			t.Errorf("got %q, want empty string", got)
		}
	})

	t.Run("invalid line beyond file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "short.go")
		writeLinesFile(t, path, 3)

		got, _, err := ExtractCodeContext(&DetectedError{
			FilePath:   path,
			LineNumber: 99,
		}, freshOpts())
		if err == nil {
			t.Fatal("expected error for out-of-range line")
		}
		if got != "" {
			t.Errorf("got %q, want empty string", got)
		}
	})

	t.Run("stale when modified after run start", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "stale.go")
		writeLinesFile(t, path, 5)

		runStarted := time.Now().Add(-time.Hour)
		newMtime := time.Now().Add(time.Minute)
		if err := os.Chtimes(path, newMtime, newMtime); err != nil {
			t.Fatalf("Chtimes() error = %v", err)
		}

		got, stale, err := ExtractCodeContext(&DetectedError{
			FilePath:   path,
			LineNumber: 3,
		}, CodeContextOpts{RunStartedAt: runStarted})
		if err != nil {
			t.Fatalf("ExtractCodeContext() error = %v", err)
		}
		if !stale {
			t.Fatal("expected stale snippet when file mtime is after run start")
		}
		if !strings.Contains(got, "3 >>") {
			t.Errorf("got %q, want error line marker", got)
		}
	})

	t.Run("stale when baseline mtime differs", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "baseline.go")
		writeLinesFile(t, path, 5)

		got, stale, err := ExtractCodeContext(&DetectedError{
			FilePath:   path,
			LineNumber: 2,
		}, CodeContextOpts{
			RunStartedAt:    time.Now().Add(-time.Hour),
			BaselineModTime: time.Now().Add(-time.Minute),
		})
		if err != nil {
			t.Fatalf("ExtractCodeContext() error = %v", err)
		}
		if !stale {
			t.Fatal("expected stale snippet when baseline mtime differs")
		}
		if !strings.Contains(got, "2 >>") {
			t.Errorf("got %q, want error line marker", got)
		}
	})

	t.Run("not stale when baseline matches and run started after mtime", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "fresh.go")
		writeLinesFile(t, path, 5)

		pastMtime := time.Now().Add(-time.Hour)
		if err := os.Chtimes(path, pastMtime, pastMtime); err != nil {
			t.Fatalf("Chtimes() error = %v", err)
		}

		modTime, err := FileModTime(path)
		if err != nil {
			t.Fatalf("FileModTime() error = %v", err)
		}

		got, stale, err := ExtractCodeContext(&DetectedError{
			FilePath:   path,
			LineNumber: 2,
		}, CodeContextOpts{
			RunStartedAt:    time.Now(),
			BaselineModTime: modTime,
		})
		if err != nil {
			t.Fatalf("ExtractCodeContext() error = %v", err)
		}
		if stale {
			t.Fatal("expected non-stale snippet when baseline matches and file predates run")
		}
		if !strings.Contains(got, "2 >>") {
			t.Errorf("got %q, want error line marker", got)
		}
	})
}
