#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
OUT_DIR="${1:-${ROOT}/evidence/phase10_1/runtime}"
rm -rf "$OUT_DIR"
mkdir -p "$OUT_DIR"

exec go run ./cmd/phase10_1 \
  -blocks 400 \
  -data-dir "$OUT_DIR/chain" \
  -jsonl "$OUT_DIR/blocks.jsonl" \
  -report "$OUT_DIR/report.md" \
  -run-id phase10.1-400-block
