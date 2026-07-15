//go:build unix

package runner

import (
	"os/exec"
	"syscall"
	"time"
)

// terminateChild asks the child to shut down gracefully before force-killing it.
// Unix processes can handle SIGTERM and release resources; after 3 seconds we
// escalate to SIGKILL via Kill().
func terminateChild(cmd *exec.Cmd, done <-chan error) {
	if cmd.Process == nil {
		<-done
		return
	}
	_ = cmd.Process.Signal(syscall.SIGTERM)
	select {
	case <-done:
		// Child exited after SIGTERM.
	case <-time.After(3 * time.Second):
		_ = cmd.Process.Kill()
		<-done
	}
}
