# Changelog

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
