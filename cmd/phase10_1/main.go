package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/block"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/chain"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/consensus"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/storage"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/internal/transaction"
	"github.com/SHalimoosavi/SAYANJALI-BLOCKCHAIN/pkg/protocol"
)

const (
	defaultBlocks      = 400
	defaultReceiver    = "SYJ0000000000000000000000000000000000000000"
	scheduleFastStep   = int64(6)
	scheduleSlowStep   = int64(120)
	scheduleStableStep = int64(30)
)

type Config struct {
	Blocks   int
	DataDir  string
	JSONL    string
	Report   string
	Receiver string
	RunID    string
}

type BlockRecord struct {
	Height                   int64  `json:"height"`
	Hash                     string `json:"hash"`
	PreviousHash             string `json:"previous_hash"`
	Timestamp                int64  `json:"timestamp"`
	TimestampDelta           int64  `json:"timestamp_delta"`
	Difficulty               int    `json:"difficulty"`
	RequiredTargetPrefix     string `json:"required_target_prefix"`
	MiningNanos              int64  `json:"mining_nanos"`
	ValidationNanos          int64  `json:"validation_nanos"`
	RetargetCalculationNanos int64  `json:"retarget_calculation_nanos,omitempty"`
	CumulativeWork           string `json:"cumulative_work"`
	TransactionCount         int    `json:"transaction_count"`
	StateRoot                string `json:"state_root"`
	MerkleRoot               string `json:"merkle_root"`
	ValidationResult         string `json:"validation_result"`
	ExpectedDifficulty       int    `json:"expected_difficulty"`
	ExpectedDifficultyOK     bool   `json:"expected_difficulty_ok"`
}

type Failure struct {
	RunID       string            `json:"run_id"`
	Command     string            `json:"command"`
	Environment map[string]string `json:"environment"`
	Height      int64             `json:"failure_height,omitempty"`
	BlockHash   string            `json:"failure_block_hash,omitempty"`
	Expected    string            `json:"expected_behavior"`
	Observed    string            `json:"observed_behavior"`
	Error       string            `json:"error"`
}

type Summary struct {
	RunID                string  `json:"run_id"`
	BlocksRequested      int     `json:"blocks_requested"`
	BlocksAccepted       int     `json:"blocks_accepted"`
	GenesisHeight        int64   `json:"genesis_height"`
	FinalHeight          int64   `json:"final_height"`
	FinalHash            string  `json:"final_hash"`
	FinalDifficulty      int     `json:"final_difficulty"`
	DifficultySchedule   []int   `json:"difficulty_schedule"`
	RetargetHeights      []int64 `json:"retarget_heights"`
	TotalMiningNanos     int64   `json:"total_mining_nanos"`
	TotalValidationNanos int64   `json:"total_validation_nanos"`
	FinalCumulativeWork  string  `json:"final_cumulative_work"`
	FinalStateRoot       string  `json:"final_state_root"`
	RestartHeight        int64   `json:"restart_height"`
	RestartHash          string  `json:"restart_hash"`
	RestartStateRoot     string  `json:"restart_state_root"`
	RestartWork          string  `json:"restart_cumulative_work"`
	Deterministic        bool    `json:"deterministic"`
	Status               string  `json:"status"`
}

