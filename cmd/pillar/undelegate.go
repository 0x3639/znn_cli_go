package pillar

import (
	"fmt"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/transaction"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
	"github.com/zenon-network/go-zenon/common/types"
)

// undelegateCmd removes delegation
var undelegateCmd = &cobra.Command{
	Use:   "undelegate",
	Short: "Remove delegation",
	Long: `Remove your delegation from any pillar.

This will stop delegating your ZNN weight to the pillar.
Your ZNN remains in your wallet.

Requires --keyStore flag to specify which wallet to use.`,
	RunE: runUndelegate,
}

func init() {
	PillarCmd.AddCommand(undelegateCmd)
}

func runUndelegate(cmdCobra *cobra.Command, args []string) error {
	cfg, keystoreName, passphrase, index, err := getConfigAndFlags(cmdCobra)
	if err != nil {
		return err
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

	parsedAddress := types.ParseAddressPanic(address)

	// Create undelegate template
	template := rpcClient.PillarApi.Undelegate()

	// Send transaction
	fmt.Println("Removing delegation")

	err = transaction.BuildAndSend(rpcClient.RpcClient, parsedAddress, template, keypair)
	if err != nil {
		return fmt.Errorf("failed to undelegate: %w", err)
	}

	fmt.Println("Done")
	fmt.Println("Successfully removed delegation")

	return nil
}
