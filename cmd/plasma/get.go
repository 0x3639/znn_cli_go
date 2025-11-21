package plasma

import (
	"fmt"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
	"github.com/zenon-network/go-zenon/common/types"
)

// getCmd displays plasma info for the current address
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get plasma info for current address",
	Long: `Display current plasma and max plasma for the current wallet address.

Shows:
  - Current plasma available
  - Maximum plasma capacity
  - Amount of QSR fused

Requires --keyStore flag to specify which wallet to use.`,
	RunE: runGet,
}

func init() {
	PlasmaCmd.AddCommand(getCmd)
}

func runGet(cmdCobra *cobra.Command, args []string) error {
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

	// Get plasma info
	plasmaInfo, err := rpcClient.PlasmaApi.Get(types.ParseAddressPanic(address))
	if err != nil {
		return fmt.Errorf("failed to get plasma info: %w", err)
	}

	// Display plasma info
	fmt.Printf("%s has %s / %s plasma with %s %s fused\n",
		format.Green(address),
		format.Green(fmt.Sprintf("%d", plasmaInfo.CurrentPlasma)),
		fmt.Sprintf("%d", plasmaInfo.MaxPlasma),
		format.Amount(plasmaInfo.QsrAmount, 8), // QSR has 8 decimals
		format.Blue("QSR"))

	return nil
}
