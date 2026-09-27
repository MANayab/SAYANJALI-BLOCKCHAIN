#!/usr/bin/env bash
set -u

# Phase 9.5 runtime harness. It never substitutes source inspection for runtime evidence.
# Set SYJ_BIN to a built syjd binary. Every executed scenario writes stdout/stderr and
# exit status into the supplied evidence directory.

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
EVIDENCE_DIR=${EVIDENCE_DIR:-"$ROOT/evidence/phase9_5/runtime"}
SYJ_BIN=${SYJ_BIN:-"$ROOT/bin/syjd"}
mkdir -p "$EVIDENCE_DIR"

run_case() {
  local id="$1"; shift
  local out="$EVIDENCE_DIR/${id}.log"
  {
    echo "CASE=$id"
    echo "UTC=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
    echo "COMMAND=$*"
    echo "--- OUTPUT ---"
    "$@"
    rc=$?
    echo "--- EXIT_CODE=$rc ---"
    return "$rc"
  } >"$out" 2>&1
}

if [[ ! -x "$SYJ_BIN" ]]; then
  cat > "$EVIDENCE_DIR/HARNESS_NOT_EXECUTED.log" <<LOG
CASE=HARNESS
STATUS=NOT EXECUTED
REASON=SYJ_BIN is not an executable built node binary.
EXPECTED=$SYJ_BIN
UTC=$(date -u +%Y-%m-%dT%H:%M:%SZ)
No runtime security finding is upgraded by this condition.
LOG
  echo "NOT EXECUTED: build syjd and set SYJ_BIN=$SYJ_BIN" >&2
  exit 2
fi

DATA_DIR=$(mktemp -d "${TMPDIR:-/tmp}/syj95-node.XXXXXX")
cleanup() { "$SYJ_BIN" stop "$DATA_DIR" >/dev/null 2>&1 || true; rm -rf "$DATA_DIR"; }
trap cleanup EXIT

run_case H09_GO_VERSION go version || true
run_case H04_INIT "$SYJ_BIN" init "$DATA_DIR" || true

# Start a real node process. Subsequent HTTP cases target this process rather than mocks.
"$SYJ_BIN" start "$DATA_DIR" >"$EVIDENCE_DIR/syjd.stdout.log" 2>"$EVIDENCE_DIR/syjd.stderr.log" &
PID=$!
echo "PID=$PID" >"$EVIDENCE_DIR/syjd.pid"

for _ in $(seq 1 60); do
  if curl -fsS --max-time 1 http://127.0.0.1:8080/health >"$EVIDENCE_DIR/health.json" 2>"$EVIDENCE_DIR/health.err"; then break; fi
  sleep 0.25
done

run_case M01_HEALTH curl -fsS --max-time 5 http://127.0.0.1:8080/health || true
run_case M02_STATUS curl -fsS --max-time 5 http://127.0.0.1:8080/status || true
run_case L01_MALFORMED_REQUEST curl -sS -o /dev/null -w 'HTTP=%{http_code} TIME=%{time_total}\n' --max-time 5 -H 'content-type: application/json' -d '{"timestamp":1.5}' http://127.0.0.1:8080/transaction/submit || true

cat > "$EVIDENCE_DIR/README.txt" <<TXT
This directory is produced by scripts/phase9_5/runtime_verify.sh.
It targets a real syjd process. It does not use mocked node behavior.
P2P-specific D-01/D-02/D-03/H-04/H-05/H-06/H-07 scenarios require the
corresponding integration cases and are not marked verified by this script alone.
TXT

"$SYJ_BIN" stop "$DATA_DIR" >"$EVIDENCE_DIR/stop.stdout.log" 2>"$EVIDENCE_DIR/stop.stderr.log" || true
wait "$PID" 2>/dev/null || true
