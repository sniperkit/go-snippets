package rioexecclient

import (
	"os/exec"
	"syscall"

	. "github.com/warpfork/go-errcat"

	"go.polydawn.net/go-timeless-api/rio"
)

func waitFor(cmd *exec.Cmd) (int, error) {
	err := cmd.Wait()
	if err == nil {
		return 0, nil
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return -1, Errorf(rio.ErrRPCBreakdown, "fork rio: unknown wait error: %s", err)
	}
	waitStatus, ok := exitErr.ProcessState.Sys().(syscall.WaitStatus)
	if !ok {
		return -1, Errorf(rio.ErrRPCBreakdown, "fork rio: unknown process state implementation %T", exitErr.ProcessState.Sys())
	}
	if waitStatus.Exited() {
		return waitStatus.ExitStatus(), nil
	} else if waitStatus.Signaled() {
		return int(waitStatus.Signal()) + 128, Errorf(rio.ErrRPCBreakdown, "fork rio: process killed with signal %d", waitStatus.Signal())
	} else {
		return -1, Errorf(rio.ErrRPCBreakdown, "fork rio: unknown process wait status (%#v)", waitStatus)
	}
}
