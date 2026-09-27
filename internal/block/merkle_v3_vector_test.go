package block

import "testing"

func TestMerkleV3ProtocolVectors(t *testing.T) {
	cases := []struct {
		name string
		in   []string
		want string
	}{
		{"empty", nil, "cf8095a19fd0b8199d8ed86aebf607fe37cbc66774495bcabebe91aec6e9242c"},
		{"one", []string{"A"}, "15b7188492895e2b0e6d6c6f11d9ad8c7d02ca6c451179b6069b1967d210b6a5"},
		{"two", []string{"A", "B"}, "dbe3ae7c87eebe5e024fcc4c1433bf3e095ead7c4129e282bddb0b65cb5fdd39"},
		{"three", []string{"A", "B", "C"}, "ae6672be6a105891a8c53e4208fca6ea0e4c4454591be8620e85bc6b1d5c8004"},
		{"four", []string{"A", "B", "C", "D"}, "be36f9ceb13b59813af22d325dea744aa49d5e001f1544301fa9a7eb1f6d9064"},
		{"five", []string{"A", "B", "C", "D", "E"}, "cdd88fdc41cd47791add7894c0a87163a6be9fc6ff19ffdf4ef13c70d666e6b6"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := MerkleRootV3(tc.in); got != tc.want {
				t.Fatalf("root=%s want=%s", got, tc.want)
			}
		})
	}
}
