package diagnose

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aibrix/ai-accelerator-tool/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// mockExecCmd replaces the real ExecCmd for testing
type mockExecCmd struct {
	commands     map[string]string
	pipeCommands map[string]string
	err          error
}

func (m *mockExecCmd) exec(_ context.Context, cmd string, args []string) (string, error) {
	cmdStr := cmd
	for _, arg := range args {
		cmdStr += " " + arg
	}
	if output, ok := m.commands[cmdStr]; ok {
		return output, m.err
	}
	return "", fmt.Errorf("command not found: %s", cmdStr)
}

func (m *mockExecCmd) execPipe(_ context.Context, cmds []string) (string, error) {
	cmdStr := strings.Join(cmds, " | ")
	if output, ok := m.pipeCommands[cmdStr]; ok {
		return output, m.err
	}
	return "", fmt.Errorf("command not found: %s", cmdStr)
}

func TestCheckNVIDIACardCount(t *testing.T) {
	tests := []struct {
		name              string
		expectedCardCount int
		actualCardCount   int
		wantHealthy       bool
		wantErr           bool
		wantErrContains   string
	}{
		{
			name:              "matching card count",
			expectedCardCount: 4,
			actualCardCount:   4,
			wantHealthy:       true,
			wantErr:           false,
		},
		{
			name:              "mismatched card count",
			expectedCardCount: 4,
			actualCardCount:   3,
			wantHealthy:       false,
			wantErr:           true,
			wantErrContains:   "GPU card count mismatch: got 3, expected 4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &controller{
				ExpectedCardCount: tt.expectedCardCount,
			}
			result, err := c.checkNVIDIACardCount(tt.actualCardCount)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrContains != "" {
					assert.Contains(t, err.Error(), tt.wantErrContains)
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantHealthy, *result.IsHealthy)
			assert.Equal(t, DiagnoseGPUCardCount, result.Name)
		})
	}
}

