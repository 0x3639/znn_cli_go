package cmd

import (
	"fmt"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
	"github.com/zenon-network/go-zenon/common/types"
)

// unreceivedCmd displays unreceived blocks
var unreceivedCmd = &cobra.Command{
	Use:   "unreceived",
	Short: "List unreceived (pending) transactions",
	Long: `List all unreceived/pending transactions for the current address.
These are incoming transactions that need to be received.

Requires --keyStore flag to specify which wallet to use.`,
	RunE: runUnreceived,
}

func init() {
	rootCmd.AddCommand(unreceivedCmd)
}

func runUnreceived(cmd *cobra.Command, args []string) error {
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
	defer func() { _ = rpcClient.Close() }()

	// Get unreceived blocks
	blocks, err := rpcClient.LedgerApi.GetUnreceivedBlocksByAddress(types.ParseAddressPanic(address), 0, 50)
	if err != nil {
		return fmt.Errorf("failed to get unreceived blocks: %w", err)
	}

	// Display results
	fmt.Printf("Unreceived blocks for %s:\n", format.Cyan(address))
	fmt.Println()

	if len(blocks.List) == 0 {
		fmt.Println("No unreceived blocks")
		return nil
	}

	fmt.Printf("Found %d unreceived block(s):\n", len(blocks.List))
	for i, block := range blocks.List {
		symbol := block.TokenInfo.TokenSymbol
		decimals := int(block.TokenInfo.Decimals)
		amount := format.Amount(block.Amount, decimals)

		fmt.Printf("\n%d. Hash: %s\n", i+1, format.Cyan(block.Hash.String()))
		fmt.Printf("   From: %s\n", format.Cyan(block.Address.String()))
		fmt.Printf("   Amount: %s %s\n", amount, format.FormatToken(block.Amount, decimals, symbol))
	}

	return nil
}
