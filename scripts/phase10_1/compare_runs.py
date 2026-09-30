#!/usr/bin/env python3
"""Compare two Phase 10.1 JSONL runs on deterministic consensus fields."""
import json
import sys
from pathlib import Path

FIELDS = (
    "height", "hash", "previous_hash", "timestamp", "timestamp_delta",
    "difficulty", "required_target_prefix", "cumulative_work",
    "transaction_count", "state_root", "merkle_root", "validation_result",
    "expected_difficulty", "expected_difficulty_ok",
)


def load(path: Path):
    return [{k: obj[k] for k in FIELDS} for obj in map(json.loads, path.read_text().splitlines())]


def main() -> int:
    if len(sys.argv) != 3:
        print(f"usage: {sys.argv[0]} RUN_A.jsonl RUN_B.jsonl", file=sys.stderr)
        return 2
    a = load(Path(sys.argv[1]))
    b = load(Path(sys.argv[2]))
    if a != b:
        for idx, (left, right) in enumerate(zip(a, b), 1):
            if left != right:
                print(f"FAIL: first consensus difference at record {idx}")
                print(json.dumps({"run_a": left, "run_b": right}, indent=2, sort_keys=True))
                return 1
        if len(a) != len(b):
            print(f"FAIL: record count differs: {len(a)} != {len(b)}")
            return 1
        return 1
    print(f"PASS: {len(a)} consensus records are identical")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
