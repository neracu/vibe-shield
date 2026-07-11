package ui

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

func PrintShielding(command string) {
	color.New(color.FgHiBlack, color.Faint).Fprintf(os.Stderr,
		"🛡️ [Vibe-Shield] Shielding your code session (running: %s)...\n", command)
}

func PrintCrashDetected(file string, line int) {
	color.New(color.FgRed, color.Bold).Fprintf(os.Stderr,
		"🚨 [Vibe-Shield] Crash detected in %s:%d!\n", filepath.Base(file), line)
}

func PrintClipboardSuccess() {
	color.New(color.FgCyan, color.Bold).Fprintln(os.Stderr,
		"📋 [Vibe-Shield] Surgical prompt successfully copied to clipboard. Paste it into your AI!")
}

func PrintSourceReadWarning() {
	color.New(color.FgYellow).Fprintln(os.Stderr,
		"⚠️ [Vibe-Shield] Could not read source file, context might be limited.")
}
