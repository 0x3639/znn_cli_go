package sentinel

import (
	"fmt"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/transaction"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
	"github.com/zenon-network/go-zenon/common/types"
)

// revokeCmd revokes the sentinel owned by the current address
var revokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke your sentinel",
	Long: `Revoke (dismantle) the sentinel owned by the current wallet address.

This will:
  - Return the deposited 5,000 ZNN and 50,000 QSR
  - Remove the sentinel from the network
  - Cannot be undone

Requires --keyStore flag to specify which wallet to use.`,
	RunE: runRevoke,
}

func init() {
	SentinelCmd.AddCommand(revokeCmd)
}

func runRevoke(cmdCobra *cobra.Command, args []string) error {
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

	// Verify sentinel exists for this address
	_, err = rpcClient.SentinelApi.GetByOwner(parsedAddress)
	if err != nil {
		return fmt.Errorf("no sentinel found for address %s", address)
	}

	// Display revoke info
	fmt.Printf("%s Revoking sentinel\n", format.Red("Warning!"))
	fmt.Printf("This will return deposited ZNN and QSR\n")
	fmt.Println()

	// Create revoke template
	template := rpcClient.SentinelApi.Revoke()

	// Send transaction
	err = transaction.BuildAndSend(rpcClient.RpcClient, parsedAddress, template, keypair)
	if err != nil {
		return fmt.Errorf("failed to revoke sentinel: %w", err)
	}

	fmt.Println("Done")
	fmt.Printf("Use %s to receive your ZNN and QSR after 2 momentums\n", format.Green("receiveAll"))

	return nil
}
