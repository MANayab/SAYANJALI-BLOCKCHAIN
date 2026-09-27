package block

import "testing"

func TestMerkleV3EliminatesDuplicateLastAmbiguity(t *testing.T) {
	abc := MerkleRootV3([]string{"A", "B", "C"})
	abcc := MerkleRootV3([]string{"A", "B", "C", "C"})
	if abc == abcc {
		t.Fatalf("Merkle V3 collision: [A,B,C] == [A,B,C,C] == %s", abc)
	}
}

func TestMerkleV3DeterministicShapes(t *testing.T) {
	cases := [][]string{nil, {"A"}, {"A", "B"}, {"A", "B", "C"}, {"A", "B", "C", "D", "E"}}
	seen := map[string]struct{}{}
	for _, in := range cases {
		a, b := MerkleRootV3(in), MerkleRootV3(in)
		if a != b {
			t.Fatalf("non-deterministic root for %v", in)
		}
		if _, ok := seen[a]; ok {
			t.Fatalf("unexpected vector root reuse for %v", in)
		}
		seen[a] = struct{}{}
	}
}
