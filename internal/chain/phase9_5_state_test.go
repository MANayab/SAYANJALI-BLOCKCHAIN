package chain

import (
	"testing"

	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/block"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/consensus"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/networkid"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/statecommitment"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/storage"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/tokenomics"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/transaction"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/wallet"
)

func cloneUint64Map(in map[string]uint64) map[string]uint64 {
	out := make(map[string]uint64, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func snapshotFromChainState(c *Chain) *statecommitment.Snapshot {
	s := statecommitment.NewSnapshot()

	s.Balances = make(map[string]uint64, len(c.balances))
	for k, v := range c.balances {
		s.Balances[k] = v
	}

	s.Nonces = cloneUint64Map(c.nonces)
	s.GenesisSupply = c.genesisSupply
	s.MiningIssued = c.miningIssued
	s.Supply = c.supply

	return s
}

func TestPhase9_5IncrementalStateMatchesLegacyReplay(t *testing.T) {
	t.Parallel()

	gs := phase7TestState(t)

	dir := t.TempDir()
	st, err := storage.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = st.Close()
	})

	c, err := OpenWithGenesisState(st, gs)
	if err != nil {
		t.Fatal(err)
	}

	// The incremental initialization performed by OpenWithGenesisState
	// must already have reconstructed the same economic state as the
	// legacy replay path.
	if len(c.active) < 1 {
		t.Fatal("active chain is empty")
	}

	incremental, ok := c.stateSnapshots[c.activeTip]
	if !ok {
		t.Fatal("missing incremental state snapshot at active tip")
	}

	if incremental.Root() != c.stateRoot {
		t.Fatalf(
			"state root mismatch inside incremental state: snapshot=%s chain=%s",
			incremental.Root(),
			c.stateRoot,
		)
	}

	// Capture the incremental result before invoking the legacy replay.
	incrementalCopy := incremental.Clone()

	// Reconstruct the legacy state from the exact same active chain.
	if err := c.replayState(c.active); err != nil {
		t.Fatalf("legacy replay failed: %v", err)
	}

	legacy := snapshotFromChainState(c)

	if got, want := legacy.GenesisSupply, incrementalCopy.GenesisSupply; got != want {
		t.Fatalf("genesis supply mismatch: legacy=%d incremental=%d", got, want)
	}

	if got, want := legacy.MiningIssued, incrementalCopy.MiningIssued; got != want {
		t.Fatalf("mining issued mismatch: legacy=%d incremental=%d", got, want)
	}

	if got, want := legacy.Supply, incrementalCopy.Supply; got != want {
		t.Fatalf("supply mismatch: legacy=%d incremental=%d", got, want)
	}

	if len(legacy.Balances) != len(incrementalCopy.Balances) {
		t.Fatalf(
			"balance-account count mismatch: legacy=%d incremental=%d",
			len(legacy.Balances),
			len(incrementalCopy.Balances),
		)
	}

	for addr, want := range incrementalCopy.Balances {
		if got := legacy.Balances[addr]; got != want {
			t.Fatalf(
				"balance mismatch for %s: legacy=%d incremental=%d",
				addr,
				got,
				want,
			)
		}
	}

	if len(legacy.Nonces) != len(incrementalCopy.Nonces) {
		t.Fatalf(
			"nonce-account count mismatch: legacy=%d incremental=%d",
			len(legacy.Nonces),
			len(incrementalCopy.Nonces),
		)
	}

	for addr, want := range incrementalCopy.Nonces {
		if got := legacy.Nonces[addr]; got != want {
			t.Fatalf(
				"nonce mismatch for %s: legacy=%d incremental=%d",
				addr,
				got,
				want,
			)
		}
	}

	if got, want := legacy.Root(), incrementalCopy.Root(); got != want {
		t.Fatalf(
			"state root mismatch: legacy=%s incremental=%s",
			got,
			want,
		)
	}
}

func TestPhase9_5IncrementalStateSurvivesCheckpointRestart(t *testing.T) {
	gs := phase7TestState(t)

	dir := t.TempDir()

	st, err := storage.Open(dir)
	if err != nil {
		t.Fatal(err)
	}

	c, err := OpenWithGenesisState(st, gs)
	if err != nil {
		_ = st.Close()
		t.Fatal(err)
	}

	if len(c.active) < 1 {
		_ = st.Close()
		t.Fatal("active chain is empty")
	}

	tip := c.activeTip
	first := c.stateSnapshots[tip]
	if first == nil {
		_ = st.Close()
		t.Fatal("missing initial state snapshot")
	}

	wantRoot := first.Root()

	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	st2, err := storage.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = st2.Close()
	})

	c2, err := OpenWithGenesisState(st2, gs)
	if err != nil {
		t.Fatal(err)
	}

	if c2.activeTip != tip {
		t.Fatalf(
			"active tip changed across restart: before=%s after=%s",
			tip,
			c2.activeTip,
		)
	}

	after := c2.stateSnapshots[c2.activeTip]
	if after == nil {
		t.Fatal("missing state snapshot after restart")
	}

	if got := after.Root(); got != wantRoot {
		t.Fatalf(
			"state root changed across restart: before=%s after=%s",
			wantRoot,
			got,
		)
	}

	if got := c2.stateRoot; got != wantRoot {
		t.Fatalf(
			"chain state root changed across restart: expected=%s got=%s",
			wantRoot,
			got,
		)
	}

	// Confirm the restarted incremental state also matches the
	// independent legacy replay path.
	if err := c2.replayState(c2.active); err != nil {
		t.Fatalf("legacy replay after restart failed: %v", err)
	}

	legacy := snapshotFromChainState(c2)

	if got := legacy.Root(); got != wantRoot {
		t.Fatalf(
			"legacy replay root differs after restart: expected=%s got=%s",
			wantRoot,
			got,
		)
	}
}

