package block

import "testing"

func TestMerkleProtocolSelectionPreservesHistoricalVersions(t *testing.T) {
	hashes := []string{"A", "B", "C"}
	old := MerkleRoot(hashes)
	got1, err := MerkleRootForProtocol(1, hashes)
	if err != nil {
		t.Fatal(err)
	}
	got2, err := MerkleRootForProtocol(2, hashes)
	if err != nil {
		t.Fatal(err)
	}
	got3, err := MerkleRootForProtocol(3, hashes)
	if err != nil {
		t.Fatal(err)
	}
	if got1 != old || got2 != old {
		t.Fatalf("historical Merkle behavior changed")
	}
	if got3 == old {
		t.Fatalf("amended root unexpectedly equals historical root")
	}
}
