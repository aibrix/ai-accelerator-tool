package app

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ai-accelerator-tool",
	Short: "ai-accelerator-tool is a simple but powerful AI accelerator detection tool",
	Long: "ai-accelerator-tool is a AI accelerator detection tool that supports detection of AI accelerators from " +
		"different manufacturers",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(NewDiagnoseCmd())
	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.AddCommand(NewMockCmd())
}