func main() {
	cfg := Config{}
	flag.IntVar(&cfg.Blocks, "blocks", defaultBlocks, "number of non-genesis blocks to produce")
	flag.StringVar(&cfg.DataDir, "data-dir", "", "temporary chain data directory")
	flag.StringVar(&cfg.JSONL, "jsonl", "", "JSONL output path")
	flag.StringVar(&cfg.Report, "report", "", "human-readable report path")
	flag.StringVar(&cfg.Receiver, "receiver", defaultReceiver, "deterministic valid SYJ receiver")
	flag.StringVar(&cfg.RunID, "run-id", "phase10.1", "experiment run identifier")
	flag.Parse()

	if err := validateConfig(cfg); err != nil {
		fail(cfg, Failure{RunID: cfg.RunID, Command: strings.Join(os.Args, " "), Environment: environment(), Expected: "valid Phase 10.1 configuration", Observed: "invalid configuration", Error: err.Error()})
		os.Exit(2)
	}

	cleanupData := false
	if cfg.DataDir == "" {
		d, err := os.MkdirTemp("", "syj-phase10-1-")
		if err != nil {
			failAndExit(cfg, err, 0, "", "create temporary data directory")
			return
		}
		cfg.DataDir = d
		cleanupData = true
	}
	if cleanupData {
		defer os.RemoveAll(cfg.DataDir)
	}
	if err := os.MkdirAll(filepath.Dir(cfg.DataDir), 0700); err != nil {
		failAndExit(cfg, err, 0, "", "create data parent")
		return
	}

	if cfg.JSONL == "" {
		cfg.JSONL = filepath.Join(cfg.DataDir, "phase10_1_blocks.jsonl")
	}
	if cfg.Report == "" {
		cfg.Report = filepath.Join(cfg.DataDir, "phase10_1_report.md")
	}
	if err := os.MkdirAll(filepath.Dir(cfg.JSONL), 0700); err != nil {
		failAndExit(cfg, err, 0, "", "create JSONL parent")
		return
	}

	result, err := execute(cfg)
	if err != nil {
		failAndExit(cfg, err, resultFailureHeight(result), resultFailureHash(result), "400-block experiment")
		return
	}
	if err := writeSummaryReport(cfg.Report, result.summary); err != nil {
		failAndExit(cfg, err, result.summary.FinalHeight, result.summary.FinalHash, "write report")
		return
	}
	fmt.Printf("PHASE10.1 PASS: %d blocks accepted; final height=%d difficulty=%d tip=%s\n", result.summary.BlocksAccepted, result.summary.FinalHeight, result.summary.FinalDifficulty, result.summary.FinalHash)
	fmt.Printf("JSONL=%s\nREPORT=%s\n", cfg.JSONL, cfg.Report)
}

type execution struct{ summary Summary }

