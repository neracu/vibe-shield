package analyzer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func TestExtractCodeContext(t *testing.T) {
	t.Run("middle line", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "sample.go")
		writeLinesFile(t, path, 40)

		got, err := ExtractCodeContext(&DetectedError{
			FilePath:   path,
			LineNumber: 20,
		})
		if err != nil {
			t.Fatalf("ExtractCodeContext() error = %v", err)
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

		got, err := ExtractCodeContext(&DetectedError{
			FilePath:   path,
			LineNumber: 1,
		})
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

		got, err := ExtractCodeContext(&DetectedError{
			FilePath:   path,
			LineNumber: 10,
		})
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

		got, err := ExtractCodeContext(&DetectedError{
			FilePath:   relName,
			LineNumber: 2,
		})
		if err != nil {
			t.Fatalf("ExtractCodeContext() error = %v", err)
		}
		if got != "1 | x\n2 >> xx\n3 | xxx" {
			t.Errorf("got %q, want %q", got, "1 | x\n2 >> xx\n3 | xxx")
		}
	})

	t.Run("missing file", func(t *testing.T) {
		got, err := ExtractCodeContext(&DetectedError{
			FilePath:   filepath.Join(t.TempDir(), "missing.go"),
			LineNumber: 1,
		})
		if err == nil {
			t.Fatal("expected error for missing file")
		}
		if got != "" {
			t.Errorf("got %q, want empty string", got)
		}
	})

	t.Run("invalid line zero", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "short.go")
		writeLinesFile(t, path, 3)

		got, err := ExtractCodeContext(&DetectedError{
			FilePath:   path,
			LineNumber: 0,
		})
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

		got, err := ExtractCodeContext(&DetectedError{
			FilePath:   path,
			LineNumber: 99,
		})
		if err == nil {
			t.Fatal("expected error for out-of-range line")
		}
		if got != "" {
			t.Errorf("got %q, want empty string", got)
		}
	})
}
