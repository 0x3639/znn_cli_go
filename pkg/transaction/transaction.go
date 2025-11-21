// Package transaction provides transaction building, signing, and publishing utilities.
// It implements critical functionality not provided by the SDK, including transaction
// autofill and plasma/PoW decision logic.
//
// NOTE: This package is a placeholder for Phase 4. The SDK API types need to be
// properly analyzed and the functions below need to be implemented correctly.
// These functions are commented out to allow Phase 3 (wallet commands) to compile.
package transaction

const (
	// MinPlasmaAmount is the minimum plasma required for a transaction
	MinPlasmaAmount = 21000

	// DefaultPoWDifficulty is the default PoW difficulty when plasma is insufficient
	DefaultPoWDifficulty = 80000
)

// TODO: Phase 4 - Implement transaction helpers
// The following functions need to be implemented with correct SDK API types:
//
// - Autofill(c *rpc_client.RpcClient, address types.Address, template *nom.AccountBlock) error
//   Sets height, previousHash, and momentumAcknowledged
//
// - EnsurePlasmaOrPoW(c *rpc_client.RpcClient, address types.Address, template *nom.AccountBlock, requiredPlasma uint64) error
//   Checks plasma availability and generates PoW if needed
//
// - Sign(template *nom.AccountBlock, keypair *wallet.KeyPair) error
//   Signs the transaction with the keypair
//
// - Publish(c *rpc_client.RpcClient, template *nom.AccountBlock) error
//   Publishes the signed transaction to the network
//
// - BuildAndSend(c *rpc_client.RpcClient, address types.Address, template *nom.AccountBlock, keypair *wallet.KeyPair) error
//   Convenience function that calls all of the above in sequence
