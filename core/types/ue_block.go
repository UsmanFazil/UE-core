package types

import (
	"time"
)

// Transaction represents a transaction in the Underground Empire network
// (If already defined elsewhere, you can remove this duplicate)
type Transaction struct {
	Hash      string
	From      Address
	To        Address
	Amount    CoinAmount
	Gas       uint64
	GasPrice  uint64
	Data      []byte
	Nonce     uint64
	Signature []byte
	Timestamp int64
}

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
