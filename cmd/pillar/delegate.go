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

// delegateCmd delegates to a pillar
var delegateCmd = &cobra.Command{
	Use:   "delegate <pillarName>",
	Short: "Delegate to a pillar",
	Long: `Delegate your ZNN weight to a pillar.

Delegation increases the pillar's weight in the network, which affects:
  - Momentum production frequency
  - Reward distribution

You can only delegate to one pillar at a time.
Your ZNN remains in your wallet while delegated.

Example:
  znn-cli pillar delegate MyPillar

Requires --keyStore flag to specify which wallet to use.`,
	Args: cobra.ExactArgs(1),
	RunE: runDelegate,
}

func init() {
	PillarCmd.AddCommand(delegateCmd)
}

func runDelegate(cmdCobra *cobra.Command, args []string) error {
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
	defer func() { _ = rpcClient.Close() }()

	parsedAddress := types.ParseAddressPanic(address)

	// Verify pillar exists
	pillar, err := rpcClient.PillarApi.GetByName(pillarName)
	if err != nil {
		return fmt.Errorf("pillar '%s' not found: %w", pillarName, err)
	}

	// Create delegate template
	template := rpcClient.PillarApi.Delegate(pillar.Name)

	// Send transaction
	fmt.Printf("Delegating to pillar %s\n", format.Green(pillarName))

	err = transaction.BuildAndSend(rpcClient.RpcClient, parsedAddress, template, keypair)
	if err != nil {
		return fmt.Errorf("failed to delegate: %w", err)
	}

	fmt.Println("Done")
	fmt.Printf("Successfully delegated to %s\n", format.Green(pillarName))

	return nil
}
