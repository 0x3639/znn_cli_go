package wallet

import (
	"fmt"
	"strconv"

	"github.com/0x3639/znn_cli_go/cmd"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
)

// deriveAddressesCmd derives a range of addresses from the wallet
var deriveAddressesCmd = &cobra.Command{
	Use:   "deriveAddresses <start> <end>",
	Short: "Derive a range of addresses from the wallet",
	Long: `Derive multiple addresses from the wallet using BIP44 derivation.
Displays addresses from index 'start' to 'end' (inclusive).

Example:
  znn-cli wallet deriveAddresses 0 10

This will display addresses at indices 0 through 10.

Requires --keyStore flag to specify which wallet to use.`,
	Args: cobra.ExactArgs(2),
	RunE: runDeriveAddresses,
}

func init() {
	walletCmd.AddCommand(deriveAddressesCmd)
}

func runDeriveAddresses(c *cobra.Command, args []string) error {
	// Parse arguments
	start, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid start index: %w", err)
	}

	end, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid end index: %w", err)
	}

	if start < 0 || end < start {
		return fmt.Errorf("invalid range: start=%d end=%d", start, end)
	}

	cfg := cmd.GetConfig()
	keystoreName := cmd.GetKeyStore()
	passphrase := cmd.GetPassphrase()

	// Load wallet
	keyStore, _, err := wallet.LoadWallet(cfg.Wallet.WalletDir, keystoreName, passphrase, 0)
	if err != nil {
		return err
	}

	// Derive addresses
	addresses, err := wallet.DeriveAddresses(keyStore, start, end)
	if err != nil {
		return fmt.Errorf("failed to derive addresses: %w", err)
	}

	// Display addresses
	fmt.Printf("Addresses for keyStore %s:\n", format.Green(keystoreName))
	for i, addr := range addresses {
		fmt.Printf("  %d\t%s\n", start+i, format.Cyan(addr))
	}

	return nil
}
