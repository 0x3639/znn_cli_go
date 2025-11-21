package cmd

import (
	"fmt"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
	"github.com/zenon-network/go-zenon/common/types"
)

// balanceCmd displays account balance
var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Display account balance",
	Long: `Display the balance for all tokens (ZNN, QSR, and custom ZTS tokens)
for the currently selected wallet address.

Requires --keyStore flag to specify which wallet to use.`,
	RunE: runBalance,
}

func init() {
	rootCmd.AddCommand(balanceCmd)
}

func runBalance(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()
	keystoreName := GetKeyStore()
	passphrase := GetPassphrase()
	index := GetIndex()

	// Load wallet to get address
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

	// Get account info
	accountInfo, err := rpcClient.LedgerApi.GetAccountInfoByAddress(types.ParseAddressPanic(address))
	if err != nil {
		return fmt.Errorf("failed to get account info: %w", err)
	}

	// Display address
	fmt.Printf("Address: %s\n", format.Cyan(address))
	fmt.Println()

	// Display balances
	if len(accountInfo.BalanceInfoMap) == 0 {
		fmt.Println("No balances found")
		return nil
	}

	fmt.Println("Balances:")
	for tokenStandard, balanceInfo := range accountInfo.BalanceInfoMap {
		symbol := balanceInfo.TokenInfo.TokenSymbol
		decimals := int(balanceInfo.TokenInfo.Decimals)
		formattedAmount := format.Amount(balanceInfo.Balance, decimals)

		// Color code based on token
		switch symbol {
		case "ZNN":
			fmt.Printf("  %s %s\n", formattedAmount, format.Green(symbol))
		case "QSR":
			fmt.Printf("  %s %s\n", formattedAmount, format.Blue(symbol))
		default:
			fmt.Printf("  %s %s (%s)\n", formattedAmount, format.Magenta(symbol), tokenStandard)
		}
	}

	return nil
}
