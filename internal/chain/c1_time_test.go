package chain

import (
	"math/rand"
	"testing"
	"time"

	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/block"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/clock"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/storage"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/transaction"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/pkg/protocol"
)

const c1TestNetworkID = "1111111111111111111111111111111111111111111111111111111111111111"

func openC1V3(t *testing.T, dir string, clk clock.Clock) *Chain {
	t.Helper()
	st, err := storage.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	gs := phase7TestState(t)
	c, err := openWithProtocolAndClock(st, &gs, 3, c1TestNetworkID, "syjnet-v2-"+c1TestNetworkID, clk)
	if err != nil {
		_ = st.Close()
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return c
}

func c1V3Block(t *testing.T, c *Chain, timestamp float64) *block.Block {
	t.Helper()
	reward, ok := uint64(5_000_000_000), true
	if !ok {
		t.Fatal("unreachable")
	}
	tx, err := transaction.NewV2Coinbase("SYJaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", c.NetworkID(), reward, timestamp)
	if err != nil {
		t.Fatal(err)
	}
	b, err := block.New(c.Tip().Index+1, c.Tip().Hash, timestamp, 0, 4, []transaction.Transaction{tx})
	if err != nil {
		t.Fatal(err)
	}
	for b.Nonce < 5_000_000 {
		if b.MeetsDifficulty(4) {
			return b
		}
		b.Nonce++
		if err := b.Recompute(); err != nil {
			t.Fatal(err)
		}
	}
	t.Fatal("unable to mine test block")
	return nil
}

func TestC1V3TimestampBoundaries(t *testing.T) {
	prefix := []*block.Block{{Header: block.Header{Timestamp: 100}}, {Header: block.Header{Timestamp: 110}}}
	cases := []struct {
		name    string
		ts      float64
		wantErr bool
	}{
		{"mtp", 110, true}, {"below_mtp", 109, true}, {"mtp_plus_one", 111, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateTimestampForProtocol(tc.ts, prefix, 3)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err=%v wantErr=%v", err, tc.wantErr)
			}
		})
	}
}

func TestC1V3IngressFutureBoundary(t *testing.T) {
	now := time.Unix(1735689900, 0)
	for _, delta := range []int64{300, 301} {
		t.Run(time.Duration(delta).String(), func(t *testing.T) {
			c := openC1V3(t, t.TempDir(), clock.FixedClock{Current: now})
			b := c1V3Block(t, c, float64(now.Unix()+delta))
			ok, _, err := c.Accept(b)
			if delta == 300 && (err != nil || !ok) {
				t.Fatalf("+300 rejected: ok=%v err=%v", ok, err)
			}
			if delta == 301 && (err == nil || ok) {
				t.Fatalf("+301 accepted: ok=%v err=%v", ok, err)
			}
		})
	}
}

func TestC1V3LegacyParent120IsNotV3Ceiling(t *testing.T) {
	prefix := []*block.Block{{Header: block.Header{Timestamp: 100}}}
	if err := validateTimestampForProtocol(401, prefix, 3); err != nil {
		t.Fatalf("V3 timestamp beyond legacy +120 was rejected: %v", err)
	}
	if err := validateTimestampForProtocol(221, prefix, 2); err == nil {
		t.Fatal("V2 legacy +120 rule was removed")
	}
}

func TestC1V3SeededTimestampDeterminism(t *testing.T) {
	const seed int64 = 20261001
	r1 := rand.New(rand.NewSource(seed))
	r2 := rand.New(rand.NewSource(seed))
	for i := 0; i < 500; i++ {
		a := int64(r1.Intn(601))
		b := int64(r2.Intn(601))
		if a != b {
			t.Fatalf("seed divergence at %d: %d != %d", i, a, b)
		}
	}
}

func TestC1V3ReplayIsIndependentOfClock(t *testing.T) {
	dir := t.TempDir()
	base := clock.FixedClock{Current: time.Unix(int64(protocol.GenesisTimestamp+1000), 0)}
	c := openC1V3(t, dir, base)
	b := c1V3Block(t, c, protocol.GenesisTimestamp+1001)
	if ok, _, err := c.Accept(b); err != nil || !ok {
		t.Fatalf("accept: %v", err)
	}
	wantHeight, wantTip, wantRoot, wantWork := c.Height(), c.TipHash(), c.StateRoot(), c.Work().String()
	_ = c.store.Close()

	for name, now := range map[string]time.Time{
		"decades_ahead":  base.Now().Add(10 * 365 * 24 * time.Hour),
		"decades_behind": base.Now().Add(-10 * 365 * 24 * time.Hour),
	} {
		t.Run(name, func(t *testing.T) {
			st, err := storage.Open(dir)
			if err != nil {
				t.Fatal(err)
			}
			gs := phase7TestState(t)
			cc, err := openWithProtocolAndClock(st, &gs, 3, c1TestNetworkID, "syjnet-v2-"+c1TestNetworkID, clock.FixedClock{Current: now})
			if err != nil {
				_ = st.Close()
				t.Fatal(err)
			}
			defer st.Close()
			if cc.Height() != wantHeight || cc.TipHash() != wantTip || cc.StateRoot() != wantRoot || cc.Work().String() != wantWork {
				t.Fatalf("replay differs: height=%d tip=%s root=%s work=%s", cc.Height(), cc.TipHash(), cc.StateRoot(), cc.Work())
			}
		})
	}
}

func TestC1V3ProgressionsAndSeededInvalidFuture(t *testing.T) {
	prefix := []*block.Block{{Header: block.Header{Timestamp: 100}}}
	for _, delta := range []float64{1, 30, 120} {
		if err := validateTimestampForProtocol(100+delta, prefix, 3); err != nil {
			t.Fatalf("+%.0fs progression rejected: %v", delta, err)
		}
	}
	const seed int64 = 730001
	r := rand.New(rand.NewSource(seed))
	c := openC1V3(t, t.TempDir(), clock.FixedClock{Current: time.Unix(400, 0)})
	accepted, rejected := 0, 0
	for i := 0; i < 128; i++ {
		ts := 100 + float64(r.Intn(601))
		if err := c.validateV3IngressTimestamp(ts, prefix); err == nil {
			accepted++
		} else {
			rejected++
		}
	}
	t.Logf("seed=%d blocks=128 accepted=%d rejected=%d", seed, accepted, rejected)
}

func TestC1V1V2TimestampCompatibility(t *testing.T) {
	prefix := []*block.Block{{Header: block.Header{Timestamp: 100}}}
	if err := validateTimestampForProtocol(221, prefix, 1); err == nil {
		t.Fatal("V1 parent+120 rule changed")
	}
	if err := validateTimestampForProtocol(221, prefix, 2); err == nil {
		t.Fatal("V2 parent+120 rule changed")
	}
	if err := validateTimestampForProtocol(220, prefix, 1); err != nil {
		t.Fatalf("V1 +120 rejected: %v", err)
	}
	if err := validateTimestampForProtocol(220, prefix, 2); err != nil {
		t.Fatalf("V2 +120 rejected: %v", err)
	}
}
