//go:build windows

package runner

import "os/exec"

// terminateChild forcefully stops the child process.
// Windows does not support POSIX signals to child processes; Kill() maps to
// TerminateProcess and is always abrupt—there is no true SIGTERM grace period.
func terminateChild(cmd *exec.Cmd, done <-chan error) {
	if cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
	<-done
}
