// Package config provides configuration management for the Zenon CLI.
// It supports loading configuration from YAML files and environment variables
// using Viper, with sensible defaults for all settings.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Node    NodeConfig    `mapstructure:"node"`
	Wallet  WalletConfig  `mapstructure:"wallet"`
	Display DisplayConfig `mapstructure:"display"`
}

// NodeConfig contains Zenon node connection settings
type NodeConfig struct {
	URL           string        `mapstructure:"url"`
	AutoReconnect bool          `mapstructure:"auto_reconnect"`
	Timeout       time.Duration `mapstructure:"timeout"`
}

// WalletConfig contains wallet-related settings
type WalletConfig struct {
	DefaultKeyStore string `mapstructure:"default_keystore"`
	DefaultIndex    int    `mapstructure:"default_index"`
	WalletDir       string `mapstructure:"wallet_dir"`
}

// DisplayConfig contains display and output settings
type DisplayConfig struct {
	Colors  bool `mapstructure:"colors"`
	Verbose bool `mapstructure:"verbose"`
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	home, _ := os.UserHomeDir()
	return &Config{
		Node: NodeConfig{
			URL:           "ws://127.0.0.1:35998",
			AutoReconnect: true,
			Timeout:       30 * time.Second,
		},
		Wallet: WalletConfig{
			DefaultKeyStore: "",
			DefaultIndex:    0,
			WalletDir:       filepath.Join(home, ".znn", "wallet"),
		},
		Display: DisplayConfig{
			Colors:  true,
			Verbose: false,
		},
	}
}

// Load reads the configuration file and returns a Config.
// If cfgFile is empty, it looks for config in default locations.
func Load(cfgFile string) (*Config, error) {
	v := viper.New()

	// Set defaults
	defaults := DefaultConfig()
	v.SetDefault("node.url", defaults.Node.URL)
	v.SetDefault("node.auto_reconnect", defaults.Node.AutoReconnect)
	v.SetDefault("node.timeout", defaults.Node.Timeout)
	v.SetDefault("wallet.default_keystore", defaults.Wallet.DefaultKeyStore)
	v.SetDefault("wallet.default_index", defaults.Wallet.DefaultIndex)
	v.SetDefault("wallet.wallet_dir", defaults.Wallet.WalletDir)
	v.SetDefault("display.colors", defaults.Display.Colors)
	v.SetDefault("display.verbose", defaults.Display.Verbose)

	if cfgFile != "" {
		// Use config file from the flag
		v.SetConfigFile(cfgFile)
	} else {
		// Search for config in default locations
		home, err := os.UserHomeDir()
		if err != nil {
			return defaults, err
		}

		// Look for config in ~/.znn/
		znnDir := filepath.Join(home, ".znn")
		v.AddConfigPath(znnDir)
		v.SetConfigName("cli-config")
		v.SetConfigType("yaml")
	}

	// Read in environment variables that match
	v.SetEnvPrefix("ZNN")
	v.AutomaticEnv()

	// If a config file is found, read it in
	if err := v.ReadInConfig(); err != nil {
		// Config file not found; use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// Save writes the configuration to a file
func (c *Config) Save(path string) error {
	v := viper.New()
	v.Set("node", c.Node)
	v.Set("wallet", c.Wallet)
	v.Set("display", c.Display)

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := v.WriteConfigAs(path); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
