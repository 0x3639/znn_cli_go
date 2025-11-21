# Zenon Go CLI (znn-cli)

A complete command-line interface wallet for the Zenon Network, written in Go. This CLI replicates all core functionality from the [Dart CLI](https://github.com/zenon-network/znn_cli_dart) with improved usability and performance.

## Features

- **42 Commands** covering all Zenon Network core operations
- **Wallet Management**: Create, import, export wallets with BIP39 mnemonic support
- **Transactions**: Send, receive, auto-receive with plasma or PoW
- **Staking**: Stake ZNN for rewards (1-12 months)
- **Plasma**: Fuse QSR to generate plasma for feeless transactions
- **Pillar Operations**: Register, delegate, collect rewards
- **Sentinel Operations**: Register, collect rewards
- **Token Management**: Issue, mint, burn, transfer ZTS tokens
- **Security**: Comprehensive input validation, secure password handling
- **Well-tested**: Go vet clean, formatted code, production-ready

## Installation

### Prerequisites

- Go 1.21 or higher
- A running Zenon node (default: `ws://127.0.0.1:35998`)

### From Source

```bash
# Clone the repository
git clone https://github.com/0x3639/znn_cli_go.git
cd znn_cli_go

# Install dependencies
make deps

# Build and install
make install
```

The binary will be installed to `$GOPATH/bin/znn-cli`.

### Manual Build

```bash
go build -o znn-cli
./znn-cli --help
```

## Quick Start

### 1. Create a Wallet

```bash
# Create a new wallet
znn-cli wallet createNew

# Or import from mnemonic
znn-cli wallet createFromMnemonic
```

### 2. Check Balance

```bash
znn-cli balance --keyStore my-wallet
```

### 3. Send Tokens

```bash
znn-cli send z1qz... 10 ZNN --keyStore my-wallet
```

## Usage

### Global Flags

```
-u, --url <URL>             WebSocket daemon URL (default: ws://127.0.0.1:35998)
-p, --passphrase <PASS>     Wallet passphrase (prompts if not provided)
-k, --keyStore <NAME>       KeyStore file name
-i, --index <INDEX>         BIP44 account index (default: 0)
-v, --verbose               Enable verbose logging
-h, --help                  Show help information
```

### Command Categories

#### Wallet Commands (6)
```bash
wallet list                                         # List all wallets
wallet createNew                                    # Create new wallet
wallet createFromMnemonic                           # Import from mnemonic
wallet dumpMnemonic                                 # Show mnemonic
wallet deriveAddresses <start> <end>                # Derive addresses
wallet export <filePath>                            # Export wallet
```

#### Query & Transaction Commands (7)
```bash
version                                             # Show version info
balance                                             # Show balances
send <address> <amount> <token>                     # Send tokens
receive <blockHash>                                 # Receive specific block
receiveAll                                          # Receive all pending
unreceived                                          # List pending transactions
unconfirmed                                         # Show unconfirmed blocks
frontierMomentum                                    # Current momentum info
```

#### Plasma Commands (4)
```bash
plasma list [pageIndex] [pageSize]                  # List fusion entries
plasma get                                          # Get plasma info
plasma fuse <address> <amount>                      # Fuse QSR
plasma cancel <id>                                  # Cancel fusion
```

#### Staking Commands (4)
```bash
stake list                                          # List stake entries
stake register <amount> <months>                    # Stake ZNN (1-12 months)
stake revoke <id>                                   # Cancel stake
stake collect                                       # Collect rewards
```

#### Pillar Commands (7)
```bash
pillar list                                         # List all pillars
pillar register <name> <producer> <reward>          # Register pillar
pillar revoke <name>                                # Revoke pillar
pillar delegate <name>                              # Delegate to pillar
pillar undelegate                                   # Remove delegation
pillar collect                                      # Collect rewards
pillar withdrawQsr                                  # Withdraw QSR
```

#### Sentinel Commands (5)
```bash
sentinel list                                       # List sentinels
sentinel register                                   # Register sentinel
sentinel revoke                                     # Revoke sentinel
sentinel collect                                    # Collect rewards
sentinel withdrawQsr                                # Withdraw QSR
```

#### Token Commands (9)
```bash
token list [page] [size]                            # List all tokens
token getByStandard <zts>                           # Get token by ZTS
token getByOwner <address>                          # Get tokens by owner
token issue <name> <symbol> <domain> <total> <max> <decimals> <mint> <burn> <utility>
token mint <zts> <amount> <address>                 # Mint tokens
token burn <zts> <amount>                           # Burn tokens
token transferOwnership <zts> <newOwner>            # Transfer ownership
token disableMint <zts>                             # Disable minting
```

## Configuration

The CLI can be configured via `~/.znn/cli-config.yaml`:

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

## Development

### Running Tests

```bash
make test              # Run all tests
make test-coverage     # Run tests with coverage report
```

### Linting

```bash
make lint              # Run golangci-lint
make security          # Run gosec security scanner
make fmt               # Format code
```

### Build

```bash
make build             # Build binary
make clean             # Clean build artifacts
make all               # Run all checks and build
```

## Project Structure

```
.
├── cmd/               # Command implementations
│   ├── root.go       # Root command
│   ├── wallet/       # Wallet subcommands
│   ├── plasma/       # Plasma subcommands
│   ├── stake/        # Staking subcommands
│   ├── pillar/       # Pillar subcommands
│   ├── sentinel/     # Sentinel subcommands
│   └── token/        # Token subcommands
├── pkg/              # Public packages
│   ├── config/       # Configuration management
│   ├── wallet/       # Wallet operations
│   ├── client/       # RPC client wrapper
│   ├── transaction/  # Transaction helpers
│   └── format/       # Formatting utilities
├── internal/         # Private packages
│   ├── prompt/       # User prompts
│   ├── tui/          # Terminal UI
│   └── validation/   # Input validation
└── main.go           # Entry point
```

## Security

This CLI:
- Uses Argon2 for key derivation (memory-hard)
- Stores wallets in `~/.znn/wallet/` with AES-256-GCM encryption
- Never logs or transmits passphrases
- Validates all user input
- Passes gosec security scanning

**⚠️ Important**: Always backup your mnemonic phrase. Store it securely offline.

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run `make all` to ensure tests pass and code is formatted
5. Submit a pull request

## References

- **Dart CLI**: [github.com/zenon-network/znn_cli_dart](https://github.com/zenon-network/znn_cli_dart)
- **Go SDK**: [github.com/0x3639/znn-sdk-go](https://github.com/0x3639/znn-sdk-go)
- **Zenon Network**: [zenon.network](https://zenon.network)

## License

[MIT License](LICENSE)

## Support

- Issues: [GitHub Issues](https://github.com/0x3639/znn_cli_go/issues)
- Telegram: [Zenon Network](https://t.me/zenonnetwork)

---

**Status**: ✅ **Production Ready** - All core functionality complete! See [ROADMAP.md](ROADMAP.md) for details.
