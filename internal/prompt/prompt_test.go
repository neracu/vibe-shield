package prompt

import (
	"strings"
	"testing"

	"github.com/neracu/vibe-shield/internal/analyzer"
)

func TestTailLogs(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  int
	}{
		{name: "empty", input: nil, want: 0},
		{name: "under limit", input: []string{"a", "b", "c"}, want: 3},
		{name: "at limit", input: makeLines(10), want: 10},
		{name: "over limit", input: makeLines(15), want: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TailLogs(tt.input)
			if len(got) != tt.want {
				t.Fatalf("TailLogs() len = %d, want %d", len(got), tt.want)
			}
			if tt.want > 0 && tt.input != nil && got[0] != tt.input[len(tt.input)-tt.want] {
				t.Errorf("TailLogs() first = %q, want %q", got[0], tt.input[len(tt.input)-tt.want])
			}
		})
	}
}

func TestGenerateMarkdownPrompt(t *testing.T) {
	err := &analyzer.DetectedError{
		ErrorType:    "ReferenceError",
		ErrorMessage: "x is not defined",
		FilePath:     `C:\proj\src\index.js`,
		LineNumber:   10,
		StackTrace: []string{
			"ReferenceError: x is not defined",
			`    at Object.<anonymous> (C:\proj\src\index.js:10:5)`,
		},
	}
	snippet := "10 | const x = 1;\n11 >> console.log(y);"
	lastLogs := []string{"[INFO] starting", "[INFO] ready"}

	md := GenerateMarkdownPrompt(err, snippet, lastLogs)

	checks := []string{
		"# Vibe-Shield Crash Report",
		"### THE ERROR",
		"### CODE SNIPPET",
		"### LAST LOGS",
		"### INSTRUCTION FOR AI",
		"**Type:** ReferenceError",
		"**Message:** x is not defined",
		"ReferenceError: x is not defined",
		"C:\\proj\\src\\index.js:10:5",
		"```javascript",
		snippet,
		"[INFO] starting",
		"[INFO] ready",
		"Fix only this specific error",
		"minimal diff",
		"existing architecture",
	}

	for _, want := range checks {
		if !strings.Contains(md, want) {
			t.Errorf("prompt missing %q", want)
		}
	}
}

func TestGenerateMarkdownPrompt_emptyOptionalSections(t *testing.T) {
	err := &analyzer.DetectedError{
		ErrorType:    "Error",
		ErrorMessage: "boom",
		FilePath:     "main.unknown",
		LineNumber:   1,
	}
	md := GenerateMarkdownPrompt(err, "", nil)

	if !strings.Contains(md, "(no stack trace available)") {
		t.Error("expected empty stack trace placeholder")
	}
	if !strings.Contains(md, "(snippet unavailable)") {
		t.Error("expected empty snippet placeholder")
	}
	if !strings.Contains(md, "(none)") {
		t.Error("expected empty last logs placeholder")
	}
	if !strings.Contains(md, "```text") {
		t.Error("expected text language tag for unknown extension")
	}
}

func makeLines(n int) []string {
	lines := make([]string, n)
	for i := range lines {
		lines[i] = strings.Repeat("l", i+1)
	}
	return lines
}
