# Go Pointers (02_c_pointers)

This project demonstrates the concept and usage of **pointers** in Go.

Pointers allow you to store and manipulate the memory addresses of variables, enabling direct modification of values and sharing data efficiently without copying.

## Key Concepts

### 1. Pointer Declaration & Zero Value
A pointer type is written as `*T` (e.g., `*int` is a pointer to an integer). The zero value of any pointer type is `nil` (it does not point to any memory location).
```go
var p *int // Declares a nil pointer to an int
```

### 2. Operators
- **Address-of (`&`)**: Returns the memory address of a variable.
- **Dereference (`*`)**: Accesses or modifies the value at the address held by the pointer.
```go
x := 42
p = &x   // p now holds the address of x
val := *p // dereferences p to get 42
*p = 100 // modifies the value of x to 100
```

### 3. Modifying Arguments in Functions
Go is a **pass-by-value** language, which means arguments are copied when passed to functions.
- If you pass a variable by value, the function gets a copy and cannot modify the original.
- If you pass a pointer, the function receives a copy of the pointer (which points to the same memory address), allowing it to modify the original variable.

```go
func modifyValue(v int)   { v = 999 }  // Does not change original
func modifyPointer(v *int) { *v = 999 } // Changes original
```

### 4. Struct Pointers & Syntactic Sugar
In Go, if you have a pointer to a struct, you can access its fields directly using `structPtr.Field` instead of writing `(*structPtr).Field`. Go automatically dereferences it for you.
```go
type Person struct { Name string }
p := Person{Name: "Alice"}
ptr := &p
ptr.Name = "Bob" // Syntactic sugar for (*ptr).Name = "Bob"
```

### 5. Memory Allocation with `new()`
The built-in `new(T)` function allocates zeroed storage for a new item of type `T` and returns its address (`*T`).
```go
p := new(int) // p is a *int pointing to a value of 0
```

### 6. Pointer Arithmetic Safety
Unlike C and C++, **Go does not allow pointer arithmetic** (like `p++` or `p + 5`). This is a safety feature designed to prevent memory corruption and bugs.

---

## Running the Program

Navigate to this directory (`02_c_pointers`) and run:

```bash
go run main.go
```
