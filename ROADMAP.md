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

### 1.4 Documentation Standards ✅
- [x] Create README.md
- [x] Create CLAUDE.md
- [ ] Set up LICENSE file

---

## Phase 2: Core Infrastructure ✅

### 2.1 Install Dependencies ✅
- [x] cobra + viper
- [x] fatih/color
- [x] golang.org/x/term (for secure password input)
- [x] SDK: github.com/0x3639/znn-sdk-go
- [ ] bubbletea + tview + lipgloss (deferred to Phase 11 - TUI)

### 2.2 Configuration Package (`pkg/config/`) ✅
- [x] Config struct with viper integration
- [x] Default values (node URL, wallet path)
- [x] Load from `~/.znn/cli-config.yaml`
- [x] godoc documentation

### 2.3 Wallet Package (`pkg/wallet/`) ✅
- [x] Wrapper around SDK KeyStoreManager
- [x] Wallet loading with passphrase prompts
- [x] Address derivation helpers
- [x] godoc documentation

### 2.4 Client Package (`pkg/client/`) ✅
- [x] RPC client wrapper
- [x] Connection management
- [x] Auto-reconnection logic
- [x] godoc documentation

### 2.5 Transaction Package (`pkg/transaction/`) ✅
- [x] Autofill helper (height, previousHash, momentum)
- [x] Plasma/PoW decision logic
- [x] Sign and publish helpers
- [x] godoc documentation

### 2.6 Format Package (`pkg/format/`) ✅
- [x] Amount formatting
- [x] Duration formatting
- [x] Token standard parsing
- [x] Color helpers
- [x] godoc documentation

### 2.7 Prompt Package (`internal/prompt/`) ✅
- [x] Secure password input (echo disabled)
- [x] Password confirmation
- [x] Yes/no confirmations
- [x] godoc documentation

### 2.8 Root Command (`cmd/root.go`) ✅
- [x] Cobra root command setup
- [x] Global flags (--url, --passphrase, --keyStore, --index, --verbose, --config)
- [x] Command registration
- [x] Version command

---

## Phase 3: Wallet Commands (6 commands) ✅

Reference: `reference/znn_cli_dart/lib/init_znn.dart` and `cli_handler.dart`

- [x] `wallet list` - List all keyStore files
- [x] `wallet createNew` - Create wallet with BIP39 mnemonic
- [x] `wallet createFromMnemonic` - Import from mnemonic
- [x] `wallet dumpMnemonic` - Display mnemonic for backup
- [x] `wallet deriveAddresses` - Show addresses from BIP44 derivation
- [x] `wallet export` - Export keyStore to file

---

## Phase 4: Query Commands (4 commands) ✅

Reference: `reference/znn_cli_dart/lib/cli_handler.dart`

- [x] `balance` - Display ZNN/QSR/ZTS balances (Line 60-82)
- [x] `unreceived` - List pending transactions (Line 247-272)
- [x] `unconfirmed` - Show unconfirmed blocks (Line 274-295)
- [x] `frontierMomentum` - Get current momentum info (Line 297-310)

---

## Phase 5: Transaction Commands (3 commands) ✅

Reference: `reference/znn_cli_dart/lib/cli_handler.dart`

- [x] `send` - Send tokens (Line 84-160)
- [x] `receive` - Receive specific block by hash (Line 162-199)
- [x] `receiveAll` - Batch receive all pending (Line 201-245)

Note: `autoreceive` daemon mode is deferred to Phase 12

---

## Phase 6: Plasma Commands (4 commands) ✅

Reference: `cli_handler.dart` lines 493-589

- [x] `plasma.list` - List fusion entries (Line 493-522)
- [x] `plasma.get` - Get plasma info (Line 524-541)
- [x] `plasma.fuse` - Fuse QSR for beneficiary (Line 543-565)
- [x] `plasma.cancel` - Cancel fusion by ID (Line 567-589)

---

## Phase 7: Staking Commands (4 commands) ✅

Reference: `cli_handler.dart` lines 591-700

- [x] `stake.list` - List stake entries (Line 591-625)
- [x] `stake.register` - Stake ZNN 1-12 months (Line 627-662)
- [x] `stake.revoke` - Cancel expired stake (Line 664-683)
- [x] `stake.collect` - Collect staking rewards (Line 685-700)

---

## Phase 8: Pillar Commands (7 commands) ✅

Reference: `cli_handler.dart` lines 312-491

