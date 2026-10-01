package protocol

import "time"

const (
	// V3FutureTimeBound is the maximum wall-clock future allowance applied
	// only when a V3 block first enters a node.
	V3FutureTimeBound time.Duration = 300 * time.Second

	// V3MiningClockSkewLimit is an operational mining-safety threshold. It is
	// not a consensus validity rule for received blocks.
	V3MiningClockSkewLimit time.Duration = 60 * time.Second
)
