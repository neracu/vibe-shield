package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/vibe-shield/vibe-shield/internal/analyzer"
)

func GenerateMarkdownPrompt(err *analyzer.DetectedError, snippet string, lastLogs []string) string {
	var b strings.Builder

	fmt.Fprintf(&b, "# Vibe-Shield Crash Report\n\n")
	fmt.Fprintf(&b, "| Field | Value |\n|-------|-------|\n")
	fmt.Fprintf(&b, "| OS | %s |\n", runtime.GOOS)
	fmt.Fprintf(&b, "| Runtime | %s |\n", runtimeLabel(err.FilePath))
	fmt.Fprintf(&b, "| Command | `%s` |\n", commandLine())
	fmt.Fprintf(&b, "| File | `%s:%d` |\n\n", err.FilePath, err.LineNumber)

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