- [x] `pillar.list` - List all pillars (Line 312-366)
- [x] `pillar.register` - Register new pillar (Line 368-405)
- [x] `pillar.revoke` - Revoke pillar (Line 407-423)
- [x] `pillar.delegate` - Delegate to pillar (Line 425-443)
- [x] `pillar.undelegate` - Remove delegation (Line 445-460)
- [x] `pillar.collect` - Collect pillar rewards (Line 462-476)
- [x] `pillar.withdrawQsr` - Withdraw deposited QSR (Line 478-491)

---

## Phase 9: Sentinel Commands (5 commands) ✅

Reference: `cli_handler.dart` lines 702-791

- [x] `sentinel.list` - List sentinel info (Line 702-732)
- [x] `sentinel.register` - Register sentinel (Line 734-752)
- [x] `sentinel.revoke` - Revoke sentinel (Line 754-768)
- [x] `sentinel.collect` - Collect rewards (Line 770-781)
- [x] `sentinel.withdrawQsr` - Withdraw QSR (Line 783-791)

---

## Phase 10: Token Commands (9 commands) ✅

Reference: `cli_handler.dart` lines 793-1043

- [x] `token.list` - List all ZTS tokens (Line 793-830)
- [x] `token.getByStandard` - Get token by ZTS (Line 832-850)
- [x] `token.getByOwner` - Get tokens by owner (Line 852-876)
- [x] `token.issue` - Issue new token (Line 878-951)
- [x] `token.mint` - Mint additional supply (Line 953-979)
- [x] `token.burn` - Burn tokens (Line 981-1003)
- [x] `token.transferOwnership` - Transfer ownership (Line 1005-1024)
- [x] `token.disableMint` - Disable minting (Line 1026-1043)

---

## Phase 11: TUI Interface ⏳

- [ ] Main menu with command categories
- [ ] Interactive forms for complex operations
- [ ] Real-time balance dashboard
- [ ] Transaction history viewer
- [ ] Wallet selector
- [ ] `autoreceive` TUI mode with live updates
- [ ] godoc documentation

---

## Phase 12: Testing & Quality ⏳

### 12.1 Linting
- [ ] Run `make lint` - fix all issues
- [ ] Configure CI to run linters

### 12.2 Security
- [ ] Run `make security` - address findings
- [ ] Review crypto operations
- [ ] Validate input sanitization

### 12.3 Testing
- [ ] Unit tests for pkg/ packages
- [ ] Integration tests with testnet
- [ ] Table-driven tests for formatting
- [ ] Error case coverage
- [ ] Target: >80% coverage

### 12.4 Documentation Review
- [ ] Verify all godoc complete
- [ ] Generate docs: `godoc -http=:6060`
- [ ] README completeness
- [ ] CLAUDE.md accuracy

### 12.5 Auto-receive Daemon
- [ ] Implement `autoreceive` daemon mode (deferred from Phase 5)
- [ ] WebSocket subscription to account blocks
- [ ] Auto-receive incoming transactions
- [ ] TUI mode with live updates

---

## Phase 13: Release Preparation ⏳

- [ ] Version tagging (v0.1.0)
- [ ] GitHub releases with binaries
- [ ] Installation instructions
- [ ] Shell completion scripts (bash/zsh/fish)
- [ ] Final documentation review
- [ ] Announce release

---

## Summary Stats

**Total Commands to Implement**: 44
- Wallet: 6 commands ✅
- Query: 4 commands ✅
- Transaction: 3 commands ✅
- Plasma: 4 commands ✅
- Staking: 4 commands ✅
- Pillar: 7 commands ✅
- Sentinel: 5 commands ✅
- Token: 9 commands ✅
- TUI: Interactive mode ⏳ (deferred to Phase 11)
- autoreceive: Daemon mode ⏳ (deferred to Phase 12)

**Current Progress**:
- Phase 1: Project Setup ✅ (100%)
- Phase 2: Core Infrastructure ✅ (100%)
- Phase 3: Wallet Commands ✅ (6/6 commands)
- Phase 4: Query Commands ✅ (4/4 commands)
- Phase 5: Transaction Commands ✅ (3/3 commands)
- Phase 6: Plasma Commands ✅ (4/4 commands)
- Phase 7: Staking Commands ✅ (4/4 commands)
- Phase 8: Pillar Commands ✅ (7/7 commands)
- Phase 9: Sentinel Commands ✅ (5/5 commands)
- Phase 10: Token Commands ✅ (9/9 commands)
- Overall: **95%** (42/44 commands implemented)
