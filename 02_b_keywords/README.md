# Go Keywords & Predeclared Identifiers (02_b_keywords)

This project demonstrates the reserved keywords, predeclared constants, types, and functions in Go.

## Reserved Keywords (25)
These words are reserved and cannot be used as identifiers (e.g., variable names, function names, etc.):
- **Declaration / Scope**: `package`, `import`, `const`, `var`, `type`, `func`
- **Control Flow**: `if`, `else`, `switch`, `case`, `default`, `fallthrough`, `for`, `range`, `break`, `continue`, `goto`, `return`
- **Concurrency**: `go`, `select`, `chan`
- **Types / Interfaces**: `struct`, `map`, `interface`
- **Special**: `defer`

## Predeclared Identifiers
These are not reserved keywords (meaning they can technically be shadowed, though it's highly discouraged):
- **Constants**: `true`, `false`, `iota`, `nil`
- **Types**: Numeric (`int`, `uint8`, etc.), string, bool, error, etc.
- **Built-in Functions**: `make`, `new`, `len`, `cap`, `append`, `copy`, `delete`, `complex`, `real`, `imag`, `panic`, `recover`, etc.

## Running the Program

Navigate to this directory (`02_b_keywords`) and run:

```bash
go run main.go
```
