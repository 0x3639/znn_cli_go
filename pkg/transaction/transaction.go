// Package transaction provides transaction building, signing, and publishing utilities.
// It implements critical functionality not provided by the SDK, including transaction
// autofill and plasma/PoW decision logic.
package transaction

import (
	"fmt"
	"math/big"

	"github.com/0x3639/znn-sdk-go/wallet"
	"github.com/zenon-network/go-zenon/chain/nom"
	"github.com/zenon-network/go-zenon/common/types"
	"github.com/zenon-network/go-zenon/pow"
	"github.com/zenon-network/go-zenon/rpc/api/embedded"

	rpc_client "github.com/0x3639/znn-sdk-go/rpc_client"
)

const (
	// MinPlasmaAmount is the minimum plasma required for a transaction
	MinPlasmaAmount = 21000

	// DefaultPoWDifficulty is the default PoW difficulty when plasma is insufficient
	DefaultPoWDifficulty = 80000
)

// Autofill sets the Height, PreviousHash, and MomentumAcknowledged fields on a transaction template.
// These fields must be set before signing and publishing a transaction.
//
// The function:
//  1. Gets current account height and increments by 1 for new block
//  2. Sets PreviousHash to frontier block hash (if height > 0)
//  3. Sets MomentumAcknowledged to current frontier momentum
//
// Parameters:
//   - c: RPC client for querying account and momentum info
//   - address: Address of the account creating the transaction
//   - template: AccountBlock template to autofill
//
// Returns an error if unable to query account info or frontier momentum.
func Autofill(c *rpc_client.RpcClient, address types.Address, template *nom.AccountBlock) error {
	// Set the address field
	template.Address = address

	// Get account info to determine height
	accountInfo, err := c.LedgerApi.GetAccountInfoByAddress(address)
	if err != nil {
		return fmt.Errorf("failed to get account info: %w", err)
	}

	// Set height (next block = current height + 1)
	template.Height = accountInfo.AccountHeight + 1

	// Set previous hash (if not first block)
	if accountInfo.AccountHeight > 0 {
		frontierBlock, err := c.LedgerApi.GetFrontierAccountBlock(address)
		if err != nil {
			return fmt.Errorf("failed to get frontier account block: %w", err)
		}
		template.PreviousHash = frontierBlock.Hash
	} else {
		// First block has zero hash as previous
		template.PreviousHash = types.ZeroHash
	}

	// Get frontier momentum for acknowledgment
	momentum, err := c.LedgerApi.GetFrontierMomentum()
	if err != nil {
		return fmt.Errorf("failed to get frontier momentum: %w", err)
	}

	template.MomentumAcknowledged = types.HashHeight{
		Hash:   momentum.Hash,
		Height: momentum.Height,
	}

	return nil
}

// EnsurePlasmaOrPoW checks if the account has sufficient plasma for the transaction.
// If not, it generates Proof of Work with the required difficulty.
//
// The function:
//  1. Computes transaction hash
//  2. Queries required PoW difficulty based on available plasma
//  3. If plasma sufficient (difficulty = 0), returns immediately
//  4. Otherwise, generates PoW with required difficulty
//
// Parameters:
//   - c: RPC client for querying plasma requirements
//   - address: Address of the account creating the transaction
//   - template: AccountBlock template (must already have hash computed)
//
// Returns an error if unable to query plasma or generate PoW.
func EnsurePlasmaOrPoW(c *rpc_client.RpcClient, address types.Address, template *nom.AccountBlock) error {
	// Check required PoW difficulty
	toAddr := &template.ToAddress
	param := embedded.GetRequiredParam{
		SelfAddr:  address,
		BlockType: template.BlockType,
		ToAddr:    toAddr,
		Data:      template.Data,
	}

	result, err := c.PlasmaApi.GetRequiredPoWForAccountBlock(param)
	if err != nil {
		return fmt.Errorf("failed to get required PoW: %w", err)
	}

	// If plasma is sufficient, no PoW needed
	if result.RequiredDifficulty == 0 {
		template.FusedPlasma = result.BasePlasma
		return nil
	}

	// Generate PoW
	difficulty := result.RequiredDifficulty
	if difficulty < DefaultPoWDifficulty {
		difficulty = DefaultPoWDifficulty
	}

	// Generate PoW nonce using SDK function
	difficultyBig := new(big.Int).SetUint64(difficulty)
	nonceBytes := pow.GetPoWNonce(difficultyBig, template.Hash)

	// Set the nonce
	copy(template.Nonce.Data[:], nonceBytes)
	template.Difficulty = difficulty

	return nil
}

// Sign signs the transaction with the provided keypair.
// The transaction hash must be computed before calling this function.
//
// Parameters:
//   - template: AccountBlock to sign (must have hash computed)
//   - keypair: Wallet keypair to use for signing
//
// Returns an error if signing fails.
func Sign(template *nom.AccountBlock, keypair *wallet.KeyPair) error {
	signature, err := keypair.Sign(template.Hash.Bytes())
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	publicKey, err := keypair.GetPublicKey()
	if err != nil {
		return fmt.Errorf("failed to get public key: %w", err)
	}

	copy(template.Signature[:], signature)
	copy(template.PublicKey[:], publicKey)

	return nil
}

// Publish publishes a signed and finalized transaction to the network.
//
// The transaction must be fully prepared:
//   - Autofilled (height, previousHash, momentumAcknowledged)
//   - Hash computed
//   - PoW/Plasma ensured
//   - Signed
//
// Parameters:
//   - c: RPC client for publishing
//   - template: Fully prepared AccountBlock
//
// Returns an error if the transaction is rejected by the node.
func Publish(c *rpc_client.RpcClient, template *nom.AccountBlock) error {
	return c.LedgerApi.PublishRawTransaction(template)
}

// BuildAndSend is a convenience function that performs the complete transaction flow:
//  1. Autofill (height, previousHash, momentumAcknowledged)
//  2. Compute hash
//  3. Ensure plasma or generate PoW
//  4. Sign with keypair
//  5. Publish to network
//
// This is the recommended way to send transactions as it handles all steps correctly.
//
// Parameters:
//   - c: RPC client for querying and publishing
//   - address: Address of the account creating the transaction
//   - template: AccountBlock template (ToAddress, Amount, TokenStandard, Data, etc.)
//   - keypair: Wallet keypair to use for signing
//
// Returns an error if any step fails.
//
// Example:
//
//	template := c.LedgerApi.SendTemplate(toAddress, types.ZnnTokenStandard, amount, nil)
//	err := transaction.BuildAndSend(c, myAddress, template, keypair)
//	if err != nil {
//	    return fmt.Errorf("failed to send: %w", err)
//	}
func BuildAndSend(c *rpc_client.RpcClient, address types.Address, template *nom.AccountBlock, keypair *wallet.KeyPair) error {
	// 1. Autofill
	if err := Autofill(c, address, template); err != nil {
		return fmt.Errorf("autofill failed: %w", err)
	}

	// 2. Compute hash
	template.Hash = template.ComputeHash()

	// 3. Ensure plasma or generate PoW
	if err := EnsurePlasmaOrPoW(c, address, template); err != nil {
		return fmt.Errorf("plasma/PoW failed: %w", err)
	}

	// 4. Sign
	if err := Sign(template, keypair); err != nil {
		return fmt.Errorf("signing failed: %w", err)
	}

	// 5. Publish
	if err := Publish(c, template); err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	return nil
}
