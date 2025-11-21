package cmd

import (
	"fmt"
	"time"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/spf13/cobra"
)

// frontierMomentumCmd displays frontier momentum information
var frontierMomentumCmd = &cobra.Command{
	Use:   "frontierMomentum",
	Short: "Display current frontier momentum information",
	Long:  `Display information about the current frontier momentum (latest confirmed block).`,
	RunE:  runFrontierMomentum,
}

func init() {
	rootCmd.AddCommand(frontierMomentumCmd)
}

func runFrontierMomentum(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()

	// Connect to node
	rpcClient, err := client.New(cfg.Node.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to node: %w", err)
	}
	defer func() { _ = rpcClient.Close() }()

	// Get frontier momentum
	momentum, err := rpcClient.LedgerApi.GetFrontierMomentum()
	if err != nil {
		return fmt.Errorf("failed to get frontier momentum: %w", err)
	}

	// Display momentum info
	fmt.Println("Frontier Momentum:")
	fmt.Printf("  Height: %s\n", format.Green(fmt.Sprintf("%d", momentum.Height)))
	fmt.Printf("  Hash: %s\n", format.Cyan(momentum.Hash.String()))
	fmt.Printf("  Producer: %s\n", format.Cyan(momentum.Producer.String()))
	if momentum.Timestamp != nil {
		fmt.Printf("  Timestamp: %s\n", momentum.Timestamp.Format(time.RFC3339))
	}

	return nil
}
