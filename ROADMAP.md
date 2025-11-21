# Zenon Go CLI - Development Roadmap

This roadmap tracks the implementation progress of the Zenon Go CLI wallet, replicating all features from the Dart CLI reference.

## Legend
- ✅ Completed
- 🚧 In Progress
- ⏳ Planned
- 🔄 Under Review

---

## Phase 1: Project Setup & Development Infrastructure ✅

### 1.1 Initialize Go Project ✅
- [x] Initialize go module: `github.com/0x3639/znn_cli_go`
- [x] Create project directory structure
- [x] Set up main.go entry point

### 1.2 Reference Code Setup ✅
- [x] Create `reference/` directory
- [x] Copy Dart CLI to `./reference/znn_cli_dart/`
- [x] Add to `.gitignore`

### 1.3 Quality Tooling ✅
- [x] Configure golangci-lint (`.golangci.yml`)
- [x] Configure gosec (`.gosec.yml`)
- [x] Create Makefile with targets:
  - [x] `make build`
  - [x] `make install`
  - [x] `make test`
  - [x] `make lint`
  - [x] `make security`
  - [x] `make fmt`
  - [x] `make clean`

### 1.4 Documentation Standards ⏳
- [ ] Create README.md
- [ ] Create CLAUDE.md
- [ ] Set up LICENSE file

---

## Phase 2: Core Infrastructure ⏳

### 2.1 Install Dependencies
- [ ] cobra + viper
- [ ] bubbletea + tview + lipgloss
- [ ] fatih/color
- [ ] manifoldco/promptui
- [ ] SDK: github.com/0x3639/znn-sdk-go

### 2.2 Configuration Package (`pkg/config/`)
- [ ] Config struct with viper integration
- [ ] Default values (node URL, wallet path)
- [ ] Load from `~/.znn/cli-config.yaml`
- [ ] godoc documentation

### 2.3 Wallet Package (`pkg/wallet/`)
- [ ] Wrapper around SDK KeyStoreManager
- [ ] Wallet loading with passphrase prompts
- [ ] Address derivation helpers
- [ ] godoc documentation

### 2.4 Client Package (`pkg/client/`)
- [ ] RPC client wrapper
- [ ] Connection management
- [ ] Auto-reconnection logic
- [ ] godoc documentation

### 2.5 Transaction Package (`pkg/transaction/`)
- [ ] Autofill helper (height, previousHash, momentum)
- [ ] Plasma/PoW decision logic
- [ ] Sign and publish helpers
- [ ] godoc documentation

### 2.6 Format Package (`pkg/format/`)
- [ ] Amount formatting
- [ ] Duration formatting
- [ ] Token standard parsing
- [ ] Color helpers
- [ ] godoc documentation

### 2.7 Root Command (`cmd/root.go`)
- [ ] Cobra root command setup
- [ ] Global flags (--url, --passphrase, --keyStore, --index, --verbose)
- [ ] Command registration
- [ ] Version command

---

## Phase 3: Wallet Commands (6 commands) ⏳

Reference: `reference/znn_cli_dart/lib/init_znn.dart` and `cli_handler.dart`

- [ ] `wallet.list` - List all keyStore files
- [ ] `wallet.createNew` - Create wallet with BIP39 mnemonic
- [ ] `wallet.createFromMnemonic` - Import from mnemonic
- [ ] `wallet.dumpMnemonic` - Display mnemonic for backup
- [ ] `wallet.deriveAddresses` - Show addresses from BIP44 derivation
- [ ] `wallet.export` - Export keyStore to file

---

## Phase 4: Basic Transaction Commands (9 commands) ⏳

Reference: `reference/znn_cli_dart/lib/cli_handler.dart`

- [ ] `version` - Show CLI/SDK/daemon versions (Line 18-27)
- [ ] `balance` - Display ZNN/QSR/ZTS balances (Line 60-82)
- [ ] `send` - Send tokens with optional message (Line 84-160)
- [ ] `receive` - Receive specific block by hash (Line 162-199)
- [ ] `receiveAll` - Batch receive all pending (Line 201-245)
- [ ] `unreceived` - List pending transactions (Line 247-272)
- [ ] `unconfirmed` - Show unconfirmed blocks (Line 274-295)
- [ ] `frontierMomentum` - Get current momentum info (Line 297-310)
- [ ] `autoreceive` - Auto-receive daemon mode (Line 1045-1161)

---

## Phase 5: Plasma Commands (4 commands) ⏳

