# Build 1 Recovery — C-1 Changes

## 1. Recovery identity

Build 1 did not exist historically in the authoritative Git history. This is a recovery implementation created from the verified Phase 10.1 mainline state.

Baseline commit:

`190f5cac5185e10d56e6f506e127e16efffcd7a7`

Baseline tree:

`2443422fd5ccc356defd7739542980c3d472a5c7`

Baseline source ZIP:

`SAYANJALI-BLOCKCHAIN-build1-source-190f5ca.zip`

Baseline source SHA-256:

`5e5a39332200e913902644052348c892135c82b8e8731ddb1903d1a51fd1d4d1`

## 2. Implementation

Implemented C-1 only:

- explicit V3 protocol selection through the existing Chain/Node architecture;
- centralized `pkg/protocol/v3params.go` parameters;
- injectable `internal/clock` abstraction;
- V3 deterministic `MTP + 1` lower bound;
- V3 first-ingress `Clock.Now() + 300s` future bound;
- V3 removal of the legacy parent+120 ceiling while preserving V1/V2;
- V3 mining timestamp `max(now, MTP+1)`;
- operational V3 mining clock safety with a 60-second peer-median threshold;
- explicit refusal when peer median is unavailable;
- no P2P wire-protocol modification;
- replay/restart tests using clocks decades ahead and behind;
- seeded adversarial timestamp tests;
- V1/V2 timestamp compatibility tests.
- deterministic `cmd/build1_c1` harness with seed, acceptance/rejection counts, elapsed time, difficulty observations, cumulative work, tip/state root, replay, and restart results.

## 3. Ingress/replay boundary

The V3 wall-clock check is performed in `Chain.Accept` after the parent chain is resolved and before deterministic block/state validation and persistence. Both direct network block reception and synchronization reach this same acceptance boundary.

Historical startup/replay continues through chain construction and deterministic validation without invoking the ingress clock check.

## 4. Peer clock limitation

The baseline P2P protocol has no legitimate peer-time observation mechanism. No timestamp field or synthetic observation was added. V3 mining conservatively refuses when the peer median is unavailable. A supplied `PeerMedian` implementation can be used by future network-time work without changing C-1 consensus semantics.

## 5. Tests

Added dedicated C-1 tests covering:

- MTP equal/lower/+1 boundaries;
- +300/+301 ingress boundary;
- V3 timestamps beyond legacy +120;
- +1/+30/+120 progression;
- seeded deterministic timestamp generation;
- seeded ingress future outcomes;
- decades-ahead and decades-behind replay;
- restart/reload/replay state identity;
- V1/V2 timestamp compatibility;
- V3 mining timestamp and clock-safety states.

The existing Phase 10.1 harness was preserved.

## 6. Validation limitation

Actual installed Go:

`go version go1.27.0 android/arm64`

The required Go 1.27.1 toolchain was unavailable in the execution environment. Therefore exact Go 1.27.1 validation is recorded as NOT EXECUTED rather than fabricated.

Validation actually executed with Go 1.27.0 android/arm64:

- `gofmt -l .` — PASS
- `go vet ./...` — PASS
- `go build ./...` — PASS
- `go test ./...` — PASS
- C-1 deterministic tests — PASS
- `go test -race ./...` — NOT EXECUTED because race is unsupported on android/arm64
- `./build/VALIDATE.sh` — stops at the required Go 1.27.1 gate, as designed.

## 7. Explicitly deferred findings

C-2 OPEN
H-1 OPEN
H-2 OPEN
H-3 OPEN
H-4 OPEN
H-5 OPEN
H-6 OPEN
H-7 OPEN
H-8 OPEN
M-1 OPEN
M-2 OPEN
M-3 OPEN
M-4 OPEN
M-5 OPEN
M-6 OPEN
M-7 OPEN
M-8 OPEN
L-1 OPEN
L-2 OPEN
L-3 OPEN

## 8. Quality boundary

This artifact does not claim production-ready, mainnet-ready, public-testnet-ready, or fully audited status.
