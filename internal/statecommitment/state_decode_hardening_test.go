package statecommitment

import (
	"bytes"
	"testing"
)

func TestDecodeRejectsTruncatedCheckpoint(t *testing.T) {
	s := NewSnapshot()
	s.Balances["SYJa"] = 10
	s.Nonces["SYJa"] = 1
	s.Supply = 10

	encoded := Encode(s)

	for i := 0; i < len(encoded); i++ {
		if _, err := Decode(encoded[:i]); err == nil {
			t.Fatalf("truncated checkpoint accepted at length %d", i)
		}
	}
}

func TestDecodeRejectsDuplicateAccount(t *testing.T) {
	var b bytes.Buffer

	b.WriteString(checkpointMagic)
	b.WriteByte(Version)

	writeU64(&b, 0)  // GenesisSupply
	writeU64(&b, 0)  // MiningIssued
	writeU64(&b, 10) // Supply

	writeU64(&b, 2) // Two entries

	writeBytes(&b, []byte("SYJa"))
	writeU64(&b, 10)
	writeU64(&b, 1)

	writeBytes(&b, []byte("SYJa"))
	writeU64(&b, 20)
	writeU64(&b, 2)

	// Root is deliberately irrelevant: duplicate detection must happen
	// before the checkpoint can be accepted.
	b.Write(make([]byte, 64))

	if _, err := Decode(b.Bytes()); err == nil {
		t.Fatal("duplicate account checkpoint accepted")
	}
}
