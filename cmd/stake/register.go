package stake

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/transaction"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
	"github.com/zenon-network/go-zenon/common/types"
)

const (
	// MinStakeAmount is the minimum ZNN that can be staked (1 ZNN)
	MinStakeAmount = 1 * 1e8

	// StakeTimeUnit is one month in seconds (30 days)
	StakeTimeUnit = 30 * 24 * 60 * 60

	// MinStakeMonths is the minimum staking duration
	MinStakeMonths = 1

	// MaxStakeMonths is the maximum staking duration
	MaxStakeMonths = 12
)

// registerCmd stakes ZNN for rewards
var registerCmd = &cobra.Command{
	Use:   "register <amount> <duration>",
	Short: "Stake ZNN for rewards",
	Long: `Stake ZNN tokens to earn rewards.

The amount is locked for the specified duration (in months).
After expiration, you can revoke the stake to get your ZNN back.

Requirements:
  - Minimum amount: 1 ZNN
  - Duration: 1-12 months
  - Sufficient ZNN balance

Example:
  znn-cli stake register 100 3    # Stake 100 ZNN for 3 months

Requires --keyStore flag to specify which wallet to use.`,
	Args: cobra.ExactArgs(2),
	RunE: runRegister,
}

func init() {
	StakeCmd.AddCommand(registerCmd)
}

func runRegister(cmdCobra *cobra.Command, args []string) error {
	cfg, keystoreName, passphrase, index, err := getConfigAndFlags(cmdCobra)
	if err != nil {
		return err
	}

	// Parse arguments
	amountStr := args[0]
	durationStr := args[1]

	// Parse amount (ZNN has 8 decimals)
	amount, err := format.ParseAmount(amountStr, 8)
	if err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}

	// Parse duration (in months)
	duration, err := strconv.ParseInt(durationStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid duration: must be a number between %d and %d", MinStakeMonths, MaxStakeMonths)
	}

	// Validate amount
	if amount.Cmp(big.NewInt(MinStakeAmount)) < 0 {
		return fmt.Errorf("invalid amount: minimum stake amount is %s ZNN",
			format.Amount(big.NewInt(MinStakeAmount), 8))
	}

	// Validate duration
	if duration < MinStakeMonths || duration > MaxStakeMonths {
		return fmt.Errorf("invalid duration: must be between %d and %d months", MinStakeMonths, MaxStakeMonths)
	}

	// Calculate duration in seconds
	durationSeconds := duration * StakeTimeUnit

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

	// Get account info to check ZNN balance
	accountInfo, err := rpcClient.LedgerApi.GetAccountInfoByAddress(types.ParseAddressPanic(address))
	if err != nil {
		return fmt.Errorf("failed to get account info: %w", err)
	}

	// Check ZNN balance
	znnBalance, found := accountInfo.BalanceInfoMap[types.ZnnTokenStandard]
	if !found || znnBalance.Balance.Cmp(amount) < 0 {
		currentBalance := big.NewInt(0)
		if found {
			currentBalance = znnBalance.Balance
		}
		return fmt.Errorf("insufficient ZNN balance. You have %s but need %s",
			format.Amount(currentBalance, 8),
			format.Amount(amount, 8))
	}

	// Create stake template
	template := rpcClient.StakeApi.Stake(durationSeconds, amount)

	// Send transaction
	fmt.Printf("Staking %s %s for %d month(s)\n",
		format.Amount(amount, 8),
		format.Green("ZNN"),
		duration)

	err = transaction.BuildAndSend(rpcClient.RpcClient, types.ParseAddressPanic(address), template, keypair)
	if err != nil {
		return fmt.Errorf("failed to stake: %w", err)
	}

	fmt.Println("Done")
	fmt.Printf("Use %s to see your stake entries\n", format.Green("stake list"))

	return nil
}
