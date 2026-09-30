package main

import (
	"strings"
	"testing"

	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/block"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/pkg/protocol"
)

func TestValidateConfigRequiresMultipleRetargets(t *testing.T) {
	cfg := Config{Blocks: 1, Receiver: defaultReceiver}
	if err := validateConfig(cfg); err == nil {
		t.Fatal("expected insufficient block count to fail")
	}
}

func TestTimestampScheduleStaysWithinConsensusFutureBound(t *testing.T) {
	parent := float64(protocol.GenesisTimestamp)
	for height := 1; height <= 400; height++ {
		next := timestampForHeight(height, parent)
		if float64(next)-parent <= 0 || float64(next)-parent > float64(4*protocol.TargetBlockTimeSeconds) {
			t.Fatalf("height %d schedule step=%d violates future bound", height, next-int64(parent))
		}
		parent = float64(next)
	}
}

func TestTimestampScheduleCrossesFastSlowAndStableEpochs(t *testing.T) {
	parent := float64(protocol.GenesisTimestamp)
	steps := map[int64]bool{}
	for height := 1; height <= 40; height++ {
		next := timestampForHeight(height, parent)
		steps[next-int64(parent)] = true
		parent = float64(next)
	}
	for _, want := range []int64{scheduleFastStep, scheduleSlowStep, scheduleStableStep} {
		if !steps[want] {
			t.Fatalf("missing schedule step %d", want)
		}
	}
}

func TestFastHashMatchesProductionHeaderHash(t *testing.T) {
	b, err := block.New(1, strings.Repeat("a", 64), protocol.GenesisTimestamp+120, 0, 4, nil)
	if err != nil {
		t.Fatal(err)
	}
	prefix, suffix := fastHeaderParts(b.Header)
	want, _, err := block.HashHeader(b.Header)
	if err != nil {
		t.Fatal(err)
	}
	if got := fastHash(prefix, suffix, 0); got != want {
		t.Fatalf("accelerated header hash mismatch: got %s want %s", got, want)
	}
}
