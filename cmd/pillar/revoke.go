package pillar

import (
	"fmt"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/transaction"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
	"github.com/zenon-network/go-zenon/common/types"
)

// revokeCmd revokes the pillar owned by the current address
var revokeCmd = &cobra.Command{
	Use:   "revoke <pillarName>",
	Short: "Revoke your pillar",
	Long: `Revoke (dismantle) the pillar owned by the current wallet address.

This will:
  - Return the deposited 15,000 ZNN and 150,000 QSR
  - Remove the pillar from the network
  - Cannot be undone

The pillar name must match a pillar owned by your current address.

Example:
  znn-cli pillar revoke MyPillar

Requires --keyStore flag to specify which wallet to use.`,
	Args: cobra.ExactArgs(1),
	RunE: runRevoke,
}

func init() {
	PillarCmd.AddCommand(revokeCmd)
}

func runRevoke(cmdCobra *cobra.Command, args []string) error {
	cfg, keystoreName, passphrase, index, err := getConfigAndFlags(cmdCobra)
	if err != nil {
		return err
	}

	// Parse pillar name
	pillarName := args[0]

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

	// Verify pillar exists
	_, err = rpcClient.PillarApi.GetByName(pillarName)
	if err != nil {
		return fmt.Errorf("pillar '%s' not found: %w", pillarName, err)
	}

	// Display revoke info
	fmt.Printf("%s Revoking pillar %s\n", format.Red("Warning!"), format.Green(pillarName))
	fmt.Printf("This will return deposited ZNN and QSR\n")
	fmt.Println()

	// Create revoke template
	template := rpcClient.PillarApi.Revoke()

	// Send transaction
	err = transaction.BuildAndSend(rpcClient.RpcClient, parsedAddress, template, keypair)
	if err != nil {
		return fmt.Errorf("failed to revoke pillar: %w", err)
	}

	fmt.Println("Done")
	fmt.Printf("Use %s to receive your ZNN and QSR after 2 momentums\n", format.Green("receiveAll"))

	return nil
}
