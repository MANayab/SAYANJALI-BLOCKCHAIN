# Build 1 Recovery — C-1 Evidence

## Baseline

Commit: `190f5cac5185e10d56e6f506e127e16efffcd7a7`

Tree: `2443422fd5ccc356defd7739542980c3d472a5c7`

Source ZIP: `SAYANJALI-BLOCKCHAIN-build1-source-190f5ca.zip`

Source ZIP SHA-256: `5e5a39332200e913902644052348c892135c82b8e8731ddb1903d1a51fd1d4d1`

Historical Build 1: **NOT FOUND**

Recovery: **IMPLEMENTED FROM AUTHORITATIVE PHASE 10.1 BASELINE**

## Implementation

Changed source files:

- `internal/chain/chain.go`
- `internal/chain/v2.go`
- `internal/clock/clock.go`
- `internal/node/api.go`
- `internal/node/node.go`
- `internal/p2pnode/network.go`
- `pkg/protocol/v3params.go`
- `internal/chain/c1_time_test.go`
- `internal/node/c1_mining_test.go`
- `cmd/build1_c1/main.go`
- `docs/protocol/V3_DECISIONS.md`
- `build/CHANGES.md`
- `build/VALIDATE.sh`
- `build/EVIDENCE.md`
- `build/changes.patch`
- `build/SHA256SUMS`

C-1 behavior:

- V3 explicitly selectable through existing protocol version architecture.
- V3 MTP lower bound is `MTP + 1`.
- V3 first-ingress future bound is `Clock.Now() + 300s`.
- V3 has no parent+120 future ceiling.
- V1/V2 retain their existing parent+120 behavior.
- Historical replay does not consult wall clock.
- V3 mining uses `max(Clock.Now(), MTP+1)`.
- V3 mining refuses when peer median is unavailable.
- V3 mining refuses when legitimate peer median skew exceeds 60 seconds.
- No P2P wire timestamp field was fabricated or added.

## Required toolchain

Command:

`go version`

Actual output:

`go version go1.27.0 android/arm64`

Required:

`go1.27.1`

Toolchain attempt:

The environment attempted to obtain the required Go `1.27.1` toolchain, but
the required toolchain was unavailable in the execution environment.

Required Go 1.27.1 validation: **NOT EXECUTED — Go 1.27.1 unavailable**.

The installed Go version used for the successful non-race validation was:

`go version go1.27.0 android/arm64`

## Validation commands

### `gofmt -l .`

Executed.

Output: **empty**.

Result: **PASS**.

### `go vet ./...`

Executed with Go 1.27.0 on android/arm64.

Result: **PASS**.

Exit code: `0`.

### `go build ./...`

Executed with Go 1.27.0 on android/arm64.

Result: **PASS**.

Exit code: `0`.

### `go test ./...`

Executed with Go 1.27.0 on android/arm64.

Result: **PASS**.

All repository packages passed.

### `go test -race ./...`

Attempted with Go 1.27.0 on android/arm64.

Result: **NOT EXECUTED — unsupported on android/arm64**.

Actual output:

`-race is not supported on android/arm64`

Exit code: `2`.

### C-1 tests

Chain command:

`go test ./internal/chain ./internal/node -run 'C1|C-1' -count=1 -v`

The C-1 chain tests passed after the injected-clock correction.

The V3 mining safety test was separately executed:

`go test ./internal/node -run 'TestC1V3MiningClockSafety' -count=1 -v`

Result: **PASS**.

### Compatibility tests

The full repository test suite included the existing compatibility packages and passed under Go 1.27.0 android/arm64.

### Restart/reload tests

The C-1 replay/restart coverage executed successfully as part of the C-1 chain tests.

Replay was tested with clocks decades ahead and decades behind.

Result: **PASS**.

### Required Go 1.27.1

Required exact toolchain:

`go1.27.1`

Actual available toolchain:

`go1.27.0 android/arm64`

Go 1.27.1 was not available in the current Android/ARM64 environment.

