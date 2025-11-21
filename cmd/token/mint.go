package token

import (
	"fmt"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/transaction"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
	"github.com/zenon-network/go-zenon/common/types"
)

// mintCmd mints additional token supply
var mintCmd = &cobra.Command{
	Use:   "mint <tokenStandard> <amount> <receiveAddress>",
	Short: "Mint additional supply",
	Long: `Mint additional supply for a token.

Requirements:
  - Must be token owner
  - Token must be mintable
  - New total supply must not exceed max supply

The minted tokens will be sent to the specified receive address.

Example:
  znn-cli token mint zts1... 1000 z1qz...

Requires --keyStore flag to specify which wallet to use.`,
	Args: cobra.ExactArgs(3),
	RunE: runMint,
}

func init() {
	TokenCmd.AddCommand(mintCmd)
}

func runMint(cmdCobra *cobra.Command, args []string) error {
	cfg, keystoreName, passphrase, index, err := getConfigAndFlags(cmdCobra)
	if err != nil {
		return err
	}

	// Parse arguments
	tokenStandardStr := args[0]
	amountStr := args[1]
	receiveAddressStr := args[2]

	// Parse token standard
	tokenStandard, err := types.ParseZTS(tokenStandardStr)
	if err != nil {
		return fmt.Errorf("invalid token standard: %w", err)
	}

	// Parse receive address
	receiveAddress, err := types.ParseAddress(receiveAddressStr)
	if err != nil {
		return fmt.Errorf("invalid receive address: %w", err)
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

	// Get token info
	token, err := rpcClient.TokenApi.GetByZts(tokenStandard)
	if err != nil {
		return fmt.Errorf("failed to get token info: %w", err)
	}

	// Verify ownership
	if token.Owner.String() != address {
		return fmt.Errorf("you do not own this token. Owner is %s", token.Owner.String())
	}

	// Verify mintable
	if !token.IsMintable {
		return fmt.Errorf("token is not mintable")
	}

	// Parse amount with token decimals
	amount, err := format.ParseAmount(amountStr, int(token.Decimals))
	if err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}

	// Display mint info
	fmt.Printf("Minting %s %s (%s)\n",
		format.Amount(amount, int(token.Decimals)),
		format.Magenta(token.TokenSymbol),
		token.ZenonTokenStandard.String())
	fmt.Printf("  Receive address: %s\n", receiveAddress.String())
	fmt.Println()

	// Create mint template
	template := rpcClient.TokenApi.Mint(tokenStandard, amount, receiveAddress)

	// Send transaction
	err = transaction.BuildAndSend(rpcClient.RpcClient, parsedAddress, template, keypair)
	if err != nil {
		return fmt.Errorf("failed to mint tokens: %w", err)
	}

	fmt.Println("Done")
	fmt.Printf("Successfully minted %s %s\n",
		format.Amount(amount, int(token.Decimals)),
		format.Magenta(token.TokenSymbol))

	return nil
}
