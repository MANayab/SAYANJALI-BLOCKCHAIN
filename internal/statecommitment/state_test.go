package statecommitment

import "testing"

func TestRootDeterministicIndependentOfMapOrder(t *testing.T) {
	a := NewSnapshot()
	a.Balances["SYJa"] = 10
	a.Balances["SYJb"] = 20
	a.Nonces["SYJa"] = 2
	a.GenesisSupply = 30
	a.MiningIssued = 5
	a.Supply = 35
	b := NewSnapshot()
	b.Balances["SYJb"] = 20
	b.Balances["SYJa"] = 10
	b.Nonces["SYJa"] = 2
	b.GenesisSupply = 30
	b.MiningIssued = 5
	b.Supply = 35
	if a.Root() != b.Root() {
		t.Fatalf("roots differ: %s %s", a.Root(), b.Root())
	}
}

func TestCheckpointRoundTripAndCorruptionDetection(t *testing.T) {
	s := NewSnapshot()
	s.Balances["SYJa"] = 123
	s.Nonces["SYJa"] = 7
	s.Supply = 123
	enc := Encode(s)
	got, err := Decode(enc)
	if err != nil {
		t.Fatal(err)
	}
	if got.Root() != s.Root() {
		t.Fatalf("root changed after round trip")
	}
	enc[len(enc)-1] ^= 1
	if _, err := Decode(enc); err == nil {
		t.Fatal("corrupted checkpoint accepted")
	}
}

func TestRootChangesForConsensusState(t *testing.T) {
	base := NewSnapshot()
	base.Balances["SYJa"] = 10
	base.Supply = 10
	variants := []*Snapshot{base.Clone(), base.Clone(), base.Clone()}
	variants[0].Balances["SYJa"] = 11
	variants[1].Nonces["SYJa"] = 1
	variants[2].Supply = 11
	roots := map[string]struct{}{base.Root(): {}}
	for _, v := range variants {
		roots[v.Root()] = struct{}{}
	}
	if len(roots) != 4 {
		t.Fatalf("expected four distinct roots, got %d", len(roots))
	}
}
