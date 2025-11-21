package wallet

import (
	"fmt"

	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/0x3639/znn_cli_go/pkg/wallet"
	"github.com/spf13/cobra"
)

// listCmd lists all available keyStores
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available keyStores",
	Long:  `Display a list of all keyStore files in the wallet directory (~/.znn/wallet/).`,
	RunE:  runList,
}

func init() {
	walletCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	// Create wallet manager
	mgr, err := wallet.NewManager("")
	if err != nil {
		return fmt.Errorf("failed to create wallet manager: %w", err)
	}

	// List all keyStores
	stores, err := mgr.List()
	if err != nil {
		return fmt.Errorf("failed to list wallets: %w", err)
	}

	if len(stores) == 0 {
		fmt.Println("No keyStores found")
		return nil
	}

	fmt.Println("Available keyStores:")
	for _, store := range stores {
		fmt.Printf("  %s\n", format.Green(store))
	}

	return nil
}
