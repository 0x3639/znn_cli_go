package sentinel

import (
	"github.com/0x3639/znn_cli_go/pkg/config"
	"github.com/spf13/cobra"
)

// getConfigAndFlags extracts configuration and flags from the command
func getConfigAndFlags(cmd *cobra.Command) (*config.Config, string, string, int, error) {
	keystoreName, _ := cmd.Flags().GetString("keyStore")
	passphrase, _ := cmd.Flags().GetString("passphrase")
	index, _ := cmd.Flags().GetInt("index")
	url, _ := cmd.Flags().GetString("url")
	configFile, _ := cmd.Flags().GetString("config")

	// Load configuration
	cfg, err := config.Load(configFile)
	if err != nil {
		cfg = config.DefaultConfig()
	}
	if url != "" {
		cfg.Node.URL = url
	}

	return cfg, keystoreName, passphrase, index, nil
}
