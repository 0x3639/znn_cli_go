package pillar

import (
	"github.com/spf13/cobra"
)

// PillarCmd is the root command for pillar operations
var PillarCmd = &cobra.Command{
	Use:   "pillar",
	Short: "Pillar operations",
	Long: `Pillar operations for network consensus and governance.

Pillars are the backbone of the Zenon Network, responsible for:
- Producing momentums (blocks)
- Validating transactions
- Network governance

Requirements to register a pillar:
  - 15,000 ZNN
  - 150,000 QSR

Available subcommands:
  list        - List all pillars
  register    - Register a new pillar
  revoke      - Revoke your pillar
  delegate    - Delegate to a pillar
  undelegate  - Remove delegation
  collect     - Collect pillar rewards
  withdrawQsr - Withdraw deposited QSR`,
}

func init() {
	// Subcommands will register themselves
}
