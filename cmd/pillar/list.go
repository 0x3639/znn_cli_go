package pillar

import (
	"fmt"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/config"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/spf13/cobra"
)

// listCmd lists all pillars
var listCmd = &cobra.Command{
	Use:   "list [pageIndex pageSize]",
	Short: "List all pillars",
	Long: `List all registered pillars in the network.

Shows:
  - Pillar name
  - Owner address
  - Producer/withdraw address
  - Weight (delegation weight)
  - Momentum produced/expected
  - Rank

Optional pagination parameters:
  pageIndex - Page number (default: 0)
  pageSize  - Items per page (default: 25)`,
	Args: cobra.RangeArgs(0, 2),
	RunE: runList,
}

func init() {
	PillarCmd.AddCommand(listCmd)
}

func runList(cmdCobra *cobra.Command, args []string) error {
	// Parse pagination
	pageIndex := uint32(0)
	pageSize := uint32(25)
	if len(args) >= 1 {
		// #nosec G104 - Default value used on parse failure
		fmt.Sscanf(args[0], "%d", &pageIndex)
	}
	if len(args) >= 2 {
		// #nosec G104 - Default value used on parse failure
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

	// Get pillar list
	pillarList, err := rpcClient.PillarApi.GetAll(pageIndex, pageSize)
	if err != nil {
		return fmt.Errorf("failed to get pillar list: %w", err)
	}

	// Display results
	if pillarList.Count == 0 {
		fmt.Println("No pillars found")
		return nil
	}

	fmt.Printf("Total pillars: %d\n", pillarList.Count)
	fmt.Println()

	for idx, pillar := range pillarList.List {
		rank := int(pageIndex)*int(pageSize) + idx + 1

		fmt.Printf("%d. Pillar %s\n", rank, format.Green(pillar.Name))
		fmt.Printf("   Producer: %s\n", pillar.BlockProducingAddress.String())
		fmt.Printf("   Reward: %s\n", pillar.RewardWithdrawAddress.String())
		fmt.Printf("   Weight: %s\n", format.Amount(pillar.Weight, 8))

		if pillar.CurrentStats.ExpectedMomentums > 0 {
			percentage := float64(pillar.CurrentStats.ProducedMomentums) * 100.0 / float64(pillar.CurrentStats.ExpectedMomentums)
			fmt.Printf("   Momentums: %d / %d (%.2f%%)\n",
				pillar.CurrentStats.ProducedMomentums,
				pillar.CurrentStats.ExpectedMomentums,
				percentage)
		}
		fmt.Println()
	}

	return nil
}
