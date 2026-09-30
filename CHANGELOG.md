# Changelog

## Phase 10.1 — 2026-09-28

- Added a reusable 400-block Go consensus/difficulty regression harness.
- Harness exercises production block hashing, proof-of-work acceptance, chain validation, incremental state commitments, persistence reload, and cumulative-work progression.
- Added deterministic timestamp scheduling that crosses 40 difficulty retarget boundaries and exercises downward, upward, and stable difficulty behavior without changing production consensus constants.
- Added machine-readable JSONL evidence output and a deterministic two-run comparison utility.
- Added Phase 10 provenance and validation documentation.
- Empirical consensus-path validation accepted 400/400 blocks and reproduced all consensus fields across independent runs.
- Production Go 1.27 build/test/vet/race validation remains NOT EXECUTED in the isolated build environment; no production readiness claim is made.
- Project remains NOT PRODUCTION-READY / NOT MAINNET-READY.

## Phase 9.5 Working Candidate — 2026-09-27

- Added deterministic `SYJ-STATE-ROOT-V1` state commitment implementation.
- Added incremental candidate state application and checkpoint persistence/recovery paths.
- Added checkpoint corruption/root-mismatch detection.
- Added Python state-root reference implementation and deterministic tests.
- Added explicit Merkle V3 construction with domain separation and leaf-count commitment.
- Added Merkle V3 deterministic vectors and ambiguity regression coverage.
- Added Phase 9.5 runtime verification harness.
- Added Phase 9.5 security/state/Merkle/audit documentation.
- Existing Phase 9.3 runtime-remediated findings remain `REMEDIATED`; no unsupported runtime verification status was added.
- Current local validation: Go `1.27.0 android/arm64` full test/vet/build gates PASS; Python `3.14.6` reports `262 passed, 1 warning`; compileall, pip-audit, and Bandit PASS. `staticcheck`, `govulncheck`, and `gosec` remain unexecuted, as do live runtime verification, supported-target race validation, fuzzing, the empirical 400-block mining/difficulty regression, and external security audit.
- Project remains NOT PRODUCTION-READY / NOT MAINNET-READY.