func TestPhase9_5StateSnapshotDeterministic(t *testing.T) {
	gs := phase7TestState(t)

	s1 := statecommitment.NewSnapshot()
	s2 := statecommitment.NewSnapshot()

	for _, allocation := range gs.Allocations {
		s1.Balances[allocation.Recipient] += allocation.AmountBaseUnits
		s2.Balances[allocation.Recipient] += allocation.AmountBaseUnits
	}

	s1.GenesisSupply = tokenomics.ExpectedGenesisTotal()
	s2.GenesisSupply = tokenomics.ExpectedGenesisTotal()
	s1.Supply = s1.GenesisSupply
	s2.Supply = s2.GenesisSupply

	if s1.Root() != s2.Root() {
		t.Fatalf("identical snapshots produced different roots")
	}

	// Exercise insertion-order independence explicitly.
	k1, err := wallet.New()
	if err != nil {
		t.Fatal(err)
	}
	k2, err := wallet.New()
	if err != nil {
		t.Fatal(err)
	}

	a := statecommitment.NewSnapshot()
	b := statecommitment.NewSnapshot()

	a.Balances[k1.Address] = 100
	a.Balances[k2.Address] = 200

	b.Balances[k2.Address] = 200
	b.Balances[k1.Address] = 100

	a.GenesisSupply = 300
	b.GenesisSupply = 300
	a.Supply = 300
	b.Supply = 300

	if a.Root() != b.Root() {
		t.Fatalf("state root depends on account insertion order")
	}
}

func minePhase9_5V2Block(t *testing.T, b *block.Block) *block.Block {
	t.Helper()

	const maxNonce = uint64(^uint32(0))

	for nonce := uint64(0); nonce <= maxNonce; nonce++ {
		b.Nonce = nonce

		hash, _, err := block.HashHeader(b.Header)
		if err != nil {
			t.Fatal(err)
		}

		b.Hash = hash

		if err := consensus.ValidatePoW(b, b.Difficulty); err == nil {
			return b
		}
	}

	t.Fatal("unable to mine Phase 9.5 V2 test block")
	return nil
}

