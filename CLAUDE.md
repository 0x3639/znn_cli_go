# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go CLI wallet for the Zenon Network (`znn-cli`), replicating all 47 commands from the Dart CLI reference with an added TUI (Terminal UI) mode. The project uses cobra for CLI structure, and will integrate with the Zenon Go SDK at `github.com/0x3639/znn-sdk-go`.

**Module Path**: `github.com/0x3639/znn_cli_go`

## Build and Development Commands

```bash
# Development workflow
make deps          # Download and tidy dependencies
make build         # Build the znn-cli binary
make install       # Install to $GOPATH/bin
make run           # Build and run

# Code quality
make fmt           # Format code with gofmt and goimports
make vet           # Run go vet
make lint          # Run golangci-lint (must be installed)
make security      # Run gosec security scanner (must be installed)

# Testing
make test          # Run tests with race detector and coverage
make test-coverage # Show coverage report

# Complete check
make all           # Runs: clean deps fmt vet lint security test build

# Clean
make clean         # Remove binaries and coverage files
```

## Architecture and Structure

### Reference Implementation
- **Dart CLI source**: Located at `./reference/znn_cli_dart/` (gitignored, local only)
- **Key file**: `reference/znn_cli_dart/lib/cli_handler.dart` contains all command implementations
- **Line references**: ROADMAP.md contains line numbers for each command in the Dart CLI
- When implementing commands, always reference the Dart CLI for exact behavior, output format, and error handling

### Project Layout
```
cmd/               # Cobra command implementations (one file/package per command group)
pkg/               # Public reusable packages
  config/          # Viper config, ~/.znn/cli-config.yaml loading
  wallet/          # SDK KeyStoreManager wrapper, wallet loading logic
  client/          # RPC client wrapper with reconnection
  transaction/     # CRITICAL: Autofill, plasma/PoW logic, sign & publish
  format/          # Amount/duration formatting, token parsing, color helpers
internal/          # Private packages
  prompt/          # Password and confirmation prompts
  tui/             # Terminal UI components (bubbletea/tview)
  validation/      # Input validation
reference/         # Dart CLI reference (not committed)
```

### Key Architectural Patterns