Reference: `cli_handler.dart` lines 493-589

- [ ] `plasma.list` - List fusion entries (Line 493-522)
- [ ] `plasma.get` - Get plasma info (Line 524-541)
- [ ] `plasma.fuse` - Fuse QSR for beneficiary (Line 543-565)
- [ ] `plasma.cancel` - Cancel fusion by ID (Line 567-589)

---

## Phase 6: Staking Commands (4 commands) ⏳

Reference: `cli_handler.dart` lines 591-700

- [ ] `stake.list` - List stake entries (Line 591-625)
- [ ] `stake.register` - Stake ZNN 1-12 months (Line 627-662)
- [ ] `stake.revoke` - Cancel expired stake (Line 664-683)
- [ ] `stake.collect` - Collect staking rewards (Line 685-700)

---

## Phase 7: Pillar Commands (7 commands) ⏳

Reference: `cli_handler.dart` lines 312-491

- [ ] `pillar.list` - List all pillars (Line 312-366)
- [ ] `pillar.register` - Register new pillar (Line 368-405)
- [ ] `pillar.revoke` - Revoke pillar (Line 407-423)
- [ ] `pillar.delegate` - Delegate to pillar (Line 425-443)
- [ ] `pillar.undelegate` - Remove delegation (Line 445-460)
- [ ] `pillar.collect` - Collect pillar rewards (Line 462-476)
- [ ] `pillar.withdrawQsr` - Withdraw deposited QSR (Line 478-491)

---

## Phase 8: Sentinel Commands (5 commands) ⏳

Reference: `cli_handler.dart` lines 702-791

- [ ] `sentinel.list` - List sentinel info (Line 702-732)
- [ ] `sentinel.register` - Register sentinel (Line 734-752)
- [ ] `sentinel.revoke` - Revoke sentinel (Line 754-768)
- [ ] `sentinel.collect` - Collect rewards (Line 770-781)
- [ ] `sentinel.withdrawQsr` - Withdraw QSR (Line 783-791)

---

## Phase 9: Token Commands (9 commands) ⏳

Reference: `cli_handler.dart` lines 793-1043

- [ ] `token.list` - List all ZTS tokens (Line 793-830)
- [ ] `token.getByStandard` - Get token by ZTS (Line 832-850)
- [ ] `token.getByOwner` - Get tokens by owner (Line 852-876)
- [ ] `token.issue` - Issue new token (Line 878-951)
- [ ] `token.mint` - Mint additional supply (Line 953-979)
- [ ] `token.burn` - Burn tokens (Line 981-1003)
- [ ] `token.transferOwnership` - Transfer ownership (Line 1005-1024)
- [ ] `token.disableMint` - Disable minting (Line 1026-1043)

---

## Phase 10: TUI Interface ⏳

- [ ] Main menu with command categories
- [ ] Interactive forms for complex operations
- [ ] Real-time balance dashboard
- [ ] Transaction history viewer
- [ ] Wallet selector
- [ ] `autoreceive` TUI mode with live updates
- [ ] godoc documentation

---

## Phase 11: Testing & Quality ⏳

### 11.1 Linting
- [ ] Run `make lint` - fix all issues
- [ ] Configure CI to run linters

### 11.2 Security
- [ ] Run `make security` - address findings
- [ ] Review crypto operations
- [ ] Validate input sanitization

### 11.3 Testing
- [ ] Unit tests for pkg/ packages
- [ ] Integration tests with testnet
- [ ] Table-driven tests for formatting
- [ ] Error case coverage
- [ ] Target: >80% coverage

### 11.4 Documentation Review
- [ ] Verify all godoc complete
- [ ] Generate docs: `godoc -http=:6060`
- [ ] README completeness
- [ ] CLAUDE.md accuracy

---

## Phase 12: Release Preparation ⏳

- [ ] Version tagging (v0.1.0)
- [ ] GitHub releases with binaries
- [ ] Installation instructions
- [ ] Shell completion scripts (bash/zsh/fish)
- [ ] Final documentation review
- [ ] Announce release

---

## Summary Stats

**Total Commands to Implement**: 47
- Wallet: 6 commands
- Basic Operations: 9 commands
- Plasma: 4 commands
- Staking: 4 commands
- Pillar: 7 commands
- Sentinel: 5 commands
- Token: 9 commands
- TUI: Interactive mode

**Current Progress**:
- Phase 1: 100% ✅
- Overall: ~8%
