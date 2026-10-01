package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/block"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/chain"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/clock"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/storage"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/tokenomics"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/transaction"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/wallet"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/pkg/protocol"
)

const seed int64 = 20261001
const networkID = "3333333333333333333333333333333333333333333333333333333333333333"
const receiver = "SYJaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

func address(seed byte) string {
	var pub [64]byte
	for i := range pub {
		pub[i] = seed
	}
	return wallet.AddressFromPublicKeyHex(fmt.Sprintf("%x", pub[:]))
}

func genesisState() tokenomics.GenesisState {
	return tokenomics.GenesisState{Version: tokenomics.Version, Sender: protocol.GenesisAllocationSender, Allocations: []tokenomics.Allocation{
		{Category: tokenomics.Presale, Recipient: address(1), AmountBaseUnits: 7_200_000_000_000_000},
		{Category: tokenomics.Treasury, Recipient: address(2), AmountBaseUnits: 5_760_000_000_000_000},
		{Category: tokenomics.Ecosystem, Recipient: address(3), AmountBaseUnits: 3_600_000_000_000_000},
		{Category: tokenomics.Liquidity, Recipient: address(4), AmountBaseUnits: 4_320_000_000_000_000},
		{Category: tokenomics.Team, Recipient: address(5), AmountBaseUnits: 7_920_000_000_000_000},
	}, TotalBaseUnits: tokenomics.ExpectedGenesisTotal()}
}

func mine(parent *block.Block, network string, timestamp float64) (*block.Block, error) {
	tx, err := transaction.NewV2Coinbase(receiver, network, 5_000_000_000, timestamp)
	if err != nil {
		return nil, err
	}
	b, err := block.New(parent.Index+1, parent.Hash, timestamp, 0, 4, []transaction.Transaction{tx})
	if err != nil {
		return nil, err
	}
	for b.Nonce < ^uint64(0) {
		if b.MeetsDifficulty(4) {
			return b, nil
		}
		b.Nonce++
		if err := b.Recompute(); err != nil {
			return nil, err
		}
	}
	return nil, fmt.Errorf("nonce exhausted")
}

func main() {
	start := time.Now()
	dir, err := os.MkdirTemp("", "syj-build1-c1-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	st, err := storage.Open(dir)
	if err != nil {
		panic(err)
	}
	gs := genesisState()
	fixed := time.Unix(int64(protocol.GenesisTimestamp)+1000, 0)
	c, err := chain.OpenV3WithGenesisState(st, gs, networkID, "syjnet-v2-"+networkID)
	if err != nil {
		panic(err)
	}
	c.SetClock(clock.FixedClock{Current: fixed})

	r := rand.New(rand.NewSource(seed))
	accepted, rejected := 0, 0
	difficultyObservations := make([]int, 0, 64)
	for i := 1; i <= 64; i++ {
		ts := fixed.Unix() + int64(i)
		if r.Intn(10) == 0 {
			ts = fixed.Unix() + 301 + int64(r.Intn(200))
		}
		b, err := mine(c.Tip(), networkID, float64(ts))
		if err != nil {
			panic(err)
		}
		ok, _, err := c.Accept(b)
		if err != nil || !ok {
			rejected++
			continue
		}
		accepted++
		difficultyObservations = append(difficultyObservations, b.Difficulty)
	}

	finalHeight := c.Height()
	finalTip := c.TipHash()
	finalStateRoot := c.StateRoot()
	finalWork := c.Work().String()

	if err := st.Close(); err != nil {
		panic(err)
	}

	// Replay check: reopen the persisted chain with a clock decades ahead.
	replayStore, err := storage.Open(dir)
	if err != nil {
		panic(err)
	}
	replayClock := clock.FixedClock{Current: fixed.Add(10 * 365 * 24 * time.Hour)}
	replayChain, err := chain.OpenV3WithGenesisState(replayStore, gs, networkID, "syjnet-v2-"+networkID)
	if err != nil {
		_ = replayStore.Close()
		panic(err)
	}
	replayChain.SetClock(replayClock)

	replayOK := replayChain.Height() == finalHeight &&
		replayChain.TipHash() == finalTip &&
		replayChain.StateRoot() == finalStateRoot &&
		replayChain.Work().String() == finalWork

	if err := replayStore.Close(); err != nil {
		panic(err)
	}

	// Restart check: reopen again with a clock decades behind and verify the
	// persisted consensus state remains identical.
	restartStore, err := storage.Open(dir)
	if err != nil {
		panic(err)
	}
	restartClock := clock.FixedClock{Current: fixed.Add(-10 * 365 * 24 * time.Hour)}
	restartChain, err := chain.OpenV3WithGenesisState(restartStore, gs, networkID, "syjnet-v2-"+networkID)
	if err != nil {
		_ = restartStore.Close()
		panic(err)
	}
	restartChain.SetClock(restartClock)

	restartOK := restartChain.Height() == finalHeight &&
		restartChain.TipHash() == finalTip &&
		restartChain.StateRoot() == finalStateRoot &&
		restartChain.Work().String() == finalWork

	if err := restartStore.Close(); err != nil {
		panic(err)
	}

	replayResult := "FAIL"
	if replayOK {
		replayResult = "PASS"
	}
	restartResult := "FAIL"
	if restartOK {
		restartResult = "PASS"
	}

	fmt.Printf("seed=%d\nblocks=%d\naccepted=%d\nrejected=%d\nelapsed=%s\ndifficulty_observations=%v\ncumulative_work=%s\nfinal_height=%d\nfinal_tip=%s\nstate_root=%s\nreplay_result=%s\nrestart_result=%s\n", seed, 64, accepted, rejected, time.Since(start), difficultyObservations, finalWork, finalHeight, finalTip, finalStateRoot, replayResult, restartResult)
	if !replayOK || !restartOK {
		os.Exit(1)
	}
}