func TestPhase9_5V2ReorgIncrementalStateMatchesLegacyReplay(t *testing.T) {
	gs, k, network := v2TestState(t)

	dir := t.TempDir()

	st, err := storage.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = st.Close()
	})

	networkName, err := networkid.WireNetworkName(network)
	if err != nil {
		t.Fatal(err)
	}

	c, err := OpenV2WithGenesisState(
		st,
		gs,
		network,
		networkName,
	)
	if err != nil {
		t.Fatal(err)
	}

	g := c.Tip()

	reward, ok := consensus.ExpectedMiningReward(0)
	if !ok {
		t.Fatal("expected initial mining reward")
	}

	receiverA := "SYJaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	receiverB := "SYJbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"

	// Branch A spends one unit from the Phase 7 allocation.
	txA := transaction.Transaction{
		Version:         2,
		NetworkID:       network,
		Sender:          k.Address,
		Receiver:        receiverA,
		AmountBaseUnits: 1,
		Timestamp:       float64(g.Timestamp + 1),
		Nonce:           0,
	}
	if err := txA.SignV2(k, network); err != nil {
		t.Fatal(err)
	}

	coinA, err := transaction.NewV2Coinbase(
		receiverA,
		network,
		reward,
		g.Timestamp+1,
	)
	if err != nil {
		t.Fatal(err)
	}

	a1, err := block.New(
		1,
		g.Hash,
		g.Timestamp+1,
		0,
		4,
		[]transaction.Transaction{coinA, txA},
	)
	if err != nil {
		t.Fatal(err)
	}
	a1 = minePhase9_5V2Block(t, a1)

	if ok, reason, err := c.Accept(a1); err != nil || !ok || reason != "best" {
		t.Fatalf("branch A accept: ok=%v reason=%s err=%v", ok, reason, err)
	}

	rootA := c.stateRoot
	snapshotA := c.stateSnapshots[a1.Hash]
	if snapshotA == nil {
		t.Fatal("missing branch A state snapshot")
	}

	// Branch B starts from genesis and spends the same source balance
	// to a different receiver. It is initially a fork.
	txB := transaction.Transaction{
		Version:         2,
		NetworkID:       network,
		Sender:          k.Address,
		Receiver:        receiverB,
		AmountBaseUnits: 1,
		Timestamp:       float64(g.Timestamp + 2),
		Nonce:           0,
	}
	if err := txB.SignV2(k, network); err != nil {
		t.Fatal(err)
	}

	coinB, err := transaction.NewV2Coinbase(
		receiverB,
		network,
		reward,
		g.Timestamp+2,
	)
	if err != nil {
		t.Fatal(err)
	}

	b1, err := block.New(
		1,
		g.Hash,
		g.Timestamp+2,
		0,
		4,
		[]transaction.Transaction{coinB, txB},
	)
	if err != nil {
		t.Fatal(err)
	}
	b1 = minePhase9_5V2Block(t, b1)

	if ok, reason, err := c.Accept(b1); err != nil || !ok || reason != "fork" {
		t.Fatalf("branch B accept: ok=%v reason=%s err=%v", ok, reason, err)
	}

	// Extend B so its cumulative work exceeds branch A.
	coinB2, err := transaction.NewV2Coinbase(
		receiverB,
		network,
		reward,
		b1.Timestamp+1,
	)
	if err != nil {
		t.Fatal(err)
	}

	b2, err := block.New(
		2,
		b1.Hash,
		b1.Timestamp+1,
		0,
		4,
		[]transaction.Transaction{coinB2},
	)
	if err != nil {
		t.Fatal(err)
	}
	b2 = minePhase9_5V2Block(t, b2)

	if ok, reason, err := c.Accept(b2); err != nil || !ok || reason != "best" {
		t.Fatalf("branch B reorg: ok=%v reason=%s err=%v", ok, reason, err)
	}

	if got := c.TipHash(); got != b2.Hash {
		t.Fatalf("unexpected active tip after reorg: got=%s want=%s", got, b2.Hash)
	}

	if got := c.Height(); got != 2 {
		t.Fatalf("unexpected active height after reorg: got=%d want=2", got)
	}

	if rootA == c.stateRoot {
		t.Fatal("state root did not change across divergent reorg state")
	}

	if snapshotA.Root() != rootA {
		t.Fatal("branch A snapshot root changed unexpectedly")
	}

	winning, ok := c.stateSnapshots[b2.Hash]
	if !ok {
		t.Fatal("missing winning branch state snapshot")
	}

	if got := winning.Root(); got != c.stateRoot {
		t.Fatalf(
			"winning snapshot/root mismatch: snapshot=%s chain=%s",
			got,
			c.stateRoot,
		)
	}

	// The winning branch must have B's state, not A's state.
	// receiverB starts with the Phase 7 genesis allocation, then receives
	// the one-unit transaction and the block-1 mining reward.
	expectedReceiverBBalance := uint64(5_000_000_000) + 1 + reward
	if got := winning.Balances[receiverB]; got != expectedReceiverBBalance {
		t.Fatalf(
			"winning receiver B balance mismatch: got=%d want=%d",
			got,
			expectedReceiverBBalance,
		)
	}

	if got := winning.Balances[receiverA]; got != 0 {
		t.Fatalf(
			"losing receiver A balance remained in winning state: got=%d",
			got,
		)
	}

	if got := winning.Nonces[k.Address]; got != 1 {
		t.Fatalf(
			"winning sender nonce mismatch: got=%d want=1",
			got,
		)
	}

	if got := winning.MiningIssued; got != 2*reward {
		t.Fatalf(
			"winning mining issuance mismatch: got=%d want=%d",
			got,
			2*reward,
		)
	}

	// Reconstruct the exact active chain using the legacy reference path.
	if err := c.replayState(c.active); err != nil {
		t.Fatalf("legacy replay after reorg failed: %v", err)
	}

	legacy := snapshotFromChainState(c)

	if got, want := legacy.Root(), winning.Root(); got != want {
		t.Fatalf(
			"legacy/incremental reorg root mismatch: legacy=%s incremental=%s",
			got,
			want,
		)
	}

	if got, want := legacy.Balances[receiverB], winning.Balances[receiverB]; got != want {
		t.Fatalf(
			"legacy/incremental receiver B mismatch: legacy=%d incremental=%d",
			got,
			want,
		)
	}

	if got, want := legacy.Balances[receiverA], winning.Balances[receiverA]; got != want {
		t.Fatalf(
			"legacy/incremental receiver A mismatch: legacy=%d incremental=%d",
			got,
			want,
		)
	}

	if got, want := legacy.Nonces[k.Address], winning.Nonces[k.Address]; got != want {
		t.Fatalf(
			"legacy/incremental nonce mismatch: legacy=%d incremental=%d",
			got,
			want,
		)
	}

	if got, want := legacy.MiningIssued, winning.MiningIssued; got != want {
		t.Fatalf(
			"legacy/incremental mining issuance mismatch: legacy=%d incremental=%d",
			got,
			want,
		)
	}

	if got, want := legacy.Supply, winning.Supply; got != want {
		t.Fatalf(
			"legacy/incremental supply mismatch: legacy=%d incremental=%d",
			got,
			want,
		)
	}
}
