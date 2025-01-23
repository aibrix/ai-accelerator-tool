package diagnose

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aibrix/ai-accelerator-tool/pkg/utils"
)

func (c *controller) checkNVIDIA(ctx context.Context) (map[GPUUID][]*DiagnoseResult, error) {
	results := map[GPUUID][]*DiagnoseResult{}

	// 1. Check driver status
	resGPUDriver, err := c.checkNVIDIAGPUDriverStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("checkNVIDIAGPUDriverStatus() failed: %s", err)
	}
	results[GPUUUIDOverall] = []*DiagnoseResult{resGPUDriver}
	if !*resGPUDriver.IsHealthy {
		return results, nil
	}

	// 2. Get GPU count
	gpuCardCount, err := c.getNVIDIAGPUCardCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("getNVIDIAGPUCardCount() failed: %s", err)
	}

	// 3. Check card count BEFORE getting UUIDs
	resCard, err := c.checkNVIDIACardCount(gpuCardCount)
	if err != nil {
		return nil, fmt.Errorf("checkNVIDIACardCount() failed: %s", err)
	}
	results[GPUUUIDOverall] = append(results[GPUUUIDOverall], resCard)
	if !*resCard.IsHealthy {
		return results, nil
	}

	// 4. Only get GPU UUIDs if card count is correct
	c.gpuIDs, err = c.getNVIDIAGPUsID(ctx, gpuCardCount)
	if err != nil {
		return nil, fmt.Errorf("getNVIDIAGPUsID() failed: %s", err)
	}

	// CHECK: GPU Link Status.
	err = c.checkNVIDIAGPUsLinkStatus(ctx, results)
	if err != nil {
		return nil, fmt.Errorf("checkNVIDIACard failed: %s", err)
	}

	// CHECK: GPU VRAM Status.
	err = c.checkNVIDIAVRAMStatus(ctx, results)
	if err != nil {
		return nil, fmt.Errorf("checkNVIDIACard failed: %s", err)
	}

	// CHECK: GPU Other Status.

	return results, nil
}

func (c *controller) checkNVIDIAGPUDriverStatus(ctx context.Context) (*DiagnoseResult, error) {
	res, err := utils.ExecCmd(ctx, "nvidia-smi", []string{"-L"})
	if err != nil {
		return &DiagnoseResult{
			Name:      DiagnoseGPUDriverStatus,
			IsHealthy: utils.BoolPtr(false),
			Message:   fmt.Sprintf("checkNVIDIAGPUDriverStatus() failed: %s", err),
		}, nil
	}

	if strings.Contains(res, "NVIDIA-SMI has failed because it couldn't communicate with the NVIDIA driver") {
		return &DiagnoseResult{
			Name:      DiagnoseGPUDriverStatus,
			IsHealthy: utils.BoolPtr(false),
			Message:   res,
		}, nil
	}

	return &DiagnoseResult{
		Name:      DiagnoseGPUDriverStatus,
		IsHealthy: utils.BoolPtr(true),
		Message:   "GPU Driver is loaded successfully",
	}, nil
}

func (c *controller) checkNVIDIACardCount(cardCount int) (*DiagnoseResult, error) {
	if c.ExpectedCardCount != cardCount {
		return &DiagnoseResult{
			Name:      DiagnoseGPUCardCount,
			IsHealthy: utils.BoolPtr(false),
			Message:   fmt.Sprintf("GPU Card Count: %d, Expected: %d", cardCount, c.ExpectedCardCount),
		}, fmt.Errorf("GPU card count mismatch: got %d, expected %d", cardCount, c.ExpectedCardCount)
	}

	return &DiagnoseResult{
		Name:      DiagnoseGPUCardCount,
		IsHealthy: utils.BoolPtr(true),
		Message:   fmt.Sprintf("GPU Card Count: %d", cardCount),
	}, nil
}

func (c *controller) checkNVIDIAGPUsLinkStatus(ctx context.Context, results map[GPUUID][]*DiagnoseResult) error {
	for i := 0; i < c.ExpectedCardCount; i++ {
		gpuID := c.gpuIDs[i]
		msg := ""

		err := c.checkNVIDIAGPULinkStatus(ctx, i)
		if err != nil {
			msg = fmt.Sprintf("Link is not OK: %s", err)
		}

		results[gpuID] = append(results[gpuID], &DiagnoseResult{
			Name:      DiagnoseGPULinkStatus,
			IsHealthy: utils.BoolPtr(err == nil),
			Message:   msg,
		})
	}

	return nil
}

func (c *controller) checkNVIDIAVRAMStatus(ctx context.Context, results map[GPUUID][]*DiagnoseResult) error {
	for i := 0; i < c.ExpectedCardCount; i++ {
		gpuID := c.gpuIDs[i]
		msg := ""

		// CHECK: VRAM Unrecoverable Errors.
		err := c.checkNVIDIAVRAMUnrecoverableErrors(ctx, i)
		if err != nil {
			msg = fmt.Sprintf("VRAM Unrecoverable Errors: %s", err)
		}
		results[gpuID] = append(results[gpuID], &DiagnoseResult{
			Name:      DiagnoseGPUnrecoverableErrors,
			IsHealthy: utils.BoolPtr(err == nil),
			Message:   msg,
		})

		// CHECK: VRAM Recoverable Errors.
		err = c.checkNVIDIAVRAMRecoverableErrors(ctx, i)
		msg = ""
		if err != nil {
			msg = fmt.Sprintf("VRAM Recoverable Errors: %s", err)
		}
		results[gpuID] = append(results[gpuID], &DiagnoseResult{
			Name:      DiagnoseGPURecoverableErrors,
			IsHealthy: utils.BoolPtr(err == nil),
			Message:   msg,
		})
	}

	return nil
}