**Global Flags** (defined in `cmd/root.go`):
- `--url` / `-u`: WebSocket URL (default: ws://127.0.0.1:35998)
- `--keyStore` / `-k`: Wallet name
- `--passphrase` / `-p`: Passphrase (prompts if not provided)
- `--index` / `-i`: BIP44 account index (default: 0)
- `--verbose` / `-v`: Enable verbose logging

**Wallet Loading Pattern**:
1. Determine which keyStore to use (from flag, config, or auto-detect if only one exists)
2. Get passphrase (from flag or secure prompt)
3. Load keyStore via SDK's KeyStoreManager
4. Derive keypair at specified BIP44 index
5. Connect to RPC client

**Transaction Flow** (must implement in `pkg/transaction/`):
1. Create template using SDK API (e.g., `client.TokenApi.IssueToken(...)`)
2. **Autofill**: Set height, previousHash, momentumAcknowledged
   - Height = accountInfo.BlockCount + 1
   - PreviousHash = last block hash (or zero hash if first block)
   - MomentumAcknowledged = current frontier momentum hash
3. **Plasma or PoW**: Check if sufficient plasma exists, else generate PoW
4. **Sign**: Use keypair.Sign(template.Hash.Bytes())
5. **Publish**: client.LedgerApi.PublishRawTransaction(template)

**Color Coding** (use `github.com/fatih/color`):
- Green: ZNN amounts, success messages
- Blue: QSR amounts
- Magenta: Custom ZTS tokens
- Red: Errors and warnings

**Amount Formatting**:
- Display human-readable amounts (e.g., "10.50000000 ZNN")
- Parse user input with decimal support
- Convert to/from base units using token decimals

## Critical Implementation Details

### Transaction Autofill (Not in SDK)
You MUST implement autofill logic in `pkg/transaction/autofill.go`:
```go
func AutofillTransaction(client *rpc_client.RpcClient, address types.Address, template *nom.AccountBlock) error {
    // Get account height
    accountInfo, err := client.LedgerApi.GetAccountInfoByAddress(address)
    template.Height = accountInfo.BlockCount + 1

    // Get previous block hash
    if accountInfo.BlockCount > 0 {
        lastBlock, err := client.LedgerApi.GetAccountBlockByHeight(address, accountInfo.BlockCount)
        template.PreviousHash = lastBlock.Hash
    }

    // Get frontier momentum
    momentum, err := client.LedgerApi.GetFrontierMomentum()
    template.MomentumAcknowledged = momentum.Hash

    return nil
}
```

### Plasma vs PoW Decision Logic
Implement in `pkg/transaction/`:
- Check `client.PlasmaApi.Get(address)` for current plasma
- If plasma >= required: Use plasma (set template.FusedPlasma, difficulty=0)
- Else: Generate PoW using `pow.GeneratePoW(template.Hash, difficulty)`

### Secure Password Input
Use `golang.org/x/term` to disable echo when reading passphrases:
```go
fd := int(os.Stdin.Fd())
oldState, _ := term.MakeRaw(fd)
defer term.Restore(fd, oldState)
password, _ := term.ReadPassword(fd)
```

### Wallet Storage
- Default directory: `~/.znn/wallet/`
- Multiple wallet support (user specifies with --keyStore)
- SDK handles Argon2 encryption + AES-256-GCM

## Dependencies to Install (Phase 2)

```bash
go get github.com/spf13/cobra
go get github.com/spf13/viper
go get github.com/charmbracelet/bubbletea
go get github.com/rivo/tview
go get github.com/charmbracelet/lipgloss
go get github.com/fatih/color
go get github.com/manifoldco/promptui
go get github.com/0x3639/znn-sdk-go
go get golang.org/x/term
```

## Command Implementation Checklist

When implementing each command:
1. **Reference Dart CLI**: Find exact line numbers in ROADMAP.md, read `reference/znn_cli_dart/lib/cli_handler.dart`
2. **Match behavior exactly**: Same arguments, same output format, same error messages
3. **Color coding**: Use green/blue/magenta/red as per Dart CLI
4. **godoc**: Document all exported functions with examples
5. **Error handling**: Wrap errors with context using `fmt.Errorf("...: %w", err)`
6. **Input validation**: Validate amounts, addresses, token standards before RPC calls
7. **Confirmation prompts**: Use for destructive/expensive operations (token issuance, pillar registration)

## Zenon Network Constants (from SDK)

```go
// Amounts
OneZnn = 100000000  // 1 ZNN = 10^8 base units
OneQsr = 100000000  // 1 QSR = 10^8 base units
CoinDecimals = 8

// Costs
PillarRegisterZnnAmount = 15000 * OneZnn
PillarRegisterQsrAmount = 150000 * OneQsr
SentinelRegisterZnnAmount = 5000 * OneZnn
SentinelRegisterQsrAmount = 50000 * OneQsr
TokenZtsIssueFeeInZnn = 1 * OneZnn

// Staking
StakeMinZnnAmount = 1 * OneZnn
StakeTimeMinSec = 30 * 24 * 60 * 60  // 1 month
StakeTimeMaxSec = 360 * 24 * 60 * 60 // 12 months

// Plasma
FuseMinQsrAmount = 10 * OneQsr
MinPlasmaAmount = 21000
```

## Configuration File Format

`~/.znn/cli-config.yaml`:
```yaml
node:
  url: ws://127.0.0.1:35998
  auto_reconnect: true
  timeout: 30s

wallet:
  default_keystore: main-wallet
  default_index: 0

display:
  colors: true
  verbose: false
```

## Testing Strategy

- Unit tests: All `pkg/` packages (format, validation logic)
- Table-driven tests: Amount parsing, duration formatting, token standard parsing
- Integration tests: Require testnet connection
- Target: >80% coverage for `pkg/` packages
- Security: All code must pass `make security` (gosec) with no HIGH/CRITICAL findings

## Important Notes

- **Never commit `reference/` directory** - it's gitignored and should stay local
- **Always run `make lint` before committing** - CI will enforce this
- **godoc comments required** for all exported symbols
- **BIP44 path**: `m/44'/73404'/account'/0'/0'` (73404 is Zenon's coin type)
- **Address format**: Starts with "z1", 40 characters total
- **Token standards**: "ZNN", "QSR", or full ZTS address (40 chars starting with "zts1")

## Progress Tracking

See ROADMAP.md for detailed implementation status. Update checkboxes as commands are completed.

## SDK Usage Examples

```go
// Create RPC client
client, err := rpc_client.NewRpcClient("ws://127.0.0.1:35998")
defer client.Stop()

// Load wallet
manager, _ := wallet.NewKeyStoreManager("~/.znn/wallet")
ks, _ := manager.ReadKeyStore(passphrase, keystoreName)
kp, _ := ks.GetKeyPair(0)
address, _ := kp.GetAddress()

// Get balance
accountInfo, _ := client.LedgerApi.GetAccountInfoByAddress(*address)
znnBalance := accountInfo.BalanceInfoMap[types.ZnnTokenStandard].Balance

// Send transaction (after autofill, plasma/PoW, sign)
err := client.LedgerApi.PublishRawTransaction(template)
```
