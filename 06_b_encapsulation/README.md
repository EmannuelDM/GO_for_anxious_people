# OOP Encapsulation in Go (06_b_encapsulation)

In Object-Oriented Programming (OOP), **encapsulation** (also known as *information hiding*) is the practice of bundling data with the methods that operate on that data, while restricting direct access to the internal state from the outside world.

Go takes a distinctly pragmatic, elegant approach to encapsulation:
- **No Access Specifier Keywords**: Go does not have `private`, `public`, `protected`, or `friend`.
- **Package-Level Boundaries**: The unit of encapsulation in Go is the **package**, not the struct or class.
- **Capitalization Rule**: An identifier's visibility is determined solely by the case of its first letter:
  - **Exported (Public)**: Starts with a Capital letter (e.g., `Account`, `Deposit`, `Balance`). Accessible from outside the package.
  - **Unexported (Private)**: Starts with a lowercase letter (e.g., `balanceCents`, `owner`, `sign`). Accessible *only* within the package where it is declared.

---

## The 3 Core Benefits of Encapsulation

```
┌─────────────────────────────────────────────────────────────┐
│                      Client Package                         │
│  (Can only see Exported types, constructors, and methods)   │
└──────────────────────────────┬──────────────────────────────┘
                               │ Calls NewAccount(), Deposit(), Balance()
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                      package bank                           │
│                                                             │
│   ┌──────────────────────────────────────────────────────┐  │
│   │                 type Account struct                  │  │
│   │   - id           (unexported)                        │  │
│   │   - owner        (unexported)                        │  │
│   │   - balanceCents (unexported - hidden representation)│  │
│   │   - transactions (unexported - defensive copy)       │  │
│   └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 1. Invariant Protection (Guaranteeing Valid State)
An **invariant** is a condition that must always remain true for a data structure to be valid (e.g., an account balance must never be negative, an email must have an `@` symbol, a mutex must protect a map).

When fields are unexported:
- Clients cannot construct arbitrary or corrupt instances using struct literals like `bank.Account{balanceCents: -9999}`.
- Clients must use an exported **Constructor function** (`NewAccount(...)`) which validates invariants before returning the object.
- Clients must use exported methods (`Deposit()`, `Withdraw()`) which enforce domain rules before mutating state.

```go
// Inside package bank:
func (a *Account) Withdraw(amount float64, note string) error {
    if amount <= 0 {
        return ErrNegativeAmount
    }
    cents := int64(math.Round(amount * 100))
    if cents > a.balanceCents {
        return ErrInsufficientFund // Invariant protected!
    }
    a.balanceCents -= cents
    return nil
}
```

### 2. Implementation Hiding & Freedom to Refactor
Clients program to the public method contracts, not the storage layout. The package maintainer can completely change the internal data representation without breaking external callers:

- **Currency representation**: Storing money as `float64` leads to cumulative IEEE 754 precision errors (e.g., `0.1 + 0.2 = 0.30000000000000004`). In `bank.Account`, money is stored internally as integer cents (`balanceCents int64`), but exposed to callers as `Balance() float64`.
- **Abstract Data Types**: In `bitset.IntSet`, callers interact with a set of integers (`Add(42)`, `Has(42)`), unaware that the internal implementation is a compact array of 64-bit unsigned integers (`[]uint64`) manipulated via bit shifts. The author can swap this for a hash map or tree without breaking callers.

### 3. Preventing Unintended Mutation & Read-Only Views
By not exporting a field, you can provide **read-only access** by exposing a getter method without a corresponding setter.

---

## Idiomatic Go Getters and Setters

In Go, getter methods **do not** use the `Get` prefix:

| Language | Idiomatic Getter | Idiomatic Setter |
| :--- | :--- | :--- |
| **Java / C#** | `account.getBalance()` | `account.setBalance(value)` |
| **Go** | `account.Balance()` | `account.SetBalance(value)` |

```go
// IDIOMATIC GO GETTER:
func (a *Account) Balance() float64 {
    return float64(a.balanceCents) / 100.0
}

