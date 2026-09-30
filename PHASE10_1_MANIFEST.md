# Phase 10.1 Delivery Manifest

## Delivery artifact

- Filename: `SYJ-BLOCKCHAIN-phase-10.1-400-block-consensus.zip`
- Exact ZIP SHA-256: **recorded in the companion `.sha256` file and final delivery report**.
- Note: embedding the exact SHA-256 of a ZIP inside that same ZIP is self-referential and changes the hash. The companion checksum is therefore the authoritative exact archive hash.

## Source baseline

- Authoritative commit: `ea26fe1db68f2e45ad7051cbedae655de36b0ca6`
- Tag: `v0.9.5-security-state-hardening`
- Source tree SHA: `a6add2f4a874dbbeb000b57488a547c0e1983736`
- Uploaded source archive SHA-256: `6db49cdd86f6283f07178b3032401f3830c47e5a9bad61e9f3da626f3fb070c`
- Source reconciliation: 305/305 path/mode/blob entries matched.

## Important Phase 10.1 files

Hashes are Git blob SHA-1 values for the final source files.

| File | Git blob SHA-1 |
|---|---|
| `cmd/phase10_1/main.go` | `144d7b1e05ed201755d65b72334d5c20e53fc502` |
| `cmd/phase10_1/main_test.go` | `b1ff85978d0af859f01d517acc354ddcbc552c8c` |
| `scripts/phase10_1/run_400_block.sh` | `e131dffd7b0d5323d02fa83d1fd811fc7500ef8b` |
| `scripts/phase10_1/compare_runs.py` | `84f0c2c5a6060d0c8676a07d6d691355df772d0f` |
| `docs/consensus/phase10-400-block-regression.md` | `a75502ccce18e066fdf9bac27581e0d6e0254010` |
| `docs/security/phase10-runtime-validation.md` | `29031eca1a0947544ce4d5ab50d123715595337c` |
| `PHASE10_PROVENANCE.md` | `559ade459f6d19c4ec6e890b9a3c3aa0612036cf` |
| `PHASE10_1_VALIDATION_REPORT.md` | `1f2c87abc3d9575c1736f2b7e00b46f3100915e3` |

## Validation summary

- 400-block consensus-path run: PASS
- 400/400 blocks accepted
- 40 retarget boundaries crossed
- Deterministic consensus comparison: PASS
- Restart/reload validation: PASS
- Python tests: 262 passed
- Production Go 1.27 build/test/vet/race: NOT EXECUTED because required toolchain was unavailable
- staticcheck/govulncheck/gosec/Bandit/pip-audit: NOT EXECUTED because tools were unavailable

## Known limitations

The empirical 400-block execution used a temporary crypto stub in a disposable validation copy solely because the production Go toolchain and pinned secp256k1 module were unavailable. The stub is not part of the delivery archive. This evidence does not constitute cryptographic security validation or a production-build PASS.

The state root remains a sidecar/checkpoint commitment and Merkle V3 remains unactivated.
