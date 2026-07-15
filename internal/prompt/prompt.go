package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/neracu/vibe-shield/internal/analyzer"
)

const snippetStaleWarning = "> **Warning:** The source file was modified after the command started or after the crash was detected. This code snippet may not reflect the state at the time of the crash.\n\n"

func GenerateMarkdownPrompt(err *analyzer.DetectedError, snippet string, lastLogs []string, snippetStale bool) string {
	var b strings.Builder

	writeReportHeader(&b, [][2]string{
		{"Runtime", runtimeLabel(err.FilePath)},
		{"Command", fmt.Sprintf("`%s`", commandLine())},
		{"File", fmt.Sprintf("`%s:%d`", err.FilePath, err.LineNumber)},
	})

	fmt.Fprintf(&b, "### THE ERROR\n\n")
	fmt.Fprintf(&b, "**Type:** %s\n", err.ErrorType)
	fmt.Fprintf(&b, "**Message:** %s\n\n", err.ErrorMessage)
	fmt.Fprintf(&b, "```\n%s\n```\n\n", formatStackTrace(err.StackTrace))

	fmt.Fprintf(&b, "### CODE SNIPPET\n\n")
	fmt.Fprintf(&b, "```%s\n", languageTag(err.FilePath))
	if snippet != "" {
		b.WriteString(snippet)
	} else {
		b.WriteString("(snippet unavailable)")
	}
	fmt.Fprintf(&b, "\n```\n\n")
	if snippetStale {
		b.WriteString(snippetStaleWarning)
	}

	fmt.Fprintf(&b, "### LAST LOGS\n\n")
	fmt.Fprintf(&b, "```\n%s\n```\n\n", formatLastLogs(lastLogs))

	fmt.Fprintf(&b, "### INSTRUCTION FOR AI\n\n")
	b.WriteString("Fix only this specific error. Return only the fixed code or a minimal diff.\n")
	b.WriteString("Do not refactor unrelated code or break the existing architecture.\n")

	return b.String()
}

func commandLine() string {
	if len(os.Args) <= 1 {
		return ""
	}
	return strings.Join(os.Args[1:], " ")
}

func languageTag(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".js", ".mjs", ".cjs":
		return "javascript"
	case ".ts":
		return "typescript"
	case ".jsx":
		return "jsx"
	case ".tsx":
		return "tsx"
	case ".py":
		return "python"
	case ".go":
		return "go"
	default:
		return "text"
	}
}

func runtimeLabel(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".py":
		return "Python"
	case ".js", ".mjs", ".cjs", ".ts", ".jsx", ".tsx":
		return "Node.js"
	case ".go":
		return "Go"
	default:
		return "unknown"
	}
}

func formatStackTrace(lines []string) string {
	if len(lines) == 0 {
		return "(no stack trace available)"
	}
	return strings.Join(lines, "\n")
}

func formatLastLogs(lines []string) string {
	if len(lines) == 0 {
		return "(none)"
	}
	return strings.Join(lines, "\n")
}

func GenerateFallbackPrompt(exitCode int, stderrTail, stdoutTail []string) string {
	var b strings.Builder

	writeReportHeader(&b, [][2]string{
		{"Command", fmt.Sprintf("`%s`", commandLine())},
		{"Exit code", fmt.Sprintf("%d", exitCode)},
	})

	fmt.Fprintf(&b, "### THE ERROR\n\n")
	b.WriteString("No structured stack trace was parseable from stderr.\n\n")

	fmt.Fprintf(&b, "### STDERR (last lines)\n\n")
	fmt.Fprintf(&b, "```\n%s\n```\n\n", formatLastLogs(stderrTail))

	fmt.Fprintf(&b, "### LAST LOGS\n\n")
	fmt.Fprintf(&b, "```\n%s\n```\n\n", formatLastLogs(stdoutTail))

	fmt.Fprintf(&b, "### INSTRUCTION FOR AI\n\n")
	b.WriteString("Diagnose this failure from the stderr output and last logs above.\n")
	b.WriteString("No source file or line could be located automatically.\n")
	b.WriteString("Return only the minimal fix needed; do not refactor unrelated code.\n")

	return b.String()
}

func writeReportHeader(b *strings.Builder, extraRows [][2]string) {
	fmt.Fprintf(b, "# Vibe-Shield Crash Report\n\n")
	fmt.Fprintf(b, "| Field | Value |\n|-------|-------|\n")
	fmt.Fprintf(b, "| OS | %s |\n", runtime.GOOS)
	for _, row := range extraRows {
		fmt.Fprintf(b, "| %s | %s |\n", row[0], row[1])
	}
	b.WriteByte('\n')
}
