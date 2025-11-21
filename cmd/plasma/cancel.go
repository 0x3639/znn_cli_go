package plasma

import (
	"fmt"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/transaction"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
	"github.com/zenon-network/go-zenon/common/types"
)

// cancelCmd cancels a fusion entry by ID
var cancelCmd = &cobra.Command{
	Use:   "cancel <id>",
	Short: "Cancel fusion by ID",
	Long: `Cancel a plasma fusion entry by its ID.

Fusion entries can only be canceled after they reach their expiration height
(10 momentums after creation, approximately 1 hour).

Use 'plasma list' to see fusion entry IDs and expiration heights.

Example:
  znn-cli plasma cancel abc123...

Requires --keyStore flag to specify which wallet to use.`,
	Args: cobra.ExactArgs(1),
	RunE: runCancel,
}

func init() {
	PlasmaCmd.AddCommand(cancelCmd)
}

func runCancel(cmdCobra *cobra.Command, args []string) error {
	cfg, keystoreName, passphrase, index, err := getConfigAndFlags(cmdCobra)
	if err != nil {
		return err
	}

	// Parse fusion ID
	idStr := args[0]
	var fusionId types.Hash
	if err := fusionId.UnmarshalText([]byte(idStr)); err != nil {
		return fmt.Errorf("invalid fusion ID: %w", err)
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

	parsedAddress := types.ParseAddressPanic(address)

	// Get current momentum height
	momentum, err := rpcClient.LedgerApi.GetFrontierMomentum()
	if err != nil {
		return fmt.Errorf("failed to get frontier momentum: %w", err)
	}

	// Search through fusion entries to find the one with matching ID
	pageIndex := uint32(0)
	found := false
	gotError := false

	for {
		fusionList, err := rpcClient.PlasmaApi.GetEntriesByAddress(parsedAddress, pageIndex, 25)
		if err != nil {
			return fmt.Errorf("failed to get fusion entries: %w", err)
		}

		if len(fusionList.Fusions) == 0 {
			break
		}

		// Check if this page contains the fusion ID
		for _, entry := range fusionList.Fusions {
			if entry.Id == fusionId {
				found = true
				// Check if it can be canceled
				if entry.ExpirationHeight > momentum.Height {
					fmt.Printf("%s Fuse entry cannot be cancelled yet\n", format.Red("Error!"))
					fmt.Printf("Can be canceled at momentum height %d (current: %d)\n",
						entry.ExpirationHeight, momentum.Height)
					gotError = true
				}
				break
			}
		}

		if found {
			break
		}

		pageIndex++
	}

	if !found {
		return fmt.Errorf("no fusion entry found with ID %s", fusionId.String())
	}

	if gotError {
		return fmt.Errorf("fusion entry not ready to cancel")
	}

	// Create cancel template
	template := rpcClient.PlasmaApi.Cancel(fusionId)

	// Send transaction
	fmt.Println("Canceling fusion entry...")
	err = transaction.BuildAndSend(rpcClient.RpcClient, parsedAddress, template, keypair)
	if err != nil {
		return fmt.Errorf("failed to cancel fusion: %w", err)
	}

	fmt.Println("Done")
	fmt.Printf("Use %s to collect your QSR after 2 momentums\n", format.Green("receiveAll"))

	return nil
}
