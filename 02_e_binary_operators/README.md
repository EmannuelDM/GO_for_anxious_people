# Go Binary Operators (02_e_binary_operators)

This project demonstrates the usage of **binary operators** in Go, including arithmetic, comparison, logical, and bitwise operators, as well as Go's specific rules for operator precedence.

---

## Key Concepts

### 1. Arithmetic Operators
Arithmetic operators operate on numeric values:
- `+` (addition)
- `-` (subtraction)
- `*` (multiplication)
- `/` (division)
- `%` (modulo/remainder)

> [!NOTE]
> Integer division in Go truncates towards zero (e.g., `10 / 3` is `3`). To get a fractional result, at least one of the operands must be a floating-point number (e.g., `10.0 / 3.0`).
> Modulo `%` is only defined for integer types.

### 2. Comparison Operators
Comparison operators compare two operands of the **same type** and yield a `bool`:
- `==` (equal)
- `!=` (not equal)
- `<` (less than)
- `<=` (less than or equal to)
- `>` (greater than)
- `>=` (greater than or equal to)

### 3. Logical Operators
Logical operators operate only on boolean values and support **short-circuit evaluation**:
- `&&` (conditional AND): Evaluates to `true` only if both operands are `true`. If the left operand is `false`, the right operand is not evaluated.
- `||` (conditional OR): Evaluates to `true` if at least one operand is `true`. If the left operand is `true`, the right operand is not evaluated.

### 4. Bitwise Operators
Bitwise operators operate on the individual bits of integer values:
- `&` (bitwise AND)
- `|` (bitwise OR)
- `^` (bitwise XOR)
- `&^` (bit clear / AND NOT): Clears bits in the first operand where the corresponding bit in the second operand is `1`.
- `<<` (left shift): Shifts bits to the left, filling empty spaces with `0` (equivalent to multiplying by $2^n$).
- `>>` (right shift): Shifts bits to the right (equivalent to integer division by $2^n$).

### 5. Operator Precedence
Go simplifies operator precedence down to **5 distinct levels** for binary operators (ordered from highest priority to lowest):

| Precedence | Operators | Category |
| :--- | :--- | :--- |
| **5** (Highest) | `*`, `/`, `%`, `<<`, `>>`, `&`, `&^` | Multiplicative / Shift / Bitwise AND / Bit Clear |
| **4** | `+`, `-`, `\|`, `^` | Additive / Bitwise OR / Bitwise XOR |
| **3** | `==`, `!=`, `<`, `<=`, `>`, `>=` | Comparison |
| **2** | `&&` | Logical AND |
| **1** (Lowest) | `\|\|` | Logical OR |

> [!WARNING]
> **Gotcha for C/C++ Developers:** 
> - In Go, bitwise shift (`<<`, `>>`) and bitwise AND (`&`) have **higher precedence** than addition/subtraction (`+`, `-`) and comparisons (`==`, `!=`). 
> - In C/C++, shifts have lower precedence than addition, and bitwise AND has lower precedence than equality.
> - **Example 1:** `1 << 2 + 1` in Go is evaluated as `(1 << 2) + 1` (evaluates to `5`). In C/C++, it is `1 << (2 + 1)` (evaluates to `8`).
> - **Example 2:** `x & y == 2` in Go is evaluated as `(x & y) == 2`. In C/C++, it is `x & (y == 2)`.

---

## Running the Program

Navigate to this directory (`02_e_binary_operators`) and run:

```bash
go run main.go
```
