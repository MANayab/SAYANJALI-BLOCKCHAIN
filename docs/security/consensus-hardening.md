# Phase 9.3 Consensus Hardening

## Implemented

- V1 transaction identities are rejected when duplicated within a block or active-chain prefix.
- High-S ECDSA signatures are rejected at the Go cryptographic verification boundary.
- Consensus block and transaction timestamps are restricted to finite integer seconds.
- PoW is checked before transaction/state processing so invalid candidates fail before expensive state validation.
- The Go miner uses actual Unix seconds (or parent+1 when the local clock has not advanced) instead of manufacturing a +30 second timestamp for every block. This exposes sustained fast mining to the existing retarget mechanism.

## Contract-sensitive items

The existing Merkle construction can make `[A,B,C]` and `[A,B,C,C]` share a root. Fixing this safely requires a protocol commitment to transaction count/tree shape and therefore requires a Phase 9 amendment; it is intentionally not silently changed in Phase 9.3.

The timestamp representation remains a float field on the wire/JSON model for compatibility, but consensus now accepts only integer seconds. A future protocol amendment can migrate the encoded type to an explicit integer.

## Difficulty

Difficulty retarget remains the frozen epoch algorithm. The mining timestamp change removes the artificial +30-second-per-block behavior that masked fast production. A 400-block long-run regression must be executed with the exact Phase 9.3 binary before C-06 can be marked closed.
