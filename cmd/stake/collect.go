package stake

import (
	"fmt"
	"math/big"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/transaction"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
	"github.com/zenon-network/go-zenon/common/types"
)

// collectCmd collects staking rewards
var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "Collect staking rewards",
	Long: `Collect accumulated staking rewards.

Rewards are automatically accumulated while your ZNN is staked.
Use this command to claim your rewards.

The rewards will be sent as pending transactions that need to be received.
Use 'receiveAll' after collection to receive the rewards.

Requires --keyStore flag to specify which wallet to use.`,
	RunE: runCollect,
}

func init() {
	StakeCmd.AddCommand(collectCmd)
}

func runCollect(cmdCobra *cobra.Command, args []string) error {
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

	// Get uncollected rewards
	rewardInfo, err := rpcClient.StakeApi.GetUncollectedReward(parsedAddress)
	if err != nil {
		return fmt.Errorf("failed to get uncollected rewards: %w", err)
	}

	// Check if there are rewards to collect
	zero := big.NewInt(0)
	if rewardInfo.Znn.Cmp(zero) == 0 && rewardInfo.Qsr.Cmp(zero) == 0 {
		fmt.Println("Nothing to collect")
		return nil
	}

	// Display rewards
	fmt.Println("Collecting rewards:")
	if rewardInfo.Znn.Cmp(zero) > 0 {
		fmt.Printf("  %s %s\n",
			format.Amount(rewardInfo.Znn, 8),
			format.Green("ZNN"))
	}
	if rewardInfo.Qsr.Cmp(zero) > 0 {
		fmt.Printf("  %s %s\n",
			format.Amount(rewardInfo.Qsr, 8),
			format.Blue("QSR"))
	}
	fmt.Println()

	// Create collect template
	template := rpcClient.StakeApi.CollectReward()

	// Send transaction
	err = transaction.BuildAndSend(rpcClient.RpcClient, parsedAddress, template, keypair)
	if err != nil {
		return fmt.Errorf("failed to collect rewards: %w", err)
	}

	fmt.Println("Done")
	fmt.Printf("Use %s to receive the rewards\n", format.Green("receiveAll"))

	return nil
}
