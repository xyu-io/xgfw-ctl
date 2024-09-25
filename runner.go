package xgfw_ctl

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"strings"
)

type Runner struct {
	path   string
	option []string
}

func NewRunner(option []string) *Runner {
	return &Runner{
		option: option,
	}
}

// getFirewallProcess returns the correct firewalld command .
func getFirewallProcess() string {
	return "firewall-cmd"
}

func (fire *Runner) Exec() (string, error) {
	cmd := getFirewallProcess()
	path, err := exec.LookPath(cmd)
	if err != nil {

		return "", err
	}
	fire.path = path

	if flag, _ := fire.isExistOldRules(fire.option); flag {
		return "success", nil
	}
	var stdout bytes.Buffer
	err = fire.runWithOutput(fire.option, &stdout)
	if err != nil {
		return "", err
	}
	out, err := Reload()
	if err != nil {
		return string(out), err
	}

	return stdout.String(), nil
}

func (fire *Runner) ExecArgs(args []string, ruleSpec ...string) (string, error) {
	cmd := getFirewallProcess()
	path, err := exec.LookPath(cmd)
	if err != nil {

		return "", err
	}
	fire.path = path

	if flag, _ := fire.isExistOldRules(args, ruleSpec...); flag {
		return "success", nil
	}
	var stdout bytes.Buffer
	err = fire.runWithOutput(args, &stdout)
	if err != nil {
		return "", err
	}
	out, err := Reload()
	if err != nil {
		return string(out), err
	}

	return stdout.String(), nil
}

// Checks if a rule specification exists for a table
func (fire *Runner) isExistOldRules(args []string, ruleSpec ...string) (bool, error) {
	var stdout bytes.Buffer

	arg := append(args, ruleSpec...)

	err := fire.runWithOutput(arg, &stdout)
	if err != nil {
		return false, err
	}

	rs := "ALREADY_ENABLED"
	return strings.Contains(stdout.String(), rs), nil
}

// runWithOutput runs an command with the given arguments,
// writing any stdout output to the given writer
func (fire *Runner) runWithOutput(args []string, stdout io.Writer) error {

	var stderr bytes.Buffer
	cmd := exec.Cmd{
		Path:   fire.path,
		Args:   append([]string{fire.path}, args...),
		Stdout: stdout,
		Stderr: &stderr,
	}

	if err := cmd.Run(); err != nil {
		var e *exec.ExitError
		switch {
		case errors.As(err, &e):
			return &Error{ExitError: *e, Cmd: cmd, Msg: stderr.String()}
		default:
			return err
		}
	}

	return nil
}

// Reload reloads firewalld for setting to take effect
func Reload() ([]byte, error) {
	cmd := exec.Command("firewall-cmd", "--reload")
	output, err := cmd.CombinedOutput()
	return output, err
}
