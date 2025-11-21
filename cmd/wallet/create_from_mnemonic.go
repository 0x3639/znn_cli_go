package wallet

import (
	"fmt"

	"github.com/0x3639/znn_cli_go/internal/prompt"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
)

// createFromMnemonicCmd imports a wallet from an existing mnemonic
var createFromMnemonicCmd = &cobra.Command{
	Use:   "createFromMnemonic <mnemonic> [name]",
	Short: "Import a wallet from a BIP39 mnemonic phrase",
	Long: `Import/restore a wallet from an existing 24-word BIP39 mnemonic phrase.
The wallet will be encrypted with the provided passphrase and stored in ~/.znn/wallet/.

The mnemonic should be provided as a quoted string with words separated by spaces.
If no name is provided, a default name will be generated.

Example:
  znn-cli wallet createFromMnemonic "word1 word2 ... word24"`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runCreateFromMnemonic,
}

func init() {
	walletCmd.AddCommand(createFromMnemonicCmd)
}

func runCreateFromMnemonic(cmd *cobra.Command, args []string) error {
	mnemonic := args[0]

	// Get passphrase
	passphrase, err := prompt.PasswordWithConfirm("Enter passphrase for new wallet: ")
	if err != nil {
		return fmt.Errorf("failed to read passphrase: %w", err)
	}

	if passphrase == "" {
		return fmt.Errorf("passphrase cannot be empty")
	}

	// Get optional wallet name
	var name string
	if len(args) > 1 {
		name = args[1]
	}

	// Create wallet manager
	mgr, err := wallet.NewManager("")
	if err != nil {
		return fmt.Errorf("failed to create wallet manager: %w", err)
	}

	// Create wallet from mnemonic
	_, err = mgr.CreateFromMnemonic(mnemonic, passphrase, name)
	if err != nil {
		return fmt.Errorf("failed to create wallet from mnemonic: %w", err)
	}

	// Get the wallet file name
	walletName := name
	if walletName == "" {
		// The SDK generates a name, we need to find it
		stores, _ := mgr.List()
		if len(stores) > 0 {
			walletName = stores[len(stores)-1] // Last created
		}
	}

	// Display success message
	format.Success(fmt.Sprintf("keyStore successfully created from mnemonic: %s", walletName))

	return nil
}