// NON-IDIOMATIC (AVOID):
func (a *Account) GetBalance() float64 { ... }
```

> [!TIP]
> In Go, favor domain-driven verbs like `Deposit()` and `Withdraw()` over generic setters like `SetBalance()`. Behavior-rich methods express intent and enforce validation far better than raw setters.

---

## Defensive Copying (The Slice & Pointer Leak Trap)

A subtle and dangerous encapsulation pitfall in Go occurs with reference types (slices, maps, and pointers):

```go
// VULNERABLE (LEAKY GETTER):
func (a *Account) LeakTransactions() []Transaction {
    return a.transactions // Returns a slice header pointing to internal array!
}
```

Because slices share the underlying backing array, an external caller can mutate elements directly:
```go
leaked := acc.LeakTransactions()
leaked[0].Amount = 999999.00 // Modifies account's actual ledger!
```

### The Solution: Defensive Copying
Always allocate a fresh copy before returning internal slices or maps:

```go
// SAFE (DEFENSIVE COPY):
func (a *Account) Transactions() []Transaction {
    copied := make([]Transaction, len(a.transactions))
    copy(copied, a.transactions)
    return copied
}
```

Now, any mutation made by the caller affects only their isolated copy; the internal state remains pristine.

---

## Encapsulating Concurrency (Thread-Safe Structures)

Encapsulation is vital for concurrent programming. Exposing raw synchronization primitives like `sync.Mutex` or `sync.RWMutex` invites race conditions and deadlocks:
- Callers might forget to lock before reading/writing.
- Callers might forget to unlock, deadlocking the application.
- Callers might acquire locks in conflicting orders.

By encapsulating both the data and the mutex inside unexported fields, thread-safety is guaranteed by construction:

```go
type Cache struct {
    mu    sync.RWMutex      // Unexported: callers cannot manipulate lock directly
    items map[string]string // Unexported: callers cannot bypass lock
}

func (c *Cache) Set(key, val string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = val
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    val, ok := c.items[key]
    return val, ok
}
```

---

## Interface-Based Encapsulation (Opaque Types)

The strongest degree of encapsulation in Go is achieved by exporting an **interface** while keeping the **concrete struct completely unexported**:

```go
// 1. Exported interface defines behavior
type Manager interface {
    Generate(subject string, ttl time.Duration) string
    Validate(tokenStr string) (string, bool)
}

// 2. Concrete struct is UNEXPORTED
type hmacManager struct {
    secretKey []byte
}

// 3. Factory function returns the interface
func NewManager(secret string) Manager {
    return &hmacManager{secretKey: []byte(secret)}
}
```

### Why this is powerful:
1. **Total Opacity**: External packages cannot even write `var m hmacManager`—it does not compile.
2. **Cannot Inspect Fields**: Secret keys, salt values, and internal algorithms are completely hidden from callers.
3. **Pluggable Implementations**: The package can introduce RSA, ECDSA, or Mock implementations without altering the client code.

---

## Go Encapsulation vs. Traditional OOP

| Feature | Java / C++ / C# | Go |
| :--- | :--- | :--- |
| **Keywords** | `public`, `private`, `protected` | None (capitalization only) |
| **Boundary** | Class-level (or file/assembly) | **Package-level** |
| **Same-package access** | `private` blocks even same package | All code in the same package can access unexported fields |
| **Getter naming** | `getX()` | `X()` |
| **Invariant enforcement** | Constructors (`new Class()`) | Factory functions (`NewType(...)`) |
| **Opaque types** | Abstract classes / Interfaces | Exported Interface + Unexported Struct |

---

## Project Structure

```
06_b_encapsulation/
├── go.mod               # Module definition (06_b_encapsulation)
├── README.md            # Conceptual guide & tutorial documentation
├── main.go              # Interactive runner demonstrating all encapsulation pillars
├── bank/
│   └── account.go       # Invariant enforcement, constructor, getters, defensive copying
├── cache/
│   └── cache.go         # Concurrency encapsulation (RWMutex + Map)
├── bitset/
│   └── intset.go        # Abstract Data Type (bitmask encapsulated behind set API)
└── token/
    └── token.go         # Interface-based encapsulation (opaque HMAC token manager)
```

---

## Running the Code

Execute the demo module using standard Go CLI:

```bash
cd 06_b_encapsulation
go run main.go
```

To run static analysis and lint checks:

```bash
go vet ./...
```
