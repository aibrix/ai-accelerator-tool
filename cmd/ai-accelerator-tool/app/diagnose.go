package app

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/aibrix/ai-accelerator-tool/pkg/diagnose"
	"github.com/aibrix/ai-accelerator-tool/pkg/utils"
)

func NewDiagnoseCmd() *cobra.Command {
	var command = &cobra.Command{
		Use:   "diagnose",
		Short: "Check whether the GPU in the machine is abnormal.",
		RunE: func(cmd *cobra.Command, args []string) error {
			gpuCardCountStr := os.Getenv(utils.GPU_CARD_COUNT)
			if gpuCardCountStr == "" {
				return fmt.Errorf("%s is not set", utils.GPU_CARD_COUNT)
			}
			gpuCardCount, err := strconv.Atoi(gpuCardCountStr)
			if err != nil {
				return fmt.Errorf("%s is not a number: %v", utils.GPU_CARD_COUNT, gpuCardCountStr)
			}

			controller, err := diagnose.NewController(&diagnose.Config{
				ExpectedCardCount: gpuCardCount,
			})
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			res, err := controller.Check(ctx)
			if err != nil {
				return err
			}
			klog.InfoS("Diagnose Results")
			utils.PrettyPrint(res)

			return nil
		},
	}

	return command
}
