package main

import (
	"fmt"
	"sync"
	"time"

	"06_b_encapsulation/bank"
	"06_b_encapsulation/bitset"
	"06_b_encapsulation/cache"
	"06_b_encapsulation/token"
)

func main() {
	fmt.Println("##################################################################")
	fmt.Println("             GO OOP: ENCAPSULATION & DATA HIDING                  ")
	fmt.Println("##################################################################")

	demoPackageVisibilityAndInvariants()
	demoDefensiveCopying()
	demoConcurrencyEncapsulation()
	demoAbstractDataTypeEncapsulation()
	demoInterfaceBasedEncapsulation()

	fmt.Println("\nAll encapsulation demonstrations completed successfully!")
}

// ============================================================================
// 1. PACKAGE VISIBILITY & INVARIANT ENFORCEMENT
// ============================================================================
// In Go, encapsulation is strictly PACKAGE-LEVEL, governed by capitalization:
// - Uppercase: Exported (Public to other packages)
// - Lowercase: Unexported (Private to the package)
//
// Benefits:
// 1. Invariants cannot be bypassed (e.g. balance cannot be set negative).
// 2. Struct literal instantiation can be blocked, forcing use of constructors.
// 3. Implementation details (cents vs dollars) can change freely without
//    breaking external API callers.
func demoPackageVisibilityAndInvariants() {
	fmt.Println("\n==================================================================")
	fmt.Println("1. PACKAGE VISIBILITY & INVARIANT ENFORCEMENT (Bank Account)")
	fmt.Println("==================================================================")

	// Step 1: Trying to create invalid accounts via constructor
	fmt.Println("--- Invariant Validation at Construction ---")
	_, err := bank.NewAccount("", 100.0)
	if err != nil {
		fmt.Printf("Rejected empty owner name: %v\n", err)
	}

	_, err = bank.NewAccount("Bob", -50.0)
	if err != nil {
		fmt.Printf("Rejected negative initial deposit: %v\n", err)
	}

	// Step 2: Creating a valid account
	acc, err := bank.NewAccount("Alice Smith", 250.00)
	if err != nil {
		fmt.Printf("Unexpected error: %v\n", err)
		return
	}

	// Idiomatic Go getters omit the "Get" prefix: acc.Balance(), NOT acc.GetBalance()
	fmt.Printf("\nCreated Account: ID=%s, Owner=%s, Balance=$%.2f\n",
		acc.ID(), acc.Owner(), acc.Balance())

	// Step 3: Demonstrating what callers CANNOT do:
	// The lines below would cause COMPILE ERRORS if uncommented:
	//
	// acc.balanceCents = 9999999 // Error: acc.balanceCents undefined (cannot refer to unexported field)
	// acc.owner = "Mallory"     // Error: acc.owner undefined (cannot refer to unexported field)
	//
	// This proves callers CANNOT bypass business rules or corrupt internal state!

	// Step 4: Testing domain invariants via methods
	fmt.Println("\n--- Invariant Enforcement via Methods ---")

	// Invariant: deposits must be positive
	err = acc.Deposit(-20.0, "Bad deposit")
	fmt.Printf("Deposit -$20.00 result: %v\n", err)

	// Invariant: withdrawals cannot overdraft the account
	err = acc.Withdraw(500.0, "Rent payment")
	fmt.Printf("Withdraw $500.00 (balance $250.00) result: %v\n", err)

	// Valid transactions
	_ = acc.Deposit(100.00, "Freelance bonus")
	_ = acc.Withdraw(75.50, "Grocery shopping")
	fmt.Printf("Balance after +$100.00 and -$75.50: $%.2f\n", acc.Balance())

	// Updating owner through validated setter
	_ = acc.SetOwner("Alice Smith-Johnson")
	fmt.Printf("Updated Owner: %s\n", acc.Owner())
}

// ============================================================================
// 2. DEFENSIVE COPYING (THE SLICE & POINTER REFERENCE TRAP)
// ============================================================================
// A subtle pitfall in Go: slices and pointers reference underlying memory.
// If a getter method directly returns an internal slice, the caller can
// modify slice elements, silently breaking encapsulation!
//
// The remedy: return a DEFENSIVE COPY.
func demoDefensiveCopying() {
	fmt.Println("\n==================================================================")
	fmt.Println("2. DEFENSIVE COPYING: PREVENTING THE SLICE LEAK TRAP")
	fmt.Println("==================================================================")

	acc, _ := bank.NewAccount("Charlie Brown", 100.00)
	_ = acc.Deposit(50.00, "Birthday gift")

	// Case A: Defensive Copying (Safe)
	fmt.Println("--- Case A: Safe Getter with Defensive Copy (acc.Transactions()) ---")
	safeHistory := acc.Transactions()
	fmt.Printf("Before tampering - First transaction note: '%s'\n", safeHistory[0].Note)

	// External caller tampers with the slice they received:
	safeHistory[0].Note = "TAMPERED: Fake transaction memo!"
	safeHistory[0].Amount = 999999.00

	// Check if the bank account's internal state was affected:
	freshCopy := acc.Transactions()
	fmt.Printf("After caller tampering - Account's actual record: '%s', Amount: $%.2f\n",
		freshCopy[0].Note, freshCopy[0].Amount)
	fmt.Println("-> Encapsulation PRESERVED! Internal ledger remained untouched.")

	// Case B: Leaky Getter (Unsafe demonstration)
	fmt.Println("\n--- Case B: Leaky Getter without Copy (acc.LeakTransactions()) ---")
	leakedHistory := acc.LeakTransactions()
	fmt.Printf("Before tampering - First transaction note: '%s'\n", leakedHistory[0].Note)

	// External caller tampers with the leaked slice:
	leakedHistory[0].Note = "CORRUPTED: Attacker modified internal ledger!"
	leakedHistory[0].Amount = 1337.00

	// Check if the internal state was corrupted:
	corruptedView := acc.LeakTransactions()
	fmt.Printf("After caller tampering - Account's actual record: '%s', Amount: $%.2f\n",
		corruptedView[0].Note, corruptedView[0].Amount)
	fmt.Println("-> Encapsulation BROKEN! Leaked internal slice allowed mutation.")
}

