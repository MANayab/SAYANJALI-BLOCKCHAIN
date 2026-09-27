# PHASE 9.5 MANIFEST

## Provenance

- Supplied archive source working-tree identity: `b35afa9468294c9649e9541c947734beab10cbea`
- Authoritative project baseline: `b821d158e2f5fb231025ada0d1417890d3fb0fcf`
- Current archive state is a working-tree snapshot and is **not** asserted to be byte-identical to either Git identity above.
- Current selected-file snapshot SHA-256: `9304989321f8e3635c1ea0c2c04f2a18a6f26b65e52991d2ebd02f2de4439a34`
- The selected-file snapshot hash is derived from the SHA-256 values listed below in deterministic path order.

## Important file SHA-256

| File | SHA-256 |
|---|---|
| `PHASE9_5_VALIDATION_REPORT.md` | `fa30fbb50e6437f6b6ca86e0bb8acbdee04d843d3c99e9a4ef28195e99eba1d9` |
| `CHANGELOG.md` | `aaf6ade0049306a419d3f19208087a071825f3316ce819e87ef98870a0a72370` |
| `ROADMAP.md` | `1c2789f161a8bcb4cc6c18ab11ed2d4270cc2a04cd8c5cf10ffdce96a7db033e` |
| `README.md` | `9477e0379184ee2029b2bf15ed6caaf950dff32f6dfa5f30e82e0fca98d53664` |
| `docs/security/phase9-5-remediation-matrix.md` | `8e2a5e0217b5fd2487427d9912da443529a74d0186343f64fcafedb89e08c159` |
| `docs/security/phase9-5-runtime-verification.md` | `fae53d9f6d30ad683ac961b0c8eda5f845ce0d383da2e4a31bbb84606b78d8e9` |
| `docs/security/state-validation-architecture.md` | `b5eaf4b65949234c5b9c87e69fef4957f06644b2eb2b258ee1866a5579bbc012` |
| `docs/security/state-root-specification.md` | `356c6101d147e91c228c6041b1f288a2363970d728fa32dd386572918ce2f2b6` |
| `docs/security/merkle-protocol-amendment.md` | `f7f5670c7789f0a09724ab3bc239b065b8dffd09d88c9f497104dda24f943867` |
| `docs/security/external-audit-plan.md` | `3dc9b4bfdfa2d0f54ad330badce65a5b09159f13874c523b8de208ba32850457` |
| `internal/statecommitment/state.go` | `72581873f446c80fffc74191a1fe8816008c16051010e2991ad80b8a7dee360f` |
| `internal/statecommitment/state_test.go` | `09645ef2a71059cac59afb02b36c4374d5930ea32587cbedcf7baa76c1f9d4da` |
| `internal/statecommitment/state_vector_test.go` | `ea2a33f27872997d7fe77db9b28fbbaf0a657309378289c0957614683a3ac104` |
| `internal/chain/chain.go` | `1ef5ed53486a44f9b132347e950aa01e40c71d0e15cd2b8b12b702848739cf78` |
| `internal/storage/storage.go` | `cc9bb3496867bd6e414557c20375e0d5de0f48fb3f111222b694c1122383743b` |
| `internal/block/merkle_v3.go` | `2f8da5e8bdb922ea1743de8ddffd893f801f551362f1ddb8fa6f525998f3ab44` |
| `internal/block/merkle_v3_test.go` | `87ac4c07eb1256dc8a754ecf74615b78d666000b4d6f2c0dfcca98b8c154b35a` |
| `internal/block/merkle_v3_vector_test.go` | `e874c2f5cc85ee5be15741e1ba7eb869b9b7ca9ab04f2dd9386edda133782a29` |
| `internal/block/merkle_protocol_test.go` | `c91cbb200f65d53811036da0f623672f61211dbf6c3ddd7f9df19c0990866f92` |
| `internal/block/block.go` | `b06000f9c868c70f0c8b6985bff243708fda38c1c9ae5a0cd749045dc7a3b4c4` |
| `blockchain/state_root.py` | `2e015e8c9d9435282172d14e003ed4b5c7c9cb7609f50c8ac6dc4663752445a2` |
| `tests/test_state_root.py` | `fc035d7793c3aea51a10260b8b9affcb1350fe522b6e1e3b8fc1337c3db3da91` |
| `protocol/test-vectors/merkle-v3.json` | `8542484453bff6b5fd2b5cb44063cbf8e5663cc7d961bbf0dad78c6bb1e8ddeb` |
| `protocol/test-vectors/state-root-v1.json` | `fe9a7fd1082195f5f096f7505ff2ad7ba56521ef1ada8d770723f883f1c44e74` |
| `scripts/phase9_5/runtime_verify.sh` | `5a997d93b50edd6141e409af1fd6203434266649b110cd79847b8f6dfd18dacc` |
| `scripts/phase9_5/README.md` | `84c47190eb79258377154c3bff57107c6936ce8df9ffb86c504d34939e99837c` |
| `go.mod` | `3e608d5a01a2dd7d2fb26364b661470303884ab22deccb85354a808a4b667523` |
| `requirements.lock.txt` | `6307c4bd4296dcbd64197dc229540583b2c2e51c8736697c8f5fb2be3fb84ddf` |
| `requirements.txt` | `6307c4bd4296dcbd64197dc229540583b2c2e51c8736697c8f5fb2be3fb84ddf` |

## Artifact

- Filename: `SYJ-BLOCKCHAIN-phase-9.5-working.zip`
- Project status: **NOT PRODUCTION-READY / NOT MAINNET-READY**
- Final ZIP SHA-256: **d8c33a3a49c094859e181e47e56058e9ab47f2d2aef91120267fcbe3815b22cd**

The final ZIP hash is intentionally recorded outside the ZIP because embedding a self-hash inside the archive is circular: changing the embedded hash changes the ZIP bytes.

## Validation provenance

- Go: `1.27.0 android/arm64`
- Python: `3.14.6`
- `GOTOOLCHAIN=local go test ./...`: PASS
- `GOTOOLCHAIN=local go vet ./...`: PASS
- `GOTOOLCHAIN=local go build ./...`: PASS
- `gofmt -l cmd internal pkg tests`: CLEAN
- Python tests: `262 passed, 1 warning`
- Python compileall: PASS
- pip-audit: PASS
- Bandit: PASS
- staticcheck: NOT EXECUTED
- govulncheck: NOT EXECUTED
- gosec: NOT EXECUTED
- Go race validation: NOT EXECUTED
- fuzzing: NOT EXECUTED
- 400-block empirical mining/difficulty regression: NOT EXECUTED
- live runtime verification: NOT EXECUTED
- external independent audit: NOT EXECUTED

## Interpretation

This manifest records provenance and file integrity for the current Phase 9.5 working snapshot. File hashes do not constitute independent security verification.

The project remains **NOT PRODUCTION-READY / NOT MAINNET-READY**.
