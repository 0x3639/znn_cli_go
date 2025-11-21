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

// receiveAllCmd receives all pending transactions
var receiveAllCmd = &cobra.Command{
	Use:   "receiveAll",
	Short: "Receive all pending transactions",
	Long: `Receive all pending (unreceived) transactions in batches.

This command will:
  1. Query all unreceived blocks
  2. Receive them in batches of 5
  3. Continue until all blocks are received

Requires --keyStore flag to specify which wallet to use.`,
	RunE: runReceiveAll,
}

func init() {
	rootCmd.AddCommand(receiveAllCmd)
}

func runReceiveAll(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()
	keystoreName := GetKeyStore()
	passphrase := GetPassphrase()
	index := GetIndex()

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

	parsedAddress := types.ParseAddressPanic(address)

	// Get initial unreceived blocks
	blocks, err := rpcClient.LedgerApi.GetUnreceivedBlocksByAddress(parsedAddress, 0, 5)
	if err != nil {
		return fmt.Errorf("failed to get unreceived blocks: %w", err)
	}

	if len(blocks.List) == 0 {
		fmt.Println("Nothing to receive")
		return nil
	}

	// Show how many transactions need to be received
	if blocks.More {
		fmt.Printf("You have %s than %s transaction(s) to receive\n",
			format.Red("more"),
			format.Green(fmt.Sprintf("%d", len(blocks.List))))
	} else {
		fmt.Printf("You have %s transaction(s) to receive\n",
			format.Green(fmt.Sprintf("%d", len(blocks.List))))
	}

	fmt.Println("Receiving transactions...")

	// Receive all blocks in batches
	receivedCount := 0
	for len(blocks.List) > 0 {
		// Receive each block in current batch
		for _, block := range blocks.List {
			template := &nom.AccountBlock{
				Version:         1,
				ChainIdentifier: 1,
				BlockType:       nom.BlockTypeUserReceive,
				FromBlockHash:   block.Hash,
				Data:            nil,
			}

			err = transaction.BuildAndSend(rpcClient.RpcClient, parsedAddress, template, keypair)
			if err != nil {
				return fmt.Errorf("failed to receive block %s: %w", block.Hash, err)
			}

			receivedCount++
			if cfg.Display.Verbose {
				fmt.Printf("  Received %s\n", format.Cyan(block.Hash.String()))
			}
		}

		// Get next batch
		blocks, err = rpcClient.LedgerApi.GetUnreceivedBlocksByAddress(parsedAddress, 0, 5)
		if err != nil {
			return fmt.Errorf("failed to get unreceived blocks: %w", err)
		}
	}

	fmt.Printf("Successfully received %s transaction(s)\n", format.Green(fmt.Sprintf("%d", receivedCount)))

	return nil
}
