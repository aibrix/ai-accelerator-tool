package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckEnv(t *testing.T) {
	tests := []struct {
		name     string
		mockCmds map[string]string
		want     *Env
		wantErr  bool
	}{
		{
			name: "nvidia gpu exists",
			mockCmds: map[string]string{
				"lspci -v -d 10de:": "NVIDIA Corporation [L20]",
				"nvidia-smi --query-gpu=name --format=csv,noheader": "NVIDIA L20",
			},
			want: &Env{
				Vendor:  NvidiaVendor,
				GPUType: "NVIDIA L20",
			},
			wantErr: false,
		},
		{
			name: "no nvidia gpu",
			mockCmds: map[string]string{
				"lspci -v -d 10de:": "",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockExecCmd{Commands: tt.mockCmds}
			cleanup := SetExecCmd(mock.Exec)
			defer cleanup()

			got, err := CheckEnv(context.Background())
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, ErrNoNvidiaDevice, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestIsSafeCommand(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    bool
	}{
		{
			name:    "safe command - nvidia-smi",
			command: "nvidia-smi",
			want:    true,
		},
		{
			name:    "safe command - lspci",
			command: "lspci",
			want:    true,
		},
		{
			name:    "safe command - wc",
			command: "wc",
			want:    true,
		},
		{
			name:    "safe command - grep",
			command: "grep",
			want:    true,
		},
		{
			name:    "unsafe command - rm",
			command: "rm",
			want:    false,
		},
		{
			name:    "unsafe command - arbitrary",
			command: "arbitrary",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isSafeCommand(tt.command)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBoolPtr(t *testing.T) {
	tests := []struct {
		name string
		val  bool
		want bool
	}{
		{
			name: "true value",
			val:  true,
			want: true,
		},
		{
			name: "false value",
			val:  false,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BoolPtr(tt.val)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want, *got)
		})
	}
}

func TestPrettyPrint(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "print struct",
			data: struct {
				Name  string
				Value int
			}{
				Name:  "test",
				Value: 123,
			},
		},
		{
			name: "print map",
			data: map[string]interface{}{
				"key1": "value1",
				"key2": 42,
			},
		},
		{
			name: "print nil",
			data: nil,
		},
		{
			name: "print env",
			data: &Env{
				Vendor:  NvidiaVendor,
				GPUType: "NVIDIA L20",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Redirect stdout to capture output
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			PrettyPrint(tt.data)

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Verify output is valid JSON
			var js interface{}
			err := json.Unmarshal([]byte(output), &js)
			assert.NoError(t, err, "Output should be valid JSON")

			// For nil input, verify empty output
			if tt.data == nil {
				assert.Equal(t, "null\n", output)
			}
		})
	}
}
