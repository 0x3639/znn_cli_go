package cmd

import (
	"fmt"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/transaction"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
	"github.com/zenon-network/go-zenon/chain/nom"
	"github.com/zenon-network/go-zenon/common/types"
)

// receiveCmd receives a specific pending transaction by hash
var receiveCmd = &cobra.Command{
	Use:   "receive <blockHash>",
	Short: "Receive a specific pending transaction",
	Long: `Receive a specific pending (unreceived) transaction by its block hash.

Use the 'unreceived' command to list pending transactions and their hashes.

Example:
  znn-cli receive abc123...

Requires --keyStore flag to specify which wallet to use.`,
	Args: cobra.ExactArgs(1),
	RunE: runReceive,
}

func init() {
	rootCmd.AddCommand(receiveCmd)
}

func runReceive(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()
	keystoreName := GetKeyStore()
	passphrase := GetPassphrase()
	index := GetIndex()

	// Parse block hash
	blockHashStr := args[0]
	var blockHash types.Hash
	err := blockHash.UnmarshalText([]byte(blockHashStr))
	if err != nil {
		return fmt.Errorf("invalid block hash: %w", err)
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
	defer func() { _ = rpcClient.Close() }()

	// Create receive template
	template := &nom.AccountBlock{
		Version:         1,
		ChainIdentifier: 1,
		BlockType:       nom.BlockTypeUserReceive,
		FromBlockHash:   blockHash,
		Data:            nil,
	}

	// Receive transaction
	fmt.Println("Receiving transaction...")
	err = transaction.BuildAndSend(rpcClient.RpcClient, types.ParseAddressPanic(address), template, keypair)
	if err != nil {
		return fmt.Errorf("failed to receive transaction: %w", err)
	}

	// Display success
	fmt.Printf("Successfully received transaction %s\n", format.Cyan(blockHash.String()))

	return nil
}
