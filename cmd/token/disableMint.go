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

// disableMintCmd disables minting for a token
var disableMintCmd = &cobra.Command{
	Use:   "disableMint <tokenStandard>",
	Short: "Disable future minting",
	Long: `Disable the ability to mint additional supply for a token.

Requirements:
  - Must be token owner
  - Token must currently be mintable
  - Cannot be undone

Once minting is disabled, no one (including the owner) can ever
mint additional tokens. The total supply becomes fixed.

Example:
  znn-cli token disableMint zts1...

Requires --keyStore flag to specify which wallet to use.`,
	Args: cobra.ExactArgs(1),
	RunE: runDisableMint,
}

func init() {
	TokenCmd.AddCommand(disableMintCmd)
}

func runDisableMint(cmdCobra *cobra.Command, args []string) error {
	cfg, keystoreName, passphrase, index, err := getConfigAndFlags(cmdCobra)
	if err != nil {
		return err
	}

	// Parse token standard
	tokenStandardStr := args[0]
	tokenStandard, err := types.ParseZTS(tokenStandardStr)
	if err != nil {
		return fmt.Errorf("invalid token standard: %w", err)
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

	// Get token info
	token, err := rpcClient.TokenApi.GetByZts(tokenStandard)
	if err != nil {
		return fmt.Errorf("failed to get token info: %w", err)
	}

	// Verify ownership
	if token.Owner.String() != address {
		return fmt.Errorf("you do not own this token. Owner is %s", token.Owner.String())
	}

	// Verify currently mintable
	if !token.IsMintable {
		return fmt.Errorf("token minting is already disabled")
	}

	// Display disable mint info
	fmt.Printf("%s Disabling minting for %s (%s)\n",
		format.Red("Warning!"),
		format.Magenta(token.TokenName),
		token.ZenonTokenStandard.String())
	fmt.Println("This will permanently fix the total supply and cannot be undone!")
	fmt.Println()

	// Create disable mint template (update token with mintable=false)
	template := rpcClient.TokenApi.UpdateToken(tokenStandard, token.Owner, false, token.IsBurnable)

	// Send transaction
	err = transaction.BuildAndSend(rpcClient.RpcClient, parsedAddress, template, keypair)
	if err != nil {
		return fmt.Errorf("failed to disable minting: %w", err)
	}

	fmt.Println("Done")
	fmt.Printf("Successfully disabled minting for %s\n", format.Magenta(token.TokenSymbol))

	return nil
}
