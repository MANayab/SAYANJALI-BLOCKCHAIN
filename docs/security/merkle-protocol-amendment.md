# Phase 9.5 C-02 Merkle Protocol Amendment

**Amendment:** C-02 / Merkle Commitment V3
**Status:** IMPLEMENTED AS AN EXPLICIT VERSIONED CONSTRUCTION; ACTIVATION DEFERRED

## Problem

The historical duplicate-last construction permits an ambiguity where the tree for `[A,B,C]` can share the same duplicated-last intermediate construction as `[A,B,C,C]`.

Historical V1/V2 Merkle behavior is frozen and is not modified by this amendment.

## V3 construction

### Leaf

```text
SHA256("SYJ-MERKLE-V3-LEAF\\x00" || u64_be(byte_length(leaf)) || leaf_bytes)
```

### Internal node

```text
SHA256("SYJ-MERKLE-V3-NODE\\x00" || left_32_bytes || right_32_bytes)
```

### Odd-node rule

If a level contains an odd number of nodes, the final node is duplicated for that level. The final root commitment additionally includes the exact leaf count, which removes the historical count ambiguity.

### Root

```text
SHA256(
  "SYJ-MERKLE-V3-ROOT\\x00" ||
  u64_be(leaf_count) ||
  tree_root_32_bytes
)
```

### Empty tree

The tree component is:

```text
SHA256("SYJ-MERKLE-V3-EMPTY\\x00")
```

and the final root still commits `leaf_count = 0`.

## Compatibility

- V1/V2 roots remain unchanged.
- Existing genesis and historical vectors are not rewritten.
- V3 is available through explicit protocol-version selection APIs.
- V3 is not silently activated for existing V2 chains.

## Activation rule

A future activation must specify an exact block/network boundary, protocol version, peer compatibility rule, historical-root treatment, and migration procedure before V3 becomes consensus-active.

## Test vectors

`protocol/test-vectors/merkle-v3.json` and the Go vector tests provide deterministic values.

The regression specifically requires:

```text
MerkleV3([A,B,C]) != MerkleV3([A,B,C,C])
```

## Acceptance

C-02 can move from REMEDIATED to VERIFIED only after the implementation and vectors are independently executed against the built implementation. It cannot move to CLOSED until activation/compatibility review is complete.
