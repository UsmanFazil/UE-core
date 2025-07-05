package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// Address represents an Underground Empire address
type Address [20]byte

// CoinAmount represents an amount of UE coins
type CoinAmount struct {
	Amount uint64
	Denom  string
}

// Transaction represents a transaction in the Underground Empire network
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

// NewAddress creates a new address from a hex string
func NewAddress(hexStr string) (Address, error) {
	var addr Address

	// Remove 0x prefix if present
	hexStr = strings.TrimPrefix(hexStr, "0x")

	// Decode hex string
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return addr, fmt.Errorf("invalid hex string: %v", err)
	}

	// Check length
	if len(bytes) != 20 {
		return addr, fmt.Errorf("invalid address length: expected %d, got %d", 20, len(bytes))
	}

	// Copy bytes to address
	copy(addr[:], bytes)
	return addr, nil
}

// String returns the hex string representation of the address
func (a Address) String() string {
	return "0x" + hex.EncodeToString(a[:])
}

// Bytes returns the byte representation of the address
func (b Address) Bytes() []byte {
	return b[:]
}

// IsZero checks if the address is the zero address
func (a Address) IsZero() bool {
	for _, b := range a {
		if b != 0 {
			return false
		}
	}
	return true
}

// NewCoinAmount creates a new coin amount
func NewCoinAmount(amount uint64, denom string) CoinAmount {
	return CoinAmount{
		Amount: amount,
		Denom:  denom,
	}
}

// NewUECoins creates a new UE coin amount
func NewUECoins(amount uint64) CoinAmount {
	return NewCoinAmount(amount, "ue")
}

// String returns the string representation of the coin amount
func (c CoinAmount) String() string {
	return fmt.Sprintf("%d%s", c.Amount, c.Denom)
}

// Add adds two coin amounts (must be same denomination)
func (c CoinAmount) Add(other CoinAmount) (CoinAmount, error) {
	if c.Denom != other.Denom {
		return CoinAmount{}, fmt.Errorf("cannot add different denominations: %s and %s", c.Denom, other.Denom)
	}

	return CoinAmount{
		Amount: c.Amount + other.Amount,
		Denom:  c.Denom,
	}, nil
}

// Sub subtracts two coin amounts (must be same denomination)
func (c CoinAmount) Sub(other CoinAmount) (CoinAmount, error) {
	if c.Denom != other.Denom {
		return CoinAmount{}, fmt.Errorf("cannot subtract different denominations: %s and %s", c.Denom, other.Denom)
	}

	if c.Amount < other.Amount {
		return CoinAmount{}, fmt.Errorf("insufficient balance: %s < %s", c.String(), other.String())
	}

	return CoinAmount{
		Amount: c.Amount - other.Amount,
		Denom:  c.Denom,
	}, nil
}

// Mul multiplies the coin amount by a factor
func (c CoinAmount) Mul(factor uint64) CoinAmount {
	return CoinAmount{
		Amount: c.Amount * factor,
		Denom:  c.Denom,
	}
}

// Div divides the coin amount by a factor
func (c CoinAmount) Div(factor uint64) CoinAmount {
	if factor == 0 {
		return CoinAmount{Amount: 0, Denom: c.Denom}
	}
	return CoinAmount{
		Amount: c.Amount / factor,
		Denom:  c.Denom,
	}
}

// IsZero checks if the coin amount is zero
func (c CoinAmount) IsZero() bool {
	return c.Amount == 0
}

// IsPositive checks if the coin amount is positive
func (c CoinAmount) IsPositive() bool {
	return c.Amount > 0
}

// NewTransaction creates a new transaction
func NewTransaction(from, to Address, amount CoinAmount, gas, gasPrice uint64, data []byte, nonce uint64) Transaction {
	return Transaction{
		From:     from,
		To:       to,
		Amount:   amount,
		Gas:      gas,
		GasPrice: gasPrice,
		Data:     data,
		Nonce:    nonce,
	}
}

// CalculateHash calculates the transaction hash
func (tx Transaction) CalculateHash() string {
	// Create a string representation for hashing
	data := fmt.Sprintf("%s%s%s%d%d%d%d",
		tx.From.String(),
		tx.To.String(),
		tx.Amount.String(),
		tx.Gas,
		tx.GasPrice,
		tx.Nonce,
		tx.Timestamp)

	// Add data if present
	if len(tx.Data) > 0 {
		data += hex.EncodeToString(tx.Data)
	}

	// Calculate hash
	hash := sha256.Sum256([]byte(data))
	return "0x" + hex.EncodeToString(hash[:])
}

// Validate validates the transaction
func (tx Transaction) Validate() error {
	// Check addresses
	if tx.From.IsZero() {
		return fmt.Errorf("from address cannot be zero")
	}

	if tx.To.IsZero() {
		return fmt.Errorf("to address cannot be zero")
	}

	// Check amount
	if tx.Amount.IsZero() {
		return fmt.Errorf("amount cannot be zero")
	}

	// Check gas
	if tx.Gas == 0 {
		return fmt.Errorf("gas cannot be zero")
	}

	// Check gas price
	if tx.GasPrice == 0 {
		return fmt.Errorf("gas price cannot be zero")
	}

	return nil
}

// CalculateGasCost calculates the total gas cost
func (tx Transaction) CalculateGasCost() uint64 {
	return tx.Gas * tx.GasPrice
}

// ParseCoinAmount parses a coin amount from string
func ParseCoinAmount(s string) (CoinAmount, error) {
	// Find the denomination (last non-digit character)
	var denom string
	var amountStr string

	for i := len(s) - 1; i >= 0; i-- {
		if s[i] >= '0' && s[i] <= '9' {
			amountStr = s[:i+1]
			denom = s[i+1:]
			break
		}
	}

	if amountStr == "" {
		return CoinAmount{}, fmt.Errorf("invalid coin amount format: %s", s)
	}

	amount, err := strconv.ParseUint(amountStr, 10, 64)
	if err != nil {
		return CoinAmount{}, fmt.Errorf("invalid amount: %v", err)
	}

	return CoinAmount{
		Amount: amount,
		Denom:  denom,
	}, nil
}
