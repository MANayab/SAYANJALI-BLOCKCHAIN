# Proposed Phase 9 Amendments — Phase 9.3

These are review items, not silent changes to the frozen Phase 9 contract.

## A-01 — P2P handshake freshness

Phase 9.3 increments the P2P protocol minor version and adds `HELLO_FINISH`. The purpose is to make a peer prove possession of its identity key over a fresh challenge before peer state is allocated. Migration impact: old P2P peers using the previous handshake cannot interoperate with the hardened handshake.

## A-02 — Integer timestamp encoding

Consensus currently rejects fractional timestamps but the serialized field remains a floating-point-compatible field. A future protocol version should encode timestamps as an explicit unsigned integer.

## A-03 — Merkle tree commitment

A production Merkle construction should commit to transaction count/tree shape using domain-separated leaf/node hashing or an explicit count commitment. This changes block identity and requires new vectors/genesis migration rules. Phase 9.3 does not silently rewrite existing vectors.
