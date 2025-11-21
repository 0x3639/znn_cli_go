package sentinel

import (
	"github.com/spf13/cobra"
)

// SentinelCmd is the root command for sentinel operations
var SentinelCmd = &cobra.Command{
	Use:   "sentinel",
	Short: "Sentinel operations",
	Long: `Sentinel operations for network security and monitoring.

Sentinels are guardian nodes that:
- Monitor network health
- Validate consensus
- Provide additional security layer

Requirements to register a sentinel:
  - 5,000 ZNN
  - 50,000 QSR

Available subcommands:
  list        - List all sentinels
  register    - Register a new sentinel
  revoke      - Revoke your sentinel
  collect     - Collect sentinel rewards
  withdrawQsr - Withdraw deposited QSR`,
}

func init() {
	// Subcommands will register themselves
}
