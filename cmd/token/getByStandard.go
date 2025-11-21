package token

import (
	"fmt"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/config"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/spf13/cobra"
	"github.com/zenon-network/go-zenon/common/types"
)

// getByStandardCmd gets token info by ZTS address
var getByStandardCmd = &cobra.Command{
	Use:   "getByStandard <tokenStandard>",
	Short: "Get token by ZTS address",
	Long: `Get detailed information about a token by its ZTS address.

The token standard can be:
  - "ZNN" or "znn" for ZNN token
  - "QSR" or "qsr" for QSR token
  - Full ZTS address (e.g., zts1...)

Shows:
  - Token name and symbol
  - Token standard (ZTS address)
  - Total supply and max supply
  - Decimals
  - Owner address
  - Mintability status

Example:
  znn-cli token getByStandard ZNN
  znn-cli token getByStandard zts1...`,
	Args: cobra.ExactArgs(1),
	RunE: runGetByStandard,
}

func init() {
	TokenCmd.AddCommand(getByStandardCmd)
}

func runGetByStandard(cmdCobra *cobra.Command, args []string) error {
	tokenStandardStr := args[0]

	// Parse token standard
	var tokenStandard types.ZenonTokenStandard
	var err error

	// Handle special cases for ZNN and QSR
	switch tokenStandardStr {
	case "ZNN", "znn":
		tokenStandard = types.ZnnTokenStandard
	case "QSR", "qsr":
		tokenStandard = types.QsrTokenStandard
	default:
		tokenStandard, err = types.ParseZTS(tokenStandardStr)
		if err != nil {
			return fmt.Errorf("invalid token standard: %w", err)
		}
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

	// Get token info
	token, err := rpcClient.TokenApi.GetByZts(tokenStandard)
	if err != nil {
		return fmt.Errorf("failed to get token info: %w", err)
	}

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

	// Display token info
	fmt.Printf("Token: %s (%s)\n",
		tokenColor(token.TokenName),
		tokenColor(token.TokenSymbol))
	fmt.Printf("ZTS: %s\n", token.ZenonTokenStandard.String())
	fmt.Printf("Domain: %s\n", token.TokenDomain)
	fmt.Printf("Total Supply: %s\n", format.Amount(token.TotalSupply, int(token.Decimals)))
	fmt.Printf("Max Supply: %s\n", format.Amount(token.MaxSupply, int(token.Decimals)))
	fmt.Printf("Decimals: %d\n", token.Decimals)
	fmt.Printf("Owner: %s\n", token.Owner.String())
	fmt.Printf("Mintable: %v\n", token.IsMintable)
	fmt.Printf("Burnable: %v\n", token.IsBurnable)
	fmt.Printf("Utility: %v\n", token.IsUtility)

	return nil
}
