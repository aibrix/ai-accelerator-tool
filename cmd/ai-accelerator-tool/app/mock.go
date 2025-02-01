package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/aibrix/ai-accelerator-tool/pkg/mock"
)

func NewMockCmd() *cobra.Command {
	var configPath string
	var gpuMockDir string

	command := &cobra.Command{
		Use:   "mock",
		Short: "Run with GPU mock injection",
		Long:  `Run GPU diagnosis with mock injection using a configuration file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("--config is required")
			}

			// Get absolute path for config
			absConfigPath, err := filepath.Abs(configPath)
			if err != nil {
				return fmt.Errorf("failed to get absolute path for config: %v", err)
			}

			// Create context with cancellation
			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			// Handle interrupt signals
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-sigCh
				klog.InfoS("Received interrupt signal, cleaning up...")
				cancel()
			}()

			// Create mock controller with embedded library
			controller, err := mock.NewController(&mock.Config{
				ConfigPath: absConfigPath,
				GPUMockDir: gpuMockDir,

				// Setting LD preload file to /etc/ld.so.preload is a hack to make the mock work.
				LDPreloadFile: "/etc/ld.so.preload",
			})
			if err != nil {
				return fmt.Errorf("failed to create mock controller: %v", err)
			}

			// Start mock environment
			if err := controller.Start(); err != nil {
				return fmt.Errorf("failed to start mock: %v", err)
			}
			defer controller.Stop()

			errCh := make(chan error, 1)

			// Wait for either completion or cancellation
			select {
			case err := <-errCh:
				if err != nil {
					return fmt.Errorf("diagnose failed: %v", err)
				}
				klog.InfoS("Mock test completed successfully")
			case <-ctx.Done():
				return fmt.Errorf("mock test cancelled")
			}

			return nil
		},
	}

	command.Flags().StringVarP(&configPath, "config", "c", "", "Path to mock configuration file (required)")
	command.MarkFlagRequired("config")

	command.Flags().StringVarP(&gpuMockDir, "gpu-mock-dir", "d", "/opt/gpu_mock", "Directory for GPU mock files")

	return command
}
