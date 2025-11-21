package wallet

import (
	"fmt"

	"github.com/0x3639/znn_cli_go/cmd"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
)

// dumpMnemonicCmd displays the mnemonic for the current wallet
var dumpMnemonicCmd = &cobra.Command{
	Use:   "dumpMnemonic",
	Short: "Display the mnemonic phrase for the current wallet",
	Long: `Display the 24-word BIP39 mnemonic phrase for the current wallet.
This is used for wallet backup and recovery.

⚠️  WARNING: Never share your mnemonic with anyone. Anyone with access
to your mnemonic can access all funds in your wallet.

Requires --keyStore flag to specify which wallet to dump.`,
	RunE: runDumpMnemonic,
}

func init() {
	walletCmd.AddCommand(dumpMnemonicCmd)
}

func runDumpMnemonic(c *cobra.Command, args []string) error {
	cfg := cmd.GetConfig()
	keystoreName := cmd.GetKeyStore()
	passphrase := cmd.GetPassphrase()

	// Load wallet
	keyStore, _, err := wallet.LoadWallet(cfg.Wallet.WalletDir, keystoreName, passphrase, 0)
	if err != nil {
		return err
	}

	// Display mnemonic
	fmt.Printf("Mnemonic for keyStore %s:\n", format.Green(keystoreName))
	fmt.Println(format.Cyan(keyStore.Mnemonic))
	fmt.Println()
	fmt.Println(format.Yellow("⚠️  Keep this mnemonic safe and never share it with anyone!"))

	return nil
}
