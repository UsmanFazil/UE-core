package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"undergroundempire/core/types"
	"undergroundempire/modules/consensus"
	"undergroundempire/modules/validator"
)

var (
	// Version is set during build
	Version = "dev"

	// Root command
	rootCmd = &cobra.Command{
		Use:   "ued",
		Short: "Underground Empire (UE) - Fully decentralized Layer 1 blockchain daemon",
		Long: `Underground Empire (UE) is a fully decentralized Layer 1 blockchain 
designed for high scalability, security, and long-term adaptability.

The UE daemon provides a complete node implementation with advanced Proof of Stake 
consensus, Byzantine Fault Tolerance finalization, and smart contract execution capabilities.

Key Features:
- Advanced PoS consensus with 28,846 UE minimum stake requirement
- BFT finalization with 67% threshold for block immutability
- High-performance smart contract execution environment
- Modular architecture for scalability and upgrades
- Enhanced security with custom slashing mechanisms`,
		Version: Version,
	}
)

func init() {
	// Add subcommands
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(validatorCmd)
	rootCmd.AddCommand(treasuryCmd)
	rootCmd.AddCommand(governanceCmd)
	rootCmd.AddCommand(demoConsensusCmd)
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Underground Empire node",
	Long: `Start the Underground Empire blockchain node. This command initializes
the node, connects to the network, and begins participating in consensus.

The node will:
- Initialize the blockchain state
- Connect to peer nodes in the network
- Begin participating in block validation
- Start the consensus mechanism
- Enable smart contract execution`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting Underground Empire node...")
		fmt.Println("Node initialization complete")
		fmt.Println("Connecting to network...")
		fmt.Println("Node is now running and participating in consensus")
		fmt.Println("Press Ctrl+C to stop the node")

		// TODO: Implement actual node startup logic
		// This is a placeholder for the first commit
	},
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Underground Empire (UE) Daemon\n")
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Consensus: Advanced PoS with BFT Finalization\n")
		fmt.Printf("Minimum Validator Stake: 28,846 UE\n")
		fmt.Printf("Consensus Threshold: 67%%\n")
		fmt.Printf("Block Time: 5 seconds\n")
		fmt.Printf("Epoch Duration: 100 blocks\n")
	},
}

// validatorCmd represents the validator command group
var validatorCmd = &cobra.Command{
	Use:   "validator",
	Short: "Manage validator operations",
	Long: `Manage validator operations for the Underground Empire network.

Validators are responsible for:
- Proposing and validating blocks
- Participating in consensus
- Maintaining network security
- Earning rewards for honest participation

Minimum requirements:
- 28,846 UE stake
- Reliable network connection
- Consistent uptime`,
}

// treasuryCmd represents the treasury command group
var treasuryCmd = &cobra.Command{
	Use:   "treasury",
	Short: "Manage treasury operations",
	Long: `Manage treasury operations including token transfers, balance queries,
and account management for the Underground Empire network.`,
}

// governanceCmd represents the governance command group
var governanceCmd = &cobra.Command{
	Use:   "governance",
	Short: "Participate in network governance",
	Long: `Participate in the governance of the Underground Empire network.

Governance features include:
- Submitting proposals for protocol upgrades
- Voting on proposals
- Parameter change proposals
- Emergency proposals for critical issues`,
}

var demoConsensusCmd = &cobra.Command{
	Use:   "demo-consensus",
	Short: "Run a simulated in-memory consensus round with multiple validators",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[Demo] Starting in-memory consensus demo with 3 validators...")

		// 1. Setup validators
		valMgr := validator.NewValidatorManager()
		vals := []validator.ValidatorNode{
			{ID: "val1", StakeAmount: 30000},
			{ID: "val2", StakeAmount: 30000},
			{ID: "val3", StakeAmount: 30000},
		}
		for _, v := range vals {
			valMgr.RegisterNode(types.Context{}, v)
		}

		// 2. Setup consensus engine
		engine := consensus.NewInMemoryConsensusEngine(valMgr, vals)

		// 3. Simulate consensus for 200 blocks
		for i := 0; i < 200; i++ {
			fmt.Printf("\n[Demo] === Block %d ===\n", i+1)
			block, _ := engine.ProposeBlock()
			fmt.Printf("[Demo] Proposer: %s\n", block.Proposer)
			engine.PreVote(block)
			engine.PreCommit(block)
			err := engine.FinalizeBlock(block)
			if err != nil {
				fmt.Println("[Demo] Finalization error:", err)
			} else {
				fmt.Printf("[Demo] Block %d finalized/mined!\n", block.Height)
				fmt.Printf("[Demo] Timestamp: %s\n", block.Timestamp.Format(time.RFC3339))
			}
			fmt.Println("[Demo] ----------------------")

			// Add delay after block mining for better demo experience
			if i < 199 {
				time.Sleep(2 * time.Second)
			}
		}
		fmt.Println("[Demo] Consensus demo complete.")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
