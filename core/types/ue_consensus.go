package types

import (
	"time"
)

// VoteType represents the type of consensus vote
// (pre-vote or pre-commit)
type VoteType string

const (
	VoteTypePreVote   VoteType = "pre_vote"
	VoteTypePreCommit VoteType = "pre_commit"
)

// Vote represents a validator's vote in the consensus process
type Vote struct {
	ValidatorID string
	BlockHash   string
	Timestamp   time.Time
	Type        VoteType
}

// ConsensusData represents consensus-related data for a block
type ConsensusData struct {
	PreVotes     []Vote
	PreCommits   []Vote
	Finalized    bool
	FinalityTime time.Time
}

// ConsensusState represents the current consensus state
type ConsensusState struct {
	CurrentHeight    uint64
	CurrentBlockHash string
	Validators       []string // Validator IDs
	ConsensusRound   uint64
	Votes            []Vote
}

// FinalityData tracks finalized blocks and votes
// (for BFT threshold logic)
type FinalityData struct {
	BlockHeight   uint64
	BlockHash     string
	Finalized     bool
	FinalityVotes []Vote
	FinalityTime  time.Time
}
