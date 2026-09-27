# Phase 9.3 Dependency Security

## Go

- Baseline: Go 1.23.2 local environment / Go 1.23.x CI baseline.
- Phase 9.3 target: Go 1.27.1, the current stable release selected for this remediation. Go 1.27.1 was released 2026-09-01.
- Prior audit evidence: 33 reachable Go 1.23 standard-library vulnerability findings.
- Required closure evidence: actual `govulncheck ./...` output on Go 1.27.1.

## Python

CVE-2024-23342 affected the `ecdsa` Python dependency used by the reference wallet/P2P identity. Phase 9.3 removes the `ecdsa` package from the dependency graph and migrates those Python cryptographic operations to the maintained `cryptography` package. Signatures remain raw R||S secp256k1 values and are normalized to low-S.

No statement of zero vulnerabilities is made until `pip-audit` executes against the locked dependency graph.

## Dependency locking

`requirements.lock.txt` records the tested dependency versions. CI installs this file and runs `pip-audit`. Hash pinning/SBOM publication remains a follow-on supply-chain gate.
