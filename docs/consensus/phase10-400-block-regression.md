# Phase 10.1 — 400-Block Consensus & Difficulty Regression

## Purpose

This document defines the Phase 10.1 empirical regression harness. It exercises the production Go chain path for 400 sequential non-genesis blocks and records machine-readable evidence for block linkage, timestamps, difficulty, proof of work, cumulative work, state commitments, Merkle roots, persistence reload, and deterministic consensus behavior.

## Baseline

- Authoritative Phase 10 baseline: `ea26fe1db68f2e45ad7051cbedae655de36b0ca6`
- Phase 9.5 tag: `v0.9.5-security-state-hardening`
- Uploaded inspection archive SHA-256: `6db49cdd86f6283f07178b3032401f3830c47e5a9bad61e9f3da626f3fb070c`
- The uploaded archive was reconciled against the authoritative commit at the complete tree level: 305/305 file paths and Git blob hashes matched.
- Historical Phase 9.5 documents retain their historical provenance statements; this document does not rewrite them.

## Production path exercised

The harness uses:

`storage.Open` → `chain.Open` → `block.New` → real header hashing / PoW → `chain.Accept` → incremental state validation → persistent state checkpoints → `storage.Close` → `storage.Open` → `chain.Open`.

The harness does not replace difficulty calculation, timestamp validation, proof-of-work validation, chain acceptance, Merkle calculation, or state commitment logic.

## Consensus parameters read from source

- Target block time: 30 seconds
- Difficulty interval: 10 blocks
- Minimum difficulty: 1
- Maximum difficulty: 32
- Maximum adjustment factor: 4
- Timestamp future bound: 4 × target block time = 120 seconds
- Median-time-past window: 11 blocks
- Genesis difficulty: 0
- First mineable difficulty: 4

## Deterministic timestamp schedule

The experiment uses three deterministic epoch modes, selected only by the harness:

- fast: +6 seconds per block
- slow: +120 seconds per block
- stable: +30 seconds per block

All remain within the production timestamp future bound. The schedule crosses 40 retarget boundaries in a 400-block run and exercises both upward and downward retarget behavior without changing production consensus constants.

The harness computes expected difficulty using the same exported `consensus.NextDifficulty` implementation used by production validation and verifies the accepted block against that value.

## Recorded fields

Each block JSONL record contains:

- height
- hash
- previous hash
- timestamp
- timestamp delta
- difficulty
- required target prefix (`0` repeated difficulty times; the implementation has no separate numeric target field)
- mining/production time
- validation time
- retarget calculation time at adjustment boundaries
- cumulative work
- transaction count
- state root
- Merkle root
- validation result
- expected difficulty and comparison result

## Acceptance criteria

The harness exits non-zero if any required invariant fails. A successful run requires all 400 non-genesis blocks to be accepted as `best`, full-chain validation to succeed after every block, timestamps to satisfy the current rules, difficulty to match the production calculation, cumulative work to increase, state roots to progress, Merkle roots to validate, and a close/reopen cycle to reproduce the final height, tip, state root, and cumulative work.

## Reproducibility

Two runs with the same source and deterministic configuration should produce identical consensus fields including block hashes, timestamps, difficulty schedule, cumulative work, state roots, and Merkle roots. Mining/validation durations are expected to differ between runs and are therefore not part of the deterministic comparison.

## Runtime limitations

A full Phase 10.1 verification requires the Go toolchain declared by `go.mod` (`go 1.27`, `toolchain go1.27.1`) and the pinned secp256k1 dependency. If the environment cannot obtain those dependencies, the 400-block production Go experiment is `NOT EXECUTED`; no PASS is inferred from static inspection or Python reference tests.