// ============================================================================
// 3. CONCURRENCY SAFETY THROUGH ENCAPSULATION
// ============================================================================
// Encapsulating synchronization primitives (sync.RWMutex) alongside data
// structures (maps) guarantees thread-safety. Callers cannot forget to lock
// or unlock, preventing data races and deadlocks.
func demoConcurrencyEncapsulation() {
	fmt.Println("\n==================================================================")
	fmt.Println("3. CONCURRENCY SAFETY THROUGH ENCAPSULATION (Cache)")
	fmt.Println("==================================================================")

	c := cache.New()
	var wg sync.WaitGroup

	numWriters := 5
	numReaders := 5

	fmt.Printf("Spawning %d writers and %d readers simultaneously...\n", numWriters, numReaders)

	// Concurrent writers
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				key := fmt.Sprintf("session-%d-%d", workerID, j)
				c.Set(key, fmt.Sprintf("val-%d", workerID*100+j))
			}
		}(i)
	}

	// Concurrent readers
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				key := fmt.Sprintf("session-%d-%d", workerID, j)
				c.Get(key)
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("All goroutines finished safely without data races!\n")
	fmt.Printf("Cache entry count: %d\n", c.Len())
	fmt.Printf("Snapshot sample keys: %v\n", c.Keys()[:3])
}

// ============================================================================
// 4. ABSTRACT DATA TYPE (ADT) ENCAPSULATION
// ============================================================================
// Encapsulation hides the physical representation of data.
// IntSet represents a mathematical set of integers. Internally, it is
// implemented using 64-bit word bitwise masks (words []uint64).
// The caller neither knows nor cares about bit shifting operations.
func demoAbstractDataTypeEncapsulation() {
	fmt.Println("\n==================================================================")
	fmt.Println("4. ABSTRACT DATA TYPE ENCAPSULATION (Bitset / IntSet)")
	fmt.Println("==================================================================")

	set := bitset.New()

	// Insert numbers into the set
	numbers := []int{1, 9, 42, 63, 64, 127, 255}
	for _, n := range numbers {
		set.Add(n)
	}

	fmt.Printf("IntSet contents (formatted): %s\n", set.String())
	fmt.Printf("IntSet element count:        %d\n", set.Len())
	fmt.Printf("Contains 42?                 %t\n", set.Has(42))
	fmt.Printf("Contains 100?                %t\n", set.Has(100))

	// Remove an element
	set.Remove(42)
	fmt.Printf("After set.Remove(42):        %s (Contains 42: %t)\n",
		set.String(), set.Has(42))
	fmt.Println("-> The caller uses intuitive Set APIs; bitwise mechanics are 100% hidden.")
}

// ============================================================================
// 5. INTERFACE-BASED ENCAPSULATION (OPAQUE IMPLEMENTATIONS)
// ============================================================================
// In Go, the strongest form of encapsulation exports an INTERFACE while keeping
// the CONCRETE STRUCT unexported.
//
// The factory function NewManager returns token.Manager (interface).
// Callers cannot inspect fields (secret key), cannot bypass verification,
// and cannot instantiate the struct without the constructor.
func demoInterfaceBasedEncapsulation() {
	fmt.Println("\n==================================================================")
	fmt.Println("5. INTERFACE-BASED ENCAPSULATION (Opaque Manager)")
	fmt.Println("==================================================================")

	// Constructor returns the exported interface 'token.Manager'.
	// The underlying struct 'token.hmacManager' is UNEXPORTED.
	mgr := token.NewManager("super-secret-cryptographic-salt-key-999")

	// Generate a token valid for 500 milliseconds
	tok := mgr.Generate("user_john_doe", 500*time.Millisecond)
	fmt.Printf("Generated signed token:\n  %s\n", tok)

	// Validate valid token
	subject, ok := mgr.Validate(tok)
	fmt.Printf("\nValidation with genuine token:\n  Subject: '%s', Valid: %t\n", subject, ok)

	// Validate tampered token
	tampered := tok + "hacked"
	subject, ok = mgr.Validate(tampered)
	fmt.Printf("\nValidation with tampered token:\n  Subject: '%s', Valid: %t\n", subject, ok)

	// Validate expired token
	fmt.Println("\nWaiting 600ms for token expiration...")
	time.Sleep(600 * time.Millisecond)
	subject, ok = mgr.Validate(tok)
	fmt.Printf("Validation with expired token:\n  Subject: '%s', Valid: %t\n", subject, ok)

	fmt.Println("-> Concrete implementation is completely invisible; only interface contract is exposed.")
}
