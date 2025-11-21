package token

import (
	"github.com/spf13/cobra"
)

// TokenCmd is the root command for token operations
var TokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Token operations",
	Long: `Token operations for ZTS (Zenon Token Standard) management.

ZTS tokens are custom tokens on the Zenon Network.
Operations include issuing, minting, burning, and transferring ownership.

Token issuance cost: 1 ZNN

Available subcommands:
  list              - List all tokens
  getByStandard     - Get token by ZTS address
  getByOwner        - Get tokens owned by address
  issue             - Issue a new token
  mint              - Mint additional supply
  burn              - Burn tokens
  transferOwnership - Transfer token ownership
  disableMint       - Disable future minting`,
}

func init() {
	// Subcommands will register themselves
}
