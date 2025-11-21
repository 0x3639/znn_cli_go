package cmd

import (
	"fmt"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
	"github.com/zenon-network/go-zenon/common/types"
)

// unconfirmedCmd displays unconfirmed blocks
var unconfirmedCmd = &cobra.Command{
	Use:   "unconfirmed",
	Short: "List unconfirmed account blocks",
	Long: `List all unconfirmed account blocks for the current address.
These are blocks that have been published but not yet confirmed by a momentum.

Requires --keyStore flag to specify which wallet to use.`,
	RunE: runUnconfirmed,
}

func init() {
	rootCmd.AddCommand(unconfirmedCmd)
}

func runUnconfirmed(cmd *cobra.Command, args []string) error {
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

	// Get unconfirmed blocks
	blocks, err := rpcClient.LedgerApi.GetUnconfirmedBlocksByAddress(types.ParseAddressPanic(address), 0, 50)
	if err != nil {
		return fmt.Errorf("failed to get unconfirmed blocks: %w", err)
	}

	// Display results
	fmt.Printf("Unconfirmed blocks for %s:\n", format.Cyan(address))
	fmt.Println()

	if len(blocks.List) == 0 {
		fmt.Println("No unconfirmed blocks")
		return nil
	}

	fmt.Printf("Found %d unconfirmed block(s):\n", len(blocks.List))
	for i, block := range blocks.List {
		fmt.Printf("\n%d. Hash: %s\n", i+1, format.Cyan(block.Hash.String()))
		fmt.Printf("   Height: %d\n", block.Height)
		fmt.Printf("   Type: %d\n", block.BlockType)
	}

	return nil
}
