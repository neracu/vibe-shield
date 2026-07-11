package runner_test

import (
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"github.com/vibe-shield/vibe-shield/internal/runner"
)

func TestRunSignalInterrupt(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("programmatic os.Interrupt is unreliable on Windows; verify Ctrl+C manually")
	}

	cmd := exec.Command("sleep", "3600")
	go func() {
		time.Sleep(200 * time.Millisecond)
		p, err := os.FindProcess(os.Getpid())
		if err != nil {
			t.Errorf("FindProcess: %v", err)
			return
		}
		_ = p.Signal(os.Interrupt)
	}()

	interrupted, err := runner.Run(cmd)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !interrupted {
		t.Fatal("expected interrupted=true")
	}
}

func TestRunNormalExit(t *testing.T) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "exit", "0")
	} else {
		cmd = exec.Command("true")
	}

	interrupted, err := runner.Run(cmd)
	if interrupted {
		t.Fatal("expected interrupted=false")
	}
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
}
