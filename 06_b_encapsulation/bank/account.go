package bank

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

// Standard domain errors returned by Account operations.
var (
	ErrInvalidOwner     = errors.New("bank: account owner cannot be empty")
	ErrNegativeAmount   = errors.New("bank: amount must be strictly greater than zero")
	ErrInsufficientFund = errors.New("bank: insufficient funds for withdrawal")
)

// Transaction represents a financial ledger entry.
// Its fields are exported so external packages can inspect statement details.
type Transaction struct {
	Type      string    // "DEPOSIT" or "WITHDRAWAL"
	Amount    float64   // Amount in currency units
	Timestamp time.Time // When the transaction occurred
	Note      string    // Description or memo
}

// Account represents a bank account.
//
// In Go, encapsulation is achieved through package-level visibility:
// - All fields below start with a lowercase letter, meaning they are UNEXPORTED.
// - Code outside package 'bank' cannot view, access, or directly mutate these fields.
// - Callers must use the exported constructor (NewAccount) and exported methods.
type Account struct {
	id           string        // Unique identifier (unexported)
	owner        string        // Account holder name (unexported)
	balanceCents int64         // Stored internally as integer cents to avoid float rounding errors! (unexported)
	transactions []Transaction // Internal ledger history (unexported)
}

// idCounter provides deterministic, sequential IDs for demonstration.
var idCounter int64 = 1000

// NewAccount is the CONSTRUCTOR function (Factory Pattern).
// Because Account's fields are unexported, external packages cannot use struct
// literals like `bank.Account{owner: "Alice"}`.
// The constructor guarantees that:
// 1. Invariants are checked BEFORE memory is allocated (owner not empty, deposit >= 0).
// 2. The struct is never left in an uninitialized or corrupt state.
func NewAccount(owner string, initialDeposit float64) (*Account, error) {
	trimmedOwner := strings.TrimSpace(owner)
	if trimmedOwner == "" {
		return nil, ErrInvalidOwner
	}
	if initialDeposit < 0 {
		return nil, ErrNegativeAmount
	}

	idCounter++
	accID := fmt.Sprintf("ACC-%d", idCounter)
	cents := int64(math.Round(initialDeposit * 100))

	acc := &Account{
		id:           accID,
		owner:        trimmedOwner,
		balanceCents: cents,
		transactions: make([]Transaction, 0, 8),
	}

	if initialDeposit > 0 {
		acc.transactions = append(acc.transactions, Transaction{
			Type:      "DEPOSIT",
			Amount:    initialDeposit,
			Timestamp: time.Now(),
			Note:      "Initial account opening deposit",
		})
	}

	return acc, nil
}

// ----------------------------------------------------------------------------
// IDIOMATIC GO GETTERS
// ----------------------------------------------------------------------------
// In Go, getter methods DO NOT use the "Get" prefix!
// Use 'Balance()' instead of 'GetBalance()'.
// Use 'Owner()' instead of 'GetOwner()'.
// Use 'ID()' instead of 'GetID()'.

// ID returns the unique account identifier.
func (a *Account) ID() string {
	return a.id
}

// Owner returns the account holder's name.
func (a *Account) Owner() string {
	return a.owner
}

// Balance returns the current balance in currency units (e.g. 10.50).
// Notice how the client doesn't need to know that we store 'balanceCents' internally!
// This is the core benefit of encapsulation: we can refactor internal representation
// without breaking caller code.
func (a *Account) Balance() float64 {
	return float64(a.balanceCents) / 100.0
}

// ----------------------------------------------------------------------------
// STATE MUTATORS & INVARIANT ENFORCEMENT
// ----------------------------------------------------------------------------

// SetOwner updates the account holder's name after validation.
func (a *Account) SetOwner(newOwner string) error {
	trimmed := strings.TrimSpace(newOwner)
	if trimmed == "" {
		return ErrInvalidOwner
	}
	a.owner = trimmed
	return nil
}

// Deposit credits the account with a specified amount.
// Enforces the invariant: deposit amount must be positive.
func (a *Account) Deposit(amount float64, note string) error {
	if amount <= 0 {
		return ErrNegativeAmount
	}

	cents := int64(math.Round(amount * 100))
	a.balanceCents += cents

	a.transactions = append(a.transactions, Transaction{
		Type:      "DEPOSIT",
		Amount:    amount,
		Timestamp: time.Now(),
		Note:      note,
	})

	return nil
}

// Withdraw debits the account with a specified amount.
// Enforces invariants:
// 1. Withdrawal amount must be positive.
// 2. Withdrawal cannot exceed current balance (prevents overdraft).
func (a *Account) Withdraw(amount float64, note string) error {
	if amount <= 0 {
		return ErrNegativeAmount
	}

	cents := int64(math.Round(amount * 100))
	if cents > a.balanceCents {
		return fmt.Errorf("%w: requested %.2f, available %.2f", ErrInsufficientFund, amount, a.Balance())
	}

	a.balanceCents -= cents

	a.transactions = append(a.transactions, Transaction{
		Type:      "WITHDRAWAL",
		Amount:    amount,
		Timestamp: time.Now(),
		Note:      note,
	})

	return nil
}

// ----------------------------------------------------------------------------
// DEFENSIVE COPYING (PREVENTING THE SLICE LEAK TRAP)
// ----------------------------------------------------------------------------
// A common vulnerability in Go: slices and pointers reference shared memory.
// If an unexported slice is returned directly, the caller can mutate internal
// state behind the struct's back!

// Transactions returns a DEFENSIVE COPY of the transaction history.
// Callers can sort, filter, or mutate the returned slice without corrupting the account's internal state.
func (a *Account) Transactions() []Transaction {
	copied := make([]Transaction, len(a.transactions))
	copy(copied, a.transactions)
	return copied
}

// LeakTransactions deliberately returns the internal slice directly.
// EDUCATIONAL DEMO ONLY: Demonstrates how returning references breaks encapsulation!
func (a *Account) LeakTransactions() []Transaction {
	return a.transactions
}
