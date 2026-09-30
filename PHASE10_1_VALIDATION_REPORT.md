# Phase 10.1 Validation Report — 400-Block Consensus & Difficulty Regression

## Executive result

**Engineering result:** the Phase 10.1 harness was implemented and the real Go consensus/chain/state path successfully accepted 400 sequential blocks in an isolated validation build. The deterministic experiment reproduced all 400 consensus records across independent runs.

**Production-build gate:** **NOT EXECUTED**. The supplied environment has Go 1.23.2 only, while `go.mod` requires Go 1.27 and `toolchain go1.27.1`; the environment could not download the required toolchain or the pinned secp256k1 dependency.

**Production status:** **NOT PRODUCTION-READY / NOT MAINNET-READY**.

The 400-block result below is therefore strong consensus-path regression evidence, but it is not a full production-toolchain or cryptographic-dependency validation claim.

## Provenance

- Authoritative baseline: `ea26fe1db68f2e45ad7051cbedae655de36b0ca6`
- Commit: `security: complete Phase 9.5 state validation hardening`
- Tag: `v0.9.5-security-state-hardening`
- Uploaded archive SHA-256: `6db49cdd86f6283f07178b3032401f3830c47e5a9bad61e9f3da626f3fb070c`
- Complete source reconciliation: 305/305 files matched the authoritative `ea26fe1` tree.

## Implementation added

- `cmd/phase10_1/main.go`
- `cmd/phase10_1/main_test.go`
- `scripts/phase10_1/run_400_block.sh`
- `scripts/phase10_1/compare_runs.py`
- `docs/consensus/phase10-400-block-regression.md`
- `docs/security/phase10-runtime-validation.md`
- `PHASE10_PROVENANCE.md`
- `PHASE10_1_VALIDATION_REPORT.md`
- `PHASE10_1_MANIFEST.md`

The harness does not modify production consensus parameters.

## Actual production code path exercised by the experiment

`storage.Open` → `chain.Open` → `block.New` → production `block.HashHeader` → production difficulty predicate → `chain.Accept` → production incremental state validation/checkpointing → persistent storage close/reopen → `chain.Open` → production chain/state reconstruction.

The final accepted blocks were revalidated by the production `chain.Accept` path, which independently checks previous-hash linkage, timestamps, required difficulty, PoW, transaction structure, state transition, Merkle root, and incremental state commitment.

The harness's only acceleration is nonce search: it uses the production `block.HashHeader` function in deterministic fixed nonce batches across CPU workers and selects the lowest valid nonce. `chain.Accept` then revalidates the resulting header through the production consensus implementation. This acceleration is harness-only and is not a production mining implementation.

## Deterministic test configuration

The 400-block experiment uses one unsigned coinbase transaction per block and a deterministic valid receiver. Timestamp steps rotate by 10-block epoch:

1. `+120s` — slow epoch, drives difficulty down from 4 to 3.
2. `+6s` — fast epoch, drives difficulty up from 3 to 4.
3. `+30s` — stable epoch, maintains difficulty.

This pattern repeats for 40 retarget boundaries. Every step is within the production deterministic future-time bound of 120 seconds.

No production difficulty or timestamp constant was changed.

## Source-defined consensus parameters

Read directly from the authoritative source:

| Parameter | Value |
|---|---:|
| Target block time | 30 seconds |
| Difficulty interval | 10 blocks |
| Minimum difficulty | 1 |
| Maximum difficulty | 32 |
| Maximum adjustment factor | 4 |
| Timestamp future bound | 120 seconds |
| Median-time-past window | 11 blocks |
| First mineable difficulty | 4 |

## 400-block empirical result

Final run ID: `phase10.1-stub-consensus-final`

- Blocks requested: 400
- Non-genesis blocks accepted: **400**
- Retarget boundaries crossed: **40**
- Final height: **400**
- Final difficulty: **3**
- Final tip: `000746e976bee97ce5c35202da451731497560e581783c977de9c110b745cf05`
- Final cumulative work: `18165761`
- Final state root: `7128583826764b3b691fb557bd90e91a39d7261b1a7ffd6422224c7b795e680e`
- Restart height: **400**
- Restart tip: same as final tip
- Restart state root: same as final state root
- Restart cumulative work: same as final cumulative work
- Per-block validation result: **PASS** for all 400 records
- Deterministic consensus fields reproduced across independent runs: **PASS**

### Retarget behavior observed

- Height 10: difficulty 4 → 3
- Height 20: difficulty 3 → 4
- Height 30: difficulty 4 → 4
- Height 40: difficulty 4 → 3
- The pattern continues through height 400.

The observed schedule exercises both downward and upward retarget behavior and the stable/no-change case.

