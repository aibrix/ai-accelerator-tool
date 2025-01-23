package diagnose

import "context"

type DiagnoseType string

const (
	DiagnoseGPUDriverStatus       DiagnoseType = "gpu_driver_status"
	DiagnoseGPUCardCount          DiagnoseType = "gpu_card_count"
	DiagnoseGPULinkStatus         DiagnoseType = "gpu_link_status"
	DiagnoseGPUnrecoverableErrors DiagnoseType = "gpu_vram_unrecoverable_errors"
	DiagnoseGPURecoverableErrors  DiagnoseType = "gpu_vram_recoverable_errors"
)

type GPUUID string

const (
	// GPUUUIDOverall is a special GPUUID that represents the overall GPU status.
	GPUUUIDOverall GPUUID = "OVERALL"
)

// Diagnoser defines a set of APIs for anomaly detection in the target environment.
type Diagnoser interface {
	// Check detects all GPU anomalies in the target environment.
	Check(context.Context) (map[GPUUID][]*DiagnoseResult, error)

	// Print outputs detection results to stdout in the specified format.
	Print(context.Context, []*DiagnoseResult) error
}

// DiagnoseResult defines the test output result.
type DiagnoseResult struct {
	Name      DiagnoseType
	IsHealthy *bool
	Message   string
}