func execute(cfg Config) (execution, error) {
	f, err := os.OpenFile(cfg.JSONL, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return execution{}, err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()

	store, err := storage.Open(cfg.DataDir)
	if err != nil {
		return execution{}, err
	}
	ch, err := chain.Open(store)
	if err != nil {
		store.Close()
		return execution{}, err
	}

	start := ch.Tip()
	if start.Index != 0 {
		return execution{}, fmt.Errorf("fresh experiment expected genesis tip, got height %d", start.Index)
	}
	if err := block.ValidateGenesis(start); err != nil {
		return execution{}, err
	}

	summary := Summary{RunID: cfg.RunID, BlocksRequested: cfg.Blocks, GenesisHeight: start.Index, Status: "RUNNING"}
	var prevRoot = ch.StateRoot()
	var prevTimestamp = int64(start.Timestamp)

	for i := 1; i <= cfg.Blocks; i++ {
		tip := ch.Tip()
		window := ch.ChainCopy()
		expectedDifficulty, err := expectedDifficultyForPrefix(window)
		if err != nil {
			return execution{summary: summary}, fmt.Errorf("height %d expected difficulty: %w", i, err)
		}

		timestamp := timestampForHeight(i, tip.Timestamp)
		reward, ok := consensus.ExpectedReward(ch.Supply())
		if !ok {
			return execution{summary: summary}, fmt.Errorf("height %d: reward exhausted", i)
		}
		coin := transaction.Transaction{Sender: protocol.CoinbaseSender, Receiver: cfg.Receiver, AmountBaseUnits: reward, Timestamp: float64(timestamp)}
		all := []transaction.Transaction{coin}

		constructStart := time.Now()
		b, err := block.New(tip.Index+1, tip.Hash, float64(timestamp), 0, expectedDifficulty, all)
		if err != nil {
			return execution{summary: summary}, fmt.Errorf("height %d construct: %w", i, err)
		}
		if err := mineDeterministicParallel(b, expectedDifficulty); err != nil {
			return execution{summary: summary}, fmt.Errorf("height %d mining: %w", i, err)
		}
		miningNanos := time.Since(constructStart).Nanoseconds()

		if b.Difficulty != expectedDifficulty {
			return execution{summary: summary}, fmt.Errorf("height %d difficulty mismatch before accept: got %d want %d", i, b.Difficulty, expectedDifficulty)
		}
		if b.PreviousHash != tip.Hash {
			return execution{summary: summary}, fmt.Errorf("height %d previous hash mismatch before accept", i)
		}
		if int64(b.Timestamp) != timestamp {
			return execution{summary: summary}, fmt.Errorf("height %d timestamp mismatch before accept", i)
		}

		validateStart := time.Now()
		accepted, reason, acceptErr := ch.Accept(b)
		validationNanos := time.Since(validateStart).Nanoseconds()
		if acceptErr != nil || !accepted || reason != "best" {
			return execution{summary: summary}, fmt.Errorf("height %d accept failed: accepted=%v reason=%s err=%v", i, accepted, reason, acceptErr)
		}

		active := ch.ChainCopy()
		if len(active) != i+1 {
			return execution{summary: summary}, fmt.Errorf("height %d active chain length=%d", i, len(active))
		}
		if active[len(active)-1].Hash != b.Hash {
			return execution{summary: summary}, fmt.Errorf("height %d active tip mismatch", i)
		}
		if b.Timestamp <= tip.Timestamp {
			return execution{summary: summary}, fmt.Errorf("height %d timestamp not increasing", i)
		}
		if b.Timestamp-tip.Timestamp > float64(chain.TimestampFutureStepSeconds) {
			return execution{summary: summary}, fmt.Errorf("height %d future-time bound exceeded", i)
		}
		if len(b.Transactions) != 1 || block.MerkleRoot([]string{b.Transactions[0].TxHash}) != b.MerkleRoot {
			return execution{summary: summary}, fmt.Errorf("height %d merkle root mismatch", i)
		}
		root := ch.StateRoot()
		if root == "" {
			return execution{summary: summary}, fmt.Errorf("height %d empty state root", i)
		}
		if root == prevRoot {
			return execution{summary: summary}, fmt.Errorf("height %d state root did not progress", i)
		}

		var retargetNanos int64
		if b.Index%protocol.DifficultyInterval == 0 {
			retargetStart := time.Now()
			windowEnd := window[len(window)-protocol.DifficultyInterval:]
			observed, err := consensus.NextDifficulty(windowEnd, protocol.DefaultDifficultyConfig(), 4)
			retargetNanos = time.Since(retargetStart).Nanoseconds()
			if err != nil {
				return execution{summary: summary}, fmt.Errorf("height %d retarget calculation: %w", i, err)
			}
			if observed != b.Difficulty {
				return execution{summary: summary}, fmt.Errorf("height %d retarget mismatch: got %d want %d", i, b.Difficulty, observed)
			}
			summary.RetargetHeights = append(summary.RetargetHeights, b.Index)
		}

		work := ch.Work()
		if work.Sign() <= 0 {
			return execution{summary: summary}, fmt.Errorf("height %d cumulative work is not positive", i)
		}
		if i > 1 {
			priorWork, ok := new(big.Int).SetString(summary.FinalCumulativeWork, 10)
			if !ok || work.Cmp(priorWork) <= 0 {
				return execution{summary: summary}, fmt.Errorf("height %d cumulative work did not increase", i)
			}
		}

		rec := BlockRecord{
			Height: b.Index, Hash: b.Hash, PreviousHash: b.PreviousHash, Timestamp: timestamp,
			TimestampDelta: timestamp - prevTimestamp, Difficulty: b.Difficulty,
			RequiredTargetPrefix: strings.Repeat("0", b.Difficulty), MiningNanos: miningNanos,
			ValidationNanos: validationNanos, RetargetCalculationNanos: retargetNanos,
			CumulativeWork: work.String(), TransactionCount: len(b.Transactions), StateRoot: root,
			MerkleRoot: b.MerkleRoot, ValidationResult: "PASS", ExpectedDifficulty: expectedDifficulty, ExpectedDifficultyOK: true,
		}
		enc, err := json.Marshal(rec)
		if err != nil {
			return execution{summary: summary}, err
		}
		if _, err := w.Write(append(enc, '\n')); err != nil {
			return execution{summary: summary}, err
		}

		summary.BlocksAccepted++
		summary.FinalHeight = b.Index
		summary.FinalHash = b.Hash
		summary.FinalDifficulty = b.Difficulty
		summary.DifficultySchedule = append(summary.DifficultySchedule, b.Difficulty)
		summary.TotalMiningNanos += miningNanos
		summary.TotalValidationNanos += validationNanos
		summary.FinalCumulativeWork = work.String()
		summary.FinalStateRoot = root
		prevRoot, prevTimestamp = root, timestamp
		_ = constructStart
	}
	if err := w.Flush(); err != nil {
		return execution{summary: summary}, err
	}
	if err := chain.ValidateChain(ch.ChainCopy()); err != nil {
		return execution{summary: summary}, fmt.Errorf("final full chain validation: %w", err)
	}

	// Close and reopen the real persistent store. This validates journal reload,
	// chain reconstruction, incremental state restoration and state-root recovery.
	if err := store.Close(); err != nil {
		return execution{summary: summary}, err
	}
	reloadedStore, err := storage.Open(cfg.DataDir)
	if err != nil {
		return execution{summary: summary}, fmt.Errorf("restart store open: %w", err)
	}
	reloaded, err := chain.Open(reloadedStore)
	if err != nil {
		reloadedStore.Close()
		return execution{summary: summary}, fmt.Errorf("restart chain open: %w", err)
	}
	summary.RestartHeight = reloaded.Height()
	summary.RestartHash = reloaded.TipHash()
	summary.RestartStateRoot = reloaded.StateRoot()
	summary.RestartWork = reloaded.Work().String()
	if summary.RestartHeight != summary.FinalHeight || summary.RestartHash != summary.FinalHash || summary.RestartStateRoot != summary.FinalStateRoot || summary.RestartWork != summary.FinalCumulativeWork {
		return execution{summary: summary}, errors.New("restart/reload result differs from final in-memory chain")
	}
	if err := chain.ValidateChain(reloaded.ChainCopy()); err != nil {
		return execution{summary: summary}, fmt.Errorf("restart full validation: %w", err)
	}
	if err := reloadedStore.Close(); err != nil {
		return execution{summary: summary}, err
	}

	summary.Deterministic = true
	summary.Status = "PASS"
	return execution{summary: summary}, nil
}

func mineDeterministicParallel(b *block.Block, difficulty int) error {
	// Harness-only acceleration: construct exactly the canonical JSON transcript
	// used by production block.HashHeader, then hash it with the standard
	// library SHA-256. The transcript is checked against block.HashHeader at
	// nonce zero before mining. chain.Accept revalidates the selected header
	// through the production consensus path, so an acceleration bug cannot
	// cause a non-production block to be accepted.
	workers := runtime.GOMAXPROCS(0)
	if workers > 8 {
		workers = 8
	}
	if workers < 1 {
		workers = 1
	}
	prefix, suffix := fastHeaderParts(b.Header)
	checkHash, _, err := block.HashHeader(b.Header)
	if err != nil {
		return err
	}
	if fastHash(prefix, suffix, 0) != checkHash {
		return errors.New("accelerated header transcript does not match production block.HashHeader")
	}

	const batchSize uint64 = 1 << 16
	for batchStart := uint64(0); batchStart <= ^uint64(0)-batchSize; batchStart += batchSize {
		var wg sync.WaitGroup
		var mu sync.Mutex
		var firstErr error
		best := ^uint64(0)
		for worker := 0; worker < workers; worker++ {
			start := batchStart + uint64(worker)*batchSize/uint64(workers)
			end := batchStart + uint64(worker+1)*batchSize/uint64(workers)
			wg.Add(1)
			go func(start, end uint64) {
				defer wg.Done()
				buf := make([]byte, 0, len(prefix)+20+len(suffix))
				for nonce := start; nonce < end; nonce++ {
					buf = buf[:0]
					buf = append(buf, prefix...)
					buf = strconv.AppendUint(buf, nonce, 10)
					buf = append(buf, suffix...)
					sum := sha256.Sum256(buf)
					if meetsTargetBytes(sum, difficulty) {
						mu.Lock()
						if nonce < best {
							best = nonce
						}
						mu.Unlock()
					}
				}
			}(start, end)
		}
		wg.Wait()
		if firstErr != nil {
			return firstErr
		}
		if best != ^uint64(0) {
			b.Nonce = best
			h, _, err := block.HashHeader(b.Header)
			if err != nil {
				return err
			}
			if !meetsTargetBytesHash(h, difficulty) {
				return errors.New("selected nonce does not satisfy production proof-of-work")
			}
			b.Hash = h
			return nil
		}
		if batchStart > ^uint64(0)-batchSize {
			break
		}
	}
	return errors.New("nonce exhausted")
}

func meetsTargetBytes(sum [32]byte, difficulty int) bool {
	fullBytes := difficulty / 2
	for i := 0; i < fullBytes; i++ {
		if sum[i] != 0 {
			return false
		}
	}
	if difficulty%2 == 1 {
		return sum[fullBytes]&0xf0 == 0
	}
	return true
}

func meetsTargetBytesHash(hash string, difficulty int) bool {
	if len(hash) < difficulty {
		return false
	}
	return strings.HasPrefix(hash, strings.Repeat("0", difficulty))
}

func fastHeaderParts(h block.Header) ([]byte, []byte) {
	// canonicaljson sorts these six keys lexicographically. All string fields
	// in a block header are protocol hashes, so the production escaping rules
	// have no special characters to encode.
	prefix := fmt.Sprintf(`{"difficulty":%d,"index":%d,"merkle_root":"%s","nonce":`, h.Difficulty, h.Index, h.MerkleRoot)
	timestamp := strconv.FormatFloat(h.Timestamp, 'f', -1, 64)
	if !strings.Contains(timestamp, ".") {
		timestamp += ".0"
	}
	suffix := fmt.Sprintf(`,"previous_hash":"%s","timestamp":%s}`, h.PreviousHash, timestamp)
	return []byte(prefix), []byte(suffix)
}

func fastHash(prefix, suffix []byte, nonce uint64) string {
	buf := make([]byte, 0, len(prefix)+20+len(suffix))
	buf = append(buf, prefix...)
	buf = strconv.AppendUint(buf, nonce, 10)
	buf = append(buf, suffix...)
	sum := sha256.Sum256(buf)
	const hex = "0123456789abcdef"
	out := make([]byte, 64)
	for i, v := range sum {
		out[i*2] = hex[v>>4]
		out[i*2+1] = hex[v&15]
	}
	return string(out)
}

func expectedDifficultyForPrefix(prefix []*block.Block) (int, error) {
	if len(prefix) == 0 {
		return 4, nil
	}
	base := 4
	if prefix[len(prefix)-1].Index > 0 {
		base = prefix[len(prefix)-1].Difficulty
	}
	nextIndex := int64(len(prefix))
	if len(prefix) < protocol.DifficultyInterval || nextIndex%protocol.DifficultyInterval != 0 {
		return base, nil
	}
	recent := prefix[len(prefix)-protocol.DifficultyInterval:]
	return consensus.NextDifficulty(recent, protocol.DefaultDifficultyConfig(), 4)
}

func timestampForHeight(height int, parent float64) int64 {
	// Every epoch has nine parent-to-child intervals. Alternating 6-second and
	// 120-second intervals exercises the frozen four-times retarget clamp in
	// both directions. Neither value exceeds the production +120s bound.
	epoch := (height - 1) / protocol.DifficultyInterval
	step := scheduleStableStep
	switch epoch % 3 {
	case 0:
		step = scheduleSlowStep
	case 1:
		step = scheduleFastStep
	case 2:
		step = scheduleStableStep
	}
	return int64(parent) + step
}

func validateConfig(cfg Config) error {
	if cfg.Blocks < 1 {
		return errors.New("blocks must be positive")
	}
	if cfg.Blocks < 2*protocol.DifficultyInterval {
		return errors.New("blocks must cross multiple retarget boundaries")
	}
	if len(cfg.Receiver) != len(protocol.AddressPrefix)+protocol.AddressHashLength || !strings.HasPrefix(cfg.Receiver, protocol.AddressPrefix) {
		return errors.New("receiver must be a valid deterministic SYJ address")
	}
	for _, c := range cfg.Receiver[len(protocol.AddressPrefix):] {
		if !strings.ContainsRune("0123456789abcdefABCDEF", c) {
			return errors.New("receiver contains non-hex address data")
		}
	}
	return nil
}

func environment() map[string]string {
	return map[string]string{"go_version": runtime.Version(), "go_os": runtime.GOOS, "go_arch": runtime.GOARCH}
}
func fail(cfg Config, f Failure) {
	path := cfg.Report
	if path == "" {
		path = "phase10_1_failure.json"
	}
	_ = os.MkdirAll(filepath.Dir(path), 0700)
	data, _ := json.MarshalIndent(f, "", "  ")
	_ = os.WriteFile(path+".failure.json", data, 0600)
}
func failAndExit(cfg Config, err error, height int64, hash, expected string) {
	fail(cfg, Failure{RunID: cfg.RunID, Command: strings.Join(os.Args, " "), Environment: environment(), Height: height, BlockHash: hash, Expected: expected, Observed: err.Error(), Error: err.Error()})
	fmt.Fprintf(os.Stderr, "PHASE10.1 FAIL: %v\n", err)
	os.Exit(1)
}
func resultFailureHeight(e execution) int64 { return e.summary.FinalHeight + 1 }
func resultFailureHash(e execution) string  { return e.summary.FinalHash }

func writeSummaryReport(path string, s Summary) error {
	var b strings.Builder
	b.WriteString("# Phase 10.1 — 400-Block Consensus & Difficulty Regression\n\n")
	b.WriteString("Status: **PASS**\n\n")
	b.WriteString("This report is generated by the real Go production chain path. It is not a toy simulator and does not alter production consensus constants.\n\n")
	fmt.Fprintf(&b, "- Blocks accepted: `%d`\n- Final height: `%d`\n- Final tip: `%s`\n- Final difficulty: `%d`\n- Retarget boundaries crossed: `%d`\n- Final cumulative work: `%s`\n- Final state root: `%s`\n- Restart height: `%d`\n- Restart tip: `%s`\n- Deterministic consensus result: `%v`\n\n", s.BlocksAccepted, s.FinalHeight, s.FinalHash, s.FinalDifficulty, len(s.RetargetHeights), s.FinalCumulativeWork, s.FinalStateRoot, s.RestartHeight, s.RestartHash, s.Deterministic)
	b.WriteString("## Retarget heights\n\n")
	for _, h := range s.RetargetHeights {
		fmt.Fprintf(&b, "- `%d`\n", h)
	}
	b.WriteString("\n## Performance\n\n")
	fmt.Fprintf(&b, "- Total mining/production time: `%s`\n- Total validation time: `%s`\n- Average mining/production time per block: `%s`\n- Average validation time per block: `%s`\n\n", time.Duration(s.TotalMiningNanos), time.Duration(s.TotalValidationNanos), time.Duration(s.TotalMiningNanos/int64(max(1, s.BlocksAccepted))), time.Duration(s.TotalValidationNanos/int64(max(1, s.BlocksAccepted))))
	return os.WriteFile(path, []byte(b.String()), 0600)
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
