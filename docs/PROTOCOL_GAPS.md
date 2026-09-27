# SAYANJALI BLOCKCHAIN --- Protocol Gaps

This is a gap register, not a redesign proposal. Missing behavior is
recorded as undefined until an explicit protocol decision is approved.

## Critical compatibility gaps

### G-001 --- Normative protocol specification

**Severity:** Critical\
**Status:** Open\
Rules are currently distributed across Python source/config/tests. A
single normative specification does not yet exist.

### G-002 --- Deterministic cross-language serialization

**Severity:** Critical\
**Status:** Open\
The current implementation relies on Python JSON serialization behavior,
including float formatting and `default=str`. Go must not guess these
rules.

### G-003 --- Protocol test vectors

**Severity:** Critical\
**Status:** Open\
No committed Python/Go compatibility vector suite currently defines
expected hashes, signatures, Merkle roots, difficulty results and
monetary transitions.

### G-004 --- Protocol version vs implementation version

**Severity:** High\
**Status:** Open\
The package and API version strings are not aligned with the current
development checkpoint. A formal protocol-version field is absent from
block/transaction serialization.

## Consensus gaps

### G-005 --- Timestamp consensus bounds

**Severity:** High\
**Status:** Open\
Only strict monotonicity is enforced. Future-time and median-time rules
are undefined.

### G-006 --- Difficulty window semantics

**Severity:** High\
**Status:** Open\
The setting is `10`, but the current implementation uses 10 blocks to
measure 9 intervals. This is existing behavior and must be vectorized
before any change.

### G-007 --- Difficulty representation

**Severity:** High\
**Status:** Open\
Difficulty is an integer count of leading hexadecimal zeros. Work is
`16 ** difficulty`. This is valid current behavior but unusual and must
be preserved exactly.

### G-008 --- Chain-work domain

**Severity:** Medium\
**Status:** Open\
The chain-work formula is simple and deterministic, but its long-term
economic/security meaning has not been formally documented beyond the
current implementation.

## Monetary/state gaps

### G-009 --- Issuance schedule

**Severity:** High\
**Status:** Open\
`50 SYJ` is the current configured block reward and `210,000` is a
reserved halving interval, but the halving rule is not enforced. No
final emission schedule has been approved.

### G-010 --- Genesis allocation

**Severity:** Critical for mainnet\
**Status:** Undefined\
No final genesis distribution/allocation is defined. Do not invent one.

### G-011 --- Transaction nonce

**Severity:** High\
**Status:** Undefined\
No sender nonce exists. This affects replay semantics and future
account/state architecture.

### G-012 --- Fees

**Severity:** Medium\
**Status:** Undefined\
No transaction fee model is currently defined.

### G-013 --- State root

**Severity:** High\
**Status:** Partially addressed by Phase 9.5; consensus activation deferred\

Phase 9.5 implements a deterministic state-root commitment as a
sidecar/checkpoint commitment. The commitment is derived from deterministic
state snapshots and is validated during incremental state processing and
checkpoint recovery.

The current block-header consensus format does **not** contain an active
state-root field. Therefore the Phase 9.5 state root must not be represented
as an already-active consensus field.

Any future activation of a state root in block-header consensus requires an
explicit protocol amendment, compatibility rules, deterministic
cross-language vectors, migration/activation rules, and network-wide
validation.

## Networking gaps

### G-014 --- Production wire protocol
**Severity:** Critical\
**Status:** Resolved / Frozen\
The production `sayanjali-p2p` v1.0 framing, message grammar, request
correlation, handshake, synchronization messages, limits and rejection
registry are frozen in `protocol/p2p-wire-spec.md`. The standalone Go
codec and deterministic wire fixtures are under `internal/p2p/` and
`tests/p2p`.

### G-015 --- Incremental synchronization

**Severity:** High\
**Status:** Open / Implementation pending\
Production synchronization is not implemented in the node yet. The v1.0
wire contract now defines bounded locator-based headers exchange and
block retrieval for the future Phase 6 sync engine.

### G-016 --- Peer reputation

**Severity:** High\
**Status:** Partial\
Failure backoff exists, but persistent reputation/quarantine/ban
semantics are absent.

### G-017 --- Transport security

**Severity:** High\
**Status:** Partial\
Peer authentication is cryptographic at the message layer, but transport
is still HTTP. Production TLS/secure transport semantics are undefined.

### G-018 --- DNS rebinding

**Severity:** High\
**Status:** Partial\
Address validation resolves hostnames once but does not pin the resolved
IP to the connection.

### G-019 --- Sybil/eclipsing strategy

**Severity:** High\
**Status:** Open\
Authentication proves key possession but does not by itself solve Sybil
or eclipse resistance.

## Operational gaps

### G-020 --- CI

**Severity:** High\
**Status:** Open\
The repository contains `.github/workflows/production-validation.yml`; security scanner coverage and CI tool versions are tracked in Phase 9.3.

### G-021 --- Dependency reproducibility

**Severity:** High\
**Status:** Open\
`requirements.txt` uses minimum-version constraints; no lockfile is
present.

### G-022 --- Wallet key exposure

**Severity:** High\
**Status:** Open\
`POST /wallet/create` does not return private-key material. The API contract
must continue to enforce this. Private-key material must not be exposed by a
production-facing node API.

### G-023 --- Documentation drift

**Severity:** Medium\
**Status:** Open\
README and badge counts do not match the verified Phase 3 baseline.

## Phase 9.5 Gap Status

Phase 9.5 addresses a subset of the state-validation and commitment gaps
without declaring the protocol production-ready.

### Addressed or materially advanced

- Deterministic state snapshots and state-root computation.
- Deterministic account ordering.
- Bounded and hardened checkpoint encoding/decoding.
- Truncated-checkpoint and duplicate-account rejection.
- State-root deterministic vectors.
- Merkle Protocol V3 implementation and regression coverage.
- Incremental state validation.
- Periodic state checkpoints and restart/recovery validation.
- Legacy replay versus incremental-state differential testing.
- V2 divergent-branch/reorganization state validation.
- Improved block → checkpoint → durable-tip persistence ordering.

### Still open or deferred

- Active state-root field in the block-header consensus format.
- Merkle V3 activation and network compatibility rules.
- Historical V1 mempool/reorganization semantics.
- Crash/fault-injection proof of fully atomic block/state/tip persistence.
- Go race validation on Android/ARM64.
- Fuzzing.
- Empirical 400-block mining/difficulty regression.
- Live Phase 9.5 runtime verification harness.
- External security audit.
- Full production protocol compatibility and migration process.

**Overall project status: NOT PRODUCTION-READY / NOT MAINNET-READY.**

## Compatibility policy

Every gap is classified as: - **Defined:** directly implemented and
testable; - **Partial:** behavior exists but is not
production-complete; - **Undefined:** no protocol rule exists and
implementation must not invent one; - **Protocol change required:**
changing it would require an ADR and compatibility decision.
