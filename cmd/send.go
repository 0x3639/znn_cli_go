package cmd

import (
	"fmt"
	"math/big"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/transaction"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
	"github.com/zenon-network/go-zenon/chain/nom"
	"github.com/zenon-network/go-zenon/common/types"
)

// sendCmd sends tokens to an address
var sendCmd = &cobra.Command{
	Use:   "send <toAddress> <amount> <token>",
	Short: "Send tokens to an address",
	Long: `Send ZNN, QSR, or custom ZTS tokens to a destination address.

Examples:
  znn-cli send z1qz... 10.5 ZNN
  znn-cli send z1qz... 100 QSR
  znn-cli send z1qz... 5.25 zts1...

Token can be:
  - ZNN (Zenon coin)
  - QSR (Quasar coin)
  - zts1... (Custom ZTS token standard)

Requires --keyStore flag to specify which wallet to use.`,
	Args: cobra.ExactArgs(3),
	RunE: runSend,
}

func init() {
	rootCmd.AddCommand(sendCmd)
}

func runSend(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()
	keystoreName := GetKeyStore()
	passphrase := GetPassphrase()
	index := GetIndex()

	// Parse arguments
	toAddressStr := args[0]
	amountStr := args[1]
	tokenStr := args[2]

	// Parse destination address
	toAddress, err := types.ParseAddress(toAddressStr)
	if err != nil {
		return fmt.Errorf("invalid destination address: %w", err)
	}

	// Parse token standard
	var tokenStandard types.ZenonTokenStandard
	switch tokenStr {
	case "ZNN", "znn":
		tokenStandard = types.ZnnTokenStandard
	case "QSR", "qsr":
		tokenStandard = types.QsrTokenStandard
	default:
		// Try to parse as ZTS
		ts, err := types.ParseZTS(tokenStr)
		if err != nil {
			return fmt.Errorf("invalid token standard (use ZNN/QSR or zts1...): %w", err)
		}
		tokenStandard = ts
	}

	// Load wallet to get address
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

	// Get account info to check balance and get token decimals
	accountInfo, err := rpcClient.LedgerApi.GetAccountInfoByAddress(types.ParseAddressPanic(address))
	if err != nil {
		return fmt.Errorf("failed to get account info: %w", err)
	}

	// Find the token balance and decimals
	var decimals int
	var balance *big.Int
	found := false
	for tokenStd, balanceInfo := range accountInfo.BalanceInfoMap {
		if tokenStd == tokenStandard {
			decimals = int(balanceInfo.TokenInfo.Decimals)
			balance = balanceInfo.Balance
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("you have no balance for token %s", tokenStandard)
	}

	// Parse amount with token decimals
	amount, err := format.ParseAmount(amountStr, decimals)
	if err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}

	// Check if balance is sufficient
	if balance.Cmp(amount) < 0 {
		return fmt.Errorf("insufficient balance. You have %s but need %s",
			format.Amount(balance, decimals),
			format.Amount(amount, decimals))
	}

	// Create send template
	template := &nom.AccountBlock{
		Version:         1,
		ChainIdentifier: 1,
		BlockType:       nom.BlockTypeUserSend,
		ToAddress:       toAddress,
		Amount:          amount,
		TokenStandard:   tokenStandard,
		Data:            nil,
	}

	// Send transaction
	fmt.Println("Sending transaction...")
	err = transaction.BuildAndSend(rpcClient.RpcClient, types.ParseAddressPanic(address), template, keypair)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	// Display success
	fmt.Printf("Successfully sent %s %s to %s\n",
		format.Amount(amount, decimals),
		format.FormatToken(amount, decimals, tokenStandard.String()),
		format.Cyan(toAddress.String()))

	return nil
}
