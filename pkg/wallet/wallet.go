// Package wallet provides a wrapper around the Zenon SDK wallet functionality
// with CLI-specific helpers for loading wallets and managing keypairs.
package wallet

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/0x3639/znn-sdk-go/wallet"
	"github.com/0x3639/znn_cli_go/internal/prompt"
	"github.com/0x3639/znn_cli_go/pkg/format"
)

// Manager wraps the SDK KeyStoreManager with CLI-specific functionality
type Manager struct {
	manager *wallet.KeyStoreManager
}

// NewManager creates a new wallet manager.
// If walletDir is empty, uses the default ~/.znn/wallet/ directory.
func NewManager(walletDir string) (*Manager, error) {
	if walletDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		walletDir = filepath.Join(home, ".znn", "wallet")
	}

	// Ensure wallet directory exists
	if err := os.MkdirAll(walletDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create wallet directory: %w", err)
	}

	mgr, err := wallet.NewKeyStoreManager(walletDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet manager: %w", err)
	}

	return &Manager{manager: mgr}, nil
}

// List returns a list of all wallet names in the wallet directory
func (m *Manager) List() ([]string, error) {
	return m.manager.ListAllKeyStores()
}

// CreateNew creates a new wallet with a random BIP39 mnemonic
func (m *Manager) CreateNew(passphrase, name string) (*wallet.KeyStore, error) {
	return m.manager.CreateNew(passphrase, name)
}

// CreateFromMnemonic creates a wallet from an existing mnemonic
func (m *Manager) CreateFromMnemonic(mnemonic, passphrase, name string) (*wallet.KeyStore, error) {
	return m.manager.CreateFromMnemonic(mnemonic, passphrase, name)
}

// Load loads a keyStore by name with the given passphrase
func (m *Manager) Load(passphrase, name string) (*wallet.KeyStore, error) {
	return m.manager.ReadKeyStore(passphrase, name)
}

// LoadWallet loads a wallet and returns a keypair at the specified index.
// This is the main entry point for CLI commands that need wallet access.
//
// Parameters:
//   - walletDir: Directory containing wallets (empty for default ~/.znn/wallet/)
//   - keystoreName: Name of the keyStore file (empty to auto-detect)
//   - passphrase: Wallet passphrase (empty to prompt)
//   - index: BIP44 account index
//
// Returns the keyStore, keypair, and any error encountered.
func LoadWallet(walletDir, keystoreName, passphrase string, index int) (*wallet.KeyStore, *wallet.KeyPair, error) {
	// Create manager
	mgr, err := NewManager(walletDir)
	if err != nil {
		return nil, nil, err
	}

	// Determine which keyStore to use
	if keystoreName == "" {
		wallets, err := mgr.List()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to list wallets: %w", err)
		}

		if len(wallets) == 0 {
			return nil, nil, fmt.Errorf("no wallets found in %s. Create one with: znn-cli wallet.createNew", walletDir)
		}

		if len(wallets) == 1 {
			keystoreName = wallets[0]
			if os.Getenv("ZNN_CLI_VERBOSE") == "1" {
				format.Info(fmt.Sprintf("Using wallet: %s", keystoreName))
			}
		} else {
			// Multiple wallets found, ask user to specify
			return nil, nil, fmt.Errorf("multiple wallets found: %v. Specify with --keyStore flag", wallets)
		}
	}

	// Get passphrase if not provided
	if passphrase == "" {
		pass, err := prompt.Password("Enter passphrase: ")
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read passphrase: %w", err)
		}
		passphrase = pass
	}

	// Load keyStore
	ks, err := mgr.Load(passphrase, keystoreName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load wallet: %w", err)
	}

	// Get keypair at index
	kp, err := ks.GetKeyPair(index)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get keypair at index %d: %w", index, err)
	}

	return ks, kp, nil
}

// GetAddress returns the address for a keypair as a string
func GetAddress(kp *wallet.KeyPair) (string, error) {
	addr, err := kp.GetAddress()
	if err != nil {
		return "", fmt.Errorf("failed to get address: %w", err)
	}
	return addr.String(), nil
}

// DeriveAddresses derives a range of addresses from a keyStore
func DeriveAddresses(ks *wallet.KeyStore, start, end int) ([]string, error) {
	if start < 0 || end < start {
		return nil, fmt.Errorf("invalid range: start=%d end=%d", start, end)
	}

	addresses := make([]string, 0, end-start+1)
	for i := start; i <= end; i++ {
		kp, err := ks.GetKeyPair(i)
		if err != nil {
			return nil, fmt.Errorf("failed to get keypair at index %d: %w", i, err)
		}

		addr, err := GetAddress(kp)
		if err != nil {
			return nil, err
		}

		addresses = append(addresses, addr)
	}

	return addresses, nil
}
