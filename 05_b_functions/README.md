# Errors, EOF, Anonymous & Variadic Functions (05_b_functions)

This module explores four essential, interrelated topics in Go function design and control flow:
1. **Idiomatic Error Handling** (Sentinel errors, wrapping with `%w`, `errors.Is`, `errors.As`, and `errors.Join`)
2. **Handling EOF (`io.EOF`)** (The `io.Reader` contract, chunked streaming, `bufio.Reader`, and `bufio.Scanner`)
3. **Anonymous Functions & Closures** (IIFEs, callbacks, stateful closures, defer telemetry, and goroutines)
4. **Variadic Functions (`...T`) & The Functional Options Pattern** (Slice unpacking, `...any`, and idiomatic API design)

---

## 1. Idiomatic Error Handling

In Go, errors are normal values that implement the standard `error` interface:

```go
type error interface {
    Error() string
}
```

### Sentinel Errors
Sentinel errors are package-level exported variables representing specific, anticipated failure scenarios. By convention, they are prefixed with `Err`:

```go
var (
    ErrNotFound         = errors.New("record not found")
    ErrPermissionDenied = errors.New("permission denied")
    ErrInvalidConfig    = errors.New("invalid configuration")
)
```

### Error Wrapping with `%w` and `errors.Is`
When bubbling errors up the call stack, add context while preserving the underlying error using `fmt.Errorf` with the `%w` verb:

```go
func fetchUserData(id int) (string, error) {
    data, err := simulateDatabaseQuery(id)
    if err != nil {
        // %w preserves the error chain for unwrapping
        return "", fmt.Errorf("fetchUserData(id=%d): %w", id, err)
    }
    return data, nil
}
```

Use `errors.Is` to check if any error in the wrap chain matches the target sentinel:

```go
if errors.Is(err, ErrNotFound) {
    // True even if err is deeply wrapped!
}
```

### Custom Typed Errors and `errors.As`
When you need structured metadata attached to an error, implement the `error` interface on a custom struct:

```go
type ValidationError struct {
    Field   string
    Value   any
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on '%s': %s", e.Field, e.Message)
}
```

Use `errors.As` to inspect and extract the structured error from an error chain:

```go
var valErr *ValidationError
if errors.As(err, &valErr) {
    fmt.Println("Faulty field:", valErr.Field)
}
```

### Combining Errors with `errors.Join` (Go 1.20+)
To aggregate multiple independent validation failures:

```go
err := errors.Join(err1, err2, err3)
// If all inputs are nil, errors.Join returns nil.
```

---

## 2. Working with EOF (`io.EOF`)

`io.EOF` is defined in the standard `io` package:
```go
var EOF = errors.New("EOF")
```

> [!IMPORTANT]
> **EOF is NOT a failure!**
> `io.EOF` is a sentinel error signaling the graceful, expected end of an input stream.

### The Critical `io.Reader` Contract
According to Go's `io.Reader` specification:
> "An `io.Reader` may return `n > 0` bytes alongside an `err == io.EOF` on the final read."

**Always process `n` bytes before checking `err == io.EOF`:**

```go
for {
    n, err := reader.Read(buf)
    if n > 0 {
        // STEP 1: Process valid bytes FIRST
        processData(buf[:n])
    }
    if err != nil {
        if errors.Is(err, io.EOF) {
            // STEP 2: Handle graceful end of stream
            break
        }
        return err // Real unexpected error
    }
}
```

### Stream Reading Comparison

| Approach | Mechanics | Handling of `io.EOF` |
| :--- | :--- | :--- |
| **`io.Reader.Read()`** | Low-level chunk buffer | Returns `io.EOF`; must check `n > 0` first. |
| **`bufio.Reader.ReadString()`** | Reads up to delimiter (e.g. `\n`) | Can return final line *and* `io.EOF` simultaneously. |
| **`bufio.Scanner`** | High-level token reader | Automatically consumes `io.EOF`; returns `false` from `Scan()`, while `scanner.Err()` returns `nil`. |

---

## 3. Anonymous Functions & Closures

An **anonymous function** (function literal) is a function defined without a name identifier.

### 1. Immediately Invoked Function Expression (IIFE)
Useful for scoping localized setup logic:

```go
version := func(major, minor int) string {
    return fmt.Sprintf("v%d.%d", major, minor)
}(1, 4)
```

### 2. Variable Assignment & Callbacks
Functions are first-class values in Go:

```go
words := []string{"elephant", "cat", "dog"}
sort.Slice(words, func(i, j int) bool {
    return len(words[i]) < len(words[j])
})
```

### 3. Stateful Closures (Lexical Scoping)
A closure captures and references variables from its enclosing scope:

```go
func makeAccumulator(initial int) func(int) int {
    sum := initial
    return func(val int) int {
        sum += val // Retains reference to 'sum'
        return sum
    }
}
```

### 4. Closures in `defer` for Telemetry
```go
defer func(start time.Time) {
    fmt.Printf("Duration: %v\n", time.Since(start))
}(time.Now())
```

---

## 4. Variadic Functions (`...T`)

A variadic function accepts zero, one, or more arguments for its final parameter:

```go
func calculateMetrics(category string, numbers ...int) (min int, max int, sum int, count int) {
    // Inside the function, 'numbers' has type []int
    count = len(numbers)
    // ...
}
```

### Invocation Modes
1. **Discrete arguments**:
   ```go
   calculateMetrics("CPU", 10, 20, 30)
   ```
2. **Zero arguments**:
   ```go
   calculateMetrics("Idle") // numbers is an empty slice (len = 0)
   ```
3. **Unpacking a slice (`...`)**:
   ```go
   data := []int{5, 15, 25}
   calculateMetrics("Batch", data...) // Spread operator
   ```

### Variadic with `any` (`...any`)
Used to accept arbitrary values of any type (similar to `fmt.Println`):

```go
func logValues(prefix string, values ...any) {
    for _, v := range values {
        switch val := v.(type) {
        case string:
            // handle string
        case error:
            // handle error
        }
    }
}
```

---

## 5. Real-World Synergy: The Functional Options Pattern

The **Functional Options Pattern** is a widely adopted Go design pattern that combines **variadic functions**, **closures / anonymous functions**, and **error handling** to create robust, extensible APIs.

```go
type ServerConfig struct {
    Host      string
    Port      int
    Timeout   time.Duration
    EnableTLS bool
}

type ServerOption func(*ServerConfig) error

func WithPort(port int) ServerOption {
    return func(s *ServerConfig) error {
        if port < 1 || port > 65535 {
            return fmt.Errorf("invalid port %d: %w", port, ErrInvalidConfig)
        }
        s.Port = port
        return nil
    }
}

func NewServerConfig(host string, opts ...ServerOption) (*ServerConfig, error) {
    cfg := &ServerConfig{
        Host:    host,
        Port:    8080, // Sensible defaults
        Timeout: 30 * time.Second,
    }

    for _, opt := range opts {
        if err := opt(cfg); err != nil {
            return nil, err
        }
    }
    return cfg, nil
}
```

**Usage:**
```go
// Clean, readable, and self-documenting:
server, err := NewServerConfig("api.service.internal",
    WithPort(8443),
    WithTLS(true),
)
```

---

## Running the Code

Navigate to this directory and run:

```bash
go run main.go
```
