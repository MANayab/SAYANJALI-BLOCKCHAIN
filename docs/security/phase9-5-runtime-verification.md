# Phase 9.5 Runtime Verification Harness

**Status:** HARNESS IMPLEMENTED; REQUIRED RUNTIME EXECUTION NOT COMPLETED IN THIS ENVIRONMENT

## Reproduction

```bash
SYJ_BIN=/absolute/path/to/syjd ./scripts/phase9_5/runtime_verify.sh
```

The harness requires a real built `syjd` binary. It starts the process, records
stdout/stderr and exit codes, probes the live HTTP API, and writes evidence under
`evidence/phase9_5/runtime/`.

A missing binary causes `NOT EXECUTED`. The harness never upgrades a finding on
that basis.

## Independent-runtime rule

A finding is eligible for `VERIFIED` only when its attack/regression is executed
against the actual built node/transport and the captured evidence satisfies the
Phase 9.5 acceptance criteria.

The current local environment provides Go `1.27.0 android/arm64` and Python
`3.14.6`. The available local validation gates have been executed separately:
the full Go test/vet/build gates pass; Python tests report `262 passed, 1 warning`;
compileall, pip-audit, and Bandit pass.

The live Phase 9.5 runtime harness itself remains **NOT EXECUTED**. No live
node/transport finding is upgraded to `VERIFIED` based solely on source-level
tests or local build validation.

The current environment does not provide `staticcheck`, `govulncheck`, or
`gosec`, so those scanner gates remain unexecuted. Go race validation on
Android/ARM64, fuzzing, the empirical 400-block mining/difficulty regression,
and external audit also remain unexecuted.

Existing Phase 9.3 source/test evidence may support a `REMEDIATED` engineering
status, but independent runtime verification is still required before a
finding can be marked `VERIFIED`.


## Current Local Validation Snapshot

**Environment**

- Go: `go1.27.0 android/arm64`
- Python: `3.14.6`

**Executed successfully**

- `GOTOOLCHAIN=local go test ./...` — PASS
- `GOTOOLCHAIN=local go vet ./...` — PASS
- `GOTOOLCHAIN=local go build ./...` — PASS
- gofmt verification — PASS
- `python3 -m pytest -q` — 262 passed, 1 warning
- `python3 -m compileall -q blockchain tests` — PASS
- `pip-audit -r requirements.lock.txt` — PASS / No known vulnerabilities
- Bandit — PASS / No issues identified

**Not executed**

- `staticcheck`
- `govulncheck`
- `gosec`
- Go race detector
- fuzzing
- empirical 400-block mining/difficulty regression
- live `scripts/phase9_5/runtime_verify.sh` execution
- external security audit

**Interpretation:** These results demonstrate local source/build/test validation only.
They do not constitute live-node security verification or production approval.

**Project status:** **NOT PRODUCTION-READY / NOT MAINNET-READY.**
