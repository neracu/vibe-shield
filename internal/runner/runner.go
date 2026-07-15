package runner

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// Run starts cmd, waits for it to finish or for SIGINT/SIGTERM.
// On signal, the child is killed and the caught signal is returned with nil error.
func Run(cmd *exec.Cmd) (sig os.Signal, err error) {
	if err = cmd.Start(); err != nil {
		return nil, err
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case sig = <-sigCh:
		// Child shutdown is platform-specific (see runner_unix.go / runner_windows.go).
		// Windows has no POSIX SIGTERM: Process.Kill() maps to TerminateProcess and is always abrupt.
		// On Unix we can deliver SIGTERM first so the child may exit cleanly before SIGKILL.
		terminateChild(cmd, done)
		return sig, nil
	case err = <-done:
		return nil, err
	}
}
