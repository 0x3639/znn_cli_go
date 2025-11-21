package token

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/0x3639/znn_cli_go/pkg/client"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/transaction"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
	"github.com/zenon-network/go-zenon/common/types"
)

const (
	// TokenIssueFee is the ZNN required to issue a token (1 ZNN)
	TokenIssueFee = 1 * 1e8
)

// issueCmd issues a new token
var issueCmd = &cobra.Command{
	Use:   "issue <name> <symbol> <domain> <totalSupply> <maxSupply> <decimals> <mintable> <burnable> <utility>",
	Short: "Issue a new token",
	Long: `Issue a new ZTS token.

Parameters:
  name        - Token name (1-40 characters)
  symbol      - Token symbol (1-10 characters, uppercase)
  domain      - Token domain/website (max 128 characters)
  totalSupply - Initial supply to mint
  maxSupply   - Maximum possible supply (must be >= totalSupply)
  decimals    - Number of decimals (0-18)
  mintable    - Can mint more tokens? (true/false)
  burnable    - Can burn tokens? (true/false)
  utility     - Is utility token? (true/false)

Cost: 1 ZNN

Example:
  znn-cli token issue MyToken MTK example.com 1000000 10000000 8 true true false

Requires --keyStore flag to specify which wallet to use.`,
	Args: cobra.ExactArgs(9),
	RunE: runIssue,
}

func init() {
	TokenCmd.AddCommand(issueCmd)
}

func runIssue(cmdCobra *cobra.Command, args []string) error {
	cfg, keystoreName, passphrase, index, err := getConfigAndFlags(cmdCobra)
	if err != nil {
		return err
	}

	// Parse arguments
	tokenName := args[0]
	tokenSymbol := strings.ToUpper(args[1])
	tokenDomain := args[2]
	totalSupplyStr := args[3]
	maxSupplyStr := args[4]
	decimalsStr := args[5]
	mintableStr := args[6]
	burnableStr := args[7]
	utilityStr := args[8]

	// Validate name and symbol
	if len(tokenName) < 1 || len(tokenName) > 40 {
		return fmt.Errorf("token name must be 1-40 characters")
	}
	if len(tokenSymbol) < 1 || len(tokenSymbol) > 10 {
		return fmt.Errorf("token symbol must be 1-10 characters")
	}
	if len(tokenDomain) > 128 {
		return fmt.Errorf("domain must be max 128 characters")
	}

	// Parse decimals
	decimals, err := strconv.ParseUint(decimalsStr, 10, 8)
	if err != nil || decimals > 18 {
		return fmt.Errorf("decimals must be 0-18")
	}

	// Parse supply amounts
	totalSupply, err := format.ParseAmount(totalSupplyStr, int(decimals))
	if err != nil {
		return fmt.Errorf("invalid total supply: %w", err)
	}

	maxSupply, err := format.ParseAmount(maxSupplyStr, int(decimals))
	if err != nil {
		return fmt.Errorf("invalid max supply: %w", err)
	}

	// Validate supplies
	if maxSupply.Cmp(totalSupply) < 0 {
		return fmt.Errorf("max supply must be >= total supply")
	}

	// Parse boolean flags
	mintable, err := strconv.ParseBool(mintableStr)
	if err != nil {
		return fmt.Errorf("mintable must be true or false")
	}

	burnable, err := strconv.ParseBool(burnableStr)
	if err != nil {
		return fmt.Errorf("burnable must be true or false")
	}

	utility, err := strconv.ParseBool(utilityStr)
	if err != nil {
		return fmt.Errorf("utility must be true or false")
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

	// Check ZNN balance for issuance fee
	accountInfo, err := rpcClient.LedgerApi.GetAccountInfoByAddress(parsedAddress)
	if err != nil {
		return fmt.Errorf("failed to get account info: %w", err)
	}

	requiredZnn := big.NewInt(TokenIssueFee)
	znnBalance, found := accountInfo.BalanceInfoMap[types.ZnnTokenStandard]
	if !found || znnBalance.Balance.Cmp(requiredZnn) < 0 {
		currentBalance := big.NewInt(0)
		if found {
			currentBalance = znnBalance.Balance
		}
		return fmt.Errorf("insufficient ZNN for issuance fee. You have %s but need %s",
			format.Amount(currentBalance, 8),
			format.Amount(requiredZnn, 8))
	}

	// Display token info
	fmt.Printf("Issuing token %s (%s)\n", format.Magenta(tokenName), format.Magenta(tokenSymbol))
	fmt.Printf("  Domain: %s\n", tokenDomain)
	fmt.Printf("  Total Supply: %s\n", format.Amount(totalSupply, int(decimals)))
	fmt.Printf("  Max Supply: %s\n", format.Amount(maxSupply, int(decimals)))
	fmt.Printf("  Decimals: %d\n", decimals)
	fmt.Printf("  Mintable: %v\n", mintable)
	fmt.Printf("  Burnable: %v\n", burnable)
	fmt.Printf("  Utility: %v\n", utility)
	fmt.Printf("  Cost: %s %s\n", format.Amount(requiredZnn, 8), format.Green("ZNN"))
	fmt.Println()

	// Create token issuance template
	template := rpcClient.TokenApi.IssueToken(
		tokenName,
		tokenSymbol,
		tokenDomain,
		totalSupply,
		maxSupply,
		uint8(decimals),
		mintable,
		burnable,
		utility,
	)

	// Send transaction
	err = transaction.BuildAndSend(rpcClient.RpcClient, parsedAddress, template, keypair)
	if err != nil {
		return fmt.Errorf("failed to issue token: %w", err)
	}

	fmt.Println("Done")
	fmt.Printf("Token %s successfully issued!\n", format.Magenta(tokenSymbol))

	return nil
}
