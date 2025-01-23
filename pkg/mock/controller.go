package mock

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aibrix/ai-accelerator-tool/pkg/mock/resources"
)

type Config struct {
	ConfigPath    string
	GPUMockDir    string
	LDPreloadFile string
}

type Controller struct {
	config *Config
	mu     sync.Mutex
	active bool
	// Track the temporary lib path if we extract the embedded one
	tempLibPath string
	// Track original preload content to restore it
	originalPreloadContent string
}

func NewController(config *Config) (*Controller, error) {
	if config.ConfigPath == "" {
		return nil, fmt.Errorf("config path is required")
	}

	return &Controller{
		config: config,
		active: false,
	}, nil
}

func (c *Controller) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.active {
		return fmt.Errorf("mock environment already started")
	}

	// Verify config file exists
	if _, err := os.Stat(c.config.ConfigPath); err != nil {
		return fmt.Errorf("config file not found: %v", err)
	}

	// Create parent directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(c.config.GPUMockDir), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %v", err)
	}

	// Extract embedded library
	err := os.Mkdir(c.config.GPUMockDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create mock dir: %v", err)
	}

	// Add cleanup in case of error
	defer func() {
		if err != nil {
			os.RemoveAll(c.config.GPUMockDir)
		}
	}()

	libData, err := resources.GetInjectionLibrary()
	if err != nil {
		return fmt.Errorf("failed to get embedded library: %v", err)
	}

	libPath := filepath.Join(c.config.GPUMockDir, "nvml_injectiond.so")
	if err := os.WriteFile(libPath, libData, 0755); err != nil {
		return fmt.Errorf("failed to write library: %v", err)
	}

	c.tempLibPath = libPath

	// Verify injection library exists
	if _, err := os.Stat(libPath); err != nil {
		return fmt.Errorf("injection library not found: %v", err)
	}

	// Save original preload state
	preloadContent, err := os.ReadFile(c.config.LDPreloadFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to read preload file: %v", err)
		}
		// File doesn't exist, that's fine
		preloadContent = []byte{}
	}
	c.originalPreloadContent = string(preloadContent)

	// Add our library to preload
	newContent := c.originalPreloadContent
	if newContent != "" && !strings.HasSuffix(newContent, "\n") {
		newContent += "\n"
	}
	newContent += libPath + "\n"

	// Write updated preload content
	if err := os.WriteFile(c.config.LDPreloadFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to update preload file: %v", err)
	}

	// Copy config file
	destPath := filepath.Join(c.config.GPUMockDir, "gpu_mock_conf.toml")
	configData, err := os.ReadFile(c.config.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	if err := os.WriteFile(destPath, configData, 0644); err != nil {
		return fmt.Errorf("failed to copy config file: %v", err)
	}

	c.active = true
	return nil
}

func (c *Controller) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.active {
		return nil
	}

	// Clean up preload file
	if c.tempLibPath != "" {
		if c.originalPreloadContent == "" {
			// If file didn't exist originally, remove it
			if err := os.Remove(c.config.LDPreloadFile); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to remove preload file: %v", err)
			}
		} else {
			// Restore original content
			if err := os.WriteFile(c.config.LDPreloadFile, []byte(c.originalPreloadContent), 0644); err != nil {
				return fmt.Errorf("failed to restore preload file: %v", err)
			}
		}

		// Clean up temporary files
		if err := os.RemoveAll(filepath.Dir(c.tempLibPath)); err != nil {
			return fmt.Errorf("failed to cleanup temp library: %v", err)
		}
		c.tempLibPath = ""
	}

	c.active = false
	return nil
}

func (c *Controller) IsActive() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.active
}
