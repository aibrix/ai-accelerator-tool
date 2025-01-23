package utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecPipeCmd(t *testing.T) {
	tests := []struct {
		name      string
		cmds      []string
		mockCmds  map[string]string
		want      string
		wantError bool
	}{
		{
			name: "valid pipe command",
			cmds: []string{"nvidia-smi -L", "wc -l"},
			mockCmds: map[string]string{
				"nvidia-smi -L | wc -l": "4",
			},
			want:      "4",
			wantError: false,
		},
		{
			name:      "empty command",
			cmds:      []string{""},
			wantError: true,
		},
		{
			name:      "unsafe command",
			cmds:      []string{"rm -rf /"},
			wantError: true,
		},
		{
			name: "multiple pipe commands",
			cmds: []string{"nvidia-smi -L", "grep GPU", "wc -l"},
			mockCmds: map[string]string{
				"nvidia-smi -L | grep GPU | wc -l": "2",
			},
			want:      "2",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockExecCmd{PipeCommands: tt.mockCmds}
			cleanup := SetExecPipeCmd(mock.ExecPipe)
			defer cleanup()

			got, err := ExecPipeCmd(context.Background(), tt.cmds)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestExecCmd(t *testing.T) {
	tests := []struct {
		name      string
		cmd       string
		args      []string
		mockCmds  map[string]string
		want      string
		wantError bool
	}{
		{
			name: "valid command",
			cmd:  "nvidia-smi",
			args: []string{"--query-gpu=name", "--format=csv,noheader"},
			mockCmds: map[string]string{
				"nvidia-smi --query-gpu=name --format=csv,noheader": "NVIDIA L20",
			},
			want:      "NVIDIA L20",
			wantError: false,
		},
		{
			name:      "empty command",
			cmd:       "",
			args:      []string{},
			wantError: true,
		},
		{
			name:      "unsafe command",
			cmd:       "rm",
			args:      []string{"-rf", "/"},
			wantError: true,
		},
		{
			name: "command with multiple args",
			cmd:  "lspci",
			args: []string{"-v", "-d", "10de:"},
			mockCmds: map[string]string{
				"lspci -v -d 10de:": "NVIDIA Corporation",
			},
			want:      "NVIDIA Corporation",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockExecCmd{Commands: tt.mockCmds}
			cleanup := SetExecCmd(mock.Exec)
			defer cleanup()

			got, err := ExecCmd(context.Background(), tt.cmd, tt.args)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCommandExists(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    bool
	}{
		{
			name:    "existing command - ls",
			command: "ls",
			want:    true,
		},
		{
			name:    "non-existent command",
			command: "non-existent-command-12345",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CommandExists(tt.command)
			assert.Equal(t, tt.want, got)
		})
	}
}
