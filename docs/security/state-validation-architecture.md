# Phase 9.5 State Validation Architecture

**Status:** IMPLEMENTED INCREMENTAL STATE TRACK WITH CHECKPOINT RECOVERY
**Protocol impact:** Historical block serialization is unchanged.

## Previous architecture

The previous candidate validation path reconstructed account state by replaying the entire prefix from genesis. That creates increasing work as chain height grows and repeated work for independent candidate validation.

## Phase 9.5 architecture

```text
validated parent block
        |
        v
parent state snapshot
        +
validated candidate transactions
        |
        v
incremental state transition
        |
        +--> balances / nonces / supply
        |
        v
SYJ-STATE-ROOT-V1
        |
        v
optional checkpoint persistence
```

### Acceptance

1. Locate the candidate parent.
2. Obtain the parent snapshot from the in-memory state index or nearest persisted checkpoint.
3. Validate structural/consensus properties for the candidate.
4. Apply only the candidate transaction set to a cloned parent snapshot.
5. Compute the deterministic state root.
6. Persist the block and, at the checkpoint interval, the state snapshot.

### Restart

The node locates the active tip, searches backward for the nearest checkpoint, validates the checkpoint checksum/root, and applies only blocks after that checkpoint. No ordinary candidate requires full-history account-state replay.

### Fork/reorg

Fork branches use the nearest branch checkpoint/common known state and apply only the remaining branch blocks. The winning branch's resulting snapshot becomes the active state. Existing V2 reorg transaction handling remains in place.

## Complexity expectations

- Ordinary candidate state transition: proportional to candidate transaction count plus deterministic account-root serialization.
- Retarget calculation: bounded by the protocol difficulty interval.
- Restart recovery: proportional to the distance from the nearest valid checkpoint, plus checkpoint decoding.
- Fork recovery: proportional to branch distance from the nearest available checkpoint.

The root computation currently serializes all committed accounts, so root calculation remains O(number of committed accounts). This is intentionally documented rather than represented as an O(1) authenticated tree.

## Measurement requirement

Phase 9.5 does not claim a measured speedup until the old and new paths are benchmarked in the same environment. Benchmark output belongs in `evidence/phase9_5/benchmarks/` and must identify the environment and sample sizes.
