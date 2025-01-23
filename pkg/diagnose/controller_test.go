package diagnose

import (
	"context"
	"testing"

	"github.com/aibrix/ai-accelerator-tool/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewController(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			cfg: &Config{
				ExpectedCardCount: 4,
			},
			wantErr: false,
		},
		{
			name:    "nil config",
			cfg:     nil,
			wantErr: true,
			errMsg:  "config cannot be nil",
		},
		{
			name: "invalid card count",
			cfg: &Config{
				ExpectedCardCount: 0,
			},
			wantErr: true,
			errMsg:  "expected card count must be positive, got 0",
		},
		{
			name: "negative card count",
			cfg: &Config{
				ExpectedCardCount: -1,
			},
			wantErr: true,
			errMsg:  "expected card count must be positive, got -1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewController(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				if tt.errMsg != "" {
					assert.Equal(t, tt.errMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				c, ok := got.(*controller)
				assert.True(t, ok)
				assert.Equal(t, tt.cfg.ExpectedCardCount, c.ExpectedCardCount)
				assert.NotNil(t, c.gpuIDs)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	tests := []struct {
		name              string
		expectedCardCount int
		mockCmds          map[string]string
		mockPipeCmds      map[string]string
		want              map[GPUUID][]*DiagnoseResult
		wantErr           bool
		wantErrContains   string
	}{
		{
			name:              "healthy nvidia system",
			expectedCardCount: 2,
			mockCmds: map[string]string{
				"nvidia-smi -L":     "GPU 0: NVIDIA A100-SXM4-40GB\nGPU 1: NVIDIA A100-SXM4-40GB",
				"lspci -v -d 10de:": "NVIDIA Corporation GA102",
				"nvidia-smi --query-gpu=name --format=csv,noheader":                         "NVIDIA A100-SXM4-40GB",
				"nvidia-smi -i 0 --query-gpu=uuid --format=csv,noheader":                    "GPU-uuid-1",
				"nvidia-smi -i 1 --query-gpu=uuid --format=csv,noheader":                    "GPU-uuid-2",
				"nvidia-smi -i 0 --query-gpu=pcie.link.width.max --format=csv,noheader":     "16",
				"nvidia-smi -i 0 --query-gpu=pcie.link.width.current --format=csv,noheader": "16",
				"nvidia-smi -i 1 --query-gpu=pcie.link.width.max --format=csv,noheader":     "16",
				"nvidia-smi -i 1 --query-gpu=pcie.link.width.current --format=csv,noheader": "16",
			},
			mockPipeCmds: map[string]string{
				"nvidia-smi -L | wc -l": "2",
			},
			wantErr: false,
		},
		{
			name:              "no nvidia gpu",
			expectedCardCount: 2,
			mockCmds: map[string]string{
				"nvidia-smi -L":     "",
				"lspci -v -d 10de:": "",
			},
			wantErr: true,
		},
		{
			name:              "card count mismatch",
			expectedCardCount: 4,
			mockCmds: map[string]string{
				"nvidia-smi -L":     "GPU 0: NVIDIA A100-SXM4-40GB\nGPU 1: NVIDIA A100-SXM4-40GB",
				"lspci -v -d 10de:": "NVIDIA Corporation GA102",
				"nvidia-smi --query-gpu=name --format=csv,noheader": "NVIDIA A100-SXM4-40GB",
			},
			mockPipeCmds: map[string]string{
				"nvidia-smi -L | wc -l": "2",
			},
			wantErr:         true,
			wantErrContains: "GPU card count mismatch: got 2, expected 4",
		},
		{
			name:              "unsupported vendor",
			expectedCardCount: 2,
			mockCmds: map[string]string{
				"nvidia-smi -L":     "GPU 0: AMD GPU",
				"lspci -v -d 10de:": "AMD Corporation",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &utils.MockExecCmd{
				Commands:     tt.mockCmds,
				PipeCommands: tt.mockPipeCmds,
			}
			cleanup := utils.SetExecCmd(mock.Exec)
			cleanupPipe := utils.SetExecPipeCmd(mock.ExecPipe)
			defer cleanup()
			defer cleanupPipe()

			c, err := NewController(&Config{ExpectedCardCount: tt.expectedCardCount})
			assert.NoError(t, err)

			results, err := c.Check(context.Background())
			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrContains != "" {
					assert.Contains(t, err.Error(), tt.wantErrContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, results)
			assert.Contains(t, results, GPUUUIDOverall)

			// Only check results for healthy system
			if tt.name == "healthy nvidia system" {
				// Check driver status first
				driverResult := results[GPUUUIDOverall][0]
				assert.Equal(t, DiagnoseGPUDriverStatus, driverResult.Name)
				assert.True(t, *driverResult.IsHealthy)

				// Then check card count
				cardCountResult := results[GPUUUIDOverall][1]
				assert.Equal(t, DiagnoseGPUCardCount, cardCountResult.Name)
				assert.True(t, *cardCountResult.IsHealthy)

				// Check individual GPU results
				assert.Contains(t, results, GPUUID("GPU-uuid-1"))
				assert.Contains(t, results, GPUUID("GPU-uuid-2"))
			}
		})
	}
}

func TestPrint(t *testing.T) {
	tests := []struct {
		name    string
		results []*DiagnoseResult
		wantErr bool
	}{
		{
			name: "print results",
			results: []*DiagnoseResult{
				{
					Name:      DiagnoseGPUCardCount,
					IsHealthy: utils.BoolPtr(true),
					Message:   "GPU count matches expected",
				},
			},
			wantErr: false,
		},
		{
			name:    "print nil results",
			results: nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewController(&Config{ExpectedCardCount: 2})
			assert.NoError(t, err)

			err = c.Print(context.Background(), tt.results)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
