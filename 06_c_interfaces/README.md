# Interfaces and Interface Types in Go (06_c_interfaces)

In Go, an **interface** is an abstract type that defines a contract of behavior: a set of method signatures. Unlike traditional object-oriented languages (such as Java, C++, or C#), Go does **not** have an `implements` keyword. Instead, Go uses **structural typing** (often called *compile-time duck typing*):

> *"If something can do this, then it can be used here."*

If a concrete type declares all methods required by an interface, it implements that interface **implicitly and automatically**.

---

## Architecture of Go Interfaces

```
┌────────────────────────────────────────────────────────────────────────┐
│                        Interface Value (16 bytes)                      │
│                                                                        │
│   ┌──────────────────────────────┐   ┌──────────────────────────────┐  │
│   │   tab (*itab)                │   │   data (unsafe.Pointer)      │  │
│   │   - Interface Type Info      │   │   - Points to underlying     │  │
│   │   - Concrete Type Descriptor │   │     concrete value or copy   │  │
│   │   - Method Dispatch Table    │   │                              │  │
│   └──────────────┬───────────────┘   └──────────────┬───────────────┘  │
└──────────────────┼──────────────────────────────────┼──────────────────┘
                   │                                  │
                   ▼                                  ▼
      ┌─────────────────────────┐        ┌─────────────────────────┐
      │  Type: *Circle          │        │  Circle{Radius: 7.5}    │
      │  Methods: [Area, Name]  │        └─────────────────────────┘
      └─────────────────────────┘
```

An interface variable is a **two-word fat pointer** (16 bytes on 64-bit systems):
1. **`tab` (`*itab`)**: Points to the interface table containing metadata about the interface, the concrete dynamic type, and the function pointers for dynamic dispatch.
2. **`data` (`unsafe.Pointer`)**: Points to the dynamic value (the concrete data).

---

## Table of Contents

1. [Implicit Implementation & Polymorphism](#1-implicit-implementation--polymorphism)
2. [Pointer vs. Value Receivers & Method Sets](#2-pointer-vs-value-receivers--method-sets)
3. [Interface Internal Representation (iface & eface)](#3-interface-internal-representation-iface--eface)
4. [The Infamous "Nil Interface" Trap](#4-the-infamous-nil-interface-trap)
5. [The Empty Interface (any / interface{})](#5-the-empty-interface-any--interface)
6. [Type Assertions & Capability Querying](#6-type-assertions--capability-querying)
7. [Type Switches](#7-type-switches)
8. [Interface Composition & Embedding](#8-interface-composition--embedding)
9. [Interface Types: Method Sets vs. Type Sets (Go 1.18+ Generics)](#9-interface-types-method-sets-vs-type-sets-go-118-generics)
10. [Idiomatic Interface Patterns & Best Practices](#10-idiomatic-interface-patterns--best-practices)

---

## 1. Implicit Implementation & Polymorphism

In Go, interfaces are satisfied **implicitly**:

```go
type Shape interface {
    Name() string
    Area() float64
    Perimeter() float64
}

type Rectangle struct {
    Width, Height float64
}

// Rectangle implicitly satisfies Shape:
func (r Rectangle) Name() string      { return "Rectangle" }
func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }
```

### Why Implicit Satisfaction is Powerful
- **Zero Coupling**: The package declaring `Rectangle` does not need to know about the package declaring `Shape`.
- **Retroactive Satisfaction**: You can declare an interface for third-party library types without modifying their source code.
- **Natural Decoupling**: You can define small, focused interfaces at the call site where they are consumed.

---

## 2. Pointer vs. Value Receivers & Method Sets

One of the most common compilation errors in Go stems from **method sets**:

| Type | Methods in Method Set |
| :--- | :--- |
| **Value `T`** | Only methods with value receiver `(t T)` |
| **Pointer `*T`** | Both pointer `(*t T)` **AND** value `(t T)` receivers |

### The Addressability Rule

```go
type Incrementer interface {
    Increment()
    Value() int
}

type Counter struct {
    count int
}

func (c *Counter) Increment() { c.count++ }
func (c Counter) Value() int  { return c.count }

// Valid:
var inc Incrementer = &Counter{} // OK: *Counter has Increment()

// INVALID (COMPILE ERROR):
// var inc Incrementer = Counter{}
// cannot use Counter{} as Incrementer: Counter does not implement Incrementer
// (method Increment has pointer receiver)
```

> [!IMPORTANT]
> **Why doesn't Go automatically take the address of `Counter{}`?**
> When a concrete value is stored in an interface, Go stores a **copy** of that value. The value stored inside an interface is **not addressable** (`&inc.count` is illegal). If Go allowed calling a pointer-receiver method on an interface containing a value copy, any mutation would modify only the copy and be discarded, causing subtle bugs.

---

## 3. Interface Internal Representation (iface & eface)

Under the hood in Go's runtime:

1. **Non-Empty Interface (`iface`)**:
   Contains `tab *itab` (interface table and concrete type descriptor) and `data unsafe.Pointer`.
2. **Empty Interface (`eface` / `any`)**:
   Contains `_type *_type` (concrete type descriptor) and `data unsafe.Pointer`.

### The Zero Value of an Interface

An interface variable is initialized to its zero value:
```go
var s Shape // s is nil: tab == nil AND data == nil
```

```go
fmt.Printf("Dynamic Type: %T, Dynamic Value: %v, Is nil: %t\n", s, s, s == nil)
// Output: Dynamic Type: <nil>, Dynamic Value: <nil>, Is nil: true
```

---

## 4. The Infamous "Nil Interface" Trap

> [!CAUTION]
> **The #1 bug in Go production code:**
> An interface is `nil` **if and only if** both its dynamic type and dynamic value are `nil`.
> Storing a `nil` pointer of a concrete type in an interface makes the interface **non-nil**!

```go
type CustomAppError struct {
    Code    int
    Message string
}

func (e *CustomAppError) Error() string {
    return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

// ❌ DANGEROUS / BUGGY FUNCTION:
func findUserBad(id string) error {
    var err *CustomAppError = nil // typed pointer initialized to nil
    if id == "" {
        err = &CustomAppError{Code: 400, Message: "empty ID"}
        return err
    }
    return err // BUG: returns (*CustomAppError)(nil) wrapped in error interface!
}

// Calling findUserBad:
err := findUserBad("valid-id")
if err != nil {
    // ❌ THIS BRANCH EXECUTES!
    // err is an interface holding: (type = *CustomAppError, value = nil)
    // Because type != nil, (err != nil) evaluates to TRUE!
}
```

### The Idiomatic Fix: Return Literal `nil`

```go
// ✅ IDIOMATIC FUNCTION:
func findUserGood(id string) error {
    if id == "" {
        return &CustomAppError{Code: 400, Message: "empty ID"}
    }
    return nil // Explicit untyped nil: interface is (type = nil, value = nil)
}
```

---

## 5. The Empty Interface (`any` / `interface{}`)

Since `interface{}` declares zero methods, **every type in Go satisfies it**. In Go 1.18+, `any` was introduced as a standard predeclared alias:

```go
type any = interface{}
```

### Slice Conversion Pitfall

You **cannot** pass `[]int` or `[]string` directly to a function expecting `[]any`:

```go
nums := []int{1, 2, 3}
// printSlice(nums) // COMPILE ERROR: cannot use nums (type []int) as []any
```

> [!NOTE]
> **Why?**
> A slice `[]int` is a contiguous block of integers (8 bytes each). A slice `[]any` is a contiguous block of interface values (16 bytes each). Converting them requires allocating a new slice and converting each element:
> ```go
> anySlice := make([]any, len(nums))
> for i, v := range nums {
>     anySlice[i] = v
> }
> ```

---

## 6. Type Assertions & Capability Querying

A **Type Assertion** extracts the underlying concrete value from an interface:

```go
var val any = "Learning Go Interfaces"

// 1. Unsafe assertion (PANICS if types don't match):
s := val.(string)

// 2. Safe "comma-ok" assertion:
if s, ok := val.(string); ok {
    fmt.Println("Value is a string:", s)
}
```

### Capability Querying (Interface-to-Interface)

You can use type assertions to dynamically check if an interface value implements *another* interface:

```go
var w io.Writer = getWriter()

// Capability Query: Does w also implement Flusher?
if flusher, ok := w.(Flusher); ok {
    flusher.Flush() // Feature detected and used safely!
}
```

This pattern is widely used in the standard library (e.g., `net/http`, `io.Copy`).

---

## 7. Type Switches

A **Type Switch** compares the dynamic type of an interface against multiple candidate types:

```go
switch v := x.(type) {
case nil:
    fmt.Println("Interface is nil")
case int:
    fmt.Println("Integer squared:", v*v) // v is typed as int
case string:
    fmt.Println("Uppercase:", strings.ToUpper(v)) // v is typed as string
case Shape:
    fmt.Println("Shape area:", v.Area()) // v is typed as Shape
default:
    fmt.Printf("Unknown type: %T\n", v)
}
```

---

## 8. Interface Composition & Embedding

Larger interfaces are built by **composing** smaller, single-method interfaces:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}

// Composed interfaces:
type ReadWriter interface {
    Reader
    Writer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}
```

> [!TIP]
> **Rob Pike's Proverb**:
> *"The bigger the interface, the weaker the abstraction."*
> Keep interfaces small (1 to 2 methods). Single-method interfaces like `io.Reader`, `io.Writer`, `fmt.Stringer`, and `error` are the foundation of Go's composability.

---

## 9. Interface Types: Method Sets vs. Type Sets (Go 1.18+ Generics)

Starting in Go 1.18, interfaces were generalized to define **type sets**:

| Interface Category | Definition | Allowed Usage |
| :--- | :--- | :--- |
| **Basic Interface** | Contains **only methods** (e.g. `io.Reader`, `Shape`) | Anywhere: variable types, parameters, return types, struct fields |
| **General Interface** | Contains **type elements** (e.g. `int \| float64`, `~string`) | **ONLY as type constraints** in generic type parameters |

```go
// General Interface (Type Set constraint)
type Number interface {
    ~int | ~int64 | ~float32 | ~float64
}

// Generic function constrained by Number:
func Sum[T Number](values ...T) T {
    var total T
    for _, v := range values {
        total += v
    }
    return total
}
```

### The Approximation Operator (`~`)

`~int` matches `int` as well as **any defined type whose underlying type is `int`**:

```go
type Celsius float64

// Because Number includes ~float64, Celsius satisfies Number!
temps := []Celsius{21.5, 23.0}
total := Sum(temps...) // Compiles cleanly!
```

---

## 10. Idiomatic Interface Patterns & Best Practices

### Pattern 1: Compile-Time Interface Verification

Ensure a struct satisfies an interface at compile-time using the blank identifier:

```go
var _ service.UserStore = (*service.MemoryStore)(nil)
```
If `*MemoryStore` fails to implement any method in `UserStore`, the compiler immediately flags this line.

### Pattern 2: Client-Defined Interfaces ("Accept Interfaces, Return Structs")

- **Java/C# approach**: The provider declares `IUserRepository`, and services depend on the provider package.
- **Go approach**: The **consumer** (`package service`) declares the `UserStore` interface with only the methods it needs. The provider package returns concrete structs (`*MemoryStore`, `*PostgresStore`).

```go
// In consumer package (service):
type UserStore interface {
    FindByID(id string) (*User, error)
    Save(u *User) error
}

func NewUserService(store UserStore) *UserService {
    return &UserService{store: store}
}
```

### Pattern 3: Effortless Mocking & Testing

Because interfaces are consumer-defined and structurally satisfied, you can write mock implementations for unit testing without third-party mocking libraries:

```go
type MockStore struct {
    FindByIDFn func(id string) (*User, error)
    SaveFn     func(user *User) error
}

func (m *MockStore) FindByID(id string) (*User, error) {
    return m.FindByIDFn(id)
}
// ...
```

---

## Running the Code

Execute the complete tutorial and interactive demonstrations:

```bash
cd 06_c_interfaces
go run main.go
```

Run the unit tests:

```bash
go test -v ./...
```
