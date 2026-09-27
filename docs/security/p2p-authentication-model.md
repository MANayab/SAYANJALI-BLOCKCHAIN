# Phase 9.3 P2P Authentication Model

## Security model

A node identity is the hash-bound public key identity already defined by the Phase 9 protocol. A HELLO is an identity assertion, not proof of liveness. Phase 9.3 therefore adds a signed `HELLO_FINISH` exchange.

Handshake:

1. Connecting node sends signed HELLO containing a fresh challenge.
2. Receiving node validates identity/signature and replies with a fresh server challenge.
3. Connecting node must sign that server challenge in HELLO_FINISH.
4. Only after HELLO_FINISH verification is the peer inserted into the authenticated peer map.

A captured HELLO cannot complete the handshake because the attacker does not know the fresh server challenge. The old 4,096-entry/10-minute HELLO challenge replay cache is no longer used for authentication admission.

## Duplicate identities

A newly authenticated session must not evict an existing authenticated session for the same node ID. The new connection is rejected. This prevents replayed identity assertions from displacing a live peer.

## Admission control

Pending TCP handshakes remain bounded globally and per source IP. No unbounded identity-keyed replay cache is allocated to arbitrary unauthenticated HELLOs. Production deployments must use TLS 1.3 with mutual certificate verification; plaintext is retained only for loopback/private development.

## Session integrity

Production P2P transport requires TLS 1.3. Every application write has a bounded write deadline. Plaintext transport is not considered authenticated session integrity and is rejected for non-loopback node listeners.

## Restart and replay

Authentication freshness is established by the fresh server challenge rather than by volatile replay-cache lifetime. Restart therefore does not make a previously captured HELLO sufficient for peer authentication.
