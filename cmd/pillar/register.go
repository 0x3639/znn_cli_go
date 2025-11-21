package pillar

import (
	"fmt"
	"math/big"
	"regexp"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/transaction"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
	"github.com/zenon-network/go-zenon/common/types"
)

const (
	// PillarRegisterZnnAmount is the ZNN required to register a pillar (15,000 ZNN)
	PillarRegisterZnnAmount = 15000 * 1e8

	// PillarRegisterQsrAmount is the QSR required to register a pillar (150,000 QSR)
	PillarRegisterQsrAmount = 150000 * 1e8
)

var pillarNameRegex = regexp.MustCompile(`^[a-zA-Z0-9]+[\-\.\_\+]{0,1}[a-zA-Z0-9]+$`)

// registerCmd registers a new pillar
var registerCmd = &cobra.Command{
	Use:   "register <name> <producerAddress> <rewardAddress>",
	Short: "Register a new pillar",
	Long: `Register a new pillar in the network.

Requirements:
  - 15,000 ZNN
  - 150,000 QSR
  - Unique pillar name (3-40 characters, alphanumeric with -._+)

The producing address will be used to sign momentums.
The reward address will receive pillar rewards.

Example:
  znn-cli pillar register MyPillar z1qz... z1qz...

Requires --keyStore flag to specify which wallet to use.`,
	Args: cobra.ExactArgs(3),
	RunE: runRegister,
}

func init() {
	PillarCmd.AddCommand(registerCmd)
}

func runRegister(cmdCobra *cobra.Command, args []string) error {
	cfg, keystoreName, passphrase, index, err := getConfigAndFlags(cmdCobra)
	if err != nil {
		return err
	}

	// Parse arguments
	pillarName := args[0]
	producerAddressStr := args[1]
	rewardAddressStr := args[2]

	// Validate pillar name
	if len(pillarName) < 3 || len(pillarName) > 40 {
		return fmt.Errorf("pillar name must be between 3 and 40 characters")
	}
	if !pillarNameRegex.MatchString(pillarName) {
		return fmt.Errorf("invalid pillar name: must be alphanumeric with optional -._+ separators")
	}

	// Parse addresses
	producerAddress, err := types.ParseAddress(producerAddressStr)
	if err != nil {
		return fmt.Errorf("invalid producer address: %w", err)
	}

	rewardAddress, err := types.ParseAddress(rewardAddressStr)
	if err != nil {
		return fmt.Errorf("invalid reward address: %w", err)
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

	// Check if pillar name is already taken
	existingPillar, err := rpcClient.PillarApi.GetByName(pillarName)
	if err == nil && existingPillar != nil {
		return fmt.Errorf("pillar name '%s' is already registered", pillarName)
	}

	// Get account info to check balances
	accountInfo, err := rpcClient.LedgerApi.GetAccountInfoByAddress(parsedAddress)
	if err != nil {
		return fmt.Errorf("failed to get account info: %w", err)
	}

	// Check ZNN balance
	requiredZnn := big.NewInt(PillarRegisterZnnAmount)
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
	requiredQsr := big.NewInt(PillarRegisterQsrAmount)
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
	fmt.Printf("Registering pillar %s\n", format.Green(pillarName))
	fmt.Printf("  Owner: %s\n", address)
	fmt.Printf("  Producer: %s\n", producerAddress.String())
	fmt.Printf("  Reward: %s\n", rewardAddress.String())
	fmt.Printf("  Cost: %s %s + %s %s\n",
		format.Amount(requiredZnn, 8), format.Green("ZNN"),
		format.Amount(requiredQsr, 8), format.Blue("QSR"))
	fmt.Println()

	// Create pillar registration template
	template := rpcClient.PillarApi.Register(pillarName, producerAddress, rewardAddress, uint8(0), uint8(100))

	// Send transaction
	err = transaction.BuildAndSend(rpcClient.RpcClient, parsedAddress, template, keypair)
	if err != nil {
		return fmt.Errorf("failed to register pillar: %w", err)
	}

	fmt.Println("Done")
	fmt.Printf("Pillar %s successfully registered!\n", format.Green(pillarName))

	return nil
}
