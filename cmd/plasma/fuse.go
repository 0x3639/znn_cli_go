package plasma

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
	// MinFuseAmount is the minimum QSR that can be fused (10 QSR)
	MinFuseAmount = 10 * 1e8
)

// fuseCmd fuses QSR for plasma
var fuseCmd = &cobra.Command{
	Use:   "fuse <beneficiaryAddress> <amount>",
	Short: "Fuse QSR for beneficiary",
	Long: `Fuse QSR tokens to generate plasma for a beneficiary address.

The beneficiary can be your own address or another address.
Fused QSR can be unfused after 10 momentums (~1 hour).

Requirements:
  - Minimum amount: 10 QSR
  - Amount must be a whole number (no decimals)
  - Sufficient QSR balance

Example:
  znn-cli plasma fuse z1qz... 50

Requires --keyStore flag to specify which wallet to use.`,
	Args: cobra.ExactArgs(2),
	RunE: runFuse,
}

func init() {
	PlasmaCmd.AddCommand(fuseCmd)
}

func runFuse(cmdCobra *cobra.Command, args []string) error {
	cfg, keystoreName, passphrase, index, err := getConfigAndFlags(cmdCobra)
	if err != nil {
		return err
	}

	// Parse arguments
	beneficiaryStr := args[0]
	amountStr := args[1]

	// Parse beneficiary address
	beneficiary, err := types.ParseAddress(beneficiaryStr)
	if err != nil {
		return fmt.Errorf("invalid beneficiary address: %w", err)
	}

	// Parse amount (QSR has 8 decimals)
	amount, err := format.ParseAmount(amountStr, 8)
	if err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}

	// Validate amount
	if amount.Cmp(big.NewInt(MinFuseAmount)) < 0 {
		return fmt.Errorf("invalid amount: %s QSR. Minimum fuse amount is %s",
			format.Amount(amount, 8),
			format.Amount(big.NewInt(MinFuseAmount), 8))
	}

	// Check if amount is a whole number (no decimals)
	oneQsr := big.NewInt(1e8)
	remainder := new(big.Int).Mod(amount, oneQsr)
	if remainder.Cmp(big.NewInt(0)) != 0 {
		return fmt.Errorf("amount must be a whole number (no decimals)")
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

	// Get account info to check QSR balance
	accountInfo, err := rpcClient.LedgerApi.GetAccountInfoByAddress(types.ParseAddressPanic(address))
	if err != nil {
		return fmt.Errorf("failed to get account info: %w", err)
	}

	// Check QSR balance
	qsrBalance, found := accountInfo.BalanceInfoMap[types.QsrTokenStandard]
	if !found || qsrBalance.Balance.Cmp(amount) < 0 {
		currentBalance := big.NewInt(0)
		if found {
			currentBalance = qsrBalance.Balance
		}
		return fmt.Errorf("insufficient QSR balance. You have %s but need %s",
			format.Amount(currentBalance, 8),
			format.Amount(amount, 8))
	}

	// Create fuse template
	template := rpcClient.PlasmaApi.Fuse(beneficiary, amount)

	// Send transaction
	fmt.Printf("Fusing %s %s to %s\n",
		format.Amount(amount, 8),
		format.Blue("QSR"),
		beneficiary.String())

	err = transaction.BuildAndSend(rpcClient.RpcClient, types.ParseAddressPanic(address), template, keypair)
	if err != nil {
		return fmt.Errorf("failed to fuse: %w", err)
	}

	fmt.Println("Done")
	fmt.Println("Plasma will be available after 1 momentum")

	return nil
}
