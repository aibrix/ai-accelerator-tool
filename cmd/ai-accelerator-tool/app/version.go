package app

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aibrix/ai-accelerator-tool/pkg/version"
)

func NewVersionCmd() *cobra.Command {
	var command = &cobra.Command{
		Use:   "version",
		Short: "Show tool version.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Version: %s\n", version.Version)
			return nil
		},
	}

	return command
}
