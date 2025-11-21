package pillar

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

// withdrawQsrCmd withdraws deposited QSR from a pillar
var withdrawQsrCmd = &cobra.Command{
	Use:   "withdrawQsr",
	Short: "Withdraw deposited QSR",
	Long: `Withdraw QSR that was deposited when registering the pillar.

When you revoke a pillar, the deposited QSR becomes available for withdrawal.
This command withdraws the available QSR.

The withdrawn QSR will be sent as a pending transaction that needs to be received.
Use 'receiveAll' after withdrawal to receive the QSR.

Requires --keyStore flag to specify which wallet to use.`,
	RunE: runWithdrawQsr,
}

func init() {
	PillarCmd.AddCommand(withdrawQsrCmd)
}

func runWithdrawQsr(cmdCobra *cobra.Command, args []string) error {
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
	defer rpcClient.Close()

	parsedAddress := types.ParseAddressPanic(address)

	// Get deposit info
	depositInfo, err := rpcClient.PillarApi.GetDepositedQsr(parsedAddress)
	if err != nil {
		return fmt.Errorf("failed to get deposited QSR: %w", err)
	}

	// Check if there is QSR to withdraw
	zero := big.NewInt(0)
	if depositInfo.Cmp(zero) == 0 {
		fmt.Println("No QSR available for withdrawal")
		return nil
	}

	// Display withdrawal info
	fmt.Printf("Withdrawing %s %s\n",
		format.Amount(depositInfo, 8),
		format.Blue("QSR"))

	// Create withdraw template
	template := rpcClient.PillarApi.WithdrawQsr()

	// Send transaction
	err = transaction.BuildAndSend(rpcClient.RpcClient, parsedAddress, template, keypair)
	if err != nil {
		return fmt.Errorf("failed to withdraw QSR: %w", err)
	}

	fmt.Println("Done")
	fmt.Printf("Use %s to receive the QSR\n", format.Green("receiveAll"))

	return nil
}
