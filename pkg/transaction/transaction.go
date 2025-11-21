// Package transaction provides transaction building, signing, and publishing utilities.
// It implements critical functionality not provided by the SDK, including transaction
// autofill and plasma/PoW decision logic.
package transaction

import (
	"fmt"

	"github.com/0x3639/znn-sdk-go/pow"
	"github.com/0x3639/znn-sdk-go/rpc_client"
	"github.com/0x3639/znn-sdk-go/wallet"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/zenon-network/go-zenon/chain/nom"
	"github.com/zenon-network/go-zenon/common/types"
)

const (
	// MinPlasmaAmount is the minimum plasma required for a transaction
	MinPlasmaAmount = 21000

	// DefaultPoWDifficulty is the default PoW difficulty when plasma is insufficient
	DefaultPoWDifficulty = 80000
)

// Autofill fills in the required fields for a transaction template.
// This includes:
//   - Height: The next block height for the account
//   - PreviousHash: Hash of the previous account block (or zero hash if first block)
//   - MomentumAcknowledged: Hash of the current frontier momentum
//
// This function must be called before signing and publishing a transaction.
func Autofill(c *rpc_client.RpcClient, address types.Address, template *nom.AccountBlock) error {
	// Get account info to determine height and previous hash
	accountInfo, err := c.LedgerApi.GetAccountInfoByAddress(address)
	if err != nil {
		return fmt.Errorf("failed to get account info: %w", err)
	}

	// Set height (next block number)
	template.Height = accountInfo.BlockCount + 1

	// Set previous hash (or zero hash if first block)
	if accountInfo.BlockCount > 0 {
		previousBlock, err := c.LedgerApi.GetAccountBlockByHeight(address, accountInfo.BlockCount)
		if err != nil {
			return fmt.Errorf("failed to get previous block: %w", err)
		}
		template.PreviousHash = previousBlock.Hash
	}
	// Note: Zero hash is automatically set by template creation

	// Get frontier momentum for acknowledgment
	momentum, err := c.LedgerApi.GetFrontierMomentum()
	if err != nil {
		return fmt.Errorf("failed to get frontier momentum: %w", err)
	}
	template.MomentumAcknowledged = momentum.Hash

	return nil
}

// EnsurePlasmaOrPoW checks if the account has sufficient plasma for the transaction.
// If not, it generates PoW. This must be called after Autofill and before Sign.
//
// Parameters:
//   - c: RPC client
//   - address: Account address
//   - template: Transaction template (must be autofilled)
//   - requiredPlasma: Amount of plasma required (0 to auto-calculate)
//
// Returns an error if PoW generation fails.
func EnsurePlasmaOrPoW(c *rpc_client.RpcClient, address types.Address, template *nom.AccountBlock, requiredPlasma uint64) error {
	if requiredPlasma == 0 {
		requiredPlasma = MinPlasmaAmount
	}

	// Check current plasma
	plasmaInfo, err := c.PlasmaApi.Get(address)
	if err != nil {
		// If we can't get plasma info, fall back to PoW
		format.Warning("Failed to check plasma, generating PoW")
		return generatePoW(template, DefaultPoWDifficulty)
	}

	if plasmaInfo.CurrentPlasma >= requiredPlasma {
		// Sufficient plasma available
		template.FusedPlasma = requiredPlasma
		template.Difficulty = 0
		if format.GetVerbose() {
			format.Info(fmt.Sprintf("Using plasma: %d (available: %d)", requiredPlasma, plasmaInfo.CurrentPlasma))
		}
		return nil
	}

	// Insufficient plasma, generate PoW
	if format.GetVerbose() {
		format.Warning(fmt.Sprintf("Insufficient plasma (have: %d, need: %d), generating PoW", plasmaInfo.CurrentPlasma, requiredPlasma))
	}
	return generatePoW(template, DefaultPoWDifficulty)
}

// generatePoW generates proof-of-work for a transaction template
func generatePoW(template *nom.AccountBlock, difficulty uint64) error {
	if format.GetVerbose() {
		format.Info(fmt.Sprintf("Generating PoW with difficulty %d...", difficulty))
	}

	nonce := pow.GeneratePoW(template.Hash.Bytes(), difficulty)

	template.Nonce = nonce
	template.Difficulty = difficulty
	template.FusedPlasma = 0

	if format.GetVerbose() {
		format.Success("PoW generation complete")
	}

	return nil
}

// Sign signs a transaction template with the given keypair.
// The template must be autofilled and have plasma/PoW set before signing.
func Sign(template *nom.AccountBlock, keypair *wallet.KeyPair) error {
	signature, err := keypair.Sign(template.Hash.Bytes())
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}
	template.Signature = signature
	return nil
}

// Publish publishes a signed transaction to the network.
func Publish(c *rpc_client.RpcClient, template *nom.AccountBlock) error {
	if err := c.LedgerApi.PublishRawTransaction(template); err != nil {
		return fmt.Errorf("failed to publish transaction: %w", err)
	}
	return nil
}

// BuildAndSend is a convenience function that autofills, ensures plasma/PoW,
// signs, and publishes a transaction in one call.
//
// This is the recommended way to send transactions from CLI commands.
func BuildAndSend(c *rpc_client.RpcClient, address types.Address, template *nom.AccountBlock, keypair *wallet.KeyPair) error {
	// Autofill transaction fields
	if err := Autofill(c, address, template); err != nil {
		return fmt.Errorf("autofill failed: %w", err)
	}

	// Ensure plasma or generate PoW
	if err := EnsurePlasmaOrPoW(c, address, template, 0); err != nil {
		return fmt.Errorf("plasma/PoW generation failed: %w", err)
	}

	// Sign transaction
	if err := Sign(template, keypair); err != nil {
		return fmt.Errorf("signing failed: %w", err)
	}

	// Publish transaction
	if err := Publish(c, template); err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	if format.GetVerbose() {
		format.Success("Transaction published successfully")
	}

	return nil
}
