package types

import (
	"context"
	"time"
)

// Context represents the Underground Empire context
type Context struct {
	context.Context
	Height    uint64
	Timestamp time.Time
	ChainID   string
}

// UEError represents Underground Empire specific errors
type UEError struct {
	Code    string
	Message string
}

func (e UEError) Error() string {
	return e.Message
}

// UE-specific constants
const (
	// Validator Requirements
	MinValidatorStake = 28846 // UE coins required to become validator

	// Consensus Parameters
	ConsensusThreshold = 67 // Percentage for block finalization

	// Network Parameters
	BlockTime     = 5   // seconds per block
	EpochDuration = 100 // blocks per epoch

	// Gas Parameters
	DefaultGasLimit = 200000
	DefaultGasPrice = 1000000000 // 1 gwei in wei

	// Chain Parameters
	DefaultChainID = "underground-empire-1"

	// Address Parameters
	AddressLength = 20 // bytes
	HashLength    = 32 // bytes
)

// NewContext creates a new UE context
func NewContext(ctx context.Context, height uint64, timestamp time.Time, chainID string) Context {
	return Context{
		Context:   ctx,
		Height:    height,
		Timestamp: timestamp,
		ChainID:   chainID,
	}
}

// WithHeight returns a new context with updated height
func (c Context) WithHeight(height uint64) Context {
	return Context{
		Context:   c.Context,
		Height:    height,
		Timestamp: c.Timestamp,
		ChainID:   c.ChainID,
	}
}

// WithTimestamp returns a new context with updated timestamp
func (c Context) WithTimestamp(timestamp time.Time) Context {
	return Context{
		Context:   c.Context,
		Height:    c.Height,
		Timestamp: timestamp,
		ChainID:   c.ChainID,
	}
}

// IsValidatorEligible checks if a stake amount meets validator requirements
func IsValidatorEligible(stakeAmount uint64) bool {
	return stakeAmount >= MinValidatorStake
}

// CalculateValidatorReward calculates reward for a validator based on stake
func CalculateValidatorReward(stakeAmount uint64, totalStake uint64, blockReward uint64) uint64 {
	if totalStake == 0 {
		return 0
	}
	return (stakeAmount * blockReward) / totalStake
}

// IsConsensusReached checks if consensus threshold is met
func IsConsensusReached(votes uint64, totalValidators uint64) bool {
	if totalValidators == 0 {
		return false
	}
	percentage := (votes * 100) / totalValidators
	return percentage >= ConsensusThreshold
}

// CalculateEpochNumber calculates the current epoch number
func CalculateEpochNumber(blockHeight uint64) uint64 {
	return blockHeight / EpochDuration
}

// IsEpochBoundary checks if the current block is at epoch boundary
func IsEpochBoundary(blockHeight uint64) bool {
	return blockHeight%EpochDuration == 0
}

// CalculateNextEpochHeight calculates the height of the next epoch
func CalculateNextEpochHeight(currentHeight uint64) uint64 {
	currentEpoch := CalculateEpochNumber(currentHeight)
	return (currentEpoch + 1) * EpochDuration
}
