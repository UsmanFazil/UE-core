package consensus

import (
	"fmt"
	"sync"
	"time"

	"undergroundempire/core/types"
	"undergroundempire/modules/validator"
)

// ConsensusEngine defines the interface for the consensus process
type ConsensusEngine interface {
	ProposeBlock() (*types.BlockData, error)
	PreVote(block *types.BlockData) error
	PreCommit(block *types.BlockData) error
	FinalizeBlock(block *types.BlockData) error
	GetState() *ConsensusState
}

// ConsensusState holds the current state of consensus
// (wraps types.ConsensusState for in-memory tracking)
type ConsensusState struct {
	CurrentHeight   uint64
	CurrentRound    uint64
	Validators      []validator.ValidatorNode
	ProposerIndex   int
	Votes           []types.Vote
	FinalizedBlocks []*types.BlockData
	Mutex           sync.Mutex
}

// InMemoryConsensusEngine is a simple, single-node consensus engine for demo/testing
// (no networking, no persistence)
type InMemoryConsensusEngine struct {
	state      *ConsensusState
	valManager *validator.ValidatorManager
}

// NewInMemoryConsensusEngine creates a new consensus engine
func NewInMemoryConsensusEngine(valManager *validator.ValidatorManager, initialValidators []validator.ValidatorNode) *InMemoryConsensusEngine {
	return &InMemoryConsensusEngine{
		state: &ConsensusState{
			CurrentHeight:   1,
			CurrentRound:    0,
			Validators:      initialValidators,
			ProposerIndex:   0,
			Votes:           []types.Vote{},
			FinalizedBlocks: []*types.BlockData{},
		},
		valManager: valManager,
	}
}

// ProposeBlock selects the next proposer (round-robin) and creates a new block
func (ce *InMemoryConsensusEngine) ProposeBlock() (*types.BlockData, error) {
	ce.state.Mutex.Lock()
	defer ce.state.Mutex.Unlock()

	if len(ce.state.Validators) == 0 {
		return nil, fmt.Errorf("no validators available")
	}
	proposer := ce.state.Validators[ce.state.ProposerIndex]
	defaultBlockTx := types.Transaction{
		Hash:      "tx1",
		From:      [20]byte{},
		To:        [20]byte{},
		Amount:    types.CoinAmount{Amount: 1, Denom: "ue"},
		Gas:       21000,
		GasPrice:  1,
		Data:      nil,
		Nonce:     1,
		Signature: nil,
		Timestamp: time.Now().Unix(),
	}
	block := &types.BlockData{
		Height:       ce.state.CurrentHeight,
		Hash:         fmt.Sprintf("block_%d", ce.state.CurrentHeight),
		Timestamp:    time.Now(),
		Proposer:     proposer.ID,
		Transactions: []types.Transaction{defaultBlockTx},
		Consensus:    types.ConsensusData{},
	}
	fmt.Printf("[Consensus] Proposer for block %d: %s\n", block.Height, proposer.ID)
	return block, nil
}

// PreVote simulates pre-vote phase for the block
func (ce *InMemoryConsensusEngine) PreVote(block *types.BlockData) error {
	ce.state.Mutex.Lock()
	defer ce.state.Mutex.Unlock()

	for _, v := range ce.state.Validators {
		vote := types.Vote{
			ValidatorID: v.ID,
			BlockHash:   block.Hash,
			Timestamp:   time.Now(),
			Type:        types.VoteTypePreVote,
		}
		ce.state.Votes = append(ce.state.Votes, vote)
		fmt.Printf("[Consensus] PreVote by %s for block %s\n", v.ID, block.Hash)
	}
	return nil
}

// PreCommit simulates pre-commit phase for the block
func (ce *InMemoryConsensusEngine) PreCommit(block *types.BlockData) error {
	ce.state.Mutex.Lock()
	defer ce.state.Mutex.Unlock()

	for _, v := range ce.state.Validators {
		vote := types.Vote{
			ValidatorID: v.ID,
			BlockHash:   block.Hash,
			Timestamp:   time.Now(),
			Type:        types.VoteTypePreCommit,
		}
		ce.state.Votes = append(ce.state.Votes, vote)
		fmt.Printf("[Consensus] PreCommit by %s for block %s\n", v.ID, block.Hash)
	}
	return nil
}

// FinalizeBlock finalizes the block if >=67% pre-commits
func (ce *InMemoryConsensusEngine) FinalizeBlock(block *types.BlockData) error {
	ce.state.Mutex.Lock()
	defer ce.state.Mutex.Unlock()

	totalValidators := len(ce.state.Validators)
	preCommits := 0
	for _, vote := range ce.state.Votes {
		if vote.BlockHash == block.Hash && vote.Type == types.VoteTypePreCommit {
			preCommits++
		}
	}
	if totalValidators == 0 {
		return fmt.Errorf("no validators available")
	}
	percentage := (preCommits * 100) / totalValidators
	if percentage >= int(types.ConsensusThreshold) {
		block.Consensus.Finalized = true
		block.Consensus.FinalityTime = time.Now()
		ce.state.FinalizedBlocks = append(ce.state.FinalizedBlocks, block)
		fmt.Printf("[Consensus] Block %d finalized with %d/%d pre-commits (>=67%%)\n", block.Height, preCommits, totalValidators)
		// Move to next height and proposer
		ce.state.CurrentHeight++
		ce.state.ProposerIndex = (ce.state.ProposerIndex + 1) % totalValidators
		ce.state.Votes = []types.Vote{}
		return nil
	}
	return fmt.Errorf("not enough pre-commits to finalize block: %d/%d", preCommits, totalValidators)
}

// GetState returns the current consensus state
func (ce *InMemoryConsensusEngine) GetState() *ConsensusState {
	return ce.state
}
