package stake

import (
	"fmt"
	"time"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
	"github.com/zenon-network/go-zenon/common/types"
)

// listCmd lists stake entries for the current address
var listCmd = &cobra.Command{
	Use:   "list [pageIndex pageSize]",
	Short: "List stake entries",
	Long: `List all staking entries for the current wallet address.

Shows:
  - ZNN amount staked
  - Start time (when stake was created)
  - Expiration time (when it can be revoked)
  - Duration (months staked)
  - Stake entry ID (for revocation)

Optional pagination parameters:
  pageIndex - Page number (default: 0)
  pageSize  - Items per page (default: 25)

Requires --keyStore flag to specify which wallet to use.`,
	Args: cobra.RangeArgs(0, 2),
	RunE: runList,
}

func init() {
	StakeCmd.AddCommand(listCmd)
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
		// Ignore error - default value used on parse failure
		_, _ = fmt.Sscanf(args[0], "%d", &pageIndex)
	}
	if len(args) >= 2 {
		// Ignore error - default value used on parse failure
		_, _ = fmt.Sscanf(args[1], "%d", &pageSize)
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

	// Get stake entries
	stakeList, err := rpcClient.StakeApi.GetEntriesByAddress(types.ParseAddressPanic(address), pageIndex, pageSize)
	if err != nil {
		return fmt.Errorf("failed to get stake entries: %w", err)
	}

	// Display results
	if stakeList.Count == 0 {
		fmt.Println("No stake entries found")
		return nil
	}

	fmt.Printf("Staking %s %s in %d entries\n",
		format.Amount(stakeList.TotalAmount, 8),
		format.Green("ZNN"),
		stakeList.Count)
	fmt.Println()

	for _, entry := range stakeList.Entries {
		startTime := time.Unix(entry.StartTimestamp, 0)
		expirationTime := time.Unix(entry.ExpirationTimestamp, 0)
		durationMonths := (entry.ExpirationTimestamp - entry.StartTimestamp) / (30 * 24 * 60 * 60)

		fmt.Printf("  %s %s for %d month(s)\n",
			format.Amount(entry.Amount, 8),
			format.Green("ZNN"),
			durationMonths)
		fmt.Printf("  Started at %s, can revoke at %s\n",
			startTime.Format("2006-01-02 15:04:05"),
			expirationTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("  ID %s\n", format.Cyan(entry.Id.String()))
		fmt.Println()
	}

	return nil
}
