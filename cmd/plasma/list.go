package plasma

import (
	"fmt"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
	"github.com/zenon-network/go-zenon/common/types"
)

// listCmd lists fusion entries for the current address
var listCmd = &cobra.Command{
	Use:   "list [pageIndex pageSize]",
	Short: "List fusion entries",
	Long: `List all plasma fusion entries for the current wallet address.

Shows:
  - QSR amount fused
  - Beneficiary address
  - Expiration height (when it can be canceled)
  - Fusion entry ID (for cancellation)

Optional pagination parameters:
  pageIndex - Page number (default: 0)
  pageSize  - Items per page (default: 25)

Requires --keyStore flag to specify which wallet to use.`,
	Args: cobra.RangeArgs(0, 2),
	RunE: runList,
}

func init() {
	PlasmaCmd.AddCommand(listCmd)
}

func runList(cmdCobra *cobra.Command, args []string) error {
	cfg, keystoreName, passphrase, index, err := getConfigAndFlags(cmdCobra)
	if err != nil {
		return err
	}

	// Parse pagination
	pageIndex := uint32(0)
	pageSize := uint32(25)
	if len(args) >= 1 {
		// #nosec G104 - Default value used on parse failure
		fmt.Sscanf(args[0], "%d", &pageIndex)
	}
	if len(args) >= 2 {
		// #nosec G104 - Default value used on parse failure
		fmt.Sscanf(args[1], "%d", &pageSize)
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
	defer rpcClient.Close()

	// Get fusion entries
	fusionList, err := rpcClient.PlasmaApi.GetEntriesByAddress(types.ParseAddressPanic(address), pageIndex, pageSize)
	if err != nil {
		return fmt.Errorf("failed to get fusion entries: %w", err)
	}

	// Display results
	if fusionList.Count == 0 {
		fmt.Println("No Plasma fusion entries found")
		return nil
	}

	fmt.Printf("Fusing %s %s for Plasma in %d entries\n",
		format.Amount(fusionList.QsrAmount, 8),
		format.Blue("QSR"),
		fusionList.Count)
	fmt.Println()

	for _, entry := range fusionList.Fusions {
		fmt.Printf("  %s %s for %s\n",
			format.Amount(entry.QsrAmount, 8),
			format.Blue("QSR"),
			entry.Beneficiary.String())
		fmt.Printf("  Can be canceled at momentum height: %d. Use id %s to cancel\n",
			entry.ExpirationHeight,
			format.Cyan(entry.Id.String()))
		fmt.Println()
	}

	return nil
}
