package block

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

const MerkleProtocolV3 uint8 = 3

// MerkleRootV3 is the Phase 9.5 amended construction. Historical V1/V2 roots
// remain unchanged; callers select this construction only at an explicit
// protocol-version boundary.
func MerkleRootV3(txHashes []string) string {
	count := uint64(len(txHashes))
	var tree [32]byte
	if len(txHashes) == 0 {
		tree = sha256.Sum256([]byte("SYJ-MERKLE-V3-EMPTY\x00"))
	} else {
		level := make([][32]byte, len(txHashes))
		for i, leaf := range txHashes {
			level[i] = merkleV3Leaf(leaf)
		}
		for len(level) > 1 {
			next := make([][32]byte, 0, (len(level)+1)/2)
			for i := 0; i < len(level); i += 2 {
				right := level[i]
				if i+1 < len(level) {
					right = level[i+1]
				}
				next = append(next, merkleV3Node(level[i], right))
			}
			level = next
		}
		tree = level[0]
	}
	h := sha256.New()
	h.Write([]byte("SYJ-MERKLE-V3-ROOT\x00"))
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], count)
	h.Write(b[:])
	h.Write(tree[:])
	return hexString(h.Sum(nil))
}

func merkleV3Leaf(value string) [32]byte {
	h := sha256.New()
	h.Write([]byte("SYJ-MERKLE-V3-LEAF\x00"))
	writeLen(h, []byte(value))
	return sum32(h.Sum(nil))
}
func merkleV3Node(left, right [32]byte) [32]byte {
	h := sha256.New()
	h.Write([]byte("SYJ-MERKLE-V3-NODE\x00"))
	h.Write(left[:])
	h.Write(right[:])
	return sum32(h.Sum(nil))
}
func writeLen(h interface{ Write([]byte) (int, error) }, value []byte) {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(len(value)))
	_, _ = h.Write(b[:])
	_, _ = h.Write(value)
}
func sum32(v []byte) [32]byte { var out [32]byte; copy(out[:], v); return out }
func hexString(v []byte) string {
	const x = "0123456789abcdef"
	out := make([]byte, len(v)*2)
	for i, b := range v {
		out[i*2] = x[b>>4]
		out[i*2+1] = x[b&15]
	}
	return string(out)
}

func MerkleRootForProtocol(protocolVersion uint8, txHashes []string) (string, error) {
	switch protocolVersion {
	case 1, 2:
		return MerkleRoot(txHashes), nil
	case MerkleProtocolV3:
		return MerkleRootV3(txHashes), nil
	default:
		return "", fmt.Errorf("unsupported Merkle protocol version %d", protocolVersion)
	}
}
