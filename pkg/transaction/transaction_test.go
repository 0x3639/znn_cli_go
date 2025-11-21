package transaction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestConstants verifies the package constants are correct
func TestConstants(t *testing.T) {
	tests := []struct {
		name     string
		actual   uint64
		expected uint64
	}{
		{
			name:     "MinPlasmaAmount",
			actual:   MinPlasmaAmount,
			expected: 21000,
		},
		{
			name:     "DefaultPoWDifficulty",
			actual:   DefaultPoWDifficulty,
			expected: 80000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.actual)
		})
	}
}

// Note: Full integration tests for Autofill, EnsurePlasmaOrPoW, Sign, Publish, and BuildAndSend
// require a live Zenon node connection and are better suited for integration test suites.
// These functions are tested indirectly through the CLI command tests with real node connections.
//
// The transaction package implements critical blockchain interaction logic that depends on:
// - Live RPC client connections to query account state
// - Real blockchain data (heights, hashes, momentum)
// - Cryptographic operations (signing, PoW generation)
// - Network communication (publishing transactions)
//
// Testing recommendations:
// 1. Use integration tests with a test node for complete transaction flow validation
// 2. Test individual commands (send, receive, etc.) which use this package
// 3. Manually test against testnet before mainnet deployment
