package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type VendorType string

const (
	NvidiaVendor VendorType = "nvidia"
)

type Env struct {
	Vendor  VendorType
	GPUType string
}

// CheckEnv checks the environment and returns the environment information.
// TODO: Support more vendors.
func CheckEnv(ctx context.Context) (*Env, error) {
	env, err := checkNVIDIA(ctx)
	if err != nil {
		return env, err
	}

	// TODO: Support more vendor.

	return env, nil
}

func isSafeCommand(cmd string) bool {
	safeCommands := []string{"lspci", "nvidia-smi", "wc", "ls", "grep", "echo", "cp"}
	for _, safeCmd := range safeCommands {
		if cmd == safeCmd {
			return true
		}
	}

	return false
}

// CommandExists checks if a command exists in the system's executable path.
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)

	return err == nil
}

func checkNVIDIA(ctx context.Context) (*Env, error) {
	args := []string{"-v", "-d", "10de:"}
	res, err := ExecCmd(ctx, "lspci", args)
	if err != nil {
		return nil, err
	}
	if res == "" {
		return nil, ErrNoNvidiaDevice
	}

	args = []string{"--query-gpu=name", "--format=csv,noheader"}
	gpuType, err := ExecCmd(ctx, "nvidia-smi", args)
	if err != nil {
		return nil, err
	}

	return &Env{
		Vendor:  NvidiaVendor,
		GPUType: strings.TrimSpace(gpuType),
	}, nil
}

func BoolPtr(b bool) *bool {
	return &b
}

func PrettyPrint(data interface{}) {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(prettyJSON))
}
