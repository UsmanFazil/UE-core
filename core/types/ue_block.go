package types

import (
	"time"
)

// BlockData represents a block in the blockchain
// Includes consensus metadata for BFT
// (ConsensusData is defined in ue_consensus.go)
type BlockData struct {
	Height       uint64
	Hash         string
	Timestamp    time.Time
	Proposer     string
	Transactions []Transaction
	Consensus    ConsensusData
}
