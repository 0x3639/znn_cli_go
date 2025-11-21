package cmd

import (
	"fmt"

	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/spf13/cobra"
)

var (
	// Version information (set via ldflags during build)
	Version   = "0.1.0-dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long:  `Display version information for the Zenon CLI, SDK, and connected daemon.`,
	RunE:  runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) error {
	// Display CLI version
	fmt.Printf("Zenon CLI Version: %s\n", format.Green(Version))
	fmt.Printf("Git Commit: %s\n", GitCommit)
	fmt.Printf("Build Date: %s\n", BuildDate)
	fmt.Printf("SDK Version: %s\n", format.Cyan("github.com/0x3639/znn-sdk-go v0.1.6"))

	// Try to get daemon version
	cfg := GetConfig()
	if cfg != nil && cfg.Node.URL != "" {
		fmt.Printf("\nConnecting to node at %s...\n", cfg.Node.URL)

		// Note: Daemon version requires RPC connection
		// This will be implemented when we add client connectivity checks
		format.Info("Daemon version check not yet implemented")
	}

	return nil
}
