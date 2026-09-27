# Phase 9.5 State Root Specification

**Status:** IMPLEMENTED AS A VERSIONED SIDECAR COMMITMENT — NOT YET A BLOCK-HEADER CONSENSUS FIELD
**State schema:** `SYJ-STATE-ROOT-V1`

## Scope

Phase 9.5 introduces a deterministic state commitment and checkpoint format without changing historical block hashes. The commitment is selected by an explicit state schema version and is not silently applied to historical blocks.

## State model

The committed state contains:

- account balances: `address -> uint64 base units`
- account nonces: `address -> uint64`
- genesis supply: `uint64`
- mining issuance: `uint64`
- total supply: `uint64`

Missing balance/nonce values serialize as zero.

## Root construction

1. Start SHA-256 with domain `SYJ-STATE-ROOT-V1\x00`.
2. Append one byte state version (`0x01`).
3. Append `genesis_supply`, `mining_issued`, and `supply` as unsigned 64-bit big-endian integers.
4. Form the union of balance and nonce keys.
5. Sort keys by their UTF-8 byte/string ordering.
6. Append the account count as uint64 big-endian.
7. For each key append domain `SYJ-STATE-ACCOUNT-V1\x00`, an 8-byte length-prefixed UTF-8 key, balance uint64, and nonce uint64.
8. SHA-256 of the resulting byte stream is the state root, represented as lowercase hexadecimal.

## Checkpoints

Checkpoint serialization is versioned and deterministic. It contains the same monetary/account state, sorted account keys, and a trailing ASCII state-root checksum. The decoder recomputes the root and rejects a mismatch.

Checkpoints are persisted by block hash. The reference implementation uses an interval of 16 blocks for new checkpoints. On restart it searches backward from the active tip for the nearest valid checkpoint and applies only subsequent blocks. If no checkpoint exists, recovery starts from deterministic genesis state.

## Block acceptance

For ordinary candidate acceptance, the node obtains the parent state snapshot and applies only the candidate block's validated transactions. It does not reconstruct state from genesis for each candidate.

Difficulty retarget validation still uses the protocol-defined retarget window; this is independent of account-state replay.

## Fork/reorg

A fork parent can be reconstructed from the nearest checkpoint on that branch. The candidate branch is then applied incrementally. When a branch becomes active, the resulting snapshot becomes the active committed state and its root is exposed through the chain state API.

## Corruption

Checkpoint decoding rejects:

- invalid checkpoint magic
- unsupported state version
- oversized account counts/fields
- malformed/truncated data
- trailing bytes
- state-root checksum mismatch

## Compatibility

Historical V1/V2 block hashes and historical Merkle roots remain unchanged. The state root is not inserted into the historical block header in Phase 9.5. Making the state root a consensus block-header field requires a separately versioned protocol amendment, activation rule, migration rule, and independent verification.

## Cross-language contract

`blockchain/state_root.py` implements the same root serialization for the Python reference. Deterministic vectors and differential verification are required before the state root can be treated as a consensus commitment.
