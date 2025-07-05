package validator

import (
	"fmt"
	"time"

	"undergroundempire/core/types"
)

// ValidatorRegistry manages validator registration and operations
type ValidatorRegistry interface {
	RegisterNode(ctx types.Context, node ValidatorNode) error
	DeregisterNode(ctx types.Context, nodeID string) error
	GetActiveValidators(ctx types.Context) []ValidatorNode
	GetValidator(ctx types.Context, nodeID string) (ValidatorNode, error)
	UpdateValidator(ctx types.Context, node ValidatorNode) error
}

// ValidatorRewardEngine manages validator rewards and penalties
type ValidatorRewardEngine interface {
	CalculateRewards(ctx types.Context, nodeID string) uint64
	DistributeRewards(ctx types.Context, nodeID string, amount uint64) error
	SlashNode(ctx types.Context, nodeID string, reason SlashReason) error
	GetSlashHistory(ctx types.Context, nodeID string) []SlashRecord
}

// ValidatorNode represents a validator in the Underground Empire network
type ValidatorNode struct {
	ID          string
	Address     types.Address
	StakeAmount uint64
	Status      ValidatorStatus
	Commission  uint64 // Commission rate in basis points (0-10000)
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Description string
	Website     string
}

// ValidatorStatus represents the status of a validator
type ValidatorStatus string

const (
	ValidatorStatusActive   ValidatorStatus = "active"
	ValidatorStatusInactive ValidatorStatus = "inactive"
	ValidatorStatusSlashed  ValidatorStatus = "slashed"
	ValidatorStatusJailed   ValidatorStatus = "jailed"
)

// SlashReason represents the reason for slashing a validator
type SlashReason string

const (
	SlashReasonDoubleSigning SlashReason = "double_signing"
	SlashReasonDowntime      SlashReason = "downtime"
	SlashReasonInvalidBlock  SlashReason = "invalid_block"
	SlashReasonEquivocation  SlashReason = "equivocation"
)

// SlashRecord represents a slash record for a validator
type SlashRecord struct {
	ValidatorID string
	Reason      SlashReason
	Amount      uint64
	Timestamp   time.Time
	Height      uint64
}

// ValidatorManager implements validator management operations
type ValidatorManager struct {
	// TODO: Add storage interface in future commits
	validators map[string]ValidatorNode
}

// NewValidatorManager creates a new validator manager
func NewValidatorManager() *ValidatorManager {
	return &ValidatorManager{
		validators: make(map[string]ValidatorNode),
	}
}

// RegisterNode registers a new validator node
func (vm *ValidatorManager) RegisterNode(ctx types.Context, node ValidatorNode) error {
	// Validate minimum stake requirement
	if !types.IsValidatorEligible(node.StakeAmount) {
		return fmt.Errorf("insufficient stake: minimum required is %d UE, got %d",
			types.MinValidatorStake, node.StakeAmount)
	}

	// Check if validator already exists
	if _, exists := vm.validators[node.ID]; exists {
		return fmt.Errorf("validator with ID %s already exists", node.ID)
	}

	// Set timestamps
	node.CreatedAt = time.Now()
	node.UpdatedAt = time.Now()

	// Set initial status
	node.Status = ValidatorStatusActive

	// Store validator
	vm.validators[node.ID] = node

	return nil
}

// DeregisterNode deregisters a validator node
func (vm *ValidatorManager) DeregisterNode(ctx types.Context, nodeID string) error {
	validator, exists := vm.validators[nodeID]
	if !exists {
		return fmt.Errorf("validator with ID %s not found", nodeID)
	}

	// Update status
	validator.Status = ValidatorStatusInactive
	validator.UpdatedAt = time.Now()

	vm.validators[nodeID] = validator
	return nil
}

// GetActiveValidators returns all active validators
func (vm *ValidatorManager) GetActiveValidators(ctx types.Context) []ValidatorNode {
	var activeValidators []ValidatorNode

	for _, validator := range vm.validators {
		if validator.Status == ValidatorStatusActive {
			activeValidators = append(activeValidators, validator)
		}
	}

	return activeValidators
}

// GetValidator returns a specific validator
func (vm *ValidatorManager) GetValidator(ctx types.Context, nodeID string) (ValidatorNode, error) {
	validator, exists := vm.validators[nodeID]
	if !exists {
		return ValidatorNode{}, fmt.Errorf("validator with ID %s not found", nodeID)
	}

	return validator, nil
}

// UpdateValidator updates a validator's information
func (vm *ValidatorManager) UpdateValidator(ctx types.Context, node ValidatorNode) error {
	if _, exists := vm.validators[node.ID]; !exists {
		return fmt.Errorf("validator with ID %s not found", node.ID)
	}

	node.UpdatedAt = time.Now()
	vm.validators[node.ID] = node

	return nil
}

// CalculateRewards calculates rewards for a validator
func (vm *ValidatorManager) CalculateRewards(ctx types.Context, nodeID string) uint64 {
	validator, err := vm.GetValidator(ctx, nodeID)
	if err != nil {
		return 0
	}

	// Get total stake from all active validators
	activeValidators := vm.GetActiveValidators(ctx)
	totalStake := uint64(0)

	for _, v := range activeValidators {
		totalStake += v.StakeAmount
	}

	// Calculate reward based on stake proportion
	// TODO: Implement actual reward calculation logic in future commits
	blockReward := uint64(1000) // Placeholder value

	return types.CalculateValidatorReward(validator.StakeAmount, totalStake, blockReward)
}

// SlashNode slashes a validator for misbehavior
func (vm *ValidatorManager) SlashNode(ctx types.Context, nodeID string, reason SlashReason) error {
	validator, err := vm.GetValidator(ctx, nodeID)
	if err != nil {
		return err
	}

	// Calculate slash amount based on reason
	slashAmount := vm.calculateSlashAmount(validator.StakeAmount, reason)

	// Update validator status and stake
	validator.Status = ValidatorStatusSlashed
	validator.StakeAmount -= slashAmount
	validator.UpdatedAt = time.Now()

	// Ensure minimum stake is maintained
	if validator.StakeAmount < types.MinValidatorStake {
		validator.StakeAmount = 0
	}

	vm.validators[nodeID] = validator

	return nil
}

// calculateSlashAmount calculates the amount to slash based on the reason
func (vm *ValidatorManager) calculateSlashAmount(stakeAmount uint64, reason SlashReason) uint64 {
	switch reason {
	case SlashReasonDoubleSigning:
		return stakeAmount / 2 // 50% slash
	case SlashReasonDowntime:
		return stakeAmount / 10 // 10% slash
	case SlashReasonInvalidBlock:
		return stakeAmount / 4 // 25% slash
	case SlashReasonEquivocation:
		return stakeAmount / 2 // 50% slash
	default:
		return stakeAmount / 10 // Default 10% slash
	}
}

// GetTotalStake returns the total stake of all active validators
func (vm *ValidatorManager) GetTotalStake(ctx types.Context) uint64 {
	activeValidators := vm.GetActiveValidators(ctx)
	totalStake := uint64(0)

	for _, validator := range activeValidators {
		totalStake += validator.StakeAmount
	}

	return totalStake
}

// GetValidatorCount returns the number of active validators
func (vm *ValidatorManager) GetValidatorCount(ctx types.Context) int {
	return len(vm.GetActiveValidators(ctx))
}
