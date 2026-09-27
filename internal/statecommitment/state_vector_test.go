package statecommitment

import "testing"

func TestStateRootProtocolVectors(t *testing.T) {
	cases := []struct {
		name                   string
		balances, nonces       map[string]uint64
		genesis, mined, supply uint64
		want                   string
	}{
		{"empty", map[string]uint64{}, map[string]uint64{}, 0, 0, 0, "91755aa11bcdbc16d735ccc1f58a755bbf40a8f9a9956bebfda97ade20aab863"},
		{"accounts", map[string]uint64{"SYJa": 10, "SYJb": 20}, map[string]uint64{"SYJa": 2}, 30, 5, 35, "c2327fa4645ef04ead0e59e5e8492977888efe346473e9def32f3e6eba423904"},
		{"transfer", map[string]uint64{"SYJa": 7, "SYJb": 23}, map[string]uint64{"SYJa": 3}, 30, 5, 35, "aebe5a161bfc0c10e64cd5f151011ac7c06b89f56f3538696a935c0a6aafc51e"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSnapshot()
			s.Balances, s.Nonces = tc.balances, tc.nonces
			s.GenesisSupply, s.MiningIssued, s.Supply = tc.genesis, tc.mined, tc.supply
			if got := s.Root(); got != tc.want {
				t.Fatalf("root=%s want=%s", got, tc.want)
			}
		})
	}
}
