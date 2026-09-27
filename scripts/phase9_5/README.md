# Phase 9.5 Verification Harness

Build the node first, then run:

```bash
SYJ_BIN=/absolute/path/to/syjd ./scripts/phase9_5/runtime_verify.sh
```

The harness starts a real node process and records commands, UTC timestamps,
stdout/stderr, and exit codes under `evidence/phase9_5/runtime/`.

A missing binary results in `NOT EXECUTED` and a non-zero harness exit. It does
not upgrade any security finding.

P2P-specific scenarios must use the real transport implementation and should
be added as integration cases before their findings are moved to `VERIFIED`.
