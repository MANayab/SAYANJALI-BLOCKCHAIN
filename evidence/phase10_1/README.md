# Phase 10.1 Evidence Bundle

This directory contains the machine-readable evidence from two independent 400-block Phase 10.1 consensus-path runs.

## Provenance

- Source baseline: `ea26fe1db68f2e45ad7051cbedae655de36b0ca6`
- Uploaded source archive SHA-256: `6db49cdd86f6283f07178b3032401f3830c47e5a9bad61e9f3da626f3fb070c`
- Environment: Linux/amd64, Go 1.23.2 with a temporary disposable crypto stub because the repository requires Go 1.27.1 and the pinned secp256k1 module was unavailable.

## Files

- `400-block-run-a.jsonl` — final 400-block run.
- `400-block-run-b.jsonl` — independent comparison run.
- `400-block-run-a-report.md` — human-readable summary.
- `determinism-comparison.txt` — comparison of consensus fields with timing fields excluded.

## Result

- 400/400 blocks accepted in each run.
- 40 retarget boundaries crossed.
- Final height: 400.
- Final difficulty: 3.
- Final cumulative work: `18165761`.
- Final state root: `7128583826764b3b691fb557bd90e91a39d7261b1a7ffd6422224c7b795e680e`.
- Deterministic consensus comparison: PASS.
- Restart/reload validation: PASS.

These artifacts are consensus-path evidence only. They are not a production Go toolchain, dependency, cryptographic, race, fuzz, or security audit result.