func TestCheckNVIDIAGPULinkStatus(t *testing.T) {
	tests := []struct {
		name     string
		mockCmds map[string]string
		wantErr  bool
		cardIdx  int
	}{
		{
			name: "healthy link status",
			mockCmds: map[string]string{
				"nvidia-smi -i 0 --query-gpu=pcie.link.width.max --format=csv,noheader":     "16",
				"nvidia-smi -i 0 --query-gpu=pcie.link.width.current --format=csv,noheader": "16",
			},
			cardIdx: 0,
			wantErr: false,
		},
		{
			name: "degraded link status",
			mockCmds: map[string]string{
				"nvidia-smi -i 0 --query-gpu=pcie.link.width.max --format=csv,noheader":     "16",
				"nvidia-smi -i 0 --query-gpu=pcie.link.width.current --format=csv,noheader": "8",
			},
			cardIdx: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockExecCmd{commands: tt.mockCmds}
			cleanup := utils.SetExecCmd(mock.exec)
			defer cleanup()

			c := &controller{}
			err := c.checkNVIDIAGPULinkStatus(context.Background(), tt.cardIdx)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckNVIDIAVRAMUnrecoverableErrors(t *testing.T) {
	tests := []struct {
		name         string
		mockCmds     map[string]string
		mockPipeCmds map[string]string
		wantErr      bool
		cardIdx      int
	}{
		{
			name: "no errors with ECC enabled",
			mockCmds: map[string]string{
				"nvidia-smi -i 0 --query-gpu=ecc.mode.current --format=csv,noheader":                      "Enabled",
				"nvidia-smi -i 0 --query-gpu=ecc.errors.uncorrected.volatile.total --format=csv,noheader": "0",
			},
			mockPipeCmds: map[string]string{
				"nvidia-smi -i 0 --query-retired-pages=retired_pages.address,retired_pages.cause --format=csv,noheader | grep -i 'Double Bit ECC'": "N/A",
			},
			cardIdx: 0,
			wantErr: false,
		},
		{
			name: "unrecoverable errors present",
			mockCmds: map[string]string{
				"nvidia-smi -i 0 --query-gpu=ecc.mode.current --format=csv,noheader":                      "Enabled",
				"nvidia-smi -i 0 --query-gpu=ecc.errors.uncorrected.volatile.total --format=csv,noheader": "1",
			},
			cardIdx: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockExecCmd{
				commands:     tt.mockCmds,
				pipeCommands: tt.mockPipeCmds,
			}
			cleanup := utils.SetExecCmd(mock.exec)
			cleanupPipe := utils.SetExecPipeCmd(mock.execPipe)
			defer cleanup()
			defer cleanupPipe()

			c := &controller{}
			err := c.checkNVIDIAVRAMUnrecoverableErrors(context.Background(), tt.cardIdx)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckNVIDIAVRAMRecoverableErrors(t *testing.T) {
	tests := []struct {
		name         string
		mockCmds     map[string]string
		mockPipeCmds map[string]string
		wantErr      bool
		cardIdx      int
	}{
		{
			name: "no recoverable errors",
			mockCmds: map[string]string{
				"nvidia-smi -i 0 --query-gpu=ecc.mode.current --format=csv,noheader":                    "Enabled",
				"nvidia-smi -i 0 --query-gpu=ecc.errors.corrected.volatile.total --format=csv,noheader": "0",
			},
			mockPipeCmds: map[string]string{
				"nvidia-smi -i 0 --query-retired-pages=retired_pages.address,retired_pages.cause --format=csv,noheader | grep -i 'Single Bit ECC'": "N/A",
			},
			cardIdx: 0,
			wantErr: false,
		},
		{
			name: "recoverable errors present",
			mockCmds: map[string]string{
				"nvidia-smi -i 0 --query-gpu=ecc.mode.current --format=csv,noheader":                    "Enabled",
				"nvidia-smi -i 0 --query-gpu=ecc.errors.corrected.volatile.total --format=csv,noheader": "5",
			},
			mockPipeCmds: map[string]string{
				"nvidia-smi -i 0 --query-retired-pages=retired_pages.address,retired_pages.cause --format=csv,noheader | grep -i 'Single Bit ECC'": "N/A",
			},
			cardIdx: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockExecCmd{
				commands:     tt.mockCmds,
				pipeCommands: tt.mockPipeCmds,
			}
			cleanup := utils.SetExecCmd(mock.exec)
			cleanupPipe := utils.SetExecPipeCmd(mock.execPipe)
			defer cleanup()
			defer cleanupPipe()

			c := &controller{}
			err := c.checkNVIDIAVRAMRecoverableErrors(context.Background(), tt.cardIdx)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckNVIDIA(t *testing.T) {
	tests := []struct {
		name              string
		mockCmds          map[string]string
		mockPipeCmds      map[string]string
		expectedCardCount int
		wantErr           bool
	}{
		{
			name: "healthy GPU system",
			mockPipeCmds: map[string]string{
				"nvidia-smi -L | wc -l": "2",
			},
			mockCmds: map[string]string{
				"nvidia-smi -i 0 --query-gpu=uuid --format=csv,noheader":                                  "GPU-uuid-1",
				"nvidia-smi -i 1 --query-gpu=uuid --format=csv,noheader":                                  "GPU-uuid-2",
				"nvidia-smi -i 0 --query-gpu=pcie.link.width.max --format=csv,noheader":                   "16",
				"nvidia-smi -i 0 --query-gpu=pcie.link.width.current --format=csv,noheader":               "16",
				"nvidia-smi -i 1 --query-gpu=pcie.link.width.max --format=csv,noheader":                   "16",
				"nvidia-smi -i 1 --query-gpu=pcie.link.width.current --format=csv,noheader":               "16",
				"nvidia-smi -i 0 --query-gpu=ecc.mode.current --format=csv,noheader":                      "Enabled",
				"nvidia-smi -i 1 --query-gpu=ecc.mode.current --format=csv,noheader":                      "Enabled",
				"nvidia-smi -i 0 --query-gpu=ecc.errors.uncorrected.volatile.total --format=csv,noheader": "0",
				"nvidia-smi -i 1 --query-gpu=ecc.errors.uncorrected.volatile.total --format=csv,noheader": "0",
				"nvidia-smi -i 0 --query-gpu=ecc.errors.corrected.volatile.total --format=csv,noheader":   "0",
				"nvidia-smi -i 1 --query-gpu=ecc.errors.corrected.volatile.total --format=csv,noheader":   "0",
			},
			expectedCardCount: 2,
			wantErr:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockExecCmd{
				commands:     tt.mockCmds,
				pipeCommands: tt.mockPipeCmds,
			}
			cleanup := utils.SetExecCmd(mock.exec)
			cleanupPipe := utils.SetExecPipeCmd(mock.execPipe)
			defer cleanup()
			defer cleanupPipe()

			c := &controller{
				ExpectedCardCount: tt.expectedCardCount,
			}
			results, err := c.checkNVIDIA(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, results)
				assert.Contains(t, results, GPUUUIDOverall)
			}
		})
	}
}

func TestCheckNVIDIAGPUDriverStatus(t *testing.T) {
	tests := []struct {
		name     string
		mockCmds map[string]string
		want     *DiagnoseResult
		wantErr  bool
	}{
		{
			name: "driver loaded successfully",
			mockCmds: map[string]string{
				"nvidia-smi -L": "GPU 0: NVIDIA A100-SXM4-40GB (UUID: GPU-deadbeef-1234-5678-90ab-cdef01234567)",
			},
			want: &DiagnoseResult{
				Name:      DiagnoseGPUDriverStatus,
				IsHealthy: utils.BoolPtr(true),
				Message:   "GPU Driver is loaded successfully",
			},
			wantErr: false,
		},
		{
			name: "driver not loaded",
			mockCmds: map[string]string{
				"nvidia-smi -L": "NVIDIA-SMI has failed because it couldn't communicate with the NVIDIA driver",
			},
			want: &DiagnoseResult{
				Name:      DiagnoseGPUDriverStatus,
				IsHealthy: utils.BoolPtr(false),
				Message:   "NVIDIA-SMI has failed because it couldn't communicate with the NVIDIA driver",
			},
			wantErr: false,
		},
		{
			name:     "command execution error",
			mockCmds: map[string]string{
				// Don't provide mock response - will trigger error
			},
			want: &DiagnoseResult{
				Name:      DiagnoseGPUDriverStatus,
				IsHealthy: utils.BoolPtr(false),
				Message:   "checkNVIDIAGPUDriverStatus() failed: command not found: nvidia-smi -L",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mock := &utils.MockExecCmd{
				Commands: tt.mockCmds,
			}
			cleanup := utils.SetExecCmd(mock.Exec)
			defer cleanup()

			// Create controller instance and type assert to concrete type
			diagnoser, err := NewController(&Config{ExpectedCardCount: 1})
			assert.NoError(t, err)

			c, ok := diagnoser.(*controller)
			assert.True(t, ok, "Expected diagnoser to be *controller")

			// Execute test
			got, err := c.checkNVIDIAGPUDriverStatus(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, *tt.want.IsHealthy, *got.IsHealthy)
			assert.Equal(t, tt.want.Message, got.Message)
		})
	}
}
