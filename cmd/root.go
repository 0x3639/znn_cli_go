// Package cmd implements the command-line interface for the Zenon Network CLI wallet.
package cmd

import (
	"fmt"
	"os"

	"github.com/0x3639/znn_cli_go/pkg/config"
	"github.com/spf13/cobra"
)

var (
	// cfgFile is the path to the configuration file
	cfgFile string

	// Global flags
	url        string
	keyStore   string
	passphrase string
	index      int
	verbose    bool

	// cfg holds the application configuration
	cfg *config.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "znn-cli",
	Short: "Zenon Network CLI Wallet",
	Long: `A command-line interface wallet for the Zenon Network.

This CLI provides complete functionality for interacting with the Zenon Network,
including wallet management, token transfers, staking, plasma operations, and more.

For more information, visit: https://github.com/0x3639/znn_cli_go`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global persistent flags (available to all subcommands)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.znn/cli-config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "WebSocket daemon URL (default: ws://127.0.0.1:35998)")
	rootCmd.PersistentFlags().StringVarP(&keyStore, "keyStore", "k", "", "keyStore file name")
	rootCmd.PersistentFlags().StringVarP(&passphrase, "passphrase", "p", "", "wallet passphrase (will prompt if not provided)")
	rootCmd.PersistentFlags().IntVarP(&index, "index", "i", 0, "address index in wallet")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error
	cfg, err = config.Load(cfgFile)
	if err != nil {
		if verbose {
			fmt.Fprintf(os.Stderr, "Warning: Failed to load config: %v\n", err)
		}
		// Use default config if loading fails
		cfg = config.DefaultConfig()
	}

	// Override config with command-line flags if provided
	if url != "" {
		cfg.Node.URL = url
	}
	if keyStore != "" {
		cfg.Wallet.DefaultKeyStore = keyStore
	}
	if index != 0 {
		cfg.Wallet.DefaultIndex = index
	}
	if verbose {
		cfg.Display.Verbose = true
	}
}

// GetConfig returns the current configuration
func GetConfig() *config.Config {
	return cfg
}

// GetPassphrase returns the passphrase from flag
func GetPassphrase() string {
	return passphrase
}

// GetKeyStore returns the keyStore name from flag or config
func GetKeyStore() string {
	if keyStore != "" {
		return keyStore
	}
	return cfg.Wallet.DefaultKeyStore
}

// GetIndex returns the address index from flag or config
func GetIndex() int {
	if index != 0 {
		return index
	}
	return cfg.Wallet.DefaultIndex
}
