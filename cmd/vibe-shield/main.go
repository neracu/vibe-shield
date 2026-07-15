package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/neracu/vibe-shield/internal/analyzer"
	"github.com/neracu/vibe-shield/internal/prompt"
	"github.com/neracu/vibe-shield/internal/runner"
	"github.com/neracu/vibe-shield/internal/stdcapture"
	"github.com/neracu/vibe-shield/internal/ui"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: vibe-shield <command> [args...]")
		os.Exit(1)
	}

	ui.PrintShielding(strings.Join(args, " "))

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin

	stderrCap := stdcapture.New()
	stdoutCap := stdcapture.NewStdout()
	cmd.Stdout = stdoutCap
	cmd.Stderr = stderrCap

	interrupted, err := runner.Run(cmd)
	if interrupted {
		os.Exit(0)
	}

	_ = stderrCap.Flush()
	_ = stdoutCap.Flush()

	if err != nil {
		exitCode := 1
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		}

		if detected, ok := analyzer.AnalyzeStderr(stderrCap.Lines()); ok {
			snippet, snippetErr := analyzer.ExtractCodeContext(detected)
			lastLogs := prompt.TailLogs(stdoutCap.Lines())
			md := prompt.GenerateMarkdownPrompt(detected, snippet, lastLogs)

			if clipErr := clipboard.WriteAll(md); clipErr != nil {
				fmt.Fprintf(os.Stderr, "vibe-shield: clipboard unavailable: %v\n", clipErr)
			}

			ui.PrintCrashDetected(detected.FilePath, detected.LineNumber)
			if snippetErr != nil {
				ui.PrintSourceReadWarning()
			}
			ui.PrintClipboardSuccess()
		} else {
			stderrTail := prompt.TailStderr(stderrCap.Lines())
			lastLogs := prompt.TailLogs(stdoutCap.Lines())
			md := prompt.GenerateFallbackPrompt(exitCode, stderrTail, lastLogs)

			if clipErr := clipboard.WriteAll(md); clipErr != nil {
				fmt.Fprintf(os.Stderr, "vibe-shield: clipboard unavailable: %v\n", clipErr)
			}

			ui.PrintFallbackCrashDetected()
			ui.PrintClipboardSuccess()
		}

		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "vibe-shield: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
