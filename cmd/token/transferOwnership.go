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

// transferOwnershipCmd transfers token ownership
var transferOwnershipCmd = &cobra.Command{
	Use:   "transferOwnership <tokenStandard> <newOwnerAddress>",
	Short: "Transfer token ownership",
	Long: `Transfer ownership of a token to a new address.

Requirements:
  - Must be current token owner
  - Cannot be undone

The new owner will have full control over the token,
including the ability to mint (if mintable) and transfer ownership again.

Example:
  znn-cli token transferOwnership zts1... z1qz...

Requires --keyStore flag to specify which wallet to use.`,
	Args: cobra.ExactArgs(2),
	RunE: runTransferOwnership,
}

func init() {
	TokenCmd.AddCommand(transferOwnershipCmd)
}

func runTransferOwnership(cmdCobra *cobra.Command, args []string) error {
	cfg, keystoreName, passphrase, index, err := getConfigAndFlags(cmdCobra)
	if err != nil {
		return err
	}

	// Parse arguments
	tokenStandardStr := args[0]
	newOwnerAddressStr := args[1]

	// Parse token standard
	tokenStandard, err := types.ParseZTS(tokenStandardStr)
	if err != nil {
		return fmt.Errorf("invalid token standard: %w", err)
	}

	// Parse new owner address
	newOwnerAddress, err := types.ParseAddress(newOwnerAddressStr)
	if err != nil {
		return fmt.Errorf("invalid new owner address: %w", err)
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

	// Display transfer info
	fmt.Printf("%s Transferring ownership of %s (%s)\n",
		format.Red("Warning!"),
		format.Magenta(token.TokenName),
		token.ZenonTokenStandard.String())
	fmt.Printf("  From: %s\n", address)
	fmt.Printf("  To: %s\n", newOwnerAddress.String())
	fmt.Println("This cannot be undone!")
	fmt.Println()

	// Create transfer ownership template
	template := rpcClient.TokenApi.UpdateToken(tokenStandard, newOwnerAddress, token.IsMintable, token.IsBurnable)

	// Send transaction
	err = transaction.BuildAndSend(rpcClient.RpcClient, parsedAddress, template, keypair)
	if err != nil {
		return fmt.Errorf("failed to transfer ownership: %w", err)
	}

	fmt.Println("Done")
	fmt.Printf("Successfully transferred ownership of %s to %s\n",
		format.Magenta(token.TokenSymbol),
		newOwnerAddress.String())

	return nil
}
