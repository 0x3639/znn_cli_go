package stake

import (
	"fmt"
	"time"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/transaction"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
	"github.com/zenon-network/go-zenon/common/types"
)

// revokeCmd cancels an expired stake entry
var revokeCmd = &cobra.Command{
	Use:   "revoke <id>",
	Short: "Cancel expired stake",
	Long: `Cancel a stake entry by its ID.

Stake entries can only be revoked after they reach their expiration time.
The staked ZNN will be returned to your account.

Use 'stake list' to see stake entry IDs and expiration times.

Example:
  znn-cli stake revoke abc123...

Requires --keyStore flag to specify which wallet to use.`,
	Args: cobra.ExactArgs(1),
	RunE: runRevoke,
}

func init() {
	StakeCmd.AddCommand(revokeCmd)
}

func runRevoke(cmdCobra *cobra.Command, args []string) error {
	cfg, keystoreName, passphrase, index, err := getConfigAndFlags(cmdCobra)
	if err != nil {
		return err
	}

	// Parse stake ID
	idStr := args[0]
	var stakeId types.Hash
	if err := stakeId.UnmarshalText([]byte(idStr)); err != nil {
		return fmt.Errorf("invalid stake ID: %w", err)
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

	// Search through stake entries to find the one with matching ID
	pageIndex := uint32(0)
	found := false
	gotError := false

	for {
		stakeList, err := rpcClient.StakeApi.GetEntriesByAddress(parsedAddress, pageIndex, 25)
		if err != nil {
			return fmt.Errorf("failed to get stake entries: %w", err)
		}

		if len(stakeList.Entries) == 0 {
			break
		}

		// Check if this page contains the stake ID
		for _, entry := range stakeList.Entries {
			if entry.Id == stakeId {
				found = true
				// Check if it can be revoked
				now := time.Now().Unix()
				if entry.ExpirationTimestamp > now {
					expirationTime := time.Unix(entry.ExpirationTimestamp, 0)
					fmt.Printf("%s Stake entry cannot be revoked yet\n", format.Red("Error!"))
					fmt.Printf("Can be revoked at %s (in %s)\n",
						expirationTime.Format("2006-01-02 15:04:05"),
						format.Duration(entry.ExpirationTimestamp-now))
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
		return fmt.Errorf("no stake entry found with ID %s", stakeId.String())
	}

	if gotError {
		return fmt.Errorf("stake entry not ready to revoke")
	}

	// Create revoke template
	template := rpcClient.StakeApi.Cancel(stakeId)

	// Send transaction
	fmt.Println("Revoking stake entry...")
	err = transaction.BuildAndSend(rpcClient.RpcClient, parsedAddress, template, keypair)
	if err != nil {
		return fmt.Errorf("failed to revoke stake: %w", err)
	}

	fmt.Println("Done")
	fmt.Printf("Use %s to collect your ZNN after 2 momentums\n", format.Green("receiveAll"))

	return nil
}
