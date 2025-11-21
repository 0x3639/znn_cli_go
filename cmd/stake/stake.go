package stake

import (
	"github.com/spf13/cobra"
)

// StakeCmd is the root command for staking operations
var StakeCmd = &cobra.Command{
	Use:   "stake",
	Short: "Staking operations",
	Long: `Staking operations for earning rewards.

Stake ZNN tokens for 1-12 months to earn rewards.
Minimum stake amount: 1 ZNN
Valid durations: 1-12 months (in 30-day increments)

Available subcommands:
  list     - List stake entries
  register - Stake ZNN for rewards
  revoke   - Cancel expired stake
  collect  - Collect staking rewards`,
}

func init() {
	// Subcommands will register themselves
}
