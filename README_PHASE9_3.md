# SYJ BLOCKCHAIN — Phase 9.3

Security remediation and consensus hardening package based on main commit `b821d158e2f5fb231025ada0d1417890d3fb0fcf`.

## Scope

This package hardens the existing Go/Python implementation. It does **not** implement PoS, staking, validator economics, governance execution, smart contracts, Cosmos SDK, CometBFT, or mainnet genesis.

## Major changes

- P2P HELLO authentication now requires a fresh server challenge and signed `HELLO_FINISH` before authenticated peer state is installed.
- Duplicate authenticated peer IDs no longer evict the existing live connection.
- The 4,096-entry HELLO replay cache is removed from admission authentication.
- Public/non-loopback P2P requires TLS 1.3 mutual certificate verification.
- P2P writes use deadlines and peer-map locks are not held during network writes.
- Periodic and orphan-triggered synchronization is added.
- V1 transaction identity replay is rejected across the active chain and public V1 is gated to loopback/private development.
- High-S ECDSA signatures are rejected.
- Consensus timestamps are integer seconds at the validation boundary.
- The Go miner no longer manufactures a +30-second timestamp for every block.
- PoW validation is moved before expensive transaction/state processing.
- Python `ecdsa` dependency is removed in favor of `cryptography`.
- Python wallet creation API no longer returns private keys.
- API timeouts and bounded rate-limit state are added.
- CI is hardened with security scanner jobs and locked Python dependencies.

## Intentionally unresolved

The ambiguous legacy Merkle construction and unbounded fork-storage persistence require protocol/storage design review and are not silently rewritten. They remain open in the remediation matrix.

## Readiness

**NOT PRODUCTION-READY.** Phase 9.3 does not constitute public-testnet, incentivized-testnet, or mainnet approval.
