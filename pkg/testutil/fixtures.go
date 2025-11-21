// Package testutil provides testing utilities and fixtures for the znn-cli.
package testutil

import (
	"math/big"

	"github.com/zenon-network/go-zenon/common/types"
)

// Test addresses for different scenarios
var (
	// ValidAddress is a valid Zenon address for testing
	ValidAddress = types.ParseAddressPanic("z1qzal6c5s9rjnnxd2z672tx3apscy5s5qqhslq5")

	// ValidAddress2 is another valid Zenon address for testing
	ValidAddress2 = types.ParseAddressPanic("z1qqjnwjjpnue8xmmpanz6csze6tcmtzzdtfsww7")

	// ProducerAddress is a valid producer address for pillar testing
	ProducerAddress = types.ParseAddressPanic("z1qr4pexnnfaexqqz8nscjjcsajy5hdqfkgadvwx")
)

// Test token standards
var (
	// TestZTS is a sample ZTS token standard for testing
	TestZTS = types.ParseZTSPanic("zts1qsrxxxxxxxxxxxxxmrhjll")
)

// Test hashes
var (
	// ValidHash is a valid transaction hash for testing
	ValidHash = types.HexToHashPanic("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")

	// ValidHash2 is another valid hash for testing
	ValidHash2 = types.HexToHashPanic("fedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321")
)

// AmountFixtures provides common amount values for testing
type AmountFixtures struct {
	OneZNN         *big.Int // 1 ZNN = 10^8 base units
	OneQSR         *big.Int // 1 QSR = 10^8 base units
	HalfZNN        *big.Int // 0.5 ZNN
	Zero           *big.Int // 0
	MaxSupply      *big.Int // Large supply for token testing
	MinStakeAmount *big.Int // 1 ZNN minimum stake
	MinFuseAmount  *big.Int // 10 QSR minimum fuse
}

// NewAmountFixtures creates a new set of amount fixtures
func NewAmountFixtures() *AmountFixtures {
	return &AmountFixtures{
		OneZNN:         big.NewInt(1e8),
		OneQSR:         big.NewInt(1e8),
		HalfZNN:        big.NewInt(5e7),
		Zero:           big.NewInt(0),
		MaxSupply:      big.NewInt(1000000 * 1e8),
		MinStakeAmount: big.NewInt(1e8),
		MinFuseAmount:  big.NewInt(10 * 1e8),
	}
}

// MnemonicFixture is a valid BIP39 mnemonic for testing wallet operations
const MnemonicFixture = "route become dream access impulse price inform obtain engage ski believe awful absent pig thing vibrant possible exotic flee pepper marble rural fire fancy"

// InvalidAddresses provides a list of invalid address formats for testing
var InvalidAddresses = []string{
	"",   // Empty
	"z1", // Too short
	"z1qzal6c5s9rjnnxd2z672tx3apscy5s5qqhslq",   // Missing last character
	"x1qzal6c5s9rjnnxd2z672tx3apscy5s5qqhslq5",  // Wrong prefix
	"z1qzal6c5s9rjnnxd2z672tx3apscy5s5qqhslq5z", // Too long
	"invalid-address",                           // Invalid format
	"z1QZAL6C5S9RJNNXD2Z672TX3APSCY5S5QQHSLQ5",  // Uppercase (invalid)
}

// InvalidTokenStandards provides invalid token standard strings for testing
var InvalidTokenStandards = []string{
	"",                           // Empty
	"ZTS",                        // Too short
	"zts1qsrxxxxxxxxxxxxxmrhjl",  // Missing character
	"znn",                        // Lowercase ZNN
	"qsr",                        // Lowercase QSR
	"invalid-token",              // Invalid format
	"zts1QSRXXXXXXXXXXXXXMRHJLL", // Uppercase (invalid)
}

// KeyStoreNameFixture is a sample keystore name for testing
const KeyStoreNameFixture = "test-wallet"

// PassphraseFixture is a sample passphrase for testing
const PassphraseFixture = "test-passphrase-123"
