package token

import (
	"fmt"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/transaction"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
	"github.com/zenon-network/go-zenon/common/types"
)

// burnCmd burns tokens
var burnCmd = &cobra.Command{
	Use:   "burn <tokenStandard> <amount>",
	Short: "Burn tokens",
	Long: `Burn (destroy) tokens from your balance.

Requirements:
  - Token must be burnable
  - Must have sufficient balance

The burned tokens are permanently destroyed and cannot be recovered.

Example:
  znn-cli token burn zts1... 1000

Requires --keyStore flag to specify which wallet to use.`,
	Args: cobra.ExactArgs(2),
	RunE: runBurn,
}

func init() {
	TokenCmd.AddCommand(burnCmd)
}

func runBurn(cmdCobra *cobra.Command, args []string) error {
	cfg, keystoreName, passphrase, index, err := getConfigAndFlags(cmdCobra)
	if err != nil {
		return err
	}

	// Parse arguments
	tokenStandardStr := args[0]
	amountStr := args[1]

	// Parse token standard
	tokenStandard, err := types.ParseZTS(tokenStandardStr)
	if err != nil {
		return fmt.Errorf("invalid token standard: %w", err)
	}

	// Load wallet
	_, keypair, err := wallet.LoadWallet(cfg.Wallet.WalletDir, keystoreName, passphrase, index)
	if err != nil {
		return err
	}

	address, err := wallet.GetAddress(keypair)
	if err != nil {
		return err
	}

	// Connect to node
	rpcClient, err := client.New(cfg.Node.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to node: %w", err)
	}
	defer rpcClient.Close()

	parsedAddress := types.ParseAddressPanic(address)

	// Get token info
	token, err := rpcClient.TokenApi.GetByZts(tokenStandard)
	if err != nil {
		return fmt.Errorf("failed to get token info: %w", err)
	}

	// Verify burnable
	if !token.IsBurnable {
		return fmt.Errorf("token is not burnable")
	}

	// Parse amount with token decimals
	amount, err := format.ParseAmount(amountStr, int(token.Decimals))
	if err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}

	// Check balance
	accountInfo, err := rpcClient.LedgerApi.GetAccountInfoByAddress(parsedAddress)
	if err != nil {
		return fmt.Errorf("failed to get account info: %w", err)
	}

	balance, found := accountInfo.BalanceInfoMap[tokenStandard]
	if !found || balance.Balance.Cmp(amount) < 0 {
		return fmt.Errorf("insufficient balance. You have %s but need %s",
			format.Amount(balance.Balance, int(token.Decimals)),
			format.Amount(amount, int(token.Decimals)))
	}

	// Display burn info
	fmt.Printf("%s Burning %s %s (%s)\n",
		format.Red("Warning!"),
		format.Amount(amount, int(token.Decimals)),
		format.Magenta(token.TokenSymbol),
		token.ZenonTokenStandard.String())
	fmt.Println("This cannot be undone!")
	fmt.Println()

	// Create burn template
	template := rpcClient.TokenApi.Burn(tokenStandard, amount)

	// Send transaction
	err = transaction.BuildAndSend(rpcClient.RpcClient, parsedAddress, template, keypair)
	if err != nil {
		return fmt.Errorf("failed to burn tokens: %w", err)
	}

	fmt.Println("Done")
	fmt.Printf("Successfully burned %s %s\n",
		format.Amount(amount, int(token.Decimals)),
		format.Magenta(token.TokenSymbol))

	return nil
}
