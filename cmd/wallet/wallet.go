// Package wallet implements wallet management commands
package wallet

import (
	"github.com/0x3639/znn_cli_go/cmd"
	"github.com/spf13/cobra"
)

// walletCmd represents the wallet command group
var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Wallet management commands",
	Long:  `Manage Zenon Network wallets: create, import, export, and derive addresses.`,
}

func init() {
	cmd.RootCmd().AddCommand(walletCmd)
}
