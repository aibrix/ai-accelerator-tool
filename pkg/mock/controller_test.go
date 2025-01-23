package mock

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewController(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		wantErr     bool
		errContains string
	}{
		{
			name:        "empty config path",
			config:      &Config{},
			wantErr:     true,
			errContains: "config path is required",
		},
		{
			name: "valid config",
			config: &Config{
				ConfigPath:    "testdata/config.toml",
				GPUMockDir:    "/tmp/gpu-mock",
				LDPreloadFile: "/tmp/ld.so.preload",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller, err := NewController(tt.config)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				if tt.errContains != "" && err != nil {
					if !strings.Contains(err.Error(), tt.errContains) {
						t.Errorf("error = %v, want error containing %v", err, tt.errContains)
					}
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if controller == nil {
				t.Error("expected controller to not be nil")
			}
		})
	}
}

func TestController_StartStop(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gpu-mock-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after test

	mockDir := filepath.Join(tempDir, "gpu_mock")
	configPath := filepath.Join(tempDir, "config.toml")
	preloadPath := filepath.Join(tempDir, "ld.so.preload")

	// Create a test config file
	if err := os.WriteFile(configPath, []byte("test config"), 0644); err != nil {
		t.Fatalf("failed to create test config: %v", err)
	}

	config := &Config{
		ConfigPath:    configPath,
		GPUMockDir:    mockDir,
		LDPreloadFile: preloadPath,
	}

	controller, err := NewController(config)
	if err != nil {
		t.Fatalf("failed to create controller: %v", err)
	}

	// Test Start
	if err := controller.Start(); err != nil {
		t.Fatalf("failed to start mock environment: %v", err)
	}

	if !controller.IsActive() {
		t.Error("controller should be active after Start()")
	}

	// Verify files were created
	if _, err := os.Stat(mockDir); err != nil {
		t.Errorf("mock directory was not created: %v", err)
	}

	if _, err := os.Stat(preloadPath); err != nil {
		t.Errorf("preload file was not created: %v", err)
	}

	// Verify config file was copied
	mockConfigPath := filepath.Join(mockDir, "gpu_mock_conf.toml")
	if _, err := os.Stat(mockConfigPath); err != nil {
		t.Errorf("mock config file was not created: %v", err)
	}

	// Test Stop
	if err := controller.Stop(); err != nil {
		t.Fatalf("failed to stop mock environment: %v", err)
	}

	if controller.IsActive() {
		t.Error("controller should not be active after Stop()")
	}

	// Verify cleanup
	if _, err := os.Stat(mockDir); !os.IsNotExist(err) {
		t.Error("mock directory was not cleaned up")
	}
}
