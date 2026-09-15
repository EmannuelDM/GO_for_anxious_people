# Functions in Go (05_a_functions)

This project demonstrates the core concepts, syntax, mechanics, and advanced patterns of **functions** in Go.

In Go, functions are primary building blocks. They support multiple return values, named returns, variadic parameters, and are **first-class citizens**—meaning they can be assigned to variables, passed as arguments, returned from other functions, and form stateful closures.

---

## Key Concepts

### 1. Function Declarations & Pass-by-Value
* **Syntax**:
  ```go
  func FunctionName(param1 Type, param2 Type) ReturnType {
      return value
  }
  ```
* **Type Collapsing**: Consecutive parameters sharing the same type can omit the type for all but the last parameter:
  ```go
  func add(a, b int) int { return a + b }
  ```
* **Strictly Pass-by-Value**: Go passes all function arguments by value (a copy is created). Modifying a primitive parameter inside a function does **not** affect the caller's variable.
* **Mutating via Pointers**: To modify the caller's variable, pass a pointer (`*T`):
  ```go
  func modify(ptr *int) { *ptr = 999 }
  ```

---

### 2. Multiple Return Values & Error Handling
Unlike languages with `try/catch` exception hierarchies, Go functions idiomatically return a pair: `(result, error)` or `(result, ok)`.

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("cannot divide by zero")
    }
    return a / b, nil
}
```

* Callers must explicitly check errors:
  ```go
  res, err := divide(10, 2)
  if err != nil {
      log.Fatal(err)
  }
  ```
* Use the blank identifier `_` to discard unwanted return values.

---

### 3. Named Return Values & Naked Returns
Return values can be named in the function signature. They act as local variables initialized to their respective zero values.

```go
func calculateStats(numbers []int) (min int, max int, avg float64) {
    // min and max default to 0, avg defaults to 0.0
    // ...
    return // Naked return: returns the current values of min, max, avg
}
```

> [!TIP]
> **Best Practice**: Use naked returns only in short, simple functions. In longer functions, naked returns harm code readability and maintainability.

#### Modifying Named Returns with `defer`
A deferred function executes *after* the `return` statement assigns values to named returns, but *before* returning control to the caller. This allows `defer` closures to inspect or modify return values:

```go
func modifyReturn() (result int) {
    defer func() {
        result *= 2 // Alters the returned value to 10
    }()
    return 5
}
```

---

### 4. Variadic Functions (`...T`)
A variadic function accepts zero or more arguments of a specified type.

* Inside the function body, the variadic parameter behaves as a slice (`[]T`).
* **Unpacking a slice**: Pass an existing slice to a variadic parameter using the spread operator `...`:

```go
func sum(nums ...int) int {
    total := 0
    for _, n := range nums { total += n }
    return total
}

// Calls:
sum(1, 2, 3)          // Passing individual arguments
items := []int{4, 5}
sum(items...)         // Unpacking a slice
```

---

### 5. First-Class & Higher-Order Functions
In Go, functions are first-class types:

1. **Assigning to variables**:
   ```go
   fn := func(x int) int { return x * 2 }
   ```
2. **Custom Function Types**:
   ```go
   type Transformer func(int) int
   ```
3. **Higher-Order Functions**: Functions that accept functions as arguments or return them:
   ```go
   func mapSlice(nums []int, fn Transformer) []int { ... }
   func createMultiplier(factor int) func(int) int { ... }
   ```

---

### 6. Closures & Lexical Scoping
A **closure** is an anonymous function that references variables from outside its body. The function retains access to these referenced variables even after the outer function has finished executing.

```go
func makeCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}
```

Each call to `makeCounter()` captures and maintains its own private, isolated copy of `count`.

---

### 7. Defer Mechanics

The `defer` statement defers the execution of a function until the surrounding function returns.

1. **LIFO Stack (Last-In, First-Out)**: Multiple deferred calls execute in reverse order of declaration.
2. **Immediate Argument Evaluation**: Arguments to a deferred function are evaluated **immediately** when the `defer` statement is encountered, not when the function is executed.
   ```go
   x := 10
   defer fmt.Println(x) // Prints 10 even if x changes later!
   x = 20
   ```
3. **Resource Cleanup Idiom**: Standard Go pattern for closing files, mutex locks, database connections, and HTTP response bodies:
   ```go
   file, err := os.Open("data.txt")
   if err != nil { return err }
   defer file.Close()
   ```

---

### 8. Panic and Recover
* **`panic`**: Stops normal execution of the current goroutine, runs all deferred functions, and crashes the program with a stack trace.
* **`recover()`**: Built-in function that regains control of a panicking goroutine. It is **only effective when called directly within a deferred function**.

```go
func safeExecute(action func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
        }
    }()
    action()
    return nil
}
```

---

### 9. Recursion
Go fully supports recursive functions. Go goroutines have dynamically resizing stacks that start very small (~2 KB) and grow as needed, making standard recursion efficient for typical depth limits.

```go
func factorial(n uint) uint {
    if n <= 1 { return 1 }
    return n * factorial(n-1)
}
```

---

## Running the Program

Navigate to this directory (`05_a_functions`) and run:

```bash
go run main.go
```
