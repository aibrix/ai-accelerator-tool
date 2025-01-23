package utils

import (
	"context"
	"os/exec"
	"strings"
)

// ExecCmdFunc is a function type for command execution
type ExecCmdFunc func(context.Context, string, []string) (string, error)

// ExecPipeCmdFunc is a function type for piped command execution
type ExecPipeCmdFunc func(context.Context, []string) (string, error)

// ExecCmd is the default implementation
var ExecCmd ExecCmdFunc = realExecCmd

// ExecPipeCmd is the default implementation for piped commands
var ExecPipeCmd ExecPipeCmdFunc = realExecPipeCmd

// realExecCmd is the actual implementation that executes commands
func realExecCmd(ctx context.Context, cmdName string, args []string) (string, error) {
	cmd := exec.CommandContext(ctx, cmdName, args...)
	stdOutStderr, err := cmd.CombinedOutput()

	return string(stdOutStderr), err
}

// realExecPipeCmd is the actual implementation that executes piped commands
func realExecPipeCmd(ctx context.Context, cmds []string) (string, error) {
	for _, cmd := range cmds {
		if cmd == "" {
			return "", ErrEmptyCommand
		}

		cmdName := strings.Split(cmd, " ")[0]
		if !isSafeCommand(cmdName) {
			return "", ErrUnsafeCommand
		}
	}

	cmdStr := strings.Join(cmds, "|")
	cmd := exec.CommandContext(ctx, "sh", "-c", cmdStr)
	stdOutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(stdOutStderr)), nil
}

// SetExecCmd allows setting a mock implementation for testing
func SetExecCmd(mock ExecCmdFunc) func() {
	original := ExecCmd
	ExecCmd = mock
	return func() {
		ExecCmd = original
	}
}

// SetExecPipeCmd allows setting a mock implementation for testing
func SetExecPipeCmd(mock ExecPipeCmdFunc) func() {
	original := ExecPipeCmd
	ExecPipeCmd = mock
	return func() {
		ExecPipeCmd = original
	}
}
