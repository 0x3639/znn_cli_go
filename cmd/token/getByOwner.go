package token

import (
	"fmt"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/config"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/spf13/cobra"
	"github.com/zenon-network/go-zenon/common/types"
)

// getByOwnerCmd gets tokens by owner address
var getByOwnerCmd = &cobra.Command{
	Use:   "getByOwner <ownerAddress> [pageIndex pageSize]",
	Short: "Get tokens owned by address",
	Long: `Get all tokens owned by a specific address.

Shows all tokens where the specified address is the owner.

Optional pagination parameters:
  pageIndex - Page number (default: 0)
  pageSize  - Items per page (default: 25)

Example:
  znn-cli token getByOwner z1qz...`,
	Args: cobra.RangeArgs(1, 3),
	RunE: runGetByOwner,
}

func init() {
	TokenCmd.AddCommand(getByOwnerCmd)
}

func runGetByOwner(cmdCobra *cobra.Command, args []string) error {
	ownerAddressStr := args[0]

	// Parse owner address
	ownerAddress, err := types.ParseAddress(ownerAddressStr)
	if err != nil {
		return fmt.Errorf("invalid owner address: %w", err)
	}

	// Parse pagination
	pageIndex := uint32(0)
	pageSize := uint32(25)
	if len(args) >= 2 {
		// Ignore error - default value used on parse failure
		_, _ = fmt.Sscanf(args[1], "%d", &pageIndex)
	}
	if len(args) >= 3 {
		// Ignore error - default value used on parse failure
		_, _ = fmt.Sscanf(args[2], "%d", &pageSize)
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
	defer func() { _ = rpcClient.Close() }()

	// Get tokens by owner
	tokenList, err := rpcClient.TokenApi.GetByOwner(ownerAddress, pageIndex, pageSize)
	if err != nil {
		return fmt.Errorf("failed to get tokens: %w", err)
	}

	// Display results
	if tokenList.Count == 0 {
		fmt.Printf("No tokens found for owner %s\n", ownerAddress.String())
		return nil
	}

	fmt.Printf("Tokens owned by %s: %d\n", ownerAddress.String(), tokenList.Count)
	fmt.Println()

	for idx, token := range tokenList.List {
		rank := int(pageIndex)*int(pageSize) + idx + 1

		// Determine token color
		var tokenColor func(...interface{}) string
		switch token.ZenonTokenStandard {
		case types.ZnnTokenStandard:
			tokenColor = format.Green
		case types.QsrTokenStandard:
			tokenColor = format.Blue
		default:
			tokenColor = format.Magenta
		}

		fmt.Printf("%d. %s (%s)\n", rank,
			tokenColor(token.TokenName),
			tokenColor(token.TokenSymbol))
		fmt.Printf("   ZTS: %s\n", token.ZenonTokenStandard.String())
		fmt.Printf("   Supply: %s / %s (max)\n",
			format.Amount(token.TotalSupply, int(token.Decimals)),
			format.Amount(token.MaxSupply, int(token.Decimals)))
		fmt.Printf("   Mintable: %v\n", token.IsMintable)
		fmt.Println()
	}

	return nil
}
