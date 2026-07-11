package runner

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// Run starts cmd, waits for it to finish or for SIGINT/SIGTERM.
// On signal, the child is killed and interrupted is true with nil error.
func Run(cmd *exec.Cmd) (interrupted bool, err error) {
	if err = cmd.Start(); err != nil {
		return false, err
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-sigCh:
		_ = cmd.Process.Kill()
		<-done
		return true, nil
	case err = <-done:
		return false, err
	}
}
