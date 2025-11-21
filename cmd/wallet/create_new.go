package wallet

import (
	"fmt"

	"github.com/0x3639/znn_cli_go/internal/prompt"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
)

// createNewCmd creates a new wallet with a random mnemonic
var createNewCmd = &cobra.Command{
	Use:   "createNew [name]",
	Short: "Create a new wallet with a random BIP39 mnemonic",
	Long: `Create a new wallet with a randomly generated 24-word BIP39 mnemonic.
The wallet will be encrypted with the provided passphrase and stored in ~/.znn/wallet/.

If no name is provided, a default name will be generated.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runCreateNew,
}

func init() {
	walletCmd.AddCommand(createNewCmd)
}

func runCreateNew(cmd *cobra.Command, args []string) error {
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
	if len(args) > 0 {
		name = args[0]
	}

	// Create wallet manager
	mgr, err := wallet.NewManager("")
	if err != nil {
		return fmt.Errorf("failed to create wallet manager: %w", err)
	}

	// Create new wallet
	keyStore, err := mgr.CreateNew(passphrase, name)
	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
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
	format.Success(fmt.Sprintf("keyStore successfully created: %s", walletName))

	// Display the mnemonic (important for backup)
	fmt.Println("\n" + format.Yellow("⚠️  IMPORTANT: Write down your mnemonic phrase and store it safely!"))
	fmt.Println(format.Yellow("This is the ONLY way to recover your wallet if you lose access."))
	fmt.Println("\nMnemonic:")
	fmt.Println(format.Cyan(keyStore.Mnemonic))
	fmt.Println()

	return nil
}
