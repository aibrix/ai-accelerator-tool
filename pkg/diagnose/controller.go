package diagnose

import (
	"context"
	"fmt"

	"github.com/aibrix/ai-accelerator-tool/pkg/utils"
)

type Config struct {
	ExpectedCardCount int
}

type controller struct {
	ExpectedCardCount int

	gpuIDs map[int]GPUUID
}

func NewController(cfg *Config) (Diagnoser, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if cfg.ExpectedCardCount <= 0 {
		return nil, fmt.Errorf("expected card count must be positive, got %d", cfg.ExpectedCardCount)
	}

	return &controller{
		ExpectedCardCount: cfg.ExpectedCardCount,
		gpuIDs:            make(map[int]GPUUID),
	}, nil
}

func (c *controller) Check(ctx context.Context) (map[GPUUID][]*DiagnoseResult, error) {
	env, err := utils.CheckEnv(ctx)
	if err != nil {
		return nil, err
	}

	switch env.Vendor {
	case utils.NvidiaVendor:
		return c.checkNVIDIA(ctx)
	default:
		return nil, fmt.Errorf("unsupported vendor: %s", env.Vendor)
	}
}

func (c *controller) Print(ctx context.Context, results []*DiagnoseResult) error {
	return nil
}