Required Go 1.27.1 validation: **NOT EXECUTED**.

## Auxiliary validation


This auxiliary attempt does not count as required validation and no production `go.mod` was modified.

## SHA256SUMS

`build/SHA256SUMS` contains hashes calculated from the final working tree for selected review artifacts. The checksum file deliberately excludes itself and the final ZIP to avoid self-referential hashes.

The final ZIP SHA-256 is calculated after packaging and is reported in the delivery report rather than embedded inside the ZIP.

## Scope status

C-1: **DONE-NOT-VERIFIED**

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

## Review limitations

The deterministic harness records seed, block count, accepted/rejected blocks, elapsed time, difficulty observations, cumulative work, final height/tip/state root, replay result, and restart result.

The implementation was not declared production-ready, mainnet-ready, public-testnet-ready, or fully audited. Required Go 1.27.1 compilation, unit, race, and vet evidence could not be produced in the current execution environment.

## `build/VALIDATE.sh` execution

The validation script intentionally requires exact Go `1.27.1`.

With the currently installed Go `1.27.0 android/arm64`, the script is
expected to fail its exact-toolchain gate before running the subsequent
commands.

This failure is not represented as a successful validation.


## Deterministic C-1 harness evidence

The deterministic Build 1 harness was executed twice with seed `20261001`.

### Run 1

```text
seed=20261001
blocks=64
accepted=19
rejected=45
elapsed=38.602743058s
difficulty_observations=[4 4 4 4 4 4 4 4 4 4 4 4 4 4 4 4 4 4 4]
cumulative_work=1245185
final_height=19
final_tip=00007df91e02a3a89ab1f445f2e024d53d5ab27c732f6303733cfd00962a4f2d
state_root=77627aa058556e65dac3e7deed04d9cf7fce3e95f3208ef59ceabfe995cf69d3
replay_result=PASS
restart_result=PASS
```

### Run 2

```text
seed=20261001
blocks=64
accepted=19
rejected=45
elapsed=33.593833633s
difficulty_observations=[4 4 4 4 4 4 4 4 4 4 4 4 4 4 4 4 4 4 4]
cumulative_work=1245185
final_height=19
final_tip=00007df91e02a3a89ab1f445f2e024d53d5ab27c732f6303733cfd00962a4f2d
state_root=77627aa058556e65dac3e7deed04d9cf7fce3e95f3208ef59ceabfe995cf69d3
replay_result=PASS
restart_result=PASS
```

The two runs were compared after removing only the nondeterministic `elapsed` line.

Result:

**PASS – deterministic consensus/result fields matched exactly.**

The harness also independently verified:

- replay with a clock ten years ahead;
- restart/reload with a clock ten years behind.

Both returned `PASS`.

## C-1 implementation correction discovered during validation

The initial C-1 implementation accepted an injected clock in
`openWithProtocolAndClock` but initialized the Chain with a real clock
instead of the supplied clock.

The defect was:

```diff
+ clock:           clock.RealClock{},
+ clock:           clk,
```

This was discovered by the `+301 seconds` future-boundary test:
the fixed-clock `+301s` case was initially accepted.

After correcting the clock injection, the C-1 future-boundary tests passed:

- `Clock.Now() + 300s`: accepted
- `Clock.Now() + 301s`: rejected

No unrelated protocol or consensus change was introduced for this
correction.

## Validation status

The repository was successfully formatted, vetted, built, and tested
with the installed Go `1.27.0 android/arm64` toolchain.

The required exact Go `1.27.1` toolchain was unavailable in the current
environment, so exact-toolchain validation remains **NOT EXECUTED**.

`go test -race ./...` is **NOT EXECUTED** because the installed
Android/ARM64 environment reports that `-race` is unsupported.

Accordingly, the final C-1 status remains:

**DONE-NOT-VERIFIED**

The implementation is not being represented as production-ready, mainnet-ready, public-testnet-ready, or fully audited.
