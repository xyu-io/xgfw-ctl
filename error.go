package xgfw_ctl

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

// Adds the output of stderr to exec.ExitError
type Error struct {
	exec.ExitError
	Cmd        exec.Cmd
	Msg        string
	ExitStatus *int //for overriding
}

func (e *Error) AlreadyExit() int {
	if e.ExitStatus != nil {
		return *e.ExitStatus
	}
	return e.Sys().(syscall.WaitStatus).ExitStatus()
}

func (e *Error) Error() string {
	return fmt.Sprintf("running %v, exit status %v: %v", e.Cmd.Args, e.AlreadyExit(), e.Msg)
}

var isNotExistPatterns = []string{
	"Bad rule (does a matching rule exist in that chain?).\n",
	"No chain/target/match by that name.\n",
	"No such file or directory",
	"does not exist",
}

// IsNotExist returns true if the error is due to the chain or rule not existing
func (e *Error) IsNotExist() bool {
	for _, str := range isNotExistPatterns {
		if strings.Contains(e.Msg, str) {
			return true
		}
	}
	return false
}