func (c *controller) checkNVIDIAGPULinkStatus(ctx context.Context, cardIdx int) error {
	maxLinkWidth, err := utils.ExecCmd(ctx, "nvidia-smi", []string{"-i", strconv.Itoa(cardIdx), "--query-gpu=pcie.link.width.max", "--format=csv,noheader"})
	if err != nil {
		return fmt.Errorf("get max link width failed: %s", err)
	}
	maxLinkWidth = strings.TrimSpace(maxLinkWidth)

	curLinkWidth, err := utils.ExecCmd(ctx, "nvidia-smi", []string{"-i", strconv.Itoa(cardIdx), "--query-gpu=pcie.link.width.current", "--format=csv,noheader"})
	if err != nil {
		return fmt.Errorf("get current link width failed: %s", err)
	}
	curLinkWidth = strings.TrimSpace(curLinkWidth)

	if maxLinkWidth != curLinkWidth {
		return fmt.Errorf("link width is not ok, max: %s, current: %s", maxLinkWidth, curLinkWidth)
	}

	return nil
}

func (c *controller) checkNVIDIAGPUECCModeEnabled(ctx context.Context, cardIdx int) (bool, error) {
	res, err := utils.ExecCmd(ctx, "nvidia-smi", []string{"-i", strconv.Itoa(cardIdx), "--query-gpu=ecc.mode.current", "--format=csv,noheader"})
	if err != nil {
		return false, fmt.Errorf("get ecc mode failed: %s", err)
	}

	if strings.TrimSpace(res) != "Enabled" {
		return false, nil
	}

	return true, nil
}

func (c *controller) getNVIDIAGPUCardCount(ctx context.Context) (int, error) {
	countStr, err := utils.ExecPipeCmd(ctx, []string{"nvidia-smi -L", "wc -l"})
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		return 0, fmt.Errorf("strconv.Atoi() failed for %q: %s", countStr, err)
	}

	return count, nil
}

func (c *controller) getNVIDIAGPUsID(ctx context.Context, cardCount int) (map[int]GPUUID, error) {
	res := map[int]GPUUID{}
	for i := 0; i < cardCount; i++ {
		idStr, err := utils.ExecCmd(ctx, "nvidia-smi", []string{"-i", strconv.Itoa(i), "--query-gpu=uuid", "--format=csv,noheader"})
		if err != nil {
			return nil, fmt.Errorf("get uuid for idx %d failed: %s", i, err)
		}
		res[i] = GPUUID(strings.TrimSpace(idStr))
	}

	return res, nil
}

func (c *controller) checkNVIDIAVRAMUnrecoverableErrors(ctx context.Context, cardIdx int) error {
	// Check VRAM Page Retirement.
	resRetirement, err := utils.ExecPipeCmd(ctx, []string{"nvidia-smi -i " + strconv.Itoa(cardIdx) + " --query-retired-pages=retired_pages.address,retired_pages.cause --format=csv,noheader", "grep -i 'Double Bit ECC'"})
	if err != nil {
		return fmt.Errorf("get retired pages failed: %s", err)
	}
	retiredPages := strings.Split(strings.TrimSpace(resRetirement), "\n")
	for _, page := range retiredPages {
		page = strings.TrimSpace(page)
		if !strings.Contains(page, "N/A") {
			return fmt.Errorf("found retired page: %s", page)
		}
	}

	// Check ECC Errors.
	enabled, err := c.checkNVIDIAGPUECCModeEnabled(ctx, cardIdx)
	if err != nil {
		return fmt.Errorf("check ecc mode failed: %s", err)
	}
	if !enabled {
		return nil
	}
	eccCounts, err := utils.ExecCmd(ctx, "nvidia-smi", []string{"-i", strconv.Itoa(cardIdx), "--query-gpu=ecc.errors.uncorrected.volatile.total", "--format=csv,noheader"})
	if err != nil {
		return fmt.Errorf("get ecc errors failed: %s", err)
	}
	counts := strings.TrimSpace(eccCounts)
	if strings.Contains(counts, "N/A") || counts == "0" {
		return nil
	}

	return fmt.Errorf("found ecc errors: %s", counts)
}

func (c *controller) checkNVIDIAVRAMRecoverableErrors(ctx context.Context, cardIdx int) error {
	// Check VRAM Page Retirement.
	resRetirement, err := utils.ExecPipeCmd(ctx, []string{
		"nvidia-smi -i " + strconv.Itoa(cardIdx) + " --query-retired-pages=retired_pages.address,retired_pages.cause --format=csv,noheader",
		"grep -i 'Single Bit ECC'"})
	if err != nil {
		return fmt.Errorf("get retired pages failed: %s", err)
	}
	retiredPages := strings.Split(strings.TrimSpace(resRetirement), "\n")
	for _, page := range retiredPages {
		page = strings.TrimSpace(page)
		if !strings.Contains(page, "N/A") {
			return fmt.Errorf("found retired page: %s", page)
		}
	}

	// Check ECC Errors.
	enabled, err := c.checkNVIDIAGPUECCModeEnabled(ctx, cardIdx)
	if err != nil {
		return fmt.Errorf("check ecc mode failed: %s", err)
	}
	if !enabled {
		return nil
	}
	eccCounts, err := utils.ExecCmd(ctx,
		"nvidia-smi",
		[]string{"-i", strconv.Itoa(cardIdx), "--query-gpu=ecc.errors.corrected.volatile.total", "--format=csv,noheader"})
	if err != nil {
		return fmt.Errorf("get ecc errors failed: %s", err)
	}
	counts := strings.TrimSpace(eccCounts)
	if strings.Contains(counts, "N/A") || counts == "0" {
		return nil
	}

	return fmt.Errorf("found ecc errors: %s", counts)
}
