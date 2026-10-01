package node

import (
	"fmt"
	"testing"
	"time"

	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/chain"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/clock"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/storage"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/tokenomics"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/wallet"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/pkg/protocol"
)

type fixedPeerMedian struct {
	t  time.Time
	ok bool
}

func (p fixedPeerMedian) MedianTime() (time.Time, bool) { return p.t, p.ok }

func c1State(t *testing.T) tokenomics.GenesisState {
	t.Helper()
	addr := func(seed byte) string {
		var pub [64]byte
		for i := range pub {
			pub[i] = seed
		}
		return wallet.AddressFromPublicKeyHex(fmt.Sprintf("%x", pub[:]))
	}
	return tokenomics.GenesisState{Version: tokenomics.Version, Sender: protocol.GenesisAllocationSender, Allocations: []tokenomics.Allocation{
		{Category: tokenomics.Presale, Recipient: addr(1), AmountBaseUnits: 7_200_000_000_000_000},
		{Category: tokenomics.Treasury, Recipient: addr(2), AmountBaseUnits: 5_760_000_000_000_000},
		{Category: tokenomics.Ecosystem, Recipient: addr(3), AmountBaseUnits: 3_600_000_000_000_000},
		{Category: tokenomics.Liquidity, Recipient: addr(4), AmountBaseUnits: 4_320_000_000_000_000},
		{Category: tokenomics.Team, Recipient: addr(5), AmountBaseUnits: 7_920_000_000_000_000},
	}, TotalBaseUnits: tokenomics.ExpectedGenesisTotal()}
}

func newC1NodeChain(t *testing.T, now time.Time) *chain.Chain {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	gs := c1State(t)
	id := "2222222222222222222222222222222222222222222222222222222222222222"
	c, err := chain.OpenV3WithGenesisState(st, gs, id, "syjnet-v2-"+id)
	if err != nil {
		t.Fatal(err)
	}
	c.SetClock(clock.FixedClock{Current: now})
	t.Cleanup(func() { _ = st.Close() })
	return c
}

func TestC1V3MiningClockSafety(t *testing.T) {
	now := time.Unix(1735690000, 0)
	c := newC1NodeChain(t, now)
	if _, err := MineNextWithClockSafety(c, "SYJaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", nil, clock.FixedClock{Current: now}, clock.UnavailablePeerMedian{}); err == nil {
		t.Fatal("mining allowed without peer median")
	}
	if _, err := MineNextWithClockSafety(c, "SYJaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", nil, clock.FixedClock{Current: now}, fixedPeerMedian{t: now.Add(-61 * time.Second), ok: true}); err == nil {
		t.Fatal("mining allowed with >60s skew")
	}
	b, err := MineNextWithClockSafety(c, "SYJaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", nil, clock.FixedClock{Current: now}, fixedPeerMedian{t: now, ok: true})
	if err != nil {
		t.Fatalf("mining at <=60s skew failed: %v", err)
	}
	if b.Timestamp != float64(now.Unix()) {
		t.Fatalf("V3 miner timestamp=%v want=%d", b.Timestamp, now.Unix())
	}
}
