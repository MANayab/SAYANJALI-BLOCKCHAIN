# SAYANJALI BLOCKCHAIN — V3 C-1 Decisions

## Build boundary

Build 1 Recovery implements **C-1 — consensus time / clock semantics only**.
It does not implement C-2, H-1, H-2, or any other audit finding.

V3 is explicitly selectable through the existing Chain/Node protocol-version mechanism. V3 uses the existing V2 economic/state validation architecture for scaffolding, while Merkle V3 activation remains deferred to M-1.

## Q1 — TIME

### Consensus timestamp

For V3 blocks:

- timestamp is an integer-second value under the existing block representation.
- timestamp must be strictly greater than Median Time Past: `timestamp >= MTP + 1`.
- the legacy V1/V2 parent-relative `+120 seconds` ceiling remains unchanged for V1/V2.
- V3 does not use the parent `+120 seconds` ceiling.

### First ingress future bound

At first block ingress only:

`block.Timestamp <= Clock.Now() + 300 seconds`

The check occurs before a V3 block can be persisted by `Chain.Accept`.
The common `Chain.Accept` boundary covers direct block reception and synchronization because both network paths enter `Chain.Accept`.

### Replay

Historical validation and replay do not consult wall clock. The deterministic timestamp validator receives protocol history only. The 300-second check is an ingress-only guard.

Replay therefore remains independent of the current wall clock.

### Clock

Consensus ingress and mining use an injectable `Clock` abstraction:

- production: system clock
- tests: fixed/adversarial clocks

Unrelated operational timers continue to use their existing time sources.

### Mining

For V3:

`timestamp = max(Clock.Now(), MTP + 1)`

The miner therefore does not use `parent.Timestamp + 1` as its V3 timestamp source.

### Peer clock safety

Mining safety threshold: **60 seconds**.

If a legitimate peer median is available and:

`abs(localTime - peerMedian) > 60s`

V3 mining is refused.

If no legitimate peer median is available, V3 mining is also refused with:

`V3 mining refused: peer median clock unavailable`

This is an operational mining-safety rule, not a consensus block-rejection rule.

### Peer-time architecture limitation

The authoritative baseline has no authenticated peer timestamp field or peer-time service. Build 1 therefore:

- does not modify the P2P wire protocol,
- does not invent peer timestamps,
- does not use receive time, TCP time, or block timestamps as peer clock observations,
- does not silently assume the local clock is correct.

A dedicated `PeerMedian` abstraction exposes an explicit unavailable state. Future network-time work belongs outside this C-1 build.

## Q2 — STATE

The approved future C-2 direction is:

- per-block deltas
- chunked checkpoints
- authenticated sparse Merkle tree

**Not implemented in Build 1. C-2 remains OPEN.**

## Q3 — ACCOUNT ECONOMICS

The approved future direction is:

- minimum account balance
- ordinary transaction fees
- no separate account-creation fee

Numeric values require later spam-cost analysis.

**Not implemented in Build 1. H-2 remains OPEN.**

## Compatibility

V1 and V2 timestamp semantics remain on their existing parent-relative `+120 seconds` rule. V3 is explicitly selected as protocol version 3.

No V3 Merkle activation is performed.

## Deferred findings

C-2, H-1, H-2, H-3, H-4, H-5, H-6, H-7, H-8, M-1, M-2, M-3, M-4, M-5, M-6, M-7, M-8, L-1, L-2, and L-3 remain OPEN.
