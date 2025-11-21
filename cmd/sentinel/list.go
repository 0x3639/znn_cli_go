package sentinel

import (
	"fmt"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/config"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/spf13/cobra"
)

// listCmd lists all sentinels
var listCmd = &cobra.Command{
	Use:   "list [pageIndex pageSize]",
	Short: "List all sentinels",
	Long: `List all registered sentinels in the network.

Shows:
  - Sentinel owner address
  - Registration timestamp
  - Status (active/revoked)

Optional pagination parameters:
  pageIndex - Page number (default: 0)
  pageSize  - Items per page (default: 25)`,
	Args: cobra.RangeArgs(0, 2),
	RunE: runList,
}

func init() {
	SentinelCmd.AddCommand(listCmd)
}

func runList(cmdCobra *cobra.Command, args []string) error {
	// Parse pagination
	pageIndex := uint32(0)
	pageSize := uint32(25)
	if len(args) >= 1 {
		fmt.Sscanf(args[0], "%d", &pageIndex)
	}
	if len(args) >= 2 {
		fmt.Sscanf(args[1], "%d", &pageSize)
	}

	// Get URL from flags or config
	url, _ := cmdCobra.Flags().GetString("url")
	configFile, _ := cmdCobra.Flags().GetString("config")

	cfg, err := config.Load(configFile)
	if err != nil {
		cfg = config.DefaultConfig()
	}
	if url != "" {
		cfg.Node.URL = url
	}

	// Connect to node
	rpcClient, err := client.New(cfg.Node.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to node: %w", err)
	}
	defer rpcClient.Close()

	// Get sentinel list
	sentinelList, err := rpcClient.SentinelApi.GetAllActive(pageIndex, pageSize)
	if err != nil {
		return fmt.Errorf("failed to get sentinel list: %w", err)
	}

	// Display results
	if sentinelList.Count == 0 {
		fmt.Println("No sentinels found")
		return nil
	}

	fmt.Printf("Total active sentinels: %d\n", sentinelList.Count)
	fmt.Println()

	for idx, sentinel := range sentinelList.List {
		rank := int(pageIndex)*int(pageSize) + idx + 1

		fmt.Printf("%d. Sentinel %s\n", rank, format.Green(sentinel.Owner.String()))
		fmt.Printf("   Registered at momentum: %d\n", sentinel.RegistrationTimestamp)
		fmt.Println()
	}

	return nil
}
