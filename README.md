# Underground Empire (UE)

A fully decentralized Layer 1 blockchain designed for high scalability, security, and long-term adaptability.

## Features

- **Advanced PoS Consensus**: Stake-weighted validator selection with 28,846 UE minimum stake requirement
- **BFT Finalization**: 67% threshold for block finalization ensuring immutability
- **Smart Contract Support**: High-performance virtualized environment for DApp execution
- **Modular Architecture**: Scalable and upgradable protocol design

## Quick Start

### Prerequisites

- Go 1.20 or higher
- Git

### Build and Run

```bash
# Clone the repository
git clone <repository-url>
cd UndergroundEmpire/Core

# Build the binary
make build

# Run the daemon
./build/ued start
```

## Usage

```bash
# Show help
./build/ued --help

# Check version
./build/ued version

# Start the node
./build/ued start
```

## Development

### Project Structure

```
Core/
├── cmd/ued/               # Command line interface
├── node/                  # Main node logic
├── modules/               # Core blockchain modules
├── core/                  # Fundamental components
├── protocol/              # Protocol layer
├── network/               # Network layer
└── storage/               # Data management
```

### Build Commands

```bash
make build        # Build the binary
make test         # Run tests
make clean        # Clean build artifacts
make help         # Show all commands
```

## License

This project is licensed under the MIT License. 