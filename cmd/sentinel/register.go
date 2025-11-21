package sentinel

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

const (
	// SentinelRegisterZnnAmount is the ZNN required to register a sentinel (5,000 ZNN)
	SentinelRegisterZnnAmount = 5000 * 1e8

	// SentinelRegisterQsrAmount is the QSR required to register a sentinel (50,000 QSR)
	SentinelRegisterQsrAmount = 50000 * 1e8
)

// registerCmd registers a new sentinel
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new sentinel",
	Long: `Register a new sentinel in the network.

Requirements:
  - 5,000 ZNN
  - 50,000 QSR

Sentinels help secure the network by monitoring and validating consensus.

Requires --keyStore flag to specify which wallet to use.`,
	RunE: runRegister,
}

func init() {
	SentinelCmd.AddCommand(registerCmd)
}

func runRegister(cmdCobra *cobra.Command, args []string) error {
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

	// Get account info to check balances
	accountInfo, err := rpcClient.LedgerApi.GetAccountInfoByAddress(parsedAddress)
	if err != nil {
		return fmt.Errorf("failed to get account info: %w", err)
	}

	// Check ZNN balance
	requiredZnn := big.NewInt(SentinelRegisterZnnAmount)
	znnBalance, found := accountInfo.BalanceInfoMap[types.ZnnTokenStandard]
	if !found || znnBalance.Balance.Cmp(requiredZnn) < 0 {
		currentBalance := big.NewInt(0)
		if found {
			currentBalance = znnBalance.Balance
		}
		return fmt.Errorf("insufficient ZNN balance. You have %s but need %s",
			format.Amount(currentBalance, 8),
			format.Amount(requiredZnn, 8))
	}

	// Check QSR balance
	requiredQsr := big.NewInt(SentinelRegisterQsrAmount)
	qsrBalance, found := accountInfo.BalanceInfoMap[types.QsrTokenStandard]
	if !found || qsrBalance.Balance.Cmp(requiredQsr) < 0 {
		currentBalance := big.NewInt(0)
		if found {
			currentBalance = qsrBalance.Balance
		}
		return fmt.Errorf("insufficient QSR balance. You have %s but need %s",
			format.Amount(currentBalance, 8),
			format.Amount(requiredQsr, 8))
	}

	// Display registration info
	fmt.Println("Registering sentinel")
	fmt.Printf("  Owner: %s\n", address)
	fmt.Printf("  Cost: %s %s + %s %s\n",
		format.Amount(requiredZnn, 8), format.Green("ZNN"),
		format.Amount(requiredQsr, 8), format.Blue("QSR"))
	fmt.Println()

	// Create sentinel registration template
	template := rpcClient.SentinelApi.Register()

	// Send transaction
	err = transaction.BuildAndSend(rpcClient.RpcClient, parsedAddress, template, keypair)
	if err != nil {
		return fmt.Errorf("failed to register sentinel: %w", err)
	}

	fmt.Println("Done")
	fmt.Println("Sentinel successfully registered!")

	return nil
}
