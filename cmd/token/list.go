package token

import (
	"fmt"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/config"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/spf13/cobra"
	"github.com/zenon-network/go-zenon/common/types"
)

// listCmd lists all tokens
var listCmd = &cobra.Command{
	Use:   "list [pageIndex pageSize]",
	Short: "List all tokens",
	Long: `List all ZTS tokens in the network.

Shows:
  - Token name and symbol
  - Token standard (ZTS address)
  - Total supply
  - Max supply
  - Decimals
  - Owner address
  - Mintability status

Optional pagination parameters:
  pageIndex - Page number (default: 0)
  pageSize  - Items per page (default: 25)`,
	Args: cobra.RangeArgs(0, 2),
	RunE: runList,
}

func init() {
	TokenCmd.AddCommand(listCmd)
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

	// Get token list
	tokenList, err := rpcClient.TokenApi.GetAll(pageIndex, pageSize)
	if err != nil {
		return fmt.Errorf("failed to get token list: %w", err)
	}

	// Display results
	if tokenList.Count == 0 {
		fmt.Println("No tokens found")
		return nil
	}

	fmt.Printf("Total tokens: %d\n", tokenList.Count)
	fmt.Println()

	for idx, token := range tokenList.List {
		rank := int(pageIndex)*int(pageSize) + idx + 1

		// Determine token color based on standard
		tokenColor := format.Magenta
		if token.ZenonTokenStandard == types.ZnnTokenStandard {
			tokenColor = format.Green
		} else if token.ZenonTokenStandard == types.QsrTokenStandard {
			tokenColor = format.Blue
		}

		fmt.Printf("%d. %s (%s)\n", rank,
			tokenColor(token.TokenName),
			tokenColor(token.TokenSymbol))
		fmt.Printf("   ZTS: %s\n", token.ZenonTokenStandard.String())
		fmt.Printf("   Supply: %s / %s (max)\n",
			format.Amount(token.TotalSupply, int(token.Decimals)),
			format.Amount(token.MaxSupply, int(token.Decimals)))
		fmt.Printf("   Decimals: %d\n", token.Decimals)
		fmt.Printf("   Owner: %s\n", token.Owner.String())
		fmt.Printf("   Mintable: %v\n", token.IsMintable)
		fmt.Println()
	}

	return nil
}
