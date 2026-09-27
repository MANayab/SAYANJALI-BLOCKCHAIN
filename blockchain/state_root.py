"""Deterministic Phase 9.5 state-root and checkpoint codec.

This sidecar commitment does not alter historical block hashes. It is selected
only by an explicit state schema version and is intended for cross-implementation
state validation and checkpoint recovery.
"""
from __future__ import annotations

import hashlib
import struct

VERSION = 1
_MAGIC = b"SYJSTATE1"


def _u64(v: int) -> bytes:
    if not 0 <= v <= (1 << 64) - 1:
        raise ValueError("uint64 out of range")
    return struct.pack(">Q", v)


def _bytes(v: bytes) -> bytes:
    return _u64(len(v)) + v


def state_root(*, balances: dict[str, int], nonces: dict[str, int], genesis_supply: int, mining_issued: int, supply: int) -> str:
    h = hashlib.sha256()
    h.update(b"SYJ-STATE-ROOT-V1\x00")
    h.update(bytes([VERSION]))
    h.update(_u64(genesis_supply)); h.update(_u64(mining_issued)); h.update(_u64(supply))
    keys = sorted(set(balances) | set(nonces))
    h.update(_u64(len(keys)))
    for key in keys:
        raw = key.encode("utf-8")
        h.update(b"SYJ-STATE-ACCOUNT-V1\x00")
        h.update(_bytes(raw))
        h.update(_u64(balances.get(key, 0)))
        h.update(_u64(nonces.get(key, 0)))
    return h.hexdigest()


def encode_checkpoint(*, balances: dict[str, int], nonces: dict[str, int], genesis_supply: int, mining_issued: int, supply: int) -> bytes:
    keys = sorted(set(balances) | set(nonces))
    out = bytearray(_MAGIC + bytes([VERSION]))
    out += _u64(genesis_supply) + _u64(mining_issued) + _u64(supply) + _u64(len(keys))
    for key in keys:
        out += _bytes(key.encode("utf-8"))
        out += _u64(balances.get(key, 0)) + _u64(nonces.get(key, 0))
    out += state_root(balances=balances, nonces=nonces, genesis_supply=genesis_supply, mining_issued=mining_issued, supply=supply).encode("ascii")
    return bytes(out)
