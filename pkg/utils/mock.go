package utils

import (
	"context"
	"fmt"
	"strings"
)

// MockExecCmd is a mock implementation of command execution for testing
type MockExecCmd struct {
	Commands     map[string]string
	PipeCommands map[string]string
	Err          error
}

// Exec mocks the execution of a single command
func (m *MockExecCmd) Exec(_ context.Context, cmd string, args []string) (string, error) {
	if cmd == "" {
		return "", ErrEmptyCommand
	}
	if !isSafeCommand(cmd) {
		return "", ErrUnsafeCommand
	}
	cmdStr := cmd
	for _, arg := range args {
		cmdStr += " " + arg
	}
	if output, ok := m.Commands[cmdStr]; ok {
		return output, m.Err
	}
	return "", fmt.Errorf("command not found: %s", cmdStr)
}

// ExecPipe mocks the execution of piped commands
func (m *MockExecCmd) ExecPipe(_ context.Context, cmds []string) (string, error) {
	if len(cmds) == 0 {
		return "", ErrEmptyCommand
	}
	for _, cmd := range cmds {
		if cmd == "" {
			return "", ErrEmptyCommand
		}
		// Extract command name (first word before any spaces)
		cmdName := strings.Fields(cmd)[0]
		if !isSafeCommand(cmdName) {
			return "", ErrUnsafeCommand
		}
	}
	cmdStr := strings.Join(cmds, " | ")
	if output, ok := m.PipeCommands[cmdStr]; ok {
		return output, m.Err
	}
	return "", fmt.Errorf("command not found: %s", cmdStr)
}