## State-root observations

- State transition succeeded for every accepted block.
- State root was non-empty for every accepted block.
- State root progressed as the mining-issued supply and receiver balance changed.
- Final state root survived store close/reopen unchanged.
- Final restarted chain passed full `chain.ValidateChain` validation.
- The Phase 9.5 state root remains a deterministic sidecar/checkpoint commitment; Phase 10.1 did not promote it into the block header consensus format.

## Merkle observations

- Current active block construction uses the existing production Merkle behavior.
- Every accepted block's recorded Merkle root was checked against its canonical transaction hash.
- Merkle V3 was not activated by Phase 10.1 and no historical V2/V3 compatibility claim is made.

## Determinism

Two independent 400-block runs were compared using:

`scripts/phase10_1/compare_runs.py`

Consensus fields were identical for all 400 records:

- block hashes
- previous hashes
- timestamps
- timestamp deltas
- difficulty
- required target prefix
- cumulative work
- transaction count
- state roots
- Merkle roots
- validation results
- expected difficulty

Mining and validation durations are intentionally excluded from the deterministic comparison.

## Performance evidence

Final isolated consensus-path run:

- total mining/production time: `2.93045451s`
- average mining/production time: `7.326136ms/block`
- total validation time: `85.901703ms`
- average validation time: `214.754µs/block`

These are absolute measurements from the isolated validation environment and are **not** a baseline comparison or production benchmark.

## Environment limitation

The repository requires:

```text
go 1.27
toolchain go1.27.1
github.com/decred/dcrd/dcrec/secp256k1/v4 v4.4.1
```

Available local Go:

```text
go1.23.2 linux/amd64
```

`go` attempted to download Go 1.27.1 and failed because this execution environment cannot resolve/reach `proxy.golang.org`. The pinned secp256k1 module was also not cached.

For the empirical consensus-path run only, a temporary local crypto stub was used in a disposable copy so the production consensus, block hashing, chain acceptance, storage, and state code could execute. The stub was **not** included in the Phase 10.1 delivery archive and no cryptographic validation result is inferred from it.

## Validation status matrix

| Validation | Result | Evidence / limitation |
|---|---|---|
| 400-block consensus-path experiment | PASS | 400/400 accepted; 40 retarget boundaries; restart validation passed |
| Deterministic two-run consensus comparison | PASS | 400/400 records identical |
| Harness unit tests in isolated Go build | PASS | `go test ./cmd/phase10_1` with temporary stub |
| Production Go build | NOT EXECUTED | Required Go 1.27.1 unavailable |
| Production Go vet | NOT EXECUTED | Required Go 1.27.1 unavailable |
| Production Go tests | NOT EXECUTED | Required Go 1.27.1 unavailable |
| Production Go race tests | NOT EXECUTED | Required Go 1.27.1 unavailable |
| staticcheck | NOT EXECUTED | Tool unavailable |
| govulncheck | NOT EXECUTED | Tool unavailable |
| gosec | NOT EXECUTED | Tool unavailable |
| Python compileall | PASS | Completed successfully |
| Python pytest | PASS | 262 passed |
| Focused consensus/mining/state-root pytest | PASS | 35 passed |
| Bandit | NOT EXECUTED | Tool unavailable |
| pip-audit | NOT EXECUTED | Tool unavailable |
| gofmt | PASS | `gofmt -l cmd internal pkg tests` returned no paths |

## Failure evidence retained during development

Two harness-only failures were encountered and corrected:

1. Merkle assertion initially used the caller's pre-normalized transaction rather than the canonical transaction copy created by `block.New`.
2. Retarget measurement initially sampled the post-acceptance window rather than the exact pre-block window used by production validation.

Neither was a consensus implementation failure. Both were corrected in the final harness and the final 400-block run passed.

## Remaining blockers

Phase 10.1 does not close the outstanding production gates. These remain unresolved unless independently demonstrated:

- target-architecture Go race testing
- fuzzing
- production-toolchain build/test/vet validation
- staticcheck
- govulncheck
- gosec
- Bandit/pip-audit execution in the target environment
- live runtime/P2P/API harness validation
- crash/fault-injection persistence validation
- independent security audit
- remaining M-03/M-07 issues
- Merkle V3 activation/compatibility review
- state-root protocol activation/compatibility review
- other Phase 9.5 open findings

## Status decision

**Phase 10.1: IMPLEMENTED / EMPIRICALLY VALIDATED ON THE CONSENSUS PATH, WITH PRODUCTION-TOOLCHAIN VALIDATION OUTSTANDING.**

**Overall project: NOT PRODUCTION-READY / NOT MAINNET-READY.**
